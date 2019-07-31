package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "PAT"
	app.Usage = "CLI Tools for pipelines."
	app.Author = "Guestlogix"
	app.Version = "0.0.1"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:      "appsettings",
			Aliases:   []string{"a"},
			Usage:     "Generates a markdown report of altered appsettings and posts a comment on the github pr",
			ArgsUsage: "<diff>",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "markdown, m", Usage: "return the results as a markdown table"},
			},
			Action: func(c *cli.Context) {
				var repo = c.Args().Get(0)
				appsettings(repo)
			},
		},
		{
			Name:      "releasenotes",
			Aliases:   []string{"r"},
			Usage:     "Generates the release notes between two tags",
			ArgsUsage: "<from> <to>",
			Action: func(c *cli.Context) {
				var from = c.Args().Get(0)
				var to = c.Args().Get(1)
				releasenotes(from, to)
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
