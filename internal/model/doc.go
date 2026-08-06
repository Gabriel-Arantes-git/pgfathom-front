// Package model holds the semantic model of an analyzed PostgreSQL schema.
//
// Two invariants shape the types here, and both are enforced by tests in this
// package rather than left to review.
//
// The package is pure: no I/O, no connections, no imports from other layers.
// The inference layer is trivial to test precisely because the types it
// operates on cannot reach a database.
//
// No exported or serializable field carries a value read from a user table.
// The model transports counts, ratios, timestamps and catalog object names.
// Column statistics are the dangerous case — most_common_vals and
// histogram_bounds from pg_stats are user data — so they live in unexported
// fields that encoding/json cannot reach.
package model
