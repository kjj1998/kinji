package shared

import "errors"

var (
	// ErrInvalidCategory is returned when a raw value is not a known Category.
	ErrInvalidCategory = errors.New("invalid category")

	// ErrInvalidDirection is returned when a raw value is neither INFLOW nor OUTFLOW.
	ErrInvalidDirection = errors.New("invalid direction")
)
