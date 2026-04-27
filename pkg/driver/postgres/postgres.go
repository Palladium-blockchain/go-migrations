package postgres

import (
	"database/sql"

	internalpostgres "github.com/Palladium-blockchain/go-migrations/internal/driver/postgres"
)

type Driver = internalpostgres.Driver

func NewDriver(db *sql.DB) *Driver {
	return internalpostgres.NewDriver(db)
}
