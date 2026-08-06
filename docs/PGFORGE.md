# PGForge

Ferramenta de linha de comando para arqueologia de schema PostgreSQL. Descobre e valida a estrutura real de bancos legados, incluindo relacionamentos que existem nos dados mas nunca foram declarados no schema.

---

## Como usar este documento

Este arquivo é a fonte de verdade do projeto e serve como contexto permanente para desenvolvimento assistido. Antes de implementar qualquer coisa, consulte a seção correspondente aqui. Se uma decisão de implementação contradiz este documento, o documento vence, ou o documento precisa ser atualizado primeiro.

Regras para quem estiver escrevendo código neste repositório:

Nunca escreva comando que altere o banco analisado. A ferramenta é estritamente read-only. Ela gera arquivos `.sql` para o usuário revisar e executar por conta própria.

Nunca imprima, logue ou persista valores de dados das tabelas do usuário. A ferramenta lê dados para comparar chaves, mas o que sai dela são contagens, proporções e nomes de objetos. Isso não é preferência estética, é requisito. O caso de uso alvo inclui bancos de gestão pública com CPF, dados de saúde e cadastro de contribuintes.

Nunca afirme uma relação como certa sem evidência nos dados. Toda inferência sai com nível de confiança e com a métrica que sustenta esse nível.

O README público e todo identificador em código são em inglês. Este documento e a documentação interna podem ficar em português.

---

## O problema

Bancos de dados PostgreSQL antigos e grandes carregam mais informação do que declaram. Ao longo de anos de manutenção, times diferentes e migrações mal executadas, o schema vai perdendo aquilo que dava sentido a ele.

O sintoma mais comum é a chave estrangeira ausente. A coluna `pedido.cliente_id` aponta pra `cliente.id` em toda linha da tabela, mas não existe constraint declarada. Isso acontece por vários motivos legítimos e ilegítimos: alguém desabilitou constraints pra acelerar uma carga e nunca religou, o ORM da época não criava, a migração de um banco antigo trouxe só as tabelas, ou simplesmente ninguém achou necessário.

O custo disso aparece depois. Quem chega no projeto não consegue entender o modelo, porque `\d tabela` não mostra relacionamento nenhum. Ferramenta de ERD gera um diagrama de caixas soltas. Gerador de código produz structs sem associação. E o pior: sem a constraint, o banco nunca impediu que dados órfãos entrassem, então provavelmente já existem, silenciosamente, há anos.

Junto disso vem uma segunda camada de conhecimento perdido. Colunas que funcionam como soft delete mas só em parte das tabelas. Colunas de tenant presentes em algumas entidades e ausentes em outras que deveriam tê-las. Campos que são enum na cabeça de todo mundo e `varchar` livre no banco, com três variações de digitação do mesmo valor. Tabelas que ninguém lê há anos e que continuam sendo mantidas, versionadas e migradas.

Nada disso está documentado em lugar nenhum. Está no banco, e dá pra extrair.

---

## O que o PGForge faz

Conecta a uma instância PostgreSQL, lê o catálogo do sistema, cruza com as estatísticas de uso e amostra os dados de forma controlada. A partir disso monta um modelo semântico do banco: as entidades, os relacionamentos reais (declarados e inferidos), e um conjunto de achados sobre a saúde estrutural do schema.

A saída principal não é um diagrama. É um relatório acionável mais artefatos prontos pra uso: DDL sugerida, queries de diagnóstico, e um modelo em JSON que outras ferramentas podem consumir.

O fluxo conceitual:

1. Inspeção do catálogo, para saber o que está declarado.
2. Geração de candidatos, para levantar hipóteses sobre o que não está.
3. Pontuação por metadados, para descartar hipóteses fracas antes de tocar em dados.
4. Validação contra os dados, para confirmar ou derrubar as hipóteses restantes.
5. Relatório e artefatos.

---

## O diferencial

Existe muita ferramenta boa no ecossistema PostgreSQL, e o PGForge não deve competir com nenhuma delas. Vale ser explícito sobre o que já está resolvido:

