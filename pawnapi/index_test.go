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

func TestLoad_ActorAPI(t *testing.T) {
	index, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	create, ok := index.ByID("native:CreateActor")
	if !ok || len(create.Availability) != 2 || create.Signature == nil || len(create.Signature.Parameters) != 5 {
		t.Fatalf("CreateActor = %+v", create)
	}
	animation, ok := index.ByID("native:GetActorAnimation")
	if !ok || len(animation.Availability) != 1 || animation.Availability[0].Profile != ProfileOpenMP {
		t.Fatalf("GetActorAnimation = %+v", animation)
	}
}

func TestLoad_CheckpointAPI(t *testing.T) {
	index, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	race, ok := index.ByID("native:SetPlayerRaceCheckpoint")
	if !ok || len(race.Availability) != 2 || race.Signature == nil || len(race.Signature.Parameters) != 9 {
		t.Fatalf("SetPlayerRaceCheckpoint = %+v", race)
	}
	active, ok := index.ByID("native:IsPlayerCheckpointActive")
	if !ok || len(active.Availability) != 1 || active.Availability[0].Profile != ProfileOpenMP {
		t.Fatalf("IsPlayerCheckpointActive = %+v", active)
	}
}

func TestLoad_DialogAPI(t *testing.T) {
	index, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	show, ok := index.ByID("native:ShowPlayerDialog")
	if !ok || len(show.Availability) != 2 || show.Signature == nil || !show.Signature.Parameters[len(show.Signature.Parameters)-1].Variadic {
		t.Fatalf("ShowPlayerDialog = %+v", show)
	}
	deprecated, ok := index.ByID("native:GetPlayerDialog")
	if !ok || deprecated.Deprecated == nil || deprecated.Deprecated.Replacement != "native:GetPlayerDialogID" {
		t.Fatalf("GetPlayerDialog = %+v", deprecated)
	}
}

func TestLoad_MenuAPI(t *testing.T) {
	index, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	create, ok := index.ByID("native:CreateMenu")
	if !ok || len(create.Availability) != 2 || create.Signature == nil || !create.Signature.Parameters[len(create.Signature.Parameters)-1].Variadic {
		t.Fatalf("CreateMenu = %+v", create)
	}
	invalid, ok := index.ByID("constant:INVALID_MENU")
	if !ok || invalid.Value == nil || invalid.Value.String() != "-1" || len(invalid.Constraints) == 0 {
		t.Fatalf("INVALID_MENU = %+v", invalid)
	}
}

func TestLoad_ObjectCoreAPI(t *testing.T) {
	index, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	attach, ok := index.ByID("native:AttachObjectToObject")
	if !ok || len(attach.Availability) != 2 || attach.Signature == nil || attach.Signature.Parameters[len(attach.Signature.Parameters)-1].Default == nil {
		t.Fatalf("AttachObjectToObject = %+v", attach)
	}
	objectType, ok := index.ByID("native:GetObjectType")
	if !ok || len(objectType.Availability) != 1 || objectType.Availability[0].Profile != ProfileOpenMP {
		t.Fatalf("GetObjectType = %+v", objectType)
	}
	deprecated, ok := index.ByID("native:SetObjectsDefaultCameraCol")
	if !ok || deprecated.Deprecated == nil || deprecated.Deprecated.Replacement != "native:SetObjectsDefaultCameraCollision" {
		t.Fatalf("SetObjectsDefaultCameraCol = %+v", deprecated)
	}
}

