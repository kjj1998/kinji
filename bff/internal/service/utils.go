package service

import (
	"math"
	"slices"
)

func percentageChange(current, previous float64) float64 {
	return roundTo2Dp(safeDivide(current-previous, previous) * 100)
}

func roundTo2Dp(value float64) float64 {
	return math.Round(value*100) / 100
}

func safeDivide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func sortByAmountDesc[T any](s []T, amount func(T) float64) {
	slices.SortFunc(s, func(a, b T) int {
		if amount(b) > amount(a) {
			return 1
		}
		if amount(b) < amount(a) {
			return -1
		}
		return 0
	})
}