// Package db owns the connection and everything that makes it safe to point at
// somebody else's production server.
package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrServerTooOld is returned when the server predates the supported floor.
var ErrServerTooOld = errors.New("server version not supported")

// Pool is a read-only connection pool.
type Pool struct {
	pool          *pgxpool.Pool
	serverVersion int
	target        string
}

// Open connects, applies the session policies and verifies the server version.
func Open(ctx context.Context, cfg Config) (*Pool, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parsing connection settings: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.Concurrency)
	poolCfg.ConnConfig.RuntimeParams["application_name"] = ApplicationName

	// Session policies belong here rather than in each query. AfterConnect runs
	// once per new connection, including the ones the pool opens on its own
	// under load, so no connection can exist without them. Applying them
	// per-query would rely on every caller remembering — and a connection
	// opened down some forgotten path is exactly the hole that only shows up in
	// someone else's production.
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return applySessionPolicies(ctx, conn, cfg)
	}

	target := Describe(poolCfg.ConnConfig.Config.Copy())

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("connecting to %s: %w", target, err)
	}

	p := &Pool{pool: pool, target: target}

	if err := p.checkServerVersion(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return p, nil
}

func applySessionPolicies(ctx context.Context, conn *pgx.Conn, cfg Config) error {
	statements := []struct {
		what string
		sql  string
	}{
		{"read-only mode", "SET default_transaction_read_only = on"},
		{"statement timeout", fmt.Sprintf("SET statement_timeout = %d", cfg.StatementTimeout.Milliseconds())},
		{"lock timeout", fmt.Sprintf("SET lock_timeout = %d", cfg.LockTimeout.Milliseconds())},
		{"idle transaction timeout", fmt.Sprintf("SET idle_in_transaction_session_timeout = %d", cfg.IdleTxTimeout.Milliseconds())},
	}

	for _, s := range statements {
		if _, err := conn.Exec(ctx, s.sql); err != nil {
			return fmt.Errorf("applying %s: %w", s.what, err)
		}
	}
	return nil
}

// checkServerVersion fails before any catalog is read. Failing early with a
// clear message beats failing mid-read with a missing-column error the user has
// no way to interpret.
func (p *Pool) checkServerVersion(ctx context.Context) error {
	// SHOW returns text for every setting, so the value is read through
	// current_setting with an explicit cast instead.
	var num int
	const query = `SELECT current_setting('server_version_num')::int`
	if err := p.pool.QueryRow(ctx, query).Scan(&num); err != nil {
		return fmt.Errorf("reading server version from %s: %w", p.target, err)
	}

	p.serverVersion = num
	if num < MinServerVersionNum {
		return fmt.Errorf("%w: %s runs %s, pgfathom requires %s or newer",
			ErrServerTooOld, p.target, formatVersion(num), formatVersion(MinServerVersionNum))
	}
	return nil
}

// ServerVersionNum is the server version as server_version_num.
func (p *Pool) ServerVersionNum() int { return p.serverVersion }

// ServerVersion is the server version rendered for humans.
func (p *Pool) ServerVersion() string { return formatVersion(p.serverVersion) }

// Target describes the server for diagnostics, with the password removed.
func (p *Pool) Target() string { return p.target }

// Query runs a read-only query.
func (p *Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

// QueryRow runs a read-only query expected to return a single row.
func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

// Close releases every connection.
func (p *Pool) Close() { p.pool.Close() }

// HasWritePrivileges reports whether the connected role can write to any table
// in scope, in a single aggregated query rather than one per table.
//
// The read-only session already blocks writes. The warning exists because
// whoever runs the tool is often not whoever chose the role, and the
// recommendation to use a dedicated read-only role has to reach that person.
func (p *Pool) HasWritePrivileges(ctx context.Context, schemas []string) (bool, error) {
	const query = `
		SELECT EXISTS (
		    SELECT 1
		    FROM pg_class c
		    JOIN pg_namespace n ON n.oid = c.relnamespace
		    WHERE c.relkind IN ('r', 'p')
		      AND n.nspname = ANY($1)
		      AND (has_table_privilege(c.oid, 'INSERT')
		        OR has_table_privilege(c.oid, 'UPDATE')
		        OR has_table_privilege(c.oid, 'DELETE'))
		)`

	var writable bool
	if err := p.pool.QueryRow(ctx, query, schemas).Scan(&writable); err != nil {
		// Visibility this query needs may itself be restricted. Degrade to less
		// information rather than to not working.
		return false, fmt.Errorf("checking write privileges: %w", err)
	}
	return writable, nil
}

// formatVersion renders server_version_num as a human version. Since
// PostgreSQL 10 the encoding is major*10000 + minor.
func formatVersion(num int) string {
	if num <= 0 {
		return "unknown"
	}
	return fmt.Sprintf("%d.%d", num/10000, num%10000)
}

// ConnConfig exposes the parsed configuration for diagnostics.
func ConnConfig(dsn string) (*pgconn.Config, error) {
	cfg, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parsing connection settings: %w", err)
	}
	return cfg, nil
}
