package profile_test

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/lvcas-dotcom/pgfathom/internal/profile"
)

func mustLoad(t *testing.T, name string) *profile.Profile {
	t.Helper()
	p, err := profile.Embedded(name)
	if err != nil {
		t.Fatalf("loading profile %q: %v", name, err)
	}
	return p
}

func TestEntityName(t *testing.T) {
	tests := []struct {
		profile string
		column  string
		want    string
	}{
		// Reference suffixes.
		{"pt-br", "cliente_id", "cliente"},
		{"pt-br", "municipio_codigo", "municipio"},
		{"pt-br", "fornecedor_cod", "fornecedor"},
		{"pt-br", "contrato_ref", "contrato"},
		{"pt-br", "empenho_fk", "empenho"},

		// Reference prefixes.
		{"pt-br", "id_fornecedor", "fornecedor"},
		{"pt-br", "cod_municipio", "municipio"},
		{"pt-br", "fk_processo", "processo"},

		// A column with no affix references the entity it is named after.
		{"pt-br", "municipio", "municipio"},
		{"pt-br", "situacao", "situacao"},

		// Longer affixes must win: with "cod" shadowing "_cod", "pedido_cod"
		// would yield "pedido_" instead of "pedido".
		{"pt-br", "pedido_cod", "pedido"},
		{"pt-br", "pedidoid", "pedido"},

		// Folding matches PostgreSQL's own handling of unquoted identifiers.
		{"pt-br", "Cliente_ID", "cliente"},
		{"pt-br", "  cliente_id  ", "cliente"},

		// Stripping must never leave nothing behind: a column that is nothing
		// but an affix keeps a usable name instead of collapsing to "".
		{"pt-br", "id", "id"},
		{"pt-br", "_id", "id"},
		{"pt-br", "cod", "cod"},

		{"en", "customer_id", "customer"},
		{"en", "order_uuid", "order"},
		{"en", "invoice", "invoice"},

		{"es", "cliente_id", "cliente"},
		{"es", "municipio_codigo", "municipio"},
		{"es", "id_proveedor", "proveedor"},
	}

	for _, tt := range tests {
		t.Run(tt.profile+"/"+tt.column, func(t *testing.T) {
			p := mustLoad(t, tt.profile)
			if got := p.EntityName(tt.column); got != tt.want {
				t.Errorf("EntityName(%q) = %q, want %q", tt.column, got, tt.want)
			}
		})
	}
}

func TestTableFormsAlwaysContainTheOriginal(t *testing.T) {
	p := mustLoad(t, "pt-br")

	for _, table := range []string{"cliente", "clientes", "tb_clientes", "opcoes", "x"} {
		forms := p.TableForms(table)
		if len(forms) == 0 {
			t.Fatalf("TableForms(%q) returned nothing", table)
		}
		if forms[0].Value != table || !forms[0].Origin.Exact() {
			t.Errorf("TableForms(%q)[0] = %+v, want the unchanged name first", table, forms[0])
		}
	}
}

// TestTableFormsExactSet pins the whole ordered set for a name that exercises
// every transformation. The order is part of the contract — the unchanged name
// must come first so an exact match always beats a normalized one — and the
// origins are what scoring later uses to tell those apart.
func TestTableFormsExactSet(t *testing.T) {
	p := mustLoad(t, "pt-br")

	want := []profile.Form{
		{Value: "tb_clientes", Origin: profile.OriginExact},
		{Value: "clientes", Origin: profile.OriginPrefixStripped},
		{Value: "tb_cliente", Origin: profile.OriginDepluralized},
		{Value: "cliente", Origin: profile.OriginPrefixAndDepluralized},
	}

	if diff := cmp.Diff(want, p.TableForms("tb_clientes")); diff != "" {
		t.Errorf("TableForms(\"tb_clientes\") mismatch (-want +got):\n%s", diff)
	}
}

