package pawnapi

import (
	"bytes"
	"reflect"
	"testing"
)

func TestRoundTrip_LoadSerializeLoad(t *testing.T) {
	original := sampleEntries()

	buf1, err := MarshalEntries(original)
	if err != nil {
		t.Fatalf("first marshal: %v", err)
	}

	loaded, err := LoadEntries(bytes.NewReader(buf1))
	if err != nil {
		t.Fatalf("first load: %v", err)
	}

	buf2, err := MarshalEntries(loaded)
	if err != nil {
		t.Fatalf("second marshal: %v", err)
	}

	reloaded, err := LoadEntries(bytes.NewReader(buf2))
	if err != nil {
		t.Fatalf("second load: %v", err)
	}

	if !reflect.DeepEqual(loaded, reloaded) {
		t.Fatalf("round trip changed data:\nfirst:  %+v\nsecond: %+v", loaded, reloaded)
	}
	if string(buf1) != string(buf2) {
		t.Fatal("round trip is not byte-stable")
	}
}

func TestMarshalEntries_Deterministic(t *testing.T) {
	entries := sampleEntries()
	reversed := make([]Entry, len(entries))
	for i, e := range entries {
		reversed[len(entries)-1-i] = e
	}

	buf1, err := MarshalEntries(entries)
	if err != nil {
		t.Fatal(err)
	}
	buf2, err := MarshalEntries(reversed)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf1) != string(buf2) {
		t.Fatal("MarshalEntries output depends on input order")
	}
}

func TestMarshalSchemaDocument_ConformsAndProjects(t *testing.T) {
	entries := sampleEntries()
	doc := ToSchemaDocument(entries)
	if doc.SchemaVersion != SchemaVersion {
		t.Fatalf("got schemaVersion %d, want %d", doc.SchemaVersion, SchemaVersion)
	}
	if len(doc.Entries) != len(entries) {
		t.Fatalf("got %d schema entries, want %d", len(doc.Entries), len(entries))
	}
	for _, se := range doc.Entries {
		if se.ID == "" || se.Kind == "" || se.Name == "" {
			t.Fatalf("schema entry missing required field: %+v", se)
		}
	}
}
