package validate

import (
	"strings"
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/model"
)

func val(method model.ValidationMethod, sampled, notNull, distinct, orphanRows, orphanVals int64) model.Validation {
	return model.Validation{
		Method:       method,
		SampledRows:  sampled,
		NotNullRows:  notNull,
		DistinctVals: distinct,
		OrphanRows:   orphanRows,
		OrphanVals:   orphanVals,
	}
}

func TestVerdictRules(t *testing.T) {
	tests := []struct {
		name string
		v    model.Validation
		want model.Verdict
	}{
		{"full total containment confirms", val(model.MethodFull, 100, 100, 10, 0, 0), model.VerdictConfirmed},
		{"sampled clean never confirms", val(model.MethodSampled, 100, 100, 10, 0, 0), model.VerdictWeak},
		{"orphans above break threshold", val(model.MethodFull, 100, 100, 10, 6, 2), model.VerdictBroken},
		{"orphans in a sample still break", val(model.MethodSampled, 100, 100, 10, 6, 2), model.VerdictBroken},
		{"low containment rejects", val(model.MethodFull, 100, 100, 10, 70, 7), model.VerdictRejected},
		{"dead zone is a human call", val(model.MethodFull, 100, 100, 10, 30, 3), model.VerdictWeak},
		{"empty child is weak", val(model.MethodFull, 0, 0, 0, 0, 0), model.VerdictWeak},
		{"single distinct value is weak", val(model.MethodFull, 100, 100, 1, 0, 0), model.VerdictWeak},
		{"null-dominated column is weak", val(model.MethodFull, 10_000, 100, 10, 0, 0), model.VerdictWeak},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, reason := verdict(tt.v, Options{})
			if got != tt.want {
				t.Fatalf("verdict = %q (%s), want %q", got, reason, tt.want)
			}
			if got != model.VerdictConfirmed && got != model.VerdictBroken && reason == "" {
				t.Error("a non-obvious verdict must carry its reason")
			}
		})
	}
}

func TestSampledCleanPointsAtFull(t *testing.T) {
	_, reason := verdict(val(model.MethodSampled, 100, 100, 10, 0, 0), Options{})
	if !strings.Contains(reason, "--full") {
		t.Errorf("reason %q must point the user at --full", reason)
	}
}

func TestSampleSelection(t *testing.T) {
	tests := []struct {
		name  string
		rows  int64
		known bool
		opts  Options
		want  sampleMethod
	}{
		{"fits the target: read whole", 50_000, true, Options{}, sampleNone},
		{"unknown estimate: read whole", 0, false, Options{}, sampleNone},
		{"full overrides everything", 10_000_000, true, Options{Full: true}, sampleNone},
		{"small table samples by row", 50_000, true, Options{TargetRows: 1000}, sampleBernoulli},
		{"large table samples by page", 10_000_000, true, Options{}, sampleSystem},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sampleFor(tt.rows, tt.known, tt.opts); got.kind != tt.want {
				t.Errorf("sample kind = %v, want %v", got.kind, tt.want)
			}
		})
	}
}

func TestSampleFractionIsBounded(t *testing.T) {
	spec := sampleFor(10_000_000_000, true, Options{TargetRows: 100})
	if spec.percent < 0.01 {
		t.Errorf("percent = %f, want the 0.01 floor: a zero fraction samples nothing", spec.percent)
	}
}

func TestQueryQuotesExoticIdentifiers(t *testing.T) {
	c := model.Candidate{
		Child:  model.ColumnRef{Schema: "public", Table: "Ordem Servico", Column: "unidade_id"},
		Parent: model.ColumnRef{Schema: "public", Table: `uni"dade`, Column: "id"},
	}

	q := buildQuery(c, sampleSpec{kind: sampleNone})

	for _, want := range []string{`"Ordem Servico"`, `"uni""dade"`, `"unidade_id"`} {
		if !strings.Contains(q, want) {
			t.Errorf("query must quote %s; got:\n%s", want, q)
		}
	}
}

// TestSampleClauseFollowsTheAlias pins the from_item grammar: the alias comes
// before TABLESAMPLE, and getting it backwards is a syntax error only a real
// server would catch.
func TestSampleClauseFollowsTheAlias(t *testing.T) {
	c := model.Candidate{
		Child:  model.ColumnRef{Schema: "public", Table: "pedido", Column: "cliente_id"},
		Parent: model.ColumnRef{Schema: "public", Table: "cliente", Column: "id"},
	}

	q := buildQuery(c, sampleSpec{kind: sampleBernoulli, percent: 12.5, seed: 42})

	if !strings.Contains(q, `"pedido" AS c TABLESAMPLE BERNOULLI (12.5000) REPEATABLE (42)`) {
		t.Errorf("sample clause misplaced or malformed:\n%s", q)
	}
}
