import { useCopy } from '../../i18n/LanguageContext'
import { Eyebrow } from '../ui/Eyebrow'
import { Reveal } from '../ui/Reveal'
import { H2, Lead, Section } from '../ui/Section'
import { GlowGrid, glowCard } from '../ui/GlowGrid'
import { SAFETY_ICONS } from '../illustrations/SafetyIcons'

export function Safety() {
  const t = useCopy()

  return (
    <Section id="safety">
      <div className="grid grid-cols-1 items-start gap-16 lg:grid-cols-[0.85fr_1.15fr]">
        <Reveal className="lg:sticky lg:top-28">
          <Eyebrow>{t.safetyEyebrow}</Eyebrow>
          <H2>{t.safetyTitle}</H2>
          <Lead>{t.safetyLead}</Lead>
        </Reveal>

        <GlowGrid
          intensity={0.6}
          className="grid content-start gap-px overflow-hidden rounded-[5px] border border-hair bg-hair"
        >
          {t.safety.map((item, i) => {
            const Icon = SAFETY_ICONS[i]
            return (
              <Reveal
                key={item.name}
                delay={i * 60}
                className={glowCard('relative flex items-start gap-5 bg-ink-800 px-6.5 py-6')}
              >
                <span className="mt-0.5 shrink-0 text-accent">{Icon && <Icon size={34} />}</span>
                <div>
                  <p className="m-0 text-[14.5px] font-medium text-bone">{item.name}</p>
                  <p className="mt-2.5 mb-0 text-[13.5px] leading-[1.7] text-pretty text-muted">
                    {item.detail}
                  </p>
                </div>
              </Reveal>
            )
          })}
        </GlowGrid>
      </div>
    </Section>
  )
}
