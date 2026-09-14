---
plan_id: web-view-001
title: Plan for Vista web de solo lectura
spec: docs/specs/spec-web-view-001.md
status: implemented
created: 2026-09-14
---

# Plan: Vista web de solo lectura (web-view-001)

Derived from `docs/specs/spec-web-view-001.md`. Each task carries a
verifiable acceptance criterion and links back to the spec.

## Tasks

### T1: Generador HTML puro (`web.go`)

- `RenderHTML(books []Book, generatedAt time.Time) []byte`: HTML autocontenido
  (CSS inline, utf-8), secciones por status con display acentuado, footer con
  total y promedio de ratings; todo escapado con `html.EscapeString`.
- **Acceptance:** test de RenderHTML verifica escape de `<script>` y presencia
  de las tres secciones.

### T2: Handler HTTP de solo lectura (`web.go`)

- `ServeHandler(path string) http.Handler`: carga el ledger por request,
  sirve `RenderHTML` con `text/html; charset=utf-8`; métodos != GET → 405.
- **Acceptance:** tests httptest: GET 200 + Content-Type, POST 405, ledger
  vacío → 200 con "(vacía)".

### T3: Comando `serve` (`main.go`)

- Flag `--addr` default `127.0.0.1:8080`; imprime la URL y llama
  `http.ListenAndServe`. Sin comando `render` (otra story).
- **Acceptance:** `go vet ./...` en verde y la URL queda registrada en usage.

### T4: Docs y gates

- README: sección Uso con `estanteria serve` (solo lectura, 127.0.0.1).
- `go test ./... && go vet ./...`, `arggon validate --json`,
  `arggon spec validate --json` en verde; spec/plan → implemented.
- **Acceptance:** todos los gates pasan en el worktree.
