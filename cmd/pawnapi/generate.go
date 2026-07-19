package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/pawnkit/pawn-api/internal/generator"
)

func runGenerate(args []string) int {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	root := fs.String("root", ".", "repository root containing data/source")
	check := fs.Bool("check", false, "verify generated files are up to date without writing them")
	if err := fs.Parse(args); err != nil {
		return exitInvalidArgs
	}

	layout := generator.DefaultLayout(*root)
	result, err := generator.Generate(layout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi generate:", err)
		return exitEnvironment
	}

	if *check {
		stale, err := checkUpToDate(layout, result)
		if err != nil {
			fmt.Fprintln(os.Stderr, "pawnapi generate --check:", err)
			return exitEnvironment
		}
		if len(stale) > 0 {
			fmt.Fprintln(os.Stderr, "pawnapi generate --check: the following files are stale, run `pawnapi generate`:")
			for _, f := range stale {
				fmt.Fprintln(os.Stderr, "  "+f)
			}
			return exitFindings
		}
		fmt.Println("ok: generated files are up to date")
		return exitOK
	}

	if err := result.Write(layout); err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi generate:", err)
		return exitEnvironment
	}

	fmt.Printf("ok: generated %d entries (%d bytes full, %d bytes interchange)\n",
		result.Manifest.EntryCount, len(result.Full), len(result.Interchange))
	return exitOK
}

func checkUpToDate(layout generator.Layout, result *generator.Result) ([]string, error) {
	pairs := []struct {
		path string
		want []byte
	}{
		{layout.FullOutput, result.Full},
		{layout.InterchangeOutput, result.Interchange},
		{layout.ManifestOutput, result.ManifestJSON},
	}

	var stale []string
	for _, p := range pairs {
		got, err := os.ReadFile(p.path)
		if err != nil {
			if os.IsNotExist(err) {
				stale = append(stale, p.path)
				continue
			}
			return nil, err
		}
		if !bytes.Equal(got, p.want) {
			stale = append(stale, p.path)
		}
	}
	return stale, nil
}
