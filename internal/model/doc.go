// Package model holds the semantic model of an analyzed PostgreSQL schema.
//
// The package is pure: it performs no I/O, opens no connections, reads no files,
// and imports no other layer of this project. A test in this package enforces
// that by inspecting the import graph, because the dependency is easy to
// introduce by accident and expensive to remove afterwards.
//
// # Provenance
//
// The model keeps three sources of knowledge about a relationship structurally
// apart: declared in the catalog (ForeignKey), evidenced by usage in a view,
// function or query log (Signal), and inferred by heuristic (Candidate). There
// is no construction that lets an inferred relationship be read as a declared
// one.
//
// # User data never leaves memory
//
// No exported or serializable field of this package may carry a value read from
// a user table. What the model transports are counts, ratios, timestamps and
// catalog object names.
//
// Column statistics are the dangerous case: most_common_vals and
// histogram_bounds from pg_stats are user data. They live in unexported fields,
// which encoding/json cannot reach, and are meant to be discarded as soon as
// they have produced their metric. Free-text diagnostic fields carry object
// names only, never values.
package model
