package cmd

import (
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/translation"
)

var removeVolumesOnRm, forceRm bool

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [SERVICE...]",
	Short: "Removes stopped service containers.",
	Long: `Removes stopped service containers.

By default, anonymous volumes attached to containers will not be removed. You
can override this with ` + "`-v`" + `. To list all volumes, use ` + "`docker volume ls`" + `.

Any data which is not in a volume will be lost.`,
	Run: func(cmd *cobra.Command, services []string) {
		issuedFlags := []string{}
		if forceRm { issuedFlags = append(issuedFlags, "--force") }
		if removeVolumesOnRm { issuedFlags = append(issuedFlags, "-v") }
		translation.BypassToCompose("rm", append(issuedFlags, services...))
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolVarP(&removeVolumesOnRm, "volumes", "v", false, "Remove any anonymous volumes attached to containers")
	rmCmd.Flags().BoolVarP(&forceRm, "force", "f", false, "Don't ask to confirm removal")
}
