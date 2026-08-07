package validate

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/lvcas-dotcom/pgfathom/internal/model"
)

type sampleMethod int

const (
	sampleNone sampleMethod = iota
	sampleBernoulli
	sampleSystem
)

// sampleSpec is how the child table will be read.
type sampleSpec struct {
	kind    sampleMethod
	percent float64
	seed    int64
}

func (s sampleSpec) method() model.ValidationMethod {
	if s.kind == sampleNone {
		return model.MethodFull
	}
	return model.MethodSampled
}

// sampleFor picks the reading mode. A table that fits the target is read
// whole, which returns the conclusive mode for free. Unknown estimates also
// read whole: guessing a fraction from a number the planner does not have
// would be inventing, and the time ceiling already bounds the cost.
func sampleFor(rows int64, known bool, opts Options) sampleSpec {
	target := opts.targetRows()
	if opts.Full || !known || rows <= target {
		return sampleSpec{kind: sampleNone}
	}

	pct := 100 * float64(target) / float64(rows)
	if pct < 0.01 {
		pct = 0.01
	}

	kind := sampleSystem
	if rows < bernoulliCeiling {
		kind = sampleBernoulli
	}
	return sampleSpec{kind: kind, percent: pct, seed: opts.Seed}
}

// buildQuery assembles the aggregation for one candidate. Identifiers are
// quoted without exception; the sampling fraction is computed here and
// rendered by strconv, never taken from user input.
//
// The scanned CTE is read once and feeds both the total examined — which is
// what makes NULL dominance measurable — and the per-value aggregation that
// turns the anti-join from row count into column cardinality.
func buildQuery(c model.Candidate, spec sampleSpec) string {
	child := pgx.Identifier{c.Child.Schema, c.Child.Table}.Sanitize()
	col := pgx.Identifier{c.Child.Column}.Sanitize()
	parent := pgx.Identifier{c.Parent.Schema, c.Parent.Table}.Sanitize()
	pk := pgx.Identifier{c.Parent.Column}.Sanitize()

	// The alias must precede TABLESAMPLE in the from_item grammar.
	from := child + " AS c" + sampleClause(spec)

	return fmt.Sprintf(`
		WITH scanned AS (
		    SELECT c.%[1]s AS v FROM %[2]s
		),
		child_vals AS (
		    SELECT v, count(*) AS n FROM scanned WHERE v IS NOT NULL GROUP BY 1
		),
		marked AS (
		    SELECT cv.n, EXISTS (SELECT 1 FROM %[3]s p WHERE p.%[4]s = cv.v) AS matched
		    FROM child_vals cv
		)
		SELECT
		    (SELECT count(*) FROM scanned)                            AS sampled_rows,
		    count(*)                                                  AS distinct_vals,
		    coalesce(sum(n), 0)::bigint                               AS not_null_rows,
		    count(*) FILTER (WHERE NOT matched)                       AS orphan_vals,
		    coalesce(sum(n) FILTER (WHERE NOT matched), 0)::bigint    AS orphan_rows,
		    coalesce(max(n), 0)::bigint                               AS max_rows_per_value
		FROM marked`,
		col, from, parent, pk)
}

func sampleClause(spec sampleSpec) string {
	var method string
	switch spec.kind {
	case sampleSystem:
		method = "SYSTEM"
	case sampleBernoulli:
		method = "BERNOULLI"
	default:
		return ""
	}

	clause := " TABLESAMPLE " + method +
		" (" + strconv.FormatFloat(spec.percent, 'f', 4, 64) + ")"
	if spec.seed != 0 {
		clause += " REPEATABLE (" + strconv.FormatInt(spec.seed, 10) + ")"
	}
	return clause
}
