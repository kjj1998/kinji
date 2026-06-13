package domain

// Period is a year together with the months within it that have transaction
// data, used to drive the date pickers in the UI.
type Period struct {
	Year   int
	Months []int
}
