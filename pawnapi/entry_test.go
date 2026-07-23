package pawnapi

import (
	"errors"
	"testing"
)

func validEntry() Entry {
	return Entry{
		ID:   "native:SetPlayerPos",
		Kind: KindNative,
		Name: "SetPlayerPos",
		Signature: &Signature{
			Parameters: []Parameter{
				{Name: "playerid"},
				{Name: "x", Tag: "Float"},
			},
			ReturnTag: "bool",
		},
		Availability: []Availability{
			{Profile: "openmp", Since: "1.0.0"},
		},
		Source: Source{
			Repository: "openmultiplayer/omp-stdlib",
			Path:       "omp_player.inc",
			Commit:     "689c824f6cc558b1ca1b36cfd1aae7f1cda16d65",
			License:    "MPL-2.0",
		},
		Confidence:   ConfidenceHigh,
		ReviewStatus: ReviewReviewed,
	}
}

func TestEntryValidate_Valid(t *testing.T) {
	if err := validEntry().Validate(); err != nil {
		t.Fatalf("expected valid entry, got error: %v", err)
	}
}

func TestEntryValidate_MissingID(t *testing.T) {
	e := validEntry()
	e.ID = ""
	if err := e.Validate(); !errors.Is(err, ErrMissingID) {
		t.Fatalf("got %v, want ErrMissingID", err)
	}
}

func TestEntryValidate_InvalidIDPattern(t *testing.T) {
	e := validEntry()
	e.ID = "not-namespaced"
	if err := e.Validate(); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("got %v, want ErrInvalidID", err)
	}
}

func TestEntryValidate_IDKindMismatch(t *testing.T) {
	e := validEntry()
	e.ID = "callback:SetPlayerPos"
	if err := e.Validate(); !errors.Is(err, ErrIDKindMismatch) {
		t.Fatalf("got %v, want ErrIDKindMismatch", err)
	}
}

func TestEntryValidate_InvalidKind(t *testing.T) {
	// id keeps a valid namespace prefix so the id-pattern check passes and
	// the (later) kind-validity check is what's actually exercised.
	e := validEntry()
	e.Kind = "bogus"
	if err := e.Validate(); !errors.Is(err, ErrInvalidKind) {
		t.Fatalf("got %v, want ErrInvalidKind", err)
	}
}

func TestEntryValidate_MissingName(t *testing.T) {
	e := validEntry()
	e.Name = ""
	if err := e.Validate(); !errors.Is(err, ErrMissingName) {
		t.Fatalf("got %v, want ErrMissingName", err)
	}
}

func TestEntryValidate_NativeMissingSignature(t *testing.T) {
	e := validEntry()
	e.Signature = nil
	if err := e.Validate(); !errors.Is(err, ErrMissingSignature) {
		t.Fatalf("got %v, want ErrMissingSignature", err)
	}
}

func TestEntryValidate_ConstantMissingValue(t *testing.T) {
	e := Entry{
		ID:           "constant:MAX_PLAYERS",
		Kind:         KindConstant,
		Name:         "MAX_PLAYERS",
		Availability: []Availability{{Profile: "openmp", Since: "1.0.0"}},
		Source:       validEntry().Source,
		Confidence:   ConfidenceHigh,
		ReviewStatus: ReviewReviewed,
	}
	if err := e.Validate(); !errors.Is(err, ErrMissingValue) {
		t.Fatalf("got %v, want ErrMissingValue", err)
	}
}

func TestEntryValidate_ConstantWithValue(t *testing.T) {
	v := NumberLiteral(1000)
	e := Entry{
		ID:           "constant:MAX_PLAYERS",
		Kind:         KindConstant,
		Name:         "MAX_PLAYERS",
		Value:        &v,
		Availability: []Availability{{Profile: "openmp", Since: "1.0.0"}},
		Source:       validEntry().Source,
		Confidence:   ConfidenceHigh,
		ReviewStatus: ReviewReviewed,
	}
	if err := e.Validate(); err != nil {
		t.Fatalf("expected valid constant, got: %v", err)
	}
}

