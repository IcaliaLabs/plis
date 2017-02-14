package action

import (
	"regexp"
	"github.com/IcaliaLabs/plis/project"
)

func Start(servicesToStart []string) []string {
  args := []string{"docker-compose", "start"}

  if len(servicesToStart) > 0 {
    servicesAlreadyCreated := []string{}
    containers := project.ContainerStates()

    for i := range servicesToStart {
      serviceName := servicesToStart[i]
      rp := regexp.MustCompile("^\\w+_" + serviceName + "_\\d+")

      for p := range containers {
        if rp.FindString(containers[p].Name) != "" {
          servicesAlreadyCreated = append(servicesAlreadyCreated, serviceName)
        }
      }
    }

    if len(servicesAlreadyCreated) != len(servicesToStart) {
      args = []string{"docker-compose", "up", "-d"}
    }
  } else if len(project.ContainerIds()) < 1 {
    args = []string{"docker-compose", "up", "-d"}
  }

  return append(args, servicesToStart...)
}
