package pawnapi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionPattern accepts Pawn's non-semver version strings.
var versionPattern = regexp.MustCompile(`^\d+(\.\d+){0,3}([+-][A-Za-z0-9.-]+)?$`)

func validateVersion(v string) error {
	if !versionPattern.MatchString(v) {
		return fmt.Errorf("%q does not look like a dotted version", v)
	}
	return nil
}

// compareVersions compares dotted numeric versions without suffixes.
func compareVersions(a, b string) int {
	pa := leadingNumericParts(a)
	pb := leadingNumericParts(b)

	for i := 0; i < len(pa) || i < len(pb); i++ {
		var x, y int
		if i < len(pa) {
			x = pa[i]
		}
		if i < len(pb) {
			y = pb[i]
		}
		if x != y {
			if x < y {
				return -1
			}
			return 1
		}
	}
	return 0
}

func leadingNumericParts(v string) []int {
	core, _, _ := strings.Cut(v, "+")
	core, _, _ = strings.Cut(core, "-")
	fields := strings.Split(core, ".")
	parts := make([]int, 0, len(fields))
	for _, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil {
			break
		}
		parts = append(parts, n)
	}
	return parts
}
