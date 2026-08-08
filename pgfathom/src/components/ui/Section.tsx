import type { ReactNode } from 'react'

export function Section({
  id,
  children,
  className = '',
}: {
  id?: string
  children: ReactNode
  className?: string
}) {
  return (
    <section id={id} className={`border-b border-hair ${className}`}>
      <div className="mx-auto max-w-[1140px] px-7 py-[72px] md:py-28">{children}</div>
    </section>
  )
}

export function H2({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <h2
      className={`mt-[22px] mb-0 text-[27px] leading-[1.12] font-medium tracking-[-0.025em] text-balance md:text-[38px] ${className}`}
    >
      {children}
    </h2>
  )
}

export function Lead({ children }: { children: ReactNode }) {
  return (
    <p className="mt-[22px] mb-0 text-[15.5px] leading-[1.7] text-pretty text-muted">{children}</p>
  )
}
