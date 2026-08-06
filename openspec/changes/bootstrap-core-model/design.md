## Context

Repositório vazio de código. Existe a especificação de produto em `docs/PGFATHOM.md` e o plano de fases em `docs/ROADMAP.md`, e nada mais.

Esta change estabelece decisões que as sete fases seguintes herdam e que ficam caras de reverter: a assinatura das funções de normalização que `internal/infer` vai consumir, os nomes de campo do `internal/model` que viram contrato do JSON público, e a disciplina de saída que todo comando futuro segue.

Restrições que vêm de fora e não estão em negociação: sem cgo, árvore de dependências pequena e justificável, `internal/model` sem I/O, e nenhum valor de dado do usuário em qualquer coisa serializável.

## Goals / Non-Goals

**Goals**

Fixar o contrato de normalização de nome antes que qualquer coisa dependa dele. Fixar os tipos do modelo com a proveniência e as invariantes de segurança embutidas na forma, não na disciplina de quem escreve. Estabelecer a disciplina de saída, cancelamento e código de saída no esqueleto CLI, de modo que nenhuma camada posterior seja escrita sem elas. Ter suíte de teste que roda sem Docker.

**Non-Goals**

Nenhuma conexão com banco. Nenhuma leitura de catálogo. Nenhuma geração de candidato — a normalização existe aqui, o casamento que a consome é da fase 3. Nenhuma renderização de tabela de resultado além do mínimo para `version` e ajuda. Detecção automática de perfil de nomenclatura, que fica para depois do MVP.

## Decisions

### Normalização de tabela devolve conjunto, não string

A alternativa óbvia é aplicar regras de plural em ordem e devolver a primeira que casa. Foi rejeitada porque a ordem vira uma fonte silenciosa de falso negativo.

O caso concreto: uma tabela `logins` produz `logim` sob a regra `ns → m`, correta para `armazens → armazem`, e produz `login` sob a regra genérica de queda de `s`. Não há informação no nome que resolva qual está certa. Com primeira-regra-vence, a ordem escolhida decide qual dos dois casos o projeto quebra, e a escolha é arbitrária.

Devolvendo um conjunto ordenado de formas candidatas — sempre incluindo o nome original — o casamento tenta todas e a ambiguidade deixa de custar recall. A forma que casou é reportada, o que dá à pontuação da fase 3 a informação para distinguir casamento exato de casamento por normalização agressiva, que é exatamente a diferença entre os sinais `SigExactName` e `SigNormalizedName`.

O custo é um conjunto pequeno por tabela em memória e a possibilidade de casamento espúrio via forma agressiva. O segundo é aceitável porque a pontuação penaliza e, principalmente, porque a validação contra dados derruba o que for coincidência.

### Perfil de nomenclatura em TOML embarcado, não em Go

Regras em constante Go tornariam a lógica mais frágil do projeto testável apenas por código, e fechariam a porta de contribuição mais fácil que o projeto tem a oferecer. Em arquivo, adicionar um idioma é um PR que não exige entender o resto.

`embed` para os oficiais evita depender de instalação de arquivo ao lado do binário, o que quebraria a promessa de binário único. `--profile` com caminho cobre o usuário que tem convenção própria de casa.

TOML sobre YAML por ser menos ambíguo e ter menos superfície de erro em arquivo escrito à mão. `pelletier/go-toml/v2` sobre `BurntSushi/toml` pela qualidade das mensagens de erro, que importa quando quem escreve o arquivo é o usuário.

### Perfil `pt-br` cobre formas sem acento

Identificador de banco é tipicamente ASCII. Uma tabela `opcoes` é muito mais comum na base alvo do que `opções`, e um perfil que só conhece a forma acentuada erra o caso majoritário.

As regras cobrem os dois conjuntos. As específicas vêm antes das genéricas, e a queda de `s` é sempre a última. A ordem importa menos do que importaria sob primeira-regra-vence, justamente pela decisão do conjunto acima, mas continua determinando a ordem de preferência reportada.

