package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pawnkit/pawn-api/importer"
	"github.com/pawnkit/pawn-api/pawnapi"
)

type coverageItem struct {
	ID   string
	Path string
}

type coverageReport struct {
	Covered int
	Total   int
	Missing []coverageItem
}

func runCoverage(args []string) int {
	flags := flag.NewFlagSet("coverage", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	profile := flags.String("profile", pawnapi.ProfileOpenMP, "availability profile")
	repository := flags.String("repository", "local", "source repository")
	commit := flags.String("commit", pawnapi.HandAuthoredCommit, "source commit")
	license := flags.String("license", "NOASSERTION", "SPDX licence")
	limit := flags.Int("limit", 100, "maximum missing entries to print; 0 prints all")
	if err := flags.Parse(args); err != nil {
		return exitInvalidArgs
	}
	if *limit < 0 {
		fmt.Fprintln(os.Stderr, "pawnapi coverage: limit must not be negative")
		return exitInvalidArgs
	}
	if flags.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "usage: pawnapi coverage [flags] <dataset.json> <include.inc>...")
		return exitInvalidArgs
	}

	datasetFile, err := os.Open(flags.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi coverage:", err)
		return exitEnvironment
	}
	dataset, err := pawnapi.LoadEntries(datasetFile)
	closeErr := datasetFile.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi coverage:", err)
		return exitFindings
	}

	imported := make(map[string]coverageItem)
	for _, path := range flags.Args()[1:] {
		text, readErr := os.ReadFile(path)
		if readErr != nil {
			fmt.Fprintln(os.Stderr, "pawnapi coverage:", readErr)
			return exitEnvironment
		}
		entries, importErr := importer.Include(text, importer.Options{
			Profile: *profile, Repository: *repository, Path: filepath.Base(path), Commit: *commit, License: *license,
		})
		if importErr != nil {
			fmt.Fprintf(os.Stderr, "pawnapi coverage: %s: %v\n", path, importErr)
			return exitFindings
		}
		for _, entry := range entries {
			imported[entry.ID] = coverageItem{ID: entry.ID, Path: path}
		}
	}

	report := compareCoverage(dataset, imported)
	fmt.Printf("%d of %d declarations covered\n", report.Covered, report.Total)
	visible := report.Missing
	if *limit != 0 && len(visible) > *limit {
		visible = visible[:*limit]
	}
	for _, item := range visible {
		fmt.Printf("missing %s (%s)\n", item.ID, item.Path)
	}
	if hidden := len(report.Missing) - len(visible); hidden > 0 {
		fmt.Printf("%d more missing entries; use --limit 0 to print all\n", hidden)
	}
	if len(report.Missing) != 0 {
		return exitFindings
	}
	return exitOK
}

func compareCoverage(dataset []pawnapi.Entry, imported map[string]coverageItem) coverageReport {
	known := make(map[string]bool, len(dataset))
	for _, entry := range dataset {
		known[entry.ID] = true
	}
	report := coverageReport{Total: len(imported)}
	for id, item := range imported {
		if known[id] {
			report.Covered++
			continue
		}
		report.Missing = append(report.Missing, item)
	}
	sort.Slice(report.Missing, func(i, j int) bool {
		if report.Missing[i].ID != report.Missing[j].ID {
			return report.Missing[i].ID < report.Missing[j].ID
		}
		return report.Missing[i].Path < report.Missing[j].Path
	})
	return report
}
