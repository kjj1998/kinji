# Feature-Based Folder Structure Refactor — Design

**Date:** 2026-06-13
**Status:** Approved (pending spec review)
**Scope:** `bff/` Go service only. Pure structural refactor — no behavior changes.

## Goal

Reorganize the `bff` service from its current **layered/hexagonal** structure (grouped by
technical layer: `model`, `service`, `handler`, `store`, `parser`) into a **feature-based**
structure (grouped by feature: `transaction`, `summary`) with a shared kernel.

The module path stays `github.com/kjj1998/kinji/bff`. No runtime behavior changes; all
existing tests pass unchanged except for package-name and import-path updates.

## Decisions (from brainstorming)

1. **Shared handling:** Feature + shared kernel. Each feature owns its handler/service/
   repository/domain; cross-cutting code lives in shared packages.
2. **Persistence:** Split per feature. Each feature owns its repository implementation and a
   segregated repository interface, both backed by a shared `*sql.DB` handle.
3. **Feature boundaries:** Two features — `transaction` (monthly transactions, periods,
   statement import, and the PDF parser) and `summary` (monthly spending analysis). The PDF
   parser becomes an internal detail of the `transaction` feature.
4. **`Transaction` type:** Lives in the shared kernel (used by both features), so features are
   true siblings with no feature-to-feature dependency.
5. **Within-feature layout:** Approach B — nested sub-packages per layer inside each feature
   (`handler/`, `service/`, `domain/`, `store/`, `parser/`). Preserves existing layer
   separation while grouping by feature.
6. **Naming:** `internal/shared` (cross-cutting domain types) and `internal/platform`
   (feature-agnostic infrastructure). Composition root in its own `internal/server` package.
7. **`ValueAndChange`** lives in `summary/domain` (only summary uses it), not in shared.

## Target structure

```
bff/
├── cmd/api/main.go                  # entrypoint (unchanged role)
└── internal/
    ├── shared/                      # cross-cutting DOMAIN types (package shared)
    │   ├── transaction.go           #   Transaction, Category, Direction, parsers, IsValid
    │   ├── money.go                 #   Money
    │   ├── month.go                 #   Month, ParseMonth, GetMonthRangeDateStrings
    │   ├── errors.go                #   ErrInvalidCategory, ErrInvalidDirection
    │   └── client_error.go          #   ClientError (app → HTTP 4xx)
    ├── platform/                    # feature-agnostic INFRA
    │   ├── database/                #   NewClient, schema (sqlite setup)
    │   ├── http/                    #   middleware, http_utils, health
    │   └── config/                  #   config.go
    ├── server/                      # composition root
    │   └── server.go                #   New(): wires features + platform, builds router
    ├── transaction/
    │   ├── domain/                  #   Statement, StatementLine, Period, statement/PDF errors
    │   ├── service/                 #   TransactionService, repo + parser interfaces
    │   ├── store/                   #   sqlite repo impl (transaction queries)
    │   ├── parser/                  #   PDF parser (+ schema)
    │   └── handler/                 #   handlers + DTOs
    └── summary/
        ├── domain/                  #   MonthlySummary, SummaryCalculator, spending types, ValueAndChange
        ├── service/                 #   SummaryService, repo interface
        ├── store/                   #   sqlite repo impl (summary queries)
        └── handler/                 #   handlers + DTOs
```

## File-by-file mapping

### `internal/model/` splits
| Current | Destination |
|---|---|
| `transaction.go` (Transaction, Category, Direction, parsers, IsValid/IsInflow/IsOutflow) | `shared/transaction.go` |
| `money.go` (Money) | `shared/money.go` |
| `month.go` (Month, ParseMonth) | `shared/month.go` |
| `errors.go` → `ErrInvalidCategory`, `ErrInvalidDirection` | `shared/errors.go` |
| `errors.go` → `ErrBalanceMismatch`, `ErrPDFPasswordRequired`, `ErrPDFWrongPassword`, `ErrPDFCorrupt` | `transaction/domain/errors.go` |
| `statement.go` (Statement, StatementLine) | `transaction/domain/statement.go` |
| `period.go` (Period) | `transaction/domain/period.go` |
| `summary.go` (MonthlySummary, SummaryInput, SummaryCalculator, math helpers) | `summary/domain/summary.go` |
| `spending.go` (MerchantSpending, CategorySpending, CategorySpendingChange, DaySpending, MonthSpending, ValueAndChange) | `summary/domain/spending.go` |
| `tests/statement_test.go` | `transaction/domain/` test |
| `tests/summary_test.go` | `summary/domain/` test |

