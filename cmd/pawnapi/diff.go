package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pawnkit/pawn-api/pawnapi"
)

func runDiff(args []string) int {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: pawnapi diff <old.json> <new.json>")
		return exitInvalidArgs
	}
	oldEntries, err := loadSnapshot(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi diff:", err)
		return exitEnvironment
	}
	newEntries, err := loadSnapshot(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi diff:", err)
		return exitEnvironment
	}
	changes := pawnapi.Diff(oldEntries, newEntries)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(changes); err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi diff:", err)
		return exitInternal
	}
	for _, change := range changes {
		if change.Class == pawnapi.ChangeBreaking || change.Class == pawnapi.ChangePotential {
			return exitFindings
		}
	}
	return exitOK
}

func loadSnapshot(path string) ([]pawnapi.Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	entries, err := pawnapi.LoadEntries(file)
	if err != nil {
		return nil, err
	}
	if err := pawnapi.ValidateDataset(entries); err != nil {
		return nil, err
	}
	return entries, nil
}
