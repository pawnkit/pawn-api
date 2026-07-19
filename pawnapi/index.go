package pawnapi

import "sort"

// Index is a validated, queryable API dataset.
type Index struct {
	entries []Entry
	byID    map[string]Entry
	byName  map[string][]Entry // name -> entries sharing that name, any kind
}

// NewIndex validates and copies entries into a sorted index.
func NewIndex(entries []Entry) (*Index, error) {
	if err := ValidateDataset(entries); err != nil {
		return nil, err
	}

	sorted := sortedByKindAndID(entries)
	ix := &Index{
		entries: sorted,
		byID:    make(map[string]Entry, len(sorted)),
		byName:  make(map[string][]Entry, len(sorted)),
	}
	for _, e := range sorted {
		ix.byID[e.ID] = e
		ix.byName[e.Name] = append(ix.byName[e.Name], e)
		for _, alias := range e.Aliases {
			if alias != e.Name {
				ix.byName[alias] = append(ix.byName[alias], e)
			}
		}
	}
	return ix, nil
}

// Len returns the number of entries in the index.
func (ix *Index) Len() int {
	return len(ix.entries)
}

// All returns every entry sorted by kind and ID. Do not mutate the result.
func (ix *Index) All() []Entry {
	return ix.entries
}

// ByID looks up a single entry by its stable id (e.g. "native:SetPlayerPos").
func (ix *Index) ByID(id string) (Entry, bool) {
	e, ok := ix.byID[id]
	return e, ok
}

// ByName returns entries matching a primary name or alias.
func (ix *Index) ByName(name string) []Entry {
	return ix.byName[name]
}

// ByKindName looks up the single entry with the given kind and name.
func (ix *Index) ByKindName(kind Kind, name string) (Entry, bool) {
	for _, e := range ix.byName[name] {
		if e.Kind == kind {
			return e, true
		}
	}
	return Entry{}, false
}

// ByKind returns every entry of the given kind, sorted by id.
func (ix *Index) ByKind(kind Kind) []Entry {
	var out []Entry
	for _, e := range ix.entries {
		if e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}

// ByProfile returns entries available under profile, sorted by kind and ID.
func (ix *Index) ByProfile(profile string) []Entry {
	var out []Entry
	for _, e := range ix.entries {
		for _, a := range e.Availability {
			if a.Profile == profile {
				out = append(out, e)
				break
			}
		}
	}
	return out
}

// Deprecated returns deprecated entries sorted by kind and ID.
func (ix *Index) Deprecated() []Entry {
	var out []Entry
	for _, e := range ix.entries {
		if e.Deprecated != nil {
			out = append(out, e)
		}
	}
	return out
}

// Profiles returns distinct profile names in alphabetical order.
func (ix *Index) Profiles() []string {
	seen := map[string]bool{}
	for _, e := range ix.entries {
		for _, a := range e.Availability {
			seen[a.Profile] = true
		}
	}
	out := make([]string, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}
