package domain

import (
	"testing"
)

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		name     string  // Descriptive name of the specific scenario
		a, b     int     // Input parameters for the function
		expected float64 // The output you expect to get
	}{
		{name: "normal division", a: 4, b: 2, expected: 2.0},
		{name: "divide by 0", a: 4, b: 0, expected: 0.0},
		{name: "negative operands", a: -10, b: 5, expected: -2.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := safeDivide(tc.a, tc.b)
			if actual != tc.expected {
				t.Errorf("safeDivide(%d, %d) failed: expected %f, got %f", tc.a, tc.b, tc.expected, actual)
			}
		})
	}
}

func TestPercentageChange(t *testing.T) {
	tests := []struct {
		name     string  // Descriptive name of the specific scenario
		a, b     int     // Input parameters for the function
		expected float64 // The output you expect to get
	}{
		{name: "increase", a: 45, b: 40, expected: 12.5},
		{name: "decrease", a: 35, b: 40, expected: -12.5},
		{name: "previous equals 0", a: 10, b: 0, expected: 0.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := percentageChange(tc.a, tc.b)
			if actual != tc.expected {
				t.Errorf("percentageChange(%d, %d) failed: expected %f, got %f", tc.a, tc.b, tc.expected, actual)
			}
		})
	}
}

func TestRoundTo2Dp(t *testing.T) {
	tests := []struct {
		name     string  // Descriptive name of the specific scenario
		input    float64 // Input parameters for the function
		expected float64 // The output you expect to get
	}{
		{name: "half up rounding", input: 45.555, expected: 45.56},
		{name: "values with one dp", input: 45.4, expected: 45.40},
		{name: "already rounded values", input: 45.40, expected: 45.40},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := roundTo2Dp(tc.input)
			if actual != tc.expected {
				t.Errorf("roundTo2Dp(%f) failed: expected %f, got %f", tc.input, tc.expected, actual)
			}
		})
	}
}

func TestSortByAmountDesc_Struct(t *testing.T) {
	type item struct {
		name   string
		amount int
	}
	got := []item{
		{"a", 10}, {"b", 50}, {"c", 30}, {"d", 50},
	}
	sortByAmountDesc(got, func(i item) int { return i.amount })

	// assert non-strictly descending by amount (don't assert order within the two 50s)
	for i := 1; i < len(got); i++ {
		if got[i-1].amount < got[i].amount {
			t.Errorf("not descending at %d: %v", i, got)
		}
	}
}