package cmd

import (
	"github.com/IcaliaLabs/plis/project"
	"github.com/IcaliaLabs/plis/translation"
	"github.com/urfave/cli"
)

func Run(c *cli.Context) {
	var cmdArgs []string

	command := c.Args()[1:len(c.Args())]
	containers := project.ContainerStates()
	serviceName := c.Args().First()
	firstContainer := project.FindFirstRunningContainer(serviceName, containers)

	if firstContainer.Name != "" && firstContainer.IsRunning {
		cmdArgs = append([]string{"docker", "exec", "-ti", firstContainer.Name}, command...)
	} else {
		cmdArgs = append([]string{"docker-compose", "run", "--rm", serviceName}, command...)
	}

	translation.Exec(cmdArgs)
}