func TestEntryValidate_NoAvailability(t *testing.T) {
	e := validEntry()
	e.Availability = nil
	if err := e.Validate(); !errors.Is(err, ErrNoAvailability) {
		t.Fatalf("got %v, want ErrNoAvailability", err)
	}
}

func TestEntryValidate_InvalidProfile(t *testing.T) {
	e := validEntry()
	e.Availability = []Availability{{Profile: "OpenMP", Since: "1.0.0"}}
	if err := e.Validate(); !errors.Is(err, ErrInvalidProfile) {
		t.Fatalf("got %v, want ErrInvalidProfile", err)
	}
}

func TestEntryValidate_InvalidVersionSince(t *testing.T) {
	e := validEntry()
	e.Availability = []Availability{{Profile: "openmp", Since: "not-a-version"}}
	if err := e.Validate(); !errors.Is(err, ErrInvalidVersion) {
		t.Fatalf("got %v, want ErrInvalidVersion", err)
	}
}

func TestEntryValidate_SinceAfterUntil(t *testing.T) {
	until := "1.0.0"
	e := validEntry()
	e.Availability = []Availability{{Profile: "openmp", Since: "2.0.0", Until: &until}}
	if err := e.Validate(); !errors.Is(err, ErrInvalidVersionRange) {
		t.Fatalf("got %v, want ErrInvalidVersionRange", err)
	}
}

func TestEntryValidate_DeprecatedMissingSince(t *testing.T) {
	e := validEntry()
	e.Deprecated = &Deprecation{Replacement: "native:Other"}
	if err := e.Validate(); !errors.Is(err, ErrMissingDeprecatedSince) {
		t.Fatalf("got %v, want ErrMissingDeprecatedSince", err)
	}
}

func TestEntryValidate_SelfReplacement(t *testing.T) {
	e := validEntry()
	e.Deprecated = &Deprecation{Since: "1.1.0", Replacement: e.ID}
	if err := e.Validate(); !errors.Is(err, ErrSelfReplacement) {
		t.Fatalf("got %v, want ErrSelfReplacement", err)
	}
}

func TestEntryValidate_MissingSourceFields(t *testing.T) {
	e := validEntry()
	e.Source = Source{}
	if err := e.Validate(); !errors.Is(err, ErrMissingRepository) {
		t.Fatalf("got %v, want ErrMissingRepository", err)
	}
}

func TestEntryValidate_InvalidCommit(t *testing.T) {
	e := validEntry()
	e.Source.Commit = "not-hex!"
	if err := e.Validate(); !errors.Is(err, ErrInvalidCommit) {
		t.Fatalf("got %v, want ErrInvalidCommit", err)
	}
}

func TestEntryValidate_MissingConfidence(t *testing.T) {
	e := validEntry()
	e.Confidence = ""
	if err := e.Validate(); !errors.Is(err, ErrMissingConfidence) {
		t.Fatalf("got %v, want ErrMissingConfidence", err)
	}
}

func TestEntryValidate_InvalidConfidence(t *testing.T) {
	e := validEntry()
	e.Confidence = "very-sure"
	if err := e.Validate(); !errors.Is(err, ErrInvalidConfidence) {
		t.Fatalf("got %v, want ErrInvalidConfidence", err)
	}
}

func TestEntryValidate_MissingReviewStatus(t *testing.T) {
	e := validEntry()
	e.ReviewStatus = ""
	if err := e.Validate(); !errors.Is(err, ErrMissingReview) {
		t.Fatalf("got %v, want ErrMissingReview", err)
	}
}

func TestEntryValidate_InvalidReviewStatus(t *testing.T) {
	e := validEntry()
	e.ReviewStatus = "maybe"
	if err := e.Validate(); !errors.Is(err, ErrInvalidReview) {
		t.Fatalf("got %v, want ErrInvalidReview", err)
	}
}

func TestEntryValidate_MissingParamName(t *testing.T) {
	e := validEntry()
	e.Signature.Parameters = append(e.Signature.Parameters, Parameter{Tag: "Float"})
	if err := e.Validate(); !errors.Is(err, ErrMissingParamName) {
		t.Fatalf("got %v, want ErrMissingParamName", err)
	}
}
