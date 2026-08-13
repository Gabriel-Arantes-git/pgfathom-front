import { useEffect, useRef } from 'react'
import { GRAPH_BOX, GRAPH_EDGES, GRAPH_TABLES } from '../../content/data'
import { useCopy } from '../../i18n/LanguageContext'
import { GlowGrid, glowCard } from '../ui/GlowGrid'

/**
 * The hero's animated element: a schema that stitches itself together.
 *
 * Tables sit still — the page of unrelated boxes the ERD tool draws. Each
 * relationship then proves itself on its own clock: it draws in, locks with a
 * pulse, stays connected for several seconds, dims out and comes back. The
 * hold times are non-harmonic, so no two lines ever move together and the
 * graph never reads as one looping animation.
 *
 * The cards used to drift on a slow sine. It never stopped reading as stepping,
 * for a reason worth recording: the drift moved well under a pixel per frame,
 * which is exactly the regime where any hitch in frame delivery is visible as a
 * jump, and no amount of smoothing fixes the underlying sampling. Static cards
 * sidestep it entirely, and cost nothing — depth still reads from scale, blur
 * and opacity, and the motion that carries meaning is in the edges.
 *
 * Because positions are constant, they are computed once per resize rather
 * than per frame, and the loop only touches the five edges.
 */

const DRAW = 0.7
const ORPHANS_PER_EDGE = 3
/** How far past the border the line stops — 1px lands it on the border line
 *  itself, so the wire visibly plugs into the card. */
const EDGE_GAP = 1

/**
 * Depth ramp. Measured on screen rather than guessed: an earlier ramp spread
 * the wired cluster across z 0.05–0.80 and faded opacity by 0.5z, which left
 * more than half the tables as unreadable smudges. The wired nine now occupy
 * 0.05–0.52 and only the ambient scenery goes deeper, so every table that
 * carries meaning stays legible while the layering still reads.
 */
const scaleOf = (z: number) => 1.06 - z * 0.3
const blurOf = (z: number) => z * 2.8
const opacityOf = (z: number) => 1 - z * 0.45
const centreOf = (table: { x: number; y: number }) => ({
  x: table.x * GRAPH_BOX.width,
  y: table.y * GRAPH_BOX.height,
})
const transformOf = (table: { x: number; y: number; z: number }) => {
  const { x, y } = centreOf(table)
  return `translate3d(${x}px, ${y}px, 0) translate(-50%, -50%) scale(${scaleOf(table.z)})`
}

const clamp01 = (n: number) => Math.min(1, Math.max(0, n))
/** Cubic ease-out — the line decelerates as it lands instead of stopping dead. */
const easeOut = (t: number) => 1 - Math.pow(1 - t, 3)

function bezier(t: number, x0: number, y0: number, cx: number, cy: number, x1: number, y1: number) {
  const u = 1 - t
  return {
    x: u * u * x0 + 2 * u * t * cx + t * t * x1,
    y: u * u * y0 + 2 * u * t * cy + t * t * y1,
  }
}

/**
 * Walks from a card's centre toward `(tx, ty)` and stops on the card's border,
 * so an edge terminates against the box instead of vanishing under it.
 */
function exitPoint(cx: number, cy: number, hw: number, hh: number, tx: number, ty: number) {
  const dx = tx - cx
  const dy = ty - cy
  if (dx === 0 && dy === 0) return { x: cx, y: cy }

  const sx = dx === 0 ? Infinity : hw / Math.abs(dx)
  const sy = dy === 0 ? Infinity : hh / Math.abs(dy)
  const scale = Math.min(sx, sy)
  const length = Math.hypot(dx, dy)

  return {
    x: cx + dx * scale + (dx / length) * EDGE_GAP,
    y: cy + dy * scale + (dy / length) * EDGE_GAP,
  }
}

