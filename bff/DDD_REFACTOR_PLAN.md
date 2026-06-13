# Refactor BFF to DDD + Onion Architecture (Ideal)

## Context

The `bff` is a layered Go service. This plan migrates it **in place** to the
ideal DDD/onion target (matching `DDD_GREENFIELD_STRUCTURE.md`), not a pragmatic
halfway point. Five things in the current code break the dependency rule or
misplace business logic:

1. **Ports are owned by adapters.** `repository.Repository`
   (`internal/repository/repository.go`) and `claude.Parser`
   (`internal/claude/parser.go`) live next to their implementations. The core
   should own the contracts.
2. **Presentation DTOs leak inward.** The sqlite adapter returns `dto.Merchant`,
   `dto.CategorySpending`, `dto.ValueAndChange`
   (`internal/repository/sqlite/sqlite.go`), coupling persistence to the wire.
3. **Anemic domain.** `models.Transaction` is a JSON-tagged struct that doubles
   as the HTTP response; all calculation lives as free functions in `service/`.
4. **A business invariant lives in an adapter.** The running-balance guard sits
   inside the Claude adapter (`internal/claude/parser.go`), not in the domain.
5. **Infrastructure leaks into the use case + presentation leaks into logic.**
   The import service calls `pdfcpu` directly
   (`internal/service/transaction.go`), and the summary logic bakes view
   concerns (`"Mon"`/`"Jan"` labels, top-3, recent-5) into `internal/service/summary.go`.

**Goal:** `domain` owns entities, value objects, the balance invariant, and the
summary calculations; `app` owns use cases + ports; `adapter/*` implement ports
and own all I/O (sqlite, Anthropic, pdfcpu, HTTP/JSON). Module path:
`github.com/kjj1998/kinji/bff`.

## The dependency rule

```
domain  ←  app  ←  adapter        cmd/api wires everything
```
`domain` imports nothing under `internal/` and no I/O libs (no SQL, HTTP,
Anthropic SDK, pdfcpu, or `json` tags). `app` imports only `domain`. `adapter/*`
import `app` + `domain`. `cmd/api` imports everything.

## Target layout

```
internal/
  domain/                   # pure model — no internal imports, no I/O, no json tags
    transaction.go            # Transaction; Category, Direction (+IsValid, IsOutflow/IsInflow)
    money.go                  # Money value object (cents) — OPTIONAL, see note
    period.go                 # Period
    statement.go              # Statement aggregate + running-balance invariant (Validate -> ErrBalanceMismatch)
    spending.go               # ValueAndChange[T], MerchantSpending, CategorySpending,
                              #   CategorySpendingChange, DaySpending{Weekday,Amount}, MonthSpending{Month,Amount}
    summary.go                # MonthlySummary read-model + SummaryCalculator (pure, raw-typed results)
    errors.go                 # ErrInvalidCategory, ErrInvalidDirection, ErrBalanceMismatch,
                              #   ErrPDFPasswordRequired, ErrPDFWrongPassword, ErrPDFCorrupt
  app/                      # use cases + ports. imports domain only
    ports.go                  # TransactionRepository, StatementParser (owned here)
    transaction_service.go    # import / save / get-monthly / get-periods
    summary_service.go        # orchestrate repo reads -> SummaryCalculator
    errors.go                 # ClientError (bad input -> 4xx)
  adapter/
    persistence/sqlite/       # implements TransactionRepository, returns domain types
    parser/claude/            # implements StatementParser: pdfcpu decrypt/validate + Anthropic extraction
    http/
      server.go               # router + middleware chain
      handler/                # transaction.go, summary.go, health.go
      dto/                    # request/response (json tags) + mappers + label formatting + top-N truncation
      middleware/             # logging, recovery, cors
  config/
cmd/api/main.go             # composition root
```

## Changes by layer

### 1. `internal/domain` (new)
- Move `internal/models/transaction.go` here; **strip JSON tags**. Keep
  `Category`/`Direction` + `IsValid()`; add `IsOutflow()`/`IsInflow()`. Add
  validating constructors `ParseCategory`/`ParseDirection` returning
  `ErrInvalidCategory`/`ErrInvalidDirection` (the claude adapter will use these
  instead of inline checks). Move `Period`.
- **`statement.go` (new — moves a real invariant into the core).** Define
  `StatementLine { Txn Transaction; Balance int }`, a `Statement` built from
  `[]StatementLine`, `Validate() error` enforcing the running-balance arithmetic
  (currently in `internal/claude/parser.go` ~L176–188) and returning
  `ErrBalanceMismatch`, plus `Transactions() []Transaction`.
- `spending.go`: move from `dto` as **value objects (no json tags, no labels)**.
  Trend types carry raw keys — `DaySpending{Weekday time.Weekday; Amount int}`,
  `MonthSpending{Month time.Time; Amount int}` — not `"Mon"`/`"Jan"` strings.
