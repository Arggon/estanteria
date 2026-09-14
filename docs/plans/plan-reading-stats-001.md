---
plan_id: reading-stats-001
title: Plan for Stats de lectura
spec: docs/specs/spec-reading-stats-001.md
status: implemented
created: 2026-09-14
---

# Plan: Stats de lectura (reading-stats-001)

Derived from `docs/specs/spec-reading-stats-001.md`. Each task carries a
verifiable acceptance criterion and links back to the spec.

## Tasks

### T1: stats.go — ComputeStats puro

- `Stats` struct: `Total int`, `ByStatus map[Status]int`, `AvgRating float64`, `PagesPerMonth map[string]int`.
- Ventana de 12 meses: `[now.AddDate(0,-11,0) truncado a mes, mes de now]`, clave `2006-01`.
- **Acceptance:** tests table-driven en stats_test.go (con/sin ratings, sin pages/finished, distribución entre meses, borde de los 12 meses) en verde.

### T2: main.go — comando `stats`

- Case `stats` en `run()` + `cmdStats`: `LoadShelf`, `ComputeStats`, impresión de totales, promedio y meses. Sin `SaveShelf`.
- **Acceptance:** e2e en main_test-style (ledger en `t.TempDir()`, `run([]string{"stats"})`) imprime totales y meses; ledger sin cambios.

### T3: docs + gates

- README sección Uso: `estanteria stats`.
- `go test ./... && go vet ./...`, `arggon validate --json`, `arggon spec validate --json` en verde; flip a `implemented`.
- **Acceptance:** todos los gates ok y specs/plans en `implemented` antes del commit.
