# Feature-Based Folder Structure Refactor — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reorganize the `bff` Go service from a layered structure (`model`/`service`/`handler`/`store`/`parser`) into a feature-based structure (`transaction`/`summary`) with `shared` and `platform` kernels, with zero behavior changes.

**Architecture:** Two feature packages (`internal/transaction`, `internal/summary`), each with nested `domain`/`service`/`store`/`handler` (transaction also `parser`) sub-packages. Cross-cutting domain types live in `internal/shared`; infrastructure in `internal/platform`; composition root in `internal/server`. Features are siblings — they never import each other; both depend only on `shared` and `platform`.

**Tech Stack:** Go 1.x, `net/http`, `modernc.org/sqlite`, standard `testing`. Module path: `github.com/kjj1998/kinji/bff`.

---

## Working notes for the implementer

This is a **pure move/rename refactor**. File *bodies* move verbatim; only `package` clauses, import paths, and type qualifiers change. The whole module must `go build ./...` and `go test ./...` green after every task.

**Run all commands from the `bff/` module root** (`/Users/junjie/Projects/kinji/bff`).

### Identifier → destination package map

When a file moves, every `model.X` reference must be requalified to the package X now lives in:

| Identifier(s) | New package | Import path |
|---|---|---|
| `Transaction`, `Category`, `Direction`, `Inflow`, `Outflow`, `Category*` consts, `ParseCategory`, `ParseDirection` | `shared` | `internal/shared` |
| `Money`, `FromCents` | `shared` | `internal/shared` |
| `Month`, `ParseMonth` | `shared` | `internal/shared` |
| `ErrInvalidCategory`, `ErrInvalidDirection` | `shared` | `internal/shared` |
| `ClientError` (was `service.ClientError`) | `shared` | `internal/shared` |
| `GetMonthRangeDateStrings` (was `store`) | `shared` | `internal/shared` |
| `Statement`, `StatementLine`, `NewStatement`, `Period` | transaction `domain` | `internal/transaction/domain` |
| `ErrBalanceMismatch`, `ErrPDFPasswordRequired`, `ErrPDFWrongPassword`, `ErrPDFCorrupt` | transaction `domain` | `internal/transaction/domain` |
| `MonthlySummary`, `SummaryInput`, `SummaryCalculator`, `NewSummaryCalculator`, `Number`, `PercentageChange`, `RoundTo2Dp`, `SafeDivide`, `SortByAmountDesc` | summary `domain` | `internal/summary/domain` |
| `MerchantSpending`, `CategorySpending`, `CategorySpendingChange`, `DaySpending`, `MonthSpending`, `ValueAndChange`, `NewValueAndChange` | summary `domain` | `internal/summary/domain` |

Both feature `domain` packages are named `package domain`. Because features never import each other, no aliasing is needed inside a feature. Where a single file imports both `shared` and its own `domain`, that's fine.

### Repository method / SQL-const split

| Transaction feature | Summary feature |
|---|---|
| `GetMonthlyTransactions`, `SaveTransactions`, `GetTransactionPeriods` (+ helpers `getTransactionsWithinDateRange`, `getPeriods`) | `GetMonthlyTopMerchants`, `GetTotalIncomeTotalSpentAndNetSavings`, `GetCategorySpendingForLastTwoMonths`, `GetMonthlyTopCategories`, `GetLastSixMonthsExpenses` |
| consts: `getAllTransactionsWithinDateRange`, `getMonthAndYearWhichTransactionsOccur`, `saveTransactions` | consts: `getTopSpendingMerchantsWithinDateRange`, `getTotalIncomeTotalSpentAndNetSavingsForTwoMonths`, `getCategorySpendingForTwoMonths`, `getTopSpendingCategoriesWithinDateRange`, `getTotalMonthlyExpensesWithinDateRange` |
| helper `GetMonthRangeDateStrings` → `shared` | helper `currentAndPreviousMonth` → `summary/store` |
| `schema` const + `NewClient` → `platform/database` | — |

---

## Task 1: Create `internal/shared` package

**Files:**
- Create: `internal/shared/transaction.go` (from `internal/model/transaction.go`)
- Create: `internal/shared/money.go` (from `internal/model/money.go`)
- Create: `internal/shared/month.go` (from `internal/model/month.go` + `GetMonthRangeDateStrings` from `internal/store/utils.go`)
- Create: `internal/shared/errors.go` (`ErrInvalidCategory`, `ErrInvalidDirection` from `internal/model/errors.go`)
- Create: `internal/shared/client_error.go` (from `internal/service/errors.go`)

- [ ] **Step 1: Move the four model files into shared and rename the package**

