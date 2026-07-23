package pawnapi

import "errors"

// Use errors.Is to test for these; they may be wrapped with extra context.
var (
	ErrMissingID        = errors.New("pawnapi: missing id")
	ErrInvalidID        = errors.New("pawnapi: invalid id")
	ErrIDKindMismatch   = errors.New("pawnapi: id prefix does not match kind")
	ErrMissingName      = errors.New("pawnapi: missing name")
	ErrInvalidKind      = errors.New("pawnapi: invalid kind")
	ErrMissingSignature = errors.New("pawnapi: missing signature")
	ErrMissingValue     = errors.New("pawnapi: missing value")
	ErrMissingParamName = errors.New("pawnapi: parameter missing name")

	ErrNoAvailability = errors.New("pawnapi: at least one availability entry is required")
	ErrInvalidProfile = errors.New("pawnapi: invalid availability profile")

	ErrInvalidVersion      = errors.New("pawnapi: invalid version")
	ErrInvalidVersionRange = errors.New("pawnapi: invalid version range")

	ErrMissingSource     = errors.New("pawnapi: missing source")
	ErrMissingRepository = errors.New("pawnapi: source missing repository")
	ErrMissingPath       = errors.New("pawnapi: source missing path")
	ErrMissingCommit     = errors.New("pawnapi: source missing commit")
	ErrInvalidCommit     = errors.New("pawnapi: source commit is not a hex object id")
	ErrMissingLicense    = errors.New("pawnapi: source missing license")

	ErrMissingConfidence = errors.New("pawnapi: missing confidence")
	ErrInvalidConfidence = errors.New("pawnapi: invalid confidence")
	ErrMissingReview     = errors.New("pawnapi: missing review status")
	ErrInvalidReview     = errors.New("pawnapi: invalid review status")

	ErrDuplicateID            = errors.New("pawnapi: duplicate entry id")
	ErrDuplicateName          = errors.New("pawnapi: duplicate (kind, name) pair")
	ErrBrokenReplacement      = errors.New("pawnapi: deprecated.replacement references an id that does not exist in the dataset")
	ErrSelfReplacement        = errors.New("pawnapi: deprecated.replacement references its own id")
	ErrMissingDeprecatedSince = errors.New("pawnapi: deprecated entry missing since")

	ErrEntryNotFound = errors.New("pawnapi: entry not found")
)
