package pawnapi

import (
	"fmt"
	"sort"
)

func sortedByKindAndID(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return out[i].Kind < out[j].Kind
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// ValidateDataset checks entries, uniqueness, and deprecation references.
// It returns every problem found.
func ValidateDataset(entries []Entry) error {
	var errs ValidationErrors

	byID := make(map[string]int, len(entries))
	byKindName := make(map[string]int, len(entries))

	for i, e := range entries {
		if err := e.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("entry %d (%s): %w", i, entryLabel(e), err))
			continue
		}

		if first, dup := byID[e.ID]; dup {
			errs = append(errs, fmt.Errorf("entry %d (%s): %w: also entry %d", i, e.ID, ErrDuplicateID, first))
		} else {
			byID[e.ID] = i
		}

		kn := string(e.Kind) + ":" + e.Name
		if first, dup := byKindName[kn]; dup {
			errs = append(errs, fmt.Errorf("entry %d (%s): %w: also entry %d", i, e.ID, ErrDuplicateName, first))
		} else {
			byKindName[kn] = i
		}
	}

	for i, e := range entries {
		if e.Deprecated == nil || e.Deprecated.Replacement == "" {
			continue
		}
		if !idPattern.MatchString(e.Deprecated.Replacement) {
			continue // free-text reason, not an id reference
		}
		if _, ok := byID[e.Deprecated.Replacement]; !ok {
			errs = append(errs, fmt.Errorf("entry %d (%s): %w: %q", i, e.ID, ErrBrokenReplacement, e.Deprecated.Replacement))
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func entryLabel(e Entry) string {
	if e.ID != "" {
		return e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return "<unnamed>"
}

// ValidationErrors collects every problem [ValidateDataset] found.
type ValidationErrors []error

func (v ValidationErrors) Error() string {
	if len(v) == 1 {
		return v[0].Error()
	}
	s := fmt.Sprintf("%d validation errors:", len(v))
	for _, e := range v {
		s += "\n  - " + e.Error()
	}
	return s
}

func (v ValidationErrors) Unwrap() []error {
	return []error(v)
}
