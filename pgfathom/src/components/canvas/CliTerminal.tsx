import { useEffect, useMemo, useRef, useState } from 'react'
import type { CSSProperties, KeyboardEvent } from 'react'
import { useCopy } from '../../i18n/LanguageContext'
import { useTypewriter } from '../../hooks/useTypewriter'
import { useLineReveal } from '../../hooks/useLineReveal'
import { GlowGrid, glowCard } from '../ui/GlowGrid'
import {
  BANNER_DESCENT,
  BANNER_MARK,
  CLI_COMMANDS,
  DEFAULT_CLI_COMMAND,
  classifyLine,
  matchesCliQuery,
  normalizeCliInput,
  resolveCliCommand,
  unknownCommandOutput,
} from '../../content/cli'
import type { CliCommand, CliCommandKey, CliTone } from '../../content/cli'

type Screen = {
  key: CliCommandKey | 'unknown'
  args: string
  body: string
}

const TONE_CLASS: Record<CliTone, string> = {
  bone: 'text-bone',
  dim: 'text-muted',
  bold: 'text-bone font-medium',
  accent: 'text-accent font-medium',
  confirm: 'text-confirm font-medium',
  warn: 'text-accent/90',
}

function screenFromCommand(cmd: CliCommand): Screen {
  return { key: cmd.key, args: cmd.args, body: cmd.body }
}

function promptFor(args: string): string {
  return args ? `pgfathom ${args}` : 'pgfathom'
}

type CliTerminalProps = {
  className?: string
  style?: CSSProperties
}

/**
 * The Hero's "pgfathom — try it" card. Types out a command, prints its real
 * output (see `content/cli.ts` for provenance), and lets a visitor pick a
 * different one from a combobox-style autocomplete: click a suggestion, or
 * type and press Enter/Tab. Four real commands — `discover`, `audit`,
 * `version`, and the root command (aliased to `help`/`status`) — cover the
 * product's actual surface rather than a synthetic fifth one.
 */
