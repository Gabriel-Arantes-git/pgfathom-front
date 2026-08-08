import { useCopy } from '../../i18n/LanguageContext'
import { Eyebrow } from '../ui/Eyebrow'
import { Reveal } from '../ui/Reveal'
import { H2, Section } from '../ui/Section'
import { GlowGrid, glowCard } from '../ui/GlowGrid'
import { useReveal } from '../../hooks/useReveal'
import { BENCHMARK, BENCHMARK_IS_PLACEHOLDER } from '../../content/data'
import type { BenchmarkRow } from '../../content/data'

/** Stacked bar: name matching, then what usage evidence adds on top. */
function Bar({ row, shown, delay }: { row: BenchmarkRow; shown: boolean; delay: number }) {
  const total = row.byName + row.byEvidence

  return (
    <div className="flex items-center gap-4">
      <div className="relative h-2 flex-1 overflow-hidden rounded-full bg-ink-600">
        <div
          className="absolute inset-y-0 left-0 rounded-full bg-accent/35 transition-[width] duration-1000 ease-[cubic-bezier(0.16,1,0.3,1)]"
          style={{ width: shown ? `${total}%` : '0%', transitionDelay: `${delay}ms` }}
        />
        <div
          className="absolute inset-y-0 left-0 rounded-full bg-accent transition-[width] duration-1000 ease-[cubic-bezier(0.16,1,0.3,1)]"
          style={{ width: shown ? `${row.byName}%` : '0%', transitionDelay: `${delay}ms` }}
        />
      </div>
      <span className="w-[72px] shrink-0 text-right font-mono text-[13px] text-bone tabular-nums">
        {total.toFixed(1)}%
      </span>
    </div>
  )
}

export function Benchmark() {
  const t = useCopy()
  const { ref, shown } = useReveal<HTMLDivElement>()

  const average =
    BENCHMARK.reduce((sum, row) => sum + row.byName + row.byEvidence, 0) / BENCHMARK.length

  return (
    <Section id="benchmark">
      <Reveal className="flex flex-wrap items-end justify-between gap-5">
        <div className="max-w-[660px]">
          <Eyebrow>{t.benchEyebrow}</Eyebrow>
          <H2>{t.benchTitle}</H2>
        </div>
        {BENCHMARK_IS_PLACEHOLDER && (
          <span className="rounded-[3px] border border-dashed border-accent/50 px-3 py-1.5 font-mono text-[10.5px] tracking-[0.12em] text-accent uppercase">
            {t.benchBadge}
          </span>
        )}
      </Reveal>

      <Reveal>
        <p className="mt-[22px] mb-0 max-w-[78ch] text-[15.5px] leading-[1.7] text-pretty text-muted">
          {t.benchLead}
        </p>
      </Reveal>

      <GlowGrid className="mt-11 grid grid-cols-1 items-start gap-8 lg:grid-cols-[1.35fr_1fr]">
        <div
          ref={ref}
          className={glowCard('relative overflow-hidden rounded-[5px] border border-hair bg-ink-800')}
        >
          <div className="grid grid-cols-[1fr_auto] items-center gap-4 border-b border-hair bg-ink-600 px-6 py-3">
            <span className="font-mono text-[10.5px] tracking-[0.13em] text-muted uppercase">
              {t.colSchema}
            </span>
            <span className="font-mono text-[10.5px] tracking-[0.13em] text-muted uppercase">
              {t.colRecovered}
            </span>
          </div>

          {BENCHMARK.map((row, i) => (
            <div
              key={row.schema}
              className={`px-6 py-4 ${i > 0 ? 'border-t border-hair' : ''}`}
            >
              <div className="mb-2.5 flex items-baseline justify-between gap-4">
                <span className="font-mono text-[13.5px] text-bone">{row.schema}</span>
                <span className="font-mono text-[11px] text-muted tabular-nums">
                  {row.tables} tables · {row.keys} keys
                </span>
              </div>
              <Bar row={row} shown={shown} delay={i * 90} />
            </div>
          ))}

          <div className="flex items-center justify-between gap-4 border-t border-hair bg-ink-700 px-6 py-4">
            <span className="font-mono text-[10.5px] tracking-[0.13em] text-muted uppercase">
              {t.benchAggregate}
            </span>
            <span className="font-mono text-[15px] text-accent tabular-nums">
              {average.toFixed(1)}%
            </span>
          </div>
        </div>

        <div className="flex flex-col gap-8">
          <div className="flex flex-wrap gap-x-6 gap-y-3">
            <span className="flex items-center gap-2.5 text-[13px] text-muted">
              <span aria-hidden="true" className="h-2.5 w-2.5 rounded-full bg-accent" />
              {t.benchLegend.byName}
            </span>
            <span className="flex items-center gap-2.5 text-[13px] text-muted">
              <span aria-hidden="true" className="h-2.5 w-2.5 rounded-full bg-accent/35" />
              {t.benchLegend.byEvidence}
            </span>
          </div>

          <div className="grid gap-px overflow-hidden rounded-[5px] border border-hair bg-hair">
            <div className={glowCard('relative bg-ink-800 px-6 py-5.5')}>
              <p className="m-0 font-mono text-[10.5px] tracking-[0.13em] text-accent uppercase">
                {t.metricLabel}
              </p>
              <p className="mt-2.5 mb-0 text-[14px] leading-[1.7] text-pretty text-muted">
                {t.metricBody}
              </p>
            </div>
            <div className={glowCard('relative bg-ink-800 px-6 py-5.5')}>
              <p className="m-0 font-mono text-[10.5px] tracking-[0.13em] text-accent uppercase">
                {t.fpLabel}
              </p>
              <p className="mt-2.5 mb-0 text-[14px] leading-[1.7] text-pretty text-muted">
                {t.fpBody}
              </p>
            </div>
          </div>
        </div>
      </GlowGrid>
    </Section>
  )
}
