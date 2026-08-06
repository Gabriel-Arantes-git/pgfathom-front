package db

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// EnvDSN is the recommended way to pass a connection string.
const EnvDSN = "PGFATHOM_DSN"

// libpqVars are the standard environment variables pgx reads on its own.
var libpqVars = []string{"PGHOST", "PGPORT", "PGUSER", "PGDATABASE", "PGPASSWORD", "PGPASSFILE", "PGSERVICE"}

// ErrNoConnectionInfo is returned when nothing identifies a server to connect to.
var ErrNoConnectionInfo = errors.New("no connection information")

// ResolveDSN picks the connection string from the environment or the flag.
//
// PGFATHOM_DSN wins over --dsn because the flag's value shows up in ps and in
// shell history, and the documentation recommends the variable; if the flag
// won, that recommendation would be defeatable by accident on a shared machine.
// The libpq variables come last: they are ambient configuration rather than
// explicit intent, and PGHOST is commonly exported in the shell of anyone who
// works with PostgreSQL.
//
// An empty result with a nil error means the libpq environment is in play and
// pgx should read it itself.
func ResolveDSN(flagDSN string, warn func(string)) (string, error) {
	envDSN := strings.TrimSpace(os.Getenv(EnvDSN))
	flagDSN = strings.TrimSpace(flagDSN)

	switch {
	case envDSN != "":
		if flagDSN != "" && warn != nil {
			warn(fmt.Sprintf("both %s and --dsn are set; using %s", EnvDSN, EnvDSN))
		}
		return envDSN, nil

	case flagDSN != "":
		if warn != nil {
			warn("--dsn is visible in ps and in shell history; prefer " + EnvDSN)
		}
		return flagDSN, nil

	case hasLibpqEnv():
		return "", nil

	default:
		return "", fmt.Errorf("%w: set %s, pass --dsn, or use the standard PG* variables",
			ErrNoConnectionInfo, EnvDSN)
	}
}

func hasLibpqEnv() bool {
	for _, v := range libpqVars {
		if os.Getenv(v) != "" {
			return true
		}
	}
	return false
}

// Describe renders a connection target for diagnostics with the password
// removed. Every message that names a connection must go through here: a DSN
// pasted verbatim into an error is how credentials end up in CI logs.
func Describe(cfg *pgconn.Config) string {
	if cfg == nil {
		return "(unknown server)"
	}

	target := cfg.Host
	if cfg.Port != 0 {
		target = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	}

	parts := make([]string, 0, 3)
	if cfg.User != "" {
		parts = append(parts, "user="+cfg.User)
	}
	if cfg.Database != "" {
		parts = append(parts, "database="+cfg.Database)
	}
	if cfg.Password != "" {
		parts = append(parts, "password=[redacted]")
	}

	if len(parts) == 0 {
		return target
	}
	return target + " (" + strings.Join(parts, ", ") + ")"
}