### `internal/service/` splits
| Current | Destination |
|---|---|
| `ports.go` → `TransactionRepository` (transaction methods) + `StatementParser` | `transaction/service/` (segregated interfaces) |
| `ports.go` → repository methods used by summary | `summary/service/` (segregated interface) |
| `transaction_service.go` | `transaction/service/service.go` |
| `summary_service.go` | `summary/service/service.go` |
| `errors.go` (ClientError) | `shared/client_error.go` |
| `mock_parser.go` | `transaction/service/` (test mock) |
| `mock_repository.go` | split into `transaction/service/` and `summary/service/` test mocks |

### `internal/store/` splits
| Current | Destination |
|---|---|
| `sqlite.go` → `NewClient` + `schema` application | `platform/database/` |
| `sqlite.go` → `Repository` transaction methods | `transaction/store/repository.go` |
| `sqlite.go` → `Repository` summary methods | `summary/store/repository.go` |
| `queries.go` → `schema` const | `platform/database/` |
| `queries.go` → transaction SQL consts | `transaction/store/` |
| `queries.go` → summary SQL consts | `summary/store/` |
| `utils.go` → `GetMonthRangeDateStrings` | `shared/month.go` |
| `utils.go` → `currentAndPreviousMonth` | `summary/store/` |

### `internal/handler/` splits
| Current | Destination |
|---|---|
| `middleware.go` | `platform/http/` |
| `http_utils.go` | `platform/http/` |
| `health.go` | `platform/http/` |
| `server.go` (New + routing) | `internal/server/server.go` |
| `transaction.go`, `transaction_dto.go`, `period_dto.go` (+ tests) | `transaction/handler/` |
| `summary.go`, `summary_dto.go` (+ tests) | `summary/handler/` |
| `mock_service.go` | split per feature alongside handler tests |

### Other
| Current | Destination |
|---|---|
| `internal/parser/parser.go`, `internal/parser/schema.go` | `transaction/parser/` |
| `internal/config/config.go` | `platform/config/` |
| `cmd/api/main.go` | unchanged role; update imports + call `server.New(...)` |

## Dependency rules (must remain acyclic)

```
cmd → server → {transaction/*, summary/*, platform/*, shared}
transaction/* → shared, platform        (handler → service → {domain, store, parser})
summary/*     → shared, platform        (handler → service → {domain, store})
shared        → (no internal deps; leaf)
platform      → shared (only; no feature imports; otherwise leaf)
```

- **Features never import each other.** `Transaction` lives in `shared`, making features siblings.
- **`platform` imports no feature.** Generic HTTP infra (middleware, utils, health) and DB setup only.
- **`server` is the only composition root** — the single package allowed to import both features.
- Each feature's repository interface is **owned by that feature's `service` package** (interface
  segregation); the implementation lives in the feature's `store` package and receives the shared
  `*sql.DB` produced by `platform/database`.

## Composition root

`handler.New(repo, parser, corsOrigin)` becomes `server.New(...)` in `internal/server`. It:
- receives the shared `*sql.DB` (or constructs feature repositories from it),
- constructs each feature's repository, parser, service, and handler,
- registers routes and applies the middleware chain from `platform/http`.

This keeps `platform` a pure leaf (no feature imports) while giving the wiring one clear home.
`cmd/api/main.go` opens the DB via `platform/database`, loads `platform/config`, and calls
`server.New(...)`.

## Non-goals

- No behavior changes, no API/route changes, no SQL changes.
- No new features, no renamed exported identifiers beyond what package moves require.
- No change to the module path or to `cmd/api` responsibilities.

## Verification

- This is a pure move/rename refactor. Module path unchanged.
- Incremental: after each feature/package is moved, `go build ./...` and `go test ./...` stay green.
- Final gate: `go vet ./...` and full `go test ./...` — all existing tests pass (only package
  names and import paths change).

## Suggested execution order

1. Create `shared/` (move Transaction, Money, Month, shared errors, ClientError, month-range helper). Build.
2. Create `platform/` (`database`, `http`, `config`). Build.
3. Move `transaction` feature (domain → service → store → parser → handler). Build + test.
4. Move `summary` feature (domain → service → store → handler). Build + test.
5. Create `internal/server` composition root; update `cmd/api/main.go`. Build + test.
6. Delete now-empty `internal/model`, `internal/service`, `internal/handler`, `internal/store`,
   `internal/parser`, `internal/config`. Final `go vet ./...` + `go test ./...`.
