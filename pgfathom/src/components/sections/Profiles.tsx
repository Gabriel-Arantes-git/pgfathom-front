import { useCopy } from '../../i18n/LanguageContext'
import { Eyebrow } from '../ui/Eyebrow'
import { Reveal } from '../ui/Reveal'
import { H2, Section } from '../ui/Section'
import { PROFILE_FILE, PROFILE_LANGS } from '../../content/data'

const S = ({ children }: { children: React.ReactNode }) => (
  <span className="text-accent">{children}</span>
)

export function Profiles() {
  const t = useCopy()

  return (
    <Section id="profiles">
      <div className="grid grid-cols-1 items-center gap-16 lg:grid-cols-2">
        <Reveal>
          <Eyebrow>{t.profilesEyebrow}</Eyebrow>
          <H2>{t.profilesTitle}</H2>
          <p className="mt-[22px] mb-0 text-[15.5px] leading-[1.7] text-pretty text-muted">
            {t.profilesP1}
          </p>
          <p className="mt-[18px] mb-0 text-[15.5px] leading-[1.7] text-pretty text-muted">
            {t.profilesP2}
          </p>

          <div className="mt-6 flex flex-wrap gap-2">
            {PROFILE_LANGS.map((lang) => (
              <span
                key={lang}
                className="rounded-[3px] border border-hair bg-ink-600 px-2.5 py-[5px] font-mono text-[12px] text-bone"
              >
                {lang}
              </span>
            ))}
            <span className="rounded-[3px] border border-dashed border-accent/50 px-2.5 py-[5px] font-mono text-[12px] text-accent">
              {t.profilesYours}
            </span>
          </div>
        </Reveal>

        <Reveal delay={100} className="overflow-hidden rounded-[5px] border border-hair">
          <div className="flex items-center justify-between border-b border-hair bg-ink-600 px-5 py-[11px]">
            <span className="font-mono text-[11px] text-muted">{PROFILE_FILE}</span>
            <span className="font-mono text-[10.5px] tracking-[0.1em] text-muted">TOML</span>
          </div>
          <pre className="m-0 overflow-x-auto bg-ink-700 px-5 py-[22px] font-mono text-[12.5px] leading-[1.75] text-bone/[0.88]">
            <code>
              {'name = '}
              <S>"pt-br"</S>
              {'\n\ncolumn_suffixes = ['}
              <S>"_id"</S>, <S>"_codigo"</S>, <S>"_cod"</S>, <S>"_key"</S>, <S>"_ref"</S>,{' '}
              <S>"_fk"</S>
              {']\ntable_prefixes  = ['}
              <S>"tb_"</S>, <S>"tbl_"</S>, <S>"sys_"</S>, <S>"cad_"</S>, <S>"mov_"</S>
              {']\n\n[[plural]]                     '}
              <span className="text-muted/70"># opcoes → opcao</span>
              {'\nsuffix = '}
              <S>"oes"</S>
              {'\nsingular = '}
              <S>"ao"</S>
              {'\n\n[[plural]]                     '}
              <span className="text-muted/70"># animais → animal</span>
              {'\nsuffix = '}
              <S>"ais"</S>
              {'\nsingular = '}
              <S>"al"</S>
            </code>
          </pre>
        </Reveal>
      </div>
    </Section>
  )
}
