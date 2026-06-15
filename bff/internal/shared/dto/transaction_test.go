package dto

import (
	"reflect"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/shared"
)

func sampleDomainTransaction() shared.Transaction {
	return shared.Transaction{
		ID:        "txn-1",
		UserID:    "user-1",
		Date:      "2026-06-13",
		Merchant:  "Cafe",
		Category:  shared.CategoryFood,
		Amount:    1250,
		Direction: shared.Outflow,
		Notes:     "lunch",
		Split:     2,
	}
}

func TestToTransaction(t *testing.T) {
	in := sampleDomainTransaction()

	out := ToTransaction(in)

	want := Transaction{
		ID:        "txn-1",
		UserID:    "user-1",
		Date:      "2026-06-13",
		Merchant:  "Cafe",
		Category:  shared.CategoryFood,
		Amount:    1250,
		Direction: shared.Outflow,
		Notes:     "lunch",
		Split:     2,
	}
	if out != want {
		t.Errorf("ToTransaction mismatch:\n got %+v\nwant %+v", out, want)
	}
}

func TestTransaction_DomainRoundTrip(t *testing.T) {
	in := sampleDomainTransaction()

	got := ToTransaction(in).Domain()

	if got != in {
		t.Errorf("round trip changed value:\n got %+v\nwant %+v", got, in)
	}
}

func TestToTransactions(t *testing.T) {
	in := []shared.Transaction{
		sampleDomainTransaction(),
		{ID: "txn-2", Direction: shared.Inflow, Amount: 500},
	}

	out := ToTransactions(in)

	want := []Transaction{
		ToTransaction(in[0]),
		ToTransaction(in[1]),
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("ToTransactions mismatch:\n got %+v\nwant %+v", out, want)
	}
}

func TestToTransactions_EmptyIsNonNil(t *testing.T) {
	out := ToTransactions(nil)
	if out == nil {
		t.Fatal("expected non-nil slice so it marshals as [] not null")
	}
	if len(out) != 0 {
		t.Errorf("expected empty slice, got len %d", len(out))
	}
}

func TestDomainTransactions(t *testing.T) {
	in := []Transaction{
		ToTransaction(sampleDomainTransaction()),
		{ID: "txn-2", Direction: shared.Inflow, Amount: 500},
	}

	out := DomainTransactions(in)

	want := []shared.Transaction{
		in[0].Domain(),
		in[1].Domain(),
	}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("DomainTransactions mismatch:\n got %+v\nwant %+v", out, want)
	}
}

func TestDomainTransactions_EmptyIsNonNil(t *testing.T) {
	out := DomainTransactions(nil)
	if out == nil {
		t.Fatal("expected non-nil slice")
	}
	if len(out) != 0 {
		t.Errorf("expected empty slice, got len %d", len(out))
	}
}
