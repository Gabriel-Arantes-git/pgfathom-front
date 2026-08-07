package infer_test

import (
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/infer"
)

func TestCompareTypes(t *testing.T) {
	tests := []struct {
		name          string
		child, parent string
		want          infer.TypeMatch
	}{
		{"same integer", "int8", "int8", infer.TypeIdentical},
		{"same text", "text", "text", infer.TypeIdentical},
		{"same uuid", "uuid", "uuid", infer.TypeIdentical},

		{"int4 into int8", "int4", "int8", infer.TypeCompatible},
		{"int2 into int8", "int2", "int8", infer.TypeCompatible},
		{"int2 into int4", "int2", "int4", infer.TypeCompatible},

		// The child would admit values the key cannot hold.
		{"int8 into int4", "int8", "int4", infer.TypeIncompatible},
		{"int4 into int2", "int4", "int2", infer.TypeIncompatible},

		{"varchar into text", "varchar", "text", infer.TypeCompatible},
		{"text into varchar", "text", "varchar", infer.TypeCompatible},
		{"bpchar into text", "bpchar", "text", infer.TypeCompatible},

		// uuid is strict: a text column holding UUIDs is a different problem
		// from a uuid column, and treating them as one invites false positives.
		{"text into uuid", "text", "uuid", infer.TypeIncompatible},
		{"uuid into text", "uuid", "text", infer.TypeIncompatible},

		{"integer into text", "int8", "text", infer.TypeIncompatible},
		{"text into integer", "text", "int8", infer.TypeIncompatible},
		{"numeric into integer", "numeric", "int8", infer.TypeIncompatible},
		{"integer into numeric", "int8", "numeric", infer.TypeIncompatible},

		{"timestamp into date", "timestamptz", "date", infer.TypeIncompatible},
		{"boolean into integer", "bool", "int4", infer.TypeIncompatible},

		{"empty child", "", "int8", infer.TypeIncompatible},
		{"empty parent", "int8", "", infer.TypeIncompatible},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := infer.CompareTypes(tt.child, tt.parent); got != tt.want {
				t.Errorf("CompareTypes(%q, %q) = %v, want %v", tt.child, tt.parent, got, tt.want)
			}
		})
	}
}

func TestTypeMatchCompatible(t *testing.T) {
	for match, want := range map[infer.TypeMatch]bool{
		infer.TypeIdentical:    true,
		infer.TypeCompatible:   true,
		infer.TypeIncompatible: false,
	} {
		if got := match.Compatible(); got != want {
			t.Errorf("%v.Compatible() = %v, want %v", match, got, want)
		}
	}
}
