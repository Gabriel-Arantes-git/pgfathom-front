import { useState } from 'react'
import { useCopy } from '../../i18n/LanguageContext'
import { useContributors } from '../../hooks/useContributors'
import { useReveal } from '../../hooks/useReveal'
import { Eyebrow } from '../ui/Eyebrow'
import { Reveal } from '../ui/Reveal'
import { H2, Lead, Section } from '../ui/Section'
import { GlowGrid, glowCard } from '../ui/GlowGrid'
import { GithubIcon } from '../ui/icons'
import type { Contributor } from '../../i18n/types'

const ROTATION_MS = 4000

/**
 * Feathered only on two sides — left and bottom — so the avatar reads as
 * docked into the panel's own top-right corner (flush, no fade there) and
 * dissolving away toward the text and footer on the other two. A diagonal
 * fade, not a vignette on all four sides.
 *
 * Percentages are relative to the image's own box, not the panel, so a
 * small photo and a large one fade by the same proportion.
 */
const AVATAR_MASK = [
  'linear-gradient(to right, transparent 0%, #000 22%)',
  'linear-gradient(to bottom, #000 0%, #000 92%, transparent 100%)',
].join(', ')

type CardProps = {
  contributor: Contributor
  index: number
  active: boolean
  /** The active card's countdown, unless the carousel is standing still. */
  countdown: boolean
  paused: boolean
  commitsLabel: string
  selectLabel: string
  onSelect: () => void
  onCountdownEnd: () => void
}

function ContributorCard({
  contributor,
  index,
  active,
  countdown,
  paused,
  commitsLabel,
  selectLabel,
  onSelect,
  onCountdownEnd,
}: CardProps) {
  const { ref, shown } = useReveal<HTMLButtonElement>()

  return (
    <button
      ref={ref}
      type="button"
      aria-pressed={active}
      aria-label={selectLabel.replace('{login}', contributor.login)}
      onClick={onSelect}
      onMouseEnter={onSelect}
      onFocus={onSelect}
      style={{ transitionDelay: `${index * 60}ms` }}
      className={glowCard(
        `reveal ${shown ? 'reveal-in' : ''} relative flex w-full cursor-pointer items-center gap-4 bg-ink-800 px-6.5 py-5.5 text-left`,
      )}
    >
      {/* The active tint is its own layer because `.reveal` already owns this
          element's `transition` property — a `transition-colors` utility here
          would replace it and kill the scroll reveal. */}
      <span
        aria-hidden="true"
        className={`absolute inset-0 z-0 bg-ink-600 transition-opacity duration-300 ${
          active ? 'opacity-100' : 'opacity-0'
        }`}
      />

      <img
        src={contributor.avatarUrl}
        alt=""
        width={44}
        height={44}
        loading="lazy"
        referrerPolicy="no-referrer"
        className={`relative z-10 h-11 w-11 shrink-0 rounded-full border object-cover transition-colors duration-300 ${
          active ? 'border-accent/50' : 'border-hair'
        }`}
      />

      <span className="relative z-10 min-w-0 flex-1">
        <span className="block truncate font-mono text-[13.5px] text-bone">
          {contributor.login}
        </span>
        {contributor.contributions > 0 && (
          <span className="mt-1 block text-[12px] text-muted">
            {contributor.contributions} {commitsLabel}
          </span>
        )}
      </span>

      <span
        className={`relative z-10 h-1.5 w-1.5 shrink-0 rounded-full transition-colors duration-300 ${
          active ? 'bg-accent' : 'bg-transparent'
        }`}
      />

      {countdown && (
        <span
          aria-hidden="true"
          onAnimationEnd={onCountdownEnd}
          style={{
            animationDuration: `${ROTATION_MS}ms`,
            animationPlayState: paused ? 'paused' : 'running',
          }}
          className="progress-sweep absolute right-0 bottom-0 left-0 z-20 h-px bg-accent"
        />
      )}
    </button>
  )
}

