# Flatten + Segregate Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reduce navigation depth and fix the repository god-interface while preserving the onion/DDD dependency rule the codebase already enforces.

**Architecture:** Keep `domain ← app ← <adapters>` with `cmd/api` as composition root. Two changes only: (1) drop the `adapter/` ring noun so adapter packages sit directly under `internal/` (`internal/http`, `internal/sqlite`, `internal/parser`); (2) split `app.TransactionRepository` (8 methods, shared by two services) into per-consumer interfaces (`LedgerRepository`, `SummaryRepository`) plus a `Repository` composite for wiring. No behavior changes, no wire-format changes.

**Tech Stack:** Go 1.25, stdlib `net/http`, sqlite adapter, Anthropic (claude) parser adapter.

---

## Why these changes (problem statement)

- **Depth:** `internal/adapter/http/handler/summary.go` is 4 levels under `internal/`. The `adapter/` and `persistence/` nouns carry no information — everything under them is obviously an adapter.
- **ISP violation:** [internal/app/ports.go](internal/app/ports.go) declares one 8-method `TransactionRepository` used by *both* `transactionService` and `summaryService`. The summary service depends on `SaveTransactions`/`GetTransactionPeriods` it never calls. Verified usage:
  - Ledger (transaction_service): `GetMonthlyTransactions`, `GetTransactionPeriods`, `SaveTransactions`
  - Summary (summary_service): `GetMonthlyTransactions`, `GetMonthlyTopMerchants`, `GetMonthlyTopCategories`, `GetTotalIncomeTotalSpentAndNetSavings`, `GetLastSixMonthsExpenses`, `GetCategorySpendingForLastTwoMonths`
  - `GetMonthlyTransactions` is the only overlap.

## Target layout

```
cmd/api/main.go                  composition root (unchanged role)
internal/
  domain/                        unchanged
  app/                           unchanged role; ports.go interface split
  http/                          was internal/adapter/http/
    dto/  handler/  middleware/  server/
  sqlite/                        was internal/adapter/persistence/sqlite/
  parser/                        was internal/adapter/parser/claude/ (package claude -> parser)
  config/                        unchanged
```

`internal/adapter/` is deleted once empty.

## Files touched (map)

- **Moved (git mv, contents unchanged except imports/package):** everything under `internal/adapter/`.
- **Import rewrites:** [cmd/api/main.go](cmd/api/main.go), [internal/http/handler/summary.go](internal/http/handler/summary.go), [internal/http/handler/transaction.go](internal/http/handler/transaction.go), [internal/http/server/server.go](internal/http/server/server.go).
- **Interface split:** [internal/app/ports.go](internal/app/ports.go), [internal/app/transaction_service.go](internal/app/transaction_service.go), [internal/app/summary_service.go](internal/app/summary_service.go), [internal/app/mock_repository.go](internal/app/mock_repository.go), [internal/http/server/server.go](internal/http/server/server.go), [cmd/api/main.go](cmd/api/main.go).
- **Docs:** [DDD_GREENFIELD_STRUCTURE.md](DDD_GREENFIELD_STRUCTURE.md), [DDD_REFACTOR_PLAN.md](DDD_REFACTOR_PLAN.md).

> **Note on `sed`:** commands below use BSD/macOS syntax (`sed -i ''`). Run from the `bff/` module root.

---

### Task 0: Capture green baseline

**Files:** none (verification only).

- [ ] **Step 1: Confirm everything builds and passes before touching anything**

Run: `cd bff && go build ./... && go vet ./... && go test ./...`
Expected: build/vet clean; tests `ok` for `internal/domain` and `internal/adapter/http/dto` (and any others). This is the regression net for every later task.

---

### Task 1: Move `internal/adapter/http` → `internal/http`

**Files:**
- Move: `internal/adapter/http/` → `internal/http/`
- Modify imports: `cmd/api/main.go`, `internal/http/handler/summary.go`, `internal/http/handler/transaction.go`, `internal/http/server/server.go`

- [ ] **Step 1: Move the directory tree (preserves git history)**

Run:
```bash
git mv internal/adapter/http internal/http
```

- [ ] **Step 2: Rewrite every import that referenced the old path**

Run:
```bash
grep -rl "internal/adapter/http" --include="*.go" . | xargs sed -i '' 's#internal/adapter/http#internal/http#g'
```

- [ ] **Step 3: Verify build + tests green**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: clean build; `internal/http/dto` tests `ok`. (sqlite/parser still under `internal/adapter/` — their imports are untouched and still valid.)

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "refactor: move adapter/http to internal/http

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 2: Move `internal/adapter/persistence/sqlite` → `internal/sqlite`

**Files:**
- Move: `internal/adapter/persistence/sqlite/` → `internal/sqlite/`
- Modify imports: `cmd/api/main.go`

