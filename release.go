package main

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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

// Generate a formatted release notes given an array of all
// commits incorporated as well as the since and until refs
func formatReleaseNotes(commits []*object.Commit, since *plumbing.Reference, until *plumbing.Reference, to semver, project string) []byte {
	var sb strings.Builder
	sb.WriteString("Release Notes \n")
	fmt.Fprintf(&sb, "Notes Since: %s \n", since.Hash())
	fmt.Fprintf(&sb, "Notes Until: %s (%s)\n", until.Hash(), to.toString())
	sb.WriteString("----------------------------------- \n")

	for _, c := range commits {
		fmt.Fprintf(&sb, "%s \n", c.Message)
	}

	issueTemplate := issueTemplateStruct{}
	issueTemplate.Fields.Issuetype.Name = "Task"
	issueTemplate.Fields.Project.Key = project
	issueTemplate.Fields.Summary = to.toString()
	issueTemplate.Fields.Description = sb.String()

	d, err := yaml.Marshal(&issueTemplate)
	checkIfError(err)

	return d
}

// Returns the last semver tag in the git history
func latestSemverTag(repo *git.Repository) (semver, error) {
	tagRefs, err := repo.Tags()
	checkIfError(err)

	var latestTagCommit *object.Commit
	var latestTagName semver
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		tagName := tagRef.Name().String()
		_, err := isSemver(tagName)
		if err == nil {
			revision := plumbing.Revision(tagName)
			tagCommitHash, err := repo.ResolveRevision(revision)
			checkIfError(err)

			commit, err := repo.CommitObject(*tagCommitHash)
			checkIfError(err)

			if latestTagCommit == nil {
				latestTagCommit = commit
				latestTagName = semverFromString(tagName)
			}

			if commit.Committer.When.After(latestTagCommit.Committer.When) {
				latestTagCommit = commit
				latestTagName = semverFromString(tagName)
			}

		}
		return nil
	})
	checkIfError(err)
	return latestTagName, nil
}

// Returns the semantic version value given commits
// and major, minor and patch values of current version
func semverBump(commits []*object.Commit, s semver) semver {
	patchRegex, err := regexp.Compile(`(build|chore|ci|docs|fix|perf|refactor|revert|style|test)(\([a-z ]+\))?: [\w ]+`)
	checkIfError(err)
	minorRegex, err := regexp.Compile(`(feat|feature)(\([a-z ]+\))?: [\w ]+`)
	checkIfError(err)
	patchFlag := false
	minorFlag := false
	for _, c := range commits {
		if patchRegex.MatchString(c.Message) {
			patchFlag = true
		}
		if minorRegex.MatchString(c.Message) {
			minorFlag = true
		}
	}
	if minorFlag {
		return s.incrementMinor()
	} else if patchFlag {
		return s.incrementPatch()
	} else {
		return s
	}
}
