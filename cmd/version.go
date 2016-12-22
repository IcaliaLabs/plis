package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var plisVersion = "0.0.0.build6"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version",
	Long: `Displays the version`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintPlisVersion()
	},
}

func PrintPlisVersion() {
	fmt.Println("Plis: " + plisVersion)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
