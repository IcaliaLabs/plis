package main

import (
	"os"
  "github.com/urfave/cli"
  "github.com/IcaliaLabs/plis/cmd"
)

func main() {
  app := cli.NewApp()
  app.Name = "Plis"
  app.Usage = "Translates common development actions into docker/docker-compose commands by asking nicely"
  app.Version = "0.0.0.build5"

  app.Commands = []cli.Command{
    {
      Name:    "start",
      Usage:   "Starts the project's containers",
      Action:  cmd.Start,
    },
    {
      Name:    "stop",
      Usage:   "Stop the project's running processes",
      Action:  cmd.Compose,
      SkipFlagParsing: true,
    },
    {
      Name:    "restart",
      Usage:   "Restarts the project's running processes",
      Action:  cmd.Compose,
      SkipFlagParsing: true,
    },
    {
      Name:    "attach",
      Usage:   "Attach the console to a running process",
      Action:  cmd.Attach,
      SkipFlagParsing: true,
    },
    {
      Name:    "run",
      Usage:   "Runs a command in a running or new container of a particular service",
      Action:  cmd.Run,
      SkipFlagParsing: true,
    },
    {
      Name:    "ps",
      Usage:   "Lists the project's running processes",
      Action:  cmd.Compose,
      SkipFlagParsing: true,
    },
    {
      Name:    "rm",
      Usage:   "Removes the project's running processes",
      Action:  cmd.Compose,
      SkipFlagParsing: true,
    },
    {
      Name:    "logs",
      Usage:   "Opens the logs of running processes",
      Action:  cmd.Compose,
      SkipFlagParsing: true,
    },
    {
      Name:    "down",
      Usage:   "Stops and removes all containers",
      Action:  cmd.Compose,
      SkipFlagParsing: true,
    },
    {
      Name:    "build",
      Usage:   "Build or rebuild services",
      Action:  cmd.Compose,
      SkipFlagParsing: true,
    },
    {
      Name:    "check",
      Usage:   "Performs several checks inside your project",
      Subcommands: []cli.Command{
        {
          Name:  "context",
          Usage: "Displays the project context",
          Action: cmd.CheckContext,
        },
      },
    },
  }

  app.Run(os.Args)
}
