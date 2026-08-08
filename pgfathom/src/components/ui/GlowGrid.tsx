import { useEffect, useRef } from 'react'
import type { CSSProperties, ReactNode } from 'react'
import { gsap } from 'gsap'
import './glow-grid.css'

/**
 * Adapted from React Bits' MagicBento, cut down hard on purpose.
 *
 * The original ships particles, 3D tilt, cursor magnetism and a click ripple.
 * All of that is off here: on a page whose argument is restraint, a card that
 * tilts and throws sparks reads as a different product. What survives is the
 * part that adds depth without motion sickness — a border glow that tracks the
 * cursor, plus a soft spotlight over the group — at roughly a third of the
 * reference opacity, in the accent instead of purple.
 */

const ACCENT = '226, 72, 61'
const RADIUS = 320
const MOBILE_BREAKPOINT = 768

function prefersReducedMotion() {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

type GlowGridProps = {
  children: ReactNode
  className?: string
  /** Softer still for dense lists; the default suits 3-up card rows. */
  intensity?: number
  /**
   * The group-wide radial wash. Worth it on a row of cards; switch it off
   * where a big glow would sit over body copy, leaving just the border trace.
   */
  spotlight?: boolean
  style?: CSSProperties
}

export function GlowGrid({
  children,
  className = '',
  intensity = 1,
  spotlight: useSpotlight = true,
  style,
}: GlowGridProps) {
  const gridRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const grid = gridRef.current
    if (!grid) return
    if (prefersReducedMotion() || window.innerWidth <= MOBILE_BREAKPOINT) return

    const spotlight = document.createElement('div')
    spotlight.className = 'glow-spotlight'
    spotlight.style.background = `radial-gradient(circle, rgba(${ACCENT},${0.1 * intensity}) 0%, rgba(${ACCENT},${0.05 * intensity}) 18%, rgba(${ACCENT},${0.02 * intensity}) 38%, transparent 68%)`
    if (useSpotlight) document.body.appendChild(spotlight)

    const proximity = RADIUS * 0.5
    const fade = RADIUS * 0.75

    const onMove = (e: MouseEvent) => {
      const rect = grid.getBoundingClientRect()
      const inside =
        e.clientX >= rect.left &&
        e.clientX <= rect.right &&
        e.clientY >= rect.top &&
        e.clientY <= rect.bottom

      const cards = grid.querySelectorAll<HTMLElement>('.glow-card')

      if (!inside) {
        gsap.to(spotlight, { opacity: 0, duration: 0.3, ease: 'power2.out' })
        cards.forEach((card) => card.style.setProperty('--glow-intensity', '0'))
        return
      }

      let nearest = Infinity

      cards.forEach((card) => {
        const box = card.getBoundingClientRect()
        const cx = box.left + box.width / 2
        const cy = box.top + box.height / 2
        const distance = Math.max(
          0,
          Math.hypot(e.clientX - cx, e.clientY - cy) - Math.max(box.width, box.height) / 2,
        )
        nearest = Math.min(nearest, distance)

        let glow = 0
        if (distance <= proximity) glow = 1
        else if (distance <= fade) glow = (fade - distance) / (fade - proximity)

        card.style.setProperty('--glow-x', `${((e.clientX - box.left) / box.width) * 100}%`)
        card.style.setProperty('--glow-y', `${((e.clientY - box.top) / box.height) * 100}%`)
        card.style.setProperty('--glow-intensity', String(glow * intensity))
        card.style.setProperty('--glow-radius', `${RADIUS}px`)
      })

      gsap.to(spotlight, { left: e.clientX, top: e.clientY, duration: 0.1, ease: 'power2.out' })

      const target =
        nearest <= proximity ? 0.7 : nearest <= fade ? ((fade - nearest) / (fade - proximity)) * 0.7 : 0

      gsap.to(spotlight, {
        opacity: target,
        duration: target > 0 ? 0.2 : 0.5,
        ease: 'power2.out',
      })
    }

    const onLeave = () => {
      grid.querySelectorAll<HTMLElement>('.glow-card').forEach((card) => {
        card.style.setProperty('--glow-intensity', '0')
      })
      gsap.to(spotlight, { opacity: 0, duration: 0.3, ease: 'power2.out' })
    }

    document.addEventListener('mousemove', onMove)
    document.addEventListener('mouseleave', onLeave)

    return () => {
      document.removeEventListener('mousemove', onMove)
      document.removeEventListener('mouseleave', onLeave)
      gsap.killTweensOf(spotlight)
      spotlight.remove()
    }
  }, [intensity, useSpotlight])

  return (
    <div ref={gridRef} className={className} style={style}>
      {children}
    </div>
  )
}

/**
 * Marks a child of <GlowGrid> as glow-tracked.
 *
 * The glow is an `::after` pinned to the card's box, so the card has to be a
 * containing block — every caller must pass a positioning class, `relative`
 * for the usual case. That is deliberately the caller's job: setting it here,
 * or in the stylesheet, means an element that positions itself ends up with
 * two competing `position` utilities and the winner comes down to Tailwind's
 * internal emission order rather than anything visible at the call site.
 */
export function glowCard(className = '') {
  return `glow-card ${className}`
}
