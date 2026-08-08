# Landing page pgfathom

## Objetivo

Landing page dark, one-page, para o projeto open source `pgfathom` — CLI Go que descobre foreign keys não declaradas em schemas PostgreSQL legados e as prova contra os dados.

## Contexto

`pgfathom` (`C:\Users\user\Documents\Projects\pgfathom`) é uma CLI Go, Apache-2.0, pré-release. Fases 1 a 3.5 do roadmap têm código (`internal/model`, `internal/profile`, `internal/db`, `internal/catalog`, `internal/infer`, `internal/report`, comandos `audit` e `discover`); fases 4 a 8 são especificação. O README ainda diz "no code yet" — está desatualizado em relação ao repositório.

O produto tem uma ética explícita que a landing precisa respeitar, senão contradiz o próprio produto:

- **"No claim without evidence"** — nada de número inventado, nada de logo de empresa que não usa, nada de depoimento fabricado.
- **"Silence is never reported as a clean bill of health"** — a landing declara o que ainda não existe em vez de esconder.
- **"Small dependency tree, on purpose"** — quatro dependências no binário. Uma landing com 40 pacotes npm contradiz o discurso.

Consequência direta no escopo: **não existe seção de pricing, nem de depoimentos, nem de logos de clientes.** Os slots equivalentes do fluxo de referência são ocupados por Roadmap, Prior Art e Contribuição. Os únicos números publicados são os que já existem medidos no `docs/ROADMAP.md` (fase 3.5), e vão rotulados como medição de fase 3, não como benchmark oficial.

O frontend hoje (`pgfathom-front/pgfathom/pgfathom`) é um scaffold Vite + React 19 + TypeScript zerado, com o template padrão do Vite em `App.tsx`, CSS puro e um sprite `public/icons.svg` com 6 símbolos.

## Referências estéticas analisadas

**`asthetic`** (Next.js + Tailwind v4 + shadcn) — o que vai ser aproveitado:

| Recurso | Arquivo | Uso na pgfathom |
|---|---|---|
| Esfera ASCII em canvas 2D | `components/landing/animated-sphere.tsx` | Hero + CTA final, retematizada |
| Navbar que encolhe e vira pill flutuante no scroll | `components/landing/navigation.tsx` | Base da navbar glassed |
| SVGs animados com `<animate>` / `<animateMotion>` puro | `components/landing/features-section.tsx` | Ícones das seções Safety e Pipeline |
| Revelação char-a-char com blur | `globals.css` (`animate-char-in`) | Headline do hero, veredito rotativo |
| Digitação de código char-a-char por linha | `how-it-works-section.tsx` | Terminal do pipeline |
| Contador animado com IntersectionObserver | `metrics-section.tsx` | Seção de medição |
| Marquee infinito | `globals.css` + `integrations-section.tsx` | Trust bar, corpus de benchmark |
| Overlay de ruído SVG | `globals.css` (`noise-overlay`) | Textura global |
| Spotlight radial seguindo o mouse | `cta-section.tsx` | CTA final |
| Grid de linhas de fundo | `hero-section.tsx` | Hero |

**`pointer-ai-landing-page`** (Next.js + Tailwind + framer-motion) — o que vai ser aproveitado:

| Recurso | Arquivo | Uso na pgfathom |
|---|---|---|
| Fluxo/ordem das seções | `app/page.tsx` | Espinha dorsal da narrativa |
| Preview flutuante sobrepondo o hero (`absolute bottom-[-400px]`) | `app/page.tsx:20` | Terminal preview sobrepondo o hero |
| `AnimatedSection` — fade + `y:20` + `scale:0.98`, `once: true` | `components/animated-section.tsx` | Padrão de entrada de toda seção |
| Bento card glassed (`rgba(...,0.08)` + `backdropFilter: blur(4px)` + borda translúcida + gradiente) | `components/bento-section.tsx:9-24` | Todos os cards glassed do site |
| Ilustração própria por card do bento | `components/bento/*` | Cards de capacidade |
| Container `max-w-[1320px]` com ritmo de `mt-16` | `app/page.tsx` | Grid e ritmo vertical |

## Solução proposta

### Decisão 1 — Stack

O projeto alvo é Vite + React 19 + TS, sem Tailwind, sem framer-motion. As duas referências são Next + Tailwind. Três caminhos:

**Resolvido: Tailwind v4.** Aprovado e já instalado.

```
npm i -D tailwindcss @tailwindcss/vite     # 12 pacotes, build-time apenas
vite.config.ts:4,11                        # plugin tailwindcss() registrado
npm run build                              # ✓ verde em 447ms
```

Uma dependência de build (que não vai para o bundle) compra fidelidade quase literal às duas referências e corta o maior custo do trabalho. O bundle final continua sendo React + CSS. Tokens declarados via `@theme` em `src/index.css` — sem `tailwind.config.js`, que o v4 não usa.

Animação: **sem framer-motion.** O `AnimatedSection` do pointer é reimplementado com `IntersectionObserver` + transição CSS — que é exatamente o que o `asthetic` já faz sem lib nenhuma. Assim o único runtime é React.

### Decisão 2 — Bilíngue EN / PT-BR com seletor de bandeira

O site é **bilíngue**, com `en` como padrão e troca por um switch de bandeira dentro da página. Não é tradução parcial: a troca altera o conteúdo inteiro do site.

**Por que os dois idiomas se justificam aqui.** O README, os docs e a CLI são em inglês e o público de DBA de banco legado é internacional — mas o argumento central do produto (`most schema tools assume English — the databases that need this tool most often aren't`) e a medição de recall vêm de bases brasileiras de gestão pública. Um site que defende naming profiles não-inglesas e só existe em inglês é incoerente com a própria tese.

#### Arquitetura de i18n — sem dependência

`react-i18next` custaria 2 dependências e um runtime de interpolação que uma landing estática não usa. Como o plano já centraliza todo o texto em `content/copy.ts`, a solução é um dicionário tipado + context:

```
src/i18n/
  types.ts        type Lang = 'en' | 'pt'
                  type Copy  = { ... }   ← forma única, derivada de copy.en
  copy.en.ts      const en: Copy = { ... }
  copy.pt.ts      const pt: Copy = { ... }   ← TS acusa chave faltando
  LanguageContext.tsx   provider + useLang() + useCopy()
```

