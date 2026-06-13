package model

import (
	"fmt"
	"time"
)

const monthLayout = "2006-01"

// Month is a specific calendar month (a year and month), independent of any day
// or time within it.
type Month struct {
	start time.Time
}

// ParseMonth builds a Month from a two-digit month ("01"-"12") and four-digit
// year ("2006").
func ParseMonth(month, year string) (Month, error) {
	start, err := time.Parse(monthLayout, year+"-"+month)
	if err != nil {
		return Month{}, fmt.Errorf("parse month %s-%s: %w", year, month, err)
	}
	return Month{start: start}, nil
}

// Start returns the first instant of the month.
func (m Month) Start() time.Time { return m.start }

// End returns the last day of the month.
func (m Month) End() time.Time { return m.start.AddDate(0, 1, -1) }

// Range returns the first and last day of the month.
func (m Month) Range() (start, end time.Time) { return m.Start(), m.End() }

// AddMonths returns the month n months from this one; n may be negative.
func (m Month) AddMonths(n int) Month { return Month{start: m.start.AddDate(0, n, 0)} }

// Previous returns the month immediately before this one.
func (m Month) Previous() Month { return m.AddMonths(-1) }

// Key returns the canonical "YYYY-MM" identifier for the month.
func (m Month) Key() string { return m.start.Format(monthLayout) }