Lint de migração está coberto pelo Squawk. Detecção de drift entre ambientes está coberta pelo Atlas e pelo migra. Saúde de índice está coberta pelo pganalyze e por dezenas de coleções de query. Diagrama ER está coberto pelo SchemaSpy, Azimutt e dbdocs. Geração de código a partir de schema está coberta pelo sqlc, jOOQ, Ent e Prisma.

Inferência de chave estrangeira não declarada, por outro lado, só existe hoje em ferramenta comercial de interface gráfica, tipo Hackolade e Oracle Data Modeler, e sempre baseada exclusivamente em metadados. No open source existem scripts pequenos que fazem casamento de nome de coluna com nome de tabela mais sufixo, e param aí.

O diferencial do PGForge é validar a inferência contra os dados reais. Casar nomes é trivial e produz muito falso positivo. O que separa hipótese de fato é rodar um anti-join e responder três perguntas: essa relação bate em quantos por cento das linhas, quantos registros órfãos existem, e a distribuição de cardinalidade indica um-para-muitos ou um-para-um.

Isso divide os candidatos em três grupos com tratamentos diferentes:

**Relação confirmada.** Contenção total, nenhum órfão. É uma FK que alguém esqueceu de declarar. Gere a DDL e siga.

**Relação real com integridade quebrada.** Contenção alta mas não total. Existem órfãos. Esse é o achado mais valioso da ferramenta, porque é um bug de dados que está em produção há anos e ninguém sabe. Aqui não basta gerar a DDL, precisa gerar também a query que lista os órfãos, porque eles têm que ser resolvidos antes.

**Coincidência de nome.** Contenção baixa. Descartar, e registrar que foi descartado, pra que o usuário não fique se perguntando por que a ferramenta ignorou uma coluna óbvia.

---

## Escopo

### v0.1 (MVP)

Um comando: `pgforge discover`. Inferência e validação de chaves estrangeiras de coluna única. Saída em tabela no terminal, JSON e SQL.

Esse recorte é intencionalmente estreito. Se a inferência funcionar bem num banco real de trezentas tabelas, o projeto já tem valor publicável. Todo o resto é extensão.

### Fora do escopo do MVP

Chaves compostas. Detecção de padrões semânticos além de FK (soft delete, tenant, enum). Geração de código. Diagrama. Suporte a outros bancos além de PostgreSQL. Modo de escrita no banco, em qualquer fase.

### Fases seguintes

**Fase 2 — Achados estruturais.** Comando `pgforge audit`. Tabelas sem leitura nem escrita desde o último reset de estatística. Colunas integralmente nulas. FK declarada sem índice na coluna filha. Colunas com nome de padrão temporal ou de exclusão lógica usadas de forma inconsistente entre tabelas. `varchar` com baixa cardinalidade que deveria ser enum ou tabela de domínio, com as variações de digitação destacadas.

**Fase 3 — Padrões transversais.** Detecção de coluna de tenant e identificação de tabelas que deveriam tê-la e não têm. Detecção de relacionamento polimórfico, o par `entidade_id` mais `entidade_tipo`, que a inferência simples nunca vai pegar corretamente.

**Fase 4 — Consumidores do modelo.** Exportação para formatos de outras ferramentas (DBML, Mermaid, PlantUML) e geração de código a partir do modelo enriquecido. Essa é a parte que o gerador original do PGForge virou. Ela vem por último de propósito: gerar código a partir de um modelo que conhece os relacionamentos reais é melhor do que qualquer gerador de mercado consegue fazer em banco legado, mas isso só é verdade se as fases anteriores forem confiáveis.

---

## Arquitetura

Camadas com dependência em sentido único. Cada uma testável isolada.

```
cmd/pgforge          entrada CLI, parsing de flag, orquestração
  |
internal/db          conexão, pool, políticas de segurança e timeout
  |
internal/catalog     leitura do pg_catalog e information_schema
  |
internal/model       modelo interno, tipos puros, sem I/O
  |
internal/infer       geração e pontuação de candidatos (só metadados)
  |
internal/validate    validação contra dados, amostragem, anti-join
  |
internal/report      renderização em terminal, JSON, SQL
```

`internal/model` não importa nada das outras camadas. `internal/infer` opera exclusivamente sobre o model e é determinístico, sem acesso a banco, o que torna o teste dele trivial. `internal/validate` é a única camada que lê dados de tabela do usuário.

---

