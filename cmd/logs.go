package cmd

import (
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/translation"
)

var followLogs, showLogTimestamps, logsWithoutColor bool
var logTail string

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [SERVICE...]",
	Short: "View output from containers",
	Long: "View output from containers",
	Run: func(cmd *cobra.Command, services []string) {
		issuedFlags := []string{}
		if followLogs { issuedFlags = append(issuedFlags, "--follow") }
		if logsWithoutColor { issuedFlags = append(issuedFlags, "--no-color") }
		if logTail != "" { issuedFlags = append(issuedFlags, "--tail", logTail) }
		if showLogTimestamps { issuedFlags = append(issuedFlags, "--timestamps") }
		translation.BypassToCompose("logs", append(issuedFlags, services...))
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolVarP(&followLogs, "follow", "f", false, "Follow log output")
	logsCmd.Flags().BoolVarP(&showLogTimestamps, "timestamps", "t", false, "Show timestamps")
	logsCmd.Flags().BoolVarP(&logsWithoutColor, "no-color", "m", false, "Produce monochrome output")
	logsCmd.Flags().StringVar(&logTail, "tail", "", `Number of lines to show from the end of the logs
		 for each container`)
}
