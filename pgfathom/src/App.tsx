import { Navbar } from './components/layout/Navbar'
import { Footer } from './components/layout/Footer'
import { Hero } from './components/sections/Hero'
import { Problem } from './components/sections/Problem'
import { Verdicts } from './components/sections/Verdicts'
import { Safety } from './components/sections/Safety'
import { Pipeline } from './components/sections/Pipeline'
import { Profiles } from './components/sections/Profiles'
import { Ci } from './components/sections/Ci'
import { PriorArt } from './components/sections/PriorArt'
import { Benchmark } from './components/sections/Benchmark'
import { FinalCta } from './components/sections/FinalCta'

function App() {
  return (
    <div className="relative min-h-screen">
      <div
        aria-hidden="true"
        className="noise pointer-events-none fixed inset-0 z-100 opacity-[0.022]"
      />

      <Navbar />

      <main>
        <Hero />
        <Problem />
        <Verdicts />
        <Safety />
        <Pipeline />
        <Profiles />
        <Ci />
        <PriorArt />
        <Benchmark />
        <FinalCta />
      </main>

      <Footer />
    </div>
  )
}

export default App
