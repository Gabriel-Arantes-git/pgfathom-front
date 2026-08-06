package cli

import "errors"

// Exit codes are part of the tool's public contract from this release onward,
// because the CI baseline command will be consumed by pipelines that branch on
// them. Changing a code later breaks builds silently, so they are fixed here.
const (
	// ExitOK means the command completed.
	ExitOK = 0

	// ExitFailure means the command could not complete: connection refused,
	// missing privilege, internal error.
	ExitFailure = 1

	// ExitUsage means the command line was wrong.
	ExitUsage = 2

	// ExitFindings means the command completed and has something to report that
	// the caller asked to be signalled — a new undeclared relationship, a
	// growing orphan count. Reserved here, used by `check`.
	ExitFindings = 3

	// ExitInterrupted means the run was cancelled by a signal. It follows the
	// shell convention of 128 plus the signal number, SIGINT being 2.
	ExitInterrupted = 130
)

// usageError marks an error as a command-line mistake rather than a runtime
// failure, so the top level can pick the right exit code.
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
