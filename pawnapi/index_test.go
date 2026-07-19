package pawnapi

import "testing"

func sampleEntries() []Entry {
	pos := validEntry()

	health := validEntry()
	health.ID = "native:SetPlayerHealth"
	health.Name = "SetPlayerHealth"
	health.Availability = []Availability{{Profile: "samp-037", Since: "0.3.7"}}

	cb := validEntry()
	cb.ID = "callback:OnPlayerConnect"
	cb.Kind = KindCallback
	cb.Name = "OnPlayerConnect"
	cb.Availability = []Availability{
		{Profile: "samp-037", Since: "0.3.7"},
		{Profile: "openmp", Since: "1.0.0"},
	}

	deprecated := validEntry()
	deprecated.ID = "native:OldSetPlayerPos"
	deprecated.Name = "OldSetPlayerPos"
	deprecated.Deprecated = &Deprecation{Since: "1.1.0", Replacement: "native:SetPlayerPos"}

	return []Entry{pos, health, cb, deprecated}
}

func TestNewIndex_RejectsInvalidDataset(t *testing.T) {
	entries := sampleEntries()
	entries = append(entries, entries[0]) // duplicate id
	if _, err := NewIndex(entries); err == nil {
		t.Fatal("expected NewIndex to reject a dataset with a duplicate id")
	}
}

func TestIndex_ByID(t *testing.T) {
	ix, err := NewIndex(sampleEntries())
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ix.ByID("native:SetPlayerHealth")
	if !ok {
		t.Fatal("expected to find native:SetPlayerHealth")
	}
	if e.Name != "SetPlayerHealth" {
		t.Fatalf("got name %q", e.Name)
	}
	if _, ok := ix.ByID("native:DoesNotExist"); ok {
		t.Fatal("expected ByID to report not found")
	}
}

func TestIndex_ByName(t *testing.T) {
	ix, err := NewIndex(sampleEntries())
	if err != nil {
		t.Fatal(err)
	}
	got := ix.ByName("SetPlayerPos")
	if len(got) != 1 || got[0].ID != "native:SetPlayerPos" {
		t.Fatalf("got %+v", got)
	}
}

func TestIndex_ByNameIncludesAliases(t *testing.T) {
	entries := sampleEntries()
	entries[0].Aliases = []string{"SetPlayerPosition"}
	ix, err := NewIndex(entries)
	if err != nil {
		t.Fatal(err)
	}
	got := ix.ByName("SetPlayerPosition")
	if len(got) != 1 || got[0].ID != "native:SetPlayerPos" {
		t.Fatalf("got %+v", got)
	}
}

func TestIndex_ByKindName(t *testing.T) {
	ix, err := NewIndex(sampleEntries())
	if err != nil {
		t.Fatal(err)
	}
	e, ok := ix.ByKindName(KindCallback, "OnPlayerConnect")
	if !ok || e.Kind != KindCallback {
		t.Fatalf("got %+v, ok=%v", e, ok)
	}
	if _, ok := ix.ByKindName(KindNative, "OnPlayerConnect"); ok {
		t.Fatal("expected no native named OnPlayerConnect")
	}
}

func TestIndex_ByKind(t *testing.T) {
	ix, err := NewIndex(sampleEntries())
	if err != nil {
		t.Fatal(err)
	}
	natives := ix.ByKind(KindNative)
	if len(natives) != 3 {
		t.Fatalf("got %d natives, want 3", len(natives))
	}
	for _, e := range natives {
		if e.Kind != KindNative {
			t.Fatalf("ByKind returned a non-native entry: %+v", e)
		}
	}
}

func TestIndex_ByProfile(t *testing.T) {
	ix, err := NewIndex(sampleEntries())
	if err != nil {
		t.Fatal(err)
	}
	samp := ix.ByProfile("samp-037")
	if len(samp) != 2 {
		t.Fatalf("got %d samp-037 entries, want 2", len(samp))
	}
	openmp := ix.ByProfile("openmp")
	if len(openmp) != 3 {
		t.Fatalf("got %d openmp entries, want 3", len(openmp))
	}
	if got := ix.ByProfile("nonexistent-profile"); got != nil {
		t.Fatalf("got %v, want nil", got)
	}
}

func TestIndex_Deprecated(t *testing.T) {
	ix, err := NewIndex(sampleEntries())
	if err != nil {
		t.Fatal(err)
	}
	dep := ix.Deprecated()
	if len(dep) != 1 || dep[0].ID != "native:OldSetPlayerPos" {
		t.Fatalf("got %+v", dep)
	}
}

func TestIndex_Profiles(t *testing.T) {
	ix, err := NewIndex(sampleEntries())
	if err != nil {
		t.Fatal(err)
	}
	profiles := ix.Profiles()
	want := []string{"openmp", "samp-037"}
	if len(profiles) != len(want) {
		t.Fatalf("got %v, want %v", profiles, want)
	}
	for i := range want {
		if profiles[i] != want[i] {
			t.Fatalf("got %v, want %v", profiles, want)
		}
	}
}

func TestIndex_AllSortedByKindAndID(t *testing.T) {
	ix, err := NewIndex(sampleEntries())
	if err != nil {
		t.Fatal(err)
	}
	all := ix.All()
	if len(all) != 4 {
		t.Fatalf("got %d entries, want 4", len(all))
	}
	for i := 1; i < len(all); i++ {
		prev, cur := all[i-1], all[i]
		if prev.Kind > cur.Kind || (prev.Kind == cur.Kind && prev.ID > cur.ID) {
			t.Fatalf("entries not sorted: %s before %s", prev.ID, cur.ID)
		}
	}
}

func TestLoad_EmbeddedDatasetIsValid(t *testing.T) {
	ix, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}
	if ix.Len() == 0 {
		t.Fatal("expected the embedded dataset to be non-empty")
	}
	if _, ok := ix.ByID("native:SetPlayerPos"); !ok {
		t.Error("expected embedded dataset to contain native:SetPlayerPos")
	}
	if _, ok := ix.ByID("callback:OnPlayerConnect"); !ok {
		t.Error("expected embedded dataset to contain callback:OnPlayerConnect")
	}
}
