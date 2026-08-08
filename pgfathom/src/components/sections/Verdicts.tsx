import { useCopy } from '../../i18n/LanguageContext'
import { Eyebrow } from '../ui/Eyebrow'
import { Reveal } from '../ui/Reveal'
import { H2, Lead, Section } from '../ui/Section'
import { GlowGrid, glowCard } from '../ui/GlowGrid'
import { VERDICT_ICONS } from '../illustrations/VerdictIcons'
import type { VerdictTag } from '../../i18n/types'

/** BROKEN is the accent — it is the finding that matters. */
const STYLE: Record<VerdictTag, { text: string; border: string }> = {
  BROKEN: { text: 'text-accent', border: 'border-accent/45' },
  CONFIRMED: { text: 'text-bone', border: 'border-bone/25' },
  WEAK: { text: 'text-muted', border: 'border-hair' },
}

export function Verdicts() {
  const t = useCopy()

  return (
    <Section id="verdicts">
      <Reveal className="max-w-[660px]">
        <Eyebrow>{t.verdictsEyebrow}</Eyebrow>
        <H2>{t.verdictsTitle}</H2>
        <Lead>{t.verdictsLead}</Lead>
      </Reveal>

      <GlowGrid className="mt-13 grid grid-cols-1 gap-px overflow-hidden rounded-[5px] border border-hair bg-hair lg:grid-cols-3">
        {t.verdicts.map((verdict, i) => {
          const style = STYLE[verdict.tag]
          const Icon = VERDICT_ICONS[verdict.tag]

          return (
            <Reveal
              key={verdict.tag}
              delay={i * 80}
              className={glowCard('relative bg-ink-800 px-6.5 pt-7 pb-7.5')}
            >
              <div className="flex items-center justify-between gap-4">
                <span
                  className={`inline-flex items-center rounded-[3px] border px-2.5 py-[3px] font-mono text-[10.5px] tracking-[0.12em] ${style.text} ${style.border}`}
                >
                  {verdict.tag}
                </span>
                <span className={style.text}>
                  <Icon size={86} />
                </span>
              </div>

              <p className="mt-5 mb-0 text-[16px] font-medium tracking-[-0.01em] text-pretty text-bone">
                {verdict.headline}
              </p>
              <p className="mt-3 mb-0 text-[13.5px] leading-[1.7] text-pretty text-muted">
                {verdict.body}
              </p>
              <p className="mt-[18px] mb-0 border-t border-hair pt-4 font-mono text-[11.5px] leading-[1.6] text-bone/75">
                {verdict.output}
              </p>
            </Reveal>
          )
        })}
      </GlowGrid>

      {/* The relationship no name-matching heuristic can reach. */}
      <Reveal className="mt-6 flex flex-wrap items-center gap-5 rounded-[5px] border border-hair bg-ink-700 p-6.5">
        <code className="font-mono text-[14px] whitespace-nowrap text-bone">
          os_servico.resp_tecnico <span className="text-accent">→</span> funcionario.id
        </code>
        <p className="m-0 min-w-[280px] flex-1 text-[14px] leading-[1.7] text-pretty text-muted">
          {t.joinMining}
        </p>
      </Reveal>
    </Section>
  )
}
