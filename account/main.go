package main

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/config"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/gateway/accountHandler"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/migrations"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/repository"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/service"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/account/utils"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("config.load.failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	postgres := utils.NewPostgresql(
		utils.WithHost(cfg.Postgresql.Host),
		utils.WithPort(cfg.Postgresql.Port),
		utils.WithUser(cfg.Postgresql.User),
		utils.WithPassword(cfg.Postgresql.Password),
		utils.WithName(cfg.Postgresql.Name),
		utils.WithMaxOpenConns(cfg.Postgresql.MaxOpenConns),
		utils.WithMaxIdleConns(cfg.Postgresql.MaxIdleConns),
		utils.WithMaxIdleTime(cfg.Postgresql.MaxIdleTime),
		utils.WithSSLMode(cfg.Postgresql.SSLMode),
		utils.WithTimeout(cfg.Postgresql.Timeout),
	)

	db, err := postgres.Connect()
	if err != nil {
		slog.Error("postgres.connect.failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			slog.Error("postgres.close.failed", slog.String("error", closeErr.Error()))
		}
	}()

	migrator, err := migrations.NewMigrator(db, postgres.Name)
	if err != nil {
		slog.Error("migrations.init.failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer func() {
		if closeErr := migrator.Close(); closeErr != nil {
			slog.Error("migrations.close.failed", slog.String("error", closeErr.Error()))
		}
	}()

	if err := migrator.Up(); err != nil {
		slog.Error("migrations.up.failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	accountRepository := repository.NewAccountRepository(db, db)
	accountService := service.NewAccountService(accountRepository)
	accountGRPCServer := accountHandler.NewGRPCServer(accountService)

	serverErrCh := make(chan error, 1)
	go func() {
		slog.Info("grpc.server.starting", slog.String("addr", cfg.Application.AccountPort))
		serverErrCh <- accountGRPCServer.Serve(cfg.Application.AccountPort)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigCh:
		slog.Info("shutdown.signal.received", slog.String("signal", sig.String()))
	case err := <-serverErrCh:
		slog.Error("grpc.server.failed", slog.String("error", err.Error()))
	}

	slog.Info("shutdown.initiating", slog.String("timeout", "30s"))

	if stopErr := accountGRPCServer.Stop(); stopErr != nil {
		slog.Error("grpc.stop.failed", slog.String("error", stopErr.Error()))
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	select {
	case <-shutdownCtx.Done():
		slog.Warn("shutdown.timeout.reached")
	case <-time.After(5 * time.Second):
		slog.Info("shutdown.drain.complete")
	}

	close(sigCh)
	slog.Info("shutdown.complete")
}
