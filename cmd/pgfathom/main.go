// Command pgfathom sounds the depth of a legacy PostgreSQL schema.
//
// It is strictly read-only. It never issues a statement that modifies the
// database under analysis, and no value read from a user table ever reaches its
// output, logs or generated files.
package main

import (
	"os"

	"github.com/lvcas-dotcom/pgfathom/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], cli.StdStreams(cli.ColorAuto)))
}
