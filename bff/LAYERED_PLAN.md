# Layered Structure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restructure the `bff` from the onion/adapter layout into a flat, conventional layered layout — one package per layer (`model`, `service`, `handler`, `store`, `parser`, `config`) — for the simplest possible mental model. No behavior or wire-format changes.

**Architecture:** Five flat packages under `internal/`, plus `cmd/api`. Data flows `handler → service → store`/`parser`; `model` holds all shared domain types + business logic. Interfaces stay where they are consumed (`service`). This deliberately trades away the onion's serialization-free-domain and adapter-ring separation for fewer, flatter packages — an accepted trade for a small single-service BFF.

**Tech Stack:** Go 1.25, stdlib `net/http`, sqlite, Anthropic (claude) parser.

> Supersedes `FLATTEN_PLAN.md` (which kept the onion). Pick one; do not run both.

---

## Package mapping (what moves where)

| Current | Target | Package rename |
|---|---|---|
| `internal/domain/` | `internal/model/` | `domain` → `model` |
| `internal/app/` | `internal/service/` | `app` → `service` |
| `internal/adapter/persistence/sqlite/` | `internal/store/` | `sqlite` → `store` |
| `internal/adapter/parser/claude/` | `internal/parser/` | `claude` → `parser` |
| `internal/adapter/http/handler/` | `internal/handler/` | (stays `handler`) |
| `internal/adapter/http/dto/` | `internal/handler/*_dto.go` | `dto` → `handler` (merged) |
| `internal/adapter/http/middleware/` | `internal/handler/middleware.go` | `middleware` → `handler` (merged) |
| `internal/adapter/http/server/` | `internal/handler/server.go` | `server` → `handler` (merged) |
| `internal/config/` | `internal/config/` | unchanged |

## Conventions for this plan (macOS / BSD)

- Run all commands from the `bff/` module root.
- **`sed`** uses BSD syntax (`sed -i ''`). Used only for path strings, package clauses, and whole-line deletes — none need word boundaries.
- **`perl`** is used for identifier (selector) rewrites because BSD `sed` lacks `\b`. All domain/app/store/parser identifiers used across packages are **exported (uppercase first letter)**, so the pattern `\bX\.([A-Z])` rewrites real qualifiers (`domain.Transaction`) without touching prose like "…the domain. The…".
- Every task ends with `go build ./... && go vet ./... && go test ./...` green, then a commit. The existing tests in `internal/domain` and `internal/adapter/http/dto` are the regression net.

---

### Task 0: Capture green baseline

**Files:** none (verification only).

- [ ] **Step 1: Confirm clean before touching anything**

Run: `cd bff && go build ./... && go vet ./... && go test ./...`
Expected: build/vet clean; tests `ok` for `internal/domain` and `internal/adapter/http/dto`.

---

### Task 1: Rename `domain` → `model`

**Files:** moves `internal/domain/` → `internal/model/`; rewrites every `domain.X` selector and `internal/domain` import path across the repo.

- [ ] **Step 1: Move the directory (keeps history)**

Run:
```bash
git mv internal/domain internal/model
```

- [ ] **Step 2: Rename the package clause in the moved files**

Run:
```bash
sed -i '' 's/^package domain$/package model/' internal/model/*.go
```

- [ ] **Step 3: Rewrite import paths repo-wide**

Run:
```bash
grep -rl "internal/domain" --include="*.go" . | xargs sed -i '' 's#internal/domain#internal/model#g'
```

- [ ] **Step 4: Rewrite `domain.` selectors repo-wide**

Run:
```bash
grep -rl "domain\." --include="*.go" . | xargs perl -i -pe 's/\bdomain\.([A-Z])/model.$1/g'
```

- [ ] **Step 5: Format, build, test**

Run: `gofmt -w ./... 2>/dev/null; go build ./... && go vet ./... && go test ./...`
Expected: clean; `internal/model` tests `ok`.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: rename domain package to model

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 2: Rename `app` → `service` (and clean the handler field name)

**Files:** moves `internal/app/` → `internal/service/`; rewrites `app.X` selectors + import path; renames the `service` struct field in two handler files to `svc` so it doesn't read as `service service.SummaryService`.

