package pawnapi

import (
	"reflect"
	"sort"
)

type ChangeClass string

const (
	ChangeCompatible ChangeClass = "source-compatible"
	ChangePotential  ChangeClass = "potentially-breaking"
	ChangeBreaking   ChangeClass = "breaking"
)

type Change struct {
	Class ChangeClass `json:"class"`
	ID    string      `json:"id"`
	Kind  string      `json:"kind"`
}

// Diff classifies API changes by stable entry ID.
func Diff(oldEntries, newEntries []Entry) []Change {
	oldByID := entriesByID(oldEntries)
	newByID := entriesByID(newEntries)
	changes := make([]Change, 0)
	for id, oldEntry := range oldByID {
		newEntry, ok := newByID[id]
		if !ok {
			changes = append(changes, Change{Class: ChangeBreaking, ID: id, Kind: "removed"})
			continue
		}
		if oldEntry.Kind != newEntry.Kind || oldEntry.Name != newEntry.Name || !reflect.DeepEqual(oldEntry.Signature, newEntry.Signature) {
			changes = append(changes, Change{Class: ChangeBreaking, ID: id, Kind: "signature"})
		} else if !reflect.DeepEqual(oldEntry.Value, newEntry.Value) || !reflect.DeepEqual(oldEntry.Availability, newEntry.Availability) {
			changes = append(changes, Change{Class: ChangePotential, ID: id, Kind: "behaviour"})
		} else if !reflect.DeepEqual(oldEntry, newEntry) {
			changes = append(changes, Change{Class: ChangeCompatible, ID: id, Kind: "metadata"})
		}
	}
	for id := range newByID {
		if _, ok := oldByID[id]; !ok {
			changes = append(changes, Change{Class: ChangeCompatible, ID: id, Kind: "added"})
		}
	}
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Class != changes[j].Class {
			return changes[i].Class < changes[j].Class
		}
		return changes[i].ID < changes[j].ID
	})
	return changes
}

func entriesByID(entries []Entry) map[string]Entry {
	result := make(map[string]Entry, len(entries))
	for _, entry := range entries {
		result[entry.ID] = entry
	}
	return result
}
