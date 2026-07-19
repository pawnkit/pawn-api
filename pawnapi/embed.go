package pawnapi

import (
	"bytes"
	_ "embed"
)

// embeddedFull is generated from data/source. Do not edit it by hand.
//
//go:embed data/pawn-api.full.json
var embeddedFull []byte

// Load returns an index over the embedded dataset.
func Load() (*Index, error) {
	entries, err := LoadEntries(embeddedReader())
	if err != nil {
		return nil, err
	}
	return NewIndex(entries)
}

func embeddedReader() *bytes.Reader {
	return bytes.NewReader(embeddedFull)
}
