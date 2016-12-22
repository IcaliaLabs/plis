package cmd

import (
	"github.com/spf13/cobra"
	"github.com/IcaliaLabs/plis/translation"
)

var pullImages, dontUseCache, removeIntermediateContainers bool

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [SERVICE...]",
	Short: "Build or rebuild services",
	Long: `Build or rebuild services.

Services are built once and then tagged as ` + "`project_service`" + `,
e.g. ` + "`composetest_db`" + `. If you change a service's ` + "`Dockerfile`" + ` or the
contents of its build directory, you can run ` + "`docker-compose build`" + ` to rebuild it.`,
	Run: func(cmd *cobra.Command, services []string) {
		issuedFlags := []string{}
		if pullImages { issuedFlags = append(issuedFlags, "--pull") }
		if dontUseCache { issuedFlags = append(issuedFlags, "--no-cache") }
		if removeIntermediateContainers { issuedFlags = append(issuedFlags, "--force-rm") }
		translation.BypassToCompose("build", append(issuedFlags, services...))
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolVar(&dontUseCache, "no-cache", false, "Do not use cache when building the image")
	buildCmd.Flags().BoolVar(&pullImages, "pull", false, "Always attempt to pull a newer version of the image")
	buildCmd.Flags().BoolVar(&removeIntermediateContainers, "force-rm", false, "Always remove intermediate containers")
}
