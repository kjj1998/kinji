# Test Plan — Go BFF (`/bff`)

## Context

The `bff` service is a clean-architecture Go backend (`net/http` + `ServeMux`) that currently has
**zero test coverage**. It exposes 6 HTTP endpoints over transactions/summaries and parses PDF bank
statements via the Claude API. The codebase is well-suited to testing: every seam is an interface
(`repository.Repository`, `claude.Parser`, `service.TransactionService`, `service.SummaryService`)
and dependencies are injected explicitly (no globals/singletons).

Goal: establish a full-stack test suite using **stdlib only** (`testing`, `net/http/httptest`,
in-memory SQLite via the already-vendored `modernc.org/sqlite`) — no new dependencies. Tests should
lock down the request/response contracts, the summary aggregation math, validation behavior, and the
SQL queries before the API grows further.

Scope: handlers, services, repository, and the parser's post-processing. **Excludes** live Anthropic
API calls (the `Parser` interface is mocked at the service layer).

## Approach

Add `_test.go` files alongside the code they test (same package for internal helpers, `_test` package
where only the exported surface matters). Use table-driven tests throughout. Hand-write small mocks
that satisfy the existing interfaces — no mocking framework.

### Shared test infrastructure (build first)

Create `internal/repository/mock_repository.go` (or a `*_test.go` test helper):
- `mockRepository` struct implementing all 8 methods of [repository.Repository](internal/repository/repository.go)
  via configurable func fields, e.g. `GetMonthlyTransactionsFn func(...) ([]models.Transaction, error)`.
  This is the workhorse for both service and (indirectly) handler tests.

Create a `mockParser` implementing [claude.Parser](internal/claude/parser.go) `ParseStatement`.

