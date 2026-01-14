package fs

import (
	"context"
	"fmt"
	"os"
	"palladium-intelligence/go-migrations/pkg/migrate"
	"path/filepath"
	"time"
)

type Creator struct {
	dir string
}

func NewCreator(dir string) *Creator {
	return &Creator{dir: dir}
}

func (c *Creator) Create(_ context.Context, name string) (migrate.MigrationsFile, error) {
	version := time.Now().UTC().Format("20060102150405")

	up := filepath.Join(c.dir, fmt.Sprintf("%s_%s.up.sql", version, name))
	if err := os.WriteFile(up, []byte("-- up\n"), 0644); err != nil {
		return migrate.MigrationsFile{}, err
	}

	down := filepath.Join(c.dir, fmt.Sprintf("%s_%s.down.sql", version, name))
	if err := os.WriteFile(down, []byte("-- down\n"), 0644); err != nil {
		return migrate.MigrationsFile{}, err
	}

	return migrate.MigrationsFile{
		Up:   up,
		Down: down,
	}, nil
}
