package cmd

import (
	"github.com/urfave/cli"
	"github.com/IcaliaLabs/plis/translation"
)

func Compose(c *cli.Context) {
  translation.Exec(append([]string{"docker-compose", c.Command.Name}, c.Args()...))
}