Vantagens sobre a lib: **zero dependência** (coerente com o discurso de `go.mod` curto do produto) e **type safety** — se uma chave existir em `en` e faltar em `pt`, o build quebra. Com `react-i18next` a chave faltando vira string crua em produção.

**Comportamento.**
- Ordem de resolução na primeira visita: `localStorage['pgfathom-lang']` → `navigator.language.startsWith('pt')` → `en`.
- A escolha persiste em `localStorage`.
- `document.documentElement.lang` é atualizado na troca (`en` / `pt-BR`) — importa para leitor de tela e para o hífen do navegador.
- Troca é instantânea, sem reload e sem perder a posição de scroll. Textos animados char-a-char **reanimam** ao trocar de idioma (a chave do React muda), o que dá uma transição de graça.
- `<title>` e `<meta name="description">` também trocam.

**O que não é traduzido, por decisão.** Saída de terminal, SQL, TOML, nomes de flag da CLI, nomes de coluna/tabela dos exemplos e o bloco `\d pedido`. São artefatos reais da ferramenta — traduzir `BROKEN` para `QUEBRADO` seria mostrar uma saída que a CLI não produz, e o site inteiro é construído em cima de não mentir sobre o output. Traduzem-se as **legendas, títulos e comentários ao redor** desses blocos.

#### O switch de bandeira

Controle segmentado glassed, no formato de pill, com as duas bandeiras em SVG inline circular de 20px:

```
┌──────────────┐
│  (🇺🇸) (🇧🇷)  │   ativo: opacidade 1 + ring 1px --accent + escala 1
└──────────────┘   inativo: opacidade .45 + grayscale(.6), hover volta a 1
```

- **SVG inline, não emoji.** Emoji de bandeira (`🇧🇷`) não renderiza no Windows — o sistema mostra as letras `BR`. Dois componentes `<FlagBR />` e `<FlagUS />` em `components/ui/flags/`, ~10 linhas de path cada, recortados em círculo por `clipPath`.
- Transição 220ms `cubic-bezier(.4,0,.2,1)`, com o ring do ativo deslizando entre as duas posições em vez de aparecer/sumir.
- `role="group"` + `aria-label`, cada bandeira é um `<button>` com `aria-pressed` e `aria-label` no idioma de destino (`Mudar para português` / `Switch to English`), e `<title>` dentro do SVG.
- Bandeira dos EUA para inglês (a documentação e as issues do projeto seguem convenção US). É uma troca de um componente se preferir a do Reino Unido.

**Onde aparece:** navbar desktop (à direita, antes do `GitHub`), rodapé do overlay mobile (em tamanho maior, 28px, com o nome do idioma escrito ao lado) e footer.

### Decisão 3 — Sistema de cor

Escala de luminosidade sobre matiz quente fixo, dark único (sem light mode).

```
--ink-900  #0f0c0a   fundo mais profundo (seções invertidas, footer)
--ink-800  #14100e   fundo alternado
--ink-700  #1b1512   base (cor pedida)
--ink-600  #221a16   surface
--ink-500  #2a201b   surface elevada
--line     rgba(244,239,232,0.10)   bordas
--line-hi  rgba(244,239,232,0.22)   bordas hover/ativa

--paper    #f4efe8   texto principal (cor da logo)
--paper-70 rgba(244,239,232,0.70)   texto secundário
--paper-45 rgba(244,239,232,0.45)   texto muted / mono eyebrow

--accent   #e0473d   acento principal
--accent-d #b0322a   acento pressionado
--accent-g rgba(224,71,61,0.14)     glow / fundo de badge
```

Semânticas dos três veredictos do produto (o único ponto onde entra matiz fora da paleta base, porque a informação exige distinção):

```
--broken     #e0473d   BROKEN     (é o próprio acento — é o achado que importa)
--confirmed  #5fa87a   CONFIRMED  verde dessaturado
--weak       #d9a441   WEAK       âmbar
```

Distribuição alvo ~55-33-12 (landing, conforme `design-ui.md`): neutros quentes dominam, acento concentrado em CTA, veredito BROKEN, sublinhados e a esfera.

Glass usado com parcimônia (navbar, bento cards, badges de status, painel de código flutuante), na fórmula do pointer adaptada ao dark:

```css
background: rgba(244, 239, 232, 0.045);
backdrop-filter: blur(14px) saturate(120%);
border: 1px solid rgba(244, 239, 232, 0.10);
box-shadow: 0 1px 0 rgba(244,239,232,0.06) inset, 0 24px 60px -30px rgba(0,0,0,0.8);
```

### Decisão 4 — Tipografia

- **Display:** Instrument Serif — headlines. O serif dá o tom editorial/instrumento de medição que combina com "sound the depth", e é o que dá caráter ao `asthetic`.
- **Mono:** JetBrains Mono — terminal, SQL, TOML, eyebrows, números. Peso alto no site inteiro: o produto *é* um terminal.
- **Sans:** Inter — corpo de texto.

Escala modular 1.25. Letter-spacing `-0.02em` em display grande, `+0.12em` em eyebrow mono uppercase.

### Decisão 5 — Ícones

O sprite `public/icons.svg` tem 6 símbolos: `bluesky-icon`, `discord-icon`, `documentation-icon`, `github-icon`, `social-icon`, `x-icon` — com `fill` hardcoded (`#08060d`, `#aa3bff`), que somem no dark. **Correção necessária:** normalizar os símbolos para `fill="currentColor"` / `stroke="currentColor"` num sprite próprio, para herdarem o tema. Uso: navbar (GitHub, Docs), footer (todos), seção Community.

`src/assets/icons/pgfathom-logo.svg` é a marca (mergulhador/sonda em `#B01E23`). Vai inline no React para o `#B01E23` virar `currentColor` e receber `--accent`, com micro-animação: os dois pontos e o círculo da "sonda" pulsam em cascata no hover (a própria logo é um sonar).

Ícones de feature **não existem** no repositório — serão SVGs autorais animados no estilo do `features-section.tsx` do `asthetic`: line-art 2px, `currentColor`, animação via `<animate>` SMIL (sem JS, sem lib), viewBox 200×160.

## Estrutura da página — seção a seção

Ordem derivada do fluxo do `pointer` (hero → preview flutuante → prova social → bento → destaque → pricing → depoimentos → FAQ → CTA → footer), com os slots comerciais substituídos por equivalentes honestos para open source.

---

### 00 — Navbar `glassed`, sempre visível

