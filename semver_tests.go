package main

import "testing"

// TestSemverFromInt - assert correct semver struct
func TestSemverFromInt(t *testing.T) {
	s := semverFromInt(7, 1, 2)
	vers := s.toString()
	if vers != "v7.1.2" {
		t.Errorf("String value incorrect, got: %d, want: %d.", vers, "v7.1.2")
	}
}
