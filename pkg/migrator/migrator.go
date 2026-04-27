package migrator

import (
	internalmigrator "github.com/Palladium-blockchain/go-migrations/internal/migrator"
	"github.com/Palladium-blockchain/go-migrations/pkg/migrate"
)

type Migrator = internalmigrator.Migrator

func NewMigrator(driver migrate.Driver, source migrate.Source) *Migrator {
	return internalmigrator.NewMigrator(driver, source)
}
