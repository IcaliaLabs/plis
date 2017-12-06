package project

import (
	"../translation"
	"strings"
)

func ContainerStates() []ContainerState {
	var (
		rawContainerStates []string
		containerStates    []ContainerState
	)

	ids := ContainerIds()

	if len(ids) > 0 {
		cmdName := "docker"
		cmdArgs := append([]string{"inspect", "--format='{{.Name}} {{.State.Running}}'"}, ids...)
		cmdOut := translation.ShellOut(cmdName, cmdArgs)

		rawContainerStates = strings.Split(cmdOut, "\n")

		for i := range rawContainerStates {
			if rawContainerStates[i] != "" {
				rawContainerState := strings.Trim(rawContainerStates[i], "'")
				fields := strings.Fields(rawContainerState)
				state := ContainerState{}
				state.Name = fields[0][1:len(fields[0])]
				state.IsRunning = fields[1] == "true"
				containerStates = append(containerStates, state)
			}
		}
	}

	return containerStates
}
