//go:build integration

package db_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

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

// TestCancellationEndsTheQueryOnTheServer proves cancellation reaches the
// server, not merely the client. A client that gives up while the backend keeps
// running leaves a query loose on someone's production server — the exact
// failure the read-only session policies exist to prevent.
func TestCancellationEndsTheQueryOnTheServer(t *testing.T) {
	ctx := context.Background()
	dsn := testutil.Postgres(t, "clean_schema")

	cfg := db.DefaultConfig()
	cfg.DSN = dsn

	worker, err := db.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("opening worker pool: %v", err)
	}
	defer worker.Close()

	observer, err := db.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("opening observer pool: %v", err)
	}
	defer observer.Close()

	// The observer's own query also mentions pg_sleep, so it excludes itself by
	// backend pid.
	const activeSleeps = `
		SELECT count(*) FROM pg_stat_activity
		WHERE application_name = 'pgfathom'
		  AND state = 'active'
		  AND query LIKE '%pg_sleep%'
		  AND pid <> pg_backend_pid()`

	countSleeps := func() int {
		var n int
		if err := observer.QueryRow(ctx, activeSleeps).Scan(&n); err != nil {
			t.Fatalf("reading pg_stat_activity: %v", err)
		}
		return n
	}

	queryCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		rows, err := worker.Query(queryCtx, "SELECT pg_sleep(30)")
		if err == nil {
			rows.Close()
			err = rows.Err()
		}
		done <- err
	}()

	// Wait until the sleep is genuinely running server-side before cancelling;
	// cancelling earlier would test the client-refusal path, which is cheaper
	// and already covered.
	deadline := time.Now().Add(10 * time.Second)
	for countSleeps() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("the sleeping query never showed up in pg_stat_activity")
		}
		time.Sleep(50 * time.Millisecond)
	}

	cancel()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("the query completed despite cancellation")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the client call did not return after cancellation")
	}

	// The backend must stop, and promptly — not at statement_timeout, 30 seconds
	// later, which would pass a lazy implementation.
	deadline = time.Now().Add(5 * time.Second)
	for countSleeps() != 0 {
		if time.Now().After(deadline) {
			t.Fatal("the query is still running on the server after cancellation")
		}
		time.Sleep(50 * time.Millisecond)
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
