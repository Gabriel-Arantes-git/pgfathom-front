## Why

Rodando a fase 3 contra bancos reais de gestão pública, a inferência recuperou **0,5%** das chaves estrangeiras declaradas num schema de 784 tabelas. Adicionar um único afixo ao perfil — `_idkey`, que aquele sistema usa e nenhum perfil oficial conhecia — levou o mesmo banco a **79,0%**.

Um segundo sistema repetiu o padrão: 75,4% depois do mesmo ajuste. Um terceiro, em Django, ficou em **0,0% mesmo com `--profile en`**, porque prefixa tabelas com o nome do app: `auth_permission`, `django_content_type`.

A conclusão dos três é a mesma: **a convenção de nomenclatura não é do idioma, é do schema.** Escolher perfil de um menu não resolve, porque o usuário não tem como saber o que escolher antes da primeira execução — e o modo de falha é silencioso. A ferramenta não erra, ela devolve quase nada, e parece que o banco não tinha o que descobrir.

Há também uma razão de método. A fase 8 mede recall no corpus de benchmark; se o corpus rodar com perfil ajustado à mão, o número publicado não representa a primeira execução de ninguém.

## What Changes

- Detecção de prefixo de tabela por frequência sobre os nomes do schema.
- Detecção do afixo de referência a partir das chaves estrangeiras já declaradas, comparando o nome da coluna filha com o da tabela alvo.
- Mesclagem do detectado sobre o perfil base, acrescentando sem substituir.
- Relato do que foi detectado na saída e no JSON.
- Flag para desligar a detecção.

Nenhuma consulta nova ao banco: tudo sai do que a fase 2 já lê. Nenhuma leitura de dado.

## Capabilities

### New Capabilities

- `naming-detection`: como o schema revela a própria convenção, o que sustenta cada detecção, e por que o resultado é sempre aditivo.

### Modified Capabilities

Nenhuma. `naming-profiles` continua com os mesmos requisitos; a detecção compõe com ela em vez de alterá-la.

## Impact

`internal/profile` passa a importar `internal/model`, que é puro. Nenhuma dependência externa nova.

Muda o comportamento padrão do `discover`: sem `--profile`, a convenção deixa de ser só o perfil embarcado. Isso é a mudança pretendida, e é por isso que o relato do detectado é requisito e não enfeite.
