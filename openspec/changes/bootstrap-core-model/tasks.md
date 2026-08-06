## 1. Esqueleto do repositório

- [ ] 1.1 Criar `go.mod` com module `github.com/lvcas-dotcom/pgfathom` e piso de versão Go 1.25
- [ ] 1.2 Criar `.gitignore` cobrindo binários, `dist/`, e artefatos de cobertura
- [ ] 1.3 Criar `Makefile` com alvos `build`, `test`, `test-integration`, `lint`, `fmt` e `cover`
- [ ] 1.4 Criar `.golangci.yml` com conjunto modesto de linters, incluindo `errcheck`, `govet`, `staticcheck` e `revive`
- [ ] 1.5 Criar workflow de CI em `.github/workflows/ci.yml` com job rápido (unit, sem Docker) em Linux, macOS e Windows
- [ ] 1.6 Adicionar ao CI um job de lint separado do job de teste

## 2. Modelo semântico

- [ ] 2.1 Criar `internal/model` com `Schema`, `Table`, `Column`, `ColumnRef` e `Index`
- [ ] 2.2 Adicionar `ForeignKey` com o campo de estado de validação vindo de `convalidated` e o indicador de índice no lado filho
- [ ] 2.3 Adicionar `TableStats` com os contadores de uso e o timestamp de reset como campo obrigatório e anulável
- [ ] 2.4 Adicionar `ColumnStats` com os campos de `pg_stats`, mantendo limites de histograma e valores comuns em campos **não exportados**
- [ ] 2.5 Adicionar `SignalKind` como enum fechado de origens e o tipo `Signal` com peso e detalhe restrito a nome de objeto
- [ ] 2.6 Adicionar `Candidate` com os campos de proveniência, score e motivo de descarte
- [ ] 2.7 Adicionar `Validation` com contenção por linha e por valor, contagens de órfãos nas duas dimensões, máximo de linhas por valor, método e duração
- [ ] 2.8 Adicionar `Verdict` como enum fechado incluindo o veredito de não avaliado
- [ ] 2.9 Adicionar `Finding` para achados que não dependem de inferência
- [ ] 2.10 Adicionar `Coverage` com totais, listas de tabelas puladas por motivo e contadores de candidato
- [ ] 2.11 Escrever teste que inspeciona o grafo de importação de `internal/model` e falha se houver dependência de camada do projeto, de `pgx` ou de pacote de I/O
- [ ] 2.12 Escrever teste que serializa todas as estruturas do modelo preenchidas e falha se qualquer campo de valor de dado aparecer no JSON

## 3. Perfis de nomenclatura

- [ ] 3.1 Definir o schema TOML do perfil: sufixos de coluna, prefixos de coluna, prefixos de tabela e lista ordenada de regras de plural
- [ ] 3.2 Criar `internal/profile` com o tipo `Profile` e o carregamento via `pelletier/go-toml/v2`
- [ ] 3.3 Implementar carregamento por nome, resolvendo contra os perfis embarcados por `embed`
- [ ] 3.4 Implementar carregamento por caminho de arquivo, com erro descritivo para arquivo ausente, TOML inválido ou campo obrigatório faltando
- [ ] 3.5 Escrever `profiles/pt-br.toml` cobrindo formas acentuadas e sem acento, com regras da mais específica para a mais genérica e queda de `s` por último
- [ ] 3.6 Escrever `profiles/en.toml`
- [ ] 3.7 Escrever `profiles/es.toml`
- [ ] 3.8 Implementar extração de nome de entidade a partir de nome de coluna, testando sufixos antes de prefixos, na ordem declarada, primeira regra vencendo
- [ ] 3.9 Garantir que coluna sem afixo devolve o próprio nome como candidato
- [ ] 3.10 Implementar normalização de nome de tabela devolvendo conjunto ordenado de formas, sempre incluindo o nome original
- [ ] 3.11 Fazer a normalização reportar, para cada forma, se veio do original, de remoção de prefixo, de despluralização ou de ambos
- [ ] 3.12 Escrever tabela de casos para `pt-br` cobrindo os plurais de `docs/ROADMAP.md`: `opcoes`, `animais`, `responsaveis`, `perfis`, `armazens`, `meses`, `órgãos`, mais o caso ambíguo `logins`
- [ ] 3.13 Escrever tabela de casos para `en` e `es`
- [ ] 3.14 Escrever teste que carrega os três perfis embarcados e valida que nenhum tem regra duplicada ou ordem inconsistente

## 4. Esqueleto da CLI

- [ ] 4.1 Criar `cmd/pgfathom` com a raiz cobra, nome, descrição curta e descrição longa em inglês
- [ ] 4.2 Implementar `version` com versão, commit e data de build injetados por `-ldflags`
- [ ] 4.3 Fazer execução sem argumento imprimir ajuda e sair com 0
- [ ] 4.4 Definir os códigos de saída como constantes documentadas: 0 sucesso, 1 falha, 2 uso incorreto, 3 reservado para achados
- [ ] 4.5 Implementar a disciplina de fluxo: resultado em stdout, diagnóstico em stderr, sem exceção
- [ ] 4.6 Implementar detecção de TTY em stdout, com respeito a `NO_COLOR` e a flag explícita de sobreposição
- [ ] 4.7 Configurar `log/slog` em stderr com nível controlado por flag global
- [ ] 4.8 Criar o `context.Context` raiz cancelável em `SIGINT` e `SIGTERM`, e propagá-lo a todo comando
- [ ] 4.9 Escrever teste de código de saída para sucesso, uso incorreto e subcomando desconhecido
- [ ] 4.10 Escrever teste que executa com stdout capturado e verifica ausência de sequência ANSI

## 5. Infraestrutura de teste

- [ ] 5.1 Adicionar `github.com/google/go-cmp` como dependência de teste
- [ ] 5.2 Criar o helper de golden file com flag `-update`, sem dependência externa
- [ ] 5.3 Criar a build tag `//go:build integration` e o arquivo de exemplo vazio, garantindo que `go test ./...` não a inclui
- [ ] 5.4 Verificar que `go test ./...` completa sem Docker e sem rede
- [ ] 5.5 Adicionar alvo de cobertura ao Makefile e reportar a cobertura de `internal/profile` no CI

## 6. Fechamento

- [ ] 6.1 Rodar `golangci-lint run` e zerar os apontamentos
- [ ] 6.2 Verificar `go build` cruzado para linux/amd64, linux/arm64, darwin/arm64 e windows/amd64 sem cgo
- [ ] 6.3 Conferir que `go.mod` tem exatamente as dependências previstas no design, sem transitiva inesperada no binário
- [ ] 6.4 Atualizar `README.md` com instrução de build e de execução dos testes
- [ ] 6.5 Rodar `openspec validate bootstrap-core-model` e corrigir o que apontar
