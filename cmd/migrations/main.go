package main

import (
	"context"
	"fmt"
	"os"
	"palladium-intelligence/go-migrations/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./migration [action] {...args}")
		os.Exit(1)
	}

	action := os.Args[1]
	args := os.Args[2:]

	var cmd cli.Command
	switch action {
	case "create":
		cmd = cli.NewCreateMigrationCommand()
	case "migrate":
		cmd = cli.NewMigrateCommand()
	default:
		fmt.Printf("Unknown action: %s\n", action)
		os.Exit(1)
	}

	if err := cmd.Execute(context.Background(), args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
