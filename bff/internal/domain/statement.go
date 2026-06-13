package domain

import "fmt"

// StatementLine is a single extracted row from a bank statement: the transaction
// together with the running account balance printed on that row (in cents). The
// balance is used only to verify extraction integrity and is not persisted.
type StatementLine struct {
	Txn     Transaction
	Balance int
}

// Statement is a short-lived aggregate over the rows extracted from one bank
// statement. It enforces the running-balance invariant before the transactions
// are accepted into the ledger.
type Statement struct {
	lines []StatementLine
}

// NewStatement builds a Statement from extracted lines, in the order they appear
// on the statement.
func NewStatement(lines []StatementLine) Statement {
	return Statement{lines: lines}
}

// Validate checks that each row's balance equals the previous row's balance plus
// the signed transaction amount (inflows add, outflows subtract). It returns an
// error wrapping ErrBalanceMismatch on the first row that does not reconcile.
func (s Statement) Validate() error {
	for i := 1; i < len(s.lines); i++ {
		prev := s.lines[i-1].Balance
		cur := s.lines[i]

		delta := cur.Txn.Amount
		if cur.Txn.IsOutflow() {
			delta = -delta
		}
		if prev+delta != cur.Balance {
			return fmt.Errorf("row %d (%q): %d %+d != %d: %w",
				i, cur.Txn.Merchant, prev, delta, cur.Balance, ErrBalanceMismatch)
		}
	}
	return nil
}

// Transactions returns the statement's transactions in order, dropping the
// balance bookkeeping.
func (s Statement) Transactions() []Transaction {
	txns := make([]Transaction, len(s.lines))
	for i, line := range s.lines {
		txns[i] = line.Txn
	}
	return txns
}
