package model

import "time"

// UnsupportedReason says why a table could not be analyzed.
type UnsupportedReason string

// Shapes that single-column inference cannot target.
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

// Coverage describes what was actually analyzed against what exists. It is
// mandatory in every result: a clean report has to mean "I looked and it is
// clean", never "I could not look".
type Coverage struct {
	TablesTotal    int `json:"tables_total"`
	TablesAnalyzed int `json:"tables_analyzed"`

	// TablesNoPrivilege lists tables the connecting role cannot SELECT. A
	// partial grant is the common case in the target environment, not the edge
	// case, so silence here would make an incomplete report look clean.
	TablesNoPrivilege []string `json:"tables_no_privilege,omitempty"`

	// TablesExcluded lists tables removed by user-supplied filters.
	TablesExcluded []string `json:"tables_excluded,omitempty"`

	// TablesUnsupported lists tables whose shape this version cannot analyze.
	TablesUnsupported []SkippedTable `json:"tables_unsupported,omitempty"`

	CandidatesFound     int `json:"candidates_found"`
	CandidatesValidated int `json:"candidates_validated"`
	CandidatesTimedOut  int `json:"candidates_timed_out"`

	// StatsResetAt is when server statistics were last reset. Nil means unknown.
	StatsResetAt *time.Time `json:"stats_reset_at"`

	// PgStatStatements reports whether the query log was available to mine.
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
