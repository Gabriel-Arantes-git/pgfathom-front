package stats_test

import (
	"strings"
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/model"
	"github.com/lvcas-dotcom/pgfathom/internal/stats"
)

func ref(table, column string) model.ColumnRef {
	return model.ColumnRef{Schema: "public", Table: table, Column: column}
}

func table(name string, rows int64) model.Table {
	return model.Table{
		Schema:     "public",
		Name:       name,
		Columns:    []model.Column{{Name: "id", Type: "bigint", BaseType: "int8"}},
		PrimaryKey: []string{"id"},
		Stats:      model.TableStats{EstimatedRows: rows},
	}
}

func schemaOf(tables ...model.Table) []model.Schema {
	return []model.Schema{{Name: "public", Tables: tables}}
}

// candidate builds a scored survivor the way generation would hand it over.
func candidate(child, parent model.ColumnRef, score float64) model.Candidate {
	return model.Candidate{
		Child:     child,
		Parent:    parent,
		Signals:   []model.Signal{{Kind: model.SigExactName, Weight: score}},
		MetaScore: score,
		Verdict:   model.VerdictUnvalidated,
	}
}

func present(nDistinct float64) model.ColumnStats {
	return model.ColumnStats{NDistinct: nDistinct, Present: true}
}

func TestViolationWithinMarginPenalizesWithoutRejecting(t *testing.T) {
	child, parent := ref("pedido", "unidade_id"), ref("unidade", "id")
	schemas := schemaOf(table("pedido", 10_000), table("unidade", 100))

	st := stats.NewStats()
	st.Add(child, present(150)) // 150 distinct vs 100 parent rows: violated, under 2x

	res := stats.Evaluate(schemas, []model.Candidate{candidate(child, parent, 0.8)}, st, stats.Options{})

	if len(res.Rejected) != 0 || len(res.Kept) != 1 {
		t.Fatalf("kept=%d rejected=%d, want the candidate kept", len(res.Kept), len(res.Rejected))
	}
	kept := res.Kept[0]
	if !kept.HasSignal(model.SigCardViolation) {
		t.Error("a violation inside the margin must leave the penalty signal")
	}
	if kept.MetaScore >= 0.8 {
		t.Errorf("score = %.2f, want it reduced by the penalty", kept.MetaScore)
	}
	if kept.Verdict != model.VerdictUnvalidated {
		t.Errorf("verdict = %q, want unvalidated: a penalty is not a conclusion", kept.Verdict)
	}
}

func TestViolationBeyondMarginRejectsWithReason(t *testing.T) {
	child, parent := ref("pedido", "unidade_id"), ref("unidade", "id")
	schemas := schemaOf(table("pedido", 10_000), table("unidade", 100))

	st := stats.NewStats()
	st.Add(child, present(250)) // 250 distinct vs 100 parent rows: beyond 2x

	res := stats.Evaluate(schemas, []model.Candidate{candidate(child, parent, 0.8)}, st, stats.Options{})

	if len(res.Rejected) != 1 {
		t.Fatalf("rejected=%d, want 1", len(res.Rejected))
	}
	rej := res.Rejected[0]
	if rej.Verdict != model.VerdictRejected {
		t.Errorf("verdict = %q, want rejected", rej.Verdict)
	}
	for _, want := range []string{"public.pedido.unidade_id", "public.unidade", "estimated"} {
		if !strings.Contains(rej.Reason, want) {
			t.Errorf("reason %q should mention %q: the estimates must be named as estimates", rej.Reason, want)
		}
	}
}

// TestNegativeNDistinctResolvesAgainstChildRows pins the pg_stats convention:
// negative n_distinct is a ratio of the row count, and misreading it as an
// absolute would let huge columns through.
func TestNegativeNDistinctResolvesAgainstChildRows(t *testing.T) {
	child, parent := ref("pedido", "unidade_id"), ref("unidade", "id")
	schemas := schemaOf(table("pedido", 1000), table("unidade", 100))

	st := stats.NewStats()
	st.Add(child, present(-0.5)) // half the 1000 rows: 500 distinct, beyond 2x100

	res := stats.Evaluate(schemas, []model.Candidate{candidate(child, parent, 0.8)}, st, stats.Options{})
	if len(res.Rejected) != 1 {
		t.Fatalf("rejected=%d, want 1: 500 estimated distinct against 100 rows", len(res.Rejected))
	}
}

