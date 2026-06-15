package domain

import "testing"

func TestNewValueAndChange_Int(t *testing.T) {
	cases := []struct {
		name       string
		values     []int
		wantValue  int
		wantChange int
	}{
		{"empty is zero", nil, 0, 0},
		{"single value has no change", []int{500}, 500, 0},
		{"two values: change is most-recent minus prior", []int{500, 300}, 500, 200},
		{"negative change when spending fell", []int{300, 500}, 300, -200},
		{"only the two most recent are used", []int{500, 300, 999}, 500, 200},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewValueAndChange(tc.values)
			if got.Value != tc.wantValue {
				t.Errorf("Value = %d, want %d", got.Value, tc.wantValue)
			}
			if got.Change != tc.wantChange {
				t.Errorf("Change = %d, want %d", got.Change, tc.wantChange)
			}
		})
	}
}

func TestNewValueAndChange_Float(t *testing.T) {
	got := NewValueAndChange([]float64{12.5, 10.0})
	if got.Value != 12.5 {
		t.Errorf("Value = %v, want 12.5", got.Value)
	}
	if got.Change != 2.5 {
		t.Errorf("Change = %v, want 2.5", got.Change)
	}
}
