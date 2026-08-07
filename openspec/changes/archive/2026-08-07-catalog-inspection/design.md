## Context

A fase 1 entregou o modelo, os perfis e o esqueleto da CLI, tudo determinístico e sem I/O. Esta fase abre a primeira conexão.

O contexto que importa não é técnico, é de confiança. A ferramenta será apontada para o banco de produção de outra pessoa, por alguém que provavelmente não é o DBA, com autorização que pode ser revogada ao primeiro incidente. Uma execução que trave um servidor não é um bug a corrigir na próxima versão — é o fim da adoção naquela empresa e provavelmente na próxima, porque essa história circula.

Restrições herdadas: read-only absoluto, nenhum valor de dado do usuário em qualquer saída, silêncio nunca reportado como ausência de problema, sem cgo.

Decisão de escopo tomada nesta fase: **PostgreSQL 13 como piso**. Cobre todas as versões com suporte da comunidade e evita caminho condicional na leitura de catálogo.

## Goals / Non-Goals

**Goals**

Concentrar as políticas de segurança num único ponto que nenhuma conexão possa contornar. Ler o catálogo por inteiro para dentro do modelo, incluindo o que precisa ser pulado. Entregar o `audit` com dois achados determinísticos. Montar a infraestrutura de teste de integração que as fases 3 a 8 vão reusar.

**Non-Goals**

Nenhuma leitura de dado de tabela do usuário — a camada de catálogo não emite consulta que toque em relação do usuário, e essa fronteira é verificável. Nenhuma inferência, nenhuma pontuação, nenhum candidato. Nenhuma geração de SQL. Chave estrangeira composta segue fora de escopo, mas passa a ser registrada em vez de ignorada. O relatório completo, com agrupamento por veredito e golden files, é da fase 7; aqui o renderizador é o mínimo que o `audit` exige.

## Decisions

### Políticas de sessão no AfterConnect, não por consulta

`pgxpool.Config.AfterConnect` roda uma vez por conexão nova, incluindo as que o pool abre sob concorrência. Aplicar `default_transaction_read_only`, `application_name` e os três timeouts ali significa que nenhuma conexão pode existir sem eles.

A alternativa — aplicar por consulta, ou num helper que todo chamador precisa lembrar de usar — depende de disciplina humana para uma garantia que o projeto anuncia no README. Uma conexão aberta em algum caminho esquecido é exatamente o tipo de furo que só aparece em produção alheia.

O teste correspondente prova a garantia pelo comportamento observável: uma tentativa de escrita numa conexão do pool falha.

### Timeouts como cinto e suspensório separados

`statement_timeout` limita a consulta. `lock_timeout` limita a espera por lock, que é um estado diferente e igualmente capaz de segurar recurso. `idle_in_transaction_session_timeout` cobre o caso em que a ferramenta abre transação e algo trava do lado do cliente — o pior dos três, porque uma transação ociosa aberta segura o `VACUUM` do banco inteiro.

Os três são configuráveis com padrão conservador. Estourar timeout marca o item como não avaliado e a execução prossegue; a ferramenta nunca deve travar num banco grande.

### Precedência de credencial invertendo o óbvio

O caminho óbvio seria a flag ter precedência sobre a variável, como em quase toda CLI. Aqui `PGFATHOM_DSN` vence `--dsn`, e a divergência é avisada.

A razão é que `--dsn` aparece em `ps` e no histórico do shell, e a documentação recomenda a variável. Se a flag vencesse silenciosamente, a recomendação seria contornável por engano em ambiente compartilhado. Avisar em vez de falhar mantém o comportamento previsível para quem passou as duas por acidente.

### Índice do lado filho exige coluna em posição inicial

Um índice composto que contenha a coluna filha em posição não inicial não serve para a busca que o `DELETE` no pai dispara. Contá-lo como cobertura produziria um falso negativo — a ferramenta diria que está tudo bem numa tabela onde o problema existe.

Reportar a FK como não indexada nesse caso pode gerar recomendação redundante em algum arranjo incomum. Falso positivo em recomendação de índice custa uma conversa; falso negativo custa o incidente que a ferramenta existia para evitar.

### Tabela particionada lida a partir do pai

Iterar partições produziria uma entrada por partição no modelo, contagens somadas em dobro e estatística que não corresponde a nada que o usuário reconheça. A leitura vem do pai, as partições não entram no modelo, e a tabela é marcada como particionada para que as fases seguintes saibam que a validação precisa de tratamento próprio.

### Teste de integração com testcontainers, atrás de tag

`testcontainers-go` sobe um PostgreSQL real. É a única forma honesta de testar leitura de catálogo: um duplo de teste testaria o duplo.

A tag `integration`, criada vazia na fase 1, passa a valer aqui. `go test ./...` continua sem Docker e sem rede, porque contribuir um perfil de nomenclatura não pode exigir runtime de contêiner.

As fixtures SQL em `testdata/` são versionadas e cada uma representa um cenário nomeado. Elas serão reusadas até a fase 8, então vale montá-las com cuidado agora: schema limpo com FK declarada, schema sem nenhuma FK, schema com constraint `NOT VALID` e órfãos preexistentes, schema com colisão de nome, schema em português com plural irregular.

### Renderizador mínimo agora, completo na fase 7

O `audit` precisa de saída, mas escrever golden files enquanto o conteúdo ainda muda cria manutenção morta. Aqui o renderizador cobre o que o `audit` emite, com a disciplina de fluxo e o bloco de cobertura já corretos. A fase 7 reescreve com agrupamento por veredito e fixa os golden files quando o conteúdo estabilizar.

### Ausência de achado é afirmada, não implícita

Uma saída vazia é ambígua: pode significar que está tudo bem ou que nada foi olhado. O comando afirma explicitamente o que analisou e que está limpo, e o bloco de cobertura acompanha toda execução.

Isso é a regra de "silêncio nunca é ausência de problema" virando comportamento observável em vez de intenção documentada.

## Risks / Trade-offs

**A verificação de privilégio de escrita custa uma consulta por tabela do escopo** → resolver em consulta única agregada sobre `information_schema.table_privileges` em vez de iterar. Em schema de trezentas tabelas a diferença é perceptível.

**`pg_stat_activity` e a checagem de privilégio dependem de visibilidade que um papel restrito pode não ter** → falha na verificação vira aviso, nunca erro fatal. A ferramenta degrada para menos informação, não para não funcionar.

**Fixtures de teste envelhecem junto com o schema alvo** → nomeá-las por cenário e não por conteúdo, e mantê-las mínimas. Fixture grande é fixture que ninguém entende quando quebra.

**Piso em 13 exclui bases em 11 e 12 que ainda existem no público-alvo** → assumido conscientemente. São versões sem patch de segurança, e suportá-las custaria caminho condicional em toda a leitura de catálogo. Registrar a versão exigida em mensagem de erro clara é o mitigador.

**O `audit` pode dar impressão de análise completa quando é só a parte determinística** → o texto da saída deixa explícito que o `audit` não faz inferência e que `discover` é outro comando, sem sugerir que um substitui o outro.

## Migration Plan

Não aplicável. Não há release nem usuário.

## Open Questions

Nenhuma. As duas questões abertas da fase 1 foram decididas: piso de PostgreSQL 13 e licença Apache-2.0, ambas incorporadas nesta change.
