package migrator

import (
	"context"
	"testing"

	"github.com/Palladium-blockchain/go-migrations/pkg/migrate"
	"github.com/stretchr/testify/require"
)

type stubDriver struct {
	applied  []string
	appliedM []migrate.Migration
}

func (d *stubDriver) Initialize(context.Context) error {
	return nil
}

func (d *stubDriver) Apply(_ context.Context, m migrate.Migration) error {
	d.appliedM = append(d.appliedM, m)
	return nil
}

func (d *stubDriver) Rollback(context.Context, migrate.Migration) error {
	return nil
}

func (d *stubDriver) ListApplied(context.Context) ([]string, error) {
	return append([]string(nil), d.applied...), nil
}

func (d *stubDriver) Lock(context.Context) error {
	return nil
}

func (d *stubDriver) Unlock(context.Context) error {
	return nil
}

type stubSource struct {
	migrations []migrate.Migration
}

func (s *stubSource) Load(context.Context) ([]migrate.Migration, error) {
	return append([]migrate.Migration(nil), s.migrations...), nil
}

func TestMigrator_Up_ErrorsOnUnknownAppliedMigration(t *testing.T) {
	driver := &stubDriver{
		applied: []string{"001_missing"},
	}
	source := &stubSource{
		migrations: []migrate.Migration{
			{ID: "002_known", Up: []byte("SELECT 1")},
		},
	}

	err := NewMigrator(driver, source).Up(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown applied migration: 001_missing")
	require.Empty(t, driver.appliedM)
}

func TestMigrator_Up_AllowsOrphanedMigrations(t *testing.T) {
	driver := &stubDriver{
		applied: []string{"001_missing"},
	}
	source := &stubSource{
		migrations: []migrate.Migration{
			{ID: "002_known", Up: []byte("SELECT 1")},
		},
	}

	err := NewMigrator(driver, source, WithAllowOrphanedMigrations()).Up(context.Background())
	require.NoError(t, err)
	require.Len(t, driver.appliedM, 1)
	require.Equal(t, "002_known", driver.appliedM[0].ID)
}

func TestMigrator_Down_StillErrorsOnUnknownAppliedMigration(t *testing.T) {
	driver := &stubDriver{
		applied: []string{"001_known", "002_missing"},
	}
	source := &stubSource{
		migrations: []migrate.Migration{
			{ID: "001_known", Up: []byte("SELECT 1"), Down: []byte("SELECT 2")},
		},
	}

	err := NewMigrator(driver, source, WithAllowOrphanedMigrations()).Down(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown applied migration: 002_missing")
}
