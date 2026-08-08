import { useLanguage } from '../../i18n/LanguageContext'
import { LanguageSwitch } from '../ui/LanguageSwitch'
import { GithubIcon } from '../ui/icons'
import { DESIGN_DOC, LICENSE_DOC, REPO, ROADMAP_DOC } from '../../content/data'
import mark from '../../assets/icons/pgfathom-logo.png'

export function Footer() {
  const { t, lang } = useLanguage()

  const links = [
    { href: REPO, label: 'GitHub' },
    { href: DESIGN_DOC[lang], label: t.linkDesign },
    { href: ROADMAP_DOC, label: t.linkRoadmap },
    { href: LICENSE_DOC, label: 'Apache-2.0' },
  ]

  return (
    <footer className="mx-auto max-w-[1140px] px-7 pt-11 pb-14">
      <div className="flex flex-wrap items-start justify-between gap-8">
        <div className="flex items-center gap-3.5">
          <img src={mark} alt="" aria-hidden="true" className="h-[34px] w-[34px] rounded" />
          <div>
            <p className="m-0 font-mono text-[13.5px] text-bone">pgfathom</p>
            <p className="mt-1.5 mb-0 text-[13px] text-muted">{t.tagline}</p>
          </div>
        </div>

        <nav aria-label="Footer" className="flex flex-wrap items-center gap-6">
          {links.map((link) => (
            <a
              key={link.label}
              href={link.href}
              target="_blank"
              rel="noopener noreferrer"
              className="text-[13px] text-muted transition-colors duration-200 hover:text-bone"
            >
              {link.label}
            </a>
          ))}
        </nav>
      </div>

      <div className="mt-9 flex flex-wrap items-center justify-between gap-4 border-t border-hair pt-6">
        <span className="font-mono text-[11.5px] text-muted/70">{t.footLeft}</span>

        <div className="flex items-center gap-4">
          <span className="hidden font-mono text-[11.5px] text-muted/70 lg:inline">
            {t.footRight}
          </span>
          <LanguageSwitch size={18} />
          <a
            href={REPO}
            target="_blank"
            rel="noopener noreferrer"
            aria-label="GitHub"
            className="text-muted transition-all duration-200 hover:-translate-y-0.5 hover:text-accent"
          >
            <GithubIcon size={17} />
          </a>
        </div>
      </div>
    </footer>
  )
}
