package model

import "errors"

var (
	// ErrInvalidCategory is returned when a raw value is not a known Category.
	ErrInvalidCategory = errors.New("invalid category")

	// ErrInvalidDirection is returned when a raw value is neither INFLOW nor OUTFLOW.
	ErrInvalidDirection = errors.New("invalid direction")

	// ErrBalanceMismatch is returned when a statement's running balances do not
	// reconcile against its transaction amounts.
	ErrBalanceMismatch = errors.New("statement balance mismatch")

	// ErrPDFPasswordRequired is returned when a statement PDF is encrypted and no
	// password was supplied.
	ErrPDFPasswordRequired = errors.New("pdf password required")

	// ErrPDFWrongPassword is returned when the supplied PDF password is incorrect.
	ErrPDFWrongPassword = errors.New("wrong pdf password")

	// ErrPDFCorrupt is returned when a statement PDF cannot be read.
	ErrPDFCorrupt = errors.New("invalid or corrupt pdf")
)
