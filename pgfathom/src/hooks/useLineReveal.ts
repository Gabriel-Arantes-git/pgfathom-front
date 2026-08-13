import { useEffect, useState } from 'react'

type Options = {
  /** ms between each newly revealed line */
  perLine?: number
}

/**
 * Reveals `total` lines one at a time — the CLI terminal's output streaming
 * in, the way a real command's stdout would, rather than the whole block
 * appearing at once. Gated on `active`: the terminal keeps this false while
 * the command line is still typing, so output only starts once the "command"
 * has finished "running".
 *
 * Keyed on `key` (the command's own text) rather than `total` alone, so
 * switching to a different command whose body happens to have the same
 * number of lines still restarts the reveal instead of silently no-op'ing.
 * Under prefers-reduced-motion every line lands immediately.
 */
export function useLineReveal(key: string, total: number, active: boolean, { perLine = 16 }: Options = {}) {
  const [visible, setVisible] = useState(0)

  useEffect(() => {
    setVisible(0)
    if (!active) return

    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      setVisible(total)
      return
    }

    const timers: number[] = []
    for (let i = 1; i <= total; i++) {
      timers.push(window.setTimeout(() => setVisible(i), i * perLine))
    }
    return () => timers.forEach(clearTimeout)
  }, [key, total, active, perLine])

  return visible
}
