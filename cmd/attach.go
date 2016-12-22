package cmd

import (
	"os"
	"fmt"
	"github.com/urfave/cli"
	"github.com/IcaliaLabs/plis/project"
	"github.com/IcaliaLabs/plis/translation"
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
