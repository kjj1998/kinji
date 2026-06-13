package claude

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/kjj1998/kinji/bff/internal/app"
	"github.com/kjj1998/kinji/bff/internal/domain"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type parser struct {
	client anthropic.Client
	model  anthropic.Model
}

// compile-time check that parser satisfies the application port.
var _ app.StatementParser = (*parser)(nil)

func NewParser(model string) app.StatementParser {
	return &parser{
		client: anthropic.NewClient(),  // reads ANTHROPIC_API_KEY from env
		model:  anthropic.Model(model), // e.g. "claude-sonnet-4-6" from cfg
	}
}

// Extract decrypts/validates the statement PDF, runs the LLM extraction, and
// returns the raw statement rows (with running balances). Balance reconciliation
// is the domain's responsibility (see domain.Statement), so it is not done here.
func (p *parser) Extract(
	ctx context.Context,
	pdf []byte,
	password string,
	onProgress func(stage string),
) ([]domain.StatementLine, error) {
	onProgress("validating")
	pdf, err := preparePDF(pdf, password)
	if err != nil {
		return nil, err
	}

	onProgress("parsing")
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	input, err := p.extract(ctx, pdf)
	if err != nil {
		return nil, err
	}

	lines := make([]domain.StatementLine, 0, len(input.Transactions))
	for _, t := range input.Transactions {
		cat, err := domain.ParseCategory(t.Category)
		if err != nil {
			return nil, fmt.Errorf("row for %q: %w", t.Merchant, err)
		}
		dir, err := domain.ParseDirection(t.Direction)
		if err != nil {
			return nil, fmt.Errorf("row for %q: %w", t.Merchant, err)
		}
		lines = append(lines, domain.StatementLine{
			Txn: domain.Transaction{
				Date:      t.Date,
				Merchant:  t.Merchant,
				Category:  cat,
				Amount:    t.Amount,
				Direction: dir,
				Notes:     t.Notes,
				// ID / UserID filled by the application service
			},
			Balance: t.Balance,
		})
	}
	return lines, nil
}

// preparePDF validates an unencrypted PDF or decrypts a password-protected one,
// mapping pdfcpu failures to domain error sentinels.
func preparePDF(pdf []byte, password string) ([]byte, error) {
	if password == "" {
		if err := api.Validate(bytes.NewReader(pdf), model.NewDefaultConfiguration()); err != nil {
			if errors.Is(err, pdfcpu.ErrWrongPassword) {
				return nil, domain.ErrPDFPasswordRequired
			}
			return nil, fmt.Errorf("%w: %v", domain.ErrPDFCorrupt, err)
		}
		return pdf, nil
	}

	var out bytes.Buffer
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password
	if err := api.Decrypt(bytes.NewReader(pdf), &out, conf); err != nil {
		if errors.Is(err, pdfcpu.ErrWrongPassword) {
			return nil, domain.ErrPDFWrongPassword
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrPDFCorrupt, err)
	}
	return out.Bytes(), nil
}

// extract runs the Anthropic tool-use call and returns the decoded tool input.
func (p *parser) extract(ctx context.Context, pdf []byte) (recordTransactionsInput, error) {
	b64 := base64.StdEncoding.EncodeToString(pdf)
	documentBlock := anthropic.NewDocumentBlock(anthropic.Base64PDFSourceParam{Data: b64})
	textBlock := anthropic.NewTextBlock(extractionPrompt)
	toolParams := anthropic.ToolParam{
		Name:        "record_transactions",
		Description: anthropic.String("Record every transaction extracted from the bank statement, in order."),
		InputSchema: RecordTransactionsInputSchema,
	}
	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(documentBlock, textBlock),
	}
	outputConfig := anthropic.OutputConfigParam{
		Effort: anthropic.OutputConfigEffortHigh,
	}
	toolChoice := anthropic.ToolChoiceUnionParam{
		OfTool: &anthropic.ToolChoiceToolParam{Name: "record_transactions"},
	}
	tools := []anthropic.ToolUnionParam{{OfTool: &toolParams}}

	slog.Info("parsing pdf statement via LLM")
	stream := p.client.Messages.NewStreaming(
		ctx,
		anthropic.MessageNewParams{
			Model:        p.model,
			MaxTokens:    16384,
			Messages:     messages,
			OutputConfig: outputConfig,
			ToolChoice:   toolChoice,
			Tools:        tools,
		},
	)

	message := anthropic.Message{}
	for stream.Next() {
		message.Accumulate(stream.Current())
	}
	if err := stream.Err(); err != nil {
		return recordTransactionsInput{}, fmt.Errorf("stream: %w", err)
	}

	slog.Info("decoding ToolUseBlock JSON")
	var input recordTransactionsInput
	found := false
	for _, block := range message.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.ToolUseBlock:
			if variant.Name != "record_transactions" {
				continue
			}
			if err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &input); err != nil {
				return recordTransactionsInput{}, fmt.Errorf("decode tool input: %w", err)
			}
			found = true
		}
	}
	if !found {
		return recordTransactionsInput{}, fmt.Errorf("model did not call record_transactions")
	}
	return input, nil
}