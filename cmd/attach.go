package cmd

import (
	"fmt"
	"github.com/IcaliaLabs/plis/project"
	"github.com/IcaliaLabs/plis/translation"
	"github.com/urfave/cli"
	"os"
)

func Attach(c *cli.Context) {
	containers := project.ContainerStates()
	serviceName := c.Args().First()
	firstContainer := project.FindFirstRunningContainer(serviceName, containers)

	if firstContainer.Name == "" {
		fmt.Println("No container running for service", serviceName)
		os.Exit(1)
	}

	command := []string{"docker", "attach", firstContainer.Name}
	translation.Exec(command)
}
