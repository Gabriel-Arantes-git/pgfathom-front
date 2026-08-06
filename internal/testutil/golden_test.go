package testutil

import (
	"strings"
	"testing"
)

func TestFirstDifference(t *testing.T) {
	tests := []struct {
		name       string
		want, got  string
		wantReport string
	}{
		{"identical", "a\nb", "a\nb", ""},
		{"changed line", "a\nb\nc", "a\nX\nc", "line 2"},
		{"extra lines", "a", "a\nb", "extra line"},
		{"missing lines", "a\nb", "a", "missing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstDifference(tt.want, tt.got)
			if tt.wantReport == "" {
				if got != "" {
					t.Errorf("firstDifference reported %q for identical input", got)
				}
				return
			}
			if !strings.Contains(got, tt.wantReport) {
				t.Errorf("firstDifference = %q, want it to mention %q", got, tt.wantReport)
			}
		})
	}
}

func TestItoa(t *testing.T) {
	for _, tt := range []struct {
		in   int
		want string
	}{{0, "0"}, {7, "7"}, {42, "42"}, {1234567, "1234567"}} {
		if got := itoa(tt.in); got != tt.want {
			t.Errorf("itoa(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
