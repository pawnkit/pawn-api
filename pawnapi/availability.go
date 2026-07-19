package pawnapi

import (
	"fmt"
	"regexp"
)

var profilePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// Well-known compatibility profiles. Other profile strings are valid.
const (
	ProfileSAMP037 = "samp-037"
	ProfileOpenMP  = "openmp"
	ProfileLegacy  = "legacy"
)

// Availability records the versions in which an entry exists for a profile.
type Availability struct {
	Profile string  `json:"profile"`
	Since   string  `json:"since"`
	Until   *string `json:"until,omitempty"`
}

func (a Availability) Validate() error {
	if !profilePattern.MatchString(a.Profile) {
		return fmt.Errorf("%w: profile %q", ErrInvalidProfile, a.Profile)
	}
	if a.Since != "" {
		if err := validateVersion(a.Since); err != nil {
			return fmt.Errorf("%w: since: %w", ErrInvalidVersion, err)
		}
	}
	if a.Until != nil && *a.Until != "" {
		if err := validateVersion(*a.Until); err != nil {
			return fmt.Errorf("%w: until: %w", ErrInvalidVersion, err)
		}
		if a.Since != "" && compareVersions(a.Since, *a.Until) > 0 {
			return fmt.Errorf("%w: since %q is after until %q", ErrInvalidVersionRange, a.Since, *a.Until)
		}
	}
	return nil
}
