package cmd

import (
	"github.com/IcaliaLabs/plis/translation"
	"github.com/urfave/cli"
)

func Compose(c *cli.Context) {
	translation.Exec(append([]string{"docker-compose", c.Command.Name}, c.Args()...))
}
