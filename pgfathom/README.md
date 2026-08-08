# pgfathom — landing page

Bilingual (EN / pt-BR) landing page for [`pgfathom`](https://github.com/lvcas-dotcom/pgfathom),
the CLI that finds the foreign keys a legacy PostgreSQL schema has but never declared.

```console
$ npm install
$ npm run dev        # http://localhost:5173
$ npm run build      # tsc -b && vite build → dist/
$ npm run lint       # oxlint
```

## Stack

Vite 8 · React 19 (with the React Compiler) · TypeScript · Tailwind v4 · GSAP.

Tailwind v4 has no `tailwind.config.js` — the design tokens are declared with `@theme`
in `src/index.css`. GSAP drives only the cursor glow in `GlowGrid`; every other
animation is CSS or `IntersectionObserver`, so the runtime stays small.

## Content and i18n

Every string lives in `src/i18n/copy.en.ts` and `src/i18n/copy.pt.ts`. The `Copy` type is
derived from the English file, so a key present in one language and missing in the other
**fails the build** rather than shipping a raw key. Language resolves from
`localStorage` → `navigator.language` → English, and the switch is the flag control in
the navbar, mobile menu and footer.

`src/content/data.ts` holds what is deliberately *not* translated: terminal output, SQL,
TOML, CLI flags. Rendering `QUEBRADO` where the CLI prints `BROKEN` would show output the
tool does not produce, and the page's whole argument is that it does not misrepresent it.

## Keeping it honest

The page mirrors the upstream README, including its restraint. Three things to preserve
when editing:

- **No invented numbers.** `BENCHMARK` in `src/content/data.ts` is a prototype with
  fictitious values, guarded by `BENCHMARK_IS_PLACEHOLDER`, which renders a
  *"prototype · placeholder numbers"* badge. Replace the constant wholesale when the
  corpus runs; never edit the values into looking real.
- **No fake social proof.** There is no pricing section, no testimonials, and the footer
  links only to GitHub — the project has no other accounts.
- **Status tracks the repository.** The pre-release notice in the hero states what runs
  today. Re-check it against the upstream README after every phase lands.

## Layout

```
src/
  i18n/           types · copy.en · copy.pt · LanguageContext
  content/        data.ts — terminal output, SQL, TOML, URLs, benchmark placeholder
  hooks/          useReveal · useTypewriter · useScrolled
  components/
    layout/       Navbar · Footer
    canvas/       SchemaGraph — the hero animation
                  AsciiSphere — ASCII sphere with a sonar ping, closing CTA only
    ui/           Reveal · Eyebrow · Section · LanguageSwitch · GlowGrid · icons · flags
    illustrations/ VerdictIcons · SafetyIcons — SMIL, no JS
    sections/     Hero · Problem · Verdicts · Safety · Pipeline · Profiles · Ci
                  PriorArt · Benchmark · FinalCta
```

`SchemaGraph` is the hero's argument in motion: table cards start disconnected, edges
draw themselves one at a time and lock with a pulse, and the three relationships that come
back broken shed orphan rows out of the line. It renders the same fictional schema the
terminal below it reports on — `GRAPH_TABLES` and `GRAPH_EDGES` in `data.ts` are kept in
sync with `REPORT` on purpose, so the animation is showing the run you then read. Tables
flagged `ambient` carry no edge, sit deepest and are allowed to run off the box, so the
frame reads as a corner of a much larger schema.

Cards are static and their position, blur, opacity and stacking are pure functions of the
table's `z`, applied in the render rather than from an effect — the first paint is already
correct. `requestAnimationFrame` only touches the five edges. Two constraints are load-
bearing and documented at `GRAPH_TABLES`: cards are placed on a collision-checked grid
(then nudged off it) so none overlaps, and every edge joins two adjacent cells so no wire
passes beneath an unrelated card.

`GlowGrid` is a cut-down adaptation of [React Bits](https://reactbits.dev)' MagicBento:
the cursor-tracking border glow and group spotlight, at roughly a third of the reference
opacity and in the project accent. Particles, 3D tilt, magnetism and the click ripple are
deliberately off — on a page whose argument is restraint, a card that tilts and throws
sparks reads as a different product.

Motion respects `prefers-reduced-motion` throughout: the typewriter jumps to its final
state, reveals are already visible, the sphere paints one static frame, and the glow is
disabled outright.
