package service

import (
	"math"
	"slices"
)

type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

func percentageChange[T Number](current, previous T) float64 {
	return roundTo2Dp(safeDivide(current-previous, previous) * 100)
}

func roundTo2Dp(value float64) float64 {
	return math.Round(value*100) / 100
}

func safeDivide[T Number](a, b T) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

func sortByAmountDesc[T any](s []T, amount func(T) int) {
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
