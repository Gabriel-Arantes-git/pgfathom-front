import { useEffect, useState } from 'react'

/** True once the page has scrolled past `threshold` — drives the navbar collapse. */
export function useScrolled(threshold = 20) {
  const [scrolled, setScrolled] = useState(false)

  useEffect(() => {
    // rAF-batched: a raw `scroll` listener can fire many times per rendered
    // frame (mouse wheel deltas, high-poll-rate trackpads), and each call
    // used to run setState synchronously — one React commit per event
    // instead of per frame. Collapsing to at most one commit per frame keeps
    // this cheap regardless of how chatty the input device is.
    let raf = 0
    const compute = () => {
      raf = 0
      setScrolled(window.scrollY > threshold)
    }
    const onScroll = () => {
      if (!raf) raf = requestAnimationFrame(compute)
    }

    compute()
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => {
      window.removeEventListener('scroll', onScroll)
      if (raf) cancelAnimationFrame(raf)
    }
  }, [threshold])

  return scrolled
}

/** Fraction of the document scrolled, 0–1. Drives the navbar progress hairline. */
export function useScrollProgress() {
  const [progress, setProgress] = useState(0)

  useEffect(() => {
    let raf = 0
    const compute = () => {
      raf = 0
      const max = document.documentElement.scrollHeight - window.innerHeight
      setProgress(max > 0 ? Math.min(window.scrollY / max, 1) : 0)
    }
    const onScroll = () => {
      if (!raf) raf = requestAnimationFrame(compute)
    }

    compute()
    window.addEventListener('scroll', onScroll, { passive: true })
    window.addEventListener('resize', onScroll)
    return () => {
      window.removeEventListener('scroll', onScroll)
      window.removeEventListener('resize', onScroll)
      if (raf) cancelAnimationFrame(raf)
    }
  }, [])

  return progress
}

/**
 * True while the page is actively scrolling, false again `idleMs` after the
 * last scroll event.
 *
 * Exists for the navbar's blur toggle: `backdrop-filter` is recomposited
 * every single frame the content behind a fixed element changes, which is
 * every frame of a scroll gesture — the single largest cost in that header.
 * Dropping the blur for the gesture's duration (nobody is looking at it
 * closely mid-scroll) and restoring it once motion settles keeps the frosted
 * look at rest while removing that cost from the part of the interaction
 * that actually needs the frame budget.
 */
export function useIsScrolling(idleMs = 150) {
  const [scrolling, setScrolling] = useState(false)

  useEffect(() => {
    let raf = 0
    let idleTimer = 0

    const onScroll = () => {
      if (!raf) {
        raf = requestAnimationFrame(() => {
          raf = 0
          setScrolling(true)
        })
      }
      window.clearTimeout(idleTimer)
      idleTimer = window.setTimeout(() => setScrolling(false), idleMs)
    }

    window.addEventListener('scroll', onScroll, { passive: true })
    return () => {
      window.removeEventListener('scroll', onScroll)
      if (raf) cancelAnimationFrame(raf)
      window.clearTimeout(idleTimer)
    }
  }, [idleMs])

  return scrolling
}
