package store

import (
	"fmt"

	"github.com/kjj1998/kinji/bff/internal/model"
)

const dateLayout = "2006-01-02"

// GetMonthRangeDateStrings returns the first and last day of the given month as
// "2006-01-02" strings, for use as SQL date-range bounds.
func GetMonthRangeDateStrings(month, year string) (string, string) {
	m, _ := model.ParseMonth(month, year)
	start, end := m.Range()
	return start.Format(dateLayout), end.Format(dateLayout)
}

// currentAndPreviousMonth returns the "2006-01" keys for the given month and the
// month immediately before it.
func currentAndPreviousMonth(month, year string) (string, string, error) {
	m, err := model.ParseMonth(month, year)
	if err != nil {
		return "", "", fmt.Errorf("computing current and previous month: %w", err)
	}
	return m.Key(), m.Previous().Key(), nil
}
