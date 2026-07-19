package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/pawnkit/pawn-api/pawnapi"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(wd, "..", "..")
}

func TestLoadSource_RepositoryData(t *testing.T) {
	entries, err := LoadSource(filepath.Join(repoRoot(t), "data", "source"))
	if err != nil {
		t.Fatalf("LoadSource: %v", err)
	}
	if len(entries) < 40 {
		t.Fatalf("got %d entries, want at least 40 (GOAL.md seed dataset size)", len(entries))
	}
	if err := pawnapi.ValidateDataset(entries); err != nil {
		t.Fatalf("checked-in source data failed validation: %v", err)
	}
}

func TestGenerate_Deterministic(t *testing.T) {
	layout := DefaultLayout(repoRoot(t))

	r1, err := Generate(layout)
	if err != nil {
		t.Fatalf("first generate: %v", err)
	}
	r2, err := Generate(layout)
	if err != nil {
		t.Fatalf("second generate: %v", err)
	}

	if string(r1.Full) != string(r2.Full) {
		t.Error("full-model output is not byte-identical across two Generate calls")
	}
	if string(r1.Interchange) != string(r2.Interchange) {
		t.Error("interchange output is not byte-identical across two Generate calls")
	}
	if string(r1.ManifestJSON) != string(r2.ManifestJSON) {
		t.Error("manifest output is not byte-identical across two Generate calls")
	}
	if r1.Manifest.FullHash != r2.Manifest.FullHash {
		t.Error("full-model content hash is not stable across two Generate calls")
	}
}

func TestGenerate_WriteThenReadBack(t *testing.T) {
	root := repoRoot(t)
	tmp := t.TempDir()

	layout := Layout{
		SourceDir:         filepath.Join(root, "data", "source"),
		FullOutput:        filepath.Join(tmp, "pawnapi", "data", "pawn-api.full.json"),
		InterchangeOutput: filepath.Join(tmp, "data", "generated", "pawn-api.json"),
		ManifestOutput:    filepath.Join(tmp, "data", "generated", "manifest.json"),
	}

	result, err := Generate(layout)
	if err != nil {
		t.Fatal(err)
	}
	if err := result.Write(layout); err != nil {
		t.Fatal(err)
	}

	full, err := os.ReadFile(layout.FullOutput)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := pawnapi.LoadEntries(bytes.NewReader(full))
	if err != nil {
		t.Fatalf("reading back written full-model output: %v", err)
	}
	if len(entries) != result.Manifest.EntryCount {
		t.Fatalf("got %d entries after write/read-back, want %d", len(entries), result.Manifest.EntryCount)
	}
}

func TestGenerate_RejectsInvalidSource(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "bad.json"), []byte(`[{"id": "native:Foo"}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	layout := Layout{SourceDir: tmp}
	if _, err := Generate(layout); err == nil {
		t.Fatal("expected Generate to reject a source entry missing required fields")
	}
}
