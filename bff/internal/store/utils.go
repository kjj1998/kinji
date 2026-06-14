package store

import (
	"fmt"

	"github.com/kjj1998/kinji/bff/internal/model"
)

// currentAndPreviousMonth returns the "2006-01" keys for the given month and the
// month immediately before it.
func currentAndPreviousMonth(month, year string) (string, string, error) {
	m, err := model.ParseMonth(month, year)
	if err != nil {
		return "", "", fmt.Errorf("computing current and previous month: %w", err)
	}
	return m.Key(), m.Previous().Key(), nil
}
