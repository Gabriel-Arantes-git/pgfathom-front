import { useCopy } from '../../i18n/LanguageContext'
import { Eyebrow } from '../ui/Eyebrow'
import { Reveal } from '../ui/Reveal'
import { H2, Section } from '../ui/Section'
import { CI_WORKFLOW_FILE } from '../../content/data'

export function Ci() {
  const t = useCopy()

  return (
    <Section id="ci">
      <Reveal className="flex flex-wrap items-end justify-between gap-5">
        <div className="max-w-[600px]">
          <Eyebrow>{t.ciEyebrow}</Eyebrow>
          <H2>{t.ciTitle}</H2>
        </div>
        <span className="rounded-[3px] border border-dashed border-accent/50 px-3 py-1.5 font-mono text-[10.5px] tracking-[0.12em] text-accent uppercase">
          {t.ciBadge}
        </span>
      </Reveal>

      <Reveal className="mt-11 grid grid-cols-1 items-start gap-8 lg:grid-cols-2">
        <div className="overflow-hidden rounded-[5px] border border-hair">
          <div className="border-b border-hair bg-ink-600 px-5 py-[11px]">
            <span className="font-mono text-[11px] text-muted">{CI_WORKFLOW_FILE}</span>
          </div>
          <pre className="m-0 overflow-x-auto bg-ink-700 px-5 py-[22px] font-mono text-[12.5px] leading-[1.75] text-bone/[0.88]">
            <code>
              <span className="text-muted">- name:</span>
              {' Schema integrity\n  '}
              <span className="text-muted">run:</span>
              {' |\n    pgfathom check --baseline schema.json \\\n      --schema public \\\n      --fail-on broken\n  '}
              <span className="text-muted">env:</span>
              {'\n    PGFATHOM_DSN: ${{ secrets.PGFATHOM_DSN }}'}
            </code>
          </pre>
        </div>

        <div className="flex flex-col gap-5">
          <p className="m-0 text-[15.5px] leading-[1.7] text-pretty text-muted">{t.ciP1}</p>
          <div className="grid gap-px overflow-hidden rounded-[5px] border border-hair bg-hair">
            {t.ciRules.map((rule) => (
              <div key={rule.detail} className="flex items-baseline gap-3.5 bg-ink-800 px-5 py-4">
                <span className="shrink-0 font-mono text-[12px] text-accent">{rule.code}</span>
                <span className="text-[13.5px] leading-[1.65] text-muted">{rule.detail}</span>
              </div>
            ))}
          </div>
        </div>
      </Reveal>
    </Section>
  )
}
