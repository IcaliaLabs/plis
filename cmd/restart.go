package cmd

import (
	"strconv"
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/translation"
)

var restartTimeout uint64
// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart running containers.",
	Long: `Restart running containers.`,
	Run: func(cmd *cobra.Command, services []string) {
		issuedFlags := []string{}

		if restartTimeout > 0 && restartTimeout != 10 {
			issuedFlags = append(issuedFlags, "--timeout", strconv.FormatUint(restartTimeout, 10))
		}

		translation.BypassToCompose("restart", append(issuedFlags, services...))
	},
}

func init() {
	RootCmd.AddCommand(restartCmd)
	restartCmd.Flags().Uint64VarP(&restartTimeout, "timeout", "t", 10, "Specify a shutdown timeout in seconds.")
}
