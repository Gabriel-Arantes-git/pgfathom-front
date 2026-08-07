//go:build integration

package stats_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/catalog"
	"github.com/lvcas-dotcom/pgfathom/internal/db"
	"github.com/lvcas-dotcom/pgfathom/internal/infer"
	"github.com/lvcas-dotcom/pgfathom/internal/model"
	"github.com/lvcas-dotcom/pgfathom/internal/profile"
	"github.com/lvcas-dotcom/pgfathom/internal/stats"
	"github.com/lvcas-dotcom/pgfathom/internal/testutil"
)

// prefilterFromFixture runs catalog → inference → prefilter against a real
// server, which is the only honest way to test that real pg_stats rows drive
// the checks the way the unit fixtures say they do.
func prefilterFromFixture(t *testing.T) *stats.Result {
	t.Helper()

	ctx := context.Background()
	cfg := db.DefaultConfig()
	cfg.DSN = testutil.Postgres(t, "stats_prefilter")

	pool, err := db.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("opening pool: %v", err)
	}
	t.Cleanup(pool.Close)

	cat, err := catalog.Read(ctx, pool, catalog.Options{Schemas: []string{"public"}})
	if err != nil {
		t.Fatalf("reading catalog: %v", err)
	}

	naming, err := profile.Embedded("pt-br")
	if err != nil {
		t.Fatalf("loading profile: %v", err)
	}

	inferred := infer.Generate(cat.Schemas, infer.Options{Profile: naming})

	res, err := stats.Prefilter(ctx, pool, cat.Schemas, inferred.Candidates, stats.Options{})
	if err != nil {
		t.Fatalf("running prefilter: %v", err)
	}
	return res
}

func lookup(list []model.Candidate, child string) (model.Candidate, bool) {
	for _, c := range list {
		if c.Child.String() == child {
			return c, true
		}
	}
	return model.Candidate{}, false
}

// TestPrefilterReducesTheFunnel is the phase's deliverable: a measurable
// reduction in what reaches validation, on a fixture built to be reduced.
func TestPrefilterReducesTheFunnel(t *testing.T) {
	res := prefilterFromFixture(t)

	if len(res.Kept) >= res.Checked {
		t.Errorf("kept %d of %d checked: the impossible candidate should have been dropped",
			len(res.Kept), res.Checked)
	}

	rejected, ok := lookup(res.Rejected, "public.leitura.unidade_id")
	if !ok {
		t.Fatal("300 distinct values against a 3-row parent must be rejected")
	}
	if rejected.Verdict != model.VerdictRejected {
		t.Errorf("verdict = %q, want rejected", rejected.Verdict)
	}
	if !strings.Contains(rejected.Reason, "estimated") {
		t.Errorf("reason %q must declare the estimates as estimates", rejected.Reason)
	}
}

func TestRangeShiftPenalizesWithoutRejecting(t *testing.T) {
	res := prefilterFromFixture(t)

	c, ok := lookup(res.Kept, "public.medicao.ponto_id")
	if !ok {
		t.Fatal("the range-shifted candidate must survive: range never rejects on its own")
	}
	if !c.HasSignal(model.SigRangeViolation) {
		t.Error("a fully shifted histogram must leave the range-violation signal")
	}
}

func TestLegitimateCandidatePassesClean(t *testing.T) {
	res := prefilterFromFixture(t)

	c, ok := lookup(res.Kept, "public.pedido.cliente_id")
	if !ok {
		t.Fatal("the contained relationship must survive the prefilter")
	}
	for _, kind := range []model.SignalKind{
		model.SigCardViolation, model.SigRangeViolation, model.SigStatsUnavailable,
	} {
		if c.HasSignal(kind) {
			t.Errorf("a clean candidate must not carry %s", kind)
		}
	}
}

func TestUnanalyzedTablePassesWithARecord(t *testing.T) {
	res := prefilterFromFixture(t)

	if res.NoStats == 0 {
		t.Error("the never-ANALYZEd table must be counted as not evaluable")
	}

	c, ok := lookup(res.Kept, "public.sem_estatistica.unidade_id")
	if !ok {
		t.Fatal("a candidate without statistics must pass through, not disappear")
	}
	if !c.HasSignal(model.SigStatsUnavailable) {
		t.Error("the silence must be recorded on the candidate")
	}
}

// TestNoStatsValueEscapesTheLayer serializes everything the layer produced and
// scans it against what the fixture planted — both the text values and the
// out-of-range ids standing in for histogram endpoints.
func TestNoStatsValueEscapesTheLayer(t *testing.T) {
	res := prefilterFromFixture(t)

	out, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("serializing the result: %v", err)
	}

	for _, planted := range []string{
		"529.318.470-11",
		"Maria Aparecida Silva",
		"Leitura Manual Bloco C",
		"9001", "9300", // medicao.ponto_id histogram endpoints
		"1001", "1300", // leitura.unidade_id histogram endpoints
	} {
		if strings.Contains(string(out), planted) {
			t.Errorf("the serialized result leaked %q", planted)
		}
	}
}
