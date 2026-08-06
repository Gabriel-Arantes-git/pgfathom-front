package model

// FindingKind identifies a structural finding that requires no inference at all.
// These come straight from the catalog: deterministic, free, and immune to
// false positives.
type FindingKind string

const (
	// FindingNotValidConstraint is a constraint created NOT VALID and never
	// validated. It blocks new violations but never checked the rows that were
	// already there, while presenting itself as a normal constraint everywhere.
	FindingNotValidConstraint FindingKind = "not_valid_constraint"

	// FindingFKWithoutIndex is a declared foreign key with no usable index on
	// the child side, which turns every parent delete into a sequential scan.
	FindingFKWithoutIndex FindingKind = "fk_without_index"

	// FindingOrphanReference is a reference column pointing at a table that no
	// longer exists. Correctly rejected by validation, but a finding in itself.
	FindingOrphanReference FindingKind = "orphan_reference"
)

// Finding is a structural observation that did not require inference.
type Finding struct {
	Kind FindingKind `json:"kind"`

	// Object is the catalog object the finding is about, schema-qualified.
	Object string `json:"object"`

	// Detail describes the finding. Object names and conditions only, never a
	// value read from a user table.
	Detail string `json:"detail,omitempty"`

	// Metrics carries counts and sizes relevant to the finding.
	Metrics map[string]int64 `json:"metrics,omitempty"`
}
