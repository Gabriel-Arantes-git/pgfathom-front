/**
 * Literal artifacts of the tool: terminal output, SQL, TOML, CLI flags.
 * Deliberately NOT translated — rendering "QUEBRADO" where the CLI prints
 * "BROKEN" would show output the tool does not produce, and the whole page
 * rests on not misrepresenting it.
 */

export const REPO = 'https://github.com/lvcas-dotcom/pgfathom'
export const ROADMAP_DOC = `${REPO}/blob/main/docs/ROADMAP.md`
export const LICENSE_DOC = `${REPO}/blob/main/LICENSE`
export const ISSUES = `${REPO}/issues`

/**
 * The repository has no single bilingual design document: `docs/PGFATHOM.md`
 * and `docs/ROADMAP.md` are written in Portuguese, while `README.md` is the
 * English description of the same design. So each language points at the
 * artifact actually written in it. If an English PGFATHOM.md ever lands, this
 * is the only place that has to change.
 */
export const DESIGN_DOC: Record<'en' | 'pt', string> = {
  en: `${REPO}#readme`,
  pt: `${REPO}/blob/main/docs/PGFATHOM.md`,
}

export const COMMAND = 'pgfathom discover --schema public'

export type ReportRow =
  | { kind: 'head'; text: string; gap?: boolean }
  | {
      kind: 'row'
      child: string
      parent: string
      ratio: string
      note: string
      tone: 'broken' | 'confirmed' | 'weak'
    }

export const REPORT: ReportRow[] = [
  { kind: 'head', text: 'BROKEN — the relationship is real, the integrity is not' },
  {
    kind: 'row',
    child: 'os_servico.resp_tecnico',
    parent: '→ funcionario.id',
    ratio: '99.7%',
    note: '1,284 orphans',
    tone: 'broken',
  },
  {
    kind: 'row',
    child: 'pedido.cliente_id',
    parent: '→ cliente.id',
    ratio: '99.9%',
    note: '37 orphans',
    tone: 'broken',
  },
  {
    kind: 'row',
    child: 'lancamento.conta_id',
    parent: '→ conta.id',
    ratio: '94.2%',
    note: '11,903 orphans',
    tone: 'broken',
  },
  { kind: 'head', text: 'CONFIRMED — undeclared foreign key, no orphans found', gap: true },
  {
    kind: 'row',
    child: 'item_pedido.pedido_id',
    parent: '→ pedido.id',
    ratio: '100.0%',
    note: '0 orphans',
    tone: 'confirmed',
  },
  {
    kind: 'row',
    child: 'endereco.municipio_id',
    parent: '→ municipio.id',
    ratio: '100.0%',
    note: '0 orphans',
    tone: 'confirmed',
  },
  { kind: 'head', text: 'WEAK — insufficient evidence to conclude', gap: true },
  {
    kind: 'row',
    child: 'documento.entidade_id',
    parent: '→ entidade.id',
    ratio: '41.0%',
    note: 'polymorphic pair',
    tone: 'weak',
  },
]

export const REPORT_PROFILE = 'profile pt-br (detected) · sampled validation (100k rows/table) · 312 tables'
export const REPORT_COVERAGE = [
  '312 tables · 298 analyzed · 14 skipped (no SELECT privilege)',
  '1,847 candidates · 226 validated · 3 timed out',
]
export const REPORT_WARNING =
  '! sampled run — nothing here is confirmed. orphan rows cluster on disk, which is exactly what page-level sampling is worst at finding. re-run with --full to prove absence.'

/**
 * The hero graph's schema. Deliberately the same fictional database the
 * terminal report below it prints, down to which relationships come back
 * broken — the animation is showing the run whose output you then read.
 *
 * `x`/`y` are normalized positions inside the graph box; `z` is depth, 0 at
 * the front and 1 at the back, and drives scale, blur, opacity and stacking.
 */
/**
 * The graph box, in CSS pixels. Fixed rather than measured so a card's
 * position is known at render time: the cards used to get their transform from
 * an effect that read `getBoundingClientRect`, which meant the very first paint
 * put all of them stacked at the box's top-left corner until the effect caught
 * up — and left them there for good if the box happened to measure zero.
 */
export const GRAPH_BOX = { width: 660, height: 520 } as const

/** Column names are localized — see `graph` in the copy files. Types are not. */
export type GraphColumn = { type: string; fk?: boolean }