export function Contributors() {
  const t = useCopy()
  const contributors = useContributors()
  const [activeIndex, setActiveIndex] = useState(0)
  const [paused, setPaused] = useState(false)

  // Read before the first render, not in an effect: the global reduced-motion
  // rule collapses the countdown to 0.01ms, so a single frame with the bar
  // mounted would fire animationend and jump the carousel.
  const [reducedMotion] = useState(
    () => window.matchMedia('(prefers-reduced-motion: reduce)').matches,
  )

  // The refreshed list can be shorter than the fallback that seeded it.
  const active = Math.min(activeIndex, contributors.length - 1)
  const featured = contributors[active]
  const rotating = !reducedMotion && contributors.length > 1

  return (
    <Section id="contributors">
      <Reveal className="max-w-[660px]">
        <Eyebrow>{t.contributorsEyebrow}</Eyebrow>
        <H2>{t.contributorsTitle}</H2>
        <Lead>{t.contributorsLead}</Lead>
      </Reveal>

      <div className="mt-13 grid grid-cols-1 gap-6 lg:grid-cols-[1.05fr_0.95fr]">
        <Reveal className="relative min-h-[440px] overflow-hidden rounded-[5px] border border-hair bg-ink-700 p-7 md:min-h-[520px]">
          {/* The avatar layer ends where the footer's clear zone begins, so the
              image can be as large as it likes inside it without ever reaching
              the commit count — the clearance is structural, not a percentage
              tuned by eye. */}
          <div aria-hidden="true" className="pointer-events-none absolute top-0 right-0 bottom-28 left-0">
            {contributors.map((contributor, i) => (
              <img
                key={contributor.login}
                src={contributor.avatarUrl}
                alt=""
                width={520}
                height={520}
                loading="lazy"
                referrerPolicy="no-referrer"
                className="absolute top-0 right-0 h-[94%] w-auto saturate-[0.85] transition-opacity duration-500"
                style={{
                  opacity: i === active ? 0.9 : 0,
                  maskImage: AVATAR_MASK,
                  WebkitMaskImage: AVATAR_MASK,
                  maskComposite: 'intersect',
                  WebkitMaskComposite: 'source-in',
                }}
              />
            ))}
            {/* A colour wash, not a mask on the photo — the avatar keeps its
                own light edge feather everywhere, and only the text column
                gets darkened underneath it so the number stays legible. */}
            <span className="absolute inset-0 bg-gradient-to-r from-ink-700 via-ink-700/55 to-transparent" />
          </div>

          <div className="relative z-10">
            <p className="m-0 font-mono text-[11px] tracking-[0.16em] text-muted uppercase">
              {t.contributorsPanelLabel}
            </p>
            <p className="mt-7 mb-0 text-[64px] leading-none font-medium tracking-[-0.03em] text-bone md:text-[76px]">
              {featured?.contributions ?? 0}
            </p>
            <p className="mt-3 mb-0 max-w-[22ch] text-[14px] leading-[1.6] text-pretty text-muted">
              {t.contributorsMainCommitsLabel}
            </p>
          </div>

          {featured && (
            <div className="absolute right-7 bottom-7 left-7 z-20 flex items-center justify-between gap-4 border-t border-hair pt-5">
              <span className="min-w-0 truncate font-mono text-[14px] text-bone">
                {featured.login}
              </span>

              <a
                href={featured.profileUrl}
                target="_blank"
                rel="noopener noreferrer"
                aria-label={t.contributorsProfileLabel.replace('{login}', featured.login)}
                className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full border border-hair text-muted transition-colors duration-200 hover:border-accent/50 hover:bg-ink-600 hover:text-bone"
              >
                <GithubIcon size={17} />
              </a>
            </div>
          )}
        </Reveal>

        {/* self-start opts out of the grid's default row stretch — without it
            this column matched the photo panel's full height regardless of
            contributor count, leaving dead space below a short card list. */}
        <div
          className="self-start"
          onMouseEnter={() => setPaused(true)}
          onMouseLeave={() => setPaused(false)}
          onFocusCapture={() => setPaused(true)}
          onBlurCapture={() => setPaused(false)}
        >
          <GlowGrid
            intensity={0.6}
            className="grid content-start gap-px overflow-hidden rounded-[5px] border border-hair bg-hair"
          >
            {contributors.map((contributor, i) => (
              <ContributorCard
                key={contributor.login}
                contributor={contributor}
                index={i}
                active={i === active}
                countdown={rotating && i === active}
                paused={paused}
                commitsLabel={t.contributorsCommitsLabel}
                selectLabel={t.contributorsSelectLabel}
                onSelect={() => setActiveIndex(i)}
                onCountdownEnd={() => setActiveIndex((prev) => (prev + 1) % contributors.length)}
              />
            ))}
          </GlowGrid>
        </div>
      </div>
    </Section>
  )
}