```bash
git mv internal/model/transaction.go internal/shared/transaction.go
git mv internal/model/money.go        internal/shared/money.go
git mv internal/model/month.go        internal/shared/month.go
sed -i '' '1s/^package model$/package shared/' \
  internal/shared/transaction.go internal/shared/money.go internal/shared/month.go
```

- [ ] **Step 2: Create `internal/shared/errors.go` with only the category/direction errors**

```go
package shared

import "errors"

var (
	// ErrInvalidCategory is returned when a raw value is not a known Category.
	ErrInvalidCategory = errors.New("invalid category")

	// ErrInvalidDirection is returned when a raw value is neither INFLOW nor OUTFLOW.
	ErrInvalidDirection = errors.New("invalid direction")
)
```

Then remove those two `var` entries from `internal/model/errors.go` (leave the four statement/PDF errors in place for Task 4).

- [ ] **Step 3: Move `ClientError` into shared**

```bash
git mv internal/service/errors.go internal/shared/client_error.go
sed -i '' '1s/^package service$/package shared/' internal/shared/client_error.go
```

- [ ] **Step 4: Move `GetMonthRangeDateStrings` into `internal/shared/month.go`**

Cut this function from `internal/store/utils.go` and append it to `internal/shared/month.go` (it needs no `model.`/`shared.` qualifier — it returns plain strings):

```go
// GetMonthRangeDateStrings returns the inclusive first/last calendar day of the
// given month as YYYY-MM-DD strings.
func GetMonthRangeDateStrings(month, year string) (string, string) {
	// ... body moved verbatim from internal/store/utils.go ...
}
```

Leave `currentAndPreviousMonth` in `internal/store/utils.go` for now (moves in Task 11).

- [ ] **Step 5: Requalify internal references within the moved files**

Within `internal/shared/*.go`, the moved code referenced sibling `model` types by bare name (same package), so no `model.` prefixes exist to change. Confirm none remain:

```bash
grep -rn 'model\.' internal/shared/ || echo "clean"
```
Expected: `clean`

- [ ] **Step 6: Build the shared package**

