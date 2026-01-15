package cli

import (
	"context"
	"errors"
	"fmt"
	"palladium-intelligence/go-migrations/internal/creator/fs"
)

type CreateMigrationCommand struct{}

func NewCreateMigrationCommand() *CreateMigrationCommand {
	return &CreateMigrationCommand{}
}

func (cmd *CreateMigrationCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: ./migration create [path] [migration-name]")
		return errors.New("not enough arguments")
	}
	path := args[0]
	name := args[1]
	fmt.Printf("Creating new migration in: %s\n", path)

	migrationFiles, err := fs.NewCreator(path).Create(ctx, name)
	if err != nil {
		fmt.Printf("Error creating migration file: %v\n", err)
		return errors.New("error creating migration file")
	}
	fmt.Printf("Migration file created:\n- %s\n- %s\n", migrationFiles.Up, migrationFiles.Down)
	return nil
}
