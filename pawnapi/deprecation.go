package pawnapi

// Deprecation records when and why an entry was deprecated.
type Deprecation struct {
	Since       string `json:"since"`
	Replacement string `json:"replacement,omitempty"`
	Reason      string `json:"reason,omitempty"`
}

func (d Deprecation) Validate() error {
	if d.Since == "" {
		return ErrMissingDeprecatedSince
	}
	if err := validateVersion(d.Since); err != nil {
		return err
	}
	return nil
}
