package model

import "time"

// SignalKind identifies where a piece of evidence about a candidate came from.
// The set is closed: every signal must be attributable to a known origin, so
// that a score can always be explained back to the facts that produced it.
type SignalKind string

// Name-based evidence: the column name, once stripped of reference affixes,
// matches a table name. Exact matches are worth more than ones obtained through
// aggressive normalization.
const (
	SigExactName      SignalKind = "exact_name"
	SigNormalizedName SignalKind = "normalized_name"
)

// Type evidence. An identical base type is the ideal case; merely compatible
// types are weaker.
const (
	SigIdenticalType  SignalKind = "identical_type"
	SigCompatibleType SignalKind = "compatible_type"
)

// Target evidence. Ambiguity — several tables answering to the same normalized
// name — is the main source of noise in name-based inference.
const (
	SigUniqueTarget    SignalKind = "unique_target"
	SigAmbiguousTarget SignalKind = "ambiguous_target"
)

// Catalog evidence. An existing index suggests somebody joins on the column; a
// comment naming the target entity is a human saying so outright.
const (
	SigChildIndexed   SignalKind = "child_indexed"
	SigCommentMention SignalKind = "comment_mention"
	SigNotNull        SignalKind = "not_null"
)

// Usage evidence: a join the real code actually executes, extracted from a view
// definition, a function body or the query log.
//
// This is the strongest signal available, and the only one that finds
// relationships whose column names bear no resemblance to the target table —
// the class of relationship that name matching cannot reach by construction.
const (
	SigJoinInView       SignalKind = "join_in_view"
	SigJoinInFunction   SignalKind = "join_in_function"
	SigJoinInStatements SignalKind = "join_in_statements"
)

// Penalties. Generic domain names match everything and mean little; the
// statistics violations come from the planner and rule out impossibilities.
const (
	SigGenericDomain  SignalKind = "generic_domain_name"
	SigRangeViolation SignalKind = "stats_range_violation"
	SigCardViolation  SignalKind = "stats_cardinality_violation"
)

// Signal is one piece of evidence for or against a candidate.
type Signal struct {
	Kind   SignalKind `json:"kind"`
	Weight float64    `json:"weight"`

	// Detail carries catalog object names only — a view name, an index name, a
	// column name. Never a value read from a user table.
	Detail string `json:"detail,omitempty"`
}

// Verdict is the conclusion about a candidate. The set is closed.
type Verdict string

const (
	// VerdictConfirmed is a forgotten foreign key: total containment, no orphans.
	VerdictConfirmed Verdict = "confirmed"

	// VerdictBroken is a real relationship with broken integrity. The most
	// valuable finding this tool produces: a data bug that has been in
	// production for years without anyone knowing.
	VerdictBroken Verdict = "broken"

	// VerdictWeak means the data cannot support a conclusion — a single
	// distinct value, an overwhelming null fraction, an empty child table.
	VerdictWeak Verdict = "weak"

	// VerdictRejected means the evidence knocked the hypothesis down. The names
	// coincided; the data did not.
	VerdictRejected Verdict = "rejected"

	// VerdictUnvalidated means no evidence was gathered — a timeout, a missing
	// privilege, an unsupported shape.
	//
	// This is NOT rejection, and conflating the two is a way of reporting
	// silence as a clean bill of health. Rejected means the evidence said no.
	// Unvalidated means there was no evidence.
	VerdictUnvalidated Verdict = "unvalidated"
)

// Candidate is a possible relationship between two columns, raised by naming
// heuristics, by usage evidence, or by both, together with everything learned
// about it.
type Candidate struct {
	Child  ColumnRef `json:"child"`
	Parent ColumnRef `json:"parent"`

	// Signals are the evidence that produced MetaScore, kept so that any score
	// can be explained rather than merely asserted.
	Signals []Signal `json:"signals"`

	// MetaScore is confidence from metadata alone, 0..1, before any data is read.
	MetaScore float64 `json:"meta_score"`

	// Validation is nil until the candidate has been checked against the data.
	Validation *Validation `json:"validation,omitempty"`

	Verdict Verdict `json:"verdict"`

	// Reason explains a non-obvious verdict: why validation was skipped, why a
	// candidate was dropped. Object names and conditions only, never values.
	Reason string `json:"reason,omitempty"`
}

// HasSignal reports whether the candidate carries a signal of the given kind.
func (c Candidate) HasSignal(k SignalKind) bool {
	for _, s := range c.Signals {
		if s.Kind == k {
			return true
		}
	}
	return false
}

// EvidencedByUsage reports whether a real join was found in a view, a function
// body or the query log. Such a candidate stands on proof that the application
// treats the columns as related, independently of what they are named.
func (c Candidate) EvidencedByUsage() bool {
	return c.HasSignal(SigJoinInView) ||
		c.HasSignal(SigJoinInFunction) ||
		c.HasSignal(SigJoinInStatements)
}

// ValidationMethod distinguishes a sampled run from a complete one.
type ValidationMethod string

const (
	// MethodFull examined every row. Slow, and the only conclusive mode.
	MethodFull ValidationMethod = "full"

	// MethodSampled examined a subset. Fast, indicative, and structurally weak
	// against orphans: orphan rows arrive in batches and end up clustered on
	// the same pages, which is exactly what page-level sampling is worst at
	// finding. A sampled run is triage, never a conclusion.
	MethodSampled ValidationMethod = "sampled"
)

// Validation is the outcome of checking a candidate against the data.
type Validation struct {
	Method      ValidationMethod `json:"method"`
	SampledRows int64            `json:"sampled_rows,omitempty"`

	// NotNullRows is the number of non-NULL rows examined in the child column.
	NotNullRows int64 `json:"not_null_rows"`

	// DistinctVals is the number of distinct non-NULL values examined.
	DistinctVals int64 `json:"distinct_vals"`

	// OrphanRows and OrphanVals count the same failure in two units. A single
	// bad value repeated a million times and a million rare bad values are
	// different problems, and neither number implies the other.
	OrphanRows int64 `json:"orphan_rows"`
	OrphanVals int64 `json:"orphan_vals"`

	// MaxRowsPerValue distinguishes cardinality: 1 suggests one-to-one, well
	// above 1 suggests one-to-many.
	MaxRowsPerValue int64 `json:"max_rows_per_value"`

	Duration time.Duration `json:"duration_ns"`
}

// ContainmentRows is the proportion of examined rows whose value exists in the
// parent key. Returns 0 when nothing was examined.
func (v Validation) ContainmentRows() float64 {
	if v.NotNullRows <= 0 {
		return 0
	}
	return float64(v.NotNullRows-v.OrphanRows) / float64(v.NotNullRows)
}

// ContainmentVals is the proportion of distinct examined values that exist in
// the parent key. Returns 0 when nothing was examined.
func (v Validation) ContainmentVals() float64 {
	if v.DistinctVals <= 0 {
		return 0
	}
	return float64(v.DistinctVals-v.OrphanVals) / float64(v.DistinctVals)
}

// Conclusive reports whether this validation can support a confirmation. A
// sampled run never can: it cannot prove the absence of an orphan.
func (v Validation) Conclusive() bool { return v.Method == MethodFull }
