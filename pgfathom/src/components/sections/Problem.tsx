import { useCopy } from '../../i18n/LanguageContext'
import { Eyebrow } from '../ui/Eyebrow'
import { Reveal } from '../ui/Reveal'
import { H2, Section } from '../ui/Section'
import { PSQL_OUTPUT } from '../../content/data'

export function Problem() {
  const t = useCopy()

  return (
    <Section id="problem">
      <Reveal className="max-w-[660px]">
        <Eyebrow>{t.problemEyebrow}</Eyebrow>
        <H2>{t.problemTitle}</H2>
      </Reveal>

      <Reveal className="mt-13 grid grid-cols-1 items-start gap-14 lg:grid-cols-2">
        <div className="overflow-hidden rounded-[5px] border border-hair">
          <div className="border-b border-hair bg-ink-600 px-5 py-[11px]">
            <span className="font-mono text-[11px] tracking-[0.08em] text-muted">psql</span>
          </div>
          <pre className="m-0 overflow-x-auto bg-ink-700 p-5 font-mono text-[12.5px] leading-[1.65] text-bone/[0.88]">
            <code>
              <span className="text-accent">=#</span>
              {PSQL_OUTPUT.slice(2)}
            </code>
          </pre>
          <div className="border-t border-hair bg-ink-600 px-5 py-3">
            <span className="font-mono text-[11px] leading-[1.5] text-muted">
              no <span className="text-accent">"Foreign-key constraints:"</span> block. that is the
              whole problem.
            </span>
          </div>
        </div>

        <div className="flex flex-col gap-[22px]">
          <p className="m-0 text-[15.5px] leading-[1.7] text-pretty text-muted">{t.problemP1}</p>
          <p className="m-0 text-[15.5px] leading-[1.7] text-pretty text-muted">{t.problemP2}</p>
          <div className="rounded border border-hair bg-ink-600 px-5 py-[18px]">
            <p className="m-0 font-mono text-[11px] tracking-[0.13em] text-accent uppercase">
              {t.notValidLabel}
            </p>
            <p className="mt-2.5 mb-0 text-[14px] leading-[1.7] text-pretty text-muted">
              {t.notValidBody}
            </p>
          </div>
        </div>
      </Reveal>
    </Section>
  )
}
