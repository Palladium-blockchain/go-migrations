package migrator

import (
	internalmigrator "github.com/Palladium-blockchain/go-migrations/internal/migrator"
	"github.com/Palladium-blockchain/go-migrations/pkg/migrate"
)

type Migrator = internalmigrator.Migrator
type Option = internalmigrator.Option

func WithAllowOrphanedMigrations() Option {
	return internalmigrator.WithAllowOrphanedMigrations()
}

func NewMigrator(driver migrate.Driver, source migrate.Source, opts ...Option) *Migrator {
	return internalmigrator.NewMigrator(driver, source, opts...)
}
