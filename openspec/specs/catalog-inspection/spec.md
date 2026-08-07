# catalog-inspection Specification

## Purpose
TBD - created by archiving change catalog-inspection. Update Purpose after archive.
## Requirements
### Requirement: Leitura completa da estrutura declarada

A ferramenta SHALL ler do catálogo, para cada schema no escopo: tabelas, colunas com tipo formatado e tipo base normalizado, nulabilidade, valor padrão e posição, chaves primárias em ordem, constraints unique, índices com suas colunas em ordem, e comentários de tabela e de coluna.

O tipo base SHALL ser normalizado para comparação. Comparar tipos formatados diretamente produz falso negativo entre grafias equivalentes do mesmo tipo.

#### Scenario: Estrutura básica

- **WHEN** um schema com tabelas, chaves e índices é lido
- **THEN** o modelo resultante contém todas as tabelas do escopo com colunas, chaves, índices e comentários

#### Scenario: Tipos equivalentes normalizam igual

- **WHEN** uma coluna é declarada como `bigint` e outra como `int8`
- **THEN** ambas produzem o mesmo tipo base no modelo

### Requirement: Estado de validação da chave estrangeira é preservado

Toda chave estrangeira lida SHALL carregar o valor de `pg_constraint.convalidated`.

Uma constraint criada `NOT VALID` e nunca validada bloqueia violações novas mas nunca verificou as linhas preexistentes, enquanto aparece idêntica a qualquer outra em `\d` e em ferramenta de diagrama. Descartar esse campo na leitura elimina a informação que sustenta um achado inteiro.

A ferramenta SHALL também determinar, para cada chave estrangeira, se existe índice utilizável na coluna filha — ou seja, um índice com a coluna em posição inicial.

#### Scenario: Constraint não validada é distinguível

- **WHEN** o schema contém uma FK criada com `NOT VALID` e nunca validada
- **THEN** o modelo a marca como não validada

#### Scenario: Índice do lado filho é detectado

- **WHEN** uma FK tem índice cuja coluna inicial é a coluna filha
- **THEN** a FK é marcada como indexada

#### Scenario: Índice em posição não inicial não conta

- **WHEN** a coluna filha aparece num índice composto, mas não em posição inicial
- **THEN** a FK é marcada como não indexada, porque esse índice não serve para a busca

### Requirement: Estatística de uso lida junto do timestamp de reset

A ferramenta SHALL ler as estatísticas de uso de tabela sempre acompanhadas do momento do último reset de estatística do banco.

Contador de uso sem esse timestamp não tem significado. Quando o momento do reset não puder ser determinado, os contadores SHALL ser marcados como não interpretáveis, e nenhum achado pode ser derivado deles.

#### Scenario: Contadores acompanhados do reset

- **WHEN** as estatísticas de uso são lidas
- **THEN** o timestamp de reset consta do modelo, ou o estado é explicitamente marcado como desconhecido

### Requirement: Nenhum dado de tabela do usuário é lido nesta fase

A camada de catálogo MUST NOT emitir consulta que leia linha de tabela do usuário. Ela lê exclusivamente `pg_catalog`, `information_schema` e as visões de estatística.

Leitura de dado começa na validação, numa fase posterior e numa camada separada. Manter essa fronteira explícita é o que permite afirmar, sobre esta fase, que nenhum valor pode vazar porque nenhum valor é lido.

#### Scenario: Consultas restritas ao catálogo

- **WHEN** um `audit` completo é executado
- **THEN** nenhuma consulta emitida referencia uma tabela do usuário

### Requirement: Tabela pulada é registrada, nunca silenciada

Tabela que não puder ser analisada SHALL ser registrada na cobertura com o motivo, e a execução SHALL continuar. Os motivos SHALL cobrir: falta de privilégio `SELECT`, chave primária composta, ausência de chave primária, tabela particionada e herança de tabela.

Tabela particionada SHALL ser lida a partir da tabela pai, e as partições MUST NOT ser iteradas separadamente.

#### Scenario: Falta de privilégio

- **WHEN** o papel conectado não tem `SELECT` numa tabela do escopo
- **THEN** a tabela aparece na lista de puladas por privilégio, e a execução prossegue

#### Scenario: Forma não suportada

- **WHEN** uma tabela tem chave primária composta
- **THEN** ela aparece na cobertura com o motivo correspondente

#### Scenario: Partições não são iteradas

- **WHEN** o schema contém uma tabela particionada com várias partições
- **THEN** o modelo contém a tabela pai e não contém uma entrada por partição

### Requirement: Escopo controlável por schema e por padrão de exclusão

A ferramenta SHALL aceitar a lista de schemas a analisar, com padrão `public`, e uma lista de padrões de tabela a excluir.

Tabela removida por filtro do usuário SHALL ser registrada em campo distinto das puladas por privilégio ou por forma não suportada. Exclusão pedida e exclusão forçada são coisas diferentes e não podem se confundir no relatório.

#### Scenario: Exclusão por padrão

- **WHEN** um padrão de exclusão casa com algumas tabelas
- **THEN** essas tabelas não são analisadas e aparecem na cobertura como excluídas

### Requirement: Cobertura preenchida em toda execução

O resultado SHALL conter a estrutura de cobertura preenchida, com o total de tabelas do escopo, o total analisado e todas as listas de puladas.

#### Scenario: Contagens fecham

- **WHEN** uma execução termina
- **THEN** o total analisado somado às tabelas puladas por todos os motivos é igual ao total do escopo

