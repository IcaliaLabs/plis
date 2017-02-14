package action

import (
	"os"
	"fmt"
	"github.com/IcaliaLabs/plis/project"
)

func Attach(serviceName string, attachArgs []string) []string {
  containers := project.ContainerStates()
  firstContainer := project.FindFirstRunningContainer(serviceName, containers)

  if firstContainer.Name == "" {
    fmt.Println("No container running for service", serviceName)
    os.Exit(1)
  }

  return []string{"docker", "attach", firstContainer.Name}
}
