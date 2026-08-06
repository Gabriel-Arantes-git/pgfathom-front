package profile

import "strings"

// Origin records how a candidate form was derived from a table name. It is what
// lets scoring tell an exact match from one obtained by aggressive
// normalization — the difference between the SigExactName and SigNormalizedName
// signals.
type Origin uint8

const (
	// OriginExact is the table name unchanged. The strongest match there is.
	OriginExact Origin = iota

	// OriginPrefixStripped had a legacy convention prefix removed.
	OriginPrefixStripped

	// OriginDepluralized had a plural rule applied.
	OriginDepluralized

	// OriginPrefixAndDepluralized had both. The weakest match, and the one most
	// likely to be a coincidence.
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

// EntityName derives the entity a column appears to reference, by stripping the
// reference affixes declared in the profile.
//
// Suffixes are tried before prefixes, each list in declared order with the first
// match winning. A column whose name matches no affix is returned unchanged: a
// column called "municipio" references municipio just as legitimately as
// "municipio_id" does.
//
// Comparison is case-insensitive and the result is lowercased, matching
// PostgreSQL's own folding of unquoted identifiers.
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

// TableForms returns the ordered set of candidate forms for a table name.
//
// It returns a set rather than a single form on purpose. Portuguese plural rules
// are genuinely ambiguous: "logins" yields "logim" under the ns→m rule, correct
// for "armazens", and "login" under the generic drop-the-s rule, and nothing in
// the name says which is right. Under first-rule-wins, whichever order was
// chosen would silently decide which of the two cases the project gets wrong.
//
// Returning every plausible form and letting the match decide costs a small set
// per table and removes an entire class of order-dependent false negative. The
// price is the occasional spurious match through an aggressive form, which
// scoring penalizes and validation against the data knocks down.
//
// The original name always comes first, so an exact match is always preferred.
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

	// A table name carries at most one legacy prefix in practice, so the first
	// match wins here.
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

// depluralize applies every applicable plural rule, in declared order, and
// returns each singular form produced.
func (p *Profile) depluralize(name string) []string {
	out := make([]string, 0, 2)
	for _, rule := range p.Plural {
		if !strings.HasSuffix(name, rule.Suffix) {
			continue
		}
		stem := name[:len(name)-len(rule.Suffix)]
		if stem == "" {
			// The whole name is the plural suffix; stripping it says nothing.
			continue
		}
		out = append(out, stem+rule.Singular)
	}
	return out
}

// Match reports whether an entity name matches any form of a table name,
// returning the form that matched. The first match wins, and because TableForms
// puts the unchanged name first, an exact match always beats a normalized one.
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
