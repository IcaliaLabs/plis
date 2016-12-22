package cmd

import (
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/project"
	"github.com/IcaliaLabs/plis/translation"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a command in a running or new container",
	Long: `Runs a command in a running or new container of a particular service`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		command := args[1:len(args)]
		doRun(serviceName, command)
	},
}

func doRun(serviceName string, command []string) {
  var cmdArgs []string

  containers := project.ContainerStates()
  firstContainer := project.FindFirstRunningContainer(serviceName, containers)

  if firstContainer.Name != "" && firstContainer.IsRunning {
    cmdArgs = append([]string{"docker", "exec", "-ti", firstContainer.Name}, command...)
  } else {
    cmdArgs = append([]string{"docker-compose", "run", "--rm", serviceName}, command...)
  }

  translation.Exec(cmdArgs)
}

func init() {
	RootCmd.AddCommand(runCmd)
}
