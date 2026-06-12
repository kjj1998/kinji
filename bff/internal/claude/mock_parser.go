package claude

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/models"
)

// MockParser implements claude.Parser
type MockParser struct {
	ParseStatementFn func(ctx context.Context, pdf []byte, onProgress func(stage string)) ([]models.Transaction, error)
}

var _ Parser = (*MockParser)(nil)

func (m *MockParser) ParseStatement(ctx context.Context, pdf []byte, onProgress func(stage string)) ([]models.Transaction, error) {
	return m.ParseStatementFn(ctx, pdf, onProgress)
}
