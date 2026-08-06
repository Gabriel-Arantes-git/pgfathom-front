// Package audit derives findings that require no inference at all.
//
// Everything here comes straight from the catalog: deterministic, free, and
// immune to false positives. Every finding it emits is a fact, not a
// hypothesis, which is also what lets the command run usefully against a
// database where inference would have nothing to say.
package audit

import (
	"sort"

	"github.com/lvcas-dotcom/pgfathom/internal/model"
)

// Findings collects the structural findings for a set of schemas.
func Findings(schemas []model.Schema) []model.Finding {
	var out []model.Finding

	for _, s := range schemas {
		for _, t := range s.Tables {
			out = append(out, notValidConstraints(t)...)
			out = append(out, unindexedForeignKeys(t, schemas)...)
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].Object < out[j].Object
	})

	return out
}

func notValidConstraints(t model.Table) []model.Finding {
	var out []model.Finding

	for _, fk := range t.ForeignKeys {
		if fk.Validated {
			continue
		}
		out = append(out, model.Finding{
			Kind:   model.FindingNotValidConstraint,
			Object: t.Ref() + "." + fk.Name,
			Detail: "constraint is NOT VALID: it blocks new violations but never " +
				"checked the rows that were already there",
			Metrics: map[string]int64{"child_estimated_rows": t.Stats.EstimatedRows},
		})
	}

	return out
}

func unindexedForeignKeys(t model.Table, schemas []model.Schema) []model.Finding {
	var out []model.Finding

	for _, fk := range t.ForeignKeys {
		if fk.HasIndex {
			continue
		}

		metrics := map[string]int64{"child_estimated_rows": t.Stats.EstimatedRows}
		if parent, ok := findTable(schemas, fk.RefSchema, fk.RefTable); ok {
			metrics["parent_estimated_rows"] = parent.Stats.EstimatedRows
		}

		out = append(out, model.Finding{
			Kind:   model.FindingFKWithoutIndex,
			Object: t.Ref() + "." + fk.Name,
			Detail: "no index leads with " + fk.Columns[0] + ": every delete on " +
				fk.RefSchema + "." + fk.RefTable + " scans this table sequentially",
			Metrics: metrics,
		})
	}

	return out
}

func findTable(schemas []model.Schema, schema, name string) (model.Table, bool) {
	for _, s := range schemas {
		if s.Name != schema {
			continue
		}
		for _, t := range s.Tables {
			if t.Name == name {
				return t, true
			}
		}
	}
	return model.Table{}, false
}
