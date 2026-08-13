import { useEffect, useState } from 'react'
import { useCopy } from '../../i18n/LanguageContext'
import { useScrolled, useScrollProgress, useIsScrolling } from '../../hooks/useScrolled'
import { LanguageSwitch } from '../ui/LanguageSwitch'
import { Wordmark } from '../ui/Wordmark'
import { GithubIcon, MenuIcon, CloseIcon } from '../ui/icons'
import { REPO } from '../../content/data'

export function Navbar() {
  const t = useCopy()
  const scrolled = useScrolled(20)
  const progress = useScrollProgress()
  const isScrolling = useIsScrolling()
  const [menuOpen, setMenuOpen] = useState(false)

  const links = [
    { href: '#problem', label: t.navProblem },
    { href: '#verdicts', label: t.navVerdicts },
    { href: '#safety', label: t.navSafety },
    { href: '#pipeline', label: t.navHow },
    { href: '#benchmark', label: t.navBenchmark },
    { href: '#contributors', label: t.navContributors },
  ]

  useEffect(() => {
    document.body.style.overflow = menuOpen ? 'hidden' : ''
    return () => {
      document.body.style.overflow = ''
    }
  }, [menuOpen])

  return (
    <>
      <header
        className={`fixed z-50 transition-all duration-500 ease-[cubic-bezier(.4,0,.2,1)] ${
          scrolled ? 'top-4 right-4 left-4' : 'top-0 right-0 left-0'
        }`}
      >
        <nav
          className={`${isScrolling ? 'glass-scrolling' : 'glass'} mx-auto transition-all duration-500 ease-[cubic-bezier(.4,0,.2,1)] ${
            scrolled
              ? 'max-w-[1120px] rounded-2xl border border-hair shadow-[0_24px_60px_-30px_rgba(0,0,0,0.9)]'
              : 'max-w-[1140px] rounded-none border-b border-hair/70'
          }`}
        >
          <div
            className={`flex items-center justify-between gap-6 px-7 transition-all duration-500 ease-[cubic-bezier(.4,0,.2,1)] ${
              scrolled ? 'h-14' : 'h-[72px]'
            }`}
          >
            <a href="#top" className="flex shrink-0 items-center" aria-label="pgfathom">
              <Wordmark size={scrolled ? 'sm' : 'md'} />
            </a>

            <div className="hidden items-center gap-[26px] md:flex">
              {links.map((link) => (
                <a
                  key={link.href}
                  href={link.href}
                  className="group relative text-[13px] text-muted transition-colors duration-200 hover:text-bone"
                >
                  {link.label}
                  <span className="absolute -bottom-1.5 left-0 h-px w-0 bg-accent transition-all duration-300 group-hover:w-full" />
                </a>
              ))}
            </div>

            <div className="flex shrink-0 items-center gap-2.5">
              <LanguageSwitch size={scrolled ? 16 : 20} />

              <span aria-hidden="true" className="hidden h-5 w-px bg-hair md:block" />

              <a
                href={REPO}
                target="_blank"
                rel="noopener noreferrer"
                className="flex h-8 items-center gap-2 rounded border border-hair px-2.5 text-[12.5px] text-muted transition-colors duration-200 hover:border-hair-hi hover:bg-ink-600 hover:text-bone"
              >
                <GithubIcon size={15} />
                <span className="hidden font-mono lg:inline">GitHub</span>
              </a>

              <button
                type="button"
                onClick={() => setMenuOpen(true)}
                aria-label={t.navMenu}
                className="flex h-8 w-8 cursor-pointer items-center justify-center rounded border border-hair text-muted transition-colors hover:text-bone md:hidden"
              >
                <MenuIcon />
              </button>
            </div>
          </div>

          {/* Scroll progress hairline pinned to the navbar's bottom edge.
              Inset once the bar rounds off, so the line starts inside the
              corner arc instead of running out past it. */}
          <div
            aria-hidden="true"
            className={`transition-[padding] duration-500 ease-[cubic-bezier(.4,0,.2,1)] ${
              scrolled ? 'px-5' : 'px-0'
            }`}
          >
            <div
              className={`h-px origin-left bg-accent transition-transform duration-150 ${
                scrolled ? 'rounded-full' : ''
              }`}
              style={{ transform: `scaleX(${progress})` }}
            />
          </div>
        </nav>
      </header>

      <div
        className={`fixed inset-0 z-60 bg-ink-800 transition-opacity duration-500 md:hidden ${
          menuOpen ? 'pointer-events-auto opacity-100' : 'pointer-events-none opacity-0'
        }`}
      >
        <div className="flex h-full flex-col px-7 pt-6 pb-10">
          <div className="flex items-center justify-between">
            <Wordmark size="md" />
            <button
              type="button"
              onClick={() => setMenuOpen(false)}
              aria-label={t.navClose}
              className="flex h-9 w-9 cursor-pointer items-center justify-center rounded border border-hair text-muted"
            >
              <CloseIcon />
            </button>
          </div>

          <div className="flex flex-1 flex-col justify-center gap-7">
            {links.map((link, i) => (
              <a
                key={link.href}
                href={link.href}
                onClick={() => setMenuOpen(false)}
                className={`text-[34px] leading-none font-medium tracking-[-0.03em] transition-all duration-500 ${
                  menuOpen ? 'translate-y-0 opacity-100' : 'translate-y-4 opacity-0'
                }`}
                style={{ transitionDelay: menuOpen ? `${i * 75}ms` : '0ms' }}
              >
                {link.label}
              </a>
            ))}
          </div>

          <div className="flex flex-col gap-5 border-t border-hair pt-7">
            <LanguageSwitch size={28} withLabels />
            <a
              href={REPO}
              target="_blank"
              rel="noopener noreferrer"
              className="flex h-11 items-center justify-center gap-2.5 rounded border border-hair font-mono text-[13px] text-bone"
            >
              <GithubIcon /> GitHub
            </a>
          </div>
        </div>
      </div>
    </>
  )
}
