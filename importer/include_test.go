package importer

import (
	"testing"

	"github.com/pawnkit/pawn-api/pawnapi"
)

func TestIncludeExtractsCallables(t *testing.T) {
	text := []byte("#define LIMIT 10\nenum PlayerState { State_None }\nnative Float:GetHealth(playerid);\nforward OnPlayerConnect(playerid);\nstock Internal() {}\n")
	entries, err := Include(text, Options{Repository: "example/repo", Path: "api.inc", Commit: "1234567", License: "MIT"})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 4 {
		t.Fatalf("entries = %+v", entries)
	}
	byID := make(map[string]pawnapi.Entry)
	for _, item := range entries {
		byID[item.ID] = item
	}
	if byID["native:GetHealth"].Signature.ReturnTag != "Float" {
		t.Fatalf("native = %+v", byID["native:GetHealth"])
	}
	if byID["callback:OnPlayerConnect"].Kind != pawnapi.KindCallback || byID["tag:PlayerState"].Kind != pawnapi.KindTag {
		t.Fatalf("entries = %+v", entries)
	}
	if byID["define:LIMIT"].Value.String() != "10" {
		t.Fatalf("define = %+v", byID["define:LIMIT"])
	}
}

func TestIncludeRequiresProvenance(t *testing.T) {
	_, err := Include([]byte("native GetValue();"), Options{Path: "api.inc"})
	if err == nil {
		t.Fatal("missing provenance accepted")
	}
}