export function SchemaGraph({ className }: { className?: string }) {
  const rootRef = useRef<HTMLDivElement>(null)
  const t = useCopy()

  useEffect(() => {
    const root = rootRef.current
    if (!root) return

    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    type CardEntry = { el: HTMLElement; hw: number; hh: number; x: number; y: number }
    const cards = new Map<string, CardEntry>()

    for (const table of GRAPH_TABLES) {
      const el = root.querySelector<HTMLElement>(`[data-table="${table.id}"]`)
      if (!el) continue
      const { x, y } = centreOf(table)
      cards.set(table.id, { el, hw: 0, hh: 0, x, y })
    }

    /**
     * Positions come from the render, not from here — this only caches each
     * card's half-extents, which the edge maths needs and which depend on the
     * laid-out box rather than on anything we can compute up front.
     */
    const measure = () => {
      for (const table of GRAPH_TABLES) {
        const entry = cards.get(table.id)
        if (!entry) continue
        const scale = scaleOf(table.z)
        entry.hw = (entry.el.offsetWidth * scale) / 2
        entry.hh = (entry.el.offsetHeight * scale) / 2
      }
    }

    const paths = GRAPH_EDGES.map((_, i) =>
      root.querySelector<SVGPathElement>(`[data-edge="${i}"]`),
    )
    const locks = GRAPH_EDGES.map((_, i) =>
      root.querySelector<SVGCircleElement>(`[data-lock="${i}"]`),
    )
    const orphanDots = GRAPH_EDGES.map((_, i) =>
      Array.from({ length: ORPHANS_PER_EDGE }, (_, o) =>
        root.querySelector<SVGCircleElement>(`[data-orphan="${i}-${o}"]`),
      ),
    )

    /**
     * Everything about an edge's curve — its control point, its two endpoints,
     * the fixed point each orphan dot falls from — is a function of card
     * centres and half-extents alone, both of which only change on resize or
     * once the webfont settles. Recomputing that geometry every animation
     * frame, and rewriting the path's `d` attribute (an SVG reparse, the
     * costliest write in this loop) 60 times a second for a curve that never
     * moves, was burning main-thread time the browser needed for scrolling —
     * felt as dropped frames while this graph was on screen. `computeGeometry`
     * does that work once per layout change instead; the per-frame loop below
     * touches only what actually animates: draw progress, fade, the lock's
     * pulse, and how far each orphan has fallen.
     */
    type EdgeGeom = {
      ctrlX: number
      ctrlY: number
      p0: { x: number; y: number }
      p1: { x: number; y: number }
      orphanBase: { x: number; y: number }[]
    }
    const edgeGeom: (EdgeGeom | null)[] = GRAPH_EDGES.map(() => null)

    const computeGeometry = () => {
      GRAPH_EDGES.forEach((edge, i) => {
        const path = paths[i]
        const lock = locks[i]
        const from = cards.get(edge.from)
        const to = cards.get(edge.to)
        if (!path || !from || !to) {
          edgeGeom[i] = null
          return
        }

        // Bow the line perpendicular to itself so parallel edges stay legible.
        const vx = to.x - from.x
        const vy = to.y - from.y
        const length = Math.hypot(vx, vy) || 1
        const bow = 20 * (i % 2 === 0 ? 1 : -1)
        const ctrlX = (from.x + to.x) / 2 + (-vy / length) * bow
        const ctrlY = (from.y + to.y) / 2 + (vx / length) * bow

        // Aim at the control point, not the far centre: that is the direction
        // the curve actually leaves in.
        const p0 = exitPoint(from.x, from.y, from.hw, from.hh, ctrlX, ctrlY)
        const p1 = exitPoint(to.x, to.y, to.hw, to.hh, ctrlX, ctrlY)

        path.setAttribute('d', `M ${p0.x} ${p0.y} Q ${ctrlX} ${ctrlY} ${p1.x} ${p1.y}`)

        if (lock) {
          lock.setAttribute('cx', String(p1.x))
          lock.setAttribute('cy', String(p1.y))
        }

        const orphanBase = Array.from({ length: ORPHANS_PER_EDGE }, (_, o) => {
          const at = bezier(0.32 + o * 0.18, p0.x, p0.y, ctrlX, ctrlY, p1.x, p1.y)
          orphanDots[i][o]?.setAttribute('cx', String(at.x))
          return at
        })

        edgeGeom[i] = { ctrlX, ctrlY, p0, p1, orphanBase }
      })
    }

    let frame = 0
    let inView = true
    let pageVisible = !document.hidden
    const start = performance.now()

    const drawEdges = (elapsed: number) => {
      GRAPH_EDGES.forEach((edge, i) => {
        const path = paths[i]
        const lock = locks[i]
        const geom = edgeGeom[i]
        if (!path || !geom) return

        // Each edge on its own clock: draw → hold → fade → pause → repeat.
        const period = DRAW + edge.hold + edge.fade + edge.pause
        const since = elapsed - edge.delay
        const phase = reduced ? DRAW + edge.hold * 0.5 : since < 0 ? -1 : since % period
        // Seconds spent connected in this cycle; negative while still drawing.
        const settled = phase - DRAW

        let progress = 1
        let alpha = 0

        if (phase < 0) {
          progress = 0
        } else if (settled < 0) {
          progress = easeOut(phase / DRAW)
          alpha = 1
        } else if (settled < edge.hold) {
          alpha = 1
        } else if (settled < edge.hold + edge.fade) {
          alpha = 1 - (settled - edge.hold) / edge.fade
        }

        path.style.strokeDashoffset = String(1 - progress)
        path.style.opacity = String(alpha)

        const connected = settled >= 0 && alpha > 0

        if (lock) {
          const pulse = clamp01(settled / 0.55)
          const showing = !reduced && connected && pulse < 1
          lock.setAttribute('r', String(3 + pulse * 15))
          lock.style.opacity = showing ? String((1 - pulse) * 0.75) : '0'
        }

        // Broken relationships shed rows out of the line they just proved.
        for (let o = 0; o < ORPHANS_PER_EDGE; o++) {
          const dot = orphanDots[i][o]
          if (!dot) continue

          if (reduced || edge.verdict !== 'broken' || !connected) {
            dot.style.opacity = '0'
            continue
          }

          const local = (settled - o * 0.6) / 2.6
          if (local < 0) {
            dot.style.opacity = '0'
            continue
          }

          const k = local % 1
          dot.setAttribute('cy', String(geom.orphanBase[o].y + k * k * 70))
          dot.style.opacity = String((1 - k) * 0.9 * alpha)
        }
      })
    }

    const loop = (now: number) => {
      if (inView && pageVisible) drawEdges((now - start) / 1000)
      frame = requestAnimationFrame(loop)
    }

    measure()
    computeGeometry()
    drawEdges(0)

    // Cards are measured before Geist Mono arrives; the fallback font gives a
    // different height, so edge endpoints computed from the stale box land
    // inside or short of the card. Re-measure once the webfont settles.
    let disposed = false
    document.fonts?.ready.then(() => {
      if (!disposed) {
        measure()
        computeGeometry()
        drawEdges((performance.now() - start) / 1000)
      }
    })

    const resizeObserver = new ResizeObserver(() => {
      measure()
      computeGeometry()
      drawEdges((performance.now() - start) / 1000)
    })
    resizeObserver.observe(root)

    if (reduced) {
      return () => {
        disposed = true
        resizeObserver.disconnect()
      }
    }

    const io = new IntersectionObserver(([entry]) => {
      inView = entry.isIntersecting
    })
    io.observe(root)

    const onVisibility = () => {
      pageVisible = !document.hidden
    }
    document.addEventListener('visibilitychange', onVisibility)

    frame = requestAnimationFrame(loop)

    return () => {
      disposed = true
      cancelAnimationFrame(frame)
      document.removeEventListener('visibilitychange', onVisibility)
      resizeObserver.disconnect()
      io.disconnect()
    }
  }, [t])

  return (
    <GlowGrid spotlight={false} intensity={0.8} className={`relative ${className ?? ''}`}>
      <div
        ref={rootRef}
        aria-hidden="true"
        className="relative"
        style={{ width: GRAPH_BOX.width, height: GRAPH_BOX.height }}
      >
        <svg className="absolute inset-0 h-full w-full">
          {GRAPH_EDGES.map((edge, i) => (
            <path
              key={`${edge.from}-${edge.to}`}
              data-edge={i}
              pathLength={1}
              fill="none"
              stroke={edge.verdict === 'broken' ? '#e2483d' : 'rgba(246,241,234,0.5)'}
              strokeWidth={1.25}
              strokeDasharray={1}
              strokeDashoffset={1}
              strokeLinecap="round"
              style={{ opacity: 0 }}
            />
          ))}

          {GRAPH_EDGES.map((edge, i) => (
            <circle
              key={`lock-${edge.from}-${edge.to}`}
              data-lock={i}
              fill="none"
              stroke={edge.verdict === 'broken' ? '#e2483d' : '#f6f1ea'}
              strokeWidth={1}
              r={3}
              style={{ opacity: 0 }}
            />
          ))}

          {GRAPH_EDGES.flatMap((_, i) =>
            Array.from({ length: ORPHANS_PER_EDGE }, (_, o) => (
              <circle
                key={`orphan-${i}-${o}`}
                data-orphan={`${i}-${o}`}
                r={1.7}
                fill="#e2483d"
                style={{ opacity: 0 }}
              />
            )),
          )}
        </svg>

        {GRAPH_TABLES.map((table) => {
          const label = t.graph[table.id]
          return (
            <div
              key={table.id}
              data-table={table.id}
              className={glowCard(
                'absolute top-0 left-0 w-[152px] overflow-hidden rounded-[5px] border border-hair bg-ink-700',
              )}
              // Position and depth are pure functions of the table's data, so
              // they belong in the render: the first paint is already correct
              // and no effect can leave a card stranded at the origin.
              style={{
                transform: transformOf(table),
                filter: table.z > 0.06 ? `blur(${blurOf(table.z).toFixed(2)}px)` : undefined,
                opacity: opacityOf(table.z),
                zIndex: Math.round((1 - table.z) * 100),
                boxShadow: '0 18px 40px -22px rgba(0,0,0,0.9)',
              }}
            >
              <div className="border-b border-hair bg-ink-600 px-2.5 py-1.5">
                <span className="font-mono text-[10.5px] tracking-[0.04em] text-accent">
                  {label?.name ?? table.id}
                </span>
              </div>
              <div className="px-2.5 py-1.5">
                {table.columns.map((column, index) => (
                  <div
                    key={index}
                    className="flex items-baseline justify-between gap-3 py-[1.5px] font-mono text-[9px] leading-[1.5]"
                  >
                    <span className={column.fk ? 'text-accent' : 'text-bone/90'}>
                      {label?.columns[index] ?? ''}
                    </span>
                    <span className="text-muted/75">{column.type}</span>
                  </div>
                ))}
              </div>
            </div>
          )
        })}
      </div>
    </GlowGrid>
  )
}
