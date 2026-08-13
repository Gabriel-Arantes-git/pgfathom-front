# Plano de ação — Terminal CLI interativo + refresh de conteúdo

Data: 2026-08-13
Comando: `/designRefactor`
Agente: Designer UX/UI

## Contexto

O front (`pgfathom-front`) foi escrito quando o CLI ainda estava incompleto. O
código em `/home/gabriel-arantes/Documentos/Projects/pgfathom` evoluiu:

- Roadmap concluído até a fase 8 (join mining, saída terminal/JSON/SQL, chaves
  compostas, corpus de benchmark, release v0.1.x). Fonte: `README.md`.
- Formato do `discover` mudou para colunas
  `relation | rows | values | orphan rows | orphan values | examined | method`
  com blocos BROKEN/CONFIRMED/WEAK/UNVALIDATED/DETECTED e bloco de cobertura.
  Fonte canônica: `internal/report/testdata/discover_all_verdicts.golden`.
- Comandos reais: `discover`, `audit`, `version`, `setup`, e a tela raiz
  (`pgfathom` sem args) que imprime o banner ASCII + descrição + lista de
  comandos. Não existe comando `status`.
- Cores de marca do CLI (`internal/report/style.go`): BrandRed `#e2483d`
  (= `--color-accent`), BrandAqua `#7ec9b8` (CONFIRMED, ausente no front),
  BrandStone `#8c7d73` (≈ muted), BrandBone `#ddd5d0` (≈ bone).
- Banner ASCII (`internal/cli/banner.go`): glifo de 14 linhas com gradiente de
  vermelhos (`descent[]`, `#e2483d`→`#a02b25`).

Decisões do usuário (perguntas respondidas):
1. Escopo = "Também refrescar conteúdo" (terminal interativo + corrigir estado
   desatualizado + benchmark com números reais do corpus, remover badge de
   placeholder).
2. Comandos no autocomplete = `discover` + tela raiz (como "status") + `audit`
   + `version`. "status" = tela raiz: banner ASCII + descrição + comandos +
   linha de estado/versão.

## Objetivo

Na seção "pgfathom: target output" (o card de terminal no Hero), transformar o
terminal estático em um terminal interativo: barra de escrita clicável com
autocomplete, onde o usuário roda um dos 4 comandos reais e a saída troca,
usando o estilo (cores, banner) do CLI que fizemos. Refrescar o conteúdo
desatualizado do restante do site.

## Arquivos

### Novos
1. `src/content/cli.ts`
   - Glifo ASCII (`mark`) e cores `descent[]` do banner.
   - Registro de comandos: `{ input, short, buildScreen(t) }`.
   - Builders de tela por comando, compondo texto literal do CLI (não traduzido)
     + identificadores localizados reaproveitando `t.graph.*` (zero mudança em
     `types.ts`):
     - `discover`: novo formato em colunas + cobertura + prefiltro + DETECTED.
     - raiz/`pgfathom`: banner + descrição longa + "Available Commands" +
       linha de estado (v0.1.x · pre-release · pipeline completo).
     - `audit`: grupos UNINDEXED + HOT COLUMN + bloco de cobertura.
     - `version`: linhas version/commit/built.
   - Saída tipada como segmentos com tom (`broken|confirmed|weak|muted|bone|
     accent`) para colorização fiel.

2. `src/components/canvas/CliTerminal.tsx`
   - Chrome do terminal (3 pontos, caption `t.terminalCaption`, badge read-only)
     — reaproveita o visual atual.
   - Área de saída: renderiza a tela do comando ativo; efeito typewriter na
     linha de comando (reusa `useTypewriter`) + fade-in do corpo (`done`).
   - Barra de escrita no rodapé: prompt `$` + input editável, foco por clique.
   - Autocomplete: dropdown com comandos que casam + `short`; chips sugeridos
     (help/status/discover/audit) para descoberta. Tab completa, Enter roda,
     ↑/↓ navega, Esc fecha. Roles `combobox`/`listbox`/`option`, `aria-*`.
   - Respeita `prefers-reduced-motion` (via `useTypewriter`).
   - Estado inicial = `discover --schema public` (preserva a primeira impressão).

### Alterados
3. `src/index.css` — adicionar token `--color-confirm: #7ec9b8` (CONFIRMED).
4. `src/components/sections/Hero.tsx` — trocar o bloco `GlowGrid` do terminal
   estático por `<CliTerminal />`. Mantém eyebrow, título, lead, botão de copiar
   o comando, CTA e o aviso. Remove imports que migram para o CliTerminal.
5. `src/content/data.ts` — `BENCHMARK` com números reais do corpus (regime
   partial): gitlab 1054/1857 (62.1→62.2), municipal pt-BR 226/277 (3.6→84.2),
   discourse 354/23 (50.0→50.0). `BENCHMARK_IS_PLACEHOLDER = false`.
