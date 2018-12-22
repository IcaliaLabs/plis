package project

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// FindFirstRunningContainer finds the state of the first running container
// matching the given service name.
func FindFirstRunningContainer(serviceName string, containers []ContainerState) ContainerState {
	var foundContainer ContainerState
	currentWorkdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dirName := filepath.Base(currentWorkdir)
	rp := regexp.MustCompile("^" + dirName + "_" + serviceName + "_\\d+")

	for i := range containers {
		if rp.FindString(containers[i].Name) != "" {
			foundContainer = containers[i]
		}
	}

	return foundContainer
}