**Comportamento.** `position: fixed`, sempre presente (diferente do `asthetic`, onde só ganha vidro depois do scroll). Dois estados:

- **Topo:** largura `max-w-[1320px]`, altura 80px, vidro bem sutil (`blur(10px)`, borda quase invisível), sem sombra.
- **Scrolled (> 20px):** encolhe para altura 56px, recolhe para `max-w-[1120px]`, ganha `top: 16px`, `border-radius: 16px`, vidro pleno (`blur(18px) saturate(130%)`), borda `--line`, sombra profunda. Transição 500ms `cubic-bezier(.4,0,.2,1)` — mesma mecânica de `navigation.tsx:27-45`.

**Conteúdo.**
- Esquerda: logo inline (24px) + wordmark `pgfathom` em mono, com um `v0.1 · pre-release` em `--paper-45`, 10px, ao lado.
- Centro: `Problem` · `How it works` · `Safety` · `Profiles` · `Roadmap` · `Docs` — 13px, `--paper-70`, hover para `--paper` com sublinhado que cresce da esquerda (`navigation.tsx:61`).
- Direita: **switch de bandeira** (🇺🇸/🇧🇷, ver Decisão 2), separador vertical 1px `--line`, `GitHub` (ícone do sprite) e botão primário **`Install`** / **`Instalar`** — pill, fundo `--accent`, texto `--ink-700`, hover `translateY(-1px)` + glow `--accent-g`.
- Barra de progresso de scroll de 1px em `--accent` colada na borda inferior da navbar.

No estado *scrolled* o switch encolhe junto (bandeiras 20px → 16px) e o separador some, mantendo a navbar respirando na altura de 56px.

**Mobile.** Overlay full-screen com links em display serif 48px entrando em cascata (75ms de delay entre eles), igual a `navigation.tsx:104-147`. O switch de bandeira fica no rodapé do overlay, em 28px, com o nome do idioma escrito ao lado (`English` / `Português`).

---

### 01 — Hero

**Estética.** `min-height: 100vh`. Fundo `--ink-700` com:
1. Grid de linhas (8 horizontais, 12 verticais) em `--paper` a 6% de opacidade — `hero-section.tsx:33-56`.
2. **Esfera ASCII** à direita, 800×800, `opacity: 0.5` — porte de `animated-sphere.tsx`, retematizada: em vez de `rgba(0,0,0,alpha)`, interpolação por profundidade entre `--accent` (pontos ao fundo) e `--paper` (pontos à frente), com o mesmo `alpha = 0.2 + (z+1)*0.4`. Charset mantido (`░▒▓█▀▄▌▐│─┤├┴┬╭╮╰╯`). Adição: um **anel de ping** concêntrico expandindo a cada 4s de `r=0` até fora do raio, `--accent` esmaecendo — a metáfora do sonar que dá nome ao projeto.
3. Vinheta radial escurecendo os cantos.
4. `noise-overlay` global a 3%.

**Textos.**

- *Eyebrow* (mono, `--paper-45`, com o traço de 32px antes): `— Open source · Apache-2.0 · read-only by construction`
- *Headline* (display serif, `clamp(3rem, 9vw, 7.5rem)`, `line-height: .9`):
  ```
  Sound the depth
  of a legacy schema.
  ```
  Entrada char-a-char com blur (`animate-char-in`, 50ms por caractere).
- *Palavra rotativa*: abaixo da headline, uma linha mono que cicla a cada 2.5s entre os três veredictos, cada um na sua cor, recomeçando a animação de blur — `hero-section.tsx:82-97`:
  `1,847 candidates → 226 validated → ` + `BROKEN` / `CONFIRMED` / `WEAK`
- *Descrição* (20px, `--paper-70`, `max-w-[560px]`):
  > pgfathom finds the relationships your database has but never declared — and proves them against the data instead of guessing from column names.
- *CTAs*: primário pill `--accent` → **`Get started`** com seta que desliza no hover; secundário outline `--line-hi` → **`Read the design doc`**.
- *Comando copiável* abaixo dos CTAs, em caixa glassed mono 13px:
  `$ go install github.com/lvcas-dotcom/pgfathom/cmd/pgfathom@latest` com botão de copiar (troca para check por 2s — padrão de `developers-section.tsx`).
- *Aviso de honestidade* (mono 11px, `--paper-45`, com um `!` em `--accent`):
  `! pre-release. terminal output on this page is the target design, not a recording.`

**Marquee inferior** (`hero-section.tsx:140-166`), 4 pares valor/label rolando:
`Go 1.25+` · *no cgo* — `4 deps` · *in the binary* — `read-only` · *structurally, no write mode* — `Apache-2.0` · *explicit patent grant*

---

### 02 — Terminal preview flutuante

Equivalente ao `DashboardPreview` do pointer, posicionado com `absolute bottom-[-400px] left-1/2 -translate-x-1/2 z-30` sobre o hero (`app/page.tsx:20-24`), entrando com o `AnimatedSection`.

**Estética.** Janela de terminal glassed, `max-w-[1000px]`, borda `--line`, `border-radius: 12px`, sombra `0 40px 120px -40px #000`. Chrome com três círculos `--paper` 20% e o título mono `pgfathom discover --schema public`. Um glow radial `--accent-g` vaza por trás da janela.

**Conteúdo.** A saída real do README (`README.md:73-98`), com as três tabelas de veredicto, cada bloco entrando em cascata e cada linha digitando char-a-char (`code-char-reveal`, 15ms/char, 80ms/linha):

```
  profile pt-br (detected) · sampled validation (100k rows/table) · 312 tables

  BROKEN — the relationship is real, the integrity is not
  ─────────────────────────────────────────────────────────────
  os_servico.resp_tecnico  → funcionario.id   99.7%   1,284 orphans
  pedido.cliente_id        → cliente.id       99.9%      37 orphans
  lancamento.conta_id      → conta.id         94.2%  11,903 orphans

  CONFIRMED — undeclared foreign key, no orphans found
  ─────────────────────────────────────────────────────────────
  item_pedido.pedido_id    → pedido.id       100.0%       0 orphans
  endereco.municipio_id    → municipio.id    100.0%       0 orphans

  WEAK — insufficient evidence to conclude
  ─────────────────────────────────────────────────────────────
  documento.entidade_id    → entidade.id      41.0%  polymorphic pair
```

