package cli

import (
	"io"
	"os"
)

// ColorMode controls ANSI emission.
type ColorMode string

// Colour modes. Auto resolves against the destination and the environment;
// the other two override that decision outright.
const (
	ColorAuto   ColorMode = "auto"
	ColorAlways ColorMode = "always"
	ColorNever  ColorMode = "never"
)

// Streams is where a command writes.
//
// The split is absolute: Out carries the result meant for consumption — the
// table, the JSON, the SQL — and Err carries everything else, including
// progress, warnings and errors. The output of this tool will be piped into
// files and CI pipelines, and one stray diagnostic line on stdout corrupts
// programmatic consumption of everything else.
type Streams struct {
	// Out receives results. Never diagnostics.
	Out io.Writer

	// Err receives diagnostics. Never results.
	Err io.Writer

	// In is reserved for future interactive prompts.
	In io.Reader

	color bool
}

// StdStreams wires a Streams to the process, resolving whether colour is
// appropriate for the destination.
func StdStreams(mode ColorMode) *Streams {
	s := &Streams{Out: os.Stdout, Err: os.Stderr, In: os.Stdin}
	s.color = resolveColor(mode, os.Stdout)
	return s
}

// Color reports whether ANSI sequences may be emitted.
func (s *Streams) Color() bool { return s.color }

// resolveColor decides whether to emit ANSI, honouring the explicit override
// first, then NO_COLOR, then whether the destination is an interactive
// terminal. Anything piped or redirected gets clean text.
func resolveColor(mode ColorMode, out *os.File) bool {
	switch mode {
	case ColorAlways:
		return true
	case ColorNever:
		return false
	}

	// https://no-color.org — presence is enough, whatever the value.
	if _, set := os.LookupEnv("NO_COLOR"); set {
		return false
	}

	// A dumb terminal cannot render escapes meaningfully.
	if os.Getenv("TERM") == "dumb" {
		return false
	}

	return isTerminal(out)
}

// isTerminal reports whether f is an interactive character device. Checking the
// file mode keeps this dependency-free; the alternative is a cgo-adjacent
// package for something the standard library already answers.
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
