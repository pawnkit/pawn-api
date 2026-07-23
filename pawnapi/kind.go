package pawnapi

// Kind classifies an [Entry]. Values match pawn-api.schema.json.
type Kind string

const (
	KindNative   Kind = "native"
	KindCallback Kind = "callback"
	KindFunction Kind = "function"
	KindConstant Kind = "constant"
	KindTag      Kind = "tag"
	KindDefine   Kind = "define"
)

func (k Kind) IsValid() bool {
	switch k {
	case KindNative, KindCallback, KindFunction, KindConstant, KindTag, KindDefine:
		return true
	default:
		return false
	}
}

func (k Kind) String() string {
	return string(k)
}

// Confidence records how thoroughly an entry was verified.
type Confidence string

const (
	// ConfidenceHigh was checked against a primary source.
	ConfidenceHigh Confidence = "high"

	// ConfidenceMedium was sourced but not independently checked.
	ConfidenceMedium Confidence = "medium"

	// ConfidenceLow is provisional and lacks primary-source verification.
	ConfidenceLow Confidence = "low"
)

func (c Confidence) IsValid() bool {
	switch c {
	case ConfidenceHigh, ConfidenceMedium, ConfidenceLow:
		return true
	default:
		return false
	}
}

func (c Confidence) String() string {
	return string(c)
}

// ReviewStatus records whether a person checked an entry.
type ReviewStatus string

const (
	ReviewGenerated ReviewStatus = "generated"
	ReviewReviewed  ReviewStatus = "reviewed"
)

func (s ReviewStatus) IsValid() bool {
	switch s {
	case ReviewGenerated, ReviewReviewed:
		return true
	default:
		return false
	}
}

func (s ReviewStatus) String() string {
	return string(s)
}