export type GraphTable = {
  id: string
  columns: GraphColumn[]
  /** Normalized against GRAPH_BOX. Values outside 0–1 bleed past the edges. */
  x: number
  y: number
  z: number
  /**
   * Ambient scenery: never carries an edge, sits deep, and is allowed to run
   * off the box. These are what make the composition feel like a corner of a
   * much larger schema rather than a diagram that happens to be nine tables.
   */
  ambient?: boolean
}

/**
 * Layout, worked out on a grid rather than by eye — scattering these by hand
 * produced overlapping cards with text bleeding through, twice.
 *
 * The box is 660×520. A card is 152 wide and, measured rather than guessed,
 * 90 tall with three columns or 74 with two. The left 64px and bottom 56px are
 * covered by the fades, so the usable area is x ∈ [64, 660], y ∈ [0, 464].
 *
 * That gives a 3×3 grid: columns at x = 175 / 360 / 545 (185 apart, against a
 * 158px widest card) and rows at y = 60 / 195 / 330, with the middle column
 * pushed 30px down so it does not read as a spreadsheet. Every gap clears the
 * largest card by at least 26px.
 *
 * Slots are assigned so that **every edge joins two adjacent cells**. That is
 * what keeps a wire from passing under a card it has nothing to do with —
 * edges paint beneath the cards, so a crossing reads as the line diving into
 * the wrong table:
 *
 *     endereco    pedido       cliente
 *     municipio   item_pedido  conta
 *     funcionario os_servico   lancamento
 *
 * The grid is the skeleton, not the look. Each card is then nudged off its
 * slot by up to 13px, which is well inside the smallest gap, so the result
 * reads as a hand-placed composition while staying provably collision-free.
 * The `ambient` tables afterwards ignore all of this on purpose.
 */
export const GRAPH_TABLES: GraphTable[] = [
  // Column 2, row 1 — the headline table, nearest the viewer.
  {
    id: 'pedido',
    x: 0.56,
    y: 0.145,
    z: 0.05,
    columns: [{ type: 'bigint' }, { type: 'bigint', fk: true }, { type: 'date' }],
  },
  // Column 3, row 1 — right of pedido, so their edge is one short hop.
  {
    id: 'cliente',
    x: 0.814,
    y: 0.1,
    z: 0.18,
    columns: [{ type: 'bigint' }, { type: 'text' }, { type: 'bool' }],
  },
  // Column 2, row 2 — directly under pedido.
  {
    id: 'item_pedido',
    x: 0.529,
    y: 0.455,
    z: 0.12,
    columns: [{ type: 'bigint' }, { type: 'bigint', fk: true }, { type: 'int' }],
  },
  // Column 2, row 3.
  {
    id: 'os_servico',
    x: 0.563,
    y: 0.672,
    z: 0.32,
    columns: [{ type: 'bigint' }, { type: 'bigint', fk: true }, { type: 'date' }],
  },
  // Column 1, row 3 — left of os_servico.
  {
    id: 'funcionario',
    x: 0.251,
    y: 0.653,
    z: 0.4,
    columns: [{ type: 'bigint' }, { type: 'text' }],
  },
  // Column 1, row 1.
  {
    id: 'endereco',
    x: 0.247,
    y: 0.135,
    z: 0.48,
    columns: [{ type: 'bigint' }, { type: 'bigint', fk: true }, { type: 'text' }],
  },
  // Column 1, row 2 — directly under endereco.
  {
    id: 'municipio',
    x: 0.285,
    y: 0.357,
    z: 0.28,
    columns: [{ type: 'bigint' }, { type: 'text' }],
  },
  // Column 3, row 2.
  {
    id: 'conta',
    x: 0.84,
    y: 0.391,
    z: 0.44,
    columns: [{ type: 'bigint' }, { type: 'text' }],
  },
  // Column 3, row 3 — directly under conta.
  {
    id: 'lancamento',
    x: 0.816,
    y: 0.613,
    z: 0.52,
    columns: [{ type: 'bigint' }, { type: 'bigint', fk: true }, { type: 'numeric' }],
  },

  // Ambient scenery. Deliberately off the grid and half off the box: the
  // schema should look like it continues past the frame. Deep z keeps them
  // blurred and dim enough that they never compete with the wired cluster.
  {
    id: 'produto',
    x: 1.06,
    y: 0.21,
    z: 0.8,
    ambient: true,
    columns: [{ type: 'bigint' }, { type: 'text' }, { type: 'numeric' }],
  },
  {
    id: 'nota_fiscal',
    x: 1.02,
    y: 0.87,
    z: 0.88,
    ambient: true,
    columns: [{ type: 'bigint' }, { type: 'text' }],
  },
  {
    id: 'usuario',
    x: 0.63,
    y: -0.1,
    z: 0.84,
    ambient: true,
    columns: [{ type: 'bigint' }, { type: 'text' }, { type: 'timestamptz' }],
  },
  {
    id: 'filial',
    x: 1.09,
    y: 0.55,
    z: 0.92,
    ambient: true,
    columns: [{ type: 'bigint' }, { type: 'text' }],
  },
  {
    id: 'movimento',
    x: 0.3,
    y: 1.02,
    z: 0.78,
    ambient: true,
    columns: [{ type: 'bigint' }, { type: 'text' }, { type: 'date' }],
  },
]

