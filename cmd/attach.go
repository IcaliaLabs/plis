package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/project"
	"github.com/IcaliaLabs/plis/translation"
)

var noStdin, sigProxy bool
var detachKeys string

// attachCmd represents the attach command
var attachCmd = &cobra.Command{
	Use:   "attach SERVICE",
	Short: "Attach the console to a running container process",
	Long: `Attach the console to a running container process`,
	Run: func(cmd *cobra.Command, args []string) {
		issuedFlags := []string{}

		if noStdin { issuedFlags = append(issuedFlags, "--no-stdin") }
		if sigProxy { issuedFlags = append(issuedFlags, "--sig-proxy") }
		if detachKeys != "" { issuedFlags = append(issuedFlags, "--detach-keys", detachKeys) }

		doAttach(args[0], issuedFlags)
	},
}

func doAttach(serviceName string, flags []string) {
  containers := project.ContainerStates()
  firstContainer := project.FindFirstRunningContainer(serviceName, containers)

  if firstContainer.Name == "" {
    fmt.Println("No container running for service", serviceName)
    os.Exit(1)
  }

  command := append([]string{"docker", "attach", firstContainer.Name}, flags...)
  translation.Exec(command)
}

func init() {
	RootCmd.AddCommand(attachCmd)
	attachCmd.Flags().StringVar(&detachKeys, "detach-keys", "", "Override the key sequence for detaching a container")
	attachCmd.Flags().BoolVar(&noStdin, "no-stdin", false, "Do not attach STDIN")
	attachCmd.Flags().BoolVar(&sigProxy, "sig-proxy", false, "Proxy all received signals to the process")
}
