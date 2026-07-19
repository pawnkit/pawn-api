package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/pawnkit/pawn-api/internal/generator"
	"github.com/pawnkit/pawn-api/internal/schema"
	"github.com/pawnkit/pawn-api/pawnapi"
)

type validateReport struct {
	OK           bool     `json:"ok"`
	EntryCount   int      `json:"entryCount"`
	SchemaErrors []string `json:"schemaErrors,omitempty"`
	DataErrors   []string `json:"dataErrors,omitempty"`
}

func runValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	root := fs.String("root", ".", "repository root containing data/source")
	output := fs.String("output", "human", "output format: human|json")
	if err := fs.Parse(args); err != nil {
		return exitInvalidArgs
	}

	layout := generator.DefaultLayout(*root)
	entries, err := generator.LoadSource(layout.SourceDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi validate:", err)
		return exitEnvironment
	}

	report := validateReport{EntryCount: len(entries)}

	if err := pawnapi.ValidateDataset(entries); err != nil {
		if verrs, ok := err.(pawnapi.ValidationErrors); ok {
			for _, e := range verrs {
				report.DataErrors = append(report.DataErrors, e.Error())
			}
		} else {
			report.DataErrors = append(report.DataErrors, err.Error())
		}
	}

	interchange, err := pawnapi.MarshalSchemaDocument(entries)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pawnapi validate:", err)
		return exitInternal
	}
	if err := schema.ValidateDocument(interchange); err != nil {
		report.SchemaErrors = append(report.SchemaErrors, err.Error())
	}

	report.OK = len(report.DataErrors) == 0 && len(report.SchemaErrors) == 0

	switch *output {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(report)
	default:
		printHumanReport(report)
	}

	if !report.OK {
		return exitFindings
	}
	return exitOK
}

func printHumanReport(r validateReport) {
	if r.OK {
		fmt.Printf("ok: %d entries validated\n", r.EntryCount)
		return
	}
	fmt.Printf("FAIL: %d entries checked, %d data error(s), %d schema error(s)\n",
		r.EntryCount, len(r.DataErrors), len(r.SchemaErrors))
	for _, e := range r.DataErrors {
		fmt.Println("  data:  ", e)
	}
	for _, e := range r.SchemaErrors {
		fmt.Println("  schema:", e)
	}
}