Rodapé da janela: `312 tables · 298 analyzed · 14 skipped (no SELECT privilege)` e um ponto pulsante `--confirmed` com `done in 42s`.
Os rótulos `BROKEN`/`CONFIRMED`/`WEAK` recebem as cores semânticas; as contagens de órfãos em `--accent`.

---

### 03 — Trust bar (substitui "social proof")

Faixa fina, fundo `--ink-800`, marquee lento em ambas as direções com cápsulas outline (`integrations-section.tsx`), texto mono 12px `--paper-45`, hover acende a borda para `--line-hi`:

`PostgreSQL 13+` · `Go 1.25+` · `no cgo` · `4 dependencies` · `read-only session` · `TABLESAMPLE` · `pg_stat_statements` · `pg_stats prefilter` · `pt-br / en / es profiles` · `JSON model` · `reviewable .sql` · `Apache-2.0`

Título acima, mono `--paper-45`, centralizado: `Built to be approved by the DBA who is nervous about it`

---

### 04 — The problem

Narrativa em duas colunas, numeração romana à esquerda no estilo `how-it-works-section.tsx`.

*Eyebrow:* `— The problem` · *Headline:* **Old databases carry more structure than they declare.**

**Coluna esquerda — texto:**
> No foreign keys. But `cliente_id` points at `cliente.id` in every single row — the ORM of the day never created the constraint, or someone dropped constraints for a bulk load and never put them back.
>
> Two things follow. Nobody can read the model, because `\d` shows nothing and your ERD tool draws a page of disconnected boxes. And because the constraint was never there, nothing ever stopped orphan rows from getting in — so they probably already did, years ago, silently, and no one has looked.

**Callout destacado**, card glassed com borda esquerda 2px `--accent`:
> **There is a nastier variant.** A foreign key can be *declared* and still guarantee nothing, if it was created `NOT VALID` and never validated. It shows up in `\d`. It draws an arrow in your diagram. It never checked a single pre-existing row.

**Coluna direita — bloco `psql`** (mono, fundo `--ink-900`, borda `--line`), com o `Indexes:` e a ausência de `Foreign-key constraints:` sublinhada em `--accent` tracejado, e a linha `cliente_id` piscando de leve:

```
=# \d pedido
      Column     |  Type  | Nullable
-----------------+--------+----------
 id              | bigint | not null
 cliente_id      | bigint | not null
 ...
Indexes:
    "pedido_pkey" PRIMARY KEY, btree (id)
```
Legenda abaixo, mono `--paper-45`: `no "Foreign-key constraints:" block. that is the whole problem.`

---

### 05 — Three verdicts

*Eyebrow:* `— What it does` · *Headline:* **Three verdicts. Three different responses.**

Três cards glassed lado a lado (fórmula do bento do pointer), cada um com barra superior de 2px na cor do veredito, e uma **ilustração SVG animada autoral** de 200×160 no topo:

| Card | Ilustração animada | Texto |
|---|---|---|
| **BROKEN** (`--accent`) | Duas tabelas ligadas por linha tracejada com `animateMotion`; alguns nós se soltam da linha e caem — os órfãos. Contador `1,284` subindo. | *A data bug that has been in production for years and nobody knows.* pgfathom writes you the query that lists the orphans, because they have to be resolved before any constraint can be added. |
| **CONFIRMED** (`--confirmed`) | Duas tabelas; um arco se desenha entre elas via `stroke-dashoffset` animado e "trava" com um pulso; um cadeado desenha o aro. | *A foreign key someone forgot to declare.* You get the DDL, with `VALIDATE CONSTRAINT` split out so the initial `ALTER TABLE` doesn't hold a heavy lock — plus `CREATE INDEX CONCURRENTLY` on the child column when it's missing. |
| **WEAK** (`--weak`) | Um nó central com três alvos possíveis; o feixe alterna entre eles sem fixar; opacidade oscilando. | *Insufficient evidence to conclude.* Weak and rejected candidates are reported too, so you never wonder why an obvious-looking column was ignored. |

**Faixa abaixo dos cards** — bloco de código com o DDL gerado, alternando por tabs (`Confirmed → DDL` / `Broken → orphan query`), digitação char-a-char:

```sql
ALTER TABLE item_pedido
  ADD CONSTRAINT item_pedido_pedido_id_fkey
  FOREIGN KEY (pedido_id) REFERENCES pedido (id) NOT VALID;

CREATE INDEX CONCURRENTLY item_pedido_pedido_id_idx
  ON item_pedido (pedido_id);

ALTER TABLE item_pedido VALIDATE CONSTRAINT item_pedido_pedido_id_fkey;
```

---

### 06 — The one nobody else finds

Seção curta e de alto impacto, largura total, fundo `--ink-900`, com a esfera reaparecendo pequena e desfocada ao fundo.

*Headline* (display, centralizado):
> `os_servico.resp_tecnico → funcionario.id`

em mono grande, com a seta em `--accent` desenhando-se da esquerda para a direita ao entrar em viewport.

*Texto abaixo:*
> No name-matching heuristic in the world finds that one — `resp_tecnico` looks nothing like `funcionario`. pgfathom finds it by reading the join predicates out of your own view and function definitions.

Ao lado, painel mono mostrando o `CREATE VIEW` de onde a evidência sai, com o `ON os_servico.resp_tecnico = f.id` destacado em `--accent` com fundo `--accent-g`.

---

### 07 — How it works (painel invertido)

Contraponto de luminosidade — a única seção em `--ink-900` puro com padrão de linhas diagonais a 3% (`how-it-works-section.tsx:74-84`).

*Eyebrow:* `— Process` · *Headline:* **Six stages. Each one cheaper than the next.**

**Pipeline horizontal no topo**, com os seis nós ligados; o nó ativo acende em `--accent` e a linha entre eles se preenche progressivamente:

```
catalog → usage evidence → candidates → scoring → stats prefilter → validation
```

**Coluna esquerda — steps clicáveis com autoplay de 5s** e barra de progresso animada no item ativo (`how-it-works-section.tsx:108-141`). Item inativo a 40% de opacidade; título desliza 8px à direita no hover.

