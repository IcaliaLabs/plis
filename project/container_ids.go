package project

import (
	"../translation"
	"strings"
)

func ContainerIds() []string {
	var (
		rawIds []string
		ids    []string
	)

	cmdName := "docker-compose"
	cmdArgs := []string{"ps", "-q"}
	cmdOut := translation.ShellOut(cmdName, cmdArgs)

	rawIds = strings.Split(cmdOut, "\n")

	ids = rawIds[:0]
	for _, x := range rawIds {
		if x != "" {
			ids = append(ids, x)
		}
	}

	return ids
}
