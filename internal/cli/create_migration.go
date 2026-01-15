package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Palladium-blockchain/go-migrations/internal/creator/fs"
)

type CreateMigrationCommand struct{}

func NewCreateMigrationCommand() *CreateMigrationCommand {
	return &CreateMigrationCommand{}
}

func (cmd *CreateMigrationCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: ./migration create [migration-name]")
		return errors.New("not enough arguments")
	}
	name := args[0]

	env, err := CreateMigrateCommandLoadEnvConfig()
	if err != nil {
		fmt.Println("Config error:", err)
		return err
	}

	fmt.Printf("Creating new migration in: %s\n", env.MigrationsPath)

	migrationFiles, err := fs.NewCreator(env.MigrationsPath).Create(ctx, name)
	if err != nil {
		fmt.Printf("Error creating migration file: %v\n", err)
		return errors.New("error creating migration file")
	}
	fmt.Printf("Migration file created:\n- %s\n- %s\n", migrationFiles.Up, migrationFiles.Down)
	return nil
}

type CreateMigrationCommandEnv struct {
	MigrationsPath string
}

func CreateMigrateCommandLoadEnvConfig() (CreateMigrationCommandEnv, error) {
	cfg := CreateMigrationCommandEnv{
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
	}

	if cfg.MigrationsPath == "" {
		return cfg, errors.New("env variable MIGRATIONS_PATH is required")
	}

	return cfg, nil
}