| # | Título | Descrição |
|---|---|---|
| I | Read the catalog | Tables, columns, keys, indexes, comments, declared foreign keys *and their validation state*, usage statistics with their reset timestamp. |
| II | Mine usage evidence | Join predicates extracted from view definitions, function bodies and `pg_stat_statements`. A view that joins two columns is proof that your code treats them as related, whatever they happen to be named. |
| III | Generate candidates | Column-name affixes are stripped and matched against depluralized table names using a naming profile — a config file, not hardcoded rules. |
| IV | Score on metadata alone | Exact name match, type identity, target ambiguity, existing index, comment mentions. Weak candidates are dropped before anything touches data. |
| V | Prefilter on planner statistics | If the child column has more distinct values than the parent has rows, full containment is arithmetically impossible. Free, from `pg_stats`, no I/O. |
| VI | Validate against the data | One aggregate per surviving candidate — never fetching rows, only counts. Containment reported by row *and* by distinct value. |

**Coluna direita — janela sticky** (`lg:sticky lg:top-32`) que troca de conteúdo conforme o step ativo, com a digitação char-a-char e o rodapé de status com ponto pulsante. Cada step mostra a query/saída real correspondente (SQL de catálogo no I, extração de join no II, lista de candidatos com score no IV, `EXPLAIN`-like no V, o agregado de contenção no VI).

---

### 08 — Safety

Seção-chave: é o argumento que decide se a ferramenta pode ser apontada para produção. Fundo `--ink-700` com uma faixa `--ink-800` sutil.

*Eyebrow:* `— Safety` · *Headline:* **Pointed at a production database owned by someone who is nervous about it.**
*Subtítulo:* Every guarantee below is a hard requirement in the spec, not a goal.

Cinco cards glassed em grid 2+3, cada um com **SVG animado autoral** 64×64 (line-art `currentColor`, animação SMIL):

| Ícone animado | Título | Texto |
|---|---|---|
| Seta de escrita batendo num escudo e ricocheteando, em loop | **Read-only, structurally** | pgfathom never issues a statement that modifies the database under analysis — there is no write mode, under any flag, in any phase. The session sets `default_transaction_read_only`. It emits `.sql` files for you to review and run yourself. |
| Caixa de memória com valores entrando e só números/porcentagens saindo; os valores se dissolvem na borda | **Your data never leaves memory** | What comes *out* are counts, ratios and object names — never a value from your tables, in any output, log, JSON field or error message. Enforced by a test that serializes every structure and scans the result, not by code review. |
| Medidor com ponteiro que sobe e é contido por um limite; ampulheta girando | **It won't take your server down** | Every validation query runs under `statement_timeout`, `lock_timeout` and `idle_in_transaction_session_timeout`. Concurrency is capped and defaults low. The connection announces itself as `pgfathom` in `pg_stat_activity`. |
| Balança com um prato "claim" e outro "evidence" equilibrando | **No claim without evidence** | Every inferred relationship carries a verdict and the metric behind it. Sampled runs can never report a *confirmed* relationship — only `--full` can prove absence of orphans. |
| Grade onde algumas células ficam vazias e se marcam como "skipped" em vez de sumirem | **Silence is never a clean bill of health** | Tables skipped for missing privileges, candidates that timed out, schemas not covered — all of it appears in the coverage block on every run. A clean report means "I looked and it's clean", never "I couldn't look". |

**Faixa final da seção** — `go.mod` renderizado num card mono estreito, com o título `The person who has to approve this will open go.mod first. We intend that to be a short read.` e o contador `4 dependencies · no cgo`.

---

### 09 — Naming profiles

*Eyebrow:* `— Naming profiles` · *Headline:* **Most schema tools assume English. The databases that need this tool most often aren't.**

**Coluna esquerda — TOML animado** (digitação char-a-char, syntax highlight mono com chaves em `--paper`, strings em `--confirmed`, comentários em `--paper-45`):

```toml
name = "pt-br"

column_suffixes = ["_id", "_codigo", "_cod", "_key", "_ref", "_fk"]
table_prefixes  = ["tb_", "tbl_", "sys_", "cad_", "mov_"]

[[plural]]                     # opcoes → opcao
suffix = "oes"
singular = "ao"

[[plural]]                     # animais → animal
suffix = "ais"
singular = "al"
```

**Coluna direita — normalização visual interativa.** Uma coluna `opcoes` entra e o componente mostra o *conjunto* de formas candidatas se abrindo em leque (`opcoes → {opcao, opcoe, opcoes}`), com a forma que casou acendendo em `--accent`. Texto:
> Normalization returns a *set* of candidate forms rather than one, so ambiguous plurals (`logins` → `logim`? `login`?) cost nothing in recall — every form is tried, and the one that matched is reported.

**Bloco de auto-detecção** (card glassed com badge `new` em `--accent`):
> **The convention is not the language's. It is the schema's.** pgfathom derives the convention from your own catalog: table prefixes by frequency, reference affixes from the foreign keys you already declared. What it detected is always printed — a profile that changes itself without saying so is worse than a wrong profile.

**Tabela de medição real** (de `docs/ROADMAP.md`, fase 3.5) — com números subindo por contador animado e rótulo de honestidade:

| Database | Tables | FKs | Recall (before) | Recall (with affix) |
|---|---|---|---|---|
| geon_pr_assai | 784 | 985 | 0.5% | 79.0% |
| tributech_2 | 148 | 114 | ~0% | 75.4% |
| sinter (Django) | 18 | 10 | 0.0% | — |

Nota mono abaixo, `--paper-45`: `measured during phase 3 against real public-sector schemas. not the benchmark corpus — those numbers get published when the corpus runs.`

**Chamada de contribuição** ao final da seção, card com borda tracejada:
> Adding a profile for your language is the easiest possible contribution. It needs no knowledge of the rest of the codebase — a TOML file and a table of test cases. → `Contribute a profile`

---

### 10 — Bento de capacidades

Grid bento 3×2 no formato exato do pointer (`bento-section.tsx:9-24`): card glassed, título em `--paper` seguido de `<br/>` e descrição em `--paper-70`, ilustração de 288px de altura ocupando a base do card.

| Título | Descrição | Ilustração animada |
|---|---|---|
| **Two commands, no configuration.** | `audit` for what the catalog already knows, `discover` for what it doesn't. | Duas prompts de terminal alternando com digitação |
| **Findings that need no inference.** | `NOT VALID` constraints never validated, and declared FKs with no index on the child side. | Constraint com selo "NOT VALID" pulsando em `--weak` |
| **Containment in two dimensions.** | By row and by distinct value — a single bad value repeated a million times and a million rare bad values are different problems. | Dois medidores preenchendo em ritmos diferentes |
| **Versioned JSON model.** | The integration point. Consume it and generate whatever you like from a schema that finally knows its own relationships. | JSON se desenhando com chaves colapsando/expandindo |
| **Reviewable `.sql` artifacts.** | Nothing generated is meant to be executed unreviewed. | Arquivo `.sql` com checkmark aparecendo linha a linha |
| **Honest coverage on every run.** | 312 tables · 298 analyzed · 14 skipped. Always printed, never omitted. | Barra de cobertura segmentada preenchendo com um segmento cinza explícito |

