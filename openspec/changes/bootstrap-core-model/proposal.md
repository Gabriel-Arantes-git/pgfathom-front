## Why

O pgfathom ainda não existe como código — há especificação de produto e nenhuma linha de Go. A lógica mais sutil e mais frágil do projeto inteiro é a normalização de nomes: extrair `cliente` de `cliente_id`, casar contra `tb_clientes`, e acertar a despluralização em português, onde `orgaos` vem de `órgão` e `responsaveis` vem de `responsável`. Errar aqui envenena todas as fases seguintes, porque candidato que nunca é gerado nunca é validado.

Essa lógica é determinística e não precisa de banco, de Docker nem de rede para ser testada exaustivamente. Começar por ela significa que a parte mais propensa a regressão fica coberta antes de qualquer coisa mais cara existir, e que um contribuidor externo consegue rodar a suíte inteira no primeiro `git clone`.

## What Changes

- Esqueleto do repositório: `go.mod` no module path `github.com/lvcas-dotcom/pgfathom`, `Makefile`, `.golangci.yml`, workflow de CI e `.gitignore`.
- `internal/model`: os tipos puros do modelo semântico, sem I/O e sem importar nenhuma outra camada do projeto.
- `internal/profile`: carregamento de perfil de nomenclatura a partir de TOML, perfis `pt-br`, `en` e `es` embarcados no binário via `embed`, e as funções de normalização de nome de coluna e de nome de tabela.
- `cmd/pgfathom`: raiz cobra, subcomando `version`, e as convenções de saída que todo comando futuro herda — detecção de TTY, disciplina de stdout/stderr e códigos de saída.
- Suíte de teste unitário sem nenhuma dependência de Docker, com tabela de casos por perfil de idioma.

Nenhum código desta change conecta em banco. Nenhum código desta change lê dado de usuário.

## Capabilities

### New Capabilities

- `domain-model`: o modelo semântico interno do schema analisado, com separação explícita entre o que foi lido do catálogo, o que foi evidenciado por uso e o que foi inferido, e com a garantia estrutural de que nenhum campo transporta valor de dado do usuário.
- `naming-profiles`: as regras de afixo e de plural que traduzem convenção de nomenclatura, carregadas de arquivo e não de código, com perfis oficiais embarcados e suporte a perfil do usuário.
- `cli-foundation`: o binário, o parsing de flag, e as convenções de saída e de código de saída que os comandos das fases seguintes herdam.

### Modified Capabilities

Nenhuma. É a primeira change do projeto e `openspec/specs/` está vazio.

## Impact

Cria a estrutura inicial do repositório. Não há código existente para quebrar.

Dependências introduzidas no binário: `github.com/spf13/cobra` e `github.com/pelletier/go-toml/v2`. Em teste: `github.com/google/go-cmp`. Nenhuma exige cgo.

Fixa decisões que as fases seguintes herdam e que ficam caras de reverter depois: o module path, o piso de versão do Go, o formato TOML dos perfis, a assinatura das funções de normalização que `internal/infer` vai consumir na fase 3, e os nomes de campo do `internal/model`, que viram contrato do JSON público na fase 7.
