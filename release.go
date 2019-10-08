package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

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

// Return only the semver portion of the tag
func parseSemver(tag string) string {
	_, err := isSemver(tag)
	checkIfError(err)

	semverRegex, _ := regexp.Compile(`v([0-9]+\.[0-9]+\.[0-9]+)`)
	return semverRegex.FindStringSubmatch(tag)[0]
}

// Determine if provided string tag contains a valid semver
func isSemver(tag string) (bool, error) {
	semverRegex, err := regexp.Compile(`v([0-9]+\.[0-9]+\.[0-9]+)`)
	if err != nil {
		return false, errors.New("Invalid semver regex")
	}
	if semverRegex.MatchString(tag) {
		return true, nil
	}
	return false, errors.New("Invalid semver regex")
}

func toSemver(major int, minor int, patch int) string {
	return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
}

// Return the int representations of the major, minor and patch values
func semverValues(tag string) (int, int, int) {
	_, err := isSemver(tag)
	checkIfError(err)

	semverRegex, _ := regexp.Compile(`v([0-9]+).([0-9]+).([0-9]+)`)
	major, err := strconv.Atoi(semverRegex.FindStringSubmatch(tag)[1])
	checkIfError(err)
	minor, err := strconv.Atoi(semverRegex.FindStringSubmatch(tag)[2])
	checkIfError(err)
	patch, err := strconv.Atoi(semverRegex.FindStringSubmatch(tag)[3])
	checkIfError(err)
	return major, minor, patch
}

// Returns the last semver tag in the git history
func latestSemverTag(repo *git.Repository) (string, error) {
	tagRefs, err := repo.Tags()
	checkIfError(err)

	var latestTagCommit *object.Commit
	var latestTagName string
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

// Returns the semantic version value given commits
// and major, minor and patch values of current version
func semverBump(commits []*object.Commit, major int, minor int, patch int) (int, int, int) {
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
		return major, minor + 1, 0
	} else if patchFlag {
		return major, minor, patch + 1
	} else {
		return major, minor, patch
	}
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

// Outputs the new semver tag if one is required, otherwise
func releaseversion(repoPath string) {
	repo := openRepo(repoPath)
	tag, _ := latestSemverTag(repo)
	fmt.Println(tag)
}

// Returns the newest Semantic Version given the current one,
// and an array of commits
func nextversion(repoPath string) {
	repo := openRepo(repoPath)
	tag, _ := latestSemverTag(repo)
	// fmt.Println("Current Version: ", tag)
	sinceRef := tagRef(repo, tag)
	commits := commitsBetweenRefs(repo, sinceRef, nil)
	major, minor, patch := semverValues(tag)
	mj, mn, pt := semverBump(commits, major, minor, patch)
	output := toSemver(mj, mn, pt)
	// fmt.Println("New Version: ", output)
	fmt.Println(output)
}
