package cli

import "errors"

// Exit codes are public contract from this release onward: the CI baseline
// command will be consumed by pipelines that branch on them, and changing one
// later breaks builds silently.
const (
	ExitOK = 0

	// ExitFailure is a run that could not complete: connection refused, missing
	// privilege, internal error.
	ExitFailure = 1

	ExitUsage = 2

	// ExitFindings signals something the caller asked to be told about — a new
	// undeclared relationship, a growing orphan count. Reserved for `check`.
	ExitFindings = 3

	// ExitInterrupted follows the shell convention of 128 plus the signal
	// number, SIGINT being 2.
	ExitInterrupted = 130
)

// usageError marks a command-line mistake rather than a runtime failure.
type usageError struct{ err error }

func (e usageError) Error() string { return e.err.Error() }
func (e usageError) Unwrap() error { return e.err }

// UsageError wraps err as a command-line usage problem.
func UsageError(err error) error { return usageError{err} }

// ExitCodeFor maps an error to the exit code that should accompany it.
func ExitCodeFor(err error) int {
	switch {
	case err == nil:
		return ExitOK
	case errors.Is(err, errInterrupted):
		return ExitInterrupted
	case errors.As(err, &usageError{}):
		return ExitUsage
	default:
		return ExitFailure
	}
}

var errInterrupted = errors.New("interrupted")
