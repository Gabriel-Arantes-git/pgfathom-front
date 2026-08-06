// Package catalog reads the PostgreSQL system catalog into the semantic model.
package catalog

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/lvcas-dotcom/pgfathom/internal/db"
	"github.com/lvcas-dotcom/pgfathom/internal/model"
)

// Options scopes what gets read.
type Options struct {
	// Schemas to analyze. Empty means public.
	Schemas []string

	// Exclude holds glob patterns matched against both "table" and
	// "schema.table".
	Exclude []string
}

func (o Options) schemas() []string {
	if len(o.Schemas) == 0 {
		return []string{"public"}
	}
	return o.Schemas
}

// Result is what a catalog read produced, together with what it could not.
type Result struct {
	Schemas  []model.Schema
	Coverage model.Coverage
}

// Read loads the catalog for the schemas in scope.
func Read(ctx context.Context, pool *db.Pool, opts Options) (*Result, error) {
	schemas := opts.schemas()

	tables, coverage, err := readRelations(ctx, pool, schemas, opts.Exclude)
	if err != nil {
		return nil, err
	}

	if err := readColumns(ctx, pool, schemas, tables); err != nil {
		return nil, err
	}
	if err := readConstraints(ctx, pool, schemas, tables); err != nil {
		return nil, err
	}
	if err := readIndexes(ctx, pool, schemas, tables); err != nil {
		return nil, err
	}
	if err := readUsage(ctx, pool, schemas, tables, &coverage); err != nil {
		return nil, err
	}

	linkForeignKeyIndexes(tables)
	classifyUnsupported(tables, &coverage)

	return &Result{Schemas: group(schemas, tables), Coverage: coverage}, nil
}

// tableKey identifies a table across the several passes.
type tableKey struct{ schema, name string }

func key(schema, name string) tableKey { return tableKey{schema, name} }

func readRelations(ctx context.Context, pool *db.Pool, schemas, exclude []string) (map[tableKey]*model.Table, model.Coverage, error) {
	var coverage model.Coverage

	rows, err := pool.Query(ctx, queryRelations, schemas)
	if err != nil {
		return nil, coverage, fmt.Errorf("reading relations: %w", err)
	}
	defer rows.Close()

	tables := make(map[tableKey]*model.Table)

	for rows.Next() {
		var (
			t             model.Table
			estimatedRows int64
			totalBytes    int64
			readable      bool
		)
		if err := rows.Scan(&t.Schema, &t.Name, &t.Partitioned, &t.Inherits,
			&estimatedRows, &totalBytes, &t.Comment, &readable); err != nil {
			return nil, coverage, fmt.Errorf("scanning relation: %w", err)
		}

		coverage.TablesTotal++
		ref := t.Schema + "." + t.Name

		if matchesAny(exclude, t.Schema, t.Name) {
			coverage.TablesExcluded = append(coverage.TablesExcluded, ref)
			continue
		}

		// A table the role cannot SELECT is recorded, never silently dropped: a
		// partial grant is the common case in the target environment, and
		// silence about it would make an incomplete report look clean.
		if !readable {
			coverage.TablesNoPrivilege = append(coverage.TablesNoPrivilege, ref)
			continue
		}

		t.Stats = model.TableStats{EstimatedRows: estimatedRows, TotalBytes: totalBytes}
		tables[key(t.Schema, t.Name)] = &t
	}
	if err := rows.Err(); err != nil {
		return nil, coverage, fmt.Errorf("reading relations: %w", err)
	}

	return tables, coverage, nil
}

