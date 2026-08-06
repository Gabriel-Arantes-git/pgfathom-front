package profile_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/profile"
)

func TestEmbeddedProfilesAreShipped(t *testing.T) {
	names := profile.Names()
	for _, want := range []string{"pt-br", "en", "es"} {
		found := false
		for _, n := range names {
			if n == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("profile %q is not embedded; got %v", want, names)
		}
	}
}

// TestEmbeddedProfilesAreConsistent guards the shipped profiles against the
// failure mode that is invisible at runtime: a rule that quietly never fires,
// or an affix ordering that strips less than intended. Both produce candidates
// that are never generated, and a candidate that is never generated is never
// validated and never reported.
func TestEmbeddedProfilesAreConsistent(t *testing.T) {
	for _, name := range profile.Names() {
		t.Run(name, func(t *testing.T) {
			p, err := profile.Embedded(name)
			if err != nil {
				t.Fatalf("loading: %v", err)
			}
			if err := p.Validate(); err != nil {
				t.Fatalf("validating: %v", err)
			}
			if p.Name != name {
				t.Errorf("profile declares name %q but is filed as %q", p.Name, name)
			}

			// The generic drop-the-s rule must come last, or it shadows every
			// specific rule that would have produced a better form.
			last := p.Plural[len(p.Plural)-1]
			if last.Suffix != "s" || last.Singular != "" {
				t.Errorf("last plural rule is %+v, want the generic {s -> \"\"} rule", last)
			}
		})
	}
}

func TestLoadResolvesNamesAndPaths(t *testing.T) {
	t.Run("bare name resolves to an embedded profile", func(t *testing.T) {
		p, err := profile.Load("pt-br")
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if p.Name != "pt-br" {
			t.Errorf("name = %q, want pt-br", p.Name)
		}
	})

	t.Run("empty name falls back to the default", func(t *testing.T) {
		p, err := profile.Load("")
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if p.Name != profile.DefaultName {
			t.Errorf("name = %q, want %q", p.Name, profile.DefaultName)
		}
	})

	t.Run("path resolves from disk", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "custom.toml")
		write(t, path, `
name = "casa"
column_suffixes = ["_id"]
[[plural]]
suffix = "s"
singular = ""
`)
		p, err := profile.Load(path)
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if p.Name != "casa" {
			t.Errorf("name = %q, want casa", p.Name)
		}
	})
}

// TestLoadFailsLoudly covers the rule that a profile must never load partially.
// A half-applied naming convention degrades recall silently, which is the worst
// kind of failure this project can have: nothing errors, results just quietly
// get worse.
func TestLoadFailsLoudly(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantIn  string
	}{
		{
			name:    "missing name",
			content: "column_suffixes = [\"_id\"]\n[[plural]]\nsuffix = \"s\"\nsingular = \"\"\n",
			wantIn:  "name",
		},
		{
			name:    "no plural rules",
			content: "name = \"x\"\ncolumn_suffixes = [\"_id\"]\n",
			wantIn:  "plural rule",
		},
		{
			name:    "malformed toml",
			content: "name = \"x\ncolumn_suffixes = [",
			wantIn:  "parsing",
		},
		{
			name:    "duplicate suffix",
			content: "name = \"x\"\ncolumn_suffixes = [\"_id\", \"_id\"]\n[[plural]]\nsuffix = \"s\"\nsingular = \"\"\n",
			wantIn:  "duplicate",
		},
		{
			name:    "duplicate plural rule",
			content: "name = \"x\"\n[[plural]]\nsuffix = \"s\"\nsingular = \"\"\n[[plural]]\nsuffix = \"s\"\nsingular = \"\"\n",
			wantIn:  "duplicate",
		},
		{
			name:    "shadowing suffix order",
			content: "name = \"x\"\ncolumn_suffixes = [\"cod\", \"_cod\"]\n[[plural]]\nsuffix = \"s\"\nsingular = \"\"\n",
			wantIn:  "shadow",
		},
		{
			name:    "shadowing prefix order",
			content: "name = \"x\"\ntable_prefixes = [\"tb\", \"tb_\"]\n[[plural]]\nsuffix = \"s\"\nsingular = \"\"\n",
			wantIn:  "shadow",
		},
		{
			name:    "empty plural suffix",
			content: "name = \"x\"\n[[plural]]\nsuffix = \"\"\nsingular = \"\"\n",
			wantIn:  "must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "p.toml")
			write(t, path, tt.content)

			p, err := profile.LoadFile(path)
			if err == nil {
				t.Fatalf("LoadFile succeeded with %+v, want an error", p)
			}
			if !strings.Contains(err.Error(), tt.wantIn) {
				t.Errorf("error = %q, want it to mention %q", err, tt.wantIn)
			}
		})
	}

	t.Run("missing file", func(t *testing.T) {
		if _, err := profile.LoadFile(filepath.Join(t.TempDir(), "absent.toml")); err == nil {
			t.Fatal("LoadFile succeeded on a missing file")
		}
	})

	t.Run("unknown embedded name", func(t *testing.T) {
		_, err := profile.Embedded("klingon")
		if !errors.Is(err, profile.ErrUnknownProfile) {
			t.Fatalf("error = %v, want ErrUnknownProfile", err)
		}
		if !strings.Contains(err.Error(), "pt-br") {
			t.Errorf("error %q should list the available profiles", err)
		}
	})
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing fixture: %v", err)
	}
}
