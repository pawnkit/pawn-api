package schema

import (
	"os"
	"testing"
)

func TestCompiled(t *testing.T) {
	sch, err := Compiled()
	if err != nil {
		t.Fatalf("Compiled() returned an error: %v", err)
	}
	if sch == nil {
		t.Fatal("Compiled() returned a nil schema with no error")
	}
}

func TestValidateDocument_ValidExample(t *testing.T) {
	raw, err := os.ReadFile("../../data/generated/pawn-api.json")
	if err != nil {
		t.Skipf("generated interchange document not present, run `pawnapi generate` first: %v", err)
	}
	if err := ValidateDocument(raw); err != nil {
		t.Fatalf("generated interchange document does not conform to the vendored schema: %v", err)
	}
}

func TestValidateDocument_RejectsMissingRequiredField(t *testing.T) {
	bad := []byte(`{"schemaVersion": 1, "entries": [{"kind": "native", "name": "Foo"}]}`)
	if err := ValidateDocument(bad); err == nil {
		t.Fatal("expected an error for an entry missing required fields (id, availability, source)")
	}
}

func TestValidateDocument_RejectsAdditionalProperties(t *testing.T) {
	bad := []byte(`{
		"schemaVersion": 1,
		"entries": [{
			"id": "native:Foo",
			"kind": "native",
			"name": "Foo",
			"notAllowed": true,
			"availability": [{"profile": "openmp", "since": "1.0.0"}],
			"source": {
				"repository": "openmultiplayer/omp-stdlib",
				"path": "fixture.inc",
				"commit": "0000000000000000000000000000000000000f",
				"license": "MPL-2.0"
			}
		}]
	}`)
	if err := ValidateDocument(bad); err == nil {
		t.Fatal("expected an error for an unknown property (schema uses additionalProperties: false)")
	}
}

func TestValidateDocument_RejectsWrongSchemaVersion(t *testing.T) {
	bad := []byte(`{"schemaVersion": 2, "entries": []}`)
	if err := ValidateDocument(bad); err == nil {
		t.Fatal("expected an error for schemaVersion != 1")
	}
}