---

### 11 — Prior art (substitui a grade de depoimentos)

*Eyebrow:* `— Prior art` · *Headline:* **pgfathom is not new science, and says so.**

*Texto:*
> Containment is known in the data-profiling literature as an **inclusion dependency** — the automatically testable part of a foreign key. There is a mature body of algorithms for discovering them (SPIDER, BINDER, MIND), implemented in Metanome. Commercial GUI modelers such as Hackolade infer relationships from metadata. Azimutt flags `_id` columns without declared relations.

**Tabela comparativa** (linhas com hover que acende a borda; a coluna pgfathom com fundo `--accent-g` sutil):

| | pgfathom | Metanome | Hackolade | Azimutt |
|---|---|---|---|---|
| Native PostgreSQL CLI | ● | ○ | ○ | ○ |
| Validates against the data | ● | ● | ○ | ○ |
| Mines evidence from the catalog | ● | ○ | ○ | ○ |
| Non-English naming conventions | ● | ○ | ○ | ○ |
| Reports its own coverage | ● | ○ | ○ | ○ |
| Hands you reviewable DDL | ● | ○ | ○ | ○ |

**Faixa "não compete com"** — marquee de cápsulas outline com o motivo em `--paper-45` no hover:
`Squawk` *migration linting* · `Atlas` *drift* · `SchemaSpy` *diagrams* · `Azimutt` *diagrams* · `sqlc` *codegen* · `jOOQ` *codegen*

*Fecho:* Code generation is explicitly **not** on the roadmap.

---

### 12 — Roadmap (substitui pricing)

*Eyebrow:* `— Roadmap` · *Headline:* **Eight phases to v0.1.**
*Subtítulo:* The ordering principle is to push risk forward. The first phases touch no user data and need no network, so they fail cheap.

**Timeline vertical** com a linha em `--line` e os marcadores preenchendo em `--accent` conforme entram em viewport (`stroke-dashoffset`). Fases entregues com marcador cheio, planejadas com marcador oco tracejado.

| Fase | Capacidade | Status |
|---|---|---|
| 1 | Core model and naming profiles | **Shipped** |
| 2 | Catalog inspection · `pgfathom audit` | **Shipped** |
| 3 | Name-based candidate inference | **Shipped** |
| 3.5 | Naming auto-detection | **Shipped** |
| 4 | Planner-statistics prefilter | In progress |
| 5 | Data validation · `pgfathom discover` | Planned |
| 6 | Join mining from views and functions | Planned |
| 7 | Terminal, JSON and SQL output | Planned |
| 8 | Benchmark corpus and release | Planned |

**Resolvido:** a tabela reflete o estado atual do repositório (fases 1–3.5 entregues, fase 4 em progressão), e não o README, que ainda diz *"no code yet"*. Como o projeto continua em movimento, a seção abre com um selo mono ao lado do headline em vez de fingir estabilidade:

`in progress · v0.1 · status reflects the repository, not the README`

Os três estados de marcador da timeline: **Shipped** (marcador cheio `--accent`), **In progress** (marcador com anel pulsando lentamente, 3s), **Planned** (marcador oco tracejado `--paper-45`).

**Bloco "after v0.1"** em card glassed separado:
> `pgfathom check --baseline` for CI — fail the build when a new undeclared relationship appears or an orphan count grows. Then structural findings, cross-cutting patterns (tenant columns, polymorphic pairs), and DBML/Mermaid/PlantUML export.

**Bloco "how correctness is measured"**, com o número em display serif gigante e o texto ao lado:
> **Zero.** The metric that has no tolerance. A missed relationship costs you a finding. A wrong one confirmed costs you the tool.
>
> The headline metric is recovery rate: take a schema with complete foreign keys, drop every one of them, run pgfathom, and count how many come back — against a public corpus (GitLab, Odoo, Discourse, Redmine, Mastodon) so anyone can reproduce it.

---

### 13 — Install & build

*Eyebrow:* `— Get started` · *Headline:* **Requires Go 1.25 or newer. No cgo, no other toolchain.**

**Tabs de código** (`Install` / `From source` / `Run` / `Contribute`) no padrão de `developers-section.tsx`, com botão de copiar e digitação por linha:

```console
# Install
$ go install github.com/lvcas-dotcom/pgfathom/cmd/pgfathom@latest

# From source
$ git clone https://github.com/lvcas-dotcom/pgfathom
$ cd pgfathom
$ make build          # → bin/pgfathom
$ make test           # unit suite: no Docker, no network

# Run
$ pgfathom audit    --dsn "$DATABASE_URL"
$ pgfathom discover --schema public --min-score 0.6
$ pgfathom discover --schema public --full --format json > model.json
```

**Quatro selos ao lado**, no padrão do `developers-section.tsx`:
`Go native` *single static binary* · `Zero config` *sensible defaults, detected profile* · `No container needed` *make test needs no Docker* · `4 dependencies` *short go.mod, on purpose*

---

### 14 — Contribute

*Eyebrow:* `— Contributing` · *Headline:* **The design is the cheapest thing to change right now.**

*Texto:*
> If you have run into this problem on a real legacy database, open an issue: what the schema looked like, what naming convention it used, and what a tool would have needed to find. Once implementation lands, the two most valuable contributions will be **naming profiles for other languages** and **real-world schemas for the benchmark corpus**.

Três cards com os ícones do sprite:
- `documentation-icon` → **Read the design doc** — `docs/PGFATHOM.md`
- `github-icon` → **Open an issue** — describe your schema
- `social-icon` → **Add a naming profile** — a TOML file and a table of test cases

Os três apontam para o repositório (docs, issues, diretório de profiles). Não há link para rede social — ver seção 16.

---

### 15 — CTA final

Bloco de largura total com borda 1px `--accent`, **spotlight radial seguindo o mouse** (`cta-section.tsx:44-50`, adaptado para `rgba(224,71,61,0.18)`), e a **esfera ASCII** reaparecendo à direita em escala menor, fechando o arco visual com o hero.

