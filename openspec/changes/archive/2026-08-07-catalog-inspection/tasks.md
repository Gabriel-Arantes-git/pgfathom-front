## 1. Licença e preparação

- [x] 1.1 Adicionar `LICENSE` com o texto Apache-2.0
- [x] 1.2 Atualizar a seção de licença do `README.md`, removendo o TBD
- [x] 1.3 Adicionar `github.com/jackc/pgx/v5` ao `go.mod`
- [x] 1.4 Adicionar `testcontainers-go` e o módulo `postgres` como dependências de teste, verificando que não entram no binário

## 2. Sessão e políticas de segurança

- [x] 2.1 Criar `internal/db` com a construção do pool `pgxpool` a partir de configuração
- [x] 2.2 Aplicar no `AfterConnect`: `default_transaction_read_only`, `application_name`, `statement_timeout`, `lock_timeout` e `idle_in_transaction_session_timeout`
- [x] 2.3 Tornar os três timeouts e a concorrência configuráveis, com padrões conservadores, e dimensionar o pool pelo limite de concorrência
- [x] 2.4 Implementar a precedência de credencial: `PGFATHOM_DSN`, variáveis do libpq, `--dsn`, com aviso quando houver divergência
- [x] 2.5 Implementar a redação de senha em toda mensagem de erro e log que reproduza um DSN
- [x] 2.6 Verificar a versão do servidor e recusar abaixo de 13, informando encontrada e exigida
- [x] 2.7 Verificar privilégio de escrita do papel numa consulta única agregada, e avisar em stderr quando houver
- [x] 2.8 Propagar o `context.Context` da raiz até a execução de consulta, garantindo cancelamento no servidor

## 3. Infraestrutura de teste de integração

- [x] 3.1 Implementar em `internal/testutil` o helper que sobe PostgreSQL via testcontainers, atrás da tag `integration`
- [x] 3.2 Criar o helper que carrega uma fixture SQL nomeada no contêiner
- [x] 3.3 Escrever a fixture `clean_schema`: schema com FKs declaradas e validadas, tudo íntegro
- [x] 3.4 Escrever a fixture `no_constraints`: mesmas tabelas, nenhuma FK declarada
- [x] 3.5 Escrever a fixture `not_valid_constraints`: FKs `NOT VALID` com órfãos preexistentes plantados
- [x] 3.6 Escrever a fixture `unindexed_fks`: FKs declaradas sem índice do lado filho, incluindo o caso de índice composto com a coluna em posição não inicial
- [x] 3.7 Escrever a fixture `unsupported_shapes`: PK composta, tabela sem PK, tabela particionada
- [x] 3.8 Escrever a fixture `restricted_privileges`: papel sem `SELECT` em parte das tabelas
- [x] 3.9 Plantar valores reconhecíveis nas fixtures, para a varredura de vazamento
- [x] 3.10 Adicionar ao CI o job de integração, em Linux com Docker, separado do job rápido

## 4. Leitura de catálogo

- [x] 4.1 Criar `internal/catalog` com a consulta de tabelas e colunas, incluindo tipo formatado, tipo base, nulabilidade, padrão e posição
- [x] 4.2 Implementar a normalização de tipo base, com teste cobrindo grafias equivalentes
- [x] 4.3 Ler chaves primárias em ordem de coluna
- [x] 4.4 Ler constraints unique
- [x] 4.5 Ler índices com colunas em ordem, marcando unique e primary
- [x] 4.6 Ler chaves estrangeiras declaradas com `convalidated` de `pg_constraint`
- [x] 4.7 Determinar, por FK, se há índice com a coluna filha em posição inicial
- [x] 4.8 Ler comentários de tabela e de coluna
- [x] 4.9 Ler estatísticas de uso de `pg_stat_user_tables` junto do reset de `pg_stat_database`, marcando como não interpretável quando ausente
- [x] 4.10 Ler tamanho e estimativa de linhas de `pg_class`
- [x] 4.11 Implementar o escopo por schema e a exclusão por padrão, registrando as excluídas em campo próprio
- [x] 4.12 Detectar e registrar as formas não suportadas: PK composta, ausência de PK, particionada, herança
- [x] 4.13 Ler tabela particionada a partir do pai, sem iterar partições
- [x] 4.14 Detectar tabela sem privilégio `SELECT` e registrá-la na cobertura
- [x] 4.15 Montar a `Coverage` e garantir que analisadas mais puladas fecham com o total
- [x] 4.16 Escrever teste que falha se alguma consulta da camada referenciar relação fora de `pg_catalog`, `information_schema` e visões de estatística

## 5. Comando audit

- [x] 5.1 Criar `internal/report` com o renderizador de terminal mínimo, respeitando a disciplina de stdout e stderr
- [x] 5.2 Renderizar os achados agrupados por tipo, com o bloco de cobertura no rodapé
- [x] 5.3 Emitir afirmação explícita quando o escopo foi analisado e nada foi encontrado
- [x] 5.4 Implementar a saída JSON a partir do `model.Result`, com `schema_version`
- [x] 5.5 Criar o comando `audit` com as flags de conexão, escopo, formato e timeouts
- [x] 5.6 Implementar o achado de constraint `NOT VALID`
- [x] 5.7 Implementar o achado de FK sem índice do lado filho
- [x] 5.8 Garantir que achados não alteram o código de saída, e que falha e erro de uso usam os códigos da fase 1

## 6. Verificação

- [x] 6.1 Teste de integração provando que escrita falha numa conexão do pool
- [x] 6.2 Teste provando que `application_name` aparece como `pgfathom` em `pg_stat_activity`
- [x] 6.3 Teste provando que cancelamento do contexto encerra a consulta no servidor sem deixar conexão pendente
- [x] 6.4 Teste de recusa de servidor abaixo do piso de versão
- [x] 6.5 Teste de precedência de credencial e de redação de senha em erro
- [x] 6.6 Teste de cobertura contra a fixture `restricted_privileges`, verificando que as tabelas sem acesso aparecem listadas
- [x] 6.7 Teste de `audit` contra `not_valid_constraints` e `unindexed_fks`, conferindo os achados esperados
- [x] 6.8 Teste que varre terminal, JSON e log procurando os valores plantados nas fixtures
- [x] 6.9 Rodar `golangci-lint run` e zerar os apontamentos
- [x] 6.10 Confirmar que `go test ./...` segue sem Docker e sem rede, e que o binário não ganhou dependência de teste
- [x] 6.11 Revisar a densidade de comentário antes de fechar, mantendo o padrão da fase 1
- [x] 6.12 Rodar `openspec validate catalog-inspection` e corrigir o que apontar
