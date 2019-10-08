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
			Name:      "releasenotes",
			Aliases:   []string{"rn"},
			Usage:     "Generates the release notes between two tags",
			ArgsUsage: "<path> <from> <to>",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "output, o", Usage: "The filepath and file name to output the yml template to.", Value: "./issue.yml"},
				cli.StringFlag{Name: "project, p", Usage: "The Jira project to create the issue on", Value: "RL"},
			},
			Action: func(c *cli.Context) {
				var repoPath = c.Args().Get(0)
				var from = c.Args().Get(1)
				var to = c.Args().Get(2)
				var filepath = c.String("output")
				var project = c.String("project")
				releasenotes(repoPath, from, to, filepath, project)
			},
		},
		{
			Name:      "releaseversion",
			Aliases:   []string{"rv"},
			Usage:     "Obtains the last semver tag in the git history",
			ArgsUsage: "<path>",
			Action: func(c *cli.Context) {
				var repoPath = c.Args().Get(0)
				releaseversion(repoPath)
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
