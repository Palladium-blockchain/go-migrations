package cli

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/Palladium-blockchain/go-migrations/internal/driver/postgres"
	"github.com/Palladium-blockchain/go-migrations/internal/migrator"
	"github.com/Palladium-blockchain/go-migrations/internal/source/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MigrateCommand struct{}

func NewMigrateCommand() *MigrateCommand {
	return &MigrateCommand{}
}

func (cmd *MigrateCommand) Execute(ctx context.Context, _ []string) error {
	// Config
	env, err := MigrateCommandLoadEnvConfig()
	if err != nil {
		fmt.Println("Config error:", err)
		return err
	}

	// Driver
	db, err := sql.Open("postgres", env.DatabaseURL)
	if err != nil {
		fmt.Println("Database error:", err)
		return err
	}
	defer func() { _ = db.Close() }()
	driver := postgres.NewDriver(db)

	// Source
	source := fs.NewSource(os.DirFS(env.MigrationsPath))

	// Migrator
	if err := migrator.NewMigrator(driver, source).Up(ctx); err != nil {
		fmt.Println("Migration error:", err)
		return err
	}

	fmt.Println("Migration done!")

	return nil
}

type MigrateCommandEnv struct {
	DatabaseURL    string
	MigrationsPath string
}

func MigrateCommandLoadEnvConfig() (MigrateCommandEnv, error) {
	cfg := MigrateCommandEnv{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
	}

	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.MigrationsPath == "" {
		cfg.MigrationsPath = "migrations"
	}

	return cfg, nil
}
