package stats_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/model"
	"github.com/lvcas-dotcom/pgfathom/internal/stats"
)

// TestNoBoundValueSurvivesSerialization is the enforcement of the project's
// hardest rule at this layer's boundary. Histogram endpoints are user data;
// they enter through AddBounds, produce signals, and must be unreachable from
// anything that serializes — including the rejection reasons, which carry
// counts and names but never values.
func TestNoBoundValueSurvivesSerialization(t *testing.T) {
	// Recognizable stand-ins for what real histogram endpoints hold: national
	// ID numbers, account codes. Chosen not to collide with any count the
	// evaluation writes into a reason string.
	const plantedLow, plantedHigh = 52931847011, 52931847999

	child, parent := ref("pedido", "unidade_id"), ref("unidade", "id")
	schemas := schemaOf(table("pedido", 10_000), table("unidade", 100))

	st := stats.NewStats()
	st.Add(child, present(250)) // beyond margin: forces the rejection path too
	st.Add(parent, present(100))
	st.AddBounds(child, plantedLow, plantedHigh)
	st.AddBounds(parent, plantedLow, plantedHigh)

	res := stats.Evaluate(schemas, []model.Candidate{candidate(child, parent, 0.8)}, st, stats.Options{})

	for name, v := range map[string]any{
		"result": res,
		"stats":  st,
	} {
		out, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("serializing %s: %v", name, err)
		}
		for _, planted := range []string{"52931847011", "52931847999"} {
			if strings.Contains(string(out), planted) {
				t.Errorf("%s serialization leaked histogram bound %s", name, planted)
			}
		}
	}

	for _, c := range append(append([]model.Candidate{}, res.Kept...), res.Rejected...) {
		for _, s := range c.Signals {
			for _, planted := range []string{"52931847011", "52931847999"} {
				if strings.Contains(s.Detail, planted) {
					t.Errorf("signal detail %q carries a histogram bound", s.Detail)
				}
			}
		}
		if strings.Contains(c.Reason, "52931847011") || strings.Contains(c.Reason, "52931847999") {
			t.Errorf("rejection reason %q carries a histogram bound", c.Reason)
		}
	}
}
