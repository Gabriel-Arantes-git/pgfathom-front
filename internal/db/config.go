package db

import (
	"fmt"
	"time"
)

// ApplicationName identifies the tool in pg_stat_activity, so the DBA who
// authorized the run can see what is running without inferring it from query
// text.
const ApplicationName = "pgfathom"

// MinServerVersionNum is PostgreSQL 13, expressed as server_version_num.
// It covers every community-supported release and keeps the catalog queries
// free of version-conditional branches.
const MinServerVersionNum = 130000

// Default session limits, all deliberately conservative. This connects to
// someone else's production server.
const (
	DefaultStatementTimeout = 30 * time.Second
	DefaultLockTimeout      = 5 * time.Second
	DefaultIdleTxTimeout    = 60 * time.Second
	DefaultConcurrency      = 4
)

// Config describes how to connect and how much the run is allowed to cost the
// server.
type Config struct {
	// DSN is empty when the libpq environment should be used instead.
	DSN string

	StatementTimeout time.Duration
	LockTimeout      time.Duration

	// IdleTxTimeout bounds an open-but-idle transaction. It is the most
	// dangerous of the three: an idle transaction holds back VACUUM for the
	// whole database, not just for the tables being read.
	IdleTxTimeout time.Duration

	// Concurrency caps simultaneous queries and sizes the pool. Firing dozens
	// of queries at a production server is the fastest way for this tool to be
	// banned from a company.
	Concurrency int
}

// DefaultConfig returns the conservative defaults.
func DefaultConfig() Config {
	return Config{
		StatementTimeout: DefaultStatementTimeout,
		LockTimeout:      DefaultLockTimeout,
		IdleTxTimeout:    DefaultIdleTxTimeout,
		Concurrency:      DefaultConcurrency,
	}
}

// Validate rejects settings that would misbehave.
func (c Config) Validate() error {
	if c.Concurrency < 1 {
		return fmt.Errorf("concurrency must be at least 1, got %d", c.Concurrency)
	}
	for _, t := range []struct {
		name  string
		value time.Duration
	}{
		{"statement timeout", c.StatementTimeout},
		{"lock timeout", c.LockTimeout},
		{"idle transaction timeout", c.IdleTxTimeout},
	} {
		if t.value <= 0 {
			return fmt.Errorf("%s must be positive, got %s", t.name, t.value)
		}
	}
	return nil
}
