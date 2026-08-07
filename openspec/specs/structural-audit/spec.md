# structural-audit Specification

## Purpose
TBD - created by archiving change catalog-inspection. Update Purpose after archive.
## Requirements
### Requirement: O comando audit não depende de inferência

`pgfathom audit` SHALL reportar apenas achados derivados diretamente do catálogo, sem heurística de nome, sem pontuação e sem leitura de dado.

Isso torna o comando determinístico e imune a falso positivo: todo achado que ele emite é um fato do catálogo, não uma hipótese. É também o que permite executá-lo em banco onde a inferência não teria nada a dizer.

#### Scenario: Nenhum candidato inferido na saída

- **WHEN** `pgfathom audit` é executado
- **THEN** a saída não contém candidato, score nem veredito de inferência

### Requirement: Constraint declarada mas nunca validada é reportada

O comando SHALL reportar toda constraint de chave estrangeira cujo estado `convalidated` seja falso.

O achado SHALL identificar a constraint, a tabela e a estimativa de linhas da tabela filha, e SHALL explicar que a constraint bloqueia violações novas mas nunca verificou as linhas preexistentes.

#### Scenario: NOT VALID encontrada

- **WHEN** o schema contém uma FK criada com `NOT VALID` e nunca validada
- **THEN** o achado correspondente aparece na saída, identificando a constraint e sua tabela

#### Scenario: Constraint validada não gera achado

- **WHEN** todas as FKs do schema foram validadas
- **THEN** nenhum achado desse tipo é emitido

### Requirement: Chave estrangeira sem índice do lado filho é reportada

O comando SHALL reportar toda FK declarada sem índice utilizável na coluna filha, acompanhada da estimativa de linhas de ambas as tabelas, que é o que indica a gravidade.

Sem esse índice, todo `DELETE` na tabela pai vira varredura sequencial da filha.

#### Scenario: FK sem índice

- **WHEN** uma FK declarada não tem índice com a coluna filha em posição inicial
- **THEN** o achado aparece na saída com as estimativas de linha de pai e filha

### Requirement: Saída em terminal e em JSON

O comando SHALL suportar saída em tabela no terminal e em JSON, selecionável por flag.

A saída em terminal SHALL agrupar os achados por tipo e SHALL terminar com o bloco de cobertura. A saída em JSON SHALL seguir o contrato versionado do modelo.

Resultado vai para stdout; aviso, progresso e erro vão para stderr, sem exceção.

#### Scenario: JSON consumível

- **WHEN** `pgfathom audit --format json` é executado com stdout redirecionado
- **THEN** o arquivo contém JSON válido e nada além dele

#### Scenario: Cobertura sempre presente

- **WHEN** o comando termina, em qualquer formato
- **THEN** o bloco de cobertura consta da saída

#### Scenario: Ausência de achado é afirmativa

- **WHEN** nenhum achado é encontrado e todas as tabelas do escopo foram analisadas
- **THEN** a saída afirma que o escopo foi analisado e está limpo, em vez de apenas não listar nada

### Requirement: Nenhum valor de dado do usuário na saída

Nenhuma saída do comando, em qualquer formato ou nível de log, SHALL conter valor lido de tabela do usuário. O que sai são nomes de objeto de catálogo, contagens e proporções.

#### Scenario: Varredura da saída

- **WHEN** o comando é executado contra as fixtures de teste, que contêm valores reconhecíveis plantados
- **THEN** nenhum desses valores aparece na saída em terminal, no JSON ou no log

### Requirement: Código de saída reflete execução, não achados

O comando SHALL sair com zero quando a execução completa, mesmo tendo encontrado achados. Achado é o produto do comando, não uma falha dele.

Falha de conexão, de privilégio ou erro interno SHALL sair com o código de falha; erro de linha de comando, com o código de uso.

#### Scenario: Achados não alteram o código de saída

- **WHEN** o comando encontra achados e completa normalmente
- **THEN** o código de saída é zero

#### Scenario: Falha de conexão

- **WHEN** a conexão não pode ser estabelecida
- **THEN** o código de saída é o de falha e a mensagem vai para stderr

