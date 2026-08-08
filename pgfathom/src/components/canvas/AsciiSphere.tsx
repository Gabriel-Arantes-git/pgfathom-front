import { useEffect, useRef } from 'react'

const CHARS = '░▒▓█▀▄▌▐│─┤├┴┬╭╮╰╯'

/** Bone #f6f1ea (near) → accent #e2483d (far), mixed by depth. */
const NEAR = [246, 241, 234] as const
const FAR = [226, 72, 61] as const

type Props = {
  /** Multiplier on the overall opacity. */
  intensity?: number
  className?: string
}

/**
 * Ported from the `asthetic` reference and retinted for the dark theme: points
 * are coloured by depth instead of a flat black, and a sonar ping expands out
 * of the centre every four seconds — the metaphor the project is named for.
 *
 * Blank-canvas guard. An earlier version tracked visibility with two writers —
 * an IntersectionObserver and a `visibilitychange` handler — into one flag, and
 * skipped painting entirely while off-screen. Whichever fired last won, so the
 * flag could end up false with the canvas on screen, and because nothing
 * repainted, the canvas stayed empty until the next resize. It also never
 * re-measured on return, so a backing store dropped while off-screen was never
 * restored. Now the two signals are separate inputs to one derived value, the
 * canvas is re-measured and repainted synchronously the moment it becomes
 * visible again, and a ResizeObserver catches layout changes the window
 * `resize` event misses.
 */
export function AsciiSphere({ intensity = 1, className }: Props) {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    let time = 0
    let frame = 0
    let inView = true
    let pageVisible = !document.hidden

    const measure = () => {
      const dpr = window.devicePixelRatio || 1
      const rect = canvas.getBoundingClientRect()
      const w = Math.max(1, Math.round(rect.width * dpr))
      const h = Math.max(1, Math.round(rect.height * dpr))
      // Assigning width/height also clears the canvas, so only do it on change.
      if (canvas.width !== w || canvas.height !== h) {
        canvas.width = w
        canvas.height = h
      }
      ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
    }

    const render = () => {
      const rect = canvas.getBoundingClientRect()
      if (rect.width === 0 || rect.height === 0) return

      ctx.clearRect(0, 0, rect.width, rect.height)

      const cx = rect.width / 2
      const cy = rect.height / 2
      const radius = Math.min(rect.width, rect.height) * 0.525

      ctx.font = '12px "Geist Mono", ui-monospace, monospace'
      ctx.textAlign = 'center'
      ctx.textBaseline = 'middle'

      const points: { x: number; y: number; z: number; char: string }[] = []

      for (let phi = 0; phi < Math.PI * 2; phi += 0.15) {
        for (let theta = 0; theta < Math.PI; theta += 0.15) {
          const x = Math.sin(theta) * Math.cos(phi + time * 0.5)
          const y = Math.sin(theta) * Math.sin(phi + time * 0.5)
          const z = Math.cos(theta)

          const rotY = time * 0.3
          const nx = x * Math.cos(rotY) - z * Math.sin(rotY)
          const nz = x * Math.sin(rotY) + z * Math.cos(rotY)

          const rotX = time * 0.2
          const ny = y * Math.cos(rotX) - nz * Math.sin(rotX)
          const fz = y * Math.sin(rotX) + nz * Math.cos(rotX)

          const depth = (fz + 1) / 2
          points.push({
            x: cx + nx * radius,
            y: cy + ny * radius,
            z: fz,
            char: CHARS[Math.floor(depth * (CHARS.length - 1))],
          })
        }
      }

      points.sort((a, b) => a.z - b.z)

      for (const p of points) {
        const depth = (p.z + 1) / 2
        const alpha = (0.14 + depth * 0.42) * intensity
        const r = Math.round(FAR[0] + (NEAR[0] - FAR[0]) * depth)
        const g = Math.round(FAR[1] + (NEAR[1] - FAR[1]) * depth)
        const b = Math.round(FAR[2] + (NEAR[2] - FAR[2]) * depth)
        ctx.fillStyle = `rgba(${r}, ${g}, ${b}, ${alpha})`
        ctx.fillText(p.char, p.x, p.y)
      }

      // Sonar ping: one ring every 4s, expanding past the sphere and fading.
      const ping = (time % 4) / 4
      ctx.strokeStyle = `rgba(226, 72, 61, ${(1 - ping) * 0.24 * intensity})`
      ctx.lineWidth = 1
      ctx.beginPath()
      ctx.arc(cx, cy, radius * (0.2 + ping * 1.25), 0, Math.PI * 2)
      ctx.stroke()
    }

    const paint = () => {
      measure()
      render()
    }

    const loop = () => {
      if (inView && pageVisible) {
        paint()
        time += 0.02
      }
      frame = requestAnimationFrame(loop)
    }

    paint()

    const onWindowResize = () => paint()
    window.addEventListener('resize', onWindowResize)

    if (reduced) {
      return () => window.removeEventListener('resize', onWindowResize)
    }

    const resizeObserver = new ResizeObserver(() => paint())
    resizeObserver.observe(canvas)

    const io = new IntersectionObserver(([entry]) => {
      const wasHidden = !inView
      inView = entry.isIntersecting
      // Repaint synchronously on return so the canvas is never left blank
      // waiting on the next frame.
      if (inView && wasHidden && pageVisible) paint()
    })
    io.observe(canvas)

    const onVisibility = () => {
      const wasHidden = !pageVisible
      pageVisible = !document.hidden
      if (pageVisible && wasHidden && inView) paint()
    }
    document.addEventListener('visibilitychange', onVisibility)

    loop()

    return () => {
      cancelAnimationFrame(frame)
      window.removeEventListener('resize', onWindowResize)
      document.removeEventListener('visibilitychange', onVisibility)
      resizeObserver.disconnect()
      io.disconnect()
    }
  }, [intensity])

  return <canvas ref={canvasRef} aria-hidden="true" className={className} />
}
