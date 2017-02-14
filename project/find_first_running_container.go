package project

import (
	"regexp"
)

func FindFirstRunningContainer(serviceName string, containers []ContainerState) ContainerState {
	var foundContainer ContainerState
	rp := regexp.MustCompile("^\\w+_" + serviceName + "_\\d+")

	for i := range containers {
		if rp.FindString(containers[i].Name) != "" {
			foundContainer = containers[i]
		}
	}

	return foundContainer
}
