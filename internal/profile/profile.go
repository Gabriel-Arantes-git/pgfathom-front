// Package profile loads naming profiles: the affix and plural rules that
// translate a schema's naming convention into entity names.
//
// The rules live in TOML, never in Go constants. This is the most fragile logic
// in the project — a wrong plural rule produces a candidate that is never
// generated and therefore never validated — and keeping it in data makes it
// testable by table of cases. It also makes teaching pgfathom a new language a
// config file rather than a patch.
package profile

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

//go:embed profiles/*.toml
var embedded embed.FS

const embeddedDir = "profiles"

// DefaultName is the profile used when none is requested. The target audience
// for this tool runs Portuguese-language schemas.
const DefaultName = "pt-br"

// PluralRule turns a plural suffix into its singular form.
type PluralRule struct {
	Suffix   string `toml:"suffix"`
	Singular string `toml:"singular"`
}

// Profile is a set of naming conventions.
type Profile struct {
	Name string `toml:"name"`

	// ColumnSuffixes and ColumnPrefixes are reference affixes stripped from a
	// column name to reveal the entity it points at. First match wins, so
	// longer affixes must precede shorter ones they would shadow; Validate
	// enforces that.
	ColumnSuffixes []string `toml:"column_suffixes"`
	ColumnPrefixes []string `toml:"column_prefixes"`

	// TablePrefixes are legacy convention prefixes stripped from table names.
	TablePrefixes []string `toml:"table_prefixes"`

	// TypeSuffixes mark the discriminator half of a polymorphic pair, as in
	// documento_tipo beside documento_id.
	TypeSuffixes []string `toml:"type_suffixes"`

	// GenericEntities are entity names that match almost anything and mean
	// little — the domain tables a legacy schema has dozens of. They are
	// language-dependent, so they live here rather than in Go.
	GenericEntities []string `toml:"generic_entities"`

	// Plural rules run from most specific to most generic. Unlike the affix
	// lists, every applicable rule contributes a form; order sets preference
	// within the set, not exclusivity.
	Plural []PluralRule `toml:"plural"`
}

// ErrUnknownProfile is returned when a name matches no embedded profile and is
// not a readable path.
var ErrUnknownProfile = errors.New("unknown naming profile")

// Load resolves a profile by embedded name or by file path. A bare name never
// touches the filesystem, which is what keeps the single-binary promise.
func Load(nameOrPath string) (*Profile, error) {
	if nameOrPath == "" {
		nameOrPath = DefaultName
	}

	if !strings.ContainsAny(nameOrPath, `/\.`) {
		return Embedded(nameOrPath)
	}
	return LoadFile(nameOrPath)
}

// Embedded returns one of the profiles compiled into the binary.
func Embedded(name string) (*Profile, error) {
	data, err := embedded.ReadFile(path.Join(embeddedDir, name+".toml"))
	if err != nil {
		return nil, fmt.Errorf("%w: %q (available: %s)", ErrUnknownProfile, name, strings.Join(Names(), ", "))
	}
	return parse(data, "embedded:"+name)
}

// LoadFile reads a profile from disk.
func LoadFile(filePath string) (*Profile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading naming profile: %w", err)
	}
	return parse(data, filePath)
}

// Names lists the embedded profiles, sorted.
func Names() []string {
	entries, err := embedded.ReadDir(embeddedDir)
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, strings.TrimSuffix(e.Name(), ".toml"))
	}
	sort.Strings(names)
	return names
}

func parse(data []byte, source string) (*Profile, error) {
	var p Profile
	if err := toml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parsing naming profile %s: %w", source, err)
	}
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("invalid naming profile %s: %w", source, err)
	}
	return &p, nil
}

// Validate rejects a profile that would misbehave silently. A quietly wrong
// profile produces candidates that are never generated, never validated and
// never reported — an invisible failure, which is why it is caught at load.
func (p *Profile) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("field \"name\" is required")
	}
	if len(p.Plural) == 0 {
		return errors.New("at least one plural rule is required")
	}

	for _, set := range []struct {
		field string
		items []string
	}{
		{"column_suffixes", p.ColumnSuffixes},
		{"column_prefixes", p.ColumnPrefixes},
		{"table_prefixes", p.TablePrefixes},
		{"type_suffixes", p.TypeSuffixes},
		{"generic_entities", p.GenericEntities},
	} {
		if err := checkNoDuplicates(set.field, set.items); err != nil {
			return err
		}
	}

	// With "cod" declared before "_cod", the column "pedido_cod" would yield
	// "pedido_" instead of "pedido".
	if err := checkSuffixOrder("column_suffixes", p.ColumnSuffixes); err != nil {
		return err
	}
	if err := checkPrefixOrder("column_prefixes", p.ColumnPrefixes); err != nil {
		return err
	}
	if err := checkPrefixOrder("table_prefixes", p.TablePrefixes); err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(p.Plural))
	for i, r := range p.Plural {
		if r.Suffix == "" {
			return fmt.Errorf("plural rule %d: suffix must not be empty", i+1)
		}
		if _, dup := seen[r.Suffix]; dup {
			return fmt.Errorf("plural rule %d: duplicate suffix %q", i+1, r.Suffix)
		}
		seen[r.Suffix] = struct{}{}
	}

	return nil
}

func checkNoDuplicates(field string, items []string) error {
	seen := make(map[string]struct{}, len(items))
	for _, it := range items {
		if it == "" {
			return fmt.Errorf("%s: empty entry", field)
		}
		if _, dup := seen[it]; dup {
			return fmt.Errorf("%s: duplicate entry %q", field, it)
		}
		seen[it] = struct{}{}
	}
	return nil
}

func checkSuffixOrder(field string, items []string) error {
	for i, short := range items {
		for _, long := range items[i+1:] {
			if len(long) > len(short) && strings.HasSuffix(long, short) {
				return fmt.Errorf(
					"%s: %q is declared before %q and would shadow it; declare the longer affix first",
					field, short, long)
			}
		}
	}
	return nil
}

func checkPrefixOrder(field string, items []string) error {
	for i, short := range items {
		for _, long := range items[i+1:] {
			if len(long) > len(short) && strings.HasPrefix(long, short) {
				return fmt.Errorf(
					"%s: %q is declared before %q and would shadow it; declare the longer affix first",
					field, short, long)
			}
		}
	}
	return nil
}
