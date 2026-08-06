package db

import (
	"errors"
	"strings"
	"testing"
)

// clearEnv removes every variable that participates in resolution, so a
// developer's own PGHOST cannot make the suite pass or fail by accident.
func clearEnv(t *testing.T) {
	t.Helper()
	t.Setenv(EnvDSN, "")
	for _, v := range libpqVars {
		t.Setenv(v, "")
	}
}

func TestResolveDSNPrecedence(t *testing.T) {
	const envValue = "postgres://env/db"
	const flagValue = "postgres://flag/db"

	t.Run("PGFATHOM_DSN beats the flag", func(t *testing.T) {
		clearEnv(t)
		t.Setenv(EnvDSN, envValue)

		var warnings []string
		got, err := ResolveDSN(flagValue, func(m string) { warnings = append(warnings, m) })
		if err != nil {
			t.Fatalf("ResolveDSN: %v", err)
		}
		if got != envValue {
			t.Errorf("dsn = %q, want %q", got, envValue)
		}
		if len(warnings) == 0 {
			t.Error("a conflict between the variable and the flag should be warned about")
		}
	})

	// Explicit intent beats ambient configuration: PGHOST is commonly exported
	// in the shell of anyone who works with PostgreSQL, and letting it override
	// a hand-written --dsn would be surprising.
	t.Run("flag beats the libpq environment", func(t *testing.T) {
		clearEnv(t)
		t.Setenv("PGHOST", "ambient.example")

		got, err := ResolveDSN(flagValue, nil)
		if err != nil {
			t.Fatalf("ResolveDSN: %v", err)
		}
		if got != flagValue {
			t.Errorf("dsn = %q, want %q", got, flagValue)
		}
	})

	t.Run("the flag warns about leaking into ps", func(t *testing.T) {
		clearEnv(t)

		var warnings []string
		if _, err := ResolveDSN(flagValue, func(m string) { warnings = append(warnings, m) }); err != nil {
			t.Fatalf("ResolveDSN: %v", err)
		}
		if len(warnings) != 1 || !strings.Contains(warnings[0], EnvDSN) {
			t.Errorf("warnings = %v, want one pointing at %s", warnings, EnvDSN)
		}
	})

	t.Run("libpq environment alone yields an empty DSN for pgx to resolve", func(t *testing.T) {
		clearEnv(t)
		t.Setenv("PGHOST", "ambient.example")

		got, err := ResolveDSN("", nil)
		if err != nil {
			t.Fatalf("ResolveDSN: %v", err)
		}
		if got != "" {
			t.Errorf("dsn = %q, want empty so pgx reads the environment", got)
		}
	})

	t.Run("nothing at all is an error", func(t *testing.T) {
		clearEnv(t)

		if _, err := ResolveDSN("", nil); !errors.Is(err, ErrNoConnectionInfo) {
			t.Fatalf("error = %v, want ErrNoConnectionInfo", err)
		}
	})

	t.Run("whitespace does not count as a value", func(t *testing.T) {
		clearEnv(t)
		t.Setenv(EnvDSN, "   ")

		if _, err := ResolveDSN("  ", nil); !errors.Is(err, ErrNoConnectionInfo) {
			t.Fatalf("error = %v, want ErrNoConnectionInfo", err)
		}
	})
}

// TestDescribeRedactsPassword covers the rule that a DSN pasted verbatim into
// an error is how credentials end up in CI logs.
func TestDescribeRedactsPassword(t *testing.T) {
	const secret = "hunter2-super-secret"

	cfg, err := ConnConfig("postgres://alice:" + secret + "@db.example:5433/prefeitura")
	if err != nil {
		t.Fatalf("ConnConfig: %v", err)
	}

	got := Describe(cfg)

	if strings.Contains(got, secret) {
		t.Fatalf("Describe leaked the password: %q", got)
	}
	for _, want := range []string{"db.example", "5433", "alice", "prefeitura", "[redacted]"} {
		if !strings.Contains(got, want) {
			t.Errorf("Describe(%q) = %q, want it to mention %q", "…", got, want)
		}
	}
}

func TestDescribeHandlesMissingConfig(t *testing.T) {
	if got := Describe(nil); got == "" {
		t.Error("Describe(nil) must still produce something printable")
	}
}

func TestConfigValidate(t *testing.T) {
	valid := DefaultConfig()
	if err := valid.Validate(); err != nil {
		t.Fatalf("the defaults must be valid: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*Config)
		wantIn string
	}{
		{"zero concurrency", func(c *Config) { c.Concurrency = 0 }, "concurrency"},
		{"negative concurrency", func(c *Config) { c.Concurrency = -1 }, "concurrency"},
		{"zero statement timeout", func(c *Config) { c.StatementTimeout = 0 }, "statement timeout"},
		{"zero lock timeout", func(c *Config) { c.LockTimeout = 0 }, "lock timeout"},
		{"zero idle transaction timeout", func(c *Config) { c.IdleTxTimeout = 0 }, "idle transaction"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			tt.mutate(&cfg)

			err := cfg.Validate()
			if err == nil {
				t.Fatal("Validate accepted a setting that would misbehave")
			}
			if !strings.Contains(err.Error(), tt.wantIn) {
				t.Errorf("error = %q, want it to mention %q", err, tt.wantIn)
			}
		})
	}
}

func TestFormatVersion(t *testing.T) {
	for _, tt := range []struct {
		num  int
		want string
	}{
		{130000, "13.0"},
		{160002, "16.2"},
		{0, "unknown"},
		{-1, "unknown"},
	} {
		if got := formatVersion(tt.num); got != tt.want {
			t.Errorf("formatVersion(%d) = %q, want %q", tt.num, got, tt.want)
		}
	}
}

// TestSupportedFloorIsThirteen pins the decision so it cannot drift silently:
// the catalog queries are written free of version-conditional branches on the
// strength of it.
func TestSupportedFloorIsThirteen(t *testing.T) {
	if MinServerVersionNum != 130000 {
		t.Errorf("MinServerVersionNum = %d, want 130000 (PostgreSQL 13)", MinServerVersionNum)
	}
}
