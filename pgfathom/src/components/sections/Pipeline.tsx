import { useCopy } from '../../i18n/LanguageContext'
import { Eyebrow } from '../ui/Eyebrow'
import { Reveal } from '../ui/Reveal'
import { H2, Section } from '../ui/Section'
import { PIPELINE_STAGES } from '../../content/data'

export function Pipeline() {
  const t = useCopy()

  return (
    <Section id="pipeline">
      <Reveal className="max-w-[660px]">
        <Eyebrow>{t.howEyebrow}</Eyebrow>
        <H2>{t.howTitle}</H2>
      </Reveal>

      <Reveal className="mt-11 flex flex-wrap items-center gap-x-3.5 gap-y-2.5 font-mono text-[12px] text-muted">
        {PIPELINE_STAGES.map((stage, i) => (
          <span key={stage} className="flex items-center gap-3.5">
            <span>{stage}</span>
            {i < PIPELINE_STAGES.length - 1 && <span className="text-accent">→</span>}
          </span>
        ))}
      </Reveal>

      <ol className="mt-7 grid list-none grid-cols-1 gap-px overflow-hidden rounded-[5px] border border-hair bg-hair p-0 sm:grid-cols-2 lg:grid-cols-3">
        {t.stages.map((stage, i) => (
          <Reveal
            key={stage.num}
            as="li"
            delay={i * 60}
            className="bg-ink-800 px-6.5 pt-7 pb-7.5"
          >
            <div className="flex items-center gap-3">
              <span className="font-mono text-[11px] tracking-[0.08em] text-accent">
                {stage.num}
              </span>
              <span aria-hidden="true" className="h-px flex-1 bg-hair" />
            </div>
            <h3 className="mt-[22px] mb-0 text-[14.5px] font-medium text-bone">{stage.name}</h3>
            <p className="mt-2.5 mb-0 text-[13.5px] leading-[1.7] text-pretty text-muted">
              {stage.detail}
            </p>
          </Reveal>
        ))}
      </ol>
    </Section>
  )
}
