package pawnapi

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func loadFixture(t *testing.T, name string) []Entry {
	t.Helper()
	path := filepath.Join("..", "testdata", "broken", name)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("opening fixture %s: %v", name, err)
	}
	defer func() { _ = f.Close() }()
	entries, err := LoadEntries(f)
	if err != nil {
		t.Fatalf("decoding fixture %s: %v", name, err)
	}
	return entries
}

func TestValidateDataset_DuplicateID(t *testing.T) {
	err := ValidateDataset(loadFixture(t, "duplicate-id.json"))
	if err == nil {
		t.Fatal("expected an error for duplicate ids")
	}
	if !errors.Is(err, ErrDuplicateID) {
		t.Fatalf("got %v, want an error wrapping ErrDuplicateID", err)
	}
}

func TestValidateDataset_DuplicateName(t *testing.T) {
	err := ValidateDataset(loadFixture(t, "duplicate-name.json"))
	if !errors.Is(err, ErrDuplicateName) {
		t.Fatalf("got %v, want an error wrapping ErrDuplicateName", err)
	}
}

func TestValidateDataset_BrokenReplacement(t *testing.T) {
	err := ValidateDataset(loadFixture(t, "broken-replacement.json"))
	if !errors.Is(err, ErrBrokenReplacement) {
		t.Fatalf("got %v, want an error wrapping ErrBrokenReplacement", err)
	}
}

func TestValidateDataset_BadVersionRange(t *testing.T) {
	err := ValidateDataset(loadFixture(t, "bad-version-range.json"))
	if !errors.Is(err, ErrInvalidVersionRange) {
		t.Fatalf("got %v, want an error wrapping ErrInvalidVersionRange", err)
	}
}

func TestValidateDataset_IDKindMismatch(t *testing.T) {
	err := ValidateDataset(loadFixture(t, "id-kind-mismatch.json"))
	if !errors.Is(err, ErrIDKindMismatch) {
		t.Fatalf("got %v, want an error wrapping ErrIDKindMismatch", err)
	}
}

func TestValidateDataset_ValidResolvedReplacement(t *testing.T) {
	replaced := validEntry()
	replaced.ID = "native:OldName"
	replaced.Name = "OldName"
	replaced.Deprecated = &Deprecation{Since: "1.1.0", Replacement: "native:SetPlayerPos"}

	entries := []Entry{validEntry(), replaced}
	if err := ValidateDataset(entries); err != nil {
		t.Fatalf("expected valid dataset, got: %v", err)
	}
}

func TestValidateDataset_FreeTextReplacementNotTreatedAsID(t *testing.T) {
	e := validEntry()
	e.Deprecated = &Deprecation{Since: "1.1.0", Replacement: "use the Y component instead"}
	if err := ValidateDataset([]Entry{e}); err != nil {
		t.Fatalf("free-text replacement should not be checked as an id reference, got: %v", err)
	}
}

func TestValidationErrors_CollectsAll(t *testing.T) {
	err := ValidateDataset(loadFixture(t, "duplicate-id.json"))
	verrs, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}
	if len(verrs) == 0 {
		t.Fatal("expected at least one collected error")
	}
}
