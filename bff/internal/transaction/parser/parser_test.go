package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfmodel "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// validPDF generates a small, valid in-memory PDF (one page with a line of text).
func validPDF(t *testing.T) []byte {
	t.Helper()
	js := `{"pages":{"1":{"content":{"text":[{"value":"statement","font":{"name":"Helvetica","size":12},"anchor":"center"}]}}}}`
	var buf bytes.Buffer
	if err := api.Create(nil, strings.NewReader(js), &buf, pdfmodel.NewDefaultConfiguration()); err != nil {
		t.Fatalf("create pdf: %v", err)
	}
	return buf.Bytes()
}

// encryptedPDF generates a valid PDF encrypted with the given user password.
func encryptedPDF(t *testing.T, password string) []byte {
	t.Helper()
	conf := pdfmodel.NewAESConfiguration(password, password, 256)
	conf.Cmd = pdfmodel.ENCRYPT
	var buf bytes.Buffer
	if err := api.Encrypt(bytes.NewReader(validPDF(t)), &buf, conf); err != nil {
		t.Fatalf("encrypt pdf: %v", err)
	}
	return buf.Bytes()
}

func TestPreparePDF(t *testing.T) {
	t.Run("valid unencrypted", func(t *testing.T) {
		if _, err := preparePDF(validPDF(t), ""); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("encrypted with correct password", func(t *testing.T) {
		if _, err := preparePDF(encryptedPDF(t, "secret"), "secret"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("encrypted with wrong password", func(t *testing.T) {
		_, err := preparePDF(encryptedPDF(t, "secret"), "nope")
		if err != domain.ErrPDFWrongPassword {
			t.Errorf("expected ErrPDFWrongPassword, got %v", err)
		}
	})

	t.Run("encrypted without password", func(t *testing.T) {
		_, err := preparePDF(encryptedPDF(t, "secret"), "")
		if err != domain.ErrPDFPasswordRequired {
			t.Errorf("expected ErrPDFPasswordRequired, got %v", err)
		}
	})

	t.Run("corrupt bytes", func(t *testing.T) {
		_, err := preparePDF([]byte("this is not a pdf"), "")
		if err == nil || !errors.Is(err, domain.ErrPDFCorrupt) {
			t.Errorf("expected ErrPDFCorrupt, got %v", err)
		}
	})
}

// fakeExtractor is a test rowExtractor that returns canned tool input.
type fakeExtractor struct {
	input recordTransactionsInput
	err   error
}

func (f fakeExtractor) extract(ctx context.Context, pdf []byte) (recordTransactionsInput, error) {
	return f.input, f.err
}

// inputFromJSON builds a recordTransactionsInput via its json tags, avoiding the
// anonymous row struct literal.
func inputFromJSON(t *testing.T, s string) recordTransactionsInput {
	t.Helper()
	var in recordTransactionsInput
	if err := json.Unmarshal([]byte(s), &in); err != nil {
		t.Fatalf("build input: %v", err)
	}
	return in
}

func TestExtract_RowMapping(t *testing.T) {
	pdf := validPDF(t)

	t.Run("valid rows map and carry balance", func(t *testing.T) {
		p := &parser{extractor: fakeExtractor{input: inputFromJSON(t, `{"transactions":[
			{"date":"2026-06-01","merchant":"Acme","category":"Food","amount":500,"direction":"OUTFLOW","balance":500,"notes":"DEBIT"},
			{"date":"2026-06-02","merchant":"Work","category":"Income","amount":2000,"direction":"INFLOW","balance":2500}
		]}`)}}

		lines, err := p.Extract(context.Background(), pdf, "", func(string) {})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d", len(lines))
		}
		if lines[0].Txn.Category != shared.CategoryFood || lines[0].Txn.Direction != shared.Outflow {
			t.Errorf("row 0 not mapped: %+v", lines[0].Txn)
		}
		if lines[0].Balance != 500 || lines[0].Txn.Merchant != "Acme" || lines[0].Txn.Notes != "DEBIT" {
			t.Errorf("row 0 fields wrong: %+v balance=%d", lines[0].Txn, lines[0].Balance)
		}
		if lines[1].Txn.Direction != shared.Inflow || lines[1].Balance != 2500 {
			t.Errorf("row 1 not mapped: %+v balance=%d", lines[1].Txn, lines[1].Balance)
		}
		if lines[0].Txn.ID != "" || lines[0].Txn.UserID != "" {
			t.Errorf("ID/UserID should be filled by the service, got %+v", lines[0].Txn)
		}
	})

	t.Run("invalid category", func(t *testing.T) {
		p := &parser{extractor: fakeExtractor{input: inputFromJSON(t, `{"transactions":[
			{"date":"2026-06-01","merchant":"Acme","category":"Bogus","amount":500,"direction":"OUTFLOW","balance":500}
		]}`)}}

		_, err := p.Extract(context.Background(), pdf, "", func(string) {})
		if !errors.Is(err, shared.ErrInvalidCategory) {
			t.Errorf("expected ErrInvalidCategory, got %v", err)
		}
	})

	t.Run("invalid direction", func(t *testing.T) {
		p := &parser{extractor: fakeExtractor{input: inputFromJSON(t, `{"transactions":[
			{"date":"2026-06-01","merchant":"Acme","category":"Food","amount":500,"direction":"SIDEWAYS","balance":500}
		]}`)}}

		_, err := p.Extract(context.Background(), pdf, "", func(string) {})
		if !errors.Is(err, shared.ErrInvalidDirection) {
			t.Errorf("expected ErrInvalidDirection, got %v", err)
		}
	})

	t.Run("empty input is non-nil", func(t *testing.T) {
		p := &parser{extractor: fakeExtractor{input: recordTransactionsInput{}}}

		lines, err := p.Extract(context.Background(), pdf, "", func(string) {})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if lines == nil {
			t.Fatal("expected non-nil slice")
		}
		if len(lines) != 0 {
			t.Errorf("expected empty slice, got %d", len(lines))
		}
	})
}
