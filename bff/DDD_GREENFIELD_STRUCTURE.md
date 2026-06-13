# Greenfield BFF Structure — DDD with Onion Architecture

## Context

Reference blueprint for how the `bff` *would* be organized if started from
scratch today, carrying the same features but designed DDD-first. Chosen shape:
**a single bounded context, organized by onion layer** (not split into
context modules). The point of this document is to define the layers, what lives
in each, the dependency rule, and the few DDD decisions that differ from the
current code. Module path: `github.com/kjj1998/kinji/bff`.

Current feature set this must cover:
- Import a bank-statement PDF (decrypt/validate, Claude extraction, running-balance guard, SSE progress)
- Save reviewed transactions
- Get a user's monthly transactions
- Get periods (years/months that have data)
- Generate a monthly insights summary (totals, top merchants/categories, trends, biggest category changes, recent txns)

## The dependency rule (the whole point)

```
domain  ←  app  ←  adapter        cmd/api wires everything
   ▲         ▲         │
   └─────────┴─────────┘  (imports always point inward)
```

- `domain` imports **nothing** under `internal/` and no I/O libs (no SQL, no HTTP, no Anthropic SDK, no `json` tags).
- `app` imports **only** `domain`. It owns the ports (interfaces) it needs.
- `adapter/*` import `app` + `domain` and implement the ports / drive the app.
- `cmd/api` (composition root) imports everything and wires concrete adapters.

## Target layout

```
bff/
  cmd/
    api/
      main.go                 # composition root: load config, build adapters, start server
  internal/
    domain/                   # innermost ring — pure model, no internal imports, no I/O
      transaction.go            # Transaction entity; Category, Direction value objects (IsValid, IsOutflow/IsInflow)
      money.go                  # Money value object (cents) — optional; see "DDD decisions"
      period.go                 # Period value object
      statement.go              # Statement / TransactionBatch + running-balance invariant (see decisions)
      spending.go               # ValueAndChange[T], MerchantSpending, CategorySpending,
                                #   CategorySpendingChange, DaySpending, MonthSpending (value objects)
      summary.go                # MonthlySummary read-model + SummaryCalculator domain service (pure calcs)
      errors.go                 # domain sentinel errors (ErrInvalidCategory, ErrInvalidDirection, ErrBalanceMismatch)
    app/                      # use cases + ports. imports domain only
      ports.go                  # TransactionRepository, StatementParser (interfaces owned here)
      transaction_service.go    # import / save / get-monthly / get-periods use cases
      summary_service.go        # orchestrates repo reads -> domain.SummaryCalculator -> MonthlySummary
      errors.go                 # ClientError (bad input -> 4xx)
    adapter/                  # outer ring
      persistence/
        sqlite/                 # TransactionRepository impl: schema, queries, date-range helpers
      parser/
        claude/                 # StatementParser impl: Anthropic client, prompt, tool schema (raw extraction only)
      http/
        server.go               # router assembly + middleware chain
        handler/                # transaction.go, summary.go, health.go (call app, map domain<->dto)
        dto/                    # request/response structs (json tags) + mappers + label formatting
        middleware/             # logging, recovery, cors
    config/                   # env loading (standalone, imported by main)
```

## DDD decisions that differ from the current code

1. **Ports live in `app`, not in adapters.** The consumer (`app`) declares
   `TransactionRepository` and `StatementParser`; `sqlite` and `claude`
   implement them. Today these interfaces sit in the adapter packages.

2. **Domain is serialization-free.** No `json` tags on `Transaction`/value
   objects. The wire format is a presentation concern: `adapter/http/dto` owns
   request/response structs + explicit mappers at the boundary.

3. **The running-balance check is a domain invariant, not adapter logic.**
   Today the balance guard lives inside the Claude adapter. Greenfield: the
   adapter only *extracts* raw rows; a domain `Statement` (or `TransactionBatch`)
   validates the running-balance invariant and returns `ErrBalanceMismatch`.
   This keeps a real business rule in the core and makes the parser a dumb
   extraction adapter.

4. **`MonthlySummary` is a read-model/projection; the math is a domain service.**
   The pure calculations (top outflow, category-change deltas, savings rate,
   trends) live in a `SummaryCalculator` in `domain`. `app.summary_service`
   only orchestrates repo reads and hands data to the calculator.

5. **Presentation labels and view windows stay out of `domain`.** Domain returns
   data keyed by `time.Weekday` / `time.Month` with raw amounts; the `dto`
   mapper renders `"Mon"`/`"Jan"` and applies view policy (top-3, recent-5,
   6-month window). This is the one spot the current code leaks presentation
   inward — greenfield fixes it.

6. **`Money` value object (optional).** Textbook DDD would model amounts as a
   `Money` value object (cents + helpers like `Dollars()`, `Add`) instead of a
   bare `int`. Worth it if money math spreads; skip and keep `int` cents if you
   want to stay lean. Listed so the choice is explicit.

7. **Aggregate roots are intentionally small.** `Transaction` is the root for
   the ledger; `Statement` is a short-lived aggregate during import that
   enforces the balance invariant. No premature aggregate design beyond these.

## Scaffolding order (if building it out)

1. `internal/domain` — entities, value objects, `SummaryCalculator`, `Statement` invariant, errors. Unit-test in isolation (no mocks needed; pure).
2. `internal/app` — ports + use-case services against the domain. Test with hand-written port mocks colocated in `app`.
3. `internal/config` — env loader.
4. `internal/adapter/persistence/sqlite` — implement `TransactionRepository`, return domain types.
5. `internal/adapter/parser/claude` — implement `StatementParser`, raw extraction only.
6. `internal/adapter/http` — handlers, dto + mappers, middleware, server assembly.
7. `cmd/api/main.go` — wire config + adapters + server.

## Future scaling note (not now)

If features grow, the natural seam is three bounded contexts — **Importing**
(PDF→transactions), **Ledger** (storage/retrieval), **Insights** (analytics) —
each promoted to its own `internal/<context>/{domain,app,adapter}`. The
layer-first layout above keeps that move cheap because the dependency rule is
already enforced; you'd split along the existing service boundaries.

## Verification (for a real build-out)

- `go build ./...`, `go vet ./...`, `go test ./...`.
- **Dependency-rule check:** confirm `internal/domain` imports no other
  `internal/*`, and `internal/app` imports only `internal/domain` (grep imports
  or `go list -deps`).
- **E2E smoke:** run `go run ./cmd/api`, then exercise
  `GET /api/v1/transactions/{id}`, `GET /api/v1/summary/{id}`,
  `POST /api/v1/transactions/{id}`, and `POST /.../import` (SSE).