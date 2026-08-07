package model

// NamingEvidence is one naming convention detected in a schema, with what
// supports it.
type NamingEvidence struct {
	Affix       string  `json:"affix"`
	Occurrences int     `json:"occurrences"`
	Share       float64 `json:"share"`
}

// NamingDetection is what a schema revealed about its own naming convention.
//
// It belongs to the result rather than to the detector: without it a run is not
// reproducible, because the reader cannot tell which conventions were in effect.
type NamingDetection struct {
	Enabled bool `json:"enabled"`

	ColumnSuffixes []NamingEvidence `json:"column_suffixes,omitempty"`
	ColumnPrefixes []NamingEvidence `json:"column_prefixes,omitempty"`
	TablePrefixes  []NamingEvidence `json:"table_prefixes,omitempty"`

	// DeclaredKeys is how many declared foreign keys the reference affixes were
	// read from. Zero means the schema had nothing to say.
	DeclaredKeys int `json:"declared_keys"`
	Tables       int `json:"tables"`
}

// Empty reports whether nothing was detected.
func (d NamingDetection) Empty() bool {
	return len(d.ColumnSuffixes) == 0 && len(d.ColumnPrefixes) == 0 && len(d.TablePrefixes) == 0
}
