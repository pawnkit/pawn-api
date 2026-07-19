package pawnapi

import (
	"fmt"
	"regexp"
	"strings"
)

var idPattern = regexp.MustCompile(`^(native|callback|function|constant|tag|define):.+$`)

// Entry describes one API entity.
type Entry struct {
	ID   string `json:"id"`
	Kind Kind   `json:"kind"`
	Name string `json:"name"`

	// Aliases lists other names for this entry.
	Aliases []string `json:"aliases,omitempty"`

	Signature *Signature `json:"signature,omitempty"`
	Value     *Literal   `json:"value,omitempty"`
	Tags      []string   `json:"tags,omitempty"`

	Availability []Availability `json:"availability"`
	Deprecated   *Deprecation   `json:"deprecated,omitempty"`

	// OwningInclude is the include used by Pawn source.
	OwningInclude string `json:"owningInclude,omitempty"`

	// Constraints records ranges, sentinels, and side effects.
	Constraints []string `json:"constraints,omitempty"`

	// CallbackContext describes when and how a callback runs.
	CallbackContext string `json:"callbackContext,omitempty"`

	DocumentationSummary string `json:"documentationSummary,omitempty"`
	DocumentationURL     string `json:"documentationUrl,omitempty"`

	Source Source `json:"source"`

	// Confidence records how thoroughly this entry was checked.
	Confidence Confidence `json:"confidence"`
	Notes      string     `json:"notes,omitempty"`
}

func kindPrefix(k Kind) string {
	if k == KindConstant {
		return "constant"
	}
	return string(k)
}

// Validate checks one entry. Use [ValidateDataset] for cross-entry checks.
func (e Entry) Validate() error {
	if e.ID == "" {
		return ErrMissingID
	}
	if !idPattern.MatchString(e.ID) {
		return fmt.Errorf("%w: %q", ErrInvalidID, e.ID)
	}
	if !e.Kind.IsValid() {
		return fmt.Errorf("%w: %q", ErrInvalidKind, e.Kind)
	}
	wantPrefix := kindPrefix(e.Kind) + ":"
	if !strings.HasPrefix(e.ID, wantPrefix) {
		return fmt.Errorf("%w: id %q does not start with %q for kind %q", ErrIDKindMismatch, e.ID, wantPrefix, e.Kind)
	}
	if e.Name == "" {
		return ErrMissingName
	}

	switch e.Kind {
	case KindNative, KindCallback, KindFunction:
		if e.Signature == nil {
			return ErrMissingSignature
		}
	case KindConstant, KindDefine:
		if e.Value == nil {
			return ErrMissingValue
		}
	}

	if e.Signature != nil {
		if err := e.Signature.Validate(); err != nil {
			return fmt.Errorf("signature: %w", err)
		}
	}

	if len(e.Availability) == 0 {
		return ErrNoAvailability
	}
	for i, a := range e.Availability {
		if err := a.Validate(); err != nil {
			return wrapIndex("availability", i, err)
		}
	}

	if e.Deprecated != nil {
		if err := e.Deprecated.Validate(); err != nil {
			return fmt.Errorf("deprecated: %w", err)
		}
		if e.Deprecated.Replacement == e.ID {
			return ErrSelfReplacement
		}
	}

	if err := e.Source.Validate(); err != nil {
		return fmt.Errorf("source: %w", err)
	}

	if e.Confidence == "" {
		return ErrMissingConfidence
	}
	if !e.Confidence.IsValid() {
		return fmt.Errorf("%w: %q", ErrInvalidConfidence, e.Confidence)
	}

	return nil
}

func wrapIndex(field string, i int, err error) error {
	return fmt.Errorf("%s[%d]: %w", field, i, err)
}
