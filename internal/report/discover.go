package report

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/lvcas-dotcom/pgfathom/internal/model"
)

// DiscoverView is everything the discover report needs, including what was
// thrown away.
type DiscoverView struct {
	Result *model.Result

	// Discarded fell below the score threshold. Shown only on request, but
	// always counted, so nobody has to wonder why an obvious-looking column was
	// ignored.
	Discarded []model.Candidate

	MinScore        float64
	ShowDiscarded   bool
	ValidationStage string

	// Detection is what the schema revealed about its own naming convention.
	// Reporting it is not decoration: a profile that changes itself without
	// saying so leaves the user unable to reproduce or to disagree.
	Detection model.NamingDetection
}

// Discover renders inferred candidates.
func Discover(w io.Writer, v DiscoverView) error {
	var b strings.Builder
	r := v.Result

	fmt.Fprintf(&b, "\n  %s %s · PostgreSQL %s · profile %s · threshold %.2f\n",
		r.Tool, r.ToolVersion, r.ServerVersion, r.Profile, v.MinScore)

	// A well-scored candidate is still a hypothesis. Presenting one without
	// this caveat would invite creating a constraint from a name match, which
	// is exactly the mistake the tool exists to prevent.
	fmt.Fprintf(&b, "  %s\n\n", v.ValidationStage)

	if len(r.Candidates) == 0 {
		fmt.Fprintf(&b, "  No candidate relationships above the threshold.\n\n")
	} else {
		writeCandidates(&b, r.Candidates)
	}

	if v.ShowDiscarded && len(v.Discarded) > 0 {
		writeDiscarded(&b, v.Discarded)
	}

	writeDetection(&b, v)
	writeObservations(&b, r.Findings)
	writeCoverage(&b, r.Coverage)

	if !v.ShowDiscarded && len(v.Discarded) > 0 {
		fmt.Fprintf(&b, "    %d candidates below the threshold — pass --include-rejected to see them\n\n",
			len(v.Discarded))
	}

	_, err := io.WriteString(w, b.String())
	return err
}

func writeCandidates(b *strings.Builder, candidates []model.Candidate) {
	fmt.Fprintf(b, "  CANDIDATES — inferred, not yet verified against data  (%d)\n", len(candidates))
	fmt.Fprintf(b, "  %s\n", strings.Repeat("─", 74))

	tw := tabwriter.NewWriter(b, 0, 0, 2, ' ', 0)
	for _, c := range candidates {
		_, _ = fmt.Fprintf(tw, "  %s\t→ %s\t%.2f\t%s\n",
			c.Child.String(), c.Parent.String(), c.MetaScore, summarize(c.Signals))
	}
	_ = tw.Flush()
	b.WriteString("\n")
}

func writeDiscarded(b *strings.Builder, discarded []model.Candidate) {
	fmt.Fprintf(b, "  BELOW THRESHOLD — generated and discarded  (%d)\n", len(discarded))
	fmt.Fprintf(b, "  %s\n", strings.Repeat("─", 74))

	tw := tabwriter.NewWriter(b, 0, 0, 2, ' ', 0)
	for _, c := range discarded {
		_, _ = fmt.Fprintf(tw, "  %s\t→ %s\t%.2f\t%s\n",
			c.Child.String(), c.Parent.String(), c.MetaScore, summarize(c.Signals))
	}
	_ = tw.Flush()
	b.WriteString("\n")
}

func writeObservations(b *strings.Builder, findings []model.Finding) {
	relevant := make([]model.Finding, 0, len(findings))
	for _, f := range findings {
		if f.Kind == model.FindingPolymorphicPair || f.Kind == model.FindingUnsupportedTarget {
			relevant = append(relevant, f)
		}
	}
	if len(relevant) == 0 {
		return
	}

	fmt.Fprintf(b, "  NOT ANALYZED — recognized, out of scope for this version  (%d)\n", len(relevant))
	fmt.Fprintf(b, "  %s\n", strings.Repeat("─", 74))

	tw := tabwriter.NewWriter(b, 0, 0, 2, ' ', 0)
	for _, f := range relevant {
		_, _ = fmt.Fprintf(tw, "  %s\t%s\n", f.Object, f.Detail)
	}
	_ = tw.Flush()
	b.WriteString("\n")
}

// summarize renders the signals compactly, strongest first, so a score can be
// argued with rather than merely accepted.
func summarize(signals []model.Signal) string {
	ordered := make([]model.Signal, len(signals))
	copy(ordered, signals)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Weight > ordered[j].Weight })

	parts := make([]string, 0, len(ordered))
	for _, s := range ordered {
		if s.Weight < 0 {
			parts = append(parts, "-"+string(s.Kind))
			continue
		}
		parts = append(parts, string(s.Kind))
	}
	return strings.Join(parts, " ")
}

func writeDetection(b *strings.Builder, v DiscoverView) {
	if !v.Detection.Enabled {
		fmt.Fprintf(b, "  naming detection is off — only the %s profile applies\n\n", v.Result.Profile)
		return
	}

	if v.Detection.Empty() {
		fmt.Fprintf(b, "  Nothing detected from %d tables and %d declared keys; the %s profile applies alone.\n\n",
			v.Detection.Tables, v.Detection.DeclaredKeys, v.Result.Profile)
		return
	}

	fmt.Fprintf(b, "  DETECTED — conventions read from the schema itself, added to %s\n", v.Result.Profile)
	fmt.Fprintf(b, "  %s\n", strings.Repeat("─", 74))

	tw := tabwriter.NewWriter(b, 0, 0, 2, ' ', 0)
	for _, group := range []struct {
		label string
		items []model.NamingEvidence
	}{
		{"reference suffix", v.Detection.ColumnSuffixes},
		{"reference prefix", v.Detection.ColumnPrefixes},
		{"table prefix", v.Detection.TablePrefixes},
	} {
		for _, e := range group.items {
			_, _ = fmt.Fprintf(tw, "  %s\t%s\t%d occurrences (%.0f%%)\n",
				group.label, e.Affix, e.Occurrences, 100*e.Share)
		}
	}
	_ = tw.Flush()
	b.WriteString("\n")
}
