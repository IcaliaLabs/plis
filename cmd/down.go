package cmd

import (
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/translation"
)

var removeVolumes, removeOrphans bool
var removeImages string

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove the current project's containers, networks, images, and volumes",
	Long: `Stops containers and removes containers, networks, volumes, and images
created by ` + "`plis start`" + ` or ` + "`docker-compose up`" + `.

By default, the only things removed are:

- Containers for services defined in the Compose file
- Networks defined in the ` + "`networks`" + ` section of the Compose file
- The default network, if one is used

Networks and volumes defined as ` + "`external`" + ` are never removed.`,
	Run: func(cmd *cobra.Command, args []string) {
		issuedFlags := []string{}

		if removeVolumes { issuedFlags = append(issuedFlags, "--volumes") }
		if removeOrphans { issuedFlags = append(issuedFlags, "--remove-orphans") }
		if removeImages != "" { issuedFlags = append(issuedFlags, "--rmi", removeImages) }

		translation.BypassToCompose("down", append(issuedFlags, args...))
	},
}

func init() {
	RootCmd.AddCommand(downCmd)

	downCmd.Flags().BoolVarP(&removeVolumes, "volumes", "v", false, `Remove named volumes declared in the 'volumes' section
		  	 of the Compose file and anonymous volumes
		  	 attached to containers.`)

	downCmd.Flags().BoolVar(&removeOrphans, "remove-orphans", false, `Remove containers for services not defined in the
			 Compose file.`)

  downCmd.Flags().StringVar(&removeImages, "rmi", "", `Remove images. Type must be one of:
			 'all': Remove all images used by any service.
			 'local': Remove only images that don't have a custom tag
			 set by the ` + "`image`" + ` field.`)
}
