package claude

import "github.com/anthropics/anthropic-sdk-go"

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