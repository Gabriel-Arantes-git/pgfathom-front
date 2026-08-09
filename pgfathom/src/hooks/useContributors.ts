import { useEffect, useState } from 'react'
import type { Contributor } from '../i18n/types'
import {
  CONTRIBUTORS_CACHE_KEY,
  CONTRIBUTORS_CACHE_TTL_MS,
  CONTRIBUTORS_ENDPOINT,
  CONTRIBUTORS_FALLBACK,
  MAX_CONTRIBUTORS,
} from '../content/contributors'

type CachedContributors = {
  savedAt: number
  list: Contributor[]
}

function isContributorPayload(value: unknown): value is {
  login: string
  html_url: string
  avatar_url: string
  contributions: number
  type: string
} {
  if (typeof value !== 'object' || value === null) return false
  const candidate = value as Record<string, unknown>
  return (
    typeof candidate.login === 'string' &&
    typeof candidate.html_url === 'string' &&
    typeof candidate.avatar_url === 'string' &&
    typeof candidate.contributions === 'number' &&
    candidate.type === 'User'
  )
}

function selectTopContributors(payload: unknown): Contributor[] {
  if (!Array.isArray(payload)) return []

  return payload
    .filter(isContributorPayload)
    .map((entry) => ({
      login: entry.login,
      profileUrl: entry.html_url,
      avatarUrl: entry.avatar_url,
      contributions: entry.contributions,
    }))
    .sort((a, b) => b.contributions - a.contributions)
    .slice(0, MAX_CONTRIBUTORS)
}

function readContributorsCache(): Contributor[] | null {
  try {
    const raw = window.sessionStorage.getItem(CONTRIBUTORS_CACHE_KEY)
    if (!raw) return null

    const cached = JSON.parse(raw) as CachedContributors
    if (typeof cached?.savedAt !== 'number' || !Array.isArray(cached.list)) return null
    if (cached.list.length === 0) return null
    if (Date.now() - cached.savedAt > CONTRIBUTORS_CACHE_TTL_MS) return null

    return cached.list
  } catch {
    // sessionStorage throws in private mode and JSON.parse throws on a
    // hand-edited entry — either way the network path is still available
    return null
  }
}

function writeContributorsCache(list: Contributor[]) {
  try {
    const cached: CachedContributors = { savedAt: Date.now(), list }
    window.sessionStorage.setItem(CONTRIBUTORS_CACHE_KEY, JSON.stringify(cached))
  } catch {
    // caching is a nicety, not a requirement
  }
}

/**
 * The repository's top contributors, capped at MAX_CONTRIBUTORS.
 *
 * Starts at the committed fallback and only ever moves to a list that is
 * non-empty after filtering, so every failure path — offline, rate limited,
 * malformed body, a response holding nothing but bots — leaves the section
 * rendering real names rather than an empty state.
 */
export function useContributors(): Contributor[] {
  const [contributors, setContributors] = useState<Contributor[]>(CONTRIBUTORS_FALLBACK)

  useEffect(() => {
    const cached = readContributorsCache()
    if (cached) {
      setContributors(cached)
      return
    }

    const controller = new AbortController()

    async function loadContributorsFromGithub() {
      try {
        const response = await fetch(CONTRIBUTORS_ENDPOINT, {
          signal: controller.signal,
          headers: { Accept: 'application/vnd.github+json' },
        })
        if (!response.ok) return

        const list = selectTopContributors(await response.json())
        if (list.length === 0) return

        writeContributorsCache(list)
        setContributors(list)
      } catch (error) {
        if (error instanceof DOMException && error.name === 'AbortError') return
        console.warn('pgfathom: contributor list kept at fallback —', error)
      }
    }

    void loadContributorsFromGithub()

    return () => controller.abort()
  }, [])

  return contributors
}
