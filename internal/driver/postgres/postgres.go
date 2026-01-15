package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Palladium-blockchain/go-migrations/pkg/migrate"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Driver struct {
	db   *sql.DB
	conn *sql.Conn
}

func NewDriver(db *sql.DB) *Driver {
	return &Driver{
		db:   db,
		conn: nil,
	}
}

func (d *Driver) Initialize(ctx context.Context) error {
	const q = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
    	version TEXT NOT NULL,
    	applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := d.db.ExecContext(ctx, q)
	return err
}

func (d *Driver) Apply(ctx context.Context, m migrate.Migration) error {
	if d.conn == nil {
		return errors.New("driver is not locked")
	}

	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, string(m.Up)); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `INSERT INTO schema_migrations(version) VALUES ($1)`, m.ID); err != nil {
		return err
	}

	return tx.Commit()
}

func (d *Driver) Rollback(ctx context.Context, m migrate.Migration) error {
	if len(m.Down) == 0 {
		return fmt.Errorf("migration %s has no down.sql", m.ID)
	}

	if d.conn == nil {
		return errors.New("driver is not locked")
	}

	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, string(m.Down)); err != nil {
		return fmt.Errorf("rollback %s failed: %w", m.ID, err)
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM schema_migrations WHERE version = $1`, m.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return fmt.Errorf("migration %s not applied", m.ID)
	}

	return tx.Commit()
}

func (d *Driver) ListApplied(ctx context.Context) ([]string, error) {
	if d.conn == nil {
		return nil, errors.New("driver is not locked")
	}

	rows, err := d.conn.QueryContext(ctx,
		`SELECT version FROM schema_migrations ORDER BY version`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []string

	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

func (d *Driver) Lock(ctx context.Context) error {
	if d.conn != nil {
		return nil
	}

	conn, err := d.db.Conn(ctx)
	if err != nil {
		return err
	}

	if _, err := conn.ExecContext(ctx, `SELECT pg_advisory_lock(424242)`); err != nil {
		conn.Close()
		return err
	}

	d.conn = conn
	return nil
}

func (d *Driver) Unlock(ctx context.Context) error {
	if d.conn == nil {
		return nil
	}

	_, err := d.conn.ExecContext(ctx, `SELECT pg_advisory_unlock(424242)`)
	_ = d.conn.Close()
	d.conn = nil

	return err
}
