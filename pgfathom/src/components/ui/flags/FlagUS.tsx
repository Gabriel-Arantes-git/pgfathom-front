import { useId } from 'react'

/** US flag, cropped to a circle. Swap this component for a UK one if the
 *  project would rather mark English with the Union Jack.
 *  See FlagBR for why the clip path id is generated per instance. */
export function FlagUS({ size = 20 }: { size?: number }) {
  const clip = useId()

  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      role="img"
      aria-hidden="true"
      className="block shrink-0"
    >
      <defs>
        <clipPath id={clip}>
          <circle cx="12" cy="12" r="12" />
        </clipPath>
      </defs>
      <g clipPath={`url(#${clip})`}>
        <rect width="24" height="24" fill="#f6f1ea" />
        {[0, 2, 4, 6, 8, 10, 12].map((y) => (
          <rect key={y} y={y * 1.846} width="24" height="1.846" fill="#b22234" />
        ))}
        <rect width="11" height="9.23" fill="#3c3b6e" />
        {[1.6, 4.4, 7.2].map((cy) =>
          [1.4, 3.6, 5.8, 8, 10.2].map((cx) => (
            <circle key={`${cx}-${cy}`} cx={cx} cy={cy} r="0.62" fill="#f6f1ea" />
          )),
        )}
      </g>
    </svg>
  )
}