func TestDepluralization(t *testing.T) {
	tests := []struct {
		profile string
		table   string
		want    string // must appear among the candidate forms
	}{
		// Portuguese, unaccented — the common case in real databases.
		{"pt-br", "clientes", "cliente"},
		{"pt-br", "opcoes", "opcao"},
		{"pt-br", "paes", "pao"},
		{"pt-br", "animais", "animal"},
		{"pt-br", "responsaveis", "responsavel"},
		{"pt-br", "lencois", "lencol"},
		{"pt-br", "perfis", "perfil"},
		{"pt-br", "armazens", "armazem"},
		{"pt-br", "mulheres", "mulher"},
		{"pt-br", "meses", "mes"},
		{"pt-br", "paises", "pais"},
		// The generic drop-the-s already produces the right form for -ãos.
		{"pt-br", "orgaos", "orgao"},

		// Portuguese, accented.
		{"pt-br", "opções", "opção"},
		{"pt-br", "órgãos", "órgão"},
		{"pt-br", "papéis", "papel"},

		// Legacy convention prefixes.
		{"pt-br", "tb_clientes", "cliente"},
		{"pt-br", "tbl_fornecedores", "fornecedor"},
		{"pt-br", "cad_municipios", "municipio"},

		{"en", "orders", "order"},
		{"en", "categories", "category"},
		{"en", "addresses", "address"},
		{"en", "boxes", "box"},
		{"en", "batches", "batch"},
		{"en", "indices", "index"},
		{"en", "tbl_invoices", "invoice"},

		{"es", "clientes", "cliente"},
		{"es", "camiones", "camion"},
		{"es", "luces", "luz"},
		{"es", "papeles", "papel"},
		{"es", "ciudades", "ciudad"},
		{"es", "planes", "plan"},
	}

	for _, tt := range tests {
		t.Run(tt.profile+"/"+tt.table, func(t *testing.T) {
			p := mustLoad(t, tt.profile)

			forms := p.TableForms(tt.table)
			values := make([]string, len(forms))
			for i, f := range forms {
				values[i] = f.Value
			}

			if !slices.Contains(values, tt.want) {
				t.Errorf("TableForms(%q) = %v, want it to contain %q", tt.table, values, tt.want)
			}
		})
	}
}

// TestAmbiguousPluralYieldsBothForms is the case that motivates returning a set
// instead of a single form. Nothing in "logins" says whether the ns→m rule
// (right for "armazens") or the generic rule applies, so both must survive.
func TestAmbiguousPluralYieldsBothForms(t *testing.T) {
	p := mustLoad(t, "pt-br")

	forms := p.TableForms("logins")
	values := make([]string, len(forms))
	for i, f := range forms {
		values[i] = f.Value
	}

	for _, want := range []string{"logim", "login"} {
		if !slices.Contains(values, want) {
			t.Errorf("TableForms(\"logins\") = %v, want it to contain %q", values, want)
		}
	}

	if _, ok := p.Match("login", "logins"); !ok {
		t.Error("entity \"login\" should match table \"logins\"")
	}
	if _, ok := p.Match("armazem", "armazens"); !ok {
		t.Error("entity \"armazem\" should match table \"armazens\"")
	}
}

func TestMatchReportsHowItMatched(t *testing.T) {
	p := mustLoad(t, "pt-br")

	tests := []struct {
		entity, table string
		wantOrigin    profile.Origin
	}{
		{"cliente", "cliente", profile.OriginExact},
		{"cliente", "clientes", profile.OriginDepluralized},
		{"clientes", "tb_clientes", profile.OriginPrefixStripped},
		{"cliente", "tb_clientes", profile.OriginPrefixAndDepluralized},
	}

	for _, tt := range tests {
		t.Run(tt.entity+"~"+tt.table, func(t *testing.T) {
			form, ok := p.Match(tt.entity, tt.table)
			if !ok {
				t.Fatalf("Match(%q, %q) did not match", tt.entity, tt.table)
			}
			if form.Origin != tt.wantOrigin {
				t.Errorf("origin = %v, want %v", form.Origin, tt.wantOrigin)
			}
		})
	}
}

func TestMatchRejectsUnrelatedNames(t *testing.T) {
	p := mustLoad(t, "pt-br")

	for _, tt := range []struct{ entity, table string }{
		{"cliente", "fornecedor"},
		{"cliente", ""},
		{"", "cliente"},
		{"status", "situacao"},
	} {
		if form, ok := p.Match(tt.entity, tt.table); ok {
			t.Errorf("Match(%q, %q) matched as %+v, want no match", tt.entity, tt.table, form)
		}
	}
}
