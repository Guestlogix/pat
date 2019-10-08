package main

import (
	"fmt"
	"regexp"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Return only the semver portion of the tag
func parseSemver(tag string) string {
	semverRegex, _ := regexp.Compile(`v([0-9]+\.[0-9]+\.[0-9]+)`)
	return semverRegex.FindStringSubmatch(tag)[0]
}

// Determine if provided string tag contains a valid semver
func isSemver(tag string) bool {
	semverRegex, _ := regexp.Compile(`v([0-9]+\.[0-9]+\.[0-9]+)`)
	return semverRegex.MatchString(tag)
}

// Returns the last semver tag in the git history
func latestSemverTag(repo *git.Repository) (string, error) {
	tagRefs, err := repo.Tags()
	checkIfError(err)

	var latestTagCommit *object.Commit
	var latestTagName string
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		tagName := tagRef.Name().String()
		if isSemver(tagName) {
			revision := plumbing.Revision(tagName)
			tagCommitHash, err := repo.ResolveRevision(revision)
			checkIfError(err)

			commit, err := repo.CommitObject(*tagCommitHash)
			checkIfError(err)

			if latestTagCommit == nil {
				latestTagCommit = commit
				latestTagName = parseSemver(tagName)
			}

			if commit.Committer.When.After(latestTagCommit.Committer.When) {
				latestTagCommit = commit
				latestTagName = parseSemver(tagName)
			}

		}
		return nil
	})
	checkIfError(err)
	return latestTagName, nil
}

// Outputs the new semver tag if one is required, otherwise
func releaseversion(repoPath string) {
	repo := openRepo(repoPath)
	tag, _ := latestSemverTag(repo)
	fmt.Println(tag)
}
