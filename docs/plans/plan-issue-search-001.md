---
plan_id: issue-search-001
title: Plan for Búsqueda por título/autor
spec: docs/specs/spec-issue-search-001.md
status: implemented
created: 2026-09-14
---

# Plan: Búsqueda por título/autor (issue-search-001)

Derived from `docs/specs/spec-issue-search-001.md`. Each task carries a
verifiable acceptance criterion and links back to the spec.

## Tasks

### T1: Dominio — `Shelf.Search`

- En `shelf.go`: `Search(query string) []Book`; fold con `foldAccents` +
  `ToLower`, contains sobre título o autor, ordenado por título, vacío/nil
  para query en blanco.
- **Acceptance:** tests table-driven en `search_test.go` (match título, match
  autor, acentos en query y datos, sin match, query vacía, orden por título) en verde.

### T2: CLI — comando `search`

- En `main.go`: `cmdSearch` con la misma fila de `list` (reusar `padRight`/
  `trunc`) + `total: N`; sin flags; nunca llama a `SaveShelf`. Actualizar `usage`.
- **Acceptance:** e2e en `main_test.go`-style: `run([]string{"search", ...})`
  imprime filas y total, exit 0, y el ledger en `t.TempDir()` queda byte a byte igual.

### T3: Gates y cierre

- `go test ./... && go vet ./...` en verde; `arggon validate --json` ok;
  `arggon spec validate --json` ok; flip de status a `implemented` en spec y plan.
- **Acceptance:** todos los gates verdes y frontmatter en `implemented`.