## Modelo interno

Esboço das estruturas centrais. Ajustar durante a implementação, mas manter a separação entre o que foi lido e o que foi inferido.

```go
type Schema struct {
    Name   string
    Tables []Table
}

type Table struct {
    Schema      string
    Name        string
    Columns     []Column
    PrimaryKey  []string        // nomes de coluna, em ordem
    Uniques     [][]string
    ForeignKeys []ForeignKey    // apenas as DECLARADAS
    Indexes     []Index
    Stats       TableStats
    Comment     string
}

type Column struct {
    Name      string
    Type      string            // tipo formatado, ex "bigint", "character varying(60)"
    BaseType  string            // tipo normalizado para comparação, ex "int8"
    Nullable  bool
    Default   string
    Position  int
    Comment   string
}

type TableStats struct {
    EstimatedRows int64         // reltuples
    SeqScans      int64
    IdxScans      int64
    Inserts       int64
    Updates       int64
    Deletes       int64
    TotalBytes    int64
}

type Candidate struct {
    Child        ColumnRef
    Parent       ColumnRef
    Signals      []Signal       // por que virou candidato
    MetaScore    float64        // 0..1, antes de olhar dados
    Validation   *Validation    // nil se não foi validado
    Verdict      Verdict
}

type Validation struct {
    Method          string      // "full" ou "sampled"
    SampledRows     int64
    ChildNotNull    int64
    ChildDistinct   int64
    Matched         int64
    Orphans         int64
    Containment     float64     // Matched / ChildNotNull
    MaxChildPerParent int64     // pra distinguir 1:1 de 1:N
    Duration        time.Duration
}

type Verdict string

const (
    VerdictConfirmed Verdict = "confirmed"  // FK esquecida, íntegra
    VerdictBroken    Verdict = "broken"     // FK real com órfãos
    VerdictWeak      Verdict = "weak"       // evidência insuficiente
    VerdictRejected  Verdict = "rejected"   // coincidência de nome
)
```

---

## Algoritmo de inferência

### Etapa 1: geração de candidatos

Para cada coluna que não seja chave primária da própria tabela e que não participe de FK já declarada, tentar extrair um nome de entidade alvo.

Normalização do nome da coluna. Remover sufixos de referência, na ordem: `_id`, `_codigo`, `_cod`, `_key`, `_ref`, `_fk`, `id`, `cod`. Remover prefixos: `id_`, `cod_`, `fk_`. O resultado é o candidato a nome de entidade.

Normalização do nome da tabela. Remover prefixos comuns de convenção antiga: `tb_`, `tbl_`, `sys_`, `cad_`, `mov_`. Aplicar despluralização, e aqui a heurística precisa cobrir português, porque o público alvo tem banco em português. Regras: `ões` vira `ão`, `ães` vira `ão`, `is` vira `il`, `res` vira `r`, `ses` vira `s`, e por último `s` cai. Guardar tanto a forma original quanto a normalizada e casar contra ambas.

Um candidato nasce quando o nome de entidade extraído da coluna casa com o nome normalizado de alguma tabela, e essa tabela tem chave primária de coluna única, e o tipo base da coluna filha é compatível com o tipo da PK.

Compatibilidade de tipo: idêntico é o caso ideal. Inteiro menor apontando pra inteiro maior é aceitável. `text` e `varchar` de qualquer tamanho são intercambiáveis entre si. `uuid` só casa com `uuid`. Numérico não casa com textual, nunca.

### Etapa 2: pontuação por metadados

Sinais que somam confiança, antes de tocar em dados:

Casamento exato do nome de entidade com o nome da tabela, sem precisar de normalização agressiva, é o sinal mais forte. Tipo base idêntico à PK alvo, forte. Apenas uma tabela candidata pra aquele nome, forte, porque ambiguidade é o principal gerador de ruído. A coluna já possuir índice, moderado, porque indica que alguém a usa em join. Comentário da coluna ou da tabela mencionando a entidade alvo, moderado. Coluna ser `NOT NULL`, fraco mas positivo.

Sinais que subtraem: a tabela alvo é muito pequena e o nome é genérico, tipo `status` ou `tipo`, porque isso costuma ser tabela de domínio onde o casamento é real mas a relação é menos interessante. Múltiplas tabelas candidatas para o mesmo nome. Tipo compatível mas não idêntico.