func readColumns(ctx context.Context, pool *db.Pool, schemas []string, tables map[tableKey]*model.Table) error {
	rows, err := pool.Query(ctx, queryColumns, schemas)
	if err != nil {
		return fmt.Errorf("reading columns: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schema, table string
		var c model.Column
		if err := rows.Scan(&schema, &table, &c.Name, &c.Type, &c.BaseType,
			&c.Nullable, &c.Default, &c.Position, &c.Comment); err != nil {
			return fmt.Errorf("scanning column: %w", err)
		}
		if t, ok := tables[key(schema, table)]; ok {
			t.Columns = append(t.Columns, c)
		}
	}
	return wrap(rows.Err(), "reading columns")
}

func readConstraints(ctx context.Context, pool *db.Pool, schemas []string, tables map[tableKey]*model.Table) error {
	rows, err := pool.Query(ctx, queryConstraints, schemas)
	if err != nil {
		return fmt.Errorf("reading constraints: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			schema, table, name, kind string
			validated                 bool
			refSchema, refTable       string
			columns, refColumns       []string
		)
		if err := rows.Scan(&schema, &table, &name, &kind, &validated,
			&refSchema, &refTable, &columns, &refColumns); err != nil {
			return fmt.Errorf("scanning constraint: %w", err)
		}

		t, ok := tables[key(schema, table)]
		if !ok {
			continue
		}

		switch kind {
		case "p":
			t.PrimaryKey = columns
		case "u":
			t.Uniques = append(t.Uniques, columns)
		case "f":
			t.ForeignKeys = append(t.ForeignKeys, model.ForeignKey{
				Name:       name,
				Columns:    columns,
				RefSchema:  refSchema,
				RefTable:   refTable,
				RefColumns: refColumns,
				Validated:  validated,
			})
		}
	}
	return wrap(rows.Err(), "reading constraints")
}

func readIndexes(ctx context.Context, pool *db.Pool, schemas []string, tables map[tableKey]*model.Table) error {
	rows, err := pool.Query(ctx, queryIndexes, schemas)
	if err != nil {
		return fmt.Errorf("reading indexes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schema, table string
		var idx model.Index
		if err := rows.Scan(&schema, &table, &idx.Name, &idx.Unique, &idx.Primary, &idx.Columns); err != nil {
			return fmt.Errorf("scanning index: %w", err)
		}
		if t, ok := tables[key(schema, table)]; ok {
			t.Indexes = append(t.Indexes, idx)
		}
	}
	return wrap(rows.Err(), "reading indexes")
}

func readUsage(ctx context.Context, pool *db.Pool, schemas []string, tables map[tableKey]*model.Table, coverage *model.Coverage) error {
	var resetAt *time.Time
	if err := pool.QueryRow(ctx, queryStatsReset).Scan(&resetAt); err != nil {
		return fmt.Errorf("reading statistics reset time: %w", err)
	}
	coverage.StatsResetAt = resetAt

	rows, err := pool.Query(ctx, queryUsageStats, schemas)
	if err != nil {
		return fmt.Errorf("reading usage statistics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schema, table string
		var c model.UsageCounters
		if err := rows.Scan(&schema, &table, &c.SeqScans, &c.IdxScans,
			&c.Inserts, &c.Updates, &c.Deletes); err != nil {
			return fmt.Errorf("scanning usage statistics: %w", err)
		}

		t, ok := tables[key(schema, table)]
		if !ok {
			continue
		}

		// Counters without their reset moment cannot support a conclusion, so
		// they are marked uninterpretable rather than quietly carried along.
		if resetAt == nil {
			t.Stats.Usage = model.NewUsageStatsUnknownReset(c)
		} else {
			t.Stats.Usage = model.NewUsageStats(c, *resetAt)
		}
	}
	return wrap(rows.Err(), "reading usage statistics")
}

// linkForeignKeyIndexes marks each foreign key as indexed only when some index
// carries the child column in leading position.
//
// A composite index holding the column anywhere else is useless for the lookup
// a parent DELETE triggers, so counting it as covered would produce a false
// negative — the tool would report all clear on a table where the problem is
// real. A redundant index suggestion costs a conversation; the false negative
// costs the incident the tool existed to prevent.
func linkForeignKeyIndexes(tables map[tableKey]*model.Table) {
	for _, t := range tables {
		for i, fk := range t.ForeignKeys {
			if len(fk.Columns) == 0 {
				continue
			}
			t.ForeignKeys[i].HasIndex = t.IsIndexedLeading(fk.Columns[0])
		}
	}
}

// classifyUnsupported records the shapes single-column inference cannot target.
// They are reported rather than passed over: a table skipped without a word
// looks identical to a table with nothing to find.
func classifyUnsupported(tables map[tableKey]*model.Table, coverage *model.Coverage) {
	for k, t := range tables {
		reason := unsupportedReason(t)
		if reason == "" {
			coverage.TablesAnalyzed++
			continue
		}
		coverage.TablesUnsupported = append(coverage.TablesUnsupported, model.SkippedTable{
			Table:  k.schema + "." + k.name,
			Reason: reason,
		})
	}

	sort.Slice(coverage.TablesUnsupported, func(i, j int) bool {
		return coverage.TablesUnsupported[i].Table < coverage.TablesUnsupported[j].Table
	})
}

func unsupportedReason(t *model.Table) model.UnsupportedReason {
	switch {
	case t.Partitioned:
		return model.ReasonPartitioned
	case t.Inherits:
		return model.ReasonInheritance
	case len(t.PrimaryKey) == 0:
		return model.ReasonNoPrimaryKey
	case len(t.PrimaryKey) > 1:
		return model.ReasonCompositePK
	default:
		return ""
	}
}

func group(schemas []string, tables map[tableKey]*model.Table) []model.Schema {
	byName := make(map[string][]model.Table, len(schemas))
	for _, t := range tables {
		byName[t.Schema] = append(byName[t.Schema], *t)
	}

	out := make([]model.Schema, 0, len(byName))
	for _, name := range schemas {
		list := byName[name]
		sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
		out = append(out, model.Schema{Name: name, Tables: list})
	}
	return out
}

// matchesAny reports whether a table matches an exclusion pattern, tried
// against both the bare name and the schema-qualified one.
func matchesAny(patterns []string, schema, table string) bool {
	qualified := schema + "." + table
	for _, p := range patterns {
		if ok, err := path.Match(p, table); err == nil && ok {
			return true
		}
		if ok, err := path.Match(p, qualified); err == nil && ok {
			return true
		}
	}
	return false
}

func wrap(err error, what string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", what, err)
}

// SortedSchemas returns the schema names in scope, normalized and deduplicated.
func SortedSchemas(raw []string) []string {
	seen := make(map[string]struct{}, len(raw))
	out := make([]string, 0, len(raw))
	for _, s := range raw {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, dup := seen[s]; dup {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
