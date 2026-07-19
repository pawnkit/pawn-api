package pawnapi

import "testing"

func TestDiffClassifiesChanges(t *testing.T) {
	base := Entry{ID: "native:One", Kind: KindNative, Name: "One", Signature: &Signature{}, DocumentationSummary: "old"}
	removed := Diff([]Entry{base}, nil)
	if len(removed) != 1 || removed[0].Class != ChangeBreaking {
		t.Fatalf("removed = %+v", removed)
	}
	added := Diff(nil, []Entry{base})
	if len(added) != 1 || added[0].Class != ChangeCompatible {
		t.Fatalf("added = %+v", added)
	}
	updated := base
	updated.DocumentationSummary = "new"
	metadata := Diff([]Entry{base}, []Entry{updated})
	if len(metadata) != 1 || metadata[0].Class != ChangeCompatible {
		t.Fatalf("metadata = %+v", metadata)
	}
	updated.Signature = &Signature{ReturnTag: "Float"}
	signature := Diff([]Entry{base}, []Entry{updated})
	if len(signature) != 1 || signature[0].Class != ChangeBreaking {
		t.Fatalf("signature = %+v", signature)
	}
}
