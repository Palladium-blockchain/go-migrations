package migrate

import "context"

type Migrator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
}

type Driver interface {
	Initialize(ctx context.Context) error
	Apply(ctx context.Context, m Migration) error
	Rollback(ctx context.Context, m Migration) error
	ListApplied(ctx context.Context) ([]string, error)
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}

type Migration struct {
	ID   string
	Up   []byte
	Down []byte
}

type Source interface {
	Load(ctx context.Context) ([]Migration, error)
}

type Creator interface {
	Create(ctx context.Context, name string) (MigrationsFile, error)
}

type MigrationsFile struct {
	Up   string
	Down string
}