*Headline* (display, `clamp(2.5rem, 6vw, 4.5rem)`):
> Your schema already knows.
> It just never told you.

*Texto:* Point it at the database nobody wants to touch. It reads, it counts, it writes nothing.

*CTAs:* `Get started` (primário `--accent`) · `Star on GitHub` (outline, ícone do sprite).

---

### 16 — Footer

Fundo `--ink-900`. **Onda ASCII animada** no topo do footer (porte de `animated-wave.tsx` do `asthetic`), em `--accent` a baixa opacidade — a superfície da água sob a qual o mergulhador da logo trabalha.

- Coluna de marca: logo inline grande + `Sound the depth of a legacy PostgreSQL schema.` + selos `Apache-2.0` · `Go 1.25+` · `PostgreSQL 13+`.
- **Product:** Problem · What it does · How it works · Safety · Naming profiles
- **Docs:** Design doc · Roadmap · CLI reference · JSON model
- **Project:** GitHub · Issues · Contributing · License · Prior art
- Barra inferior: `© 2026 pgfathom · Apache-2.0`, o **switch de bandeira** e um único ícone `github-icon` do sprite, com hover `translateY(-2px)` + cor `--accent`.

**Resolvido:** os símbolos `x-icon`, `bluesky-icon`, `discord-icon` e `social-icon` do sprite **não são usados no footer** — o projeto não tem essas contas, e ícone social apontando para `#` é exatamente o tipo de fachada vazia que o resto da página se compromete a não ter. Os símbolos permanecem no sprite, disponíveis se as contas passarem a existir. O `social-icon` continua em uso na seção 14, onde é ilustrativo e não é link para rede.

---

## Sistema de animação (transversal)

| Padrão | Implementação | Onde |
|---|---|---|
| Entrada de seção | `IntersectionObserver` (`threshold: .1`, `once`) + `opacity/translateY(20px)/scale(.98)`, 800ms `cubic-bezier(.33,1,.68,1)` | Todas as seções |
| Entrada em cascata | `transitionDelay: index * 100ms` | Cards, steps, links do mobile menu |
| Char-in com blur | `@keyframes char-in` (opacity 0→1, blur 40px→0, translateY 100%→0) | Headline do hero, veredito rotativo |
| Digitação de código | `code-line-reveal` + `code-char-reveal` (15ms/char, 80ms/linha) | Terminais, SQL, TOML |
| Contador | `requestAnimationFrame` com easing `1-(1-p)³`, 2000ms | Recall, órfãos, tabelas |
| Marquee | `translateX(0 → -50%)` linear infinito, 30s / 25s reverso | Trust bar, "não compete com" |
| Hover lift | `translateY(-2px)` + sombra, `cubic-bezier(.34,1.56,.64,1)` | Cards, botões pequenos |
| Sublinhado de nav | `width: 0 → 100%` a partir da esquerda, 300ms | Links da navbar |
| SVG de feature | SMIL (`<animate>`, `<animateMotion>`) — sem JS, sem lib | Safety, veredictos, bento |
| Esfera / onda | `<canvas>` 2D + `requestAnimationFrame`, com `devicePixelRatio` | Hero, CTA, footer |

**Acessibilidade e performance.** `prefers-reduced-motion: reduce` desliga marquees, digitação, esfera e onda (canvas para no primeiro frame, estado final visível). Canvas pausa via `IntersectionObserver` quando fora da viewport e no `visibilitychange`. Foco visível em `--accent` com offset. Contraste: `--paper` sobre `--ink-700` ≈ 13:1; `--paper-45` reservado a texto ≥12px não essencial; `--accent` sobre `--ink-700` ≈ 4.8:1 (usado em texto ≥16px ou bold).

## Arquitetura de arquivos

```
src/
  App.tsx                        composição das seções
  index.css                      reset + tokens + @theme do Tailwind + utilities
  components/
    layout/  Navbar.tsx  Footer.tsx  Section.tsx  Reveal.tsx
    canvas/  AsciiSphere.tsx  AsciiWave.tsx
    ui/      Button.tsx  GlassCard.tsx  Terminal.tsx  CodeBlock.tsx
             CopyCommand.tsx  Counter.tsx  Marquee.tsx  Icon.tsx  Logo.tsx
             LanguageSwitch.tsx  flags/FlagBR.tsx  flags/FlagUS.tsx
    sections/  Hero  TerminalPreview  TrustBar  Problem  Verdicts  HiddenRelation
               Pipeline  Safety  NamingProfiles  Bento  PriorArt  Roadmap
               Install  Contribute  FinalCta
    illustrations/  (SVGs animados autorais, um arquivo por ilustração)
  i18n/
    types.ts                     Lang, Copy (forma derivada de copy.en)
    copy.en.ts                   texto completo em inglês
    copy.pt.ts                   texto completo em pt-BR
    LanguageContext.tsx          provider, useLang(), useCopy()
  content/
    terminal-output.ts           saída de terminal, SQL, TOML — não traduzidos
  hooks/
    useReveal.ts  useCounter.ts  useScrolled.ts  useReducedMotion.ts
public/
  icons.svg                      sprite normalizado para currentColor
```

Todo o texto sai centralizado em `src/i18n/copy.*.ts` — a landing de um projeto em progressão vai mudar de texto muitas vezes, e não deve exigir mexer em JSX para isso. Os artefatos literais da ferramenta ficam separados em `content/terminal-output.ts`, justamente porque não têm versão traduzida.

## O que não muda

- `pgfathom-front/pgfathom/pgfathom/public/icons.svg` mantém os mesmos 6 `symbol id`; só os `fill`/`stroke` são normalizados. Nenhum símbolo é removido, mesmo os que saem do footer.
- `src/assets/icons/pgfathom-logo.svg` e `pgfathom.png` permanecem; a versão inline é derivada, não substitui.
- `tsconfig.*` e `.oxlintrc.json` — inalterados.
- O repositório `C:\Users\user\Documents\Projects\pgfathom` (o Go) **não é tocado**. É fonte de conteúdo, apenas.

## Decisões resolvidas

| # | Questão | Resolução |
|---|---|---|
| 1 | Stack | Tailwind v4 — instalado e verificado (`npm run build` verde) |
| 2 | Status do roadmap | Estado atual do repositório, com selo `in progress` explícito |
| 3 | Pricing / depoimentos | Removidos; slots ocupados por Roadmap, Prior Art e Contribute |
| 4 | Redes sociais | Removidas do footer; sobra só GitHub. Símbolos ficam no sprite |
| 5 | Idioma | Bilíngue EN/PT-BR com switch de bandeira, EN padrão, zero dependência |

