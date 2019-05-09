package project

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

	containerName := strings.ToLower(dirName + "_" + serviceName)

	rp := regexp.MustCompile("^" + containerName + "_\\d+")

	for i := range containers {
		if rp.FindString(containers[i].Name) != "" {
			foundContainer = containers[i]
		}
	}

	return foundContainer
}
