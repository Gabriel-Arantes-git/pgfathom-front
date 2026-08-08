import { useCopy } from '../../i18n/LanguageContext'
import { Eyebrow } from '../ui/Eyebrow'
import { Reveal } from '../ui/Reveal'
import { H2, Lead, Section } from '../ui/Section'

export function PriorArt() {
  const t = useCopy()

  return (
    <Section id="prior-art">
      <Reveal className="max-w-[700px]">
        <Eyebrow>{t.priorEyebrow}</Eyebrow>
        <H2>{t.priorTitle}</H2>
        <Lead>{t.priorLead}</Lead>
      </Reveal>

      <Reveal className="mt-12 overflow-hidden rounded-[5px] border border-hair">
        <div className="hidden grid-cols-[0.8fr_1.2fr_1fr] gap-px border-b border-hair bg-hair lg:grid">
          {[t.colTool, t.colDoes, t.colOverlap].map((col) => (
            <div key={col} className="bg-ink-600 px-6.5 py-3.5">
              <span className="font-mono text-[10.5px] tracking-[0.13em] text-muted uppercase">
                {col}
              </span>
            </div>
          ))}
        </div>

        {t.priorRows.map((row, i) => (
          <div
            key={row.tool}
            className={`grid grid-cols-1 gap-x-6 gap-y-2 px-6.5 py-5.5 transition-colors duration-200 lg:grid-cols-[0.8fr_1.2fr_1fr] ${
              i > 0 ? 'border-t border-hair' : ''
            } ${row.highlight ? 'bg-accent/[0.06]' : 'bg-ink-800 hover:bg-ink-700'}`}
          >
            <div>
              <span
                className={`font-mono text-[13.5px] ${row.highlight ? 'text-accent' : 'text-bone'}`}
              >
                {row.tool}
              </span>
            </div>
            <div>
              <span className="text-[13.5px] leading-[1.7] text-bone/[0.86]">{row.does}</span>
            </div>
            <div>
              <span className="text-[13.5px] leading-[1.7] text-muted">{row.overlap}</span>
            </div>
          </div>
        ))}
      </Reveal>

      <Reveal>
        <p className="mt-6 mb-0 max-w-[78ch] text-[14px] leading-[1.75] text-pretty text-muted">
          {t.priorNote}
        </p>
      </Reveal>
    </Section>
  )
}
