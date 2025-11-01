package main

import (
	"context"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/config"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/gateway/catalogHandler"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/repository"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/service"
	"github.com/saleh-ghazimoradi/MircoEcoMarket/catalog/utils"
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

	elastic := utils.NewElasticSearch(
		utils.WithHost(cfg.ElasticSearch.Host),
		utils.WithPort(cfg.ElasticSearch.Port),
		utils.WithUsername(cfg.ElasticSearch.Username),
		utils.WithPassword(cfg.ElasticSearch.Password),
		utils.WithTimeout(cfg.ElasticSearch.Timeout),
	)

	client, err := elastic.Connect()
	if err != nil {
		slog.Error("elastic.connect.failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	catalogRepository := repository.NewCatalogRepository(client, "catalogs")
	catalogService := service.NewCatalogService(catalogRepository)
	catalogGRPCServer := catalogHandler.NewGRPCCatalogServer(catalogService)

	serverErrCh := make(chan error, 1)
	go func() {
		slog.Info("grpc.server.starting", slog.String("addr", cfg.Application.CatalogPort))
		serverErrCh <- catalogGRPCServer.Serve(cfg.Application.CatalogPort)
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

	if stopErr := catalogGRPCServer.Stop(); stopErr != nil {
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
