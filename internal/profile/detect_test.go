package profile_test

import (
	"slices"
	"testing"

	"github.com/lvcas-dotcom/pgfathom/internal/model"
)

func table(name string, fks ...model.ForeignKey) model.Table {
	return model.Table{
		Schema:      "public",
		Name:        name,
		Columns:     []model.Column{{Name: "idkey", BaseType: "int8"}},
		PrimaryKey:  []string{"idkey"},
		ForeignKeys: fks,
	}
}

func fk(column, refTable string) model.ForeignKey {
	return model.ForeignKey{
		Columns:    []string{column},
		RefSchema:  "public",
		RefTable:   refTable,
		RefColumns: []string{"idkey"},
		Validated:  true,
	}
}

func schemaOf(tables ...model.Table) []model.Schema {
	return []model.Schema{{Name: "public", Tables: tables}}
}

func affixes(ev []model.NamingEvidence) []string {
	out := make([]string, len(ev))
	for i, e := range ev {
		out[i] = e.Affix
	}
	return out
}

// TestDetectsSuffixFromDeclaredKeys reproduces the convention that made a real
// 784-table schema score 0.5% recall: a suffix no shipped profile knew about.
func TestDetectsSuffixFromDeclaredKeys(t *testing.T) {
	p := mustLoad(t, "pt-br")

	d := p.Detect(schemaOf(
		table("lote"),
		table("bairro"),
		table("logradouro"),
		table("imovel",
			fk("lote_idkey", "lote"),
			fk("bairro_idkey", "bairro"),
			fk("logradouro_idkey", "logradouro"),
		),
	))

	if got := affixes(d.ColumnSuffixes); !slices.Contains(got, "_idkey") {
		t.Errorf("suffixes = %v, want _idkey read from the declared keys", got)
	}
	if d.DeclaredKeys != 3 {
		t.Errorf("DeclaredKeys = %d, want 3", d.DeclaredKeys)
	}
}

func TestDetectsPrefixFromDeclaredKeys(t *testing.T) {
	p := mustLoad(t, "pt-br")

	d := p.Detect(schemaOf(
		table("lote"),
		table("bairro"),
		table("operador"),
		table("calculo",
			fk("idkey_lote", "lote"),
			fk("idkey_bairro", "bairro"),
			fk("idkey_operador", "operador"),
		),
	))

	if got := affixes(d.ColumnPrefixes); !slices.Contains(got, "idkey_") {
		t.Errorf("prefixes = %v, want idkey_", got)
	}
}

// TestDetectsApplicationTablePrefix is the Django case: the prefix belongs to
// the application, not to a language, so no shipped profile could carry it.
func TestDetectsApplicationTablePrefix(t *testing.T) {
	p := mustLoad(t, "en")

	d := p.Detect(schemaOf(
		table("auth_user"),
		table("auth_group"),
		table("auth_permission"),
		table("django_content_type"),
		table("django_migrations"),
	))

	if got := affixes(d.TablePrefixes); !slices.Contains(got, "auth_") {
		t.Errorf("table prefixes = %v, want auth_", got)
	}
}

// TestSharedSubjectIsNotAPrefix is why the separator is required: these tables
// share a topic, not a convention, and stripping "ato" would produce nonsense.
func TestSharedSubjectIsNotAPrefix(t *testing.T) {
	p := mustLoad(t, "pt-br")

	d := p.Detect(schemaOf(
		table("atotipo"), table("atoanexo"), table("atotramite"),
		table("atoconteudo"), table("atoresponsavel"),
	))

	for _, e := range d.TablePrefixes {
		if e.Affix == "ato" {
			t.Errorf("detected %q as a prefix; without a separator it is a subject, not a convention", e.Affix)
		}
	}
}

func TestIsolatedOccurrenceIsNotAConvention(t *testing.T) {
	p := mustLoad(t, "pt-br")

	tables := []model.Table{table("lote")}
	for i := 0; i < 40; i++ {
		tables = append(tables, table("t"+string(rune('a'+i%26))+string(rune('a'+i/26))))
	}
	tables = append(tables, table("imovel", fk("lote_esquisito", "lote")))

	d := p.Detect(schemaOf(tables...))

	if got := affixes(d.ColumnSuffixes); slices.Contains(got, "_esquisito") {
		t.Errorf("suffixes = %v, want a single occurrence treated as noise", got)
	}
}

