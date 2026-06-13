package app

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/domain"
)

// MockParser is a function-backed test double for StatementParser.
type MockParser struct {
	ExtractFn func(ctx context.Context, pdf []byte, password string, onProgress func(stage string)) ([]domain.StatementLine, error)
}

// compile-time check that MockParser satisfies the interface.
var _ StatementParser = (*MockParser)(nil)

func (m *MockParser) Extract(ctx context.Context, pdf []byte, password string, onProgress func(stage string)) ([]domain.StatementLine, error) {
	return m.ExtractFn(ctx, pdf, password, onProgress)
}