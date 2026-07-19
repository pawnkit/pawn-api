package pawnapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func marshalIndentNoEscape(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// LoadEntries decodes entries from r without validating them.
func LoadEntries(r io.Reader) ([]Entry, error) {
	var entries []Entry
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&entries); err != nil {
		return nil, fmt.Errorf("pawnapi: decoding entries: %w", err)
	}
	return entries, nil
}

// MarshalEntries returns deterministic, indented JSON sorted by kind and ID.
func MarshalEntries(entries []Entry) ([]byte, error) {
	sorted := sortedByKindAndID(entries)
	buf, err := marshalIndentNoEscape(sorted)
	if err != nil {
		return nil, fmt.Errorf("pawnapi: encoding entries: %w", err)
	}
	return buf, nil
}

// MarshalSchemaDocument returns a deterministic schema document.
func MarshalSchemaDocument(entries []Entry) ([]byte, error) {
	doc := ToSchemaDocument(entries)
	buf, err := marshalIndentNoEscape(doc)
	if err != nil {
		return nil, fmt.Errorf("pawnapi: encoding schema document: %w", err)
	}
	return buf, nil
}
