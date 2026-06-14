package domain

import (
	"time"

	"github.com/kjj1998/kinji/bff/internal/shared"
)

// ValueAndChange pairs a current value with its change relative to a prior
// period (e.g. this month's total spent and the delta from last month).
type ValueAndChange[T int | float64] struct {
	Value  T
	Change T
}

// NewValueAndChange builds a ValueAndChange from period values ordered
// most-recent-first. With a single value the change is zero; with two or more
// the change is the difference between the two most recent.
func NewValueAndChange[T int | float64](values []T) ValueAndChange[T] {
	if len(values) == 0 {
		return ValueAndChange[T]{}
	}
	v := ValueAndChange[T]{Value: values[0]}
	if len(values) > 1 {
		v.Change = values[0] - values[1]
	}
	return v
}

// MerchantSpending is the total outflow to a single merchant over a period.
type MerchantSpending struct {
	Name     string
	Amount   int
	Category shared.Category
}

// CategorySpending is the total outflow within a single category over a period.
type CategorySpending struct {
	Category shared.Category
	Amount   int
}

// CategorySpendingChange describes how spending in a category moved between two
// periods. IsNew reports that the category had no spending in the prior period.
type CategorySpendingChange struct {
	Category         shared.Category
	Amount           int
	Change           int
	PercentageChange int
	IsNew            bool
}

// DaySpending is the total outflow for a weekday, keyed by the raw weekday so
// the presentation layer owns label formatting.
type DaySpending struct {
	Weekday time.Weekday
	Amount  int
}

// MonthSpending is the total outflow for a calendar month, keyed by the month's
// first instant so the presentation layer owns label formatting.
type MonthSpending struct {
	Month  time.Time
	Amount int
}
