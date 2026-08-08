import type { Copy } from './types'

export const en: Copy = {
  docTitle: 'pgfathom: sound the depth of a legacy PostgreSQL schema',
  docDescription:
    'pgfathom finds the relationships your PostgreSQL database has but never declared, and proves them against the data instead of guessing from column names.',

  navProblem: 'Problem',
  navVerdicts: 'What it does',
  navSafety: 'Safety',
  navHow: 'How it works',
  navBenchmark: 'Benchmark',
  navMenu: 'Open menu',
  navClose: 'Close menu',

  langSwitchLabel: 'Language',
  toEnglish: 'Switch to English',
  toPortuguese: 'Mudar para português',

  heroTitle: 'Sound the depth of a legacy PostgreSQL schema.',
  heroLead:
    'pgfathom finds the relationships your database has but never declared, then proves them against the data instead of guessing from column names.',
  ctaDesign: 'Read the design',
  copyCommand: 'Copy command',
  copiedCommand: 'Copied',
  noticeTitle: 'Pre-release. Under active development.',
  noticeBody:
    'What runs today: pgfathom audit end to end, and pgfathom discover through candidate generation, scoring, the statistical prefilter and validation against the data. Join mining and the final output formats are still ahead. Terminal output shown here is the target design, not a recording. Recovery-rate benchmarks will be published once the tool runs against the reference corpus. No numbers are claimed until then.',
  facts: [
    { label: 'Read-only', detail: 'No write mode, under any flag, in any phase' },
    { label: 'Stays in memory', detail: 'Counts and names out, never a value from your tables' },
    { label: 'Four dependencies', detail: 'No cgo. go.mod is a short read' },
    { label: 'Apache-2.0', detail: 'Chosen for the explicit patent grant' },
  ],
  terminalCaption: 'pgfathom: target output',
  // work_order.tech_lead → employee.id keeps the point of the Portuguese
  // original: the column name bears no resemblance to the table it points at.
  graph: {
    pedido: { name: 'order', columns: ['id', 'customer_id', 'issued_at'] },
    cliente: { name: 'customer', columns: ['id', 'legal_name', 'active'] },
    item_pedido: { name: 'order_item', columns: ['id', 'order_id', 'quantity'] },
    os_servico: { name: 'work_order', columns: ['id', 'tech_lead', 'opened_at'] },
    funcionario: { name: 'employee', columns: ['id', 'badge_no'] },
    endereco: { name: 'address', columns: ['id', 'city_id', 'postcode'] },
    municipio: { name: 'city', columns: ['id', 'name'] },
    conta: { name: 'account', columns: ['id', 'code'] },
    lancamento: { name: 'ledger_entry', columns: ['id', 'account_id', 'amount'] },
    produto: { name: 'product', columns: ['id', 'sku', 'unit_price'] },
    nota_fiscal: { name: 'invoice', columns: ['id', 'serial_no'] },
    usuario: { name: 'app_user', columns: ['id', 'login', 'created_at'] },
    filial: { name: 'branch', columns: ['id', 'tax_id'] },
    movimento: { name: 'stock_move', columns: ['id', 'kind', 'moved_on'] },
    // The two trailing cards in the hero, past the schema graph's own box —
    // not part of GRAPH_TABLES, rendered directly in Hero.tsx.
    sessao: { name: 'session', columns: ['id', 'user_id', 'expires_at'] },
    log_evento: { name: 'audit_log', columns: ['id', 'actor_id'] },
  },
  reportRows: [
    { child: 'work_order.tech_lead', parent: '→ employee.id' },
    { child: 'order.customer_id', parent: '→ customer.id' },
    { child: 'ledger_entry.account_id', parent: '→ account.id' },
    { child: 'order_item.order_id', parent: '→ order.id' },
    { child: 'address.city_id', parent: '→ city.id' },
    { child: 'document.entity_id', parent: '→ entity.id' },
  ],
  joinMiningExample: { child: 'work_order.tech_lead', parent: 'employee.id' },

  problemEyebrow: 'The problem',
  problemTitle: 'Old PostgreSQL databases carry more structure than they declare.',
  problemP1:
    'No foreign keys. But customer_id points at customer.id in every single row. The ORM of the day never created the constraint, or someone dropped constraints for a bulk load and never put them back.',
  problemP2:
    'Two things follow. Nobody can read the model, because \\d shows nothing and your ERD tool draws a page of disconnected boxes. And because the constraint was never there, nothing ever stopped orphan rows from getting in, so they probably already did, years ago, silently, and no one has looked.',
  notValidLabel: 'The nastier variant',
  notValidBody:
    'A foreign key can be declared and still guarantee nothing, if it was created NOT VALID and never validated. It shows up in \\d. It draws an arrow in your diagram. It never checked a single pre-existing row.',

  verdictsEyebrow: 'What it does',
  verdictsTitle: 'Three verdicts, three different responses.',
  verdictsLead:
    'Every inferred relationship carries a verdict and the metric behind it. Weak and rejected candidates are reported too, so you never wonder why an obvious-looking column was ignored.',
  verdicts: [
    {
      tag: 'BROKEN',
      headline: 'The point of the tool.',
      body: 'A data bug that has been in production for years and nobody knows.',
      output:
        '→ the query that lists the orphans, because they have to be resolved before any constraint can be added.',
    },
    {
      tag: 'CONFIRMED',
      headline: 'A foreign key someone forgot to declare.',
      body: 'Containment holds by row and by distinct value, with no orphans found.',
      output:
        "→ the DDL, with VALIDATE CONSTRAINT split out so the initial ALTER TABLE doesn't hold a heavy lock, plus CREATE INDEX CONCURRENTLY when the child column is unindexed.",
    },
    {
      tag: 'WEAK',
      headline: 'Insufficient evidence to conclude.',
      body: 'Reported rather than dropped, with the reason attached.',
      output:
        '→ the metric that fell short, and any pattern detected: a polymorphic pair, an ambiguous target.',
    },
  ],
  joinMining:
    'No name-matching heuristic in the world finds that one: tech_lead looks nothing like employee. pgfathom finds it by reading the join predicates out of your own view and function definitions.',

  safetyEyebrow: 'Safety',
  safetyTitle:
    'Designed to be pointed at a production database owned by someone who is nervous about it.',
  safetyLead: 'Every guarantee below is a hard requirement in the spec, not a goal.',
  safety: [
    {
      name: 'Read-only, structurally',
      detail:
        'pgfathom never issues a statement that modifies the database under analysis. There is no write mode, under any flag, in any phase. The session sets default_transaction_read_only, and a read-only role is the recommended setup. It emits .sql files for you to review and run yourself.',
    },
    {
      name: 'Your data never leaves memory',
      detail:
        'What comes out are counts, ratios, and object names, never a value from your tables, in any output, log, JSON field, or error message. Enforced by a test that serializes every structure and scans the result, not by code review.',
    },
    {
      name: "It won't take your server down",
      detail:
        'Every validation query runs under statement_timeout, lock_timeout, and idle_in_transaction_session_timeout. Concurrency is capped and defaults low. The connection announces itself as pgfathom in pg_stat_activity.',
    },
    {
      name: 'No claim without evidence',
      detail:
        'Sampled runs can never report a confirmed relationship: only --full can prove absence of orphans.',
    },
    {
      name: 'Silence is never a clean bill of health',
      detail:
        'Tables skipped for missing privileges, candidates that timed out, schemas not covered: all of it appears in the coverage block on every run. A clean report means "I looked and it\'s clean", never "I couldn\'t look".',
    },
  ],

  howEyebrow: 'How it works',
  howTitle: 'Six passes, from catalog to validation.',
  stages: [
    {
      num: '01',
      name: 'Read the catalog',
      detail:
        'Tables, columns, keys, indexes, comments, declared foreign keys and their validation state, usage statistics with their reset timestamp.',
    },
    {
      num: '02',
      name: 'Mine usage evidence',
      detail:
        'Join predicates extracted from view definitions, function bodies, and pg_stat_statements. A view that joins two columns is proof your code treats them as related. Pure catalog: no user data, no cost.',
    },
    {
      num: '03',
      name: 'Generate candidates',
      detail:
        'Column-name affixes are stripped and matched against depluralized table names using a naming profile: a config file, not hardcoded rules.',
    },
    {
      num: '04',
      name: 'Score on metadata alone',
      detail:
        'Exact name match, type identity, target ambiguity, existing index, comment mentions. Weak candidates are dropped before anything touches data.',
    },
    {
      num: '05',
      name: 'Prefilter on planner statistics',
      detail:
        'If the child column has more distinct values than the parent has rows, full containment is arithmetically impossible. Free, from pg_stats, no I/O.',
    },
    {
      num: '06',
      name: 'Validate against the data',
      detail:
        'One aggregate per surviving candidate, never fetching rows, only counts. Containment is reported by row and by distinct value, because one bad value repeated a million times and a million rare bad values are different problems.',
    },
  ],

  profilesEyebrow: 'Naming profiles',
  profilesTitle:
    "Most schema tools assume English. The databases that need this tool most often aren't.",
  profilesP1:
    'Affix and plural rules live in TOML, not in Go, so teaching pgfathom a new convention is a config file rather than a patch.',
  profilesP2:
    'Normalization returns a set of candidate forms rather than one, so ambiguous plurals cost nothing in recall. Every form is tried, and the one that matched is reported. Adding a profile for your language is the easiest possible contribution: a TOML file and a table of test cases.',
  profilesYours: 'your language',

  ciEyebrow: 'Continuous integration',
  ciTitle: 'Fail the build when the schema drifts.',
  ciBadge: 'planned · after v0.1',
  ciP1:
    'pgfathom check --baseline compares a run against a committed model and exits non-zero when the schema has moved. Not available yet: this is the shape it is specified to take.',
  ciRules: [
    { code: 'exit 1', detail: 'A new undeclared relationship appeared since the baseline' },
    { code: 'exit 1', detail: 'An orphan count grew on a known broken relationship' },
    { code: 'exit 0', detail: 'Coverage reported, nothing changed' },
  ],

  priorEyebrow: 'Prior art',
  priorTitle: 'Not new science, and says so.',
  priorLead:
    'Containment is known in the data-profiling literature as an inclusion dependency, the automatically testable part of a foreign key. pgfathom deliberately does not compete with the tools that already solved their problems well.',
  colTool: 'Tool',
  colDoes: 'What it does well',
  colOverlap: 'Where pgfathom differs',
  priorRows: [
    {
      tool: 'Atlas',
      does: 'Schema drift and migration management',
      overlap:
        'Compares against a declared desired state; pgfathom infers the state that was never declared',
    },
    {
      tool: 'SchemaSpy · Azimutt',
      does: 'Diagrams and schema exploration',
      overlap:
        'Draw what the catalog holds; Azimutt flags _id columns, but nothing is validated against rows',
    },
    {
      tool: 'Squawk',
      does: 'Migration linting',
      overlap: 'Reviews the DDL you write; pgfathom writes the DDL you review',
    },
    {
      tool: 'Metanome (SPIDER, BINDER, MIND)',
      does: 'Inclusion-dependency discovery research',
      overlap:
        'Academic, general-purpose, offline; pgfathom is a PostgreSQL-native CLI that mines the catalog itself',
    },
    {
      tool: 'pgfathom',
      does: 'Validates its inferences against the actual data and hands you reviewable DDL',
      overlap: 'Speaks non-English naming conventions and reports its own coverage honestly',
      highlight: true,
    },
  ],
  priorNote:
    'The versioned JSON model is the integration point: consume it and generate whatever you like from a schema that finally knows its own relationships. Code generation is explicitly not on the roadmap.',

  benchEyebrow: 'How correctness is measured',
  benchTitle: 'Drop every foreign key, then count how many come back.',
  benchBadge: 'prototype · placeholder numbers',
  benchLead:
    'The corpus is public (GitLab, Odoo, Discourse, Redmine, Mastodon), so anyone can reproduce the run. Results are split into what name matching recovers alone and what usage evidence adds on top, because that gap is precisely what join mining exists to close.',
  colSchema: 'Schema',
  colRecovered: 'Foreign keys recovered',
  benchLegend: {
    byName: 'Name matching',
    byEvidence: 'Usage evidence',
  },
  benchAggregate: 'corpus average',
  metricLabel: 'Recovery rate',
  metricBody:
    'Take a schema with complete foreign keys, drop every one of them, run pgfathom, and count how many come back. Against a public corpus (GitLab, Odoo, Discourse, Redmine, Mastodon), so anyone can reproduce it.',
  fpLabel: 'Zero confirmed false positives',
  fpBody:
    'Recall will settle well below 100%, and that is expected. The metric with no tolerance is the other one: a missed relationship costs you a finding, a wrong one confirmed costs you the tool.',

  ctaTitle: 'The design is the cheapest moment to change it.',
  ctaBody:
    'If you have run into this problem on a real legacy database, open an issue: what the schema looked like, what naming convention it used, and what a tool would have needed to find. The two most valuable contributions are naming profiles for other languages and real-world schemas for the benchmark corpus.',
  ctaIssue: 'Open an issue',
  ctaRepo: 'Browse the repository',

  tagline: 'PostgreSQL schema archaeology',
  linkDesign: 'Design doc',
  linkRoadmap: 'Roadmap',
  footLeft: 'Apache-2.0 · pre-release',
  footRight: 'Not affiliated with the PostgreSQL Global Development Group.',
}
