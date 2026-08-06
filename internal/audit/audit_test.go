package audit_test

import (
	"strings"
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/audit"
	"github.com/lvcas-dotcom/pgfathom/internal/model"
)

func schemaWith(tables ...model.Table) []model.Schema {
	return []model.Schema{{Name: "public", Tables: tables}}
}

func table(name string, fks ...model.ForeignKey) model.Table {
	return model.Table{
		Schema:      "public",
		Name:        name,
		PrimaryKey:  []string{"id"},
		ForeignKeys: fks,
		Stats:       model.TableStats{EstimatedRows: 4_000_000},
	}
}

func fk(name, column, refTable string, validated, hasIndex bool) model.ForeignKey {
	return model.ForeignKey{
		Name:       name,
		Columns:    []string{column},
		RefSchema:  "public",
		RefTable:   refTable,
		RefColumns: []string{"id"},
		Validated:  validated,
		HasIndex:   hasIndex,
	}
}

func kinds(findings []model.Finding) []model.FindingKind {
	out := make([]model.FindingKind, len(findings))
	for i, f := range findings {
		out[i] = f.Kind
	}
	return out
}

func countKind(findings []model.Finding, k model.FindingKind) int {
	n := 0
	for _, f := range findings {
		if f.Kind == k {
			n++
		}
	}
	return n
}

func TestNotValidConstraintIsReported(t *testing.T) {
	schemas := schemaWith(table("pedido", fk("pedido_cliente_fk", "cliente_id", "cliente", false, true)))

	findings := audit.Findings(schemas)

	if got := countKind(findings, model.FindingNotValidConstraint); got != 1 {
		t.Fatalf("got %d NOT VALID findings, want 1 (kinds: %v)", got, kinds(findings))
	}

	f := findings[0]
	if !strings.Contains(f.Object, "pedido_cliente_fk") {
		t.Errorf("object = %q, want it to name the constraint", f.Object)
	}
	if f.Metrics["child_estimated_rows"] != 4_000_000 {
		t.Errorf("metrics = %v, want the child row estimate", f.Metrics)
	}
}

func TestValidatedConstraintIsQuiet(t *testing.T) {
	schemas := schemaWith(table("pedido", fk("pedido_cliente_fk", "cliente_id", "cliente", true, true)))

	if got := countKind(audit.Findings(schemas), model.FindingNotValidConstraint); got != 0 {
		t.Errorf("got %d NOT VALID findings on a fully validated schema, want 0", got)
	}
}

func TestUnindexedForeignKeyIsReported(t *testing.T) {
	parent := table("cliente")
	parent.Stats.EstimatedRows = 90_000

	schemas := schemaWith(
		table("pedido", fk("pedido_cliente_fk", "cliente_id", "cliente", true, false)),
		parent,
	)

	findings := audit.Findings(schemas)

	if got := countKind(findings, model.FindingFKWithoutIndex); got != 1 {
		t.Fatalf("got %d unindexed findings, want 1 (kinds: %v)", got, kinds(findings))
	}

	f := findings[0]
	if f.Metrics["child_estimated_rows"] != 4_000_000 {
		t.Errorf("metrics = %v, want the child row estimate", f.Metrics)
	}
	if f.Metrics["parent_estimated_rows"] != 90_000 {
		t.Errorf("metrics = %v, want the parent row estimate, which is what shows the severity", f.Metrics)
	}
	if !strings.Contains(f.Detail, "cliente_id") {
		t.Errorf("detail = %q, want it to name the child column", f.Detail)
	}
}

func TestIndexedForeignKeyIsQuiet(t *testing.T) {
	schemas := schemaWith(table("pedido", fk("pedido_cliente_fk", "cliente_id", "cliente", true, true)))

	if got := countKind(audit.Findings(schemas), model.FindingFKWithoutIndex); got != 0 {
		t.Errorf("got %d unindexed findings on an indexed schema, want 0", got)
	}
}

func TestOneConstraintCanProduceBothFindings(t *testing.T) {
	schemas := schemaWith(table("pedido", fk("pedido_cliente_fk", "cliente_id", "cliente", false, false)))

	findings := audit.Findings(schemas)

	if len(findings) != 2 {
		t.Fatalf("got %d findings, want 2 (kinds: %v)", len(findings), kinds(findings))
	}
	if countKind(findings, model.FindingNotValidConstraint) != 1 ||
		countKind(findings, model.FindingFKWithoutIndex) != 1 {
		t.Errorf("kinds = %v, want one of each", kinds(findings))
	}
}

// TestFindingsAreOrderStable matters because the output is compared against
// golden files from the next phase onward, and map iteration would otherwise
// shuffle it between runs.
func TestFindingsAreOrderStable(t *testing.T) {
	schemas := schemaWith(
		table("zeta", fk("zeta_fk", "a_id", "alpha", false, false)),
		table("alpha", fk("alpha_fk", "b_id", "beta", false, false)),
	)

	first := audit.Findings(schemas)
	for i := 0; i < 20; i++ {
		again := audit.Findings(schemas)
		for j := range first {
			if first[j].Object != again[j].Object || first[j].Kind != again[j].Kind {
				t.Fatalf("ordering changed between runs at %d: %v vs %v", j, first[j], again[j])
			}
		}
	}
}

func TestEmptySchemaProducesNothing(t *testing.T) {
	if findings := audit.Findings(nil); len(findings) != 0 {
		t.Errorf("got %d findings from no schemas, want 0", len(findings))
	}
	if findings := audit.Findings(schemaWith()); len(findings) != 0 {
		t.Errorf("got %d findings from an empty schema, want 0", len(findings))
	}
}
