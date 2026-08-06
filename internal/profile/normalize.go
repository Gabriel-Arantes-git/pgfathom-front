package profile

import "strings"

// Origin records how a candidate form was derived, which is what lets scoring
// tell an exact match from one obtained by aggressive normalization.
type Origin uint8

// Ordered from strongest to weakest match.
const (
	// OriginExact is the table name unchanged.
	OriginExact Origin = iota

	// OriginPrefixStripped had a legacy convention prefix removed.
	OriginPrefixStripped

	// OriginDepluralized had a plural rule applied.
	OriginDepluralized

	// OriginPrefixAndDepluralized had both, and is the most likely coincidence.
	OriginPrefixAndDepluralized
)

// String renders the origin for diagnostics.
func (o Origin) String() string {
	switch o {
	case OriginExact:
		return "exact"
	case OriginPrefixStripped:
		return "prefix_stripped"
	case OriginDepluralized:
		return "depluralized"
	case OriginPrefixAndDepluralized:
		return "prefix_stripped+depluralized"
	default:
		return "unknown"
	}
}

// Exact reports whether the form is the table name unchanged.
func (o Origin) Exact() bool { return o == OriginExact }

// Form is one candidate form of a table name, with the transformation that
// produced it.
type Form struct {
	Value  string
	Origin Origin
}

// EntityName derives the entity a column appears to reference. Suffixes are
// tried before prefixes, each in declared order with the first match winning.
// A name matching no affix is returned unchanged: "municipio" references
// municipio just as legitimately as "municipio_id" does.
//
// The result is lowercased, matching PostgreSQL's folding of unquoted
// identifiers.
func (p *Profile) EntityName(column string) string {
	name := strings.ToLower(strings.TrimSpace(column))
	if name == "" {
		return ""
	}

	for _, suffix := range p.ColumnSuffixes {
		if trimmed, ok := trimSuffix(name, suffix); ok {
			name = trimmed
			break
		}
	}

	for _, prefix := range p.ColumnPrefixes {
		if trimmed, ok := trimPrefix(name, prefix); ok {
			name = trimmed
			break
		}
	}

	return strings.Trim(name, "_")
}

// TableForms returns the ordered set of candidate forms for a table name, with
// the unchanged name always first so an exact match wins.
//
// It returns a set rather than one form because plural rules are genuinely
// ambiguous: "logins" yields "logim" under the ns→m rule, right for "armazens",
// and "login" under the generic rule, and nothing in the name says which. Under
// first-rule-wins, the chosen order would silently decide which case breaks.
//
// The price is the occasional spurious match through an aggressive form, which
// scoring penalizes and validation knocks down.
func (p *Profile) TableForms(table string) []Form {
	name := strings.ToLower(strings.TrimSpace(table))
	if name == "" {
		return nil
	}

	forms := make([]Form, 0, len(p.Plural)+2)
	seen := make(map[string]struct{}, len(p.Plural)+2)

	add := func(value string, origin Origin) {
		if value == "" {
			return
		}
		if _, dup := seen[value]; dup {
			return
		}
		seen[value] = struct{}{}
		forms = append(forms, Form{Value: value, Origin: origin})
	}

	add(name, OriginExact)

	// A table name carries at most one legacy prefix in practice.
	stripped := ""
	for _, prefix := range p.TablePrefixes {
		if trimmed, ok := trimPrefix(name, prefix); ok {
			stripped = trimmed
			add(stripped, OriginPrefixStripped)
			break
		}
	}

	for _, singular := range p.depluralize(name) {
		add(singular, OriginDepluralized)
	}
	if stripped != "" {
		for _, singular := range p.depluralize(stripped) {
			add(singular, OriginPrefixAndDepluralized)
		}
	}

	return forms
}

// depluralize applies every applicable rule and returns each form produced.
func (p *Profile) depluralize(name string) []string {
	out := make([]string, 0, 2)
	for _, rule := range p.Plural {
		if !strings.HasSuffix(name, rule.Suffix) {
			continue
		}
		stem := name[:len(name)-len(rule.Suffix)]
		if stem == "" {
			// The whole name is the suffix; stripping it says nothing.
			continue
		}
		out = append(out, stem+rule.Singular)
	}
	return out
}

// Match reports whether an entity name matches any form of a table name, and
// which form it was.
func (p *Profile) Match(entity, table string) (Form, bool) {
	entity = strings.ToLower(strings.TrimSpace(entity))
	if entity == "" {
		return Form{}, false
	}

	for _, form := range p.TableForms(table) {
		if form.Value == entity {
			return form, true
		}
	}
	return Form{}, false
}

// trimSuffix removes suffix from name, refusing to leave nothing behind.
func trimSuffix(name, suffix string) (string, bool) {
	if suffix == "" || !strings.HasSuffix(name, suffix) {
		return name, false
	}
	trimmed := name[:len(name)-len(suffix)]
	if strings.Trim(trimmed, "_") == "" {
		return name, false
	}
	return trimmed, true
}

// trimPrefix removes prefix from name, refusing to leave nothing behind.
func trimPrefix(name, prefix string) (string, bool) {
	if prefix == "" || !strings.HasPrefix(name, prefix) {
		return name, false
	}
	trimmed := name[len(prefix):]
	if strings.Trim(trimmed, "_") == "" {
		return name, false
	}
	return trimmed, true
}
