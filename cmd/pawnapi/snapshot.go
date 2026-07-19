package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pawnkit/pawn-api/importer"
	"github.com/pawnkit/pawn-api/pawnapi"
)

func runSnapshot(args []string) int {
	flags := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	profile := flags.String("profile", pawnapi.ProfileOpenMP, "availability profile")
	repository := flags.String("repository", "local", "source repository")
	commit := flags.String("commit", pawnapi.HandAuthoredCommit, "source commit")
	license := flags.String("license", "NOASSERTION", "SPDX licence")
	if err := flags.Parse(args); err != nil {
		return exitInvalidArgs
	}
	if flags.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: pawnapi snapshot [flags] <input> <snapshot.json>")
		return exitInvalidArgs
	}
	inputPath, outputPath := flags.Arg(0), flags.Arg(1)
	input, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi snapshot:", err)
		return exitEnvironment
	}
	var entries []pawnapi.Entry
	if filepath.Ext(inputPath) == ".inc" {
		content, readErr := os.ReadFile(inputPath)
		if readErr != nil {
			err = readErr
		} else {
			entries, err = importer.Include(content, importer.Options{
				Profile: *profile, Repository: *repository, Path: filepath.Base(inputPath), Commit: *commit, License: *license,
			})
		}
	} else {
		entries, err = pawnapi.LoadEntries(input)
	}
	_ = input.Close()
	if err == nil {
		err = pawnapi.ValidateDataset(entries)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi snapshot:", err)
		return exitFindings
	}
	content, err := pawnapi.MarshalEntries(entries)
	if err == nil {
		err = os.WriteFile(outputPath, content, 0o644)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi snapshot:", err)
		return exitEnvironment
	}
	return exitOK
}
