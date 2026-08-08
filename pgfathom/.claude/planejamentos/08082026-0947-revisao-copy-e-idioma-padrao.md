# Revisão de copy e idioma padrão em inglês

## Objetivo
Remover marcas de escrita que soam geradas por IA e eliminar a mistura de português dentro da versão em inglês da landing, além de abrir o site em inglês por padrão.

## Contexto
O site (`src/i18n/copy.en.ts` / `copy.pt.ts`) usa travessão (`—`) como conector de frase em 21 strings visíveis, sempre no mesmo padrão "frase — esclarecimento técnico", repetido em ambos os idiomas. Isoladamente cada frase é específica e técnica, mas a repetição do mesmo recurso de pontuação 21 vezes é o tipo de uniformidade que lê como texto gerado por IA.

Além disso, o painel "target output" do Hero e o exemplo fixo da seção Verdicts mostram nomes de tabela/coluna sempre em português (`os_servico.resp_tecnico → funcionario.id`, `pedido.cliente_id → cliente.id`, `lancamento.conta_id → conta.id`), independente do idioma selecionado. O mesmo texto em português aparece citado dentro de frases em inglês fluido (`problemP1`, `joinMining`), o que é exatamente a mistura de idioma reportada.

Por fim, o idioma inicial hoje é decidido por `navigator.language`: visitantes com navegador em português caem automaticamente em pt-BR. O pedido é abrir sempre em inglês por padrão.

## Solução proposta
**Idioma padrão:** remover a detecção por `navigator.language` — o fallback passa a ser sempre `'en'`, mantendo a preferência salva em `localStorage` para quem já trocou de idioma antes (evita forçar de volta para inglês quem escolheu português conscientemente).

**Painel de output localizado:** os identificadores de tabela/coluna do painel do Hero e do exemplo da seção Verdicts passam a vir do dicionário de copy (traduzidos), em vez de um dado único fixo em português. Os nomes em inglês reaproveitam o mapeamento que já existe no dicionário `graph` (ex.: `os_servico.resp_tecnico` → `work_order.tech_lead`, `funcionario` → `employee`), preservando a mesma "pegada" pedagógica do exemplo original (nomes que não se parecem entre si). Números (percentuais, contagem de órfãos) e o vocabulário fixo da CLI (`BROKEN`/`CONFIRMED`/`WEAK`, linhas de cobertura, aviso de amostragem) continuam invariantes — são o "vocabulário" literal da ferramenta, não copy da página, e não foi isso que foi reportado como misturado.

Trade-off assumido: o comentário atual em `data.ts` ("Deliberately NOT translated") descreve uma decisão de design anterior (painel = transcript literal e imutável de uma análise real). Essa decisão é revertida especificamente para os identificadores de tabela/coluna, a pedido explícito do usuário — o comentário será atualizado para não ficar desatualizado.

**Travessões:** os 21 pontos localizados em `copy.en.ts` (e os 21 equivalentes em `copy.pt.ts`) trocam o travessão por ponto final, vírgula, dois-pontos ou parênteses, conforme o papel gramatical de cada frase — variando a pontuação em vez de repetir sempre o mesmo recurso. Nenhum texto novo é adicionado; é reescrita de pontuação e, no caso de `problemP1`/`joinMining` (inglês), troca dos identificadores citados.

## Mudanças

### 1. [LanguageContext.tsx:18-27](src/i18n/LanguageContext.tsx#L18) — idioma padrão sempre inglês

**O que muda:** `readInitialLang` deixa de consultar `navigator.language`; o fallback (quando não há preferência salva) passa a ser sempre `'en'`.

**Por quê:** pedido explícito de abrir a landing em inglês por padrão, independente do idioma do navegador.

**Efeito:** primeira visita sempre carrega em inglês. Quem já trocou de idioma antes continua vendo a própria escolha salva em `localStorage`.

### 2. [types.ts:83-85](src/i18n/types.ts#L83) — novos campos no tipo `Copy`

**O que muda:** adiciona `reportRows` (lista de pares `child`/`parent` para as 6 linhas do painel de output) e `joinMiningExample` (par `child`/`parent` para o exemplo fixo da seção Verdicts) ao tipo `Copy`.

**Por quê:** esses identificadores passam a ser texto localizado, então precisam de um lugar no dicionário de copy — hoje moram como dado único em `data.ts`.

**Efeito:** nenhum, é só o contrato de tipos; o TypeScript passa a exigir esses campos nos dois dicionários.

### 3. [copy.en.ts](src/i18n/copy.en.ts) — novos valores + reescrita de pontuação

**O que muda:** adiciona `reportRows` e `joinMiningExample` com os nomes em inglês (`work_order.tech_lead → employee.id`, `order.customer_id → customer.id`, `ledger_entry.account_id → account.id`, `order_item.order_id → order.id`, `address.city_id → city.id`, `document.entity_id → entity.id`). Reescreve as 21 strings que usam travessão (`docTitle`, `heroLead`, `noticeBody`, `facts[1].detail`, `terminalCaption`, `problemP1`, `problemP2`, `verdicts[1].output`, `verdicts[2].output`, `joinMining`, `safety[0/1/3/4].detail`, `stages[2/5].detail`, `profilesP2`, `ciP1`, `priorLead`, `benchLead`, `metricBody`), trocando o travessão por pontuação variada. `problemP1` e `joinMining` também trocam `cliente_id`/`resp_tecnico`/`funcionario` pelos equivalentes em inglês (`customer_id`/`tech_lead`/`employee`), já que essas frases citam o mesmo exemplo mostrado no painel.

