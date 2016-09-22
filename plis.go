package main

import (
  "fmt"
  "strings"
  "syscall"
  "regexp"
	"os"
	"os/exec"
  "github.com/fatih/color"
  "github.com/urfave/cli"
  "net/url"
  "gopkg.in/yaml.v2"
  "io/ioutil"
)

type ContainerState struct {
  Name string
  IsRunning bool
}

type ComposeService struct {
  EnvFile []string
}

type Compose struct {
  Services map[string]ComposeService
}

func PanicIfError(e error) {
  if e != nil { panic(e) }
}

func GetProjectContainerIds() []string {
  var (
		rawIds []string
    ids    []string
	)

  cmdName := "docker-compose"
	cmdArgs := []string{"ps", "-q"}
  cmdOut  := ShellOutGeneratedCommand(cmdName, cmdArgs)

  rawIds = strings.Split(cmdOut, "\n")

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
		rawContainerStates []string
    containerStates []ContainerState
	)

  ids := GetProjectContainerIds()

  if len(ids) > 0 {
    cmdName := "docker"
  	cmdArgs := append([]string{"inspect", "--format='{{.Name}} {{.State.Running}}'"}, ids...)
    cmdOut  := ShellOutGeneratedCommand(cmdName, cmdArgs)

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

func FindFirstContainer(serviceName string, containers []ContainerState) ContainerState {
  var foundContainer ContainerState
  rp := regexp.MustCompile("^\\w+_" + serviceName + "_\\d+")

  for i := range containers {
    if rp.FindString(containers[i].Name) != "" { foundContainer = containers[i] }
  }

  return foundContainer
}

func ShellOutGeneratedCommand(binName string, args []string) string {
  var (
    output []byte
    error  error
  )
  output, error = exec.Command(binName, args...).Output()

  if error != nil {
    color.Red("Command:", binName + " " + strings.Join(args, " "))
    color.Red("Error:", error)
    os.Exit(1)
	}

  return string(output)
}

func RunGeneratedCommand(command []string) {
  color.Cyan(strings.Join(command, " "))

  binary, lookErr := exec.LookPath(command[0])
  if lookErr != nil {  panic(lookErr) }

  execErr := syscall.Exec(binary, command, os.Environ())
  if execErr != nil { panic(execErr) }
}

func BypassToCompose(cmd string, args []string) {
  RunGeneratedCommand(append([]string{"docker-compose", cmd}, args...))
}

func Start(c *cli.Context) {
  EnsureEnvFilesExist()
  // StartServices(c.Args())
}

func EnsureEnvFilesExist() {
  source, err := ioutil.ReadFile("docker-compose.yml")
  PanicIfError(err)

  var compose Compose

  err = yaml.Unmarshal(source, &compose)
  PanicIfError(err)

  fmt.Printf("Value: %#v\n", compose.Services[1])
}

func StartServices(requestedServices []string) {
  args := []string{"docker-compose", "start"}
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

func Clun(c *cli.Context) {
  CloneProject(c)
  os.Chdir(ClunDirName(c))
  StartServices([]string{})
}

func CloneProject(c *cli.Context) {
  repoUri := ClunRepoUri(c)
  dirName := ClunDirName(c)
  fmt.Println("Cloning into '" + dirName + "'...")
  cloneCommand := exec.Command("git", "clone", repoUri.String(), dirName)
  _, cloneErr := cloneCommand.Output()
  if cloneErr != nil { panic(cloneErr) }

}

func ClunDirName(c *cli.Context) string {
  dirName := c.Args()[1:].First()
  if dirName == "" {
    repoUriPaths := strings.Split(ClunRepoUri(c).Path, "/")
    dirName = strings.Trim(repoUriPaths[len(repoUriPaths)-1], ".git")
  }
  return dirName
}

func ClunRepoUri(c *cli.Context) *url.URL {
  repoUri, repoUriParseError := url.Parse(c.Args().First())
  if repoUriParseError != nil { panic(repoUriParseError) }
  return repoUri
}

func main() {
  app := cli.NewApp()
  app.Name = "Plis"
  app.Usage = "Translates common development actions into docker/docker-compose commands by asking nicely"
  app.Version = "0.0.0.build5"

  app.Commands = []cli.Command{
    {
      Name:    "start",
      Usage:   "Starts the project's containers",
      Action:  Start,
    },
    {
      Name:    "stop",
      Usage:   "Stop the project's running processes",
      Action:  func (c *cli.Context) { BypassToCompose("stop", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "restart",
      Usage:   "Restarts the project's running processes",
      Action:  func (c *cli.Context) { BypassToCompose("restart", c.Args()) },
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
      Name:    "clun",
      Usage:   "Clones a Git project, copies/generates the project's required dotenv files and starts the whole project",
      Action:  Clun,
    },
    {
      Name:    "ps",
      Usage:   "Lists the project's running processes",
      Action:  func (c *cli.Context) { BypassToCompose("ps", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "rm",
      Usage:   "Removes the project's running processes",
      Action:  func (c *cli.Context) { BypassToCompose("rm", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "logs",
      Usage:   "Opens the logs of running processes",
      Action:  func (c *cli.Context) { BypassToCompose("logs", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "down",
      Usage:   "Stops and removes all containers",
      Action:  func (c *cli.Context) { BypassToCompose("down", c.Args()) },
      SkipFlagParsing: true,
    },
    {
      Name:    "build",
      Usage:   "Build or rebuild services",
      Action:  func (c *cli.Context) { BypassToCompose("build", c.Args()) },
      SkipFlagParsing: true,
    },
  }

  app.Run(os.Args)
}
