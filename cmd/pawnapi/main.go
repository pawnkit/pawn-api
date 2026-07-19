// Command pawnapi maintains and compares Pawn API datasets.
package main

import (
	"fmt"
	"os"
)

const (
	exitOK          = 0
	exitFindings    = 1
	exitInvalidArgs = 2
	exitEnvironment = 3
	exitInternal    = 4
)

var version = "dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage(os.Stderr)
		return exitInvalidArgs
	}

	switch args[0] {
	case "-h", "--help", "help":
		usage(os.Stdout)
		return exitOK
	case "-v", "--version", "version":
		fmt.Println("pawnapi", version)
		return exitOK
	case "validate":
		return runValidate(args[1:])
	case "generate":
		return runGenerate(args[1:])
	case "snapshot":
		return runSnapshot(args[1:])
	case "diff":
		return runDiff(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "pawnapi: unknown command %q\n\n", args[0])
		usage(os.Stderr)
		return exitInvalidArgs
	}
}

func usage(w *os.File) {
	_, _ = fmt.Fprint(w, `pawnapi - SA-MP/open.mp Pawn API metadata tool

Usage:
  pawnapi validate [--root <dir>] [--output human|json]
  pawnapi generate [--root <dir>] [--check]
  pawnapi snapshot [flags] <entries.json|include.inc> <snapshot.json>
  pawnapi diff <old.json> <new.json>
  pawnapi version
  pawnapi help

Commands:
  validate   Validate data/source/*.json: schema conformance, duplicate
             ids/names, broken deprecation references, version ranges.
  generate   Regenerate the embedded Go dataset and the schema-conformant
             interchange JSON from data/source/*.json.
  snapshot   Validate and record a deterministic API snapshot.
  diff       Classify differences between two snapshots as source-compatible,
             potentially breaking, or breaking.
`)
}
