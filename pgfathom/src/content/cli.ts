/**
 * The interactive terminal in the Hero runs four real pgfathom commands.
 * Every block below is copied, byte for byte where possible, from the actual
 * tool: the banner glyph and its gradient from `internal/cli/banner.go`, the
 * `discover` output from the reproducible demo in the README (itself a real
 * run against `docs/DEMO.md`, not a mockup — see CHANGELOG 0.1.2), the
 * `audit` output composed from the exact object names and suggestion strings
 * pinned by `internal/report/terminal_test.go`, and `version` from the real
 * `pgfathom version` format string in `internal/cli/root.go`.
 *
 * Nothing here is translated. Same rule as the rest of this file's siblings
 * in `data.ts`: rendering a CLI's own words in a language it does not speak
 * would show output the tool does not produce.
 */

const RULE = '─'.repeat(96)

/** The project's glyph, copied from `internal/cli/banner.go`. */
export const BANNER_MARK = [
  ' █████          ███          █████',
  '████████████ █████████ ████████████',
  '███████████ ███████████ ███████████',
  '███████████ ███████████ ███████████',
  '██████████  ███ ███ ███  ██████████',
  ' █████████  ███████████  █████████',
  '   ███████  ███████████   ██████',
  '            ███████████',
  '            ███████████',
  '              ███████',
  '               █████',
  '               █████',
  '               █████',
  '                ███',
]

/** The colour the glyph falls through, top to bottom — same array as `descent` in banner.go. */
export const BANNER_DESCENT = [
  '#e2483d', '#dd463b', '#d84439', '#d34137', '#ce3f36',
  '#c93d34', '#c43b32', '#bf3830', '#b9362e', '#b4342c',
  '#af322a', '#aa3029', '#a52d27', '#a02b25',
]

export type CliTone = 'bone' | 'dim' | 'bold' | 'accent' | 'confirm' | 'warn'

/**
 * Classifies a line of literal CLI output by the same rules
 * `internal/report/style.go` uses to pick Bold/Alert/Confirm/Warn/Dim: the
 * verdict and finding headings, the divider rules, "none" placeholders,
 * `!`-prefixed cautions, and the tally line. Doing it by pattern rather than
 * hand-tagging every line keeps each block a faithful, checkable copy of the
 * tool's own text instead of an annotated transcription.
 */
export function classifyLine(raw: string): CliTone {
  const s = raw.trim()
  if (s === '' || s === 'none') return s === '' ? 'bone' : 'dim'
  if (/^─+$/.test(s)) return 'dim'
  if (s.startsWith('!') || s.startsWith('Error:')) return 'warn'
  if (/^BROKEN\b/.test(s)) return 'accent'
  if (/^CONFIRMED\b/.test(s)) return 'confirm'
  if (
    /^(WEAK|UNVALIDATED|DETECTED|NOT ANALYZED|DISCARDED|UNINDEXED|NOT VALID|NO PRIMARY KEY|HOT COLUMN|DANGLING)\b/.test(
      s,
    )
  ) {
    return 'bold'
  }
  if (/^relation\s/.test(s)) return 'dim'
  if (/^\d+ broken · /.test(s)) return 'bold'
  if (/^(Usage:|Available Commands:|Flags:)$/.test(s)) return 'bold'
  return 'bone'
}

/**
 * `pgfathom discover --schema public --full`, reproduced exactly from the
 * README, which is itself a real run against the demo schema in
 * `docs/DEMO.md`: nine service orders name a technician absent from
 * `funcionario`, so `os_servico.resp_tecnico → funcionario.id` — a
 * relationship no name-matching heuristic finds, recovered instead from a
 * join predicate already sitting in a view — comes back BROKEN rather than
 * confirmed.
 */
export const DISCOVER_OUTPUT = `
  pgfathom v0.1.2 · PostgreSQL 16.14 · profile pt-br · threshold 0.50
  full validation — every row was examined; verdicts are conclusive

  BROKEN — the relationship is real; its integrity is not  (1)
  ${RULE}
  relation                                                rows   values  orphan rows  orphan values  examined  method
  public.os_servico.resp_tecnico → public.funcionario.id  97.5%  60.0%   200          200            8.0k      full

  CONFIRMED — total containment, verified row by row  (2)
  ${RULE}
  relation                                         rows    values  orphan rows  orphan values  examined  method
  public.item_pedido.pedido_id → public.pedido.id  100.0%  100.0%  0            0              60.0k     full
  public.pedido.cliente_id → public.cliente.id     100.0%  100.0%  0            0              20.0k     full

  WEAK — the data supports no conclusion either way  (0)
  ${RULE}
  none

  UNVALIDATED — no evidence gathered; not the same as clean  (0)
  ${RULE}
  none

  Nothing detected from 6 tables and 1 declared key; the pt-br profile applies alone.

  1 broken · 2 confirmed · 0 weak · 0 unvalidated · 0 discarded · 134ms
  6 tables · 6 analyzed (100%)
  stats prefilter: 3 checked · 0 rejected · 0 without statistics
  ! statistics reset time unknown — usage counters carry no meaning
`

