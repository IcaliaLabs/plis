package cmd

import (
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/translation"
)

var quietPs bool

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps [SERVICE...]",
	Short: "List project containers.",
	Long: "List project containers.",
	Run: func(cmd *cobra.Command, services []string) {
		issuedFlags := []string{}
		if quietPs { issuedFlags = append(issuedFlags, "-q") }
		translation.BypassToCompose("ps", append(issuedFlags, services...))
	},
}

func init() {
	RootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolVarP(&quietPs, "quiet", "q", false, "Only display IDs")

}