**Por quê:** eliminar a mistura de idioma na versão em inglês e reduzir o padrão repetitivo de pontuação.

**Efeito:** nenhuma mudança de sentido; muda pontuação e (em dois pontos) os identificadores citados.

### 4. [copy.pt.ts](src/i18n/copy.pt.ts) — novos valores + reescrita de pontuação

**O que muda:** adiciona `reportRows` e `joinMiningExample` com os nomes já usados hoje em português (`os_servico.resp_tecnico → funcionario.id` etc.). Reescreve as 21 strings equivalentes às da versão em inglês, trocando o travessão por pontuação variada. Identificadores em português não mudam (já estavam corretos no próprio idioma).

**Por quê:** manter os dois dicionários no mesmo padrão de escrita depois da mudança em inglês.

**Efeito:** nenhuma mudança de sentido; muda só pontuação.

### 5. [data.ts:1-98](src/content/data.ts#L1) — painel de output vira parcialmente localizado

**O que muda:** o comentário de cabeçalho é atualizado para explicar que só o vocabulário fixo da CLI (tags, linhas de cobertura, aviso de amostragem) permanece invariante — os identificadores de tabela/coluna agora vêm do dicionário de copy. O tipo `ReportRow` (variante `row`) perde os campos `child`/`parent`; cada linha de `REPORT` mantém só `ratio`, `note`, `tone` (e `kind`).

**Por quê:** os nomes de tabela/coluna deixam de ser dado único fixo e passam a ser texto localizado — evita ter a mesma informação duplicada em dois lugares com fontes de verdade diferentes.

**Efeito:** `REPORT` sozinho não é mais suficiente para renderizar o painel; passa a exigir o `reportRows` do idioma ativo (mudança acoplada ao item 6).

### 6. [Hero.tsx:210-229](src/components/sections/Hero.tsx#L210) — painel lê nomes do dicionário de copy

**O que muda:** ao percorrer `REPORT`, cada linha do tipo `row` passa a exibir `child`/`parent` a partir de `t.reportRows`, pareado pela posição entre as linhas de dado (ignorando as linhas de cabeçalho `BROKEN`/`CONFIRMED`/`WEAK`, que continuam vindo de `REPORT`).

**Por quê:** ligar a UI ao novo campo localizado.

**Efeito:** em inglês, o painel mostra `work_order.tech_lead → employee.id` etc.; em português, continua mostrando `os_servico.resp_tecnico → funcionario.id` etc., como hoje.

### 7. [Verdicts.tsx:64-71](src/components/sections/Verdicts.tsx#L64) — exemplo fixo vira localizado

**O que muda:** o texto fixo `os_servico.resp_tecnico → funcionario.id`, hoje hardcoded em JSX, passa a vir de `t.joinMiningExample.child` / `t.joinMiningExample.parent`.

**Por quê:** era o outro ponto onde português aparecia mesmo com o site em inglês.

**Efeito:** em inglês mostra `work_order.tech_lead → employee.id`; em português continua mostrando o texto atual.

### 8. [index.html:18](index.html#L18) — título da aba sem travessão

**O que muda:** troca o travessão do `<title>` por dois-pontos, para acompanhar a reescrita do `docTitle` em `copy.en.ts`.

**Por quê:** consistência entre o título estático (antes do JS montar) e o título aplicado depois pelo `LanguageProvider`.

**Efeito:** título da aba igual em ambos os momentos (antes e depois da hidratação).

## O que não muda
- Vocabulário fixo da CLI (`BROKEN`/`CONFIRMED`/`WEAK`, `REPORT_PROFILE`, `REPORT_COVERAGE`, `REPORT_WARNING`) continua sempre em inglês, nos dois idiomas do site — é o output literal da ferramenta, não foi apontado como mistura de idioma.
- Números do painel (percentuais, contagem de órfãos) não mudam.
- Nomes de tabela usados na animação do grafo (`SchemaGraph`/`TrailingTable`) já eram localizados via o dicionário `graph` existente — não precisam de mudança.
- Texto de comentários no código (`.tsx`, `data.ts`) não é copy visível na página e não entra nesta revisão.
- Nenhuma nova biblioteca ou dependência.
- Switcher de idioma, persistência em `localStorage` e o restante da lógica de `LanguageContext` continuam iguais.

## Fluxo (idioma padrão)
```
Primeira visita
  └─ localStorage tem 'pgfathom-lang'?
       sim → usa o valor salvo (en ou pt)
       não → usa 'en' (independente do navigator.language)
```

[Planejamento](.claude/planejamentos/08082026-0947-revisao-copy-e-idioma-padrao.md)
[Plano de ação](.claude/planos-de-acao/08082026-0947-revisao-copy-e-idioma-padrao.md)