/**
 * `pgfathom audit`, composed from the exact object names, metrics and
 * suggestion strings pinned in `internal/report/terminal_test.go` — real
 * fixtures the renderer is tested against, not invented ones.
 */
export const AUDIT_OUTPUT = `
  pgfathom v0.1.2 · PostgreSQL 16.14 · 6 tables in scope

  NOT VALID — declared, never verified against existing rows  (1)
  ${RULE}
  public.pedido.pedido_cliente_fk  child_estimated_rows=4.0M

  NO PRIMARY KEY — no row identity, sequential scans on every per-row write  (2)
  ${RULE}
  public.cadastro     promote UNIQUE(cpf) to primary key
  public.log_evento   create synthetic column idkey as primary key; schema convention: "idkey" names the primary key in 300 of 338 single-column-PK tables (89%)

  HOT COLUMN — named repeatedly in real predicates, no index leading it  (1)
  ${RULE}
  public.evento.dados  create index using gin (dados)

  6 tables · 6 analyzed (100%)
  ! statistics reset time unknown — usage counters carry no meaning
`

/**
 * `pgfathom` with no arguments — the real root command, whose only special
 * case is printing the banner before its own `--help`. Text below the glyph
 * is copied from `Long` in `internal/cli/root.go` and from cobra's generated
 * help for the command tree actually registered there: `audit`, `discover`,
 * `setup`, `version`. `completion` and `help` are cobra's own scaffolding
 * commands and are left out to keep this to the product's real surface.
 */
export const ROOT_OUTPUT = `
  pgfathom finds the relationships your database has but never declared,
  and proves them against the data instead of guessing from column names.

  It is strictly read-only: it never issues a statement that modifies the
  database under analysis, and no value from your tables ever appears in its
  output, logs or generated files.

  Usage:
    pgfathom [command]

  Available Commands:
    audit       Report structural findings that need no inference
    discover    Infer relationships the schema never declared
    setup       Guide a first run and print the command it composes
    version     Print version, commit and build date

  Flags:
        --color string      colour output: auto, always or never (default "auto")
    -h, --help               help for pgfathom
        --log-level string   diagnostic verbosity on stderr: debug, info, warn or error (default "warn")

  Use "pgfathom [command] --help" for more information about a command.
`

/**
 * `pgfathom version`, from the literal format string in
 * `internal/cli/root.go`. Version and commit are the project's real current
 * state (`CHANGELOG.md` 0.1.2, `git log -1`), not a placeholder.
 */
export const VERSION_OUTPUT = `
  pgfathom v0.1.2
  commit  869400f
  built   2026-08-13
`

export function unknownCommandOutput(input: string): string {
  return `
  Error: unknown command "${input}" for "pgfathom"
  Run 'pgfathom help' for usage.
`
}

export type CliCommandKey = 'root' | 'discover' | 'audit' | 'version'

export type CliCommand = {
  key: CliCommandKey
  /** What follows `pgfathom ` in the prompt. Empty for the root command. */
  args: string
  /** Extra words that also resolve to this command in the picker. */
  aliases: string[]
  /** One-line gloss shown next to the command in the autocomplete list. */
  hint: string
  body: string
}

export const CLI_COMMANDS: CliCommand[] = [
  {
    key: 'root',
    args: '',
    aliases: ['help', 'status', '--help'],
    hint: 'overview & available commands',
    body: ROOT_OUTPUT,
  },
  {
    key: 'discover',
    args: 'discover --schema public --full',
    aliases: [],
    hint: 'infer relationships the schema never declared',
    body: DISCOVER_OUTPUT,
  },
  {
    key: 'audit',
    args: 'audit',
    aliases: [],
    hint: 'structural findings that need no inference',
    body: AUDIT_OUTPUT,
  },
  {
    key: 'version',
    args: 'version',
    aliases: [],
    hint: 'version, commit and build date',
    body: VERSION_OUTPUT,
  },
]

export const DEFAULT_CLI_COMMAND = CLI_COMMANDS[1]

/** Strips a leading `pgfathom` a visitor might type out of habit. */
export function normalizeCliInput(raw: string): string {
  return raw.trim().replace(/^pgfathom\s+/, '').replace(/^pgfathom$/, '')
}

export function matchesCliQuery(cmd: CliCommand, query: string): boolean {
  const q = normalizeCliInput(query).toLowerCase()
  if (q === '') return true
  if (cmd.args.toLowerCase().startsWith(q)) return true
  return cmd.aliases.some((a) => a.toLowerCase().startsWith(q))
}

/** Exact resolution for Enter with no suggestion highlighted. */
export function resolveCliCommand(query: string): CliCommand | undefined {
  const q = normalizeCliInput(query).toLowerCase()
  if (q === '') return CLI_COMMANDS[0]
  return CLI_COMMANDS.find(
    (cmd) => cmd.args.toLowerCase() === q || cmd.aliases.some((a) => a.toLowerCase() === q),
  )
}
