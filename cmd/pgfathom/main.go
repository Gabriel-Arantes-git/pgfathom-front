// Command pgfathom sounds the depth of a legacy PostgreSQL schema.
package main

import (
	"os"

	"github.com/lvcas-dotcom/pgfathom/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], cli.StdStreams(cli.ColorAuto)))
}