- [ ] **Step 1: Move the directory**

Run:
```bash
git mv internal/app internal/service
```

- [ ] **Step 2: Rename the package clause**

Run:
```bash
sed -i '' 's/^package app$/package service/' internal/service/*.go
```

- [ ] **Step 3: Rewrite import paths repo-wide**

Run:
```bash
grep -rl "internal/app" --include="*.go" . | xargs sed -i '' 's#internal/app#internal/service#g'
```

- [ ] **Step 4: Rewrite `app.` selectors repo-wide**

Run:
```bash
grep -rl "app\." --include="*.go" . | xargs perl -i -pe 's/\bapp\.([A-Z])/service.$1/g'
```

After this, the two handler structs read `service service.SummaryService` / `service service.TransactionService`. The next step fixes that.

- [ ] **Step 5: Rename the handler `service` field to `svc`**

In `internal/adapter/http/handler/summary.go`, apply these exact replacements:
```go
type SummaryHandler struct {
	service service.SummaryService
}
```
→
```go
type SummaryHandler struct {
	svc service.SummaryService
}
```
and
```go
	return &SummaryHandler{service: svc}
```
→
```go
	return &SummaryHandler{svc: svc}
```
and
```go
	summary, err := h.service.GenerateMonthlySummary(r.Context(), id, month, year)
```
→
```go
	summary, err := h.svc.GenerateMonthlySummary(r.Context(), id, month, year)
```

In `internal/adapter/http/handler/transaction.go`, apply:
```go
type TransactionHandler struct {
	service service.TransactionService
}
```
→
```go
type TransactionHandler struct {
	svc service.TransactionService
}
```
and
```go
	return &TransactionHandler{service: svc}
```
→
```go
	return &TransactionHandler{svc: svc}
```
Then replace every `h.service.` with `h.svc.` in that file (4 call sites: `GetMonthlyTransactions`, `ImportStatement`, `SaveTransactions`, `GetPeriods`):
```bash
perl -i -pe 's/\bh\.service\./h.svc./g' internal/adapter/http/handler/transaction.go internal/adapter/http/handler/summary.go
```

- [ ] **Step 6: Build, test**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: clean.

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "refactor: rename app package to service

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 3: Rename `adapter/persistence/sqlite` → `store`

**Files:** moves the dir up to `internal/store/`; rewrites the import path and the `sqlite.X` selectors in `cmd/api/main.go`.

- [ ] **Step 1: Move the directory and clean empty parent**

Run:
```bash
git mv internal/adapter/persistence/sqlite internal/store
rmdir internal/adapter/persistence
```

- [ ] **Step 2: Rename the package clause**

Run:
```bash
sed -i '' 's/^package sqlite$/package store/' internal/store/*.go
```

- [ ] **Step 3: Rewrite import path repo-wide**

Run:
```bash
grep -rl "internal/adapter/persistence/sqlite" --include="*.go" . | xargs sed -i '' 's#internal/adapter/persistence/sqlite#internal/store#g'
```

- [ ] **Step 4: Rewrite `sqlite.` selectors (call sites in main.go)**

Run:
```bash
grep -rl "sqlite\." --include="*.go" . | xargs perl -i -pe 's/\bsqlite\.([A-Z])/store.$1/g'
```
(This rewrites `sqlite.NewClient`/`sqlite.NewRepository`; it does not touch the blank-imported `mattn/go-sqlite3` driver, which has no `sqlite.` selector.)

- [ ] **Step 5: Build, test**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: clean.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: rename sqlite adapter to store

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 4: Rename `adapter/parser/claude` → `parser`

**Files:** moves the dir to `internal/parser/`; rewrites import path and `claude.X` selectors in `cmd/api/main.go`.

- [ ] **Step 1: Move the directory and clean empty parents**

Run:
```bash
git mv internal/adapter/parser/claude internal/parser
rmdir internal/adapter/parser
rmdir internal/adapter
```

- [ ] **Step 2: Rename the package clause**

Run:
```bash
sed -i '' 's/^package claude$/package parser/' internal/parser/*.go
```

- [ ] **Step 3: Rewrite import path repo-wide**

