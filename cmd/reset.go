package cmd

import (
	"flag"

	"github.com/urfave/cli"
)

// Reset will invoke the `docker-compose up` command with the `--force-recreate`
// flag
func Reset(c *cli.Context) {
	newFlags := flag.NewFlagSet("contrive", 0)
	newContext := cli.NewContext(c.App, newFlags, c)
	Start(newContext)
}
