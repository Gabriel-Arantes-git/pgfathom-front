## ADDED Requirements

### Requirement: O pré-filtro lê apenas catálogo, e apenas o necessário

A camada de pré-filtro SHALL consultar exclusivamente `pg_stats` e visões de catálogo, e SHALL restringir a consulta às colunas dos candidatos que sobreviveram ao corte por limiar.

Nenhuma linha de tabela do usuário é lida. O custo da camada precisa permanecer desprezível para que ela nunca seja motivo de desligar a ferramenta em produção.

#### Scenario: Nenhuma relação fora do catálogo

- **WHEN** as consultas da camada são inspecionadas
- **THEN** toda relação referenciada pertence a `pg_catalog` ou a visões de estatística do servidor

#### Scenario: Consulta dirigida

- **WHEN** o pré-filtro roda sobre um conjunto de candidatos
- **THEN** apenas as colunas envolvidas nesses candidatos são consultadas em `pg_stats`

### Requirement: Cardinalidade impossível penaliza, e rejeita só além da margem

Quando o número estimado de valores distintos da coluna filha excede o número estimado de linhas da tabela pai, o candidato SHALL receber o sinal negativo de violação de cardinalidade e ter o score recomposto. A rejeição direta SHALL ocorrer apenas quando a violação excede uma margem de tolerância larga, e SHALL registrar o motivo com as duas estimativas declaradas como estimativas.

As duas pontas da conta vêm do `ANALYZE` e podem estar velhas. Penalidade tem recurso — o candidato segue visível com o sinal — enquanto rejeição é final, e por isso exige folga.

#### Scenario: Violação dentro da margem penaliza

- **WHEN** a filha tem mais valores distintos estimados que as linhas do pai, mas menos que o dobro
- **THEN** o candidato carrega o sinal de violação de cardinalidade, tem score reduzido e não é rejeitado

#### Scenario: Violação além da margem rejeita

- **WHEN** os valores distintos estimados da filha excedem o dobro das linhas estimadas do pai
- **THEN** o candidato é rejeitado com motivo que nomeia as estimativas envolvidas

#### Scenario: Rejeição estatística é reportável

- **WHEN** o usuário pede para ver os descartados
- **THEN** os rejeitados pela estatística aparecem com o motivo

### Requirement: Faixa penaliza e nunca rejeita sozinha

Para colunas da família numérica, quando os limites do histograma da filha caem fora dos limites da chave do pai, o candidato SHALL receber o sinal negativo de violação de faixa. Este sinal MUST NOT causar rejeição por si só, qualquer que seja a magnitude.

Limites de histograma são amostra, e tabela que cresceu após o `ANALYZE` tem valores fora dos limites antigos por construção. Para tipos fora da família numérica a checagem de faixa SHALL ser omitida.

#### Scenario: Faixa deslocada penaliza

- **WHEN** os limites da filha caem inteiramente fora dos limites da chave do pai
- **THEN** o candidato carrega o sinal de violação de faixa e tem score reduzido, sem ser rejeitado

#### Scenario: Tipo não numérico não opina sobre faixa

- **WHEN** o candidato liga colunas de tipo textual, `uuid` ou data
- **THEN** nenhum sinal de faixa é emitido, em nenhuma direção

### Requirement: Estatística ausente não gera opinião, e o silêncio fica registrado

Quando a estatística de uma coluna envolvida está ausente ou não interpretável, o pré-filtro MUST NOT emitir penalidade nem rejeição para o candidato, e SHALL registrar que não pôde opinar — no próprio candidato, por sinal de peso zero, e na cobertura, por contagem.

Inventar rejeição a partir de dado que não existe é a definição de falso negativo evitável. E ausência de penalidade sem registro seria indistinguível de aprovação.

#### Scenario: Tabela nunca analisada

- **WHEN** uma das colunas do candidato não tem linha em `pg_stats`
- **THEN** o candidato atravessa o pré-filtro sem alteração de score e carrega o sinal de estatística indisponível

#### Scenario: Cobertura contabiliza o silêncio

- **WHEN** o pré-filtro termina
- **THEN** a cobertura informa quantos candidatos não puderam ser avaliados por falta de estatística

### Requirement: Nenhum valor de pg_stats sobrevive à camada

Os valores lidos de `pg_stats` — incluindo `most_common_vals` e `histogram_bounds` e qualquer derivado deles — MUST NOT aparecer em struct serializável, saída de terminal, JSON, log ou mensagem de erro. Eles produzem sinais e contagens em memória e morrem ali.

São valores reais de tabelas do usuário, lidos do catálogo em vez da tabela, o que não muda nada sobre o que são.

#### Scenario: Serialização completa é limpa

- **WHEN** o resultado inteiro da camada é serializado e varrido contra valores plantados na fixture
- **THEN** nenhum valor plantado aparece

#### Scenario: Motivo de rejeição carrega contagens, nunca valores

- **WHEN** um candidato é rejeitado por cardinalidade
- **THEN** o motivo contém estimativas de contagem e nomes de objeto, e nenhum valor de dado

### Requirement: O pré-filtro é desligável por inteiro

O comando `discover` SHALL aceitar uma flag que remove a camada estatística da execução, e a saída SHALL declarar quando o pré-filtro não rodou.

Penalidade derivada de estimativa só é auditável se der para comparar a execução com e sem ela.

#### Scenario: Desligado não toca nos candidatos

- **WHEN** `discover` roda com o pré-filtro desligado
- **THEN** nenhum candidato carrega sinal estatístico e a cobertura declara que a camada não rodou

### Requirement: O funil entra na cobertura

A cobertura SHALL informar quantos candidatos o pré-filtro avaliou, quantos rejeitou e quantos não pôde avaliar.

É o número que prova que a camada paga o próprio custo — e o que a fase 5 usa para prever quantos anti-joins uma execução vai disparar.

#### Scenario: Funil visível

- **WHEN** `discover` termina com o pré-filtro ligado
- **THEN** a cobertura informa avaliados, rejeitados por estatística e sem estatística
