export type Lang = 'en' | 'pt'

export type VerdictTag = 'BROKEN' | 'CONFIRMED' | 'WEAK'

export type Fact = {
  label: string
  detail: string
}

export type Verdict = {
  tag: VerdictTag
  headline: string
  body: string
  output: string
}

export type SafetyItem = {
  name: string
  detail: string
}

export type Stage = {
  num: string
  name: string
  detail: string
}

export type CiRule = {
  code: string
  detail: string
}

export type PriorRow = {
  tool: string
  does: string
  overlap: string
  highlight?: boolean
}

export type BenchmarkLegend = {
  byName: string
  byEvidence: string
}

/**
 * Localized identifiers for the hero graph. `columns` is positional: entry i
 * names the column whose type sits at index i of the matching GRAPH_TABLES
 * entry, so the two arrays must stay the same length.
 */
export type GraphLabel = {
  name: string
  columns: string[]
}

/**
 * A repository contributor, already normalized from the GitHub API's shape.
 * `login` is an identifier, never translated copy.
 */
export type Contributor = {
  login: string
  profileUrl: string
  avatarUrl: string
  contributions: number
}

/**
 * The shape of a full translation. `copy.en` is the source of truth and
 * `copy.pt` is typed against it, so a missing key fails the build instead of
 * shipping a raw key to production.
 */
export type Copy = {
  docTitle: string
  docDescription: string

  navProblem: string
  navVerdicts: string
  navSafety: string
  navHow: string
  navBenchmark: string
  navContributors: string
  navMenu: string
  navClose: string

  langSwitchLabel: string
  toEnglish: string
  toPortuguese: string

  heroTitle: string
  heroLead: string
  ctaDesign: string
  copyCommand: string
  copiedCommand: string
  noticeTitle: string
  noticeBody: string
  facts: Fact[]
  terminalCaption: string
  cliPlaceholder: string
  cliInputAria: string
  cliRunAria: string
  cliEmptyHint: string
  graph: Record<string, GraphLabel>
  joinMiningExample: { child: string; parent: string }

  problemEyebrow: string
  problemTitle: string
  problemP1: string
  problemP2: string
  notValidLabel: string
  notValidBody: string

  verdictsEyebrow: string
  verdictsTitle: string
  verdictsLead: string
  verdicts: Verdict[]
  joinMining: string

  safetyEyebrow: string
  safetyTitle: string
  safetyLead: string
  safety: SafetyItem[]

  howEyebrow: string
  howTitle: string
  stages: Stage[]

  profilesEyebrow: string
  profilesTitle: string
  profilesP1: string
  profilesP2: string
  profilesYours: string

  ciEyebrow: string
  ciTitle: string
  ciBadge: string
  ciP1: string
  ciRules: CiRule[]

  priorEyebrow: string
  priorTitle: string
  priorLead: string
  colTool: string
  colDoes: string
  colOverlap: string
  priorRows: PriorRow[]
  priorNote: string

  benchEyebrow: string
  benchTitle: string
  benchBadge: string
  benchLead: string
  colSchema: string
  colRecovered: string
  benchLegend: BenchmarkLegend
  benchAggregate: string
  metricLabel: string
  metricBody: string
  fpLabel: string
  fpBody: string

  contributorsEyebrow: string
  contributorsTitle: string
  contributorsLead: string
  contributorsPanelLabel: string
  contributorsMainCommitsLabel: string
  contributorsCommitsLabel: string
  /** Contains `{login}`, replaced with the contributor's login at render. */
  contributorsProfileLabel: string
  /** Contains `{login}`. Labels the card, which selects rather than navigates. */
  contributorsSelectLabel: string

  ctaTitle: string
  ctaBody: string
  ctaIssue: string
  ctaRepo: string

  tagline: string
  linkDesign: string
  linkRoadmap: string
  footLeft: string
  footRight: string
}
