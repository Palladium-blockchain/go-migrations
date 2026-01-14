package migrator

import (
	"context"
	"fmt"
	"palladium-intelligence/go-migrations/pkg/migrate"
)

type Migrator struct {
	driver migrate.Driver
	source migrate.Source
}

func NewMigrator(
	driver migrate.Driver,
	source migrate.Source,
) *Migrator {
	return &Migrator{
		driver: driver,
		source: source,
	}
}

func (m *Migrator) Up(ctx context.Context) error {
	migrations, err := m.source.Load(ctx)
	if err != nil {
		return err
	}

	if err := m.driver.Lock(ctx); err != nil {
		return err
	}
	defer m.driver.Unlock(ctx)

	if err := m.driver.Initialize(ctx); err != nil {
		return err
	}

	applied, err := m.driver.ListApplied(ctx)
	if err != nil {
		return err
	}
	appliedMap := make(map[string]struct{}, len(applied))
	for _, applied := range applied {
		appliedMap[applied] = struct{}{}
	}

	known := make(map[string]struct{}, len(migrations))
	for _, mig := range migrations {
		known[mig.ID] = struct{}{}
	}
	for _, a := range applied {
		if _, ok := known[a]; !ok {
			return fmt.Errorf("unknown applied migration: %s", a)
		}
	}

	for _, migration := range migrations {
		if _, ok := appliedMap[migration.ID]; ok {
			continue
		}

		if err := m.driver.Apply(ctx, migration); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) Down(ctx context.Context) error {
	migrations, err := m.source.Load(ctx)
	if err != nil {
		return err
	}

	if err := m.driver.Lock(ctx); err != nil {
		return err
	}
	defer m.driver.Unlock(ctx)

	if err := m.driver.Initialize(ctx); err != nil {
		return err
	}

	applied, err := m.driver.ListApplied(ctx)
	if err != nil {
		return err
	}

	if len(applied) == 0 {
		return migrate.ErrNoChange
	}

	known := make(map[string]migrate.Migration, len(migrations))
	for _, mig := range migrations {
		known[mig.ID] = mig
	}

	lastID := applied[len(applied)-1]

	mig, ok := known[lastID]
	if !ok {
		return fmt.Errorf("unknown applied migration: %s", lastID)
	}

	if err := m.driver.Rollback(ctx, mig); err != nil {
		return err
	}

	return nil
}