Run:
```bash
grep -rl "internal/adapter/parser/claude" --include="*.go" . | xargs sed -i '' 's#internal/adapter/parser/claude#internal/parser#g'
```

- [ ] **Step 4: Rewrite `claude.` selectors**

Run:
```bash
grep -rl "claude\." --include="*.go" . | xargs perl -i -pe 's/\bclaude\.([A-Z])/parser.$1/g'
```

- [ ] **Step 5: Build, test**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: clean. `internal/adapter/` no longer exists except `internal/adapter/http` (handled next).

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "refactor: rename claude parser adapter to parser

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 5: Merge `adapter/http/{handler,dto,middleware,server}` → `internal/handler`

**Files:** collapses four packages into one `handler` package. The `dto` types/funcs lose their qualifier (`dto.ToTransactions` → `ToTransactions`); `server.go` loses its `handler.`/`middleware.` qualifiers.

> Verified safe: no identifier collides across the four packages. `dto` contributes type `Transaction`/`Period`/etc. and `To*` funcs; `handler` has `*Handler` types and `write*`/`require*`/`parse*` helpers; `middleware` has `Logging`/`Recovery`/`CORS`/`Chain`; `server` has `New`. All distinct.

- [ ] **Step 1: Move handler dir up, pull the other three in (renaming dto files)**

Run:
```bash
git mv internal/adapter/http/handler internal/handler
git mv internal/adapter/http/dto/period.go        internal/handler/period_dto.go
git mv internal/adapter/http/dto/summary.go       internal/handler/summary_dto.go
git mv internal/adapter/http/dto/transaction.go   internal/handler/transaction_dto.go
git mv internal/adapter/http/dto/period_test.go      internal/handler/period_dto_test.go
git mv internal/adapter/http/dto/summary_test.go     internal/handler/summary_dto_test.go
git mv internal/adapter/http/dto/transaction_test.go internal/handler/transaction_dto_test.go
git mv internal/adapter/http/middleware/middleware.go internal/handler/middleware.go
git mv internal/adapter/http/server/server.go         internal/handler/server.go
rmdir internal/adapter/http/dto internal/adapter/http/middleware internal/adapter/http/server internal/adapter/http internal/adapter
```

- [ ] **Step 2: Set every moved file's package to `handler`**

Run:
```bash
sed -i '' 's/^package dto$/package handler/' internal/handler/period_dto.go internal/handler/summary_dto.go internal/handler/transaction_dto.go internal/handler/period_dto_test.go internal/handler/summary_dto_test.go internal/handler/transaction_dto_test.go
sed -i '' 's/^package middleware$/package handler/' internal/handler/middleware.go
sed -i '' 's/^package server$/package handler/' internal/handler/server.go
```

- [ ] **Step 3: Drop the now-internal qualifiers**

Run:
```bash
# dto. -> nothing, in the two real handlers
perl -i -pe 's/\bdto\.//g' internal/handler/summary.go internal/handler/transaction.go
# handler. and middleware. -> nothing, in server.go
perl -i -pe 's/\bhandler\.//g; s/\bmiddleware\.//g' internal/handler/server.go
```

- [ ] **Step 4: Remove the now-unused imports**

Run:
```bash
sed -i '' '\#internal/adapter/http/dto#d' internal/handler/summary.go internal/handler/transaction.go
sed -i '' '\#internal/adapter/http/handler#d' internal/handler/server.go
sed -i '' '\#internal/adapter/http/middleware#d' internal/handler/server.go
```

- [ ] **Step 5: Point `main.go` at the merged package**

Run:
```bash
sed -i '' 's#internal/adapter/http/server#internal/handler#g' cmd/api/main.go
perl -i -pe 's/\bserver\.([A-Z])/handler.$1/g' cmd/api/main.go
```
This changes the import to `internal/handler` and `server.New(...)` → `handler.New(...)`.

- [ ] **Step 6: Format, build, test**

Run: `gofmt -w internal/handler cmd/api 2>/dev/null; go build ./... && go vet ./... && go test ./...`
Expected: clean. The dto tests (now `period_dto_test.go` etc., package `handler`) still pass; they reference `maxBiggestChanges`/`maxRecentTransactions`, which are now in the same `handler` package.

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "refactor: merge http dto/handler/middleware/server into internal/handler

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 6: Update the architecture docs

