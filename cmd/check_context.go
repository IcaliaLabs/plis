package cmd

import (
	"fmt"
	"../project"
	"github.com/urfave/cli"
	"strconv"
)

/*
Displays the project context that will be sent to the ` + "`docker build`" + ` command,
ommitting the files matching the rules in the ` + "`.dockerignore`" + ` file.
*/

func CheckContext(c *cli.Context) {
	contextFiles, err := project.ContextFiles(".")
	if err != nil {
		fmt.Print(err)
	}

	for _, contextFile := range contextFiles {
		fmt.Println(contextFile)
	}

	fmt.Println("")
	fmt.Println("============================================================")
	fmt.Println("")
	fmt.Println("Total files: " + strconv.Itoa(len(contextFiles)))
}