func TestSchemaWithoutDeclaredKeysDetectsNoAffix(t *testing.T) {
	p := mustLoad(t, "pt-br")

	d := p.Detect(schemaOf(table("lote"), table("bairro"), table("imovel")))

	if len(d.ColumnSuffixes) != 0 || len(d.ColumnPrefixes) != 0 {
		t.Errorf("detected %v / %v with nothing to read from", d.ColumnSuffixes, d.ColumnPrefixes)
	}
	if d.DeclaredKeys != 0 {
		t.Errorf("DeclaredKeys = %d, want 0", d.DeclaredKeys)
	}
}

// TestMergeKeepsEveryBaseRule is the safety property the whole feature rests on:
// detection may only add. Dropping a base rule would cost a candidate that never
// gets generated, and one that never gets generated is invisible.
func TestMergeKeepsEveryBaseRule(t *testing.T) {
	base := mustLoad(t, "pt-br")

	merged := base.WithDetected(model.NamingDetection{
		ColumnSuffixes: []model.NamingEvidence{{Affix: "_esquisito", Occurrences: 9}},
		TablePrefixes:  []model.NamingEvidence{{Affix: "app_", Occurrences: 9}},
	})

	for _, want := range base.ColumnSuffixes {
		if !slices.Contains(merged.ColumnSuffixes, want) {
			t.Errorf("merge dropped the base suffix %q", want)
		}
	}
	for _, want := range base.TablePrefixes {
		if !slices.Contains(merged.TablePrefixes, want) {
			t.Errorf("merge dropped the base table prefix %q", want)
		}
	}
	if !slices.Contains(merged.ColumnSuffixes, "_esquisito") {
		t.Error("merge did not add the detected suffix")
	}
}

func TestMergedProfileStillValidates(t *testing.T) {
	base := mustLoad(t, "pt-br")

	// "key" would shadow "_key" if it landed before it.
	merged := base.WithDetected(model.NamingDetection{
		ColumnSuffixes: []model.NamingEvidence{{Affix: "key", Occurrences: 9}},
		ColumnPrefixes: []model.NamingEvidence{{Affix: "id", Occurrences: 9}},
	})

	if err := merged.Validate(); err != nil {
		t.Fatalf("merged profile does not validate: %v", err)
	}
}

func TestMergeIsIdempotent(t *testing.T) {
	base := mustLoad(t, "pt-br")
	d := model.NamingDetection{ColumnSuffixes: []model.NamingEvidence{{Affix: "_idkey", Occurrences: 9}}}

	once := base.WithDetected(d)
	twice := once.WithDetected(d)

	if len(once.ColumnSuffixes) != len(twice.ColumnSuffixes) {
		t.Errorf("applying the same detection twice changed the profile: %d vs %d",
			len(once.ColumnSuffixes), len(twice.ColumnSuffixes))
	}
}

// TestDetectionRecoversTheRealCase is the end-to-end point of the phase: the
// same schema shape that scored 0.5% must resolve without anyone hand-editing a
// profile first.
func TestDetectionRecoversTheRealCase(t *testing.T) {
	base := mustLoad(t, "en")

	schemas := schemaOf(
		table("lote"),
		table("bairro"),
		table("logradouro"),
		table("imovel",
			fk("lote_idkey", "lote"),
			fk("bairro_idkey", "bairro"),
			fk("logradouro_idkey", "logradouro"),
		),
	)

	if _, ok := base.Match(base.EntityName("lote_idkey"), "lote"); ok {
		t.Fatal("the base profile already matches; the test proves nothing")
	}

	detected := base.WithDetected(base.Detect(schemas))

	if _, ok := detected.Match(detected.EntityName("lote_idkey"), "lote"); !ok {
		t.Error("after detection, lote_idkey should reach lote")
	}
}

func TestEmptyReportsNothingDetected(t *testing.T) {
	if !(model.NamingDetection{}).Empty() {
		t.Error("an empty detection must report itself as empty")
	}
	d := model.NamingDetection{ColumnSuffixes: []model.NamingEvidence{{Affix: "_id"}}}
	if d.Empty() {
		t.Error("a detection with an affix is not empty")
	}
}