6. `src/i18n/copy.en.ts` / `copy.pt.ts`:
   - `noticeBody`: refletir estado atual (pipeline completo incl. join mining e
     saídas terminal/JSON/SQL na v0.1; benchmarks medidos contra corpus público
     e reproduzíveis com `make benchmark`; terminal abaixo é dirigível). Mantém
     honestidade de pré-release.
   - `benchBadge`: "public corpus · make benchmark" (deixa de ser placeholder).
   - `benchLead`, `benchLegend` (→ "Profile alone" / "+ schema detection"),
     `metricBody`, `benchAggregate` alinhados ao corpus real.
   - `ciBadge`: mantém "planned · after v0.1" (ainda verdadeiro).

## Fora de escopo (não tocar)
- Lógica de `discover`/dados do grafo do Hero (`SchemaGraph`, `GRAPH_*`).
- Seções Problem, Verdicts, Safety, Pipeline, Profiles, PriorArt, Contributors,
  FinalCta (copy atual agrada; sem mudança de estado necessária).
- CI: continua "planejado" — é honesto.

## Validação — CONCLUÍDA (2026-08-13)
- `npm run build` (tsc + vite): limpo. `npm run lint` (oxlint): limpo, sem
  warnings novos.
- Verificação end-to-end via Playwright contra o dev server real (não apenas
  leitura de código): estado inicial (discover), abertura do autocomplete,
  clique em sugestão (audit), digitação+Enter (help/status → banner, version,
  comando inválido → erro estilo cobra), navegação por teclado (↓ + Enter),
  reabertura do dropdown após rodar um comando, locale PT, viewport mobile.
  Zero erros de console em todas as passagens.
- Bug real encontrado e corrigido durante a verificação: depois de rodar um
  comando com Enter, o input permanecia focado e um clique subsequente não
  reabria o dropdown (nenhum evento `focus` novo disparava). Corrigido
  adicionando `onClick` no input e no wrapper, além do `onFocus` existente.
- Bug latente corrigido em `useTypewriter`: o hook não resetava `typed`/`done`
  ao trocar o texto, o que faria o terminal (que agora reusa o hook entre
  comandos diferentes) pular a animação de digitação ao trocar de comando.
- Fidelidade conferida linha a linha contra `internal/report/terminal_test.go`,
  os golden files, `internal/cli/banner.go`, `internal/cli/root.go` e o README.

## Ajuste pós-entrega (2026-08-13, mesmo dia) — tamanho fixo + streaming

Feedback do usuário: o card redimensionava (encolhia/crescia) conforme o
comando trocava, e isso lag ava a página ao rolar por perto do terminal —
a troca de altura do card empurrava o resto da página em pleno scroll.
Também pediu que o conteúdo "apareça sendo escrito em tempo real" em vez de
só o bloco inteiro dar fade-in de uma vez.

- `src/hooks/useLineReveal.ts` (novo): revela `screen.body` linha a linha
  (16ms/linha, ~450ms para os 30 linhas do discover), gated no `done` do
  typewriter do prompt. Reduced-motion revela tudo de uma vez, mesmo padrão
  do `useTypewriter`.
- `CliTerminal.tsx`: o painel de saída agora tem altura fixa responsiva
  (`h-[280px] sm:h-[360px] lg:h-[440px]`) com scroll interno próprio
  (`.cli-scroll`, thin scrollbar na cor do tema) — trocar de comando nunca
  mais redimensiona o card. Auto-scroll para o topo ao iniciar um novo
  comando, e auto-scroll acompanhando as linhas conforme entram (como um
  terminal real). Cada linha ganha uma pequena animação de entrada
  (`cli-line-in`, 180ms).
- Verificado via Playwright: `boundingBox()` do card idêntico (1084×621)
  através de discover → version → help → discover de novo, inclusive
  durante o streaming. Confirmado empiricamente (polling do DOM) que as
  linhas realmente entram progressivamente (~15ms/linha), não tudo de uma
  vez. Build/lint limpos, zero erros de console.

## Ajuste pós-entrega #2 (2026-08-13) — cursor na barra de input + SchemaGraph mais leve

Usuário confirmou que a sensação de lag ao rolar a página não vinha do
terminal, e apontou a causa provável: a animação do `SchemaGraph` (as
"tabelas" do hero). Pediu também o cursor piscante (já usado no terminal)
também na barra de input do autocomplete.

- **Cursor na barra de input** (`CliTerminal.tsx`): input envolto num wrapper
  `relative`; `caret-transparent` esconde o cursor nativo do navegador; um
  `<span className="caret">▌</span>` absolutamente posicionado em
  `left: {query.length}ch` (seguro porque a fonte é monoespaçada) reproduz o
  mesmo glifo/piscar do prompt de saída. Confirmado via Playwright: opacidade
  realmente pisca (`["1","1","0","0","0","1",...]`) e a posição acompanha o
  texto digitado linearmente (0px → 8.6px → 43px → 267px para 0/1/5/31
  caracteres).
