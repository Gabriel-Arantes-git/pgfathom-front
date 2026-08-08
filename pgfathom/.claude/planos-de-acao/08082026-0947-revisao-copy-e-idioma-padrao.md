# Revisão de copy e idioma padrão em inglês — Plano de Ação

## Ordem

### 1. src/i18n/LanguageContext.tsx

**Ação:** editar

**Critério:** `readInitialLang` mantém a leitura de `localStorage[STORAGE_KEY]` como primeira prioridade (retorna `'en'`/`'pt'` se salvo). Remove o fallback por `navigator.language`. Quando não há valor salvo (ou `window` é `undefined`), retorna sempre `'en'`.

**Localização:** função `readInitialLang`, linhas 18-27.

---

### 2. src/i18n/types.ts

**Ação:** editar

**Critério:** adiciona ao tipo `Copy` dois campos novos, próximos de `graph` (linha 85): `reportRows`, um array de exatamente 6 objetos `{ child: string; parent: string }` (um por linha de dado do painel de output, na mesma ordem de `REPORT` em `data.ts`); e `joinMiningExample`, um único objeto `{ child: string; parent: string }` (o exemplo fixo da seção Verdicts).

**Localização:** após o campo `graph` (linha 85), antes de `problemEyebrow`.

---

### 3. src/content/data.ts

**Ação:** editar

**Critério:** o tipo `ReportRow` (variante `kind: 'row'`) perde os campos `child` e `parent`; mantém `kind`, `ratio`, `note`, `tone`. Cada uma das 6 entradas `kind: 'row'` em `REPORT` (linhas 40-89) remove os campos `child`/`parent`, mantendo `ratio`, `note`, `tone` com os valores atuais e na mesma ordem (BROKEN: os_servico/pedido/lancamento; CONFIRMED: item_pedido/endereco; WEAK: documento). O comentário de cabeçalho do arquivo (linhas 1-6) é reescrito para deixar explícito que apenas o vocabulário fixo da CLI (tags `BROKEN`/`CONFIRMED`/`WEAK`, `REPORT_PROFILE`, `REPORT_COVERAGE`, `REPORT_WARNING`) é literal e invariante entre idiomas — os identificadores de tabela/coluna passaram a viver no dicionário de copy (`reportRows`/`joinMiningExample`), localizados por idioma.

**Localização:** comentário no topo (linhas 1-6), tipo `ReportRow` (linhas 27-36), array `REPORT` (linhas 38-90).

**Depende de:** passo 2 (o novo formato só faz sentido com `reportRows` existindo no tipo `Copy`).

---

### 4. src/i18n/copy.en.ts

**Ação:** editar

**Critério — novos campos:**
`reportRows` (mesma ordem das 6 linhas de `REPORT`):
1. `{ child: 'work_order.tech_lead', parent: '→ employee.id' }`
2. `{ child: 'order.customer_id', parent: '→ customer.id' }`
3. `{ child: 'ledger_entry.account_id', parent: '→ account.id' }`
4. `{ child: 'order_item.order_id', parent: '→ order.id' }`
5. `{ child: 'address.city_id', parent: '→ city.id' }`
6. `{ child: 'document.entity_id', parent: '→ entity.id' }`

`joinMiningExample`: `{ child: 'work_order.tech_lead', parent: '→ employee.id' }` (mesmo conteúdo da primeira linha de `reportRows`).

**Critério — reescrita de pontuação (remover travessão, variar entre ponto, vírgula, dois-pontos ou parênteses conforme o papel gramatical da frase, sem alterar o sentido):**
- `docTitle` (linha 4): travessão → dois-pontos.
- `heroLead` (linha 22): travessão → vírgula + "then" no lugar de "and".
- `noticeBody` (linha 28): travessão antes de "no numbers are claimed until then" → ponto final, nova frase.
- `facts[1].detail` "Stays in memory" (linha 31): travessão → vírgula.
- `terminalCaption` (linha 35): travessão → dois-pontos.
- `problemP1` (linha 62): travessão → ponto final, nova frase; troca `cliente_id`/`cliente.id` por `customer_id`/`customer.id`.
- `problemP2` (linha 64): travessão antes de "so they probably already did" → vírgula.
- `verdicts[1].output`, o de tag `CONFIRMED` (linha 86): travessão antes de "plus CREATE INDEX CONCURRENTLY" → vírgula.
- `verdicts[2].output`, o de tag `WEAK` (linha 93): travessão antes de "a polymorphic pair" → dois-pontos.
- `joinMining` (linha 97): travessão → dois-pontos; troca `resp_tecnico`/`funcionario` por `tech_lead`/`employee`.
- `safety[0].detail` "Read-only, structurally" (linha 107): travessão → ponto final, nova frase ("There is no write mode...").
- `safety[1].detail` "Your data never leaves memory" (linha 112): travessão → vírgula.
- `safety[3].detail` "No claim without evidence" (linha 122): travessão → dois-pontos.
- `safety[4].detail` "Silence is never a clean bill of health" (linha 127): travessão → dois-pontos.
- `stages[2].detail` "Generate candidates" (linha 150): travessão → dois-pontos.
- `stages[5].detail` "Validate against the data" (linha 168): travessão → vírgula.
- `profilesP2` (linha 178): travessão antes de "every form is tried" → ponto final, nova frase.
- `ciP1` (linha 185): travessão antes de "this is the shape" → dois-pontos.
- `priorLead` (linha 195): travessão → vírgula.
- `benchLead` (linha 237): as duas ocorrências de travessão ao redor de "GitLab, Odoo, Discourse, Redmine, Mastodon" → parênteses.
- `metricBody` (linha 247): as duas ocorrências de travessão ao redor de "GitLab, Odoo, Discourse, Redmine, Mastodon" → parênteses.

