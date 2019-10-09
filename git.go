package main

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/revlist"
)

// Returns a go-git tag object should a specified tag exist within
// the specified repository
func tagRef(r *git.Repository, s semver) *plumbing.Reference {
	ref, err := r.Tag(s.toString())
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
