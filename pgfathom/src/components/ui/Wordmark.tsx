import { LogoMark } from './LogoMark'

type Props = {
  /** Icon + text size — the navbar shrinks this on scroll. */
  size?: 'sm' | 'md'
  className?: string
}

const ICON = { sm: 'h-[22px] w-[22px]', md: 'h-[26px] w-[26px]' } as const
const TEXT = { sm: 'text-[16px]', md: 'text-[19px]' } as const

/**
 * Real text and a real SVG, replacing the PNG wordmark. The PNG baked "pg" in
 * accent and "fathom" in bone into a raster with its own dark background
 * square — same background-mismatch problem as `LogoMark`, plus it couldn't
 * resize without going soft. This reproduces the same two-tone split with an
 * actual span in the site's own type family, so it stays crisp at any size
 * and never carries a background of its own.
 */
export function Wordmark({ size = 'md', className = '' }: Props) {
  return (
    <span className={`flex items-center gap-2 ${className}`}>
      <LogoMark className={`text-accent transition-all duration-500 ${ICON[size]}`} />
      <span
        className={`font-sans font-bold tracking-tight transition-all duration-500 ${TEXT[size]}`}
      >
        <span className="text-accent">pg</span>
        <span className="text-bone">fathom</span>
      </span>
    </span>
  )
}
