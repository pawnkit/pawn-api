// Package generator builds pawn-api datasets from reviewed JSON sources.
package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pawnkit/pawn-api/pawnapi"
	"github.com/pawnkit/pawnkit-core/hash"
)

// Layout contains source and output paths.
type Layout struct {
	SourceDir         string // data/source
	FullOutput        string // pawnapi/data/pawn-api.full.json (embedded)
	InterchangeOutput string // data/generated/pawn-api.json (schema-conformant)
	ManifestOutput    string // data/generated/manifest.json
}

// DefaultLayout returns the standard layout rooted at root.
func DefaultLayout(root string) Layout {
	return Layout{
		SourceDir:         filepath.Join(root, "data", "source"),
		FullOutput:        filepath.Join(root, "pawnapi", "data", "pawn-api.full.json"),
		InterchangeOutput: filepath.Join(root, "data", "generated", "pawn-api.json"),
		ManifestOutput:    filepath.Join(root, "data", "generated", "manifest.json"),
	}
}

// LoadSource reads JSON source files in filename order.
func LoadSource(dir string) ([]pawnapi.Entry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("generator: reading source dir %s: %w", dir, err)
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)

	var all []pawnapi.Entry
	for _, name := range names {
		path := filepath.Join(dir, name)
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("generator: opening %s: %w", path, err)
		}
		got, err := pawnapi.LoadEntries(f)
		closeErr := f.Close()
		if err != nil {
			return nil, fmt.Errorf("generator: decoding %s: %w", path, err)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("generator: closing %s: %w", path, closeErr)
		}
		all = append(all, got...)
	}
	return all, nil
}

// Manifest records one generation run.
type Manifest struct {
	SchemaVersion   int            `json:"schemaVersion"`
	EntryCount      int            `json:"entryCount"`
	CountsByKind    map[string]int `json:"countsByKind"`
	FullHash        string         `json:"fullHash"`
	InterchangeHash string         `json:"interchangeHash"`
}

// Result is the outcome of a Generate call.
type Result struct {
	Entries      []pawnapi.Entry
	Full         []byte
	Interchange  []byte
	Manifest     Manifest
	ManifestJSON []byte
}

// Generate validates source data and renders each output in memory.
func Generate(layout Layout) (*Result, error) {
	entries, err := LoadSource(layout.SourceDir)
	if err != nil {
		return nil, err
	}

	if err := pawnapi.ValidateDataset(entries); err != nil {
		return nil, fmt.Errorf("generator: source data failed validation: %w", err)
	}

	full, err := pawnapi.MarshalEntries(entries)
	if err != nil {
		return nil, err
	}

	interchange, err := pawnapi.MarshalSchemaDocument(entries)
	if err != nil {
		return nil, err
	}

	counts := map[string]int{}
	for _, e := range entries {
		counts[string(e.Kind)]++
	}

	manifest := Manifest{
		SchemaVersion:   pawnapi.SchemaVersion,
		EntryCount:      len(entries),
		CountsByKind:    counts,
		FullHash:        hash.Content(full),
		InterchangeHash: hash.Content(interchange),
	}
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("generator: encoding manifest: %w", err)
	}
	manifestJSON = append(manifestJSON, '\n')

	return &Result{
		Entries:      entries,
		Full:         full,
		Interchange:  interchange,
		Manifest:     manifest,
		ManifestJSON: manifestJSON,
	}, nil
}

// Write stores generated artifacts at the configured paths.
func (r *Result) Write(layout Layout) error {
	writes := []struct {
		path string
		data []byte
	}{
		{layout.FullOutput, r.Full},
		{layout.InterchangeOutput, r.Interchange},
		{layout.ManifestOutput, r.ManifestJSON},
	}
	for _, w := range writes {
		if err := os.MkdirAll(filepath.Dir(w.path), 0o755); err != nil {
			return fmt.Errorf("generator: creating directory for %s: %w", w.path, err)
		}
		if err := os.WriteFile(w.path, w.data, 0o644); err != nil {
			return fmt.Errorf("generator: writing %s: %w", w.path, err)
		}
	}
	return nil
}
