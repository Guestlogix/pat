package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

// CheckIfError should be used to naively panics if an error is not nil.
func checkIfError(err error) {
	if err == nil {
		return
	}
	throwError(err.Error())
}

// Helper function to write byte array to file on disk
func writeFile(contents []byte, filename string) {
	err := ioutil.WriteFile(filename, contents, 0644)
	checkIfError(err)
}

func throwError(msg string) {
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", msg))
	os.Exit(1)
}

// Warning should be used to display a warning
func warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Returns a go-git repo object should a repository
// exist at the specified filepath.
func openRepo(path string) *git.Repository {
	repo, err := git.PlainOpen(path)
	checkIfError(err)
	return repo
}
