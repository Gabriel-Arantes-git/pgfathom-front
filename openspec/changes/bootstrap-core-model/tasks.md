## 1. Esqueleto do repositório

- [x] 1.1 Criar `go.mod` com module `github.com/lvcas-dotcom/pgfathom` e piso de versão Go 1.25
- [x] 1.2 Criar `.gitignore` cobrindo binários, `dist/`, e artefatos de cobertura
- [x] 1.3 Criar `Makefile` com alvos `build`, `test`, `test-integration`, `lint`, `fmt` e `cover`
- [x] 1.4 Criar `.golangci.yml` com conjunto modesto de linters, incluindo `errcheck`, `govet`, `staticcheck` e `revive`
- [x] 1.5 Criar workflow de CI em `.github/workflows/ci.yml` com job rápido (unit, sem Docker) em Linux, macOS e Windows
- [x] 1.6 Adicionar ao CI um job de lint separado do job de teste

## 2. Modelo semântico

- [x] 2.1 Criar `internal/model` com `Schema`, `Table`, `Column`, `ColumnRef` e `Index`
- [x] 2.2 Adicionar `ForeignKey` com o campo de estado de validação vindo de `convalidated` e o indicador de índice no lado filho
- [x] 2.3 Adicionar `TableStats` com os contadores de uso e o timestamp de reset como campo obrigatório e anulável
- [x] 2.4 Adicionar `ColumnStats` com os campos de `pg_stats`, mantendo limites de histograma e valores comuns em campos **não exportados**
- [x] 2.5 Adicionar `SignalKind` como enum fechado de origens e o tipo `Signal` com peso e detalhe restrito a nome de objeto
- [x] 2.6 Adicionar `Candidate` com os campos de proveniência, score e motivo de descarte
- [x] 2.7 Adicionar `Validation` com contenção por linha e por valor, contagens de órfãos nas duas dimensões, máximo de linhas por valor, método e duração
- [x] 2.8 Adicionar `Verdict` como enum fechado incluindo o veredito de não avaliado
- [x] 2.9 Adicionar `Finding` para achados que não dependem de inferência
- [x] 2.10 Adicionar `Coverage` com totais, listas de tabelas puladas por motivo e contadores de candidato
- [x] 2.11 Escrever teste que inspeciona o grafo de importação de `internal/model` e falha se houver dependência de camada do projeto, de `pgx` ou de pacote de I/O
- [x] 2.12 Escrever teste que serializa todas as estruturas do modelo preenchidas e falha se qualquer campo de valor de dado aparecer no JSON

## 3. Perfis de nomenclatura

- [x] 3.1 Definir o schema TOML do perfil: sufixos de coluna, prefixos de coluna, prefixos de tabela e lista ordenada de regras de plural
- [x] 3.2 Criar `internal/profile` com o tipo `Profile` e o carregamento via `pelletier/go-toml/v2`
- [x] 3.3 Implementar carregamento por nome, resolvendo contra os perfis embarcados por `embed`
- [x] 3.4 Implementar carregamento por caminho de arquivo, com erro descritivo para arquivo ausente, TOML inválido ou campo obrigatório faltando
- [x] 3.5 Escrever `profiles/pt-br.toml` cobrindo formas acentuadas e sem acento, com regras da mais específica para a mais genérica e queda de `s` por último
- [x] 3.6 Escrever `profiles/en.toml`
- [x] 3.7 Escrever `profiles/es.toml`
- [x] 3.8 Implementar extração de nome de entidade a partir de nome de coluna, testando sufixos antes de prefixos, na ordem declarada, primeira regra vencendo
- [x] 3.9 Garantir que coluna sem afixo devolve o próprio nome como candidato
- [x] 3.10 Implementar normalização de nome de tabela devolvendo conjunto ordenado de formas, sempre incluindo o nome original
- [x] 3.11 Fazer a normalização reportar, para cada forma, se veio do original, de remoção de prefixo, de despluralização ou de ambos
- [x] 3.12 Escrever tabela de casos para `pt-br` cobrindo os plurais de `docs/ROADMAP.md`: `opcoes`, `animais`, `responsaveis`, `perfis`, `armazens`, `meses`, `órgãos`, mais o caso ambíguo `logins`
- [x] 3.13 Escrever tabela de casos para `en` e `es`
- [x] 3.14 Escrever teste que carrega os três perfis embarcados e valida que nenhum tem regra duplicada ou ordem inconsistente

## 4. Esqueleto da CLI

- [x] 4.1 Criar `cmd/pgfathom` com a raiz cobra, nome, descrição curta e descrição longa em inglês
- [x] 4.2 Implementar `version` com versão, commit e data de build injetados por `-ldflags`
- [x] 4.3 Fazer execução sem argumento imprimir ajuda e sair com 0
- [x] 4.4 Definir os códigos de saída como constantes documentadas: 0 sucesso, 1 falha, 2 uso incorreto, 3 reservado para achados
- [x] 4.5 Implementar a disciplina de fluxo: resultado em stdout, diagnóstico em stderr, sem exceção
- [x] 4.6 Implementar detecção de TTY em stdout, com respeito a `NO_COLOR` e a flag explícita de sobreposição
- [x] 4.7 Configurar `log/slog` em stderr com nível controlado por flag global
- [x] 4.8 Criar o `context.Context` raiz cancelável em `SIGINT` e `SIGTERM`, e propagá-lo a todo comando
- [x] 4.9 Escrever teste de código de saída para sucesso, uso incorreto e subcomando desconhecido
- [x] 4.10 Escrever teste que executa com stdout capturado e verifica ausência de sequência ANSI

## 5. Infraestrutura de teste

- [x] 5.1 Adicionar `github.com/google/go-cmp` como dependência de teste
- [x] 5.2 Criar o helper de golden file com flag `-update`, sem dependência externa
- [x] 5.3 Criar a build tag `//go:build integration` e o arquivo de exemplo vazio, garantindo que `go test ./...` não a inclui
- [x] 5.4 Verificar que `go test ./...` completa sem Docker e sem rede
- [x] 5.5 Adicionar alvo de cobertura ao Makefile e reportar a cobertura de `internal/profile` no CI

## 6. Fechamento

- [x] 6.1 Rodar `golangci-lint run` e zerar os apontamentos
- [x] 6.2 Verificar `go build` cruzado para linux/amd64, linux/arm64, darwin/arm64 e windows/amd64 sem cgo
- [x] 6.3 Conferir que `go.mod` tem exatamente as dependências previstas no design, sem transitiva inesperada no binário
- [x] 6.4 Atualizar `README.md` com instrução de build e de execução dos testes
- [x] 6.5 Rodar `openspec validate bootstrap-core-model` e corrigir o que apontar
