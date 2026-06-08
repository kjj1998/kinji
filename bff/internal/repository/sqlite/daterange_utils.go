package sqlite

import (
	"fmt"
	"time"
)

const monthLayout = "2006-01"

func GetMonthRangeDateStrings(month, year string) (string, string) {
	first, _ := time.Parse(monthLayout, year+"-"+month)
	from := first.Format("2006-01-02")
	to := first.AddDate(0, 1, -1).Format("2006-01-02")

	return from, to
}

func currentAndPreviousMonth(month, year string) (string, string, error) {
	curMonth := year + "-" + month
	t, err := time.Parse(monthLayout, curMonth)
	if err != nil {
		return "", "", fmt.Errorf("parsing %s as month: %w", curMonth, err)
	}

	return curMonth, t.AddDate(0, -1, 0).Format(monthLayout), nil
}
