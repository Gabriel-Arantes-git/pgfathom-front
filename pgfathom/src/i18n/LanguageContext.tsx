import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import type { Copy, Lang } from './types'
import { en } from './copy.en'
import { pt } from './copy.pt'

const DICT: Record<Lang, Copy> = { en, pt }
const STORAGE_KEY = 'pgfathom-lang'

type LanguageValue = {
  lang: Lang
  setLang: (lang: Lang) => void
  t: Copy
}

const LanguageContext = createContext<LanguageValue | null>(null)

function readInitialLang(): Lang {
  if (typeof window === 'undefined') return 'en'
  try {
    const saved = window.localStorage.getItem(STORAGE_KEY)
    if (saved === 'en' || saved === 'pt') return saved
  } catch {
    // localStorage can throw in private mode — fall through to default
  }
  return 'en'
}

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [lang, setLangState] = useState<Lang>(readInitialLang)

  const setLang = useCallback((next: Lang) => {
    setLangState(next)
    try {
      window.localStorage.setItem(STORAGE_KEY, next)
    } catch {
      // persistence is a nicety, not a requirement
    }
  }, [])

  const t = DICT[lang]

  useEffect(() => {
    document.documentElement.lang = lang === 'pt' ? 'pt-BR' : 'en'
    document.title = t.docTitle

    let meta = document.querySelector<HTMLMetaElement>('meta[name="description"]')
    if (!meta) {
      meta = document.createElement('meta')
      meta.name = 'description'
      document.head.appendChild(meta)
    }
    meta.content = t.docDescription
  }, [lang, t])

  const value = useMemo(() => ({ lang, setLang, t }), [lang, setLang, t])

  return <LanguageContext value={value}>{children}</LanguageContext>
}

export function useLanguage(): LanguageValue {
  const ctx = useContext(LanguageContext)
  if (!ctx) throw new Error('useLanguage must be used inside <LanguageProvider>')
  return ctx
}

/** Shorthand for the common case of only needing the strings. */
export function useCopy(): Copy {
  return useLanguage().t
}
