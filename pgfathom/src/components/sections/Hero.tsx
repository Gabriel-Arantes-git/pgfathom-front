import { useState } from 'react'
import { useLanguage } from '../../i18n/LanguageContext'
import { SchemaGraph } from '../canvas/SchemaGraph'
import { TrailingTable } from '../canvas/TrailingTable'
import { CliTerminal } from '../canvas/CliTerminal'
import { ArrowRightIcon, CheckIcon, CopyIcon, InfoIcon } from '../ui/icons'
import { COMMAND, DESIGN_DOC } from '../../content/data'

export function Hero() {
  const { t, lang } = useLanguage()
  const [copied, setCopied] = useState(false)

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(COMMAND)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
    }
  }

  return (
    <section id="top" className="relative overflow-hidden border-b border-hair">
      {/* Hairline grid, faded out radially from the top edge. */}
      <div
        aria-hidden="true"
        className="grid-veil pointer-events-none absolute inset-0 opacity-35"
        style={{
          maskImage: 'radial-gradient(ellipse 85% 55% at 50% 0%, #000 25%, transparent 100%)',
          WebkitMaskImage: 'radial-gradient(ellipse 85% 55% at 50% 0%, #000 25%, transparent 100%)',
        }}
      />

      {/* The schema stitching itself together. The soft edges are gradient
          overlays, NOT a CSS mask: a mask on this wrapper forced the browser to
          re-rasterize the whole animated subtree (nine blurred cards) on every
          frame, which quantized the slow drift to whole pixels and made it
          step. Plain overlays keep each card on its own compositor layer, so
          sub-pixel motion stays smooth. */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute top-[104px] right-0 hidden h-[520px] w-[660px] isolate xl:block"
      >
        <SchemaGraph className="h-full w-full" />
        {/* Hand-written gradients rather than Tailwind's. CSS `transparent` is
            transparent *black*, and Tailwind v4 interpolates gradient stops in
            oklab — `to-transparent` and even `to-ink-800/0` both compile down
            to it, dragging the ramp through colours darker than the page and
            leaving a visible dark block. Repeating the base colour at alpha 0
            keeps the ramp on one hue. */}
        <div
          className="absolute inset-y-0 left-0 z-[120] w-16"
          style={{ background: 'linear-gradient(to right, #1b1512 0%, rgba(27,21,18,0) 100%)' }}
        />
        <div
          className="absolute inset-x-0 bottom-0 z-[120] h-14"
          style={{ background: 'linear-gradient(to top, #1b1512 0%, rgba(27,21,18,0) 100%)' }}
        />
      </div>

      {/* Two lone cards trailing further down than the schema graph's own box,
          descending toward the facts strip. Positions are pixel-measured
          against the actual rendered page (Playwright, 1280–1920px, both
          languages) rather than guessed: the facts strip lands as high as
          1398px from the section top at the narrowest supported width (1280,
          English copy), so the lower card stops at 1300 — clear of that in
          every measured combination, not just the one on screen at the time. */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 hidden xl:block"
      >
        <TrailingTable
          label={t.graph.sessao}
          columns={[{ type: 'bigint' }, { type: 'bigint', fk: true }, { type: 'timestamptz' }]}
          top={760}
          right={210}
          depth={0.55}
        />
        <TrailingTable
          label={t.graph.log_evento}
          columns={[{ type: 'bigint' }, { type: 'bigint', fk: true }]}
          top={1042}
          right={60}
          depth={0.75}
        />
      </div>

      <div className="relative mx-auto max-w-[1140px] px-7 pt-[112px] pb-20 md:pt-[136px]">
        <div className="fade-up flex flex-wrap items-center gap-2.5 font-mono text-[11px] tracking-[0.15em] text-muted uppercase">
          <span aria-hidden="true" className="h-[5px] w-[5px] bg-accent" />
          <span>Go 1.25+</span>
          <span className="text-hair">/</span>
          <span>PostgreSQL 13+</span>
          <span className="text-hair">/</span>
          <span>read-only</span>
          <span className="text-hair">/</span>
          <span>Apache-2.0</span>
        </div>

        <h1
          className="fade-up mt-7 mb-0 max-w-[17ch] text-[38px] leading-[1.02] font-medium tracking-[-0.035em] text-balance md:text-[62px]"
          style={{ animationDelay: '70ms' }}
        >
          {t.heroTitle}
        </h1>

        <p
          className="fade-up mt-6 mb-0 max-w-[60ch] text-[18px] leading-[1.6] text-pretty text-muted"
          style={{ animationDelay: '140ms' }}
        >
          {t.heroLead}
        </p>

        <div
          className="fade-up mt-9 flex flex-wrap items-center gap-3"
          style={{ animationDelay: '210ms' }}
        >
          <button
            type="button"
            onClick={copy}
            aria-label={copied ? t.copiedCommand : t.copyCommand}
            className="flex h-11 cursor-pointer items-center gap-3.5 rounded border border-hair bg-ink-600 pr-2 pl-4 font-mono text-[13.5px] text-bone transition-colors duration-200 hover:border-hair-hi"
          >
            <span className="flex items-center gap-2.5">
              <span className="text-accent">$</span>
              {COMMAND}
            </span>
            <span className="flex h-7 w-7 items-center justify-center rounded-[3px] text-muted">
              {copied ? <CheckIcon className="text-accent" /> : <CopyIcon />}
            </span>
          </button>

          <a
            href={DESIGN_DOC[lang]}
            target="_blank"
            rel="noopener noreferrer"
            className="flex h-11 items-center gap-2 rounded bg-accent px-[18px] text-[13.5px] font-medium text-ink-800 transition-colors duration-200 hover:bg-accent-hi"
          >
            {t.ctaDesign}
            <ArrowRightIcon />
          </a>
        </div>

        {/* No fill: on the dark base even a 5% tint reads as an unexplained
            panel floating over a continuous background. The accent rule alone
            carries the callout. */}
        <div
          className="fade-up mt-7 flex max-w-[74ch] gap-3.5 border-l-2 border-accent py-1 pr-5 pl-5"
          style={{ animationDelay: '270ms' }}
        >
          <InfoIcon className="mt-0.5 shrink-0 text-accent" />
          <p className="m-0 text-[13.5px] leading-[1.65] text-pretty text-muted">
            <strong className="font-medium text-bone">{t.noticeTitle}</strong> {t.noticeBody}
          </p>
        </div>

        <CliTerminal className="fade-up mt-14" style={{ animationDelay: '330ms' }} />
      </div>

      <div className="border-t border-hair">
        <div className="mx-auto grid max-w-[1140px] grid-cols-1 sm:grid-cols-2 lg:grid-cols-4">
          {t.facts.map((fact, i) => (
            <div
              key={fact.label}
              className={`px-7 py-5 ${i < t.facts.length - 1 ? 'border-b border-hair lg:border-r lg:border-b-0' : ''}`}
            >
              <p className="m-0 font-mono text-[11px] tracking-[0.13em] text-accent uppercase">
                {fact.label}
              </p>
              <p className="mt-[7px] mb-0 text-[13px] leading-[1.55] text-muted">{fact.detail}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
