package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// SemverRegex is regex used to parse a string representation of a semver
const semverRegex = `v([0-9]+).([0-9]+).([0-9]+)`

// semver represents a semantic version that can be used to identify
// the version of a given commit hash.
type semver struct {
	major int
	minor int
	patch int
}

// semverFromInt constructs a new semver from int values
func semverFromInt(major int, minor int, patch int) semver {
	return semver{major, minor, patch}
}

// semverFromString constructs a new semver from string value
func semverFromString(s string) semver {
	_, err := isSemver(s)
	checkIfError(err)
	r, err := regexp.Compile(semverRegex)
	major, err := strconv.Atoi(r.FindStringSubmatch(s)[1])
	checkIfError(err)
	minor, err := strconv.Atoi(r.FindStringSubmatch(s)[2])
	checkIfError(err)
	patch, err := strconv.Atoi(r.FindStringSubmatch(s)[3])
	checkIfError(err)
	return semver{major, minor, patch}
}

// toString converts a semver struct to a string
func (s *semver) toString() string {
	return fmt.Sprintf("v%d.%d.%d", s.major, s.minor, s.patch)
}

// incrementMajor increments the major column of the version
func (s *semver) incrementMajor() semver {
	return semver{s.major + 1, 0, 0}
}

// incrementMinor increments the minor column of the version
func (s *semver) incrementMinor() semver {
	return semver{s.major, s.minor + 1, 0}
}

// incrementPatch increments the patch column of the version
func (s *semver) incrementPatch() semver {
	return semver{s.major, s.minor, s.patch + 1}
}

// Determine if provided string represents a valid semver
func isSemver(s string) (bool, error) {
	r, err := regexp.Compile(semverRegex)
	if err != nil {
		return false, errors.New("Invalid semver, must be of format `vx.y.z`")
	}
	if r.MatchString(s) {
		return true, nil
	}
	return false, errors.New("Invalid semver, must be of format `vx.y.z`")
}