Run: `go build ./internal/shared/`
Expected: builds with no output. (The rest of the module will NOT build yet — that's fixed in later tasks. Do not run `go build ./...` here.)

- [ ] **Step 7: Commit**

```bash
git add internal/shared/ internal/model/errors.go internal/store/utils.go internal/service/errors.go
git commit -m "refactor: extract shared kernel (Transaction, Money, Month, errors, ClientError)"
```

---

## Task 2: Create `internal/platform/database`

**Files:**
- Create: `internal/platform/database/database.go` (`NewClient` from `internal/store/sqlite.go`)
- Create: `internal/platform/database/schema.go` (`schema` const from `internal/store/queries.go`)

- [ ] **Step 1: Move `NewClient` and the schema const**

Create `internal/platform/database/database.go` containing `package database`, the `NewClient` function moved verbatim from `internal/store/sqlite.go` (it references `schema`), and its imports (`database/sql`, `fmt`, `_ "modernc.org/sqlite"`). Remove `NewClient` from `internal/store/sqlite.go`.

Create `internal/platform/database/schema.go`:

```go
package database

// schema is the SQLite DDL applied on connection.
const schema = `
` // ← move the full backtick body from internal/store/queries.go verbatim
```

Remove the `schema` const from `internal/store/queries.go`.

- [ ] **Step 2: Build the database package**

Run: `go build ./internal/platform/database/`
Expected: builds with no output.

- [ ] **Step 3: Commit**

```bash
git add internal/platform/database/ internal/store/
git commit -m "refactor: move sqlite client + schema to platform/database"
```

---

## Task 3: Create `internal/platform/http` and `internal/platform/config`

**Files:**
- Create: `internal/platform/http/middleware.go` (from `internal/handler/middleware.go`)
- Create: `internal/platform/http/http_utils.go` (from `internal/handler/http_utils.go`)
- Create: `internal/platform/http/health.go` (from `internal/handler/health.go`)
- Create: `internal/platform/config/config.go` (from `internal/config/config.go`)

- [ ] **Step 1: Move the HTTP infra and config files**

```bash
git mv internal/handler/middleware.go  internal/platform/http/middleware.go
git mv internal/handler/http_utils.go  internal/platform/http/http_utils.go
git mv internal/handler/health.go      internal/platform/http/health.go
git mv internal/config/config.go       internal/platform/config/config.go
sed -i '' '1s/^package handler$/package http/' \
  internal/platform/http/middleware.go internal/platform/http/http_utils.go internal/platform/http/health.go
sed -i '' '1s/^package config$/package config/' internal/platform/config/config.go
```

- [ ] **Step 2: Requalify `ClientError` references in `http_utils.go`**

If `internal/platform/http/http_utils.go` (or `middleware.go`) references `ClientError` (formerly same-package `service`/`handler` type), add the import and qualify it:

```bash
grep -n 'ClientError' internal/platform/http/*.go
```
For each hit, change `ClientError` → `shared.ClientError` and add `"github.com/kjj1998/kinji/bff/internal/shared"` to that file's import block. If there are zero hits, skip.

- [ ] **Step 3: Build the platform packages**

Run: `go build ./internal/platform/...`
Expected: builds with no output.

- [ ] **Step 4: Commit**

```bash
git add internal/platform/ internal/handler/ internal/config/
git commit -m "refactor: move http middleware/utils/health and config to platform"
```

---

## Task 4: Create `internal/transaction/domain`

**Files:**
- Create: `internal/transaction/domain/statement.go` (from `internal/model/statement.go`)
- Create: `internal/transaction/domain/period.go` (from `internal/model/period.go`)
- Create: `internal/transaction/domain/errors.go` (4 statement/PDF errors from `internal/model/errors.go`)
- Create: `internal/transaction/domain/statement_test.go` (from `internal/model/tests/statement_test.go`)

- [ ] **Step 1: Move statement + period and rename package**

```bash
git mv internal/model/statement.go internal/transaction/domain/statement.go
git mv internal/model/period.go     internal/transaction/domain/period.go
sed -i '' '1s/^package model$/package domain/' \
  internal/transaction/domain/statement.go internal/transaction/domain/period.go
```

- [ ] **Step 2: Requalify `Transaction` in statement.go**

`Statement.Transactions()` returns `[]Transaction`. `Transaction` now lives in `shared`. In `internal/transaction/domain/statement.go`, change `[]Transaction` → `[]shared.Transaction` (and any other `Transaction`/`Direction`/`Inflow`/`Outflow` references), and add the import:

```go
import "github.com/kjj1998/kinji/bff/internal/shared"
```
`ErrBalanceMismatch` stays bare (it moves into this same package in Step 3).

- [ ] **Step 3: Move the four statement/PDF errors into the domain package**

Create `internal/transaction/domain/errors.go`:

```go
package domain

import "errors"

var (
	// ErrBalanceMismatch is returned when a statement's running balances do not
	// reconcile against its transaction amounts.
	ErrBalanceMismatch = errors.New("statement balance mismatch")

	// ErrPDFPasswordRequired is returned when a statement PDF is encrypted and no
	// password was supplied.
	ErrPDFPasswordRequired = errors.New("pdf password required")

	// ErrPDFWrongPassword is returned when the supplied PDF password is incorrect.
	ErrPDFWrongPassword = errors.New("wrong pdf password")

	// ErrPDFCorrupt is returned when a statement PDF cannot be read.
	ErrPDFCorrupt = errors.New("invalid or corrupt pdf")
)
```

Then delete `internal/model/errors.go` entirely (now empty of declarations).

- [ ] **Step 4: Move and re-package the statement test**

```bash
git mv internal/model/tests/statement_test.go internal/transaction/domain/statement_test.go
sed -i '' '1s/^package tests$/package domain_test/' internal/transaction/domain/statement_test.go
```
In that test, requalify: `model.StatementLine`/`model.NewStatement`/`model.ErrBalanceMismatch` → `domain.*`; `model.Transaction`/`model.Inflow`/`model.Outflow` → `shared.*`. Update the import block to:

```go
import (
	"errors"
	"testing"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
)
```

- [ ] **Step 5: Build and test the domain package**

Run: `go test ./internal/transaction/domain/`
Expected: `ok  github.com/kjj1998/kinji/bff/internal/transaction/domain`

- [ ] **Step 6: Commit**

```bash
git add internal/transaction/domain/ internal/model/
git commit -m "refactor: move statement, period, statement/pdf errors to transaction/domain"
```

---

## Task 5: Create `internal/transaction/service`

**Files:**
- Create: `internal/transaction/service/service.go` (from `internal/service/transaction_service.go`)
- Create: `internal/transaction/service/ports.go` (transaction repo interface + `StatementParser` from `internal/service/ports.go`)
- Create: `internal/transaction/service/mock_repository.go` (transaction-method subset from `internal/service/mock_repository.go`)
- Create: `internal/transaction/service/mock_parser.go` (from `internal/service/mock_parser.go`)

- [ ] **Step 1: Move the service and rename package**

```bash
git mv internal/service/transaction_service.go internal/transaction/service/service.go
git mv internal/service/mock_parser.go         internal/transaction/service/mock_parser.go
sed -i '' '1s/^package service$/package service/' internal/transaction/service/service.go internal/transaction/service/mock_parser.go
```
(The package name stays `service` — it's now `internal/transaction/service`.)

- [ ] **Step 2: Create the segregated ports file**

Create `internal/transaction/service/ports.go` with `package service`, holding the **transaction** repository interface and the parser interface, requalified:

```go
package service

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
)

// TransactionRepository is the persistence port for the transaction feature.
type TransactionRepository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error)
	SaveTransactions(ctx context.Context, userId string, transactions []shared.Transaction) error
	GetTransactionPeriods(ctx context.Context, userId string) ([]domain.Period, error)
}

// StatementParser turns a bank-statement PDF into raw extracted rows.
type StatementParser interface {
	Extract(ctx context.Context, pdf []byte, password string, onProgress func(stage string)) ([]domain.StatementLine, error)
}
```
Remove the transaction repo methods and `StatementParser` from `internal/service/ports.go` (the summary subset stays there until Task 10).

- [ ] **Step 3: Requalify `service.go`**

In `internal/transaction/service/service.go`: `model.Transaction` → `shared.Transaction`; `model.Statement`/`model.StatementLine`/`model.NewStatement`/`model.Period`/`model.ErrPDF*` → `domain.*`. Replace the `internal/model` import with `internal/shared` and `internal/transaction/domain` as used.

- [ ] **Step 4: Split the repository mock**

Create `internal/transaction/service/mock_repository.go` (`package service`) containing a mock implementing **only** the three-method `TransactionRepository` above, modeled on the relevant fields of the old `internal/service/mock_repository.go`. Requalify return types to `shared.*`/`domain.*`. Leave the original `internal/service/mock_repository.go` in place (trimmed or untouched) for the summary side until Task 10.

- [ ] **Step 5: Requalify `mock_parser.go`**

In `internal/transaction/service/mock_parser.go`: `model.StatementLine` → `domain.StatementLine`; fix imports.

- [ ] **Step 6: Build and test**

Run: `go test ./internal/transaction/service/`
Expected: `ok` (or `no test files` — then run `go build ./internal/transaction/service/`, expect no output).

- [ ] **Step 7: Commit**

```bash
git add internal/transaction/service/ internal/service/ports.go
git commit -m "refactor: move transaction service + segregated repo/parser ports"
```

---

## Task 6: Create `internal/transaction/store`

**Files:**
- Create: `internal/transaction/store/repository.go` (transaction methods from `internal/store/sqlite.go`)
- Create: `internal/transaction/store/queries.go` (3 transaction SQL consts from `internal/store/queries.go`)

- [ ] **Step 1: Create the transaction repository implementation**

Create `internal/transaction/store/repository.go` (`package store`) holding a `Repository` struct wrapping `*sql.DB`, a `NewRepository(*sql.DB) *Repository` constructor, the compile-time assertion `var _ service.TransactionRepository = (*Repository)(nil)`, and the methods `GetMonthlyTransactions`, `SaveTransactions`, `GetTransactionPeriods` plus helpers `getTransactionsWithinDateRange`, `getPeriods` — all moved verbatim from `internal/store/sqlite.go`. Requalify `model.Transaction`→`shared.Transaction`, `model.Period`→`domain.Period`. Imports:

```go
import (
	"context"
	"database/sql"
	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/transaction/domain"
	"github.com/kjj1998/kinji/bff/internal/transaction/service"
)
```
Calls to `GetMonthRangeDateStrings(...)` → `shared.GetMonthRangeDateStrings(...)`.

- [ ] **Step 2: Move the transaction SQL consts**

Create `internal/transaction/store/queries.go` (`package store`) with `getAllTransactionsWithinDateRange`, `getMonthAndYearWhichTransactionsOccur`, `saveTransactions` moved verbatim. Remove them from `internal/store/queries.go`.

- [ ] **Step 3: Build**

Run: `go build ./internal/transaction/store/`
Expected: builds with no output.

- [ ] **Step 4: Commit**

```bash
git add internal/transaction/store/ internal/store/queries.go
git commit -m "refactor: move transaction repository + queries to transaction/store"
```

---

## Task 7: Create `internal/transaction/parser`

**Files:**
- Create: `internal/transaction/parser/parser.go` (from `internal/parser/parser.go`)
- Create: `internal/transaction/parser/schema.go` (from `internal/parser/schema.go`)

- [ ] **Step 1: Move the parser files**

```bash
git mv internal/parser/parser.go internal/transaction/parser/parser.go
git mv internal/parser/schema.go internal/transaction/parser/schema.go
sed -i '' '1s/^package parser$/package parser/' internal/transaction/parser/parser.go internal/transaction/parser/schema.go
```

- [ ] **Step 2: Requalify model references**

In `internal/transaction/parser/parser.go`: `model.StatementLine` → `domain.StatementLine`; `model.ErrPDF*` → `domain.ErrPDF*`; `model.ParseCategory`/`model.ParseDirection`/`model.Category`/`model.Direction` → `shared.*`. Update imports to `internal/shared` and `internal/transaction/domain` as used.

- [ ] **Step 3: Build**

Run: `go build ./internal/transaction/parser/`
Expected: builds with no output.

- [ ] **Step 4: Commit**

```bash
git add internal/transaction/parser/ internal/parser/
git commit -m "refactor: move pdf parser into transaction/parser"
```

---

## Task 8: Create `internal/transaction/handler`

**Files:**
- Create: `internal/transaction/handler/handler.go` (from `internal/handler/transaction.go`)
- Create: `internal/transaction/handler/dto.go` (from `internal/handler/transaction_dto.go`)
- Create: `internal/transaction/handler/period_dto.go` (from `internal/handler/period_dto.go`)
- Create: `internal/transaction/handler/dto_test.go` (from `internal/handler/transaction_dto_test.go`)
- Create: `internal/transaction/handler/period_dto_test.go` (from `internal/handler/period_dto_test.go`)
- Create: `internal/transaction/handler/mock_service.go` (transaction subset of `internal/handler/mock_service.go`)

- [ ] **Step 1: Move handler + DTO + tests**

```bash
git mv internal/handler/transaction.go          internal/transaction/handler/handler.go
git mv internal/handler/transaction_dto.go      internal/transaction/handler/dto.go
git mv internal/handler/period_dto.go           internal/transaction/handler/period_dto.go
git mv internal/handler/transaction_dto_test.go internal/transaction/handler/dto_test.go
git mv internal/handler/period_dto_test.go      internal/transaction/handler/period_dto_test.go
sed -i '' '1s/^package handler$/package handler/' internal/transaction/handler/*.go
```

- [ ] **Step 2: Requalify references**

In every moved file: `model.Transaction`/`model.Category`/`model.Direction` → `shared.*`; `model.Period`/`model.Statement*` → `internal/transaction/domain`; `service.TransactionService` → `internal/transaction/service`; HTTP helpers formerly in `handler` (e.g. `WriteJSON`, `ClientError` handling) now come from `internal/platform/http` and `internal/shared`. Update import blocks accordingly. Test files use `package handler_test` if they were external; preserve whatever package suffix they currently have (`grep -h '^package' internal/transaction/handler/*_test.go`).

- [ ] **Step 3: Move the transaction-relevant service mock**

Create `internal/transaction/handler/mock_service.go` implementing the `TransactionService` interface (from `internal/transaction/service`), based on the transaction portion of `internal/handler/mock_service.go`. Leave the original for the summary handler until Task 12.

- [ ] **Step 4: Build and test**

Run: `go test ./internal/transaction/handler/`
Expected: `ok github.com/kjj1998/kinji/bff/internal/transaction/handler`

- [ ] **Step 5: Commit**

```bash
git add internal/transaction/handler/ internal/handler/mock_service.go
git commit -m "refactor: move transaction handler + dtos to transaction/handler"
```

---

## Task 9: Create `internal/summary/domain`

**Files:**
- Create: `internal/summary/domain/summary.go` (from `internal/model/summary.go`)
- Create: `internal/summary/domain/spending.go` (from `internal/model/spending.go`)
- Create: `internal/summary/domain/summary_test.go` (from `internal/model/tests/summary_test.go`)

- [ ] **Step 1: Move and rename package**

```bash
git mv internal/model/summary.go  internal/summary/domain/summary.go
git mv internal/model/spending.go internal/summary/domain/spending.go
sed -i '' '1s/^package model$/package domain/' internal/summary/domain/summary.go internal/summary/domain/spending.go
```

- [ ] **Step 2: Requalify `Transaction`/`Category` references**

In `internal/summary/domain/summary.go` and `spending.go`: `Transaction`/`Category`/`Direction`/`Inflow`/`Outflow` → `shared.*` (these came from the same old `model` package). The spending types, `ValueAndChange`, `SummaryCalculator`, and the math helpers stay bare (same package). Add:

```go
import "github.com/kjj1998/kinji/bff/internal/shared"
```
to both files as needed.

- [ ] **Step 3: Move and re-package the summary test**

```bash
git mv internal/model/tests/summary_test.go internal/summary/domain/summary_test.go
sed -i '' '1s/^package tests$/package domain_test/' internal/summary/domain/summary_test.go
```
In the test: `model.SummaryCalculator`/`model.DaySpending`/`model.Category*` (the spending/summary types) and `model.SafeDivide`/`model.PercentageChange`/`model.RoundTo2Dp`/`model.SortByAmountDesc` → `domain.*`; `model.Transaction`/`model.Category`/`model.Inflow`/`model.Outflow`/`model.CategoryEntertainment` etc. → `shared.*`. Set imports:

```go
import (
	"reflect"
	"testing"
	"time"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/summary/domain"
)
```
Note: `Category*` enum constants (e.g. `CategoryEntertainment`, `CategoryFood`) live in `shared` (they're part of `transaction.go`). The `CategorySpendingChange`/`DaySpending`/`MonthSpending` types live in `domain`. Qualify each accordingly.

- [ ] **Step 4: Delete the now-empty tests dir**

```bash
rmdir internal/model/tests 2>/dev/null || true
```

- [ ] **Step 5: Build and test**

Run: `go test ./internal/summary/domain/`
Expected: `ok github.com/kjj1998/kinji/bff/internal/summary/domain`

- [ ] **Step 6: Commit**

```bash
git add internal/summary/domain/ internal/model/
git commit -m "refactor: move summary calculator + spending types to summary/domain"
```

---

## Task 10: Create `internal/summary/service`

**Files:**
- Create: `internal/summary/service/service.go` (from `internal/service/summary_service.go`)
- Create: `internal/summary/service/ports.go` (summary repo interface from remaining `internal/service/ports.go`)
- Create: `internal/summary/service/mock_repository.go` (from remaining `internal/service/mock_repository.go`)

- [ ] **Step 1: Move the service**

```bash
git mv internal/service/summary_service.go internal/summary/service/service.go
git mv internal/service/mock_repository.go  internal/summary/service/mock_repository.go
sed -i '' '1s/^package service$/package service/' internal/summary/service/service.go internal/summary/service/mock_repository.go
```

- [ ] **Step 2: Create the summary repo interface**

Create `internal/summary/service/ports.go` (`package service`) with the summary repository port, requalified to `shared`/`domain`:

```go
package service

import (
	"context"

	"github.com/kjj1998/kinji/bff/internal/shared"
	"github.com/kjj1998/kinji/bff/internal/summary/domain"
)

// TransactionRepository is the read port the summary feature needs.
type TransactionRepository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]shared.Transaction, error)
	GetMonthlyTopMerchants(ctx context.Context, userId, month, year string, limit int) ([]domain.MerchantSpending, error)
	GetMonthlyTopCategories(ctx context.Context, userId, month, year string, limit int) ([]domain.CategorySpending, error)
	GetTotalIncomeTotalSpentAndNetSavings(ctx context.Context, userId, month, year string) (
		totalIncome domain.ValueAndChange[int],
		totalSpent domain.ValueAndChange[int],
		netSavings domain.ValueAndChange[int],
		lastMonthSpent int,
		err error,
	)
	GetLastSixMonthsExpenses(ctx context.Context, userId, month, year string) (map[string]int, error)
	GetCategorySpendingForLastTwoMonths(ctx context.Context, userId, month, year string) (
		current map[shared.Category]int,
		previous map[shared.Category]int,
		err error,
	)
}
```
Then delete `internal/service/ports.go` (now empty).

- [ ] **Step 3: Requalify the service + mock**

In `service.go`: `model.MonthlySummary`/`model.SummaryInput`/`model.SummaryCalculator`/`model.NewSummaryCalculator` → `domain.*`; `model.Transaction` → `shared.Transaction`; spending types → `domain.*`. In `mock_repository.go`: requalify return types to match the interface above. Fix imports.

- [ ] **Step 4: Build and test**

Run: `go test ./internal/summary/service/`
Expected: `ok` (or `no test files` → `go build ./internal/summary/service/`, no output).

- [ ] **Step 5: Commit**

```bash
git add internal/summary/service/ internal/service/
git commit -m "refactor: move summary service + segregated read port"
```

---

## Task 11: Create `internal/summary/store`

**Files:**
- Create: `internal/summary/store/repository.go` (summary methods from `internal/store/sqlite.go`)
- Create: `internal/summary/store/queries.go` (5 summary SQL consts from `internal/store/queries.go`)
- Create: `internal/summary/store/utils.go` (`currentAndPreviousMonth` from `internal/store/utils.go`)

- [ ] **Step 1: Create the summary repository implementation**

Create `internal/summary/store/repository.go` (`package store`) with a `Repository` wrapping `*sql.DB`, `NewRepository(*sql.DB) *Repository`, the assertion `var _ service.TransactionRepository = (*Repository)(nil)` (pointing at `internal/summary/service`), and the methods `GetMonthlyTopMerchants`, `GetTotalIncomeTotalSpentAndNetSavings`, `GetCategorySpendingForLastTwoMonths`, `GetMonthlyTopCategories`, `GetLastSixMonthsExpenses` moved verbatim from `internal/store/sqlite.go`. Requalify `model.Transaction`→`shared.Transaction`, spending/`ValueAndChange` types → `internal/summary/domain`, `model.Category`→`shared.Category`. `GetMonthRangeDateStrings`→`shared.GetMonthRangeDateStrings`; `currentAndPreviousMonth` stays bare (moves in Step 3).

- [ ] **Step 2: Move the summary SQL consts**

Create `internal/summary/store/queries.go` (`package store`) with `getTopSpendingMerchantsWithinDateRange`, `getTotalIncomeTotalSpentAndNetSavingsForTwoMonths`, `getCategorySpendingForTwoMonths`, `getTopSpendingCategoriesWithinDateRange`, `getTotalMonthlyExpensesWithinDateRange` moved verbatim. `internal/store/queries.go` should now be empty — delete it.

- [ ] **Step 3: Move `currentAndPreviousMonth`**

```bash
git mv internal/store/utils.go internal/summary/store/utils.go
sed -i '' '1s/^package store$/package store/' internal/summary/store/utils.go
```
`internal/store/utils.go` now contains only `currentAndPreviousMonth` (`GetMonthRangeDateStrings` left in Task 1). Confirm `internal/store/` has only `sqlite.go` left (now empty of decls) and delete it:

```bash
git rm internal/store/sqlite.go
rmdir internal/store 2>/dev/null || true
```

- [ ] **Step 4: Build**

Run: `go build ./internal/summary/store/`
Expected: builds with no output.

- [ ] **Step 5: Commit**

```bash
git add internal/summary/store/ internal/store/
git commit -m "refactor: move summary repository + queries to summary/store"
```

---

## Task 12: Create `internal/summary/handler`

**Files:**
- Create: `internal/summary/handler/handler.go` (from `internal/handler/summary.go`)
- Create: `internal/summary/handler/dto.go` (from `internal/handler/summary_dto.go`)
- Create: `internal/summary/handler/dto_test.go` (from `internal/handler/summary_dto_test.go`)
- Create: `internal/summary/handler/mock_service.go` (summary subset of `internal/handler/mock_service.go`)

- [ ] **Step 1: Move handler + DTO + test**

```bash
git mv internal/handler/summary.go          internal/summary/handler/handler.go
git mv internal/handler/summary_dto.go      internal/summary/handler/dto.go
git mv internal/handler/summary_dto_test.go internal/summary/handler/dto_test.go
sed -i '' '1s/^package handler$/package handler/' internal/summary/handler/*.go
```

- [ ] **Step 2: Move the summary service mock, delete the old shared mock file**

Move whatever remains of `internal/handler/mock_service.go` (the summary `SummaryService` mock) into `internal/summary/handler/mock_service.go` (`package handler`), then:

```bash
git rm internal/handler/mock_service.go 2>/dev/null || true
```

- [ ] **Step 3: Requalify references**

In every moved summary handler file: `model.MonthlySummary`/spending types → `internal/summary/domain`; `model.Transaction`/`model.Category` → `internal/shared`; `service.SummaryService` → `internal/summary/service`; HTTP helpers → `internal/platform/http`. Update imports. Preserve test package suffix.

- [ ] **Step 4: Confirm `internal/handler/` is empty and remove it**

```bash
ls internal/handler/ 2>/dev/null && echo "NOT EMPTY — investigate" || rmdir internal/handler 2>/dev/null || true
```
Only `server.go` should remain — it moves in Task 13, so `internal/handler/` is NOT empty yet. Do not delete it here.

- [ ] **Step 5: Build and test**

Run: `go test ./internal/summary/handler/`
Expected: `ok github.com/kjj1998/kinji/bff/internal/summary/handler`

- [ ] **Step 6: Commit**

```bash
git add internal/summary/handler/ internal/handler/
git commit -m "refactor: move summary handler + dto to summary/handler"
```

---

## Task 13: Create `internal/server` composition root and update `cmd/api`

**Files:**
- Create: `internal/server/server.go` (from `internal/handler/server.go`)
- Modify: `cmd/api/main.go`

- [ ] **Step 1: Move the server wiring**

```bash
git mv internal/handler/server.go internal/server/server.go
sed -i '' '1s/^package handler$/package server/' internal/server/server.go
rmdir internal/handler 2>/dev/null || true
```

- [ ] **Step 2: Rewrite `server.New` against the new packages**

In `internal/server/server.go`, update the body so it constructs both feature repositories from a shared `*sql.DB`, builds each feature's service + handler, and applies the middleware chain from `internal/platform/http`. Target shape:

```go
package server

import (
	"database/sql"
	"net/http"

	platformhttp "github.com/kjj1998/kinji/bff/internal/platform/http"
	summaryhandler "github.com/kjj1998/kinji/bff/internal/summary/handler"
	summaryservice "github.com/kjj1998/kinji/bff/internal/summary/service"
	summarystore "github.com/kjj1998/kinji/bff/internal/summary/store"
	txnhandler "github.com/kjj1998/kinji/bff/internal/transaction/handler"
	txnservice "github.com/kjj1998/kinji/bff/internal/transaction/service"
	txnstore "github.com/kjj1998/kinji/bff/internal/transaction/store"
	txnparser "github.com/kjj1998/kinji/bff/internal/transaction/parser"
)

func New(db *sql.DB, parser txnservice.StatementParser, corsOrigin string) http.Handler {
	mux := http.NewServeMux()

	txnRepo := txnstore.NewRepository(db)
	summaryRepo := summarystore.NewRepository(db)

	summarySvc := summaryservice.NewSummaryService(summaryRepo)
	txnSvc := txnservice.NewTransactionService(txnRepo, parser)

	txnHandler := txnhandler.NewTransactionHandler(txnSvc)
	summaryHandler := summaryhandler.NewSummaryHandler(summarySvc)

	mux.HandleFunc("GET /health", platformhttp.Health)
	mux.HandleFunc("GET /api/v1/transactions/{id}", txnHandler.GetMonthlyTransactions)
	mux.HandleFunc("POST /api/v1/transactions/{id}", txnHandler.SaveTransactions)
	mux.HandleFunc("GET /api/v1/transactions/{id}/periods", txnHandler.GetPeriods)
	mux.HandleFunc("GET /api/v1/summary/{id}", summaryHandler.Summary)
	mux.HandleFunc("POST /api/v1/transactions/{id}/import", txnHandler.ImportStatement)

	return platformhttp.Chain(mux,
		platformhttp.Recovery,
		platformhttp.Logging,
		platformhttp.CORS(corsOrigin),
	)
}
```
Adjust the exact constructor names / parser-injection shape to match the real signatures from Tasks 5 and 7 (`txnparser` is imported only if `main.go` previously built the parser via the server; if `main` builds the parser itself, drop that import). Keep the route table byte-for-byte identical to the original.

- [ ] **Step 3: Update `cmd/api/main.go`**

Repoint imports: `internal/config`→`internal/platform/config`, `internal/store`→`internal/platform/database`, `internal/handler`→`internal/server`, `internal/parser`→`internal/transaction/parser`. Update call sites: DB via `database.NewClient(...)`, server via `server.New(db, parser, corsOrigin)`. Match the `New` signature chosen in Step 2.

- [ ] **Step 4: Build the whole module**

Run: `go build ./...`
Expected: builds with no output.

- [ ] **Step 5: Commit**

```bash
git add internal/server/ cmd/api/main.go
git commit -m "refactor: add internal/server composition root, rewire cmd/api"
```

---

## Task 14: Final verification and cleanup

- [ ] **Step 1: Confirm the old layered dirs are gone**

```bash
ls internal/model internal/service internal/handler internal/store internal/parser internal/config 2>&1
```
Expected: all report "No such file or directory".

- [ ] **Step 2: Vet and test the whole module**

Run: `go vet ./... && go test ./...`
Expected: no vet output; every package reports `ok` or `no test files`.

- [ ] **Step 3: Verify no stale references remain**

```bash
grep -rn 'internal/model\|internal/service\b\|internal/handler\|internal/store\b\|internal/parser\b\|internal/config\b' --include='*.go' . || echo "clean"
```
Expected: `clean`.

- [ ] **Step 4: Final structure check**

```bash
find internal -type d | sort
```
Expected: matches the target tree in the design doc (`shared`, `platform/{database,http,config}`, `server`, `transaction/{domain,service,store,parser,handler}`, `summary/{domain,service,store,handler}`).

- [ ] **Step 5: Commit any residual cleanup**

```bash
git add -A && git commit -m "refactor: feature-based structure cleanup + final verification" || echo "nothing to commit"
```

---

## Self-Review

**Spec coverage:** Every spec mapping row is covered — shared (T1), platform/database (T2), platform/http+config (T3), transaction domain/service/store/parser/handler (T4–T8), summary domain/service/store/handler (T9–T12), server composition root + cmd/api (T13), old-dir cleanup + acyclic-deps verification (T14). The dependency rules are enforced structurally: features import only `shared`/`platform`; `platform/http` imports `shared` only; `server` is the sole importer of both features.

**Placeholder scan:** No TBD/TODO/"handle edge cases". Verbatim file bodies are moved via explicit `git mv` + `sed` package-rename + an enumerated requalification list, which is the complete actionable content for a move refactor (reproducing unchanged 300-line SQL/handler bodies would add no information).

**Type consistency:** Interface method sets in Tasks 5 and 10 match the original `service.TransactionRepository`, split by the table in the working notes. The summary `TransactionRepository` (T10) uses `domain.MerchantSpending`/`domain.CategorySpending`/`domain.ValueAndChange` and `shared.Category`/`shared.Transaction`, matching the type homes set in T9. `server.New` (T13) references `NewRepository`/`NewSummaryService`/`NewTransactionService`/`NewSummaryHandler`/`NewTransactionHandler` consistent with Tasks 5/8/10/11/12; verify the two handler constructor names against the originals during T8/T12.
