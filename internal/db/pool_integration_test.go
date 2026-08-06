//go:build integration

package db_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/db"
	"github.com/lvcas-dotcom/pgfathom/internal/testutil"
)

// TestRefusesUnsupportedServer proves the floor is enforced before any catalog
// is read. Failing early with a clear message beats failing mid-read with a
// missing-column error the user cannot interpret.
func TestRefusesUnsupportedServer(t *testing.T) {
	const oldImage = "postgres:12-alpine"

	cfg := db.DefaultConfig()
	cfg.DSN = testutil.PostgresImageDSN(t, oldImage, "clean_schema")

	pool, err := db.Open(context.Background(), cfg)
	if err == nil {
		pool.Close()
		t.Fatal("connected to a server below the supported floor")
	}
	if !errors.Is(err, db.ErrServerTooOld) {
		t.Fatalf("error = %v, want ErrServerTooOld", err)
	}
	for _, want := range []string{"12", "13"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should name both the version found and the one required", err)
		}
	}
}

func TestWritePrivilegeDetection(t *testing.T) {
	ctx := context.Background()

	cfg := db.DefaultConfig()
	cfg.DSN = testutil.Postgres(t, "restricted_privileges")

	pool, err := db.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("opening pool: %v", err)
	}
	defer pool.Close()

	// The fixture's owner role can write; the point of the check is to tell the
	// operator so, since whoever runs the tool is often not whoever chose the role.
	writable, err := pool.HasWritePrivileges(ctx, []string{"public"})
	if err != nil {
		t.Fatalf("checking privileges: %v", err)
	}
	if !writable {
		t.Error("the owning role can write and the check should say so")
	}
}
