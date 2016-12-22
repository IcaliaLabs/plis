package cmd

import (
	"regexp"
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/project"
	"github.com/IcaliaLabs/plis/translation"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the project's containers",
	Long: `Starts the project's containers.`,
	Run: doStart,
}

func doStart(cmd *cobra.Command, requestedServices []string) {
  args := []string{"docker-compose", "start"}

  if len(requestedServices) > 0 {
    servicesAlreadyCreated := []string{}
    containers := project.ContainerStates()

    for i := range requestedServices {
      serviceName := requestedServices[i]
      rp := regexp.MustCompile("^\\w+_" + serviceName + "_\\d+")

      for p := range containers {
        if rp.FindString(containers[p].Name) != "" {
          servicesAlreadyCreated = append(servicesAlreadyCreated, serviceName)
        }
      }
    }

    if len(servicesAlreadyCreated) != len(requestedServices) {
      args = []string{"docker-compose", "up", "-d"}
    }
  } else if len(project.ContainerIds()) < 1 {
    args = []string{"docker-compose", "up", "-d"}
  }

  command := append(args, requestedServices...)
  translation.Exec(command)
}

func init() {
	RootCmd.AddCommand(startCmd)
}