---

# Implementação — o que foi construído

## Terceira referência: `Pgfathom landing page redesign`

Entregue depois do planejamento acima, e adotada como **fundação** em vez de mais uma inspiração. `pgfathom Landing v3.dc.html` é um design finalizado, construído a partir do README real, já bilíngue EN/PT, com paleta fechada e 1.100 linhas de markup. Isso muda a natureza do trabalho: em vez de inventar copy e paleta, o trabalho vira portar o v3 para React e sobrepor o que o usuário pediu das outras duas referências.

**Adotado do v3, sem alteração:**

| | |
|---|---|
| Paleta | `#1b1512` base · `#f6f1ea` texto · `#e2483d` acento · `#3a2e27` hairline · `#a79a90` muted · `#221a16` / `#1e1714` surfaces · `#2b211c` grid · `#57453b` borda hover · `#ef5a4f` acento hover |
| Tipografia | Geist + Geist Mono (Google Fonts) — não Instrument Serif como no plano original |
| Copy EN/PT | Os dois dicionários completos, portados verbatim |
| Estrutura | hero → problem → verdicts → safety → pipeline → profiles → ci → prior-art → roadmap → cta → footer |
| Grid | `max-w-[1140px]`, padding `28px`, seções de `112px` |
| Animações | `fade-up` na entrada, digitação do comando (42ms/char), reveal por `IntersectionObserver` a `threshold .1` / `rootMargin -40px` |
| Chrome | Barras de terminal com 3 pontos, blocos `psql`/TOML/YAML, badges tracejados |

**Divergências deliberadas do v3:**

1. **Aviso de status.** O v3 diz *"Pre-release. No code yet."* e marca as 8 fases como `Specified`/`Planned`. Trocado pelo estado real do repositório: fases 1–3.5 `Shipped`, fase 4 `In progress`, 5–8 `Planned`, e o aviso reescrito para *"Pre-release. Still in progress."* explicando que a validação de dados é a fase 5 e por isso a saída de terminal continua sendo design alvo.
2. **Switch de idioma.** Os botões de texto `EN`/`PT` do v3 viraram o controle de bandeiras com ring deslizante.
3. **Navbar.** O `position: sticky` simples do v3 virou a navbar glassed flutuante que encolhe no scroll, com barra de progresso.
4. **Esfera e ícones animados.** Não existiam no v3 — vieram do `asthetic`, conforme pedido.

## Seções finais

| # | Componente | O que ganhou além do v3 |
|---|---|---|
| — | `Navbar` | Glassed sempre visível, 72px→56px, pill `top-4` no scroll, progresso de scroll 1px, switch de bandeira, overlay mobile em cascata |
| 01 | `Hero` | **Esfera ASCII** 720px à direita, grid com máscara radial, terminal com digitação e faixa de 4 fatos |
| 02 | `Problem` | Rodapé no bloco `psql` apontando a ausência do `Foreign-key constraints:` |
| 03 | `Verdicts` | **3 ilustrações SMIL** — órfãos caindo do vínculo, arco que se desenha e trava, feixe que oscila entre alvos |
| 04 | `Safety` | **5 ícones SMIL** — escrita ricocheteando no escudo, valores dissolvendo na memória, ponteiro contido pelo timeout, balança claim/evidence, grade com célula pulada marcada |
| 05 | `Pipeline` | — |
| 06 | `Profiles` | — |
| 07 | `Ci` | — |
| 08 | `PriorArt` | Hover que acende a linha |
| 09 | `Roadmap` | Selo `in progress`, marcadores de 3 estados (cheio / anel pulsando / tracejado) |
| 10 | `FinalCta` | **Esfera** menor fechando o arco com o hero + spotlight radial seguindo o mouse |
| — | `Footer` | Só GitHub, mais o switch de bandeira |

## Arquivos

```
src/
  index.css                       @theme com a paleta, keyframes, utilities (.reveal .glass .grid-veil .noise)
  main.tsx                        LanguageProvider no topo
  App.tsx                         composição + overlay de ruído global
  i18n/     types.ts · copy.en.ts · copy.pt.ts · LanguageContext.tsx
  content/  data.ts               terminal, SQL, TOML, URLs — não traduzidos
  hooks/    useReveal · useTypewriter · useScrolled (+ useScrollProgress)
  components/
    layout/         Navbar · Footer
    canvas/         AsciiSphere
    ui/             Reveal · Eyebrow · Section (+H2, Lead) · LanguageSwitch · icons · flags/FlagBR · flags/FlagUS
    illustrations/  VerdictIcons · SafetyIcons
    sections/       Hero · Problem · Verdicts · Safety · Pipeline · Profiles · Ci · PriorArt · Roadmap · FinalCta
public/favicon.svg                marca do pgfathom em #e2483d (substitui a do Vite)
```

## Verificação

- `npm run build` — verde, 49 módulos, 1.55s · CSS 27.3 kB (6.35 kB gzip) · JS 286.7 kB (90.3 kB gzip)
- `npm run lint` — sem erros; 4 avisos `only-export-components`, todos sobre granularidade de HMR em arquivos que exportam componente + hook/constante. Tradeoff aceito.
- Utilities conferidas no CSS gerado, incluindo as fracionárias (`px-6.5`, `py-5.5`, `pb-7.5`), as arbitrárias (`text-[13.5px]`, `opacity-[0.022]`, `text-bone/[0.88]`) e as responsivas (`md:py-28`, `lg:grid-cols-[0.85fr_1.15fr]`, `lg:sticky`).
- **Não verificado visualmente em navegador** — a extensão do Chrome não foi conectada nesta sessão. `npm run dev` sobe em `localhost:5173`.

## Já executado

- `npm i -D tailwindcss @tailwindcss/vite` — 12 pacotes, 0 vulnerabilidades
- `vite.config.ts:4,11` — plugin `tailwindcss()` registrado
- `index.html` — Geist/Geist Mono, `theme-color`, meta description
- `src/App.css` removido (template do Vite); `src/index.css` reescrito
- `assets/pgfathom-logo.png` copiado do redesign para `src/assets/icons/` (o `pgfathom.png` que já existia é o wordmark)
