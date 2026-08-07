//go:build integration

package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/cli"
	"github.com/lvcas-dotcom/pgfathom/internal/testutil"
)

// plantedValues is the union of the recognizable strings every fixture plants
// in user tables. The commands under test run against a subset, but scanning
// for all of them costs nothing and survives fixture reshuffling.
var plantedValues = []string{
	"529.318.470-11",
	"145.892.663-04",
	"Maria Aparecida Silva",
	"Joao Carlos Pereira",
	"Rua das Acacias 42",
	"Construtora Horizonte LTDA",
	"CT-2019-0041",
	"Sao Bernardo do Campo",
	"servico de manutencao",
	"peca de reposicao",
}

// TestNoUserDataInAnyStream runs the real commands end to end and scans every
// byte the process would emit — stdout, and stderr carrying both diagnostics
// and the debug log. The layer-level scans prove the renderers are clean; this
// one proves no other path around them leaks, because a user pointing the tool
// at a production database gets this composition, not the layers.
func TestNoUserDataInAnyStream(t *testing.T) {
	cases := []struct {
		name    string
		fixture string
		args    []string
	}{
		{"audit table", "not_valid_constraints", []string{"audit"}},
		{"audit json", "not_valid_constraints", []string{"audit", "--format", "json"}},
		{"discover table", "inferable", []string{"discover", "--include-rejected"}},
		{"discover json", "inferable", []string{"discover", "--format", "json"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dsn := testutil.Postgres(t, tc.fixture)

			var out, errOut bytes.Buffer
			streams := &cli.Streams{Out: &out, Err: &errOut, In: strings.NewReader("")}

			// Debug is the most talkative level the tool has; if a value can
			// reach a log line, this is where it surfaces.
			args := append(tc.args, "--dsn", dsn, "--log-level", "debug", "--color", "never")

			if code := cli.Run(args, streams); code != 0 {
				t.Fatalf("exit code %d, stderr:\n%s", code, errOut.String())
			}

			for name, text := range map[string]string{
				"stdout": out.String(),
				"stderr": errOut.String(),
			} {
				for _, planted := range plantedValues {
					if strings.Contains(text, planted) {
						t.Errorf("%s leaked %q from a user table", name, planted)
					}
				}
			}
		})
	}
}
