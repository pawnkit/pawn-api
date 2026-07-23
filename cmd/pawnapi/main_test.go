package main

import (
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

func TestRun_ValidateOK(t *testing.T) {
	code := run([]string{"validate", "--root", repoRoot(t)})
	if code != exitOK {
		t.Fatalf("got exit code %d, want %d", code, exitOK)
	}
}

func TestRun_ValidateJSONOutput(t *testing.T) {
	code := run([]string{"validate", "--root", repoRoot(t), "--output", "json"})
	if code != exitOK {
		t.Fatalf("got exit code %d, want %d", code, exitOK)
	}
}

func TestRun_ValidateBadRoot(t *testing.T) {
	code := run([]string{"validate", "--root", filepath.Join(t.TempDir(), "does-not-exist")})
	if code != exitEnvironment {
		t.Fatalf("got exit code %d, want %d", code, exitEnvironment)
	}
}

func TestRun_GenerateCheckUpToDate(t *testing.T) {
	code := run([]string{"generate", "--root", repoRoot(t), "--check"})
	if code != exitOK {
		t.Fatalf("got exit code %d, want %d (generated files should already be up to date; run `pawnapi generate` if this fails)", code, exitOK)
	}
}

func TestRun_GenerateIntoTempDir(t *testing.T) {
	root := repoRoot(t)
	tmp := t.TempDir()

	// Generate needs data/source under --root; point it at the real one.
	if err := os.MkdirAll(filepath.Join(tmp, "data"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "data", "source"), filepath.Join(tmp, "data", "source")); err != nil {
		t.Skip("symlinks unavailable in this environment")
	}

	code := run([]string{"generate", "--root", tmp})
	if code != exitOK {
		t.Fatalf("got exit code %d, want %d", code, exitOK)
	}

	if _, err := os.Stat(filepath.Join(tmp, "pawnapi", "data", "pawn-api.full.json")); err != nil {
		t.Errorf("expected full-model output to be written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "data", "generated", "pawn-api.json")); err != nil {
		t.Errorf("expected interchange output to be written: %v", err)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	code := run([]string{"bogus"})
	if code != exitInvalidArgs {
		t.Fatalf("got exit code %d, want %d", code, exitInvalidArgs)
	}
}

func TestRun_NoArgs(t *testing.T) {
	code := run(nil)
	if code != exitInvalidArgs {
		t.Fatalf("got exit code %d, want %d", code, exitInvalidArgs)
	}
}

func TestRun_Version(t *testing.T) {
	if code := run([]string{"version"}); code != exitOK {
		t.Fatalf("got exit code %d, want %d", code, exitOK)
	}
}

func TestRun_Help(t *testing.T) {
	if code := run([]string{"help"}); code != exitOK {
		t.Fatalf("got exit code %d, want %d", code, exitOK)
	}
}

func TestRun_SnapshotAndDiff(t *testing.T) {
	input := filepath.Join(repoRoot(t), "pawnapi", "data", "pawn-api.full.json")
	snapshot := filepath.Join(t.TempDir(), "snapshot.json")
	if code := run([]string{"snapshot", input, snapshot}); code != exitOK {
		t.Fatalf("snapshot exit code = %d", code)
	}
	if code := run([]string{"diff", snapshot, snapshot}); code != exitOK {
		t.Fatalf("diff exit code = %d", code)
	}
}

func TestRun_SnapshotInclude(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "api.inc")
	output := filepath.Join(dir, "snapshot.json")
	if err := os.WriteFile(input, []byte("native GetValue();\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"snapshot", input, output}); code != exitOK {
		t.Fatalf("snapshot exit code = %d", code)
	}
}

func TestCompareCoverage(t *testing.T) {
	dataset := []pawnapi.Entry{{ID: "native:GetValue"}}
	imported := map[string]coverageItem{
		"native:GetValue": {ID: "native:GetValue", Path: "core.inc"},
		"native:SetValue": {ID: "native:SetValue", Path: "core.inc"},
	}
	report := compareCoverage(dataset, imported)
	if report.Covered != 1 || report.Total != 2 || len(report.Missing) != 1 || report.Missing[0].ID != "native:SetValue" {
		t.Fatalf("report = %+v", report)
	}
}
