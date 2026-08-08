import type { ReactNode } from 'react'
import { useLanguage } from '../../i18n/LanguageContext'
import type { Lang } from '../../i18n/types'
import { FlagBR } from './flags/FlagBR'
import { FlagUS } from './flags/FlagUS'

type Props = {
  size?: number
  /** Writes the language name next to each flag — used in the mobile menu. */
  withLabels?: boolean
}

const LABEL: Record<Lang, string> = { en: 'English', pt: 'Português' }

/** Padding around the flag inside its slot, per side. */
const PAD = 4

export function LanguageSwitch({ size = 20, withLabels = false }: Props) {
  const { lang, setLang, t } = useLanguage()

  const options: { value: Lang; flag: ReactNode; aria: string }[] = [
    { value: 'en', flag: <FlagUS size={size} />, aria: t.toEnglish },
    { value: 'pt', flag: <FlagBR size={size} />, aria: t.toPortuguese },
  ]

  if (withLabels) {
    return (
      <div role="group" aria-label={t.langSwitchLabel} className="flex items-center gap-2">
        {options.map((option) => {
          const active = lang === option.value
          return (
            <button
              key={option.value}
              type="button"
              onClick={() => setLang(option.value)}
              aria-pressed={active}
              aria-label={option.aria}
              className={`flex flex-1 cursor-pointer items-center justify-center gap-2.5 rounded-full border px-4 py-2.5 transition-all duration-[220ms] ease-[cubic-bezier(.4,0,.2,1)] ${
                active
                  ? 'border-accent opacity-100'
                  : 'border-hair opacity-50 grayscale-[.6] hover:opacity-90 hover:grayscale-0'
              }`}
            >
              {option.flag}
              <span className="font-mono text-[13px] text-bone">{LABEL[option.value]}</span>
            </button>
          )
        })}
      </div>
    )
  }

  // Slot is the flag plus PAD on every side; the ring is exactly one slot, so
  // it stays concentric with the flag at any size.
  const slot = size + PAD * 2

  return (
    <div
      role="group"
      aria-label={t.langSwitchLabel}
      className="glass relative flex items-center rounded-full border border-hair p-1"
    >
      <span
        aria-hidden="true"
        className="pointer-events-none absolute top-1 left-1 rounded-full border border-accent transition-transform duration-[220ms] ease-[cubic-bezier(.4,0,.2,1)]"
        style={{
          width: slot,
          height: slot,
          transform: `translateX(${lang === 'pt' ? slot : 0}px)`,
        }}
      />
      {options.map((option) => {
        const active = lang === option.value
        return (
          <button
            key={option.value}
            type="button"
            onClick={() => setLang(option.value)}
            aria-pressed={active}
            aria-label={option.aria}
            title={LABEL[option.value]}
            className={`relative z-10 flex cursor-pointer items-center justify-center rounded-full transition-all duration-[220ms] ease-[cubic-bezier(.4,0,.2,1)] ${
              active ? 'opacity-100' : 'opacity-45 grayscale-[.6] hover:opacity-80 hover:grayscale-0'
            }`}
            style={{ width: slot, height: slot }}
          >
            {option.flag}
          </button>
        )
      })}
    </div>
  )
}
