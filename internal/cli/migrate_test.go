package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMigrateCommandArgs_AllowOrphanedMigrations(t *testing.T) {
	opts, err := parseMigrateCommandArgs([]string{"--allow-orphaned-migrations"})
	require.NoError(t, err)
	require.True(t, opts.AllowOrphanedMigrations)
}

func TestParseMigrateCommandArgs_UnexpectedArgs(t *testing.T) {
	_, err := parseMigrateCommandArgs([]string{"extra"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected arguments")
}