export function CliTerminal({ className, style }: CliTerminalProps) {
  const t = useCopy()
  const inputRef = useRef<HTMLInputElement>(null)
  const scrollRef = useRef<HTMLDivElement>(null)

  const [screen, setScreen] = useState<Screen>(() => screenFromCommand(DEFAULT_CLI_COMMAND))
  const [query, setQuery] = useState('')
  const [open, setOpen] = useState(false)
  const [highlight, setHighlight] = useState(-1)

  const { value: typed, done } = useTypewriter(screen.args)

  const lines = useMemo(() => screen.body.split('\n'), [screen.body])
  const visibleLines = useLineReveal(screen.body, lines.length, done)

  // New command: jump the pane back to the top before its output streams in,
  // rather than leaving it wherever the previous, possibly taller, output
  // left the scroll position.
  useEffect(() => {
    if (scrollRef.current) scrollRef.current.scrollTop = 0
  }, [screen.body])

  // Follow the stream: as each line lands, keep the newest one in view —
  // the same reason a real terminal scrolls to the bottom as output prints.
  useEffect(() => {
    const el = scrollRef.current
    if (el) el.scrollTop = el.scrollHeight
  }, [visibleLines])

  const filtered = useMemo(() => CLI_COMMANDS.filter((cmd) => matchesCliQuery(cmd, query)), [query])

  const runCommand = (cmd: CliCommand) => {
    setScreen(screenFromCommand(cmd))
    setQuery('')
    setOpen(false)
    setHighlight(-1)
  }

  const runFreeform = () => {
    if (highlight >= 0 && filtered[highlight]) {
      runCommand(filtered[highlight])
      return
    }
    const match = resolveCliCommand(query)
    if (match) {
      runCommand(match)
      return
    }
    const input = normalizeCliInput(query)
    if (input === '') {
      runCommand(CLI_COMMANDS[0])
      return
    }
    setScreen({ key: 'unknown', args: input, body: unknownCommandOutput(input) })
    setQuery('')
    setOpen(false)
    setHighlight(-1)
  }

  const onKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault()
        setOpen(true)
        setHighlight((h) => Math.min(h + 1, filtered.length - 1))
        break
      case 'ArrowUp':
        e.preventDefault()
        setHighlight((h) => Math.max(h - 1, 0))
        break
      case 'Tab':
        if (open && highlight >= 0 && filtered[highlight]) {
          e.preventDefault()
          setQuery(filtered[highlight].args)
        }
        break
      case 'Enter':
        e.preventDefault()
        runFreeform()
        break
      case 'Escape':
        setOpen(false)
        setHighlight(-1)
        break
      default:
        break
    }
  }

  const activeId = highlight >= 0 && filtered[highlight] ? `cli-opt-${filtered[highlight].key}` : undefined

  return (
    <GlowGrid className={className} intensity={0.75} style={style}>
      <div
        className={glowCard(
          'relative overflow-hidden rounded-[5px] border border-hair bg-ink-700 shadow-[0_32px_64px_-24px_rgba(0,0,0,0.7)]',
        )}
      >
        <div className="flex h-10 items-center justify-between gap-4 border-b border-hair bg-ink-600 px-3.5">
          <div className="flex items-center gap-[7px]">
            <span className="h-[9px] w-[9px] rounded-full bg-hair" />
            <span className="h-[9px] w-[9px] rounded-full bg-hair" />
            <span className="h-[9px] w-[9px] rounded-full bg-hair" />
          </div>
          <span className="hidden font-mono text-[11px] text-muted md:inline">{t.terminalCaption}</span>
          <span className="flex items-center gap-1.5 font-mono text-[10.5px] tracking-[0.08em] text-muted uppercase">
            <span aria-hidden="true" className="h-[5px] w-[5px] rounded-full bg-accent" />
            read-only
          </span>
        </div>

        <div className="px-6 pt-6 pb-5 font-mono text-[13px] leading-[1.7]">
          <div className="flex items-start gap-2.5">
            <span className="text-accent">$</span>
            <span className="whitespace-pre text-bone">
              {typed ? promptFor(typed) : 'pgfathom'}
              <span className="caret text-accent">▌</span>
            </span>
          </div>

          {/* Fixed, responsive height: switching commands never resizes this
              card. Shorter output just leaves space below it, the way a real
              terminal window doesn't shrink when you clear the screen; taller
              output scrolls inside this pane instead of growing the card and
              shoving the rest of the page down — that resize mid-scroll was
              what read as the page lagging. */}
          <div ref={scrollRef} className="cli-scroll mt-5 h-[280px] overflow-auto sm:h-[360px] lg:h-[440px]">
            <div className="min-w-[600px]">
              {done && screen.key === 'root' && (
                <div aria-hidden="true" className="cli-line-in mb-1">
                  {BANNER_MARK.map((line, i) => (
                    <div
                      key={i}
                      className="whitespace-pre text-[11px] leading-[1.15]"
                      style={{ color: BANNER_DESCENT[Math.min(i, BANNER_DESCENT.length - 1)] }}
                    >
                      {line}
                    </div>
                  ))}
                  <div className="mt-2">
                    <span className="font-semibold" style={{ color: BANNER_DESCENT[0] }}>
                      pgfathom
                    </span>
                    <span className="text-muted">{'  relationships your database has but never declared'}</span>
                  </div>
                </div>
              )}

              {lines.slice(0, visibleLines).map((line, i) => (
                <div key={i} className={`cli-line-in whitespace-pre ${TONE_CLASS[classifyLine(line)]}`}>
                {line.length ? line : ' '}
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="relative border-t border-hair">
          {open && filtered.length > 0 && (
            <ul
              id="cli-terminal-listbox"
              role="listbox"
              className="absolute right-3 bottom-full left-3 z-10 mb-1 overflow-hidden rounded-[5px] border border-hair bg-ink-600 shadow-[0_16px_40px_-12px_rgba(0,0,0,0.6)]"
            >
              {filtered.map((cmd, i) => (
                <li
                  key={cmd.key}
                  id={`cli-opt-${cmd.key}`}
                  role="option"
                  aria-selected={i === highlight}
                  onMouseDown={(e) => {
                    e.preventDefault()
                    runCommand(cmd)
                  }}
                  onMouseEnter={() => setHighlight(i)}
                  className={`flex cursor-pointer items-center justify-between gap-4 px-4 py-2.5 font-mono text-[12.5px] ${
                    i === highlight ? 'bg-ink-500 text-bone' : 'text-muted'
                  }`}
                >
                  <span className={i === highlight ? 'text-bone' : 'text-muted'}>{promptFor(cmd.args)}</span>
                  <span className="text-[11px] text-muted">{cmd.hint}</span>
                </li>
              ))}
            </ul>
          )}

          <div
            className="flex items-center gap-2.5 px-6 py-3"
            onClick={() => {
              setOpen(true)
              inputRef.current?.focus()
            }}
          >
            <span className="font-mono text-[13px] text-accent">$</span>
            <span className="shrink-0 font-mono text-[13px] text-muted select-none">pgfathom</span>
            <div className="relative min-w-0 flex-1">
              <input
                ref={inputRef}
                type="text"
                value={query}
                onChange={(e) => {
                  setQuery(e.target.value)
                  setOpen(true)
                  setHighlight(-1)
                }}
                onFocus={() => setOpen(true)}
                onClick={() => setOpen(true)}
                onBlur={() => setOpen(false)}
                onKeyDown={onKeyDown}
                placeholder={t.cliPlaceholder}
                aria-label={t.cliInputAria}
                role="combobox"
                aria-expanded={open}
                aria-controls="cli-terminal-listbox"
                aria-autocomplete="list"
                aria-activedescendant={activeId}
                className="w-full bg-transparent font-mono text-[13px] text-bone caret-transparent outline-none placeholder:text-muted"
              />
              {/* The output pane's prompt caret is a plain inline `▌`, easy
                  because it sits after static typed text. An `<input>` has no
                  such seam to drop a glyph into, so this one is an absolutely
                  positioned overlay instead, parked `query.length` characters
                  in — safe only because the font is monospace, so `ch` lines
                  up with the native caret position it replaces
                  (`caret-transparent` above hides that native one). */}
              <span
                aria-hidden="true"
                className="caret pointer-events-none absolute top-1/2 -translate-y-1/2 text-[13px] text-accent"
                style={{ left: `${query.length}ch` }}
              >
                ▌
              </span>
            </div>
            {query === '' && (
              <span className="hidden shrink-0 font-mono text-[10.5px] text-muted sm:inline">
                {t.cliEmptyHint}
              </span>
            )}
            <button
              type="button"
              onClick={runFreeform}
              aria-label={t.cliRunAria}
              className="flex h-7 w-7 shrink-0 items-center justify-center rounded-[3px] font-mono text-muted transition-colors duration-200 hover:text-bone"
            >
              ↵
            </button>
          </div>
        </div>
      </div>
    </GlowGrid>
  )
}
