package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/revlist"
	"gopkg.in/yaml.v2"
)

// YML Structure of a new Jira Issue
type issueTemplateStruct struct {
	Fields struct {
		Project struct {
			Key string
		}
		Issuetype struct {
			Name string
		}
		Summary     string
		Description string
	}
}

// Helper function to write byte array to file on disk
func writeFile(contents []byte, filename string) {
	err := ioutil.WriteFile(filename, contents, 0644)
	checkIfError(err)
}

// Returns a go-git tag object should a specified tag exist within
// the specified repository
func tagRef(r *git.Repository, tag string) *plumbing.Reference {
	ref, err := r.Tag(tag)
	checkIfError(err)
	return ref
}

// Uses git rev-list to determine all the commits between two
// specified commit references in the specified repository
func commitsBetweenRefs(repo *git.Repository, since *plumbing.Reference, until *plumbing.Reference) []*object.Commit {

	commits := make([]*object.Commit, 0)

	// throw error if no 'since' tag is specified
	if since == nil {
		throwError("Since tag must be specified")
	}

	// if no until is specified, use HEAD
	if until == nil {
		head, err := repo.Head()
		checkIfError(err)
		until = head
	}

	ref1hist, err := revlist.Objects(repo.Storer, []plumbing.Hash{since.Hash()}, nil)
	checkIfError(err)
	ref2hist, err := revlist.Objects(repo.Storer, []plumbing.Hash{until.Hash()}, ref1hist)
	checkIfError(err)

	for _, h := range ref2hist {
		c, err := repo.CommitObject(h)
		if err != nil {
			continue
		}
		commits = append(commits, c)
	}
	return commits
}

// Generate a formatted release notes given an
// array of all commits incorporated as well as the
// since and until refs
func formatReleaseNotes(commits []*object.Commit, since *plumbing.Reference, until *plumbing.Reference, to string, project string) []byte {
	var sb strings.Builder
	sb.WriteString("Release Notes \n")
	fmt.Fprintf(&sb, "Notes Since: %s \n", since.Hash())
	fmt.Fprintf(&sb, "Notes Until: %s (%s)\n", until.Hash(), to)
	sb.WriteString("----------------------------------- \n")

	for _, c := range commits {
		fmt.Fprintf(&sb, "%s \n", c.Message)
	}

	issueTemplate := issueTemplateStruct{}
	issueTemplate.Fields.Issuetype.Name = "Task"
	issueTemplate.Fields.Project.Key = project
	issueTemplate.Fields.Summary = to
	issueTemplate.Fields.Description = sb.String()

	d, err := yaml.Marshal(&issueTemplate)
	checkIfError(err)

	return d
}

// Outputs release notes between two specified tags to stdout
func releasenotes(repoPath string, sinceTag string, untilTag string, filepath string, project string) {
	repo := openRepo(repoPath)
	sinceRef := tagRef(repo, sinceTag)
	untilRef := tagRef(repo, untilTag)
	commits := commitsBetweenRefs(repo, sinceRef, untilRef)
	notes := formatReleaseNotes(commits, sinceRef, untilRef, untilTag, project)
	writeFile(notes, filepath)
}
