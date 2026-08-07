## Context

O pré-filtro fica entre a pontuação por metadados e a validação contra dados, e existe por causa da assimetria de custo entre as duas: um sinal a mais na fase 3 custa nada; um candidato a mais na fase 5 custa um anti-join num banco de produção.

A matéria-prima aqui é diferente de tudo que a ferramenta consumiu até agora: **estimativas**. `n_distinct` e `histogram_bounds` vêm do último `ANALYZE`, que pode nunca ter rodado, ter rodado há dois anos, ou ter amostrado mal. Todo o desenho decorre de levar isso a sério.

Restrição herdada que domina tudo: `most_common_vals` e `histogram_bounds` são valores reais de tabelas do usuário. A regra 2 do projeto se aplica a eles integralmente.

## Goals / Non-Goals

**Goals**

Eliminar candidatos aritmeticamente impossíveis antes de qualquer I/O em tabela. Penalizar os improváveis sem os esconder. Registrar quando não pôde opinar. Provar por teste que nenhum valor lido sobrevive à camada.

**Non-Goals**

Nenhuma leitura de linha de tabela — a camada consulta `pg_catalog`, e só. Nenhuma confirmação: estatística só derruba ou enfraquece hipótese, nunca fortalece a ponto de confirmar. Nenhuma tentativa de medir idade da estatística nesta fase — a margem larga é a mitigação, e `last_analyze` fica como refinamento futuro se o corpus mostrar necessidade.

## Decisions

### Leitura dirigida, não varredura

A alternativa seria o catálogo da fase 2 carregar `pg_stats` para todas as colunas do escopo. Seria mais simples de orquestrar e erraria duas vezes: pagaria por estatística de colunas que nunca virarão candidato, e manteria valores de dados do usuário em memória durante toda a execução, espalhados pelo modelo que o resto do programa serializa.

A camada busca estatística apenas para as colunas dos candidatos que sobreviveram ao limiar — a filha e a chave do pai — numa consulta única. Os valores vivem num mapa interno do pacote, produzem sinais, e saem de escopo.

### Cardinalidade rejeita além de margem larga; faixa nunca rejeita sozinha

As duas checagens têm naturezas diferentes e regimes diferentes.

Cardinalidade é aritmética sobre duas estimativas: se a filha tem mais valores distintos do que o pai tem linhas, contenção total é impossível. Ainda assim as duas pontas são estimativas, então a penalidade é o caminho padrão e a rejeição direta exige que a violação exceda uma margem de tolerância larga — o dobro, por padrão. Uma coluna com o dobro de valores distintos do que o alvo tem de linhas não é uma estatística velha, é outra relação.

Faixa é mais frágil: os limites do histograma são amostra, e uma tabela que cresceu depois do `ANALYZE` tem valores novos fora dos limites antigos por construção. Limites da filha fora dos limites do pai penalizam com sinal forte e param aí. Rejeitar por faixa seria inventar certeza a partir do pedaço mais impreciso da estatística.

### Faixa só para a família numérica

A comparação de limites precisa da ordem do tipo, e `histogram_bounds` chega como `anyarray`. Para inteiros e numéricos, o servidor converte os limites para `float8` na própria consulta e a comparação é exata o bastante para uma heurística de improbabilidade. Para texto, `uuid` e datas, esta fase não opina — colação e semântica de ordenação viram fonte de falso sinal, e o custo de perder a checagem é só um anti-join a mais na fase 5.

A limitação fica declarada na spec como comportamento, não escondida como detalhe.

### Estatística ausente vira sinal de peso zero, não silêncio

Tabela nunca analisada não tem linha em `pg_stats`, e `n_distinct` zero significa desconhecido. Nos dois casos o pré-filtro não opina — mas "não opinei" precisa ficar registrado, senão a ausência de penalidade fica indistinguível de aprovação.

O registro é um sinal de peso zero no próprio candidato, mais um contador na cobertura. O sinal não move o score; existe para que quem audita um candidato veja que a estatística foi procurada e não estava lá.

### O score é recomposto pelo dono da regra de saturação

Os pesos das penalidades estatísticas moram em `internal/stats`, porque são julgamento sobre estimativas e a auditoria deles é local. Mas a combinação de sinais em score — a saturação em zero e um — tem um único dono, `internal/infer`, e passa a ser exportada para que esta camada recomponha o score sem duplicar a regra.

Duplicar seria o erro clássico: dois lugares que saturam diferente e um limiar que muda de significado dependendo de qual camada tocou o candidato por último.

### Desligável por inteiro

`--no-stats` remove a camada da execução. Não existe meio-termo configurável de pesos por flag — quem discorda da penalidade desliga o pré-filtro e reporta, e o ajuste de peso acontece no código, revisável, como na fase 3.

## Risks / Trade-offs

**Estatística velha derruba candidato legítimo** → a rejeição direta exige violação além do dobro; abaixo disso é penalidade, e candidato penalizado ainda aparece na saída com o sinal visível. O usuário que discorda roda com `--no-stats` e compara.

**`n_distinct` negativo mal resolvido produz absurdo** → a convenção (negativo é razão sobre o total de linhas) já está encapsulada em `model.ColumnStats.EstimatedDistinct`, com tabela de casos própria. Esta camada não reimplementa a conta.

**A consulta a `pg_stats` só vê o que o papel corrente pode ler** → coluna invisível é estatística ausente, e cai no regime de não opinar com registro. A cobertura da fase 2 já lista as tabelas sem privilégio; os dois relatos se somam em vez de se contradizer.

**Vazamento por caminho novo** — a consulta traz `float8` derivados de `histogram_bounds` para a memória → os valores ficam em structs não exportadas, sem tag de serialização, e o teste de vazamento serializa o resultado inteiro da camada e varre contra valores plantados na fixture. O leak test do modelo já prova que os bounds de `ColumnStats` não serializam; o desta fase prova o mesmo do que a camada nova produz.

## Migration Plan

Não aplicável. Campos novos de cobertura são aditivos no JSON.

## Open Questions

Nenhuma bloqueante. A margem de tolerância padrão (2x) é estimativa declarada como tal e será revista com o corpus da fase 8, pelo mesmo caminho do limiar de score.
