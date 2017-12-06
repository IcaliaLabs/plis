package cmd

import (
	"../project"
	"../translation"
	"github.com/urfave/cli"
	"regexp"

	"../grouping"
	"path/filepath"
	"log"
)

func Start(c *cli.Context) {
	args := []string{"docker-compose", "start"}
	requestedServicesOrGroups := c.Args()
	requestedServices := []string{}

	composefile := "docker-compose.yml"

	if len(requestedServicesOrGroups) > 0 {
		serviceGroups, err := grouping.GetServiceGroupingFrom(composefile)
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

		absPath, err := filepath.Abs(filepath.Dir(composefile))
		if err != nil { log.Fatalf("error: %v", err) }
		_, dir := filepath.Split(absPath)
		rr := regexp.MustCompile("[^a-zA-Z0-9]")
		dir = rr.ReplaceAllString(dir, "")

		for _, serviceName := range requestedServices {
			rp := regexp.MustCompile("^" + dir + "_" + serviceName + "_\\d+")

			for p := range containers {
				if rp.FindString(containers[p].Name) != "" {
					servicesAlreadyCreated = append(servicesAlreadyCreated, serviceName)
				}
			}
		}

		shouldUseUp := false
		for _, requestedService := range requestedServices {
			requestedServiceIsUp := false
			for _, serviceAlreadyCreated := range servicesAlreadyCreated {
		    requestedServiceIsUp = (serviceAlreadyCreated == requestedService)
			}
			shouldUseUp = shouldUseUp || !requestedServiceIsUp
		}

		if shouldUseUp {
			args = []string{"docker-compose", "up", "-d"}
		}
	} else if len(project.ContainerIds()) < 1 {
		args = []string{"docker-compose", "up", "-d"}
	}

	command := append(args, requestedServices...)
	translation.Exec(command)
}
