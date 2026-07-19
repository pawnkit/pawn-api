package pawnapi

import "testing"

func BenchmarkLoad(b *testing.B) {
	for range b.N {
		if _, err := Load(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewIndex(b *testing.B) {
	entries, err := LoadEntries(embeddedReader())
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for range b.N {
		if _, err := NewIndex(entries); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateDataset(b *testing.B) {
	entries, err := LoadEntries(embeddedReader())
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for range b.N {
		if err := ValidateDataset(entries); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIndex_ByID(b *testing.B) {
	ix, err := Load()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for range b.N {
		ix.ByID("native:SetPlayerPos")
	}
}

func BenchmarkMarshalEntries(b *testing.B) {
	entries, err := LoadEntries(embeddedReader())
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for range b.N {
		if _, err := MarshalEntries(entries); err != nil {
			b.Fatal(err)
		}
	}
}
