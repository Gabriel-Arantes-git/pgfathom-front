import type { GraphLabel } from '../../i18n/types'

type Column = { type: string; fk?: boolean }

type Props = {
  label: GraphLabel
  columns: Column[]
  /** Pixels from the top of the nearest positioned ancestor. */
  top: number
  /** Pixels from the right edge, matching how the schema graph itself anchors. */
  right: number
  /** 0 (sharp, opaque) – 1 (deep, faint). Same ramp as SchemaGraph's ambient tables. */
  depth: number
}

/**
 * A lone card trailing below the hero's schema graph, unconnected to any edge.
 * Not part of SchemaGraph / GRAPH_TABLES: that component's coordinate system
 * is pinned to its own 660×520 box, and stretching that box to reach this far
 * down would have re-scaled the whole wired grid. Position is measured, not
 * guessed — see the call site in Hero.tsx for where the numbers come from.
 */
export function TrailingTable({ label, columns, top, right, depth }: Props) {
  return (
    <div
      aria-hidden="true"
      className="pointer-events-none absolute w-[136px] overflow-hidden rounded-[5px] border border-hair bg-ink-700"
      style={{
        top,
        right,
        opacity: 1 - depth * 0.55,
        filter: `blur(${(depth * 2.6).toFixed(2)}px)`,
        boxShadow: '0 14px 30px -18px rgba(0,0,0,0.85)',
      }}
    >
      <div className="border-b border-hair bg-ink-600 px-2.5 py-1.5">
        <span className="font-mono text-[10px] tracking-[0.04em] text-accent">{label.name}</span>
      </div>
      <div className="px-2.5 py-1.5">
        {columns.map((column, index) => (
          <div
            key={index}
            className="flex items-baseline justify-between gap-3 py-[1.5px] font-mono text-[8.5px] leading-[1.5]"
          >
            <span className={column.fk ? 'text-accent' : 'text-bone/85'}>
              {label.columns[index] ?? ''}
            </span>
            <span className="text-muted/70">{column.type}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
