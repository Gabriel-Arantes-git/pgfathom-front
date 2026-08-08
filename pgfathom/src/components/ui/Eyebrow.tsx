export function Eyebrow({ children }: { children: React.ReactNode }) {
  return (
    <p className="m-0 flex items-center gap-2.5 font-mono text-[11px] tracking-[0.16em] text-muted uppercase">
      <span aria-hidden="true" className="h-[5px] w-[5px] bg-accent" />
      {children}
    </p>
  )
}