### Invariantes de segurança na forma do tipo, não na disciplina

A regra de não vazar dado do usuário é a mais fácil de quebrar por acidente e a mais cara quando quebra. Confiar em revisão humana é insuficiente.

Estatística de coluna vinda de `pg_stats` carrega `most_common_vals` e `histogram_bounds`, que são dados do usuário. Esses campos ficam **não exportados**, de modo que serialização em JSON não os alcança por construção, e não por lembrança de quem escreveu a struct. Campo de diagnóstico de sinal aceita apenas nome de objeto.

Isso é reforçado por teste que serializa todas as estruturas do modelo e varre o resultado. O teste existe desde esta change, quando ainda não há dado nenhum para vazar, para que a fase 4 encontre a rede já armada.

### `context.Context` desde a raiz, antes de existir I/O

Não há nada para cancelar nesta change. A propagação é estabelecida agora precisamente por isso: retrofit de contexto depois que três camadas já existem é refatoração larga e propensa a deixar caminho sem cobertura. Toda função de camada nasce recebendo contexto.

### `text/tabwriter` em vez de biblioteca de renderização

A fase 7 tem golden files para a saída de terminal. Renderizador estável significa golden file que não regride por bump de dependência. Lipgloss e similares entram depois, se o design da saída exigir cor e caixa, e nesse momento o custo já é conhecido.

### Sem viper e sem testify

Viper resolveria um problema de precedência de configuração que este projeto não tem: flag, env e um TOML lido uma vez. `go-cmp` dá diff estruturado nas structs do modelo, que é o caso difícil de teste; o resto resolve com `testing` puro. Duas dependências a menos numa árvore que o DBA vai auditar antes de autorizar a execução contra produção.

### `testcontainers` atrás de build tag desde já

Nesta change não há teste de integração. A tag `//go:build integration` e o alvo separado no Makefile são criados mesmo assim, para que a fase 2 não precise reorganizar a suíte, e para que `go test ./...` permaneça rápido e sem Docker para quem só quer mandar um perfil de idioma novo.

## Risks / Trade-offs

**Conjunto de formas pode gerar casamento espúrio** → a pontuação da fase 3 penaliza forma obtida por normalização agressiva, e a validação da fase 5 derruba coincidência. O risco é de ruído, não de erro final, e o projeto aceita falso negativo mas não falso positivo confirmado.

**Nomes de campo do modelo viram contrato público do JSON na fase 7** → assumido conscientemente. Mudar depois exige incremento de `schema_version`. Por isso os nomes são decididos aqui com a fase 7 em mente, e não improvisados.

**Perfil `pt-br` nunca vai cobrir todo plural irregular do português** → aceito. A tabela de casos cresce por contribuição, e o conjunto de formas sempre inclui o nome original, então o pior caso de uma regra ausente é o candidato não nascer, nunca um candidato errado nascer.

**Regra de vazamento verificada por teste que varre serialização pode dar falsa sensação de segurança**, porque só cobre o que a suíte exercita → mitigado tornando a varredura parte também do teste de integração da fase 2 em diante, sobre dados reais das fixtures, e não apenas sobre structs montadas à mão.

**Piso de versão do Go fecha porta para quem está em distro antiga** → piso em 1.25, build na estável corrente. Não usar recurso de linguagem recente sem necessidade real.

## Migration Plan

Não aplicável. Repositório sem código, sem usuário e sem release.

## Open Questions

A licença ainda está indefinida — o README marca TBD. Não bloqueia esta change, mas precisa estar resolvida antes de qualquer artefato público. A recomendação é Apache-2.0 pela concessão explícita de patente, que é o que jurídico de empresa procura antes de autorizar uma ferramenta em pipeline.

O conjunto exato de sinais e seus pesos é da fase 3. Esta change define apenas o tipo `Signal` e o enum de origens, sem fixar peso.
