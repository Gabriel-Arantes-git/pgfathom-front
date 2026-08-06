## ADDED Requirements

### Requirement: Perfis carregados de arquivo, nunca de código

As regras de afixo e de plural SHALL residir em arquivo TOML. Nenhuma regra de nomenclatura pode estar codificada como constante em Go. O binário SHALL embarcar os perfis oficiais `pt-br`, `en` e `es` via `embed`, e SHALL aceitar caminho para arquivo do usuário.

Um perfil define: sufixos de coluna a remover, prefixos de coluna a remover, prefixos de tabela a remover, e uma lista ordenada de regras de despluralização.

#### Scenario: Perfil oficial embarcado

- **WHEN** o carregador recebe o nome `pt-br`, `en` ou `es`
- **THEN** retorna o perfil embarcado correspondente sem tocar no sistema de arquivos

#### Scenario: Perfil do usuário

- **WHEN** o carregador recebe um caminho de arquivo existente e válido
- **THEN** retorna o perfil daquele arquivo

#### Scenario: Perfil inválido falha explicitamente

- **WHEN** o arquivo não existe, não é TOML válido, ou omite campo obrigatório
- **THEN** o carregamento retorna erro descrevendo o problema, e o processo não continua com um perfil vazio ou parcial

### Requirement: Extração de nome de entidade a partir do nome de coluna

O sistema SHALL derivar um nome de entidade candidato removendo os afixos de referência declarados no perfil. Sufixos SHALL ser testados antes de prefixos, e ambos na ordem declarada, com a primeira regra que casa vencendo.

Quando nenhum afixo casa, o nome da coluna SHALL ser devolvido inalterado como candidato — uma coluna chamada `cliente` referencia `cliente` tão legitimamente quanto `cliente_id`.

#### Scenario: Sufixo de referência

- **WHEN** a coluna se chama `cliente_id` e o perfil `pt-br` está ativo
- **THEN** o nome de entidade extraído é `cliente`

#### Scenario: Prefixo de referência

- **WHEN** a coluna se chama `id_fornecedor`
- **THEN** o nome de entidade extraído é `fornecedor`

#### Scenario: Coluna sem afixo

- **WHEN** a coluna se chama `municipio` e nenhum afixo do perfil casa
- **THEN** o nome de entidade extraído é `municipio`

#### Scenario: Sufixo mais longo tem precedência

- **WHEN** a coluna se chama `pedido_codigo` e o perfil lista tanto `_codigo` quanto `_cod`
- **THEN** o nome de entidade extraído é `pedido`, não `pedido_go`

### Requirement: Normalização de nome de tabela produz conjunto de formas

A normalização de nome de tabela SHALL retornar um **conjunto ordenado** de formas candidatas, não uma única forma. O conjunto SHALL incluir sempre o nome original, e SHALL incluir cada forma singular plausível produzida pelas regras de despluralização aplicáveis.

Casamento entre nome de entidade e nome de tabela SHALL ter sucesso quando qualquer forma do conjunto casa. A forma que casou SHALL ser reportada, para que a pontuação da fase de inferência possa distinguir casamento exato de casamento obtido por normalização agressiva.

Esta é uma decisão deliberada contra a alternativa de primeira-regra-vence. Regras de plural em português são ambíguas: `logins` produz `logim` sob a regra `ns → m` e `login` sob a regra genérica, e não há informação no nome que resolva qual está certa. Devolver as duas e deixar o casamento decidir elimina uma classe inteira de falso negativo dependente da ordem das regras.

#### Scenario: Conjunto inclui o original

- **WHEN** a tabela se chama `cliente`, que já é singular
- **THEN** o conjunto de formas contém `cliente`

#### Scenario: Plural regular

- **WHEN** a tabela se chama `clientes`
- **THEN** o conjunto de formas contém `cliente`

#### Scenario: Ambiguidade resolvida por conjunto

- **WHEN** a tabela se chama `logins` e o perfil tem tanto `ns → m` quanto a regra genérica de `s`
- **THEN** o conjunto de formas contém tanto `logim` quanto `login`, e o casamento com a entidade `login` tem sucesso

#### Scenario: Prefixo de convenção antiga

- **WHEN** a tabela se chama `tb_clientes` e o perfil lista `tb_` entre os prefixos de tabela
- **THEN** o conjunto de formas contém `cliente`, e contém também `tb_clientes` inalterado

#### Scenario: Origem do casamento é reportada

- **WHEN** a entidade `cliente` casa com a tabela `tb_clientes` apenas após remoção de prefixo e despluralização
- **THEN** o resultado indica que o casamento veio de forma normalizada, não da forma original

### Requirement: Despluralização em português cobre identificadores sem acento

O perfil `pt-br` SHALL cobrir tanto as formas acentuadas quanto as formas sem acento das regras de plural, porque identificadores de banco são tipicamente ASCII: uma tabela chamada `opcoes` é muito mais comum do que `opções`.

As regras SHALL ser declaradas da mais específica para a mais genérica, e a regra genérica de queda de `s` SHALL ser a última.

#### Scenario: Plural em -ões sem acento

- **WHEN** a tabela se chama `opcoes`
- **THEN** o conjunto de formas contém `opcao`

#### Scenario: Plural em -ais

- **WHEN** a tabela se chama `animais`
- **THEN** o conjunto de formas contém `animal`

#### Scenario: Plural em -eis

- **WHEN** a tabela se chama `responsaveis`
- **THEN** o conjunto de formas contém `responsavel`

#### Scenario: Plural em -is

- **WHEN** a tabela se chama `perfis`
- **THEN** o conjunto de formas contém `perfil`

#### Scenario: Plural em -ns

- **WHEN** a tabela se chama `armazens`
- **THEN** o conjunto de formas contém `armazem`

#### Scenario: Plural em -es após consoante

- **WHEN** a tabela se chama `meses`
- **THEN** o conjunto de formas contém `mes`

#### Scenario: Forma acentuada também funciona

- **WHEN** a tabela se chama `órgãos`
- **THEN** o conjunto de formas contém `órgão`

### Requirement: Perfil ativo é visível ao usuário

O perfil em uso SHALL ser reportado em toda execução. Um resultado de inferência sem saber qual convenção de nomenclatura foi aplicada não é interpretável.

O perfil padrão SHALL ser `pt-br`.

#### Scenario: Perfil aparece na saída

- **WHEN** qualquer comando que use inferência é executado
- **THEN** o nome do perfil ativo aparece no cabeçalho da saída
