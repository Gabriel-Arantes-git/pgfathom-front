import { useId } from 'react'

/** Brazilian flag, cropped to a circle. Inline SVG because flag emoji do not
 *  render on Windows — the OS falls back to the letters "BR".
 *  The clip path id is generated per instance: the switch renders in the
 *  navbar, the mobile menu and the footer at once, and duplicate DOM ids make
 *  every copy after the first resolve against the wrong element. */
export function FlagBR({ size = 20 }: { size?: number }) {
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
        <rect width="24" height="24" fill="#009b3a" />
        <path d="M12 3.6 22.2 12 12 20.4 1.8 12Z" fill="#fedf00" />
        <circle cx="12" cy="12" r="4.4" fill="#002776" />
        <path
          d="M7.9 10.4a10 10 0 0 1 8.3 2.5"
          fill="none"
          stroke="#fff"
          strokeWidth="1.15"
          strokeLinecap="round"
        />
      </g>
    </svg>
  )
}
