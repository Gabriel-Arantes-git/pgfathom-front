## ADDED Requirements

### Requirement: Binário único com subcomandos

O projeto SHALL produzir um binário chamado `pgfathom`, sem dependência de runtime externo e sem cgo, com subcomandos. Nesta change SHALL existir o subcomando `version`; os demais chegam nas fases seguintes.

`pgfathom` sem argumento SHALL exibir a ajuda e sair com código zero.

#### Scenario: Versão

- **WHEN** `pgfathom version` é executado
- **THEN** a versão, o commit e a data de build são impressos em stdout e o código de saída é 0

#### Scenario: Sem argumento

- **WHEN** `pgfathom` é executado sem argumento
- **THEN** a ajuda é impressa e o código de saída é 0

#### Scenario: Subcomando desconhecido

- **WHEN** `pgfathom naoexiste` é executado
- **THEN** a mensagem de erro vai para stderr e o código de saída é 2

### Requirement: Separação entre resultado e diagnóstico

Resultado destinado a consumo — tabela, JSON, SQL — SHALL ser escrito em stdout. Diagnóstico, aviso, progresso e erro SHALL ser escritos em stderr.

Esta separação SHALL valer sem exceção, porque a saída da ferramenta será canalizada para arquivo e para pipeline de CI, e diagnóstico misturado em stdout corrompe o consumo programático.

#### Scenario: Redirecionamento preserva o resultado

- **WHEN** um comando é executado com stdout redirecionado para arquivo
- **THEN** o arquivo contém apenas o resultado, sem nenhuma linha de progresso ou aviso

#### Scenario: Erro não polui stdout

- **WHEN** um comando falha
- **THEN** a mensagem de erro está em stderr e stdout está vazio ou contém apenas resultado parcial válido

### Requirement: Formatação adapta-se ao destino

O sistema SHALL detectar se stdout é um terminal interativo. Quando não for, MUST NOT emitir sequências ANSI de cor ou de controle.

Detecção SHALL respeitar também as convenções de ambiente: `NO_COLOR` definido desliga cor, e uma flag explícita de cor sobrepõe a detecção automática.

#### Scenario: Saída em pipe é limpa

- **WHEN** a saída é canalizada para outro processo ou para arquivo
- **THEN** nenhuma sequência ANSI é emitida

#### Scenario: NO_COLOR é respeitado

- **WHEN** a variável de ambiente `NO_COLOR` está definida e stdout é um terminal
- **THEN** nenhuma sequência ANSI é emitida

### Requirement: Códigos de saída são estáveis e documentados

Os códigos de saída SHALL ser estáveis desde este release, porque serão consumidos por pipeline de CI a partir da fase de `check`.

O contrato SHALL ser: 0 para execução bem-sucedida, 1 para falha de execução — conexão, privilégio, erro interno —, 2 para uso incorreto da linha de comando, e 3 reservado para "achados presentes" nos comandos que precisarem sinalizar regressão sem caracterizar falha.

#### Scenario: Sucesso

- **WHEN** um comando completa sem erro
- **THEN** o código de saída é 0

#### Scenario: Uso incorreto

- **WHEN** uma flag desconhecida ou um valor inválido é fornecido
- **THEN** o código de saída é 2 e a mensagem indica o uso correto

### Requirement: Cancelamento é honrado

Todo comando SHALL propagar um `context.Context` cancelável a partir da raiz, e SHALL cancelá-lo em `SIGINT` e `SIGTERM`.

Interrupção MUST NOT deixar trabalho pendurado no servidor de banco. Esta requisição existe nesta change, antes de haver qualquer acesso a banco, precisamente para que nenhuma camada posterior seja escrita sem receber contexto.

#### Scenario: Interrupção encerra limpo

- **WHEN** o processo recebe `SIGINT` durante a execução
- **THEN** o contexto raiz é cancelado, o processo encerra sem panic, e o código de saída indica interrupção

### Requirement: Logging estruturado sem vazamento

O sistema SHALL usar `log/slog` para diagnóstico, com nível configurável e destino em stderr.

Nenhum registro de log SHALL conter valor de dado do usuário. Atributo de log SHALL carregar apenas nome de objeto de catálogo, contagem, proporção ou duração.

#### Scenario: Nível configurável

- **WHEN** o nível de log é elevado para depuração
- **THEN** registros adicionais aparecem em stderr, e stdout permanece inalterado

#### Scenario: Log não vaza dado

- **WHEN** qualquer registro de log é emitido em qualquer nível
- **THEN** nenhum valor originado de tabela do usuário aparece nele
