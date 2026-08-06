package model

import "time"

// SchemaVersion is the version of the JSON contract this build emits.
//
// The serialized model is the integration point with the CI baseline, with
// later phases and with third-party tooling. It is a public API: any
// incompatible change to field names or meanings requires incrementing this.
const SchemaVersion = "1"

// Result is a complete analysis run.
//
// Coverage is a required field, not an optional decoration. There is no
// constructor that produces a Result without it.
type Result struct {
	SchemaVersion string    `json:"schema_version"`
	Tool          string    `json:"tool"`
	ToolVersion   string    `json:"tool_version"`
	GeneratedAt   time.Time `json:"generated_at"`

	// Profile is the naming profile in effect. An inference result is not
	// interpretable without knowing which naming convention produced it.
	Profile string `json:"profile"`

	Schemas    []Schema    `json:"schemas"`
	Candidates []Candidate `json:"candidates,omitempty"`
	Findings   []Finding   `json:"findings,omitempty"`

	Coverage Coverage `json:"coverage"`
}

// NewResult starts a result with the contract fields filled in. Coverage is
// required at construction so that no code path can produce a result that
// reports nothing about what it failed to look at.
func NewResult(toolVersion, profile string, generatedAt time.Time, coverage Coverage) *Result {
	return &Result{
		SchemaVersion: SchemaVersion,
		Tool:          "pgfathom",
		ToolVersion:   toolVersion,
		GeneratedAt:   generatedAt,
		Profile:       profile,
		Coverage:      coverage,
	}
}

// CandidatesByVerdict returns the candidates carrying the given verdict, in the
// order they were added.
func (r *Result) CandidatesByVerdict(v Verdict) []Candidate {
	var out []Candidate
	for _, c := range r.Candidates {
		if c.Verdict == v {
			out = append(out, c)
		}
	}
	return out
}

// CountByVerdict tallies candidates per verdict.
func (r *Result) CountByVerdict() map[Verdict]int {
	out := make(map[Verdict]int, 5)
	for _, c := range r.Candidates {
		out[c.Verdict]++
	}
	return out
}

// Sampled reports whether any validation in this run was sampled, which means
// no result in it can be treated as confirmed.
func (r *Result) Sampled() bool {
	for _, c := range r.Candidates {
		if c.Validation != nil && c.Validation.Method == MethodSampled {
			return true
		}
	}
	return false
}
