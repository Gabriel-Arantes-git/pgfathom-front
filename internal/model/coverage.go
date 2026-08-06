package model

import "time"

// UnsupportedReason says why a table could not be analyzed.
type UnsupportedReason string

// Shapes that single-column inference cannot target. Each must be reported
// rather than silently passed over: a table skipped without a word looks
// identical to a table with nothing to find.
const (
	ReasonCompositePK  UnsupportedReason = "composite_primary_key"
	ReasonNoPrimaryKey UnsupportedReason = "no_primary_key"
	ReasonPartitioned  UnsupportedReason = "partitioned"
	ReasonInheritance  UnsupportedReason = "table_inheritance"
)

// SkippedTable records one table that did not make it into the analysis, and why.
type SkippedTable struct {
	Table  string            `json:"table"`
	Reason UnsupportedReason `json:"reason,omitempty"`
}

// Coverage describes what was actually analyzed against what exists.
//
// It is mandatory in every result and is never assembled at rendering time.
// A clean report has to mean "I looked and it is clean", never "I could not
// look". Tables skipped for missing privileges, candidates that timed out and
// schemas left out are all findings about the analysis itself, and hiding them
// turns absence of evidence into evidence of absence.
type Coverage struct {
	TablesTotal    int `json:"tables_total"`
	TablesAnalyzed int `json:"tables_analyzed"`

	// TablesNoPrivilege lists tables the connecting role cannot SELECT. In the
	// target environment — public sector databases where obtaining privileges
	// is a political process — a partial grant is the common case, not the edge
	// case, and silence about it would make an incomplete report look clean.
	TablesNoPrivilege []string `json:"tables_no_privilege,omitempty"`

	// TablesExcluded lists tables removed by user-supplied filters.
	TablesExcluded []string `json:"tables_excluded,omitempty"`

	// TablesUnsupported lists tables whose shape this version cannot analyze.
	TablesUnsupported []SkippedTable `json:"tables_unsupported,omitempty"`

	CandidatesFound     int `json:"candidates_found"`
	CandidatesValidated int `json:"candidates_validated"`
	CandidatesTimedOut  int `json:"candidates_timed_out"`

	// StatsResetAt is when the server's statistics were last reset. Nil means
	// unknown, and every usage-based finding must be read in that light.
	StatsResetAt *time.Time `json:"stats_reset_at"`

	// PgStatStatements reports whether the extension was available to mine
	// join predicates from the query log.
	PgStatStatements bool `json:"pg_stat_statements"`
}

// Complete reports whether every table in scope was analyzed and every
// candidate resolved.
func (c Coverage) Complete() bool {
	return c.TablesAnalyzed == c.TablesTotal &&
		len(c.TablesNoPrivilege) == 0 &&
		c.CandidatesTimedOut == 0
}

// SkippedCount is how many tables did not make it into the analysis.
func (c Coverage) SkippedCount() int {
	return len(c.TablesNoPrivilege) + len(c.TablesExcluded) + len(c.TablesUnsupported)
}
