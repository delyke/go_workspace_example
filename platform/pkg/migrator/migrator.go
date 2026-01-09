package migrator

import (
	"context"
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"

	"github.com/delyke/go_workspace_example/platform/pkg/logger"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

type GooseLoggerAdapter struct{}

func (g *GooseLoggerAdapter) Write(p []byte) (n int, err error) {
	logger.Info(context.Background(), string(p))
	return len(p), nil
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up(ctx context.Context) error {
	logger.Info(ctx, "Migrator Up")
	goose.SetLogger(log.New(&GooseLoggerAdapter{}, "", 0))
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}
	return nil
}
