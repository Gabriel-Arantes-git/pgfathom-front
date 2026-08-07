// Package infer turns catalog facts into hypotheses about undeclared
// relationships, and scores them.
//
// It is pure and deterministic: it reads a model and a naming profile and
// touches nothing else. That is what keeps the judgment concentrated here cheap
// to review and cheap to change.
//
// Nothing in this package concludes anything. It answers "what is worth asking
// the data", not "what is true".
package infer

// TypeMatch is how well a child column's type fits a target key's type.
type TypeMatch int

const (
	// TypeIncompatible means no relationship is possible.
	TypeIncompatible TypeMatch = iota

	// TypeCompatible means the values could fit, through widening or through
	// equivalent text types.
	TypeCompatible

	// TypeIdentical is the ideal case: the same base type on both sides.
	TypeIdentical
)

// integerWidth ranks the integer family. A child may point at a key of equal or
// greater width; the reverse admits values the key cannot hold.
var integerWidth = map[string]int{
	"int2": 1,
	"int4": 2,
	"int8": 3,
}

// textualTypes are interchangeable among themselves. Length is not part of the
// base type, so varchar(60) and text compare equal here.
var textualTypes = map[string]bool{
	"text":    true,
	"varchar": true,
	"bpchar":  true,
}

// CompareTypes reports how well a child column type fits a target key type.
//
// The rules are deliberately strict. A missed relationship costs a finding; a
// wrong one that survives to confirmation costs the tool its credibility, so
// anything not clearly allowed is rejected.
func CompareTypes(child, parent string) TypeMatch {
	if child == "" || parent == "" {
		return TypeIncompatible
	}
	if child == parent {
		return TypeIdentical
	}

	if cw, ok := integerWidth[child]; ok {
		pw, ok := integerWidth[parent]
		if ok && cw <= pw {
			return TypeCompatible
		}
		return TypeIncompatible
	}

	if textualTypes[child] && textualTypes[parent] {
		return TypeCompatible
	}

	// uuid, numeric, timestamps and everything else match only themselves,
	// which the equality above already covered.
	return TypeIncompatible
}

// Compatible reports whether a relationship is possible at all.
func (m TypeMatch) Compatible() bool { return m != TypeIncompatible }

// String renders the match for diagnostics.
func (m TypeMatch) String() string {
	switch m {
	case TypeIdentical:
		return "identical"
	case TypeCompatible:
		return "compatible"
	default:
		return "incompatible"
	}
}
