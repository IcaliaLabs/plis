package cmd

import (
	"strconv"
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/translation"
)

// stopCmd represents the stop command
var timeout uint64

var stopCmd = &cobra.Command{
	Use:   "stop [SERVICE...]",
	Short: "Stop project's running containers without removing them.",
	Long: `Stops project's running containers without removing them.

They can be started again with ` + "`plis start`" + ` or ` + "`docker-compose start`.",
	Run: func(cmd *cobra.Command, services []string) {
		issuedFlags := []string{}

		if timeout > 0 && timeout != 10 {
			issuedFlags = append(issuedFlags, "--timeout", strconv.FormatUint(timeout, 10))
		}

		translation.BypassToCompose("stop", append(issuedFlags, services...))
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	stopCmd.Flags().Uint64VarP(&timeout, "timeout", "t", 10, "Specify a shutdown timeout in seconds.")

}
