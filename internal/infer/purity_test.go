package infer_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// allowedProjectImports is the whole dependency surface this package may have.
// Inference is trivial to test precisely because it cannot reach a database,
// and that property is easy to destroy with one convenient import.
var allowedProjectImports = map[string]bool{
	"github.com/lvcas-dotcom/pgfathom/internal/model":   true,
	"github.com/lvcas-dotcom/pgfathom/internal/profile": true,
}

func TestPackageStaysPure(t *testing.T) {
	sources, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("listing sources: %v", err)
	}

	fset := token.NewFileSet()
	checked := 0

	for _, path := range sources {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		checked++

		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading %s: %v", path, err)
		}

		file, err := parser.ParseFile(fset, path, src, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parsing %s: %v", path, err)
		}

		for _, imp := range file.Imports {
			p, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				t.Fatalf("%s: unquoting %s: %v", path, imp.Path.Value, err)
			}

			if strings.HasPrefix(p, "github.com/lvcas-dotcom/pgfathom/") {
				if !allowedProjectImports[p] {
					t.Errorf("%s imports %q: inference may depend only on model and profile", path, p)
				}
				continue
			}

			if first, _, _ := strings.Cut(p, "/"); strings.Contains(first, ".") {
				t.Errorf("%s imports %q: inference must carry no external dependency", path, p)
			}
		}
	}

	if checked == 0 {
		t.Fatal("no non-test sources found; the check would pass vacuously")
	}
}
