// Package testutil holds helpers shared across the test suites. It is
// deliberately dependency-free: whoever approves running this against
// production will read go.mod, and golden files are not worth a line in it.
package testutil

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite golden files instead of comparing against them")

// Golden compares got against testdata/<name>.golden, or rewrites the file
// when -update is set. Formatting regresses silently; a golden file turns that
// into a reviewable diff.
func Golden(t *testing.T, name, got string) {
	t.Helper()

	path := filepath.Join("testdata", name+".golden")

	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("creating testdata directory: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o600); err != nil {
			t.Fatalf("writing golden file %s: %v", path, err)
		}
		t.Logf("updated %s", path)
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("golden file %s is missing; run: go test ./... -update", path)
		}
		t.Fatalf("reading golden file %s: %v", path, err)
	}

	if got != string(want) {
		t.Errorf("output does not match %s\n--- want ---\n%s\n--- got ---\n%s\n%s",
			path, want, got, firstDifference(string(want), got))
	}
}

// firstDifference points at the first diverging line, so nobody has to eyeball
// two long blocks for one changed character.
func firstDifference(want, got string) string {
	wantLines, gotLines := strings.Split(want, "\n"), strings.Split(got, "\n")

	for i := 0; i < len(wantLines) && i < len(gotLines); i++ {
		if wantLines[i] != gotLines[i] {
			return "--- first difference ---\n" +
				"line " + itoa(i+1) + ":\n" +
				"  want: " + wantLines[i] + "\n" +
				"  got:  " + gotLines[i]
		}
	}

	switch {
	case len(gotLines) > len(wantLines):
		return "--- first difference ---\ngot has " + itoa(len(gotLines)-len(wantLines)) + " extra line(s)"
	case len(wantLines) > len(gotLines):
		return "--- first difference ---\ngot is missing " + itoa(len(wantLines)-len(gotLines)) + " line(s)"
	default:
		return ""
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
