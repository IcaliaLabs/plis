package cmd

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/IcaliaLabs/plis/action"
	"github.com/IcaliaLabs/plis/translation"
)

func Start(c *cli.Context) {

	fmt.Println("plis start called...")

  servicesToStart := c.Args()
	startOneService := (len(servicesToStart) == 1)

  command := action.Start(servicesToStart)
  translation.Exec(command)

	if startOneService {
		serviceToStart := servicesToStart[0]
		attachArgs := []string{}
		fmt.Println("Expressely asked to start one service (with dependencies): ", serviceToStart)
		command = action.Attach(serviceToStart, attachArgs)
		translation.Exec(command)
	} else {
		fmt.Println("No service was expressely asked to start...")
	}
}