- [ ] **Step 1: Move the directory**

Run:
```bash
git mv internal/adapter/persistence/sqlite internal/sqlite
rmdir internal/adapter/persistence
```

- [ ] **Step 2: Rewrite imports**

Run:
```bash
grep -rl "internal/adapter/persistence/sqlite" --include="*.go" . | xargs sed -i '' 's#internal/adapter/persistence/sqlite#internal/sqlite#g'
```

Package name stays `sqlite`; only the import path changes, so `sqlite.NewClient` / `sqlite.NewRepository` call sites are unaffected.

- [ ] **Step 3: Verify build + tests green**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: clean.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "refactor: move adapter/persistence/sqlite to internal/sqlite

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 3: Move `internal/adapter/parser/claude` → `internal/parser` (rename package)

**Files:**
- Move: `internal/adapter/parser/claude/` → `internal/parser/`
- Modify package decls: `internal/parser/parser.go`, `internal/parser/schema.go`
- Modify imports + call site: `cmd/api/main.go`

> Rationale: package name must match its directory leaf to stay idiomatic. The dir becomes `parser`, so package `claude` → `parser`. There is only one parser implementation today (YAGNI); if a second is ever added, reintroduce `internal/parser/<impl>/`.

- [ ] **Step 1: Move the directory**

Run:
```bash
git mv internal/adapter/parser/claude internal/parser
rmdir internal/adapter/parser
rmdir internal/adapter
```

- [ ] **Step 2: Rewrite the import path**

Run:
```bash
grep -rl "internal/adapter/parser/claude" --include="*.go" . | xargs sed -i '' 's#internal/adapter/parser/claude#internal/parser#g'
```

- [ ] **Step 3: Rename the package in both files**

In `internal/parser/parser.go` and `internal/parser/schema.go`, change the first line:

```go
package claude
```
to:
```go
package parser
```

Run:
```bash
sed -i '' 's/^package claude$/package parser/' internal/parser/parser.go internal/parser/schema.go
```

- [ ] **Step 4: Update the call site in main.go**

In `cmd/api/main.go`, change:
```go
	parser := claude.NewParser(cfg.AnthropicModel)
```
to:
```go
	parser := parser.NewParser(cfg.AnthropicModel)
```

Note: the local variable is also named `parser`; that shadows the package name only *after* this line, so the call itself resolves to the package. Verified safe because the package reference and assignment are on the same line.

- [ ] **Step 5: Verify build + tests green**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: clean. `internal/adapter/` no longer exists.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: move adapter/parser/claude to internal/parser

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 4: Split the repository god-interface (ISP)

**Files:**
- Modify: `internal/app/ports.go`
- Modify: `internal/app/transaction_service.go:23` (field type)
- Modify: `internal/app/summary_service.go:20` (field type) and `:24` (ctor param)
- Modify: `internal/app/transaction_service.go:27` (ctor param)
- Modify: `internal/app/mock_repository.go:22` (compile-time assertion)
- Modify: `internal/http/server/server.go:11` (param type)
- Modify: `cmd/api/main.go:25` (var type)

- [ ] **Step 1: Replace the single interface in `ports.go` with three**

In [internal/app/ports.go](internal/app/ports.go), replace the entire `TransactionRepository` interface block (lines 9–30) with:

```go
// LedgerRepository is the persistence port the transaction use cases need:
// storing reviewed transactions and reading them back by period.
type LedgerRepository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]domain.Transaction, error)
	SaveTransactions(ctx context.Context, userId string, transactions []domain.Transaction) error
	GetTransactionPeriods(ctx context.Context, userId string) ([]domain.Period, error)
}

// SummaryRepository is the persistence port the summary use case needs: the
// aggregate reads that feed the monthly insights projection.
type SummaryRepository interface {
	GetMonthlyTransactions(ctx context.Context, userId, month, year string) ([]domain.Transaction, error)
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
		current map[domain.Category]int,
		previous map[domain.Category]int,
		err error,
	)
}

// Repository is the full persistence port a single concrete store satisfies.
// It exists only for the composition root and server wiring, which hold one
// repository value and hand it to both use cases as the narrower interface.
type Repository interface {
	LedgerRepository
	SummaryRepository
}
```

Leave the `StatementParser` interface below it unchanged.

- [ ] **Step 2: Narrow the service field + ctor types**

In `internal/app/transaction_service.go`, change the field (line ~23) and constructor (line ~27):
```go
	repo   TransactionRepository
```
→
```go
	repo   LedgerRepository
```
and
```go
func NewTransactionService(repo TransactionRepository, parser StatementParser) TransactionService {
```
→
```go
func NewTransactionService(repo LedgerRepository, parser StatementParser) TransactionService {
```

