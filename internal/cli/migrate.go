package cli

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"

	migratepostgres "github.com/Palladium-blockchain/go-migrations/pkg/driver/postgres"
	"github.com/Palladium-blockchain/go-migrations/pkg/migrator"
	migratefs "github.com/Palladium-blockchain/go-migrations/pkg/source/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MigrateCommand struct{}

func NewMigrateCommand() *MigrateCommand {
	return &MigrateCommand{}
}

func (cmd *MigrateCommand) Execute(ctx context.Context, args []string) error {
	opts, err := parseMigrateCommandArgs(args)
	if err != nil {
		return err
	}

	// Config
	env, err := MigrateCommandLoadEnvConfig()
	if err != nil {
		fmt.Println("Config error:", err)
		return err
	}

	// Driver
	db, err := sql.Open("pgx", env.DatabaseURL)
	if err != nil {
		fmt.Println("Database error:", err)
		return err
	}
	defer func() { _ = db.Close() }()
	driver := migratepostgres.NewDriver(db)

	// Source
	source := migratefs.NewSource(os.DirFS(env.MigrationsPath))

	// Migrator
	migratorOptions := make([]migrator.Option, 0, 1)
	if opts.AllowOrphanedMigrations {
		migratorOptions = append(migratorOptions, migrator.WithAllowOrphanedMigrations())
	}

	if err := migrator.NewMigrator(driver, source, migratorOptions...).Up(ctx); err != nil {
		fmt.Println("Migration error:", err)
		return err
	}

	fmt.Println("Migration done!")

	return nil
}

type MigrateCommandOptions struct {
	AllowOrphanedMigrations bool
}

func parseMigrateCommandArgs(args []string) (MigrateCommandOptions, error) {
	var opts MigrateCommandOptions

	flags := flag.NewFlagSet("migrate", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.BoolVar(
		&opts.AllowOrphanedMigrations,
		"allow-orphaned-migrations",
		false,
		"ignore applied migrations that are missing locally",
	)

	if err := flags.Parse(args); err != nil {
		return opts, err
	}

	if flags.NArg() != 0 {
		return opts, fmt.Errorf("unexpected arguments: %v", flags.Args())
	}

	return opts, nil
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
