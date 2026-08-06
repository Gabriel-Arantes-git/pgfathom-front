package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/lvcas-dotcom/pgfathom/internal/audit"
	"github.com/lvcas-dotcom/pgfathom/internal/buildinfo"
	"github.com/lvcas-dotcom/pgfathom/internal/catalog"
	"github.com/lvcas-dotcom/pgfathom/internal/db"
	"github.com/lvcas-dotcom/pgfathom/internal/model"
	"github.com/lvcas-dotcom/pgfathom/internal/report"
)

type auditOptions struct {
	dsn              string
	schemas          []string
	exclude          []string
	format           string
	statementTimeout time.Duration
	lockTimeout      time.Duration
	idleTxTimeout    time.Duration
	concurrency      int
}

func newAuditCommand(streams *Streams) *cobra.Command {
	opts := &auditOptions{
		schemas:          []string{"public"},
		format:           "table",
		statementTimeout: db.DefaultStatementTimeout,
		lockTimeout:      db.DefaultLockTimeout,
		idleTxTimeout:    db.DefaultIdleTxTimeout,
		concurrency:      db.DefaultConcurrency,
	}

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Report structural findings that need no inference",
		Long: `Report structural findings taken straight from the catalog.

audit makes no inferences: every finding it emits is a fact about the schema,
not a hypothesis about the data. It reports foreign keys declared NOT VALID and
never verified, and foreign keys with no index on the child side.

Relationship inference is a separate command; one does not substitute for the
other.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runAudit(cmd.Context(), streams, opts)
		},
	}

	f := cmd.Flags()
	f.StringVar(&opts.dsn, "dsn", "",
		"connection string (visible in ps and shell history; prefer "+db.EnvDSN+")")
	f.StringSliceVar(&opts.schemas, "schema", opts.schemas, "schemas to analyze")
	f.StringSliceVar(&opts.exclude, "exclude", nil, "glob patterns of tables to skip")
	f.StringVar(&opts.format, "format", opts.format, "output format: table or json")
	f.DurationVar(&opts.statementTimeout, "timeout", opts.statementTimeout, "statement timeout per query")
	f.DurationVar(&opts.lockTimeout, "lock-timeout", opts.lockTimeout, "lock timeout per query")
	f.DurationVar(&opts.idleTxTimeout, "idle-tx-timeout", opts.idleTxTimeout, "idle transaction timeout")
	f.IntVar(&opts.concurrency, "concurrency", opts.concurrency, "maximum simultaneous queries")

	return cmd
}

func runAudit(ctx context.Context, streams *Streams, opts *auditOptions) error {
	if opts.format != "table" && opts.format != "json" {
		return UsageError(fmt.Errorf("invalid --format %q: want table or json", opts.format))
	}

	warn := func(msg string) { _, _ = fmt.Fprintln(streams.Err, "warning: "+msg) }

	dsn, err := db.ResolveDSN(opts.dsn, warn)
	if err != nil {
		return err
	}

	cfg := db.Config{
		DSN:              dsn,
		StatementTimeout: opts.statementTimeout,
		LockTimeout:      opts.lockTimeout,
		IdleTxTimeout:    opts.idleTxTimeout,
		Concurrency:      opts.concurrency,
	}

	pool, err := db.Open(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	schemas := catalog.SortedSchemas(opts.schemas)

	// A failure here means the role cannot see the privilege catalog, which is
	// not a reason to stop — degrade to less information, not to not working.
	if writable, err := pool.HasWritePrivileges(ctx, schemas); err != nil {
		warn("could not verify privileges: " + err.Error())
	} else if writable {
		warn("the connected role can write to tables in scope; " +
			"a dedicated read-only role is recommended")
	}

	cat, err := catalog.Read(ctx, pool, catalog.Options{Schemas: schemas, Exclude: opts.exclude})
	if err != nil {
		return err
	}

	version, _, _ := buildinfo.Resolve()
	result := model.NewResult(version, "", time.Now().UTC(), cat.Coverage)
	result.ServerVersion = pool.ServerVersion()
	result.Schemas = cat.Schemas
	result.Findings = audit.Findings(cat.Schemas)

	if opts.format == "json" {
		return report.JSON(streams.Out, result)
	}
	return report.Terminal(streams.Out, result)
}
