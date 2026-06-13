package handler

import (
	"reflect"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/model"
)

func TestToPeriod(t *testing.T) {
	in := model.Period{Year: 2026, Months: []int{1, 6, 12}}

	out := ToPeriod(in)

	if out.Year != 2026 {
		t.Errorf("Year not mapped, got %d", out.Year)
	}
	if !reflect.DeepEqual(out.Months, []int{1, 6, 12}) {
		t.Errorf("Months not mapped, got %v", out.Months)
	}
}

func TestToPeriods(t *testing.T) {
	in := []model.Period{
		{Year: 2025, Months: []int{11, 12}},
		{Year: 2026, Months: []int{1}},
	}

	out := ToPeriods(in)

	want := []Period{
		{Year: 2025, Months: []int{11, 12}},
		{Year: 2026, Months: []int{1}},
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("ToPeriods mismatch:\n got %+v\nwant %+v", out, want)
	}
}

func TestToPeriods_EmptyIsNonNil(t *testing.T) {
	out := ToPeriods(nil)
	if out == nil {
		t.Fatal("expected non-nil slice so it marshals as [] not null")
	}
	if len(out) != 0 {
		t.Errorf("expected empty slice, got len %d", len(out))
	}
}
