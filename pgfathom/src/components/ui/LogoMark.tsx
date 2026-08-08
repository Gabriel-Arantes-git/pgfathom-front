import { useId } from 'react'

/**
 * The diver/sonar mark, inlined from `assets/icons/pgfathom-logo.svg` rather
 * than rendered as an `<img>`. The PNG exports of that mark (used elsewhere,
 * e.g. Footer/FinalCta) are flattened onto an opaque dark square — on the
 * navbar's glass background that square never quite matches the blur/opacity
 * underneath it, so it reads as a floating block. This version has no
 * background at all and takes `currentColor`, so it disappears into whatever
 * it sits on and recolors with a text-color utility like any other icon.
 *
 * Mask id is generated per instance — the same reason the flag components do
 * it — because this renders in the desktop navbar and the mobile overlay at
 * once, and a duplicated `id` resolves every copy after the first against the
 * wrong mask.
 */
export function LogoMark({ className }: { className?: string }) {
  const maskId = useId()

  return (
    <svg viewBox="0 0 100 120" aria-hidden="true" className={className}>
      <mask id={maskId} maskUnits="userSpaceOnUse" x="0" y="0" width="100" height="120">
        <rect x="0" y="0" width="100" height="120" fill="#fff" />
        <g transform="translate(2,2) scale(0.96)" fill="#000">
          <ellipse cx="44" cy="33" rx="3.4" ry="2.5" />
          <ellipse cx="56" cy="33" rx="3.4" ry="2.5" />
        </g>
      </mask>
      <g mask={`url(#${maskId})`} fill="currentColor">
        <g transform="translate(2,2) scale(0.96)">
          <g stroke="currentColor" strokeWidth={4} strokeLinejoin="round" strokeLinecap="round">
            <path transform="rotate(12 29 14)" d="M29 14L2 14Q2 48 29 48Z" />
            <path transform="rotate(-12 71 14)" d="M71 14L98 14Q98 48 71 48Z" />
            <path d="M36 26Q36 10 50 10Q64 10 64 26L64 50Q64 60 57 62L43 62Q36 60 36 50Z" />
          </g>
          <path d="M50 58L50 86" fill="none" stroke="currentColor" strokeWidth={14} strokeLinecap="round" />
        </g>
        <circle cx="50" cy="96" r="2.4" />
        <circle cx="50" cy="103" r="2.4" />
        <circle cx="50" cy="114" r="5.5" />
      </g>
    </svg>
  )
}
