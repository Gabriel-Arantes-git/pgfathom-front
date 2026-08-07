package model

// FindingKind identifies a structural finding that requires no inference:
// straight from the catalog, deterministic, immune to false positives.
type FindingKind string

const (
	// FindingNotValidConstraint is a constraint created NOT VALID and never
	// validated. See ForeignKey.Validated.
	FindingNotValidConstraint FindingKind = "not_valid_constraint"

	// FindingFKWithoutIndex is a declared foreign key with no usable index on
	// the child side.
	FindingFKWithoutIndex FindingKind = "fk_without_index"

	// FindingOrphanReference is a reference column pointing at a table that no
	// longer exists. Correctly rejected by validation, but a finding in itself.
	FindingOrphanReference FindingKind = "orphan_reference"

	// FindingPolymorphicPair is a reference column that only makes sense beside
	// a discriminator column. Validation would reject it for low containment,
	// which reads as "name coincidence" where a real relationship exists that
	// this version does not model.
	FindingPolymorphicPair FindingKind = "polymorphic_pair"

	// FindingUnsupportedTarget is a name that matched a table this version
	// cannot point at, typically because its key is composite.
	FindingUnsupportedTarget FindingKind = "unsupported_target"
)

// Finding is a structural observation that did not require inference.
type Finding struct {
	Kind FindingKind `json:"kind"`

	// Object is the schema-qualified catalog object the finding is about.
	Object string `json:"object"`

	// Detail describes the finding. Object names and conditions only.
	Detail string `json:"detail,omitempty"`

	Metrics map[string]int64 `json:"metrics,omitempty"`
}
