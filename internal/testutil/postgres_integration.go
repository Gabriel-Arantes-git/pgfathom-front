//go:build integration

package testutil

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresImage is the server the integration suite runs against. It sits at
// the supported floor rather than the newest release, so a query that silently
// depends on a newer catalog fails here instead of in a user's database.
const PostgresImage = "postgres:13-alpine"

// Postgres starts a throwaway server, loads a named fixture and returns a DSN.
//
// A real server is the only honest way to test catalog reading: a fake would
// test the fake. The container is torn down when the test ends.
func Postgres(t *testing.T, fixture string) string {
	t.Helper()
	return PostgresImageDSN(t, PostgresImage, fixture)
}

// PostgresImageDSN is Postgres with the server image chosen by the caller, so a
// test can assert behaviour against a version the tool does not support.
func PostgresImageDSN(t *testing.T, image, fixture string) string {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.Run(ctx, image,
		postgres.WithDatabase("pgfathom_test"),
		postgres.WithUsername("pgfathom_test"),
		postgres.WithPassword("pgfathom_test"),
		postgres.WithInitScripts(fixturePath(t, fixture)),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(2*time.Minute),
		),
	)
	if err != nil {
		t.Fatalf("starting PostgreSQL: %v", err)
	}

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Logf("terminating container: %v", err)
		}
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("resolving connection string: %v", err)
	}
	return dsn
}

// fixturePath resolves a fixture by name, failing loudly when it is missing so
// a typo cannot turn into a silently empty schema.
func fixturePath(t *testing.T, fixture string) string {
	t.Helper()

	path := filepath.Join("testdata", fixture+".sql")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("fixture %q not found at %s: %v", fixture, path, err)
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("resolving fixture path: %v", err)
	}
	return abs
}
