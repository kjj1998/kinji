package claude

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/kjj1998/kinji/bff/internal/models"
)

var RecordTransactionsInputSchema = anthropic.ToolInputSchemaParam{
	Properties: map[string]any{
		"transactions": map[string]any{
			"type":        "array",
			"description": "Every transaction on the statement, in order.",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"date": map[string]any{
						"type":        "string",
						"description": "Transaction date, ISO YYYY-MM-DD (prefer the purchase date when present).",
					},
					"merchant": map[string]any{
						"type":        "string",
						"description": "Cleaned merchant name; strip any xx-NNNN card mask.",
					},
					"category": map[string]any{
						"type": "string",
						"enum": []string{
							"Entertainment", "Food", "Groceries", "Health",
							"Income", "Shopping", "Subscriptions", "Transport", "Utilities", "Credit",
						},
					},
					"amount": map[string]any{
						"type":        "integer",
						"description": "Amount in cents (365.70 → 36570).",
					},
					"direction": map[string]any{
						"type": "string",
						"enum": []string{"INFLOW", "OUTFLOW"},
					},
					"balance": map[string]any{
						"type":        "integer",
						"description": "Running account balance after this transaction, in cents (for the balance check).",
					},
					"notes": map[string]any{
						"type":        "string",
						"description": "Optional, e.g. the transaction-type line like DEBIT PURCHASE.",
					},
				},
				"required":             []string{"date", "merchant", "category", "amount", "direction", "balance"},
				"additionalProperties": false,
			},
		},
	},
	Required: []string{"transactions"},
}

type recordTransactionsInput struct {
	Transactions []struct {
		Date      string `json:"date"`
		Merchant  string `json:"merchant"`
		Category  string `json:"category"`
		Amount    int    `json:"amount"`    // cents
		Direction string `json:"direction"` // INFLOW | OUTFLOW
		Balance   int    `json:"balance"`   // cents, for the guard
		Notes     string `json:"notes"`
	} `json:"transactions"`
}

const extractionPrompt = `You are extracting transactions from a bank or credit-card statement PDF for the Kinji expense tracker.

Call the record_transactions tool exactly once, with every transaction visible in the statement, in the order they appear.

For each transaction, provide:

- date: ISO YYYY-MM-DD. When the description contains an embedded purchase date (e.g. "28/06/24"), use that — it is the actual transaction date. Otherwise use the row's posting date. If the statement abbreviates the year, infer the full year from the statement period shown elsewhere on the page.
- merchant: the cleaned merchant name. Strip leading card-mask tokens like "xx-4070" or "xx-6790" and trailing single-letter country codes like "S" (Singapore). Preserve the company casing as it appears (e.g. "KOUFU PTE LTD", "APPLE.COM/SG"). If the company name has digits behind like "McDonalds 930144", truncate and take "McDonalds" only.
- category: exactly one of "Credit", "Entertainment", "Food", "Groceries", "Health", "Income", "Shopping", "Subscriptions", "Transport", "Utilities". Pick the closest match using these cues — bank interest, rebate, cashback, credit → Credit, incoming salary → Income; ride-hailing, transit, EZ-Link, fuel → Transport; restaurants, cafes, food courts → Food; supermarkets and grocery stores → Groceries; clinics, pharmacies, fitness → Health; streaming services, SaaS, app subscriptions → Subscriptions; utility bills, telcos → Utilities; cinemas, games, events → Entertainment; everything else retail → Shopping.
- amount: integer cents — always positive. "365.70" → 36570, "1,234.56" → 123456. The sign lives in direction, not in amount.
- direction: "OUTFLOW" if the value sits in a Withdrawal / Debit column. "INFLOW" if it sits in a Deposit / Credit column.
- balance: the running account balance printed on that row, in integer cents. This is verified by arithmetic, so it must match the statement exactly.
- notes: optional. If the row has a transaction-type label such as "DEBIT PURCHASE", "FUND TRANSFER", or "GIRO", put it here. Empty string if none.

Skip rows that are not actual transactions: opening balance ("BALANCE B/F"), closing balance, page subtotals, column headers, marketing or legal text, and footer notes.

Accuracy matters more than coverage. If a row is genuinely unreadable, omit it rather than guess. But extract every row you can read with confidence — the balances will be checked against running arithmetic, so amount, direction, and balance must be precise.`

type Parser interface {
	ParseStatement(ctx context.Context, pdf []byte) ([]models.Transaction, error)
}

type parser struct {
	client anthropic.Client
	model  anthropic.Model
}

func NewParser(model string) Parser {
	return &parser{
		client: anthropic.NewClient(),  // reads ANTHROPIC_API_KEY from env
		model:  anthropic.Model(model), // e.g. "claude-sonnet-4-6" from cfg
	}
}

func (p *parser) ParseStatement(ctx context.Context, pdf []byte) ([]models.Transaction, error) {
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
		return nil, fmt.Errorf("stream: %w", err)
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
				return nil, fmt.Errorf("decode tool input: %w", err)
			}
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("model did not call record_transactions")
	}

	slog.Info("running balance guard checks")
	for i := 1; i < len(input.Transactions); i++ {
		prev, cur := input.Transactions[i-1].Balance, input.Transactions[i]
		delta := cur.Amount
		if cur.Direction == "OUTFLOW" {
			delta = -delta
		}
		if prev+delta != cur.Balance {
			return nil, fmt.Errorf("balance mismatch at row %d (%q): %d %+d != %d",
				i, cur.Merchant, prev, delta, cur.Balance)
		}
	}

	slog.Info("creating transactions")
	txns := make([]models.Transaction, 0, len(input.Transactions))
	for _, t := range input.Transactions {
		cat := models.Category(t.Category)
		if !cat.IsValid() {
			return nil, fmt.Errorf("unknown category %q on row for %q", t.Category, t.Merchant)
		}

		dir := models.Direction(t.Direction)
		if !dir.IsValid() {
			return nil, fmt.Errorf("unknown direction %q on row for %q", t.Direction, t.Merchant)
		}

		txns = append(txns, models.Transaction{
			Date:      t.Date,
			Merchant:  t.Merchant,
			Category:  cat,
			Amount:    t.Amount,
			Direction: dir,
			Notes:     t.Notes,
			// ID / UserID filled by the service
		})
	}
	return txns, nil
}
