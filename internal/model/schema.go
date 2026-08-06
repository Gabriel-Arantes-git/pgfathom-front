package model

import "strings"

// Schema is a namespace of tables read from the catalog.
type Schema struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

// Table is a relation and everything the catalog says about it.
type Table struct {
	Schema  string   `json:"schema"`
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`

	// PrimaryKey lists column names in key order. Empty when there is none.
	PrimaryKey []string `json:"primary_key,omitempty"`

	Uniques [][]string `json:"uniques,omitempty"`

	// ForeignKeys holds DECLARED constraints only. Inferred relationships live
	// in Candidate, so no consumer can mistake one for the other.
	ForeignKeys []ForeignKey `json:"foreign_keys,omitempty"`

	Indexes []Index    `json:"indexes,omitempty"`
	Stats   TableStats `json:"stats"`

	// Partitioned marks a partitioned parent: read from it, never iterate the
	// partitions, and expect statistics to behave differently.
	Partitioned bool `json:"partitioned,omitempty"`

	Inherits bool `json:"inherits,omitempty"`

	Comment string `json:"comment,omitempty"`
}

// Ref returns the schema-qualified name.
func (t Table) Ref() string { return t.Schema + "." + t.Name }

// HasSingleColumnPK reports whether the primary key is exactly one column,
// the only shape single-column inference can target.
func (t Table) HasSingleColumnPK() bool { return len(t.PrimaryKey) == 1 }

// Column looks up a column by name. The second result is false when absent.
func (t Table) Column(name string) (Column, bool) {
	for _, c := range t.Columns {
		if strings.EqualFold(c.Name, name) {
			return c, true
		}
	}
	return Column{}, false
}

// IsIndexedLeading reports whether an index has the column in leading position,
// which is what makes it usable for a foreign key lookup.
func (t Table) IsIndexedLeading(column string) bool {
	for _, idx := range t.Indexes {
		if len(idx.Columns) > 0 && strings.EqualFold(idx.Columns[0], column) {
			return true
		}
	}
	return false
}

// Column is an attribute of a table.
type Column struct {
	Name string `json:"name"`

	// Type is the type as PostgreSQL formats it, e.g. "character varying(60)".
	Type string `json:"type"`

	// BaseType is normalized for comparison, e.g. "int8". Comparing formatted
	// types directly produces false negatives between equivalent spellings.
	BaseType string `json:"base_type"`

	Nullable bool   `json:"nullable"`
	Default  string `json:"default,omitempty"`
	Position int    `json:"position"`
	Comment  string `json:"comment,omitempty"`
}

// ColumnRef points at one column of one table.
type ColumnRef struct {
	Schema string `json:"schema"`
	Table  string `json:"table"`
	Column string `json:"column"`
}

// String renders the reference as schema.table.column.
func (r ColumnRef) String() string { return r.Schema + "." + r.Table + "." + r.Column }

// TableRef renders the owning table as schema.table.
func (r ColumnRef) TableRef() string { return r.Schema + "." + r.Table }

// ForeignKey is a constraint that exists in the catalog. That is not the same
// as a guarantee of integrity: see Validated.
type ForeignKey struct {
	Name       string   `json:"name"`
	Columns    []string `json:"columns"`
	RefSchema  string   `json:"ref_schema"`
	RefTable   string   `json:"ref_table"`
	RefColumns []string `json:"ref_columns"`

	// Validated mirrors pg_constraint.convalidated. False means NOT VALID: the
	// constraint blocks new violations but never checked the rows already
	// there, while looking identical in \d and in any ERD tool.
	Validated bool `json:"validated"`

	// HasIndex reports a usable index on the child side. Without one, every
	// parent delete becomes a sequential scan of the child.
	HasIndex bool `json:"has_index"`
}

// Index is an index as declared in the catalog.
type Index struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
	Primary bool     `json:"primary"`
}
