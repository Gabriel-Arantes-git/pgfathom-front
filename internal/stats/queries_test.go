package stats

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// userRelation matches a FROM or JOIN against anything outside the system
// catalogs. Deliberately blunt, like the catalog layer's version: the point is
// to fail loudly if this layer ever grows a read of table data, because "no
// table I/O" is the property that makes the prefilter free to run anywhere.
var userRelation = regexp.MustCompile(`(?is)\b(from|join)\s+(?:only\s+)?([a-z_][a-z0-9_.]*)`)

var allowedRelations = map[string]bool{
	"pg_stats": true,
	"unnest":   true,
}

func TestQueriesTouchOnlyStatisticsViews(t *testing.T) {
	for name, sql := range queryLiterals(t, "read.go") {
		for _, m := range userRelation.FindAllStringSubmatch(sql, -1) {
			relation := strings.ToLower(m[2])
			if allowedRelations[relation] {
				continue
			}
			t.Errorf("%s reads %q: the prefilter must touch only planner statistics, "+
				"because reading data belongs to validation", name, relation)
		}
	}
}

// queryLiterals returns the SQL string constants, by name, as the catalog
// layer's test does.
func queryLiterals(t *testing.T, path string) map[string]string {
	t.Helper()

	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	file, err := parser.ParseFile(token.NewFileSet(), path, src, 0)
	if err != nil {
		t.Fatalf("parsing %s: %v", path, err)
	}

	out := make(map[string]string)
	ast.Inspect(file, func(n ast.Node) bool {
		spec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for i, name := range spec.Names {
			if i >= len(spec.Values) {
				continue
			}
			lit, ok := spec.Values[i].(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			value, err := strconv.Unquote(lit.Value)
			if err != nil {
				t.Fatalf("unquoting %s: %v", name.Name, err)
			}
			if !strings.Contains(strings.ToUpper(value), "SELECT") {
				continue
			}
			out[name.Name] = value
		}
		return true
	})

	if len(out) == 0 {
		t.Fatalf("no SQL constants found in %s; the check would pass vacuously", path)
	}
	return out
}
