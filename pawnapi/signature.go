package pawnapi

// Parameter describes a callable parameter.
type Parameter struct {
	Name    string   `json:"name"`
	Tag     string   `json:"tag,omitempty"`
	Default *Literal `json:"default,omitempty"`

	Const           bool  `json:"const,omitempty"`
	Reference       bool  `json:"reference,omitempty"`
	ArrayDimensions []int `json:"arrayDimensions,omitempty"`
	Variadic        bool  `json:"variadic,omitempty"`
}

func (p Parameter) Validate() error {
	if p.Name == "" {
		return ErrMissingParamName
	}
	return nil
}

// Signature describes a callable entry's parameter list and return
// semantics.
type Signature struct {
	Parameters []Parameter `json:"parameters,omitempty"`
	ReturnTag  string      `json:"returnTag,omitempty"`

	// ReturnSemantics explains the meaning of a return value.
	ReturnSemantics string `json:"returnSemantics,omitempty"`
}

func (s Signature) Validate() error {
	for i, p := range s.Parameters {
		if err := p.Validate(); err != nil {
			return wrapIndex("parameters", i, err)
		}
	}
	return nil
}
