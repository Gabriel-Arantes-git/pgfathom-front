package model_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// forbiddenStdlib are standard library packages that would mean this package
// grew an I/O dependency. encoding/json is deliberately absent from the list:
// serializing the model is part of its job.
var forbiddenStdlib = []string{
	"os",
	"net",
	"net/http",
	"io/fs",
	"os/exec",
	"database/sql",
	"path/filepath",
	"log",
	"log/slog",
}

// TestPackageIsPure enforces that internal/model imports nothing but a narrow
// slice of the standard library.
//
// The whole testability of the project rests on this package staying pure: the
// inference layer is trivial to test precisely because the types it operates on
// cannot reach a database. That property is easy to destroy with one convenient
// import and expensive to restore afterwards, so it is checked by a test rather
// than trusted to review.
func TestPackageIsPure(t *testing.T) {
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
				t.Fatalf("%s: unquoting import %s: %v", path, imp.Path.Value, err)
			}

			// A dot in the first path segment means a module host, which means
			// an external dependency.
			if first, _, _ := strings.Cut(p, "/"); strings.Contains(first, ".") {
				t.Errorf("%s imports %q: internal/model must depend on the standard library only", path, p)
				continue
			}

			for _, forbidden := range forbiddenStdlib {
				if p == forbidden {
					t.Errorf("%s imports %q: internal/model must perform no I/O", path, p)
				}
			}
		}
	}

	if checked == 0 {
		t.Fatal("no non-test source files found; the purity check would pass vacuously")
	}
}
