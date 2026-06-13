package handler

import "github.com/kjj1998/kinji/bff/internal/model"

// Period is the wire representation of a model.Period.
type Period struct {
	Year   int   `json:"year"`
	Months []int `json:"months"`
}

// ToPeriod maps a domain period to its wire representation.
func ToPeriod(p model.Period) Period {
	return Period{Year: p.Year, Months: p.Months}
}

// ToPeriods maps a slice of domain periods to wire representations. The result is
// never nil so it marshals as [] rather than null.
func ToPeriods(periods []model.Period) []Period {
	out := make([]Period, len(periods))
	for i, p := range periods {
		out[i] = ToPeriod(p)
	}
	return out
}