func TestLoad_ObjectMaterialEditingAPI(t *testing.T) {
	index, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	materialText, ok := index.ByID("native:SetObjectMaterialText")
	if !ok || materialText.Signature == nil || !materialText.Signature.Parameters[len(materialText.Signature.Parameters)-1].Variadic {
		t.Fatalf("SetObjectMaterialText = %+v", materialText)
	}
	if materialText.Signature.Parameters[2].Default == nil || materialText.Signature.Parameters[2].Default.String() != "0" {
		t.Fatalf("SetObjectMaterialText material index = %+v", materialText.Signature.Parameters[2])
	}
	begin, ok := index.ByID("native:BeginObjectEditing")
	if !ok || len(begin.Availability) != 1 || begin.Availability[0].Profile != ProfileOpenMP {
		t.Fatalf("BeginObjectEditing = %+v", begin)
	}
	edit, ok := index.ByID("native:EditObject")
	if !ok || len(edit.Availability) != 2 {
		t.Fatalf("EditObject = %+v", edit)
	}
	response, ok := index.ByID("constant:EDIT_RESPONSE_UPDATE")
	if !ok || response.Value == nil || response.Value.String() != "2" {
		t.Fatalf("EDIT_RESPONSE_UPDATE = %+v", response)
	}
}

func TestLoad_PlayerObjectAPI(t *testing.T) {
	index, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	materialText, ok := index.ByID("native:SetPlayerObjectMaterialText")
	if !ok || materialText.Signature == nil || !materialText.Signature.Parameters[len(materialText.Signature.Parameters)-1].Variadic {
		t.Fatalf("SetPlayerObjectMaterialText = %+v", materialText)
	}
	canonical, ok := index.ByID("native:SetPlayerObjectNoCameraCollision")
	if !ok || len(canonical.Availability) != 1 || canonical.Availability[0].Profile != ProfileOpenMP {
		t.Fatalf("SetPlayerObjectNoCameraCollision = %+v", canonical)
	}
	legacy, ok := index.ByID("native:SetPlayerObjectNoCameraCol")
	if !ok || legacy.Deprecated == nil || legacy.Deprecated.Replacement != canonical.ID {
		t.Fatalf("SetPlayerObjectNoCameraCol = %+v", legacy)
	}
	moved, ok := index.ByID("callback:OnPlayerObjectMoved")
	if !ok || len(moved.Availability) != 2 {
		t.Fatalf("OnPlayerObjectMoved = %+v", moved)
	}
}

func TestLoad_ObjectQueriesAttachmentsDLAPI(t *testing.T) {
	index, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	attached, ok := index.ByID("native:SetPlayerAttachedObject")
	if !ok || len(attached.Availability) != 2 || attached.Signature == nil || attached.Signature.Parameters[12].Default == nil || attached.Signature.Parameters[12].Default.String() != "1" {
		t.Fatalf("SetPlayerAttachedObject = %+v", attached)
	}
	dl, ok := index.ByID("native:AddCharModel")
	if !ok || len(dl.Availability) != 1 || dl.Availability[0].Profile != ProfileOpenMP || len(dl.Constraints) == 0 {
		t.Fatalf("AddCharModel = %+v", dl)
	}
	target, ok := index.ByID("native:GetObjectTarget")
	if !ok || target.Deprecated == nil || target.Deprecated.Replacement != "native:GetObjectMovingTargetPos" {
		t.Fatalf("GetObjectTarget = %+v", target)
	}
	callback, ok := index.ByID("callback:OnPlayerEditAttachedObject")
	if !ok || len(callback.Availability) != 2 {
		t.Fatalf("OnPlayerEditAttachedObject = %+v", callback)
	}
}

func TestLoad_ClassAPI(t *testing.T) {
	index, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	class, ok := index.ByID("native:AddPlayerClass")
	if !ok || len(class.Availability) != 2 || class.Signature == nil || class.Signature.Parameters[5].Default == nil || class.Signature.Parameters[5].Default.String() != "WEAPON_FIST" {
		t.Fatalf("AddPlayerClass = %+v", class)
	}
	get, ok := index.ByID("native:GetPlayerClass")
	if !ok || len(get.Availability) != 1 || get.Availability[0].Profile != ProfileOpenMP {
		t.Fatalf("GetPlayerClass = %+v", get)
	}
	teamCount, ok := index.ByID("function:SetTeamCount")
	if !ok || teamCount.Kind != KindFunction {
		t.Fatalf("SetTeamCount = %+v", teamCount)
	}
}