export type GraphEdge = {
  from: string
  to: string
  /** `broken` edges shed orphan rows while they are connected. */
  verdict: 'broken' | 'confirmed'
  /** Seconds before this edge draws for the first time. */
  delay: number
  /** Seconds it stays fully connected before starting to dim. */
  hold: number
  /** Seconds spent fading out. */
  fade: number
  /** Seconds absent before it draws itself again. */
  pause: number
}

/**
 * Every edge runs on its own clock. The hold values are deliberately
 * non-harmonic — 7.4 / 9.3 / 6.2 / 10.7 / 8.1 — so the cycles never line up
 * and no two lines ever fade together. Each relationship reads as something
 * being proved on its own schedule rather than one looping animation.
 */
export const GRAPH_EDGES: GraphEdge[] = [
  { from: 'item_pedido', to: 'pedido', verdict: 'confirmed', delay: 0.7, hold: 7.4, fade: 1.5, pause: 0.6 },
  { from: 'pedido', to: 'cliente', verdict: 'broken', delay: 1.6, hold: 9.3, fade: 1.8, pause: 0.9 },
  { from: 'os_servico', to: 'funcionario', verdict: 'broken', delay: 2.4, hold: 6.2, fade: 1.3, pause: 1.3 },
  { from: 'endereco', to: 'municipio', verdict: 'confirmed', delay: 3.5, hold: 10.7, fade: 1.6, pause: 0.7 },
  { from: 'lancamento', to: 'conta', verdict: 'broken', delay: 4.3, hold: 8.1, fade: 1.9, pause: 1.1 },
]

export const PSQL_OUTPUT = `=# \\d pedido
      Column     |  Type  | Nullable
-----------------+--------+----------
 id              | bigint | not null
 cliente_id      | bigint | not null
 ...
Indexes:
    "pedido_pkey" PRIMARY KEY, btree (id)`

export const PIPELINE_STAGES = [
  'catalog',
  'usage evidence',
  'candidates',
  'scoring',
  'stats prefilter',
  'validation',
]

export const PROFILE_LANGS = ['pt-br', 'en', 'es']

export const CI_WORKFLOW_FILE = '.github/workflows/schema.yml'
export const PROFILE_FILE = 'profiles/pt-br.toml'

/**
 * PROTOTYPE ONLY — every number below is invented.
 *
 * This is the shape the recovery-rate table is specified to take once the tool
 * runs against the reference corpus; it is on the page so the layout exists
 * before the data does. The README is explicit that no numbers are claimed
 * until the corpus runs, so the section renders a loud placeholder badge and
 * this constant must be replaced wholesale, never edited into looking real.
 */
export const BENCHMARK_IS_PLACEHOLDER = true

export type BenchmarkRow = {
  schema: string
  tables: number
  keys: number
  /** % recovered by name matching alone */
  byName: number
  /** additional % recovered once join mining is enabled */
  byEvidence: number
}

export const BENCHMARK: BenchmarkRow[] = [
  { schema: 'gitlab', tables: 612, keys: 1487, byName: 71.4, byEvidence: 12.8 },
  { schema: 'odoo', tables: 894, keys: 2103, byName: 64.9, byEvidence: 17.2 },
  { schema: 'discourse', tables: 241, keys: 396, byName: 78.1, byEvidence: 9.4 },
  { schema: 'redmine', tables: 78, keys: 141, byName: 82.6, byEvidence: 6.1 },
  { schema: 'mastodon', tables: 156, keys: 288, byName: 69.3, byEvidence: 14.5 },
]
