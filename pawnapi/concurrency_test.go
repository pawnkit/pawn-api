package pawnapi

import (
	"sync"
	"testing"
)

// TestIndex_ConcurrentReads checks immutable index access under -race.
func TestIndex_ConcurrentReads(t *testing.T) {
	ix, err := NewIndex(sampleEntries())
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = ix.ByID("native:SetPlayerPos")
			_ = ix.ByName("SetPlayerHealth")
			_ = ix.ByKind(KindNative)
			_ = ix.ByProfile("openmp")
			_ = ix.Deprecated()
			_ = ix.Profiles()
			_ = ix.All()
		}()
	}
	wg.Wait()
}
