package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/Palladium-blockchain/go-migrations/internal/cli"
)

func main() {
	os.Exit(run(os.Args, os.Stdout))
}

func run(args []string, out io.Writer) int {
	if isVersionRequest(args) {
		fmt.Fprintln(out, currentVersion())
		return 0
	}

	if len(args) < 2 {
		fmt.Fprintln(out, "Usage: ./migration [action] {...args}")
		return 1
	}

	action := args[1]
	cmdArgs := args[2:]

	var cmd cli.Command
	switch action {
	case "create":
		cmd = cli.NewCreateMigrationCommand()
	case "migrate":
		cmd = cli.NewMigrateCommand()
	default:
		fmt.Fprintf(out, "Unknown action: %s\n", action)
		return 1
	}

	if err := cmd.Execute(context.Background(), cmdArgs); err != nil {
		fmt.Fprintln(out, err)
		return 1
	}

	return 0
}

func isVersionRequest(args []string) bool {
	return len(args) == 2 && (args[1] == "-v" || args[1] == "--version")
}

func currentVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	if info.Main.Version == "" {
		return "unknown"
	}

	return info.Main.Version
}