In `internal/app/summary_service.go`, change the field (line ~20) and constructor (line ~24):
```go
	repo TransactionRepository
```
→
```go
	repo SummaryRepository
```
and
```go
func NewSummaryService(repo TransactionRepository) SummaryService {
```
→
```go
func NewSummaryService(repo SummaryRepository) SummaryService {
```

- [ ] **Step 3: Update the mock compile-time assertion**

In `internal/app/mock_repository.go` (line ~22), change:
```go
var _ TransactionRepository = (*MockRepository)(nil)
```
→
```go
var _ Repository = (*MockRepository)(nil)
```

`MockRepository` keeps all 8 `…Fn` fields and methods; satisfying `Repository` proves it satisfies both narrow interfaces.

- [ ] **Step 4: Update the wiring types**

In `internal/http/server/server.go` (line 11), change:
```go
func New(repo app.TransactionRepository, parser app.StatementParser, corsOrigin string) http.Handler {
```
→
```go
func New(repo app.Repository, parser app.StatementParser, corsOrigin string) http.Handler {
```

In `cmd/api/main.go` (line 25), change:
```go
	var repo app.TransactionRepository
```
→
```go
	var repo app.Repository
```

- [ ] **Step 5: Verify build + tests green**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: clean. The single `sqlite.Repository` value still satisfies `app.Repository`, so wiring is unchanged at runtime.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: split repository port per consumer (ISP)

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 5: Update the architecture docs to match

**Files:**
- Modify: `DDD_GREENFIELD_STRUCTURE.md` (Target layout section, lines ~34–65)
- Modify: `DDD_REFACTOR_PLAN.md` (Target layout section, lines ~42–69)

- [ ] **Step 1: Update both "Target layout" code blocks**

In each doc, replace the `internal/adapter/{persistence/sqlite,parser/claude,http}` tree with the flattened layout:
```
internal/
  domain/
  app/                  # ports split: LedgerRepository, SummaryRepository, Repository
  http/                 # server, handler/, dto/, middleware/
  sqlite/               # implements LedgerRepository + SummaryRepository
  parser/               # implements StatementParser (package parser)
  config/
cmd/api/main.go
```

- [ ] **Step 2: Add a short note recording the deviation**

Append to `DDD_GREENFIELD_STRUCTURE.md` under a new heading:
```markdown
## Update (flattening)

The `adapter/` ring noun was dropped for navigability: adapters now sit directly
under `internal/` (`internal/http`, `internal/sqlite`, `internal/parser`). The
dependency rule is unchanged — these are still the outer ring; only the path
depth changed. The single `TransactionRepository` was split per consumer
(`LedgerRepository`, `SummaryRepository`) with a `Repository` composite for wiring.
```

- [ ] **Step 3: Commit**

```bash
git add DDD_GREENFIELD_STRUCTURE.md DDD_REFACTOR_PLAN.md
git commit -m "docs: update DDD layout docs for flattened structure

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 6: Final verification

**Files:** none (verification only).

- [ ] **Step 1: Full build/vet/test**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: all clean.

- [ ] **Step 2: Dependency-rule check (the onion must still hold)**

Run:
```bash
go list -deps ./internal/domain | grep "kinji/bff/internal" || echo "domain imports no internal packages: OK"
go list -f '{{join .Imports "\n"}}' ./internal/app | grep "kinji/bff/internal" | grep -v "internal/domain" && echo "LEAK: app imports non-domain internal" || echo "app imports only domain: OK"
```
Expected: `domain imports no internal packages: OK` and `app imports only domain: OK`.

- [ ] **Step 3: Confirm the old tree is gone**

Run: `test ! -d internal/adapter && echo "internal/adapter removed: OK"`
Expected: `internal/adapter removed: OK`.

- [ ] **Step 4: End-to-end smoke (manual)**

Run `go run ./cmd/api`, then exercise:
- `GET /api/v1/transactions/{id}`
- `GET /api/v1/summary/{id}` — JSON must still emit `"Mon"`/`"Jan"` labels and top-3/recent-5 truncation
- `POST /api/v1/transactions/{id}`
- `POST /api/v1/transactions/{id}/import` (SSE progress)

Expected: identical responses to pre-refactor; no wire-format change.

---

## Self-review notes

- **Scope:** pure structural move + interface split; zero behavior/wire changes — existing tests in `internal/domain` and `internal/http/dto` are the regression net.
- **Type consistency:** `LedgerRepository` / `SummaryRepository` / `Repository` names used identically across `ports.go`, both services, the mock assertion, `server.New`, and `main.go`.
- **No new tests required:** no new behavior is introduced. If desired later, add `internal/app` service tests using the existing `MockRepository`/`MockParser` (currently unused) — out of scope here.
- **Risk:** the `parser` package rename (Task 3) and the same-line variable shadow in `main.go` are the only non-mechanical spots; Step 5 of Task 3 catches any mistake at build time.
