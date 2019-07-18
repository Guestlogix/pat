package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"github.com/waigani/diffparser"
	"golang.org/x/oauth2"
)

// Iterates through all projects in sln and builds out a table
// of all projects with appsettings (rows) and the types of appsettings
// in each project (columns)
func constructAppsettingTable(srcRoot string) map[string]map[string]bool {
	var appsettingTable map[string]map[string]bool
	appsettingTable = make(map[string]map[string]bool)
	appsettingRegex, _ := regexp.Compile("(appsettings.*.json)") // Appsetting regex
	projects, err := ioutil.ReadDir(srcRoot)                     // Get all the .NET Projects
	checkIfError(err)
	for _, project := range projects {
		if project.IsDir() {
			files, err := ioutil.ReadDir(srcRoot + "/" + project.Name()) // Get all the top level files in each project
			checkIfError(err)
			for _, file := range files {
				if appsettingRegex.MatchString(file.Name()) {
					var row map[string]bool
					if appsettingTable[project.Name()] == nil {
						row = make(map[string]bool)
					} else {
						row = appsettingTable[project.Name()]
					}
					row[file.Name()] = false
					appsettingTable[project.Name()] = row
				}
			}
		}
	}
	return appsettingTable
}

// Looks through the git diff and updates the change table accordingly
func populateTable(appsettingTable *map[string]map[string]bool, diffPath string) {
	byt, _ := ioutil.ReadFile("example.diff")
	diff, _ := diffparser.Parse(string(byt))
	r, _ := regexp.Compile("(appsettings.*.json)")

	for _, file := range diff.Files {
		if r.MatchString(file.NewName) {
			//TODO: Regex match the project and appsetting then update the appsetting table accordinly
			fmt.Printf("File: %q -> %q [%v]\n", file.OrigName, file.NewName, file.Mode)
		}
	}
}

// Given a map of projects and appsettings map[string]map[string]bool
// return the equivalent markdown text
func printTable(appsettingTable map[string]map[string]bool) string {
	var markdownTable strings.Builder

	// Header
	markdownTable.WriteString("| Startup | ")
	for _, column := range appsettingTable {
		for appsettingName := range column {
			markdownTable.WriteString(appsettingName)
			markdownTable.WriteString(" | ")
		}
		break
	}
	markdownTable.WriteString("\n")

	// Title Row
	markdownTable.WriteString("|--|")
	for _, column := range appsettingTable {
		for appsettingName := range column {
			markdownTable.WriteString("--|")
			appsettingName = appsettingName + ""
		}
		break
	}
	markdownTable.WriteString("\n")

	// Body
	for row, column := range appsettingTable {
		markdownTable.WriteString("| ")
		markdownTable.WriteString(row)
		markdownTable.WriteString(" | ")
		for _, value := range column {
			if value {
				markdownTable.WriteString(":new:")
			} else {
				markdownTable.WriteString(":new:")
			}
			markdownTable.WriteString(" | ")
		}
		markdownTable.WriteString("\n")
	}

	fmt.Println(markdownTable.String())

	return markdownTable.String()
}

func commentOnPr(message string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "deefe72874879e0a0765dc1fd0c01a60a9fc9958"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Pull Requests are underlying issues, therefore to do a comment on
	// the PR as a whole, we use the Issue API.
	_, _, err := client.Issues.CreateComment(
		ctx,
		"Guestlogix",
		"fifa-ladder",
		1,
		&github.IssueComment{Body: github.String(message)})

	checkIfError(err)
}

func appsettings(repo string) {
	appsettingTable := constructAppsettingTable(repo + "/src")
	populateTable(&appsettingTable, repo)
	//printTable(updatedAppsettingTable)
	//commentOnPr("Testing")
}
