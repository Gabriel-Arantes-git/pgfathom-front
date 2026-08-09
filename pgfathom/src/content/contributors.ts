import type { Contributor } from '../i18n/types'
import { REPO } from './data'

/**
 * The contributor list is public data, so it is fetched without a token — a PAT
 * shipped in a front-end bundle is public by definition.
 *
 * The endpoint is derived from REPO rather than written out, so the page and
 * the API can never end up pointing at different repositories.
 */
export const CONTRIBUTORS_ENDPOINT = `${REPO.replace(
  'https://github.com/',
  'https://api.github.com/repos/',
)}/contributors?per_page=10`

export const MAX_CONTRIBUTORS = 5

/**
 * Deliberately larger than MAX_CONTRIBUTORS above: bots count as contributors
 * in the response and are dropped client-side, so asking for exactly five would
 * return fewer than five people whenever one lands in the top of the list.
 */

export const CONTRIBUTORS_CACHE_KEY = 'pgfathom-contributors'
export const CONTRIBUTORS_CACHE_TTL_MS = 6 * 60 * 60 * 1000

/**
 * Rendered on the first paint and kept whenever the fetch fails, so the section
 * has no empty state and no loading skeleton. Commit counts are zero because
 * the refresh overwrites them — the names are what has to be right here.
 *
 * `github.com/<login>.png` is a stable avatar route keyed by login, so the
 * fallback needs no prior API response to show a face.
 */
export const CONTRIBUTORS_FALLBACK: Contributor[] = [
  {
    login: 'lvcas-dotcom',
    profileUrl: 'https://github.com/lvcas-dotcom',
    avatarUrl: 'https://github.com/lvcas-dotcom.png?size=400',
    contributions: 0,
  },
  {
    login: 'Gabriel-Arantes-git',
    profileUrl: 'https://github.com/Gabriel-Arantes-git',
    avatarUrl: 'https://github.com/Gabriel-Arantes-git.png?size=400',
    contributions: 0,
  },
]