**Localização:** campos indicados acima, por número de linha (linhas do arquivo antes desta edição).

**Depende de:** passo 2.

---

### 5. src/i18n/copy.pt.ts

**Ação:** editar

**Critério — novos campos** (mesma posição de `reportRows`/`joinMiningExample` que em `copy.en.ts`, valores em português — mantêm os identificadores já usados hoje em `data.ts`):
`reportRows`:
1. `{ child: 'os_servico.resp_tecnico', parent: '→ funcionario.id' }`
2. `{ child: 'pedido.cliente_id', parent: '→ cliente.id' }`
3. `{ child: 'lancamento.conta_id', parent: '→ conta.id' }`
4. `{ child: 'item_pedido.pedido_id', parent: '→ pedido.id' }`
5. `{ child: 'endereco.municipio_id', parent: '→ municipio.id' }`
6. `{ child: 'documento.entidade_id', parent: '→ entidade.id' }`

`joinMiningExample`: `{ child: 'os_servico.resp_tecnico', parent: '→ funcionario.id' }`.

**Critério — reescrita de pontuação** (mesmo padrão do passo 4, aplicado às 21 strings equivalentes; identificadores em português não mudam, só a pontuação):
- `docTitle` (linha 4), `heroLead` (linha 22), `noticeBody` (linha 28), `facts[1].detail` "Fica em memória" (linha 31), `terminalCaption` (linha 35), `problemP1` (linha 58), `problemP2` (linha 60), `verdicts[1].output` CONFIRMED (linha 82), `verdicts[2].output` WEAK (linha 89), `joinMining` (linha 93), `safety[0].detail` (linha 103), `safety[1].detail` (linha 108), `safety[3].detail` (linha 118), `safety[4].detail` (linha 123), `stages[2].detail` (linha 146), `stages[5].detail` (linha 164), `profilesP2` (linha 174), `ciP1` (linha 181), `priorLead` (linha 191), `benchLead` (linha 233), `metricBody` (linha 243) — mesma técnica de substituição do passo 4 (ponto, vírgula, dois-pontos ou parênteses conforme o papel gramatical de cada trecho).

**Localização:** campos indicados acima, por número de linha (linhas do arquivo antes desta edição).

**Depende de:** passo 2.

---

### 6. src/components/sections/Hero.tsx

**Ação:** editar

**Critério:** no bloco que percorre `REPORT` (`REPORT.map`, linhas 210-229), ao renderizar uma linha `kind === 'row'`, o `child`/`parent` exibidos passam a vir de `t.reportRows`, indexado pela posição da linha entre as linhas de dado (contando só as de `kind: 'row'`, ignorando as de `kind: 'head'`), em vez de `row.child`/`row.parent` (que deixaram de existir em `data.ts`). `ratio`, `note` e `tone` continuam vindo de `row` normalmente. As linhas `head` (`BROKEN`/`CONFIRMED`/`WEAK` etc.) continuam exibindo `row.text` sem alteração.

**Localização:** dentro do `REPORT.map` em `Hero.tsx`, linhas 210-229.

**Depende de:** passos 3, 4, 5.

---

### 7. src/components/sections/Verdicts.tsx

**Ação:** editar

**Critério:** o texto hoje hardcoded (`os_servico.resp_tecnico → funcionario.id`, linha 66) passa a ler `t.joinMiningExample.child` para o trecho antes da seta e `t.joinMiningExample.parent` para o trecho da seta em diante, mantendo a mesma estrutura visual (span com a seta destacada em `text-accent`).

**Localização:** linha 66, dentro do bloco de `code` da seção "no name-matching heuristic".

**Depende de:** passos 4, 5.

---

### 8. index.html

**Ação:** editar

**Critério:** o conteúdo do `<title>` (linha 18) passa a usar dois-pontos no lugar do travessão, ficando idêntico ao novo valor de `docTitle` em `copy.en.ts`.

**Localização:** linha 18.

**Depende de:** passo 4.

---

## Validação
Após implementar: rodar o build/typecheck do projeto (Vite + TypeScript) para garantir que `copy.en.ts` e `copy.pt.ts` satisfazem o tipo `Copy` atualizado (o compilador acusa qualquer campo faltando em um dos dois dicionários). Conferir visualmente no navegador, nos dois idiomas: painel do Hero, exemplo da seção Verdicts, e que a primeira carga (sem `localStorage` prévio) abre em inglês.
