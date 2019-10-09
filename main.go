package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "PAT"
	app.Usage = "CLI Tools for pipelines."
	app.Author = "Guestlogix"
	app.Version = "0.0.2"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "release",
			Aliases: []string{"r"},
			Usage:   "Tools to aid in a release.",
			Subcommands: []cli.Command{
				{
					Name:      "notes",
					Aliases:   []string{"rn"},
					Usage:     "Generates the release notes between two tags",
					ArgsUsage: "<path> <from> <to>",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "output, o", Usage: "The filepath and file name to output the yml template to.", Value: "./issue.yml"},
						cli.StringFlag{Name: "project, p", Usage: "The Jira project to create the issue on", Value: "RL"},
					},
					Action: func(c *cli.Context) {
						repo := openRepo(c.Args().Get(0))
						sinceTag := semverFromString(c.Args().Get(1))
						untilTag := semverFromString(c.Args().Get(2))
						sinceRef := tagRef(repo, sinceTag)
						untilRef := tagRef(repo, untilTag)
						commits := commitsBetweenRefs(repo, sinceRef, untilRef)
						notes := formatReleaseNotes(commits, sinceRef, untilRef, untilTag, c.String("project"))
						writeFile(notes, c.String("output"))
					},
				},
				{
					Name:      "version",
					Aliases:   []string{"rv"},
					Usage:     "Obtains the last semver tag in the git history",
					ArgsUsage: "<path>",
					Action: func(c *cli.Context) {
						repo := openRepo(c.Args().Get(0))
						tag, _ := latestSemverTag(repo)
						fmt.Println(tag.toString())
					},
				},
				{
					Name:      "next",
					Aliases:   []string{"nv"},
					Usage:     "Returns the next version number based on commit names",
					ArgsUsage: "<path>",
					Action: func(c *cli.Context) {
						repo := openRepo(c.Args().Get(0))
						currentSemver, _ := latestSemverTag(repo)
						sinceRef := tagRef(repo, currentSemver)
						commits := commitsBetweenRefs(repo, sinceRef, nil)
						newSemver := semverBump(commits, currentSemver)
						fmt.Println(newSemver.toString())
					},
				},
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
