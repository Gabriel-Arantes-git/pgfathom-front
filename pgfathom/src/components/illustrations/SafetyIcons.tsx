/**
 * Line-art safety icons animated with SMIL — no JS, no animation library.
 * `currentColor` throughout so they inherit the accent from the card.
 * The `prefers-reduced-motion` rule in index.css collapses the durations.
 */

type Props = { size?: number }

const wrap = (size: number) => ({
  width: size,
  height: size,
  viewBox: '0 0 40 40',
  fill: 'none' as const,
  stroke: 'currentColor' as const,
  strokeWidth: 1.5,
  strokeLinecap: 'round' as const,
  strokeLinejoin: 'round' as const,
  'aria-hidden': true as const,
})

/** A write attempt bouncing off the shield, forever. */
export function IconReadOnly({ size = 34 }: Props) {
  return (
    <svg {...wrap(size)}>
      <path d="M20 5 31 9.5v9c0 7.5-5.2 11.6-11 13.5C14.2 30.1 9 26 9 18.5v-9Z" opacity="0.55" />
      <path d="M15.5 19.5 19 23l6-7" opacity="0.9" />
      <g>
        <path d="M2 20h5" opacity="0.8">
          <animate
            attributeName="opacity"
            values="0;0.9;0.9;0"
            dur="2.6s"
            keyTimes="0;0.25;0.55;0.75"
            repeatCount="indefinite"
          />
        </path>
        <circle cx="4" cy="20" r="1.4" fill="currentColor" stroke="none">
          <animate
            attributeName="cx"
            values="2;8.4;8.4;2"
            dur="2.6s"
            keyTimes="0;0.4;0.5;0.85"
            calcMode="spline"
            keySplines="0.4 0 0.2 1;0 0 1 1;0.4 0 0.2 1"
            repeatCount="indefinite"
          />
        </circle>
      </g>
    </svg>
  )
}

/** Values fall into memory; only counts come out the other side. */
export function IconInMemory({ size = 34 }: Props) {
  return (
    <svg {...wrap(size)}>
      <rect x="11" y="12" width="18" height="16" rx="2" opacity="0.55" />
      <path d="M11 17h18" opacity="0.3" />
      {[14, 20, 26].map((x, i) => (
        <circle key={x} cx={x} cy="6" r="1.3" fill="currentColor" stroke="none" opacity="0">
          <animate
            attributeName="cy"
            values="4;15"
            dur="2.4s"
            begin={`${i * 0.35}s`}
            repeatCount="indefinite"
          />
          <animate
            attributeName="opacity"
            values="0;0.85;0"
            dur="2.4s"
            begin={`${i * 0.35}s`}
            repeatCount="indefinite"
          />
        </circle>
      ))}
      <path d="M15 32h10" opacity="0.5" />
      <path d="M16.5 23h7" opacity="0.85">
        <animate attributeName="opacity" values="0.35;0.95;0.35" dur="2.4s" repeatCount="indefinite" />
      </path>
    </svg>
  )
}

/** A needle that rises and is stopped short by the timeout mark. */
export function IconBounded({ size = 34 }: Props) {
  return (
    <svg {...wrap(size)}>
      <path d="M7 27a13 13 0 0 1 26 0" opacity="0.55" />
      <path d="M30.5 16.5 33 14" opacity="0.85" stroke="currentColor" />
      <g transform="translate(20 27)">
        <path d="M0 0 L0 -11" opacity="0.95">
          <animateTransform
            attributeName="transform"
            type="rotate"
            values="-62;44;44;-62"
            dur="3s"
            keyTimes="0;0.45;0.6;1"
            calcMode="spline"
            keySplines="0.3 0 0.2 1;0 0 1 1;0.4 0 0.3 1"
            repeatCount="indefinite"
          />
        </path>
      </g>
      <circle cx="20" cy="27" r="1.7" fill="currentColor" stroke="none" />
    </svg>
  )
}

/** Claim and evidence settling into balance. */
export function IconEvidence({ size = 34 }: Props) {
  return (
    <svg {...wrap(size)}>
      <path d="M20 8v22" opacity="0.5" />
      <path d="M12 30h16" opacity="0.5" />
      <g>
        <path d="M8 13h24" opacity="0.85" />
        <path d="M8 13v4" opacity="0.7" />
        <path d="M32 13v4" opacity="0.7" />
        <animateTransform
          attributeName="transform"
          type="rotate"
          values="-9 20 13;7 20 13;0 20 13;0 20 13"
          dur="3.4s"
          keyTimes="0;0.35;0.6;1"
          calcMode="spline"
          keySplines="0.4 0 0.2 1;0.4 0 0.2 1;0 0 1 1"
          repeatCount="indefinite"
        />
      </g>
      <circle cx="20" cy="8" r="1.5" fill="currentColor" stroke="none" opacity="0.8" />
    </svg>
  )
}

/** A coverage grid where the cell it could not read is marked, not hidden. */
export function IconCoverage({ size = 34 }: Props) {
  const cells = [0, 1, 2, 3, 4, 5, 6, 7, 8]
  return (
    <svg {...wrap(size)}>
      {cells.map((i) => {
        const x = 8 + (i % 3) * 8.5
        const y = 8 + Math.floor(i / 3) * 8.5
        const skipped = i === 4
        return (
          <rect
            key={i}
            x={x}
            y={y}
            width="7"
            height="7"
            rx="1"
            opacity={skipped ? 0.35 : 0.5}
            strokeDasharray={skipped ? '2 2' : undefined}
          >
            {!skipped && (
              <animate
                attributeName="opacity"
                values="0.25;0.85;0.25"
                dur="2.8s"
                begin={`${i * 0.12}s`}
                repeatCount="indefinite"
              />
            )}
          </rect>
        )
      })}
      <path d="M17.5 17.5 23 23" opacity="0.7" />
      <path d="M23 17.5 17.5 23" opacity="0.7" />
    </svg>
  )
}

export const SAFETY_ICONS = [
  IconReadOnly,
  IconInMemory,
  IconBounded,
  IconEvidence,
  IconCoverage,
] as const
