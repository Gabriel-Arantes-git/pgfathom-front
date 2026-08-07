# database-session Specification

## Purpose
TBD - created by archiving change catalog-inspection. Update Purpose after archive.
## Requirements
### Requirement: Sessão somente-leitura por construção

Toda conexão do pool SHALL ter `default_transaction_read_only` ativado antes de ser entregue a qualquer chamador. A configuração SHALL ser aplicada no hook de conexão do pool, não em cada consulta, de modo que uma conexão nova não possa existir sem ela.

O código MUST NOT conter nenhum caminho que emita `INSERT`, `UPDATE`, `DELETE`, `TRUNCATE`, `ALTER`, `CREATE`, `DROP` ou `GRANT` contra o banco analisado, sob qualquer flag.

#### Scenario: Escrita falha na sessão

- **WHEN** uma tentativa de escrita é executada numa conexão do pool
- **THEN** o servidor a rejeita, provando que a política está ativa

#### Scenario: Toda conexão nova herda a política

- **WHEN** o pool abre uma conexão adicional sob concorrência
- **THEN** essa conexão também tem a política aplicada, sem depender de o chamador lembrar

### Requirement: A ferramenta se identifica no servidor

Toda conexão SHALL definir `application_name` como `pgfathom`, para que apareça identificada em `pg_stat_activity`.

O DBA que autorizou a execução precisa conseguir ver o que está rodando e de onde vem, sem adivinhar a partir do texto da consulta.

#### Scenario: Identificação visível

- **WHEN** uma sessão está ativa
- **THEN** `pg_stat_activity` a mostra com `application_name` igual a `pgfathom`

### Requirement: Nenhuma consulta pode ficar pendurada

Toda conexão SHALL ter `statement_timeout`, `lock_timeout` e `idle_in_transaction_session_timeout` definidos, todos configuráveis e com padrão conservador.

Uma execução que trave um servidor de produção alheio encerra o projeto, independentemente da qualidade do resto.

#### Scenario: Consulta longa é interrompida pelo servidor

- **WHEN** uma consulta excede o `statement_timeout` configurado
- **THEN** o servidor a cancela e a ferramenta registra o motivo, sem derrubar a execução inteira

#### Scenario: Cancelamento não deixa trabalho no servidor

- **WHEN** o contexto raiz é cancelado durante uma consulta em andamento
- **THEN** a consulta é cancelada no servidor e o pool é encerrado sem conexão pendente

### Requirement: Precedência de credencial documentada e sem vazamento

A resolução de credencial SHALL seguir, nesta ordem: a variável `PGFATHOM_DSN`, depois a flag `--dsn`, e por último as variáveis padrão do libpq (`PGHOST`, `PGPORT`, `PGUSER`, `PGDATABASE`, `PGPASSWORD`, `PGPASSFILE`).

`PGFATHOM_DSN` vence a flag porque `--dsn` aparece em `ps` e no histórico do shell, e a documentação recomenda a variável — se a flag vencesse, a recomendação seria contornável por engano em máquina compartilhada.

As variáveis do libpq vêm por último porque são configuração ambiente, não intenção explícita: `PGHOST` e `PGUSER` costumam estar exportadas no shell de quem trabalha com PostgreSQL, e deixá-las sobrepor uma `--dsn` escrita à mão seria surpreendente.

A ajuda da flag `--dsn` SHALL alertar que o valor aparece em `ps` e no histórico do shell.

Nenhuma mensagem de erro, log ou saída SHALL conter a senha. Um DSN reproduzido em diagnóstico SHALL ter a senha substituída por um marcador.

#### Scenario: Variável tem precedência sobre a flag

- **WHEN** `PGFATHOM_DSN` está definida e `--dsn` também é fornecida
- **THEN** a conexão usa `PGFATHOM_DSN`, e a divergência é avisada em stderr

#### Scenario: Flag tem precedência sobre o ambiente libpq

- **WHEN** `PGHOST` está exportada no shell e `--dsn` é fornecida
- **THEN** a conexão usa a `--dsn`, porque intenção explícita vence configuração ambiente

#### Scenario: Ambiente libpq sozinho funciona

- **WHEN** nenhuma variável `PGFATHOM_DSN` e nenhuma flag são fornecidas, mas as variáveis do libpq estão definidas
- **THEN** a conexão é estabelecida a partir delas

#### Scenario: Senha nunca aparece em diagnóstico

- **WHEN** a conexão falha com um DSN que contém senha
- **THEN** a mensagem de erro identifica host, porta, banco e usuário, e a senha aparece apenas como marcador

### Requirement: Versão do servidor verificada antes de qualquer leitura

A ferramenta SHALL verificar a versão do servidor antes de ler o catálogo, e SHALL recusar operação abaixo de PostgreSQL 13 com mensagem indicando a versão encontrada e a mínima exigida.

Falhar cedo com mensagem clara é melhor do que falhar no meio da leitura com erro de coluna inexistente, que o usuário não tem como interpretar.

#### Scenario: Servidor abaixo do piso

- **WHEN** o servidor reporta versão anterior a 13
- **THEN** a execução falha antes de ler o catálogo, informando a versão encontrada e a exigida

#### Scenario: Versão registrada na saída

- **WHEN** a conexão é estabelecida
- **THEN** a versão do servidor consta do resultado, porque a interpretação de um achado depende dela

### Requirement: Privilégio de escrita é avisado

A ferramenta SHALL verificar se o papel conectado possui `INSERT`, `UPDATE` ou `DELETE` sobre as tabelas do escopo, e SHALL avisar em stderr quando possuir, recomendando um papel exclusivo somente-leitura.

A sessão somente-leitura já impede a escrita. O aviso existe porque quem executa a ferramenta muitas vezes não é quem escolheu o papel, e a recomendação precisa chegar a essa pessoa.

#### Scenario: Papel com privilégio de escrita

- **WHEN** o papel conectado pode escrever em alguma tabela do escopo
- **THEN** um aviso vai para stderr, a execução continua, e o resultado permanece em stdout intacto

#### Scenario: Papel somente-leitura não gera ruído

- **WHEN** o papel conectado não tem privilégio de escrita em nenhuma tabela do escopo
- **THEN** nenhum aviso é emitido

### Requirement: Concorrência limitada e configurável

O número de consultas simultâneas SHALL ser limitado por configuração, com padrão conservador, e o tamanho do pool SHALL acompanhar esse limite.

Disparar dezenas de consultas simultâneas num banco de produção é a forma mais rápida de a ferramenta ser proibida na empresa.

#### Scenario: Limite respeitado

- **WHEN** há mais trabalho pendente do que o limite de concorrência
- **THEN** o excedente aguarda, e o número de consultas simultâneas nunca ultrapassa o limite

