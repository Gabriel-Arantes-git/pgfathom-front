/**
 * One animated mark per verdict, in the same line-art language as the safety
 * icons. Each one dramatises what the verdict means rather than decorating it.
 */

type Props = { size?: number }

const wrap = (size: number) => ({
  width: size,
  height: (size * 40) / 96,
  viewBox: '0 0 96 40',
  fill: 'none' as const,
  stroke: 'currentColor' as const,
  strokeWidth: 1.5,
  strokeLinecap: 'round' as const,
  strokeLinejoin: 'round' as const,
  'aria-hidden': true as const,
})

const Table = ({ x }: { x: number }) => (
  <g opacity="0.55">
    <rect x={x} y="8" width="22" height="24" rx="2" />
    <path d={`M${x} 15h22`} />
  </g>
)

/** The link holds, but rows fall out of it — the orphans. */
export function VerdictBroken({ size = 96 }: Props) {
  return (
    <svg {...wrap(size)}>
      <Table x={6} />
      <Table x={68} />
      <path d="M28 20h40" strokeDasharray="3 3" opacity="0.6" />
      {[0, 1, 2].map((i) => (
        <circle key={i} cx="34" cy="20" r="1.6" fill="currentColor" stroke="none" opacity="0">
          <animate
            attributeName="cx"
            values="32;62"
            dur="2.8s"
            begin={`${i * 0.9}s`}
            repeatCount="indefinite"
          />
          <animate
            attributeName="cy"
            values="20;20;36"
            keyTimes="0;0.45;1"
            dur="2.8s"
            begin={`${i * 0.9}s`}
            repeatCount="indefinite"
          />
          <animate
            attributeName="opacity"
            values="0;0.95;0.95;0"
            keyTimes="0;0.15;0.5;1"
            dur="2.8s"
            begin={`${i * 0.9}s`}
            repeatCount="indefinite"
          />
        </circle>
      ))}
    </svg>
  )
}

/** The arc draws itself and locks. */
export function VerdictConfirmed({ size = 96 }: Props) {
  return (
    <svg {...wrap(size)}>
      <Table x={6} />
      <Table x={68} />
      <path d="M28 20q20 -12 40 0" opacity="0.9" strokeDasharray="52" strokeDashoffset="52">
        <animate
          attributeName="stroke-dashoffset"
          values="52;0;0;52"
          keyTimes="0;0.4;0.85;1"
          dur="3.2s"
          repeatCount="indefinite"
        />
      </path>
      <circle cx="48" cy="14" r="2.2" fill="currentColor" stroke="none" opacity="0">
        <animate
          attributeName="opacity"
          values="0;0;1;1;0"
          keyTimes="0;0.4;0.5;0.85;1"
          dur="3.2s"
          repeatCount="indefinite"
        />
      </circle>
    </svg>
  )
}

/** The beam sweeps between candidate targets and never settles. */
export function VerdictWeak({ size = 96 }: Props) {
  return (
    <svg {...wrap(size)}>
      <Table x={6} />
      {[6, 20, 34].map((y, i) => (
        <rect key={y} x="72" y={y - 2} width="18" height="8" rx="1.5" opacity="0.4">
          <animate
            attributeName="opacity"
            values="0.2;0.75;0.2"
            dur="3.6s"
            begin={`${i * 1.2}s`}
            repeatCount="indefinite"
          />
        </rect>
      ))}
      {[8, 22, 36].map((y, i) => (
        <path key={y} d={`M28 20 L72 ${y}`} strokeDasharray="2 3" opacity="0">
          <animate
            attributeName="opacity"
            values="0;0.7;0"
            dur="3.6s"
            begin={`${i * 1.2}s`}
            repeatCount="indefinite"
          />
        </path>
      ))}
    </svg>
  )
}

export const VERDICT_ICONS = {
  BROKEN: VerdictBroken,
  CONFIRMED: VerdictConfirmed,
  WEAK: VerdictWeak,
} as const
