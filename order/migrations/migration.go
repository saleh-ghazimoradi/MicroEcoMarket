package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrationsFS embed.FS

type Migrator struct {
	db        *sql.DB
	dbName    string
	migration *migrate.Migrate
}

func (m *Migrator) Up() error {
	if err := m.migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}

func (m *Migrator) Rollback() error {
	if err := m.migration.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply rollback: %w", err)
	}

	fmt.Println("Rollback applied successfully")
	return nil
}

func (m *Migrator) Close() error {
	if m.migration != nil {
		source, driver := m.migration.Close()
		if source != nil {
			return fmt.Errorf("failed to close source: %w", source)
		}
		if driver != nil {
			return fmt.Errorf("failed to close driver: %w", driver)
		}
	}
	return nil
}

func NewMigrator(db *sql.DB, dbName string) (*Migrator, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{DatabaseName: dbName})

	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	source, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to load migration files: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, dbName, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrate: %w", err)
	}

	return &Migrator{
		db:        db,
		dbName:    dbName,
		migration: m,
	}, nil
}
