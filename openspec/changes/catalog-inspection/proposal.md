## Why

A fase 1 entregou tipos e regras que ninguém pode executar contra nada. O modelo existe, os perfis funcionam, e o binário só sabe imprimir a própria versão.

Esta fase é onde a ferramenta encosta num banco pela primeira vez, e por isso é onde todo o risco operacional entra de uma vez: conexão, privilégio, timeout, carga em servidor de terceiro. Tratar isso cedo, enquanto ainda não há inferência nem validação por cima, significa que quando algo quebrar dá para saber que quebrou aqui.

Ela também entrega o primeiro valor real. O comando `audit` reporta dois achados que **não dependem de inferência nenhuma** — constraint declarada `NOT VALID` e nunca validada, e chave estrangeira sem índice do lado filho. São determinísticos, saem direto do catálogo, custam quase nada de código e são imunes a falso positivo. Um banco onde a inferência não encontrar nada ainda assim tem o que auditar, e uma ferramenta que sempre tem algo a dizer tem chance de ser executada uma segunda vez.

## What Changes

- `internal/db`: pool `pgx`, políticas de segurança da sessão concentradas no `AfterConnect`, precedência de credencial, checagem de versão do servidor e de privilégio de escrita do papel conectado.
- `internal/catalog`: leitura de tabelas, colunas, tipos, chaves primárias, uniques, índices, comentários, chaves estrangeiras declaradas com o estado `convalidated`, e estatísticas de uso acompanhadas do timestamp de reset.
- Preenchimento do bloco de `Coverage`, incluindo as tabelas puladas por falta de privilégio e por forma não suportada.
- `internal/report`: renderizador de terminal mínimo, o suficiente para o `audit`. O relatório completo é da fase 7.
- Comando `pgfathom audit` com os dois achados estruturais, saída em terminal e em JSON.
- Infraestrutura de teste de integração com `testcontainers-go` atrás da tag `integration`, e as primeiras fixtures SQL em `testdata/`.
- `LICENSE` Apache-2.0, resolvendo a questão em aberto registrada no design da fase 1.

**Piso de PostgreSQL fixado em 13.** Cobre todas as versões com suporte da comunidade e evita caminho condicional no código de catálogo.

## Capabilities

### New Capabilities

- `database-session`: a conexão e tudo que a torna segura de apontar para produção alheia — modo somente-leitura, identificação em `pg_stat_activity`, os três timeouts, precedência de credencial e verificação de versão e de privilégio.
- `catalog-inspection`: a leitura do catálogo do sistema para dentro do modelo, incluindo o que precisa ser pulado e registrado em vez de analisado.
- `structural-audit`: o comando `audit` e os achados que não dependem de inferência.

### Modified Capabilities

Nenhuma. As capabilities da fase 1 — `domain-model`, `naming-profiles`, `cli-foundation` — são consumidas sem alteração de requisito.

## Impact

Primeira dependência de rede e primeira dependência de Docker no projeto.

Dependências novas no binário: `github.com/jackc/pgx/v5`. Em teste: `github.com/testcontainers/testcontainers-go` e o módulo `postgres`, ambos atrás da tag `integration`, de modo que `go test ./...` continua sem Docker e sem rede.

`internal/db` e `internal/catalog` são as únicas camadas desta fase que falam com o servidor. `internal/catalog` lê exclusivamente catálogo e estatística — nenhuma linha de tabela de usuário é lida nesta fase, o que só acontece na fase 5.

Fixa decisões que as fases seguintes herdam: a superfície de flags de conexão e escopo, a forma como a `Coverage` é montada, e o formato do JSON de saída, que a partir daqui é contrato público.