**Files:** `DDD_GREENFIELD_STRUCTURE.md`, `DDD_REFACTOR_PLAN.md`, `FLATTEN_PLAN.md`.

- [ ] **Step 1: Replace each "Target layout" block with the layered tree**

Use this layout in both `DDD_*.md` docs:
```
cmd/api/main.go
internal/
  model/      # all domain types + business logic (no json tags kept here is no longer enforced)
  service/    # use cases + repo/parser interfaces
  handler/    # http: handlers, *_dto.go mappers, middleware, server
  store/      # sqlite
  parser/     # claude
  config/
```

- [ ] **Step 2: Add a deviation note to `DDD_GREENFIELD_STRUCTURE.md`**

Append:
```markdown
## Update (layered restructure)

The onion/adapter layout was flattened to a conventional layered layout
(`model` / `service` / `handler` / `store` / `parser` / `config`) for a simpler
mental model. This intentionally relaxes two earlier DDD decisions: the
domain↔wire separation (dto mappers now live in `handler` beside the domain
types in `model`) and the adapter ring. The `service` package still owns its
repo/parser interfaces. See LAYERED_PLAN.md.
```

- [ ] **Step 3: Mark FLATTEN_PLAN.md superseded**

Add at the top of `FLATTEN_PLAN.md`:
```markdown
> SUPERSEDED by LAYERED_PLAN.md — the project chose the flat layered layout instead.
```

- [ ] **Step 4: Commit**

```bash
git add DDD_GREENFIELD_STRUCTURE.md DDD_REFACTOR_PLAN.md FLATTEN_PLAN.md
git commit -m "docs: update layout docs for layered restructure

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>"
```

---

### Task 7: Final verification

**Files:** none (verification only).

- [ ] **Step 1: Full build/vet/test**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: all clean.

- [ ] **Step 2: Confirm the final package set**

Run: `ls internal/`
Expected exactly: `config  handler  model  parser  service  store`.

- [ ] **Step 3: Confirm the old trees are gone**

Run: `test ! -d internal/adapter && test ! -d internal/domain && test ! -d internal/app && echo "old dirs removed: OK"`
Expected: `old dirs removed: OK`.

- [ ] **Step 4: Sanity-check dependency direction (model is still the leaf)**

Run:
```bash
go list -deps ./internal/model | grep "kinji/bff/internal" || echo "model imports no internal packages: OK"
```
Expected: `model imports no internal packages: OK`. (Note: `service` now legitimately may be imported by `handler`/`store`/`parser`; the only hard rule kept is that `model` stays a leaf.)

- [ ] **Step 5: End-to-end smoke (manual)**

Run `go run ./cmd/api`, then exercise:
- `GET /api/v1/transactions/{id}`
- `GET /api/v1/summary/{id}` — JSON must still emit `"Mon"`/`"Jan"` labels and top-3/recent-5 truncation
- `POST /api/v1/transactions/{id}`
- `POST /api/v1/transactions/{id}/import` (SSE progress)

Expected: identical responses to pre-refactor.

---

## Self-review notes

- **Coverage:** every current package has a target (mapping table). `config` is untouched.
- **Type consistency:** package identifiers `model` / `service` / `store` / `parser` / `handler` are used identically across selectors, imports, and `main.go` wiring. The handler field rename (`service` → `svc`) avoids `service service.X`.
- **Collisions:** the four-package merge (Task 5) was verified collision-free against the actual identifiers in `summary.go`, `transaction.go`, `health.go`, `http_utils.go`, `mock_service.go`, `middleware.go`, `server.go`, and the three dto files.
- **macOS gotcha:** identifier rewrites use `perl` (BSD `sed` has no `\b`); the `([A-Z])` anchor avoids rewriting prose comments.
- **Trade-off acknowledged:** this drops the domain-purity / ports-ring rigor from the DDD docs on purpose, in exchange for a flatter, simpler structure. Tests are unchanged and remain the safety net.
- **No new tests:** pure restructure; no behavior added.
```
