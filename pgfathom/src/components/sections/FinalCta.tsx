import { useRef, useState } from 'react'
import { useCopy } from '../../i18n/LanguageContext'
import { Reveal } from '../ui/Reveal'
import { ArrowRightIcon, GithubIcon } from '../ui/icons'
import { ISSUES, REPO } from '../../content/data'
import mark from '../../assets/icons/pgfathom-logo.png'

export function FinalCta() {
  const t = useCopy()
  const [spot, setSpot] = useState({ x: 50, y: 50 })
  const sectionRef = useRef<HTMLElement>(null)

  const onMove = (e: React.MouseEvent<HTMLElement>) => {
    const rect = e.currentTarget.getBoundingClientRect()
    setSpot({
      x: ((e.clientX - rect.left) / rect.width) * 100,
      y: ((e.clientY - rect.top) / rect.height) * 100,
    })
  }

  return (
    <section
      ref={sectionRef}
      onMouseMove={onMove}
      className="relative overflow-hidden border-b border-hair"
    >
      <div
        aria-hidden="true"
        className="grid-veil pointer-events-none absolute inset-0 opacity-30"
        style={{
          maskImage: 'radial-gradient(ellipse 65% 100% at 50% 50%, #000 5%, transparent 78%)',
          WebkitMaskImage: 'radial-gradient(ellipse 65% 100% at 50% 50%, #000 5%, transparent 78%)',
        }}
      />

      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 transition-opacity duration-300"
        style={{
          background: `radial-gradient(600px circle at ${spot.x}% ${spot.y}%, rgba(226,72,61,0.10), transparent 45%)`,
        }}
      />

      <Reveal className="relative z-10 mx-auto flex max-w-[1140px] flex-col items-center px-7 py-[72px] text-center md:py-26">
        <img src={mark} alt="" aria-hidden="true" className="h-13 w-13 rounded-md" />

        <h2 className="mt-7 mb-0 max-w-[22ch] text-[27px] leading-[1.12] font-medium tracking-[-0.025em] text-balance md:text-[38px]">
          {t.ctaTitle}
        </h2>

        <p className="mt-5 mb-0 max-w-[60ch] text-[15.5px] leading-[1.7] text-pretty text-muted">
          {t.ctaBody}
        </p>

        <div className="mt-8 flex flex-wrap items-center justify-center gap-3">
          <a
            href={ISSUES}
            target="_blank"
            rel="noopener noreferrer"
            className="flex h-11 items-center gap-2 rounded bg-accent px-[18px] text-[13.5px] font-medium text-ink-800 transition-colors duration-200 hover:bg-accent-hi"
          >
            {t.ctaIssue}
            <ArrowRightIcon />
          </a>
          <a
            href={REPO}
            target="_blank"
            rel="noopener noreferrer"
            className="flex h-11 items-center gap-2.5 rounded border border-hair px-[18px] text-[13.5px] font-medium text-bone transition-colors duration-200 hover:border-hair-hi hover:bg-ink-600"
          >
            <GithubIcon />
            {t.ctaRepo}
          </a>
        </div>
      </Reveal>
    </section>
  )
}
