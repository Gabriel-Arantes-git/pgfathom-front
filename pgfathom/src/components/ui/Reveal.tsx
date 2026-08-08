import type { ElementType, ReactNode } from 'react'
import { useReveal } from '../../hooks/useReveal'

type Props = {
  children: ReactNode
  className?: string
  /** ms — staggers siblings inside a grid. */
  delay?: number
  as?: ElementType
}

export function Reveal({ children, className = '', delay = 0, as: Tag = 'div' }: Props) {
  const { ref, shown } = useReveal()

  return (
    <Tag
      ref={ref}
      className={`reveal ${shown ? 'reveal-in' : ''} ${className}`}
      style={delay ? { transitionDelay: `${delay}ms` } : undefined}
    >
      {children}
    </Tag>
  )
}