- **`SchemaGraph.tsx` mais leve**: o `requestAnimationFrame` recomputava a
  geometria completa de cada uma das 5 edges — incluindo reescrever o atributo
  `d` do `<path>` SVG (um reparse de geometria, a escrita mais cara do loop) —
  a cada frame (60×/s), mesmo com as tabelas paradas (geometria 100% estática
  entre resizes). Extraído para `computeGeometry()`, chamado só quando a
  geometria pode ter mudado de fato (montagem, fonte web carregada, resize).
  O loop por frame agora só toca o que realmente anima:
  `strokeDashoffset`/`opacity` do path, `r`/`opacity` do pulso do lock,
  `cy`/`opacity` dos pontos de órfão (a posição base deles também virou
  estática, cacheada). Comportamento visual idêntico — é puro corte de
  trabalho redundante, sem mudar timing nem easing.
- Verificado via Playwright: `MutationObserver` no atributo `d` de todos os
  paths confirma **zero reescritas em 3s** de animação rodando em regime
  permanente (antes: 60/s × 5 edges). Nenhuma *long task* detectada durante
  scroll simulado no mesmo teste. Build/lint limpos, zero regressão na suíte
  de interação (autocomplete/click/teclado) já existente.

## Auditoria de desempenho #3 (2026-08-13) — navbar, GlowGrid, scroll hooks

Usuário confirmou que a sensação de travamento não vinha do `SchemaGraph` e
pediu uma auditoria completa do código por problemas de desempenho, além do
cursor piscante também na barra de input (já coberto acima). Achados,
por ordem de impacto, e correções aplicadas:

1. **`Navbar` com `backdrop-filter: blur(14px) saturate(130%)` num header
   `fixed`, sempre montado.** Recomposto a cada frame em que o conteúdo atrás
   muda — todo frame de scroll, por definição. Fix: `src/hooks/useScrolled.ts`
   ganhou `useIsScrolling(idleMs=150)` (rAF-batched, detecta gesto de scroll
   ativo); `Navbar.tsx` troca a classe `glass` → `glass-scrolling` durante o
   gesto (`index.css`: `.glass-scrolling` é auto-contida, sem
   `backdrop-filter`, fundo mais opaco pra compensar) e volta a `glass` ~150ms
   após o scroll parar. Verificado via Playwright: `glass` em repouso,
   `glass-scrolling` em pleno gesto, `glass` de novo após o settle.
2. **Blur duplicado**: `LanguageSwitch` tinha seu próprio `backdrop-filter`
   aninhado dentro do navbar já borrado. Trocado por `bg-ink-600/85` sólido —
   sem diferença visual perceptível (pill pequeno sobre fundo já escurecido),
   sem o segundo passe de blur.
3. **6 instâncias de `GlowGrid`** (Hero×2, Verdicts, Benchmark, Contributors,
   Safety), cada uma com seu próprio listener `mousemove` no `document`,
   sempre ativo, fazendo `getBoundingClientRect`/`querySelectorAll` a cada
   evento — inclusive para seções fora da tela — sem nenhum throttling.
   `GlowGrid.tsx`: adicionado `IntersectionObserver` (rootMargin 200px) que
   mantém uma flag `nearViewport` barata, e o handler agora só agenda
   trabalho (via `requestAnimationFrame`, um por frame, não um por evento
   bruto) quando a grid está perto da tela; a limpeza ao sair também só roda
   uma vez (na transição dentro→fora), não a cada mousemove subsequente fora.
   Verificado: 50 `mousemove` síncronos despachados em um único frame
   resultam em 16 leituras de `getBoundingClientRect` no total (antes seria
   proporcional a 50 × 6 grids), e o mecanismo do glow continua correto
   (`--glow-intensity` vai a 0.8 no hover do centro do card e volta a 0 ao
   sair).
4. **`useScrolled`/`useScrollProgress` sem throttling**, disparando `setState`
   (e portanto re-render do `Navbar`, que está sob o `backdrop-filter`) a
   cada evento `scroll` bruto. Ambos agora batched via `requestAnimationFrame`
   — no máximo um commit por frame renderizado, independente da taxa do
   dispositivo de entrada.

Descartados após revisão: carrossel do `Contributors` (CSS puro +
`animationend`), `useContributors` (fetch único, cacheado em
`sessionStorage`), overlay de ruído (raster estático, decodificado uma vez).

Build/lint limpos. Regressão completa (autocomplete, teclado, tamanho fixo do
terminal, PT/mobile) re-executada sem erros de console.