Candidatos abaixo de um limiar configurável são descartados sem chegar na validação. Isso é o que evita disparar milhares de anti-joins.

### Etapa 3: validação contra dados

Para cada candidato sobrevivente, uma query de agregação. Nunca trazer linhas, só contagens.

```sql
SELECT
    count(*) FILTER (WHERE c.<col> IS NOT NULL)                      AS not_null,
    count(DISTINCT c.<col>)                                          AS distinct_vals,
    count(*) FILTER (WHERE c.<col> IS NOT NULL AND p.<pk> IS NULL)   AS orphans
FROM <child_sample> c
LEFT JOIN <parent> p ON p.<pk> = c.<col>;
```

Em modo amostrado, `<child_sample>` é a tabela filha com `TABLESAMPLE SYSTEM` calibrado para atingir aproximadamente o número de linhas alvo, com fallback para `TABLESAMPLE BERNOULLI` em tabela pequena, onde `SYSTEM` amostra por página e distorce demais. Em modo completo, é a tabela direta.

Antes de cada validação, aplicar `statement_timeout`. Candidato que estourar o timeout é marcado como não validado, com o motivo registrado, e a execução continua. A ferramenta nunca deve travar num banco grande, e nunca deve deixar uma query pendurada.

### Etapa 4: veredito

Contenção total e mais de uma linha distinta na coluna filha resulta em confirmada.

Contenção acima do limiar de quebra, por padrão noventa por cento, mas abaixo de total, resulta em quebrada. Registrar a contagem de órfãos, que é o dado que interessa.

Contenção abaixo do limiar de rejeição, por padrão cinquenta por cento, resulta em rejeitada.

Coluna com um único valor distinto, ou com proporção de nulos muito alta, resulta em fraca independente da contenção, porque a evidência estatística não sustenta conclusão.

Em modo amostrado, nenhum candidato pode ser marcado como confirmada com certeza plena. A saída deve deixar explícito que a validação foi por amostra e que a confirmação definitiva exige `--full`. Isso é questão de honestidade da ferramenta.

---

## Casos difíceis a documentar no código

Não precisam ser resolvidos no MVP, mas precisam ser detectados e reportados como "não analisado", em vez de gerar resultado errado em silêncio.

Relacionamento polimórfico, onde `documento_id` só faz sentido junto com `documento_tipo`. A validação vai encontrar contenção baixa e rejeitar, o que é o comportamento correto, mas a ferramenta deveria idealmente reconhecer o padrão pelo nome das colunas vizinhas.

Chave estrangeira composta, fora do escopo, mas se a tabela alvo tem PK composta ela deve ser pulada com nota, não ignorada silenciosamente.

Tabelas particionadas, onde estatística e contagem se comportam diferente. Ler da tabela pai e não iterar partições.

Herança de tabela, raro mas existe em base antiga.

Múltiplos schemas, onde a mesma tabela pode existir em vários. O casamento de nome deve considerar o schema, e o comportamento entre schemas precisa ser configurável.

Colunas de referência que apontam para tabelas que não existem mais. Contenção zero, rejeição correta, mas vale reportar separado porque é um achado em si.

---

## Interface de linha de comando

```
pgforge discover [flags]

  --dsn string           string de conexão PostgreSQL (ou variável PGFORGE_DSN)
  --schema strings       schemas a analisar (padrão: public)
  --exclude strings      padrões de tabela a ignorar
  --full                 validar contra todas as linhas, sem amostragem
  --sample int           linhas alvo por amostra (padrão: 100000)
  --min-score float      limiar de metadado para validar (padrão: 0.5)
  --timeout duration     statement_timeout por query de validação (padrão: 30s)
  --concurrency int      validações simultâneas (padrão: 4)
  --format string        table | json | sql (padrão: table)
  --out string           diretório de saída para artefatos
  --include-rejected     mostrar também os candidatos descartados
```

A conexão deve falhar com mensagem clara se o usuário fornecido tiver permissão de escrita, ou pelo menos avisar. O recomendado na documentação é criar um usuário exclusivo somente-leitura.

---

## Saídas

### Terminal

