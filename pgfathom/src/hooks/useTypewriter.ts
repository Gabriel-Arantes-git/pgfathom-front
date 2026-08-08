import { useEffect, useState } from 'react'

type Options = {
  /** ms before the first character lands */
  start?: number
  /** ms between characters */
  step?: number
  /** ms after the last character before `done` flips */
  settle?: number
}

/**
 * Types `text` out one character at a time. Returns the visible slice and a
 * `done` flag used to fade the report block in after the command finishes.
 * Under prefers-reduced-motion the final state is applied immediately.
 */
export function useTypewriter(text: string, { start = 600, step = 42, settle = 500 }: Options = {}) {
  const [typed, setTyped] = useState(0)
  const [done, setDone] = useState(false)

  useEffect(() => {
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      setTyped(text.length)
      setDone(true)
      return
    }

    const timers: number[] = []
    for (let i = 1; i <= text.length; i++) {
      timers.push(window.setTimeout(() => setTyped(i), start + i * step))
    }
    timers.push(window.setTimeout(() => setDone(true), start + text.length * step + settle))

    return () => timers.forEach(clearTimeout)
  }, [text, start, step, settle])

  return { value: text.slice(0, typed), done }
}