Create `mockTransactionService` / `mockSummaryService` implementing the two service interfaces for
handler tests (handlers depend on service interfaces, not concrete structs — see
[server.go:17-22](internal/server/server.go#L17-L22)).

A small fixture helper that builds `[]models.Transaction` (varied dates/categories/directions) keeps
the summary tests readable.

---

### Layer 1 — Pure functions (highest value, easiest)

These are deterministic and need no mocks. Test in-package (`package service` / `package sqlite`).

**`internal/service/math_utils_test.go`** — [math_utils.go](internal/service/math_utils.go)
- `safeDivide`: divide by zero returns 0; normal division; negative operands.
- `percentageChange`: increase/decrease, `previous == 0` (returns 0 via safeDivide), rounding to 2dp.
- `roundTo2Dp`: half-up rounding, already-rounded values.
- `sortByAmountDesc`: descending order, stable on ties, empty slice.

**`internal/service/summary_test.go`** (helpers) — [summary.go](internal/service/summary.go)
- `generateMonthlySummary`: `topTransaction == nil` → `""`; `hasPrevMonth == false` → suffix only;
  `difference < 0` → "less"; `difference > 0` → "more"; verify formatted string content.
- `computeDailySpendingTrend`: skips `INFLOW`; skips unparseable dates; groups by weekday Mon→Sun;
  always returns 7 buckets with 3-letter day labels.
- `computeCategoriesWithBiggestSpendingChange`: `IsNew` when prev==0; `len(prev)==0` falls back to
  sort-by-amount; with baseline sorts by abs(percentageChange); caps at top 3; `<3` returns all.
- `recentTransactions`: nil → empty slice (not nil); sorts date-desc; `n > len` clamps; original
  slice not mutated (it clones).
- `buildMonthlyTrend`: returns 6 buckets ending at the target month; `Jan`-style labels; missing
  months default to 0; invalid month/year string → error.
- `getTopTransaction`: empty → nil; ignores `INFLOW`; picks max-amount `OUTFLOW`.

**`internal/repository/sqlite/daterange_utils_test.go`** — [daterange_utils.go](internal/repository/sqlite/daterange_utils.go)
- `GetMonthRangeDateStrings`: correct first/last day; Feb (28), Dec→year boundary on the `to` side.
- `currentAndPreviousMonth`: Jan→prev Dec of prior year; mid-year; invalid input → error.

**`internal/models/transaction_test.go`** — [transaction.go](internal/models/transaction.go)
- `Category.IsValid` / `Direction.IsValid`: each valid value true; unknown/empty false.

**`internal/dto/dto_test.go`** — `NewValueAndChange([current, previous])`: `Change = current - previous`;
JSON omits `change` when zero (omitempty).

---

### Layer 2 — Services (mock repository + parser)

**`internal/service/summary_service_test.go`** — `GenerateMonthlySummary`
- Happy path: wire `mockRepository` to return canned data across all 7 calls; assert the assembled
  `TransactionSummary` (savings rate = `netSavings/totalIncome*100` rounded; trend/changes wired
  through). Confirms the orchestration in [summary.go:33-98](internal/service/summary.go#L33-L98).
- Error propagation: make each of the 7 repo calls fail in turn → assert wrapped error returned and
  no panic. (Table-driven over "which call fails".)
- `totalIncome.Value == 0` → savings rate 0 (no divide-by-zero).

**`internal/service/transaction_service_test.go`**
- `GetMonthlyTransactions` / `GetPeriods` / `SaveTransactions`: delegate to repo; pass-through on
  success; wrapped error on repo failure.
- `ImportStatement` (the meatiest — [transaction.go:50-104](internal/service/transaction.go#L50-L104)):
  - Unencrypted valid PDF: needs a real tiny valid PDF fixture (generate once with `pdfcpu` or commit
    a minimal `testdata/sample.pdf`). Assert `onProgress` called with `uploaded → validating →
    parsing` in order, parser invoked, and every returned txn gets a non-empty `UserID` and UUIDv7 `ID`.
  - Encrypted PDF + correct password → decrypts then parses; wrong password → `ClientError{"wrong pdf
    password given"}`; missing password on encrypted PDF → `ClientError{"pdf password required"}`.
  - Corrupt bytes → `ClientError{"invalid/corrupt pdf file"}`.
  - Parser returns error → wrapped `parse statement` error.
  - Note: the ID/UserID loop runs before the parser-error check (lines 94-101); a test documenting
    this ordering is worthwhile.
  - `testdata/` should hold: one valid PDF, one password-protected PDF, one garbage file.

---

### Layer 3 — Repository (in-memory SQLite, real SQL)

**`internal/repository/sqlite/sqlite_test.go`** — use `NewClient(":memory:")` (or a temp-file DB);
seed via `SaveTransactions`, then assert queries. Verifies the hand-written SQL in
[queries.go](internal/repository/sqlite/queries.go) and [sqlite.go](internal/repository/sqlite/sqlite.go).

Per-method cases:
- `SaveTransactions`: batch insert succeeds; round-trips via `GetMonthlyTransactions`; the `direction`
  CHECK constraint rejects invalid values; verify it's transactional (a bad row rolls back the batch).
- `GetMonthlyTransactions`: only rows within the month's date range; ordered date-desc; empty month →
  empty (non-nil) slice; correct `user_id` scoping (other users excluded).
- `GetMonthlyTopMerchants` / `GetMonthlyTopCategories`: SUM grouped + ordered desc; honors `limit`.
- `GetTotalIncomeTotalSpentAndNetSavings`: income vs spent split by direction; `Change` vs previous
  month; `lastMonthSpent` correct; month with no data → zeros.
- `GetCategorySpendingForLastTwoMonths`: returns two maps split correctly across the month boundary.
- `GetLastSixMonthsExpenses`: keys are `YYYY-MM`, outflow only, only months with data present.
- `GetTransactionPeriods`: distinct year→months grouping; multiple years.

Use a shared seed fixture spanning ≥2 months and ≥2 users so date-range and scoping are exercised.

---

### Layer 4 — HTTP handlers (`httptest`, mocked services)

Drive handlers through the real router so `{id}` path values and method routing are exercised. Build
the handler via [server.New(...)](internal/server/server.go#L13) with a mock repo/parser, OR
construct `handler.NewTransactionHandler(mockTxService)` directly and register on a `ServeMux` with the
`GET /api/v1/transactions/{id}` patterns so `r.PathValue("id")` resolves.

**`internal/handler/transaction_test.go`** & **`summary_test.go`** & **`health_test.go`**

Cross-cutting (table-driven, applies to most endpoints):
- Missing/empty `{id}` → 400 `{"error":"User ID not provided"}` (via
  [requireUserId](internal/handler/http_utils.go#L25)).
- `parseMonthYear` validation ([http_utils.go:34-58](internal/handler/http_utils.go#L34-L58)) on
  `GetMonthlyTransactions` and `Summary`: invalid month (`0`,`13`,`abc`) → 400 "invalid month…";
  invalid year (`999`,`10000`,`abc`) → 400 "invalid year…"; both empty → defaults to current month/year
  (assert service called with `time.Now` formatted values — inject via the mock and assert args).
- Service returns error → 500 with the endpoint's message ("failed to get monthly transactions", etc.).

Per endpoint:
- `Health`: 200, body `{"status":"ok"}`, `Content-Type: application/json`.
- `GetMonthlyTransactions`: 200 + JSON array body; empty result still valid JSON.
- `SaveTransactions`: valid `[]Transaction` JSON → 200 echoing saved txns; invalid JSON → 400
  "invalid request body"; **unknown field** → 400 (DisallowUnknownFields); body over 1MB rejected.
- `GetPeriods`: 200 + `[]Period`.
- `Summary`: 200 + `TransactionSummary` JSON shape.
- `ImportStatement` (SSE — [transaction.go:51-104](internal/handler/transaction.go#L51-L104)):
  - Missing `{id}` → 400 before streaming.
  - Non-multipart / malformed form → SSE `event: error` `{"message":"invalid multipart form"}`.
  - Missing `statement` file part → `{"message":"missing statement file"}`.
  - Happy path (mock service returns txns): response `Content-Type: text/event-stream`; stream contains
    `event: progress` frames then a terminal `event: done` whose `data` is the txn JSON.
  - Service `ClientError` vs generic error → corresponding `event: error` message. Parse the SSE stream
    from `httptest.ResponseRecorder.Body` (split on `\n\n`, assert `event:`/`data:` lines).

**`internal/middleware/middleware_test.go`** — [middleware.go](internal/middleware/middleware.go)
- CORS: `Access-Control-Allow-Origin` echoes configured origin; OPTIONS → 204; methods/headers set.
- Recovery: a handler that panics → 500 "internal server error", no crash.
- Logging: wraps and calls through (assert next handler invoked; status captured).

---

### Layer 5 — Claude parser internals (optional / lower priority)

[parser.go](internal/claude/parser.go) `ParseStatement` calls `anthropic.Client` directly, so it
can't be unit-tested without an HTTP transport stub. The valuable logic is the **balance guard** and
category/direction validation in the post-processing block (lines ~176-213). Two options, in order of
preference:
1. **Refactor**: extract post-processing (`toolInput → []Transaction` + balance verification) into a
   pure function and table-test it directly (balance matches; off-by-one mismatch → error with row #;
   invalid category/direction rejected). Lowest cost, highest value.
2. If no refactor: point the anthropic SDK at an `httptest.Server` returning a canned tool-use response
   and assert the mapping. Heavier; defer.

Treat this layer as a follow-up; the `Parser` interface is already mocked at Layer 2, so the import
flow is covered without it.

---

## Suggested order of execution

1. Shared mocks + fixtures.
2. Layer 1 pure functions (fast wins, no mocks).
3. Layer 3 repository (in-memory SQLite) — locks down the SQL.
4. Layer 2 services.
5. Layer 4 handlers + middleware.
6. Layer 5 parser refactor + tests (optional follow-up).

## Verification

- Run the full suite: `cd bff && go test ./...`
- With race detector (SSE/goroutines, server shutdown): `go test -race ./...`
- Coverage report: `go test -coverprofile=cover.out ./... && go tool cover -func=cover.out`
  (target meaningful coverage on `internal/service`, `internal/repository/sqlite`, `internal/handler`).
- `go vet ./...` stays clean.
- Confirm no new entries in `go.mod`'s `require` block (stdlib-only constraint). `modernc.org/sqlite`
  is already present for the in-memory repo tests.
