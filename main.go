package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "TPP Pr Tools"
	app.Usage = "CLI Tools for evaluating the PR submission."
	app.Author = "Guestlogix"
	app.Version = "0.0.1"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "appsettings",
			Aliases: []string{"a"},
			Usage:   "Generates a markdown report of altered appsettings and posts a comment on the github pr",
			Action: func(c *cli.Context) {
				var repo = c.Args().Get(0)
				appsettings(repo)
			},
		},
	}
}

func main() {
	info()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