func TestDisjointRangePenalizesAndNeverRejects(t *testing.T) {
	child, parent := ref("pedido", "unidade_id"), ref("unidade", "id")
	schemas := schemaOf(table("pedido", 10_000), table("unidade", 100))

	st := stats.NewStats()
	st.Add(child, present(50)) // cardinality fine
	st.Add(parent, present(100))
	st.AddBounds(child, 9000, 9999)
	st.AddBounds(parent, 1, 100)

	res := stats.Evaluate(schemas, []model.Candidate{candidate(child, parent, 0.8)}, st, stats.Options{})

	if len(res.Kept) != 1 || len(res.Rejected) != 0 {
		t.Fatalf("kept=%d rejected=%d, want kept: range never rejects on its own", len(res.Kept), len(res.Rejected))
	}
	kept := res.Kept[0]
	if !kept.HasSignal(model.SigRangeViolation) {
		t.Error("disjoint bounds must leave the range-violation signal")
	}
	if kept.MetaScore >= 0.8 {
		t.Errorf("score = %.2f, want it reduced", kept.MetaScore)
	}
}

func TestOverlappingRangeStaysQuiet(t *testing.T) {
	child, parent := ref("pedido", "unidade_id"), ref("unidade", "id")
	schemas := schemaOf(table("pedido", 10_000), table("unidade", 100))

	st := stats.NewStats()
	st.Add(child, present(50))
	st.Add(parent, present(100))
	st.AddBounds(child, 40, 120) // partially outside: not a violation
	st.AddBounds(parent, 1, 100)

	res := stats.Evaluate(schemas, []model.Candidate{candidate(child, parent, 0.8)}, st, stats.Options{})
	if res.Kept[0].HasSignal(model.SigRangeViolation) {
		t.Error("overlapping bounds must not be penalized: growth after ANALYZE looks exactly like this")
	}
}

func TestMissingStatsPassThroughWithARecord(t *testing.T) {
	child, parent := ref("pedido", "unidade_id"), ref("unidade", "id")
	schemas := schemaOf(table("pedido", 10_000), table("unidade", 100))

	res := stats.Evaluate(schemas, []model.Candidate{candidate(child, parent, 0.8)}, stats.NewStats(), stats.Options{})

	if res.NoStats != 1 {
		t.Fatalf("NoStats = %d, want 1", res.NoStats)
	}
	kept := res.Kept[0]
	if kept.MetaScore != 0.8 {
		t.Errorf("score = %.2f, want untouched: no opinion without statistics", kept.MetaScore)
	}
	if !kept.HasSignal(model.SigStatsUnavailable) {
		t.Error("the silence must be recorded on the candidate itself")
	}
	for _, s := range kept.Signals {
		if s.Kind == model.SigStatsUnavailable && s.Weight != 0 {
			t.Errorf("the marker signal weighs %.2f, want zero", s.Weight)
		}
	}
}

// TestUnanalyzedParentGivesNoOpinion pins the reltuples sentinel: -1 means
// unknown, and treating it as zero rows would reject every candidate whose
// parent was never ANALYZEd.
func TestUnanalyzedParentGivesNoOpinion(t *testing.T) {
	child, parent := ref("pedido", "unidade_id"), ref("unidade", "id")
	schemas := schemaOf(table("pedido", 10_000), table("unidade", model.RowsUnknown))

	st := stats.NewStats()
	st.Add(child, present(500))

	res := stats.Evaluate(schemas, []model.Candidate{candidate(child, parent, 0.8)}, st, stats.Options{})
	if len(res.Rejected) != 0 {
		t.Fatal("an unanalyzed parent must not produce a rejection")
	}
	if res.NoStats != 1 {
		t.Errorf("NoStats = %d, want 1", res.NoStats)
	}
}

func TestFunnelBalances(t *testing.T) {
	child1, child2, child3 := ref("a", "x_id"), ref("b", "x_id"), ref("c", "x_id")
	parent := ref("x", "id")
	schemas := schemaOf(table("a", 1000), table("b", 1000), table("c", 1000), table("x", 100))

	st := stats.NewStats()
	st.Add(child1, present(250)) // rejected
	st.Add(child2, present(50))  // kept clean
	// child3: no stats

	res := stats.Evaluate(schemas, []model.Candidate{
		candidate(child1, parent, 0.8),
		candidate(child2, parent, 0.8),
		candidate(child3, parent, 0.8),
	}, st, stats.Options{})

	if got := len(res.Kept) + len(res.Rejected); got != res.Checked {
		t.Errorf("kept(%d) + rejected(%d) = %d, want checked = %d",
			len(res.Kept), len(res.Rejected), got, res.Checked)
	}
	if res.NoStats != 1 || len(res.Rejected) != 1 {
		t.Errorf("funnel = rejected %d, nostats %d; want 1 and 1", len(res.Rejected), res.NoStats)
	}
}
