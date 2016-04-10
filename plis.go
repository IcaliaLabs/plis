package main

import (
  "fmt"
  "strings"
  "syscall"
  "regexp"
	"os"
	"os/exec"
  "github.com/fatih/color"
  "github.com/codegangsta/cli"
)

type ContainerState struct {
  Name string
  IsRunning bool
}

func GetProjectContainerIds() []string {
  var (
		cmdOut []byte
		err    error
    rawIds []string
    ids    []string
	)

  cmdName := "docker-compose"
	cmdArgs := []string{"ps", "-q"}

  if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
    fmt.Println("Errrr")
    os.Exit(1)
	}

  rawIds = strings.Split(string(cmdOut), "\n")

  ids = rawIds[:0]
  for _, x := range rawIds {
      if x != "" {
          ids = append(ids, x)
      }
  }

  return ids
}

func GetProjectContainerStates() []ContainerState {
  var (
		cmdOut []byte
		err    error
    rawContainerStates []string
    containerStates []ContainerState
	)

  ids := GetProjectContainerIds()

  cmdName := "docker"
	cmdArgs := append([]string{"inspect", "--format='{{.Name}} {{.State.Running}}'"}, ids...)

  if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
    fmt.Println("Errrr")
    os.Exit(1)
	}

  rawContainerStates = strings.Split(string(cmdOut), "\n")

  for i := range rawContainerStates {
    if rawContainerStates[i] != "" {
      fields := strings.Fields(rawContainerStates[i])
      state := ContainerState{}
      state.Name = fields[0][1:len(fields[0])]
      state.IsRunning = fields[1] == "true"
      containerStates = append(containerStates, state)
    }
  }

  return containerStates
}

func FindFirstContainer(serviceName string, containers []ContainerState) ContainerState {
  var foundContainer ContainerState
  rp := regexp.MustCompile("^\\w+_" + serviceName + "_\\d+")

  for i := range containers {
    if rp.FindString(containers[i].Name) != "" { foundContainer = containers[i] }
  }

  return foundContainer
}

func RunGeneratedCommand(command []string) {
  color.Cyan(strings.Join(command, " "))

  binary, lookErr := exec.LookPath(command[0])
  if lookErr != nil {  panic(lookErr) }

  execErr := syscall.Exec(binary, command, os.Environ())
  if execErr != nil { panic(execErr) }
}

func Bypass(cmd string, args []string) {
  RunGeneratedCommand(append([]string{"docker-compose", cmd}, args...))
}

func Start(c *cli.Context) {
  args := []string{"docker-compose", "start"}
  requestedServices := c.Args()

  if len(requestedServices) > 0 {
    servicesAlreadyCreated := []string{}
    containers := GetProjectContainerStates()

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
  } else if len(GetProjectContainerIds()) < 1 {
    args = []string{"docker-compose", "up", "-d"}
  }

  command := append(args, requestedServices...)
  RunGeneratedCommand(command)
}

func Attach(c *cli.Context) {
  containers := GetProjectContainerStates()
  serviceName := c.Args().First()
  firstContainer := FindFirstContainer(serviceName, containers)

  if firstContainer.Name == "" {
    fmt.Println("No container running for service", serviceName)
    os.Exit(1)
  }

  command := []string{"docker", "attach", firstContainer.Name}
  RunGeneratedCommand(command)
}


func Run(c *cli.Context) {
  var cmdArgs []string

  command := c.Args()[1:len(c.Args())]
  containers := GetProjectContainerStates()
  serviceName := c.Args().First()
  firstContainer := FindFirstContainer(serviceName, containers)

  if firstContainer.Name != "" && firstContainer.IsRunning {
    cmdArgs = append([]string{"docker", "exec", "-ti", firstContainer.Name}, command...)
  } else {
    cmdArgs = append([]string{"docker-compose", "run", "--rm", serviceName}, command...)
  }

  RunGeneratedCommand(cmdArgs)
}


func main() {
  app := cli.NewApp()
  app.Name = "Plis"
  app.Usage = "Translates common actions into docker/docker-compose commands by asking nicely"
  app.Version = "0.0.0.build1"

  app.Commands = []cli.Command{
    {
      Name:    "start",
      Usage:   "Starts the project's containers",
      Action:  Start,
    },
    {
      Name:    "stop",
      Usage:   "Stop the project's running processes",
      Action:  func (c *cli.Context) { Bypass("stop", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "restart",
      Usage:   "Restarts the project's running processes",
      Action:  func (c *cli.Context) { Bypass("restart", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "attach",
      Usage:   "Attach the console to a running process",
      Action:  Attach,
      SkipFlagParsing: true,
    },
    {
      Name:    "run",
      Usage:   "Runs a command in a running or new container of a particular service",
      Action:  Run,
      SkipFlagParsing: true,
    },
    {
      Name:    "ps",
      Usage:   "Lists the project's running processes",
      Action:  func (c *cli.Context) { Bypass("ps", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "rm",
      Usage:   "Removes the project's running processes",
      Action:  func (c *cli.Context) { Bypass("rm", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "logs",
      Usage:   "Opens the logs of running processes",
      Action:  func (c *cli.Context) { Bypass("logs", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "down",
      Usage:   "Stops and removes all containers",
      Action:  func (c *cli.Context) { Bypass("down", c.Args()) },
      SkipFlagParsing: true,
    },
  }

  app.Run(os.Args)
}
