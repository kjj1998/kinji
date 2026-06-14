package model

import "errors"

var (
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