- `summary.go`: `MonthlySummary` read-model (raw-typed, full ordered lists) +
  `SummaryCalculator` holding the moved pure functions from
  `internal/service/summary.go` (`generateMonthlySummary`,
  `computeDailySpendingTrend`, `computeCategoriesWithBiggestSpendingChange`,
  `recentTransactions`, `buildMonthlyTrend`, `getTopTransaction`) and the math
  helpers from `internal/service/math_utils.go`. **Ranking stays here; the
  truncation counts (top-3, recent-5) and labels move to the dto mapper.**
- `errors.go`: the sentinel errors listed above.

### 2. `internal/app` (new)
- `ports.go`: `TransactionRepository` (was `repository.Repository`, returning
  `domain` types not `dto`) and `StatementParser` — signature
  `Extract(ctx, pdf []byte, password string, onProgress func(string)) ([]domain.StatementLine, error)`.
  The parser returns **raw rows incl. balance**; it no longer validates balances.
- `transaction_service.go` (was `internal/service/transaction.go`): import use
  case becomes — read bytes → `parser.Extract(...)` → `domain.NewStatement(lines)`
  → `stmt.Validate()` (ErrBalanceMismatch) → stamp `UserID` + `uuid.NewV7()` on
  `stmt.Transactions()`. **No `pdfcpu` import here anymore** (moved to the
  adapter). Save / get-monthly / get-periods unchanged except types.
- `summary_service.go`: orchestration only — repo reads → `SummaryCalculator` →
  `*domain.MonthlySummary`.
- `errors.go`: `ClientError` (moved from `internal/service/errors.go`).

### 3. `internal/adapter/persistence/sqlite` (from `internal/repository/sqlite`)
- Implements `app.TransactionRepository`; replace all `dto.*` scans/returns with
  domain value objects. **Delete** `internal/repository/repository.go`. Keep
  `queries.go` / `daterange_utils.go`.

### 4. `internal/adapter/parser/claude` (from `internal/claude`)
- Owns **both** PDF handling and extraction: move the `pdfcpu`
  validate/decrypt block out of the service into here, mapping failures to
  `domain.ErrPDFPasswordRequired` / `ErrPDFWrongPassword` / `ErrPDFCorrupt`.
  Run the Anthropic extraction, map rows via `domain.ParseCategory/ParseDirection`,
  and return `[]domain.StatementLine` (Transaction + Balance). **Remove the
  balance-arithmetic guard** (now `domain.Statement.Validate`). Prompt + tool
  schema unchanged.

### 5. `internal/adapter/http` (from handler/dto/middleware/server)
- `dto/`: presentation-only. `TransactionDTO` + `To/FromTransactionDTO`;
  `TransactionSummary` + `ValueAndChange` DTOs + mapper from
  `domain.MonthlySummary`. **The mapper now owns** weekday/month label rendering
  (`"Mon"`, `"Jan"`) and the view truncation (top-3 changes, recent-5,
  6-month window) — preserving today's exact JSON field names/shape.
- `handler/`: call `app` services, map domain→dto. `SaveTransactions` decodes
  into `TransactionDTO` → domain. Map errors to status: `domain.ErrPDF*` and
  `app.ClientError` → 4xx (via `errors.As`/`errors.Is`), else 500.
- `middleware/`, `server/`: move as-is; `server.New` wires `app` services.

### 6. `cmd/api/main.go`
- Wire `adapter/persistence/sqlite`, `adapter/parser/claude`,
  `adapter/http/server`. Shape unchanged.

### 7. Tests & mocks (currently untracked)
- `summary_test.go` + `math_utils_test.go` → `internal/domain`, updated to the
  raw-typed (`time.Weekday`/`time.Month`) results.
- **New** `internal/domain` test for `Statement.Validate` covering the
  balance-mismatch case (migrated from the adapter's old guard behavior).
- **New** `adapter/http/dto` test asserting the mapper still emits `"Mon"`/`"Jan"`
  labels and top-3/recent-5 truncation (locks the wire format).
- `mock_repository.go`, `mock_parser.go` → `app`; `mock_service.go` →
  `adapter/http/handler`. Update signatures to the domain-typed ports.

## Money value object (optional)

`amount`/`balance` could become a `domain.Money` (cents + `Dollars()`, `Add`,
etc.) for full tactical DDD. Adds mapping churn at every boundary; defer unless
money math spreads. Listed so the choice is explicit, not accidental.

## Verification

- `cd bff && go build ./...`, `go vet ./...`, `go test ./...`.
- **Dependency-rule check:** confirm `internal/domain` imports no other
  `internal/*` (and no `pdfcpu`/`anthropic`/`database/sql`), and `internal/app`
  imports only `internal/domain`.
- **Invariant check:** `Statement.Validate` rejects a tampered balance sequence
  (unit test) — proves the guard survived the move into the domain.
- **End-to-end smoke:** run `go run ./cmd/api`, exercise
  `GET /api/v1/transactions/{id}`, `GET /api/v1/summary/{id}`,
  `POST /api/v1/transactions/{id}`, and `POST /.../import` (SSE). The summary
  JSON (labels, top-3, recent-5) and all wire shapes must be byte-for-byte
  unchanged after introducing the dto mappers.
```