Tabela agrupada por veredito, quebradas primeiro porque são o achado mais urgente, depois confirmadas, depois fracas. Colunas: relação, contenção, órfãos, linhas analisadas, método. Resumo final com contagens por veredito e tempo total.

### JSON

O modelo completo, incluindo o que foi lido do catálogo e o que foi inferido, com todos os sinais e métricas de validação preservados. Esse arquivo é o contrato de integração com as fases futuras e com ferramentas de terceiros. Versionar o formato desde o primeiro release.

### SQL

Um arquivo por categoria.

Para confirmadas, DDL com `NOT VALID` e o `VALIDATE CONSTRAINT` separado, comentado, porque `NOT VALID` não trava a tabela e a validação posterior pega um lock mais leve. Incluir também o `CREATE INDEX CONCURRENTLY` na coluna filha quando ela não tiver índice, porque FK sem índice do lado filho é armadilha clássica em delete.

Para quebradas, primeiro a query que lista os órfãos, depois a DDL comentada, com aviso de que ela só vai passar depois da limpeza.

Nenhum arquivo gerado deve ser executável sem revisão humana. Cabeçalho explícito em cada um.

---

## Decisões técnicas

Go. Binário único, sem runtime, distribuível por `go install` e por release no GitHub. Casa com o resto do ecossistema de ferramenta de banco e com o objetivo de aprendizado do projeto.

`pgx/v5` direto, sem ORM e sem camada de abstração. A ferramenta lê catálogo, não faz CRUD.

`cobra` para CLI. É o padrão de fato e reduz atrito de contribuição.

Nenhuma dependência que exija cgo. Cross-compile precisa ser trivial.

Concorrência limitada e configurável nas validações. Rodar quarenta anti-joins simultâneos num banco de produção é a forma mais rápida de a ferramenta ser banida da empresa.

`default_transaction_read_only` ativado na sessão. Cinto e suspensório junto com o usuário read-only.

---

## Testes

Unitário puro para `internal/infer`, que é determinístico. É onde mora a maior parte da lógica sutil, especialmente a despluralização em português, e é o que mais vai quebrar com mudança. Cobertura alta aqui é obrigatória.

Integração com `testcontainers-go` subindo PostgreSQL real. Fixtures SQL em `testdata/`, cada uma representando um cenário: schema limpo com FK declarada, schema sem nenhuma FK, schema com órfãos plantados de propósito, schema com colisão de nome, schema em português com plural irregular.

Golden files para as saídas de terminal e SQL, porque formatação regride fácil.

Um cenário de teste deve ser explicitamente uma armadilha: coluna chamada `status_id` numa tabela onde existe uma tabela `status` mas a coluna guarda outra coisa. A ferramenta tem que rejeitar.

---

## Critério de aceite do MVP

O comando roda contra um banco real de porte relevante, na casa das centenas de tabelas, sem travar, sem exceder o timeout global e sem impacto perceptível na carga do servidor.

Toda relação declarada existente no banco é redescoberta pela inferência quando as constraints são removidas de uma cópia. Esse é o melhor teste de qualidade disponível: pegue um schema que tem FKs, apague todas, rode a ferramenta e meça quantas ela recupera. A taxa de recuperação é a métrica principal do projeto e deve estar no README.

Nenhum falso positivo confirmado. Falso negativo é aceitável, falso positivo confirmado destrói a confiança na ferramenta.

A saída SQL é executável sem edição manual num banco de teste.

---

## Glossário

**Candidato.** Par de colunas que pode representar uma relação, levantado por heurística de nome e tipo, antes de qualquer verificação nos dados.

**Contenção.** Proporção dos valores não nulos da coluna filha que existem na chave da tabela pai. É a métrica central da validação.

**Órfão.** Linha da tabela filha cujo valor de referência não existe na tabela pai. Em banco com FK declarada é impossível. Em banco sem FK é o que a ferramenta procura.

**Relação declarada.** Chave estrangeira que existe como constraint no catálogo.

**Relação inferida.** Relação que existe nos dados mas não no catálogo.

**Modo amostrado.** Validação sobre subconjunto das linhas, rápido, indicativo, incapaz de provar ausência de órfão.

**Modo completo.** Validação sobre todas as linhas, lento, conclusivo.
