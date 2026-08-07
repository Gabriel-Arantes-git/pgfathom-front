# naming-detection Specification

## Purpose
TBD - created by archiving change naming-detection. Update Purpose after archive.
## Requirements
### Requirement: Afixo de referência derivado das chaves declaradas

O sistema SHALL derivar o afixo de referência comparando, em cada chave estrangeira declarada de coluna única, o nome da coluna filha com as formas do nome da tabela alvo.

Quando o nome da coluna contém uma forma do nome da tabela, o resíduo antes dela é um prefixo candidato e o resíduo depois dela é um sufixo candidato.

Um banco com centenas de chaves declaradas já está dizendo qual é a sua convenção. Ler isso é mais confiável do que qualquer lista que se possa escrever de antemão.

#### Scenario: Sufixo derivado

- **WHEN** existem chaves declaradas como `lote_idkey` apontando para `lote` e `bairro_idkey` para `bairro`
- **THEN** `_idkey` é detectado como sufixo de referência

#### Scenario: Prefixo derivado

- **WHEN** existem chaves declaradas como `idkey_lote` apontando para `lote`
- **THEN** `idkey_` é detectado como prefixo de referência

#### Scenario: Coincidência isolada não vira convenção

- **WHEN** apenas uma chave entre centenas exibe um afixo
- **THEN** esse afixo não é detectado, porque uma ocorrência é ruído e não convenção

#### Scenario: Schema sem chaves declaradas

- **WHEN** o schema não tem nenhuma chave estrangeira declarada
- **THEN** nenhum afixo de referência é detectado, e o perfil base segue valendo sozinho

### Requirement: Prefixo de tabela derivado por frequência

O sistema SHALL detectar prefixos de nome de tabela por frequência, considerando apenas prefixos terminados em separador.

Exigir o separador é o que distingue convenção de coincidência semântica: `auth_` e `tpl_` são convenção, enquanto dezenas de tabelas começando com `ato` apenas compartilham um assunto.

#### Scenario: Prefixo de aplicação

- **WHEN** um schema tem `auth_user`, `auth_group` e `auth_permission` entre poucas tabelas
- **THEN** `auth_` é detectado como prefixo de tabela

#### Scenario: Assunto comum não é prefixo

- **WHEN** um schema tem `atotipo`, `atoanexo` e `atotramite`, sem separador após `ato`
- **THEN** `ato` não é detectado como prefixo

#### Scenario: Prefixo raro é ignorado

- **WHEN** um prefixo aparece em proporção desprezível das tabelas
- **THEN** ele não é detectado

### Requirement: Detecção é aditiva, nunca substitutiva

O que for detectado SHALL ser acrescentado ao perfil base, e MUST NOT remover ou sobrepor nenhuma regra dele.

A normalização já devolve conjunto de formas, então um afixo detectado errado custa uma forma candidata a mais — ruído que a pontuação penaliza e a validação derruba. Já remover uma regra do perfil custaria um candidato que nunca nasce, que é invisível.

#### Scenario: Regras do perfil sobrevivem

- **WHEN** a detecção acrescenta um sufixo ao perfil `pt-br`
- **THEN** todos os sufixos originais do `pt-br` continuam ativos

#### Scenario: Ordem de sombreamento é preservada

- **WHEN** um afixo detectado sombrearia um do perfil base
- **THEN** o perfil resultante ainda passa na validação de ordem

### Requirement: O detectado é reportado

Toda execução que use detecção SHALL informar o que foi detectado e a evidência que sustenta cada item.

Um perfil que muda sozinho e não avisa é pior do que um perfil errado: o usuário não tem como reproduzir nem contestar o resultado.

#### Scenario: Detecção aparece na saída

- **WHEN** a detecção acrescenta afixos
- **THEN** a saída lista cada um com quantas ocorrências o sustentam

#### Scenario: Nada detectado é dito

- **WHEN** a detecção não encontra convenção alguma
- **THEN** a saída informa que o perfil base foi usado sozinho

### Requirement: Detecção pode ser desligada

O sistema SHALL oferecer forma de desligar a detecção e usar apenas o perfil declarado.

#### Scenario: Desligada

- **WHEN** a detecção é desligada
- **THEN** apenas as regras do perfil declarado são aplicadas, e a saída diz isso

### Requirement: Nenhuma consulta nova ao banco

A detecção MUST NOT emitir consulta alguma. Ela opera sobre o modelo que a leitura de catálogo já produziu.

#### Scenario: Sem custo de I/O

- **WHEN** a detecção roda
- **THEN** nenhuma consulta adicional é enviada ao servidor

