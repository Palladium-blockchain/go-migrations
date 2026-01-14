package main

import (
	"context"
	"fmt"
	"os"
	"palladium-intelligence/go-migrations/internal/creator/fs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./migration [action] {...args}")
		os.Exit(1)
	}

	action := os.Args[1]
	args := os.Args[2:]

	switch action {
	case "create":
		handleCreateMigration(args)
	default:
		fmt.Printf("Unknown action: %s\n", action)
	}
}

func handleCreateMigration(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: ./migration create [path] [migration-name]")
		return
	}
	path := args[0]
	name := args[1]
	fmt.Printf("Creating new migration in: %s\n", path)

	migrationFiles, err := fs.NewCreator(path).Create(context.Background(), name)
	if err != nil {
		fmt.Printf("Error creating migration file: %v\n", err)
		return
	}
	fmt.Printf("Migration file created:\n- %s\n- %s\n", migrationFiles.Up, migrationFiles.Down)
}
