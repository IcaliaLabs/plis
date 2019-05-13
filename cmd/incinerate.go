package cmd

import (
	"github.com/IcaliaLabs/plis/translation"
	"github.com/urfave/cli"
  "fmt"
)

func Incinerate(c *cli.Context) {
  cmdArgs := []string{"docker", "system", "prune", "-a", "-f", "--volumes"}
	translation.Exec(cmdArgs)

  fmt.Printf("All done! Your OS has been incinerated\n")
}
