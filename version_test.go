package main

import (
	"regexp"
	"testing"
)

func TestVersionFormat(t *testing.T) {
	if !regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`).MatchString(Version) {
		t.Fatalf("Version %q must be like 1.2.3 so the release workflow can parse it", Version)
	}
}
