package model

import "fmt"

// Money is a monetary amount stored as an integer number of cents, avoiding the
// rounding errors of floating-point arithmetic. A negative value is allowed and
// represents a debit relative to some baseline.
type Money int64

// FromCents builds a Money from a raw cent value (e.g. 36570 -> $365.70).
func FromCents(cents int64) Money { return Money(cents) }

// Cents returns the amount as an integer number of cents.
func (m Money) Cents() int64 { return int64(m) }

// Dollars returns the amount as a floating-point dollar value, for display only.
func (m Money) Dollars() float64 { return float64(m) / 100 }

// Add returns the sum of two amounts.
func (m Money) Add(other Money) Money { return m + other }

// Sub returns the difference of two amounts.
func (m Money) Sub(other Money) Money { return m - other }

// Negate returns the amount with its sign flipped.
func (m Money) Negate() Money { return -m }

// Abs returns the absolute value of the amount.
func (m Money) Abs() Money {
	if m < 0 {
		return -m
	}
	return m
}

// IsZero reports whether the amount is exactly zero.
func (m Money) IsZero() bool { return m == 0 }

// String formats the amount as a dollar value with two decimal places.
func (m Money) String() string { return fmt.Sprintf("$%.2f", m.Dollars()) }
