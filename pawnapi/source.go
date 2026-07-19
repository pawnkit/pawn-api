package pawnapi

import (
	"fmt"
	"regexp"
)

var commitPattern = regexp.MustCompile(`^[0-9a-f]{7,40}$`)

// Source records an entry's origin and licence.
type Source struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
	Commit     string `json:"commit"`
	License    string `json:"license"`
}

func (s Source) Validate() error {
	if s.Repository == "" {
		return ErrMissingRepository
	}
	if s.Path == "" {
		return ErrMissingPath
	}
	if s.Commit == "" {
		return ErrMissingCommit
	}
	if !commitPattern.MatchString(s.Commit) {
		return fmt.Errorf("%w: %q", ErrInvalidCommit, s.Commit)
	}
	if s.License == "" {
		return ErrMissingLicense
	}
	return nil
}

// HandAuthoredCommit marks entries that were not taken from a pinned commit.
const HandAuthoredCommit = "0000000000000000000000000000000000000000"
