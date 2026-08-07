## Why

A fase 3 entrega candidatos pontuados por fatos do catálogo, e a fase 5 vai pagar um anti-join por cada um que sobreviver ao limiar. Entre as duas existe uma fonte de informação que ainda não foi usada e que custa quase nada: as estatísticas do planner. Se a coluna filha tem mais valores distintos estimados do que a tabela pai tem linhas, contenção total é aritmeticamente impossível — e esse candidato não precisa custar uma query em produção para morrer.

O pré-filtro é a última camada barata antes da única camada cara. Cada candidato que ele elimina é um anti-join que a fase 5 não dispara contra o banco de alguém.

## What Changes

- `internal/stats`: leitura dirigida de `pg_stats` e consumo do `reltuples` que o catálogo já carrega, restrita às colunas dos candidatos sobreviventes ao limiar.
- Checagem de cardinalidade: mais valores distintos na filha do que linhas na pai penaliza forte, e rejeita direto apenas quando a violação excede uma margem de tolerância larga.
- Checagem de faixa para a família numérica: limites do histograma da filha fora dos limites da chave do pai penalizam, nunca rejeitam sozinhos.
- Estatística ausente: o pré-filtro não opina e registra que não opinou, por candidato e na cobertura.
- Flag para desligar o pré-filtro por inteiro, porque penalidade vinda de estimativa precisa ser desligável para ser auditável.
- Cobertura ganha o funil: quantos candidatos a estatística rejeitou e quantos ela não pôde avaliar.

Os valores lidos de `pg_stats` — em especial `histogram_bounds` — **são dados do usuário**. Eles produzem um número em memória e morrem ali; nenhum campo, log, erro ou JSON os carrega.

## Capabilities

### New Capabilities

- `stats-prefilter`: o que as estimativas do planner podem e não podem fazer com um candidato — as duas checagens, o regime de penalidade contra rejeição, o silêncio obrigatório diante de estatística ausente, e a proibição absoluta de vazamento dos valores lidos.

### Modified Capabilities

- `candidate-scoring`: o score passa a poder ser recomposto após a geração — os sinais estatísticos entram com peso negativo e o score é recalculado pelo mesmo mecanismo de saturação, para que o limiar continue significando a mesma coisa.

## Impact

Nenhuma dependência nova. `internal/stats` importa `internal/model`, `internal/db` e `internal/infer` — este último apenas pela função exportada de recomposição de score, para que a regra de saturação continue tendo um único dono.

A separação entre esta camada e a fase 3 é o ponto do desenho: a fase 3 opera sobre fatos do catálogo, esta opera sobre estimativas do `ANALYZE`, que podem estar velhas ou grosseiras. Misturadas, a pontuação viraria uma sopa onde ninguém sabe o que é fato e o que é palpite. Separadas, a penalidade estatística é identificável nos sinais, auditável e desligável.

O JSON ganha campos aditivos na cobertura. Nenhum campo existente muda de forma.
