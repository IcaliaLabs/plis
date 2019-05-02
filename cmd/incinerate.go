package cmd

import (
	"github.com/IcaliaLabs/plis/translation"
	"github.com/urfave/cli"
  "strings"
  "fmt"
)

func IdsGetter(cmdOut string) []string {
  var (
    rawIds []string
    ids []string
  )

	rawIds = strings.Split(cmdOut, "\n")

	ids = rawIds[:0]
	for _, x := range rawIds {
		if x != "" {
			ids = append(ids, x)
		}
	}

	return ids
}

func DockerSystemNetworksIds() []string {
	cmdName := "docker"
	cmdArgs := []string{"network", "ls", "-q"}
	cmdOut := translation.ShellOut(cmdName, cmdArgs)
  ids := IdsGetter(cmdOut)

  return ids
}

func DockerSystemVolumesIds() []string {
	cmdName := "docker"
	cmdArgs := []string{"volume", "ls", "-q"}
	cmdOut := translation.ShellOut(cmdName, cmdArgs)
  ids := IdsGetter(cmdOut)

  return ids
}

func ImagesIds() []string {
	cmdName := "docker"
	cmdArgs := []string{"image", "ls", "-aq"}
	cmdOut := translation.ShellOut(cmdName, cmdArgs)
  ids := IdsGetter(cmdOut)

  return ids
}

func DockerRunningContainersIds() []string {
	cmdName := "docker"
	cmdArgs := []string{"ps", "-aq"}
	cmdOut := translation.ShellOut(cmdName, cmdArgs)
  ids := IdsGetter(cmdOut)

  return ids
}

func IncinerateImages() {
  fmt.Printf("Incinerating the images...\n")
  imagesIds := ImagesIds()
  cmdArgs := append([]string{"rmi", "-f"}, imagesIds...)
  translation.ShellOut("docker", cmdArgs)
}

func IncinerateVolumes() {
  fmt.Printf("Incinerating the volumes...\n")
  volumesIds := DockerSystemVolumesIds()
  cmdArgs := append([]string{"volume", "rm", "-f"}, volumesIds...)
  translation.ShellOut("docker", cmdArgs)
}

func StopRunningContainers() {
  fmt.Printf("Stopping the running containers...\n")
  dockerContainerIds := DockerRunningContainersIds()
  cmdArgs := append([]string{"stop"}, dockerContainerIds...)
  translation.ShellOut("docker", cmdArgs)
}

func RemoveRunningContainers() {
  fmt.Printf("Removing the running containers...\n")
  dockerContainerIds := DockerRunningContainersIds()
  cmdArgs := append([]string{"rm"}, dockerContainerIds...)
  translation.ShellOut("docker", cmdArgs)
}

func Incinerate(c *cli.Context) {
  StopRunningContainers()
  RemoveRunningContainers()
  IncinerateImages()
  IncinerateVolumes()

  fmt.Printf("All done! Your OS has been incinerated\n")
}
