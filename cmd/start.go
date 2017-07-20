package cmd

import (
	"github.com/IcaliaLabs/plis/project"
	"github.com/IcaliaLabs/plis/translation"
	"github.com/urfave/cli"
	"regexp"

	"github.com/IcaliaLabs/plis/grouping"
	"log"
)

func Start(c *cli.Context) {
	args := []string{"docker-compose", "start"}
	requestedServicesOrGroups := c.Args()
	requestedServices := []string{}

	if len(requestedServicesOrGroups) > 0 {
		serviceGroups, err := grouping.GetServiceGroupingFrom("docker-compose.yml")
		if err != nil { log.Fatalf("error: %v", err) }

		for requestedServiceOrGroupIndex := range requestedServicesOrGroups {
			serviceOrGroupName := requestedServicesOrGroups[requestedServiceOrGroupIndex]
			if groupedServices, ok := serviceGroups[serviceOrGroupName]; ok {
				requestedServices = append(requestedServices, groupedServices...)
			} else {
				requestedServices = append(requestedServices, serviceOrGroupName)
			}
		}

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
