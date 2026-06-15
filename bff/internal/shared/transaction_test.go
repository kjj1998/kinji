package shared

import (
	"errors"
	"testing"
)

func TestParseCategory(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		for _, raw := range []string{
			"Entertainment", "Food", "Groceries", "Health", "Income",
			"Shopping", "Subscriptions", "Transport", "Utilities", "Credit",
		} {
			got, err := ParseCategory(raw)
			if err != nil {
				t.Errorf("ParseCategory(%q): unexpected error %v", raw, err)
			}
			if string(got) != raw {
				t.Errorf("ParseCategory(%q) = %q", raw, got)
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for _, raw := range []string{"", "food", "Bogus", "INCOME"} {
			_, err := ParseCategory(raw)
			if !errors.Is(err, ErrInvalidCategory) {
				t.Errorf("ParseCategory(%q): expected ErrInvalidCategory, got %v", raw, err)
			}
		}
	})
}

func TestParseDirection(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		for _, raw := range []string{"INFLOW", "OUTFLOW"} {
			got, err := ParseDirection(raw)
			if err != nil || string(got) != raw {
				t.Errorf("ParseDirection(%q) = (%q, %v)", raw, got, err)
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for _, raw := range []string{"", "inflow", "SIDEWAYS"} {
			_, err := ParseDirection(raw)
			if !errors.Is(err, ErrInvalidDirection) {
				t.Errorf("ParseDirection(%q): expected ErrInvalidDirection, got %v", raw, err)
			}
		}
	})
}

func TestCategoryIsValid(t *testing.T) {
	if !CategoryFood.IsValid() {
		t.Error("CategoryFood should be valid")
	}
	if Category("Nope").IsValid() {
		t.Error("unknown category should be invalid")
	}
}

func TestDirectionIsValid(t *testing.T) {
	if !Inflow.IsValid() || !Outflow.IsValid() {
		t.Error("Inflow/Outflow should be valid")
	}
	if Direction("X").IsValid() {
		t.Error("unknown direction should be invalid")
	}
}

func TestTransactionInflowOutflow(t *testing.T) {
	in := Transaction{Direction: Inflow}
	if !in.IsInflow() || in.IsOutflow() {
		t.Errorf("inflow txn: IsInflow=%v IsOutflow=%v", in.IsInflow(), in.IsOutflow())
	}
	out := Transaction{Direction: Outflow}
	if out.IsInflow() || !out.IsOutflow() {
		t.Errorf("outflow txn: IsInflow=%v IsOutflow=%v", out.IsInflow(), out.IsOutflow())
	}
}
