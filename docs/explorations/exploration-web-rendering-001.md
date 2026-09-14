---
exploration_id: web-rendering-001
title: web-rendering
status: decided
created: 2026-09-14
---

# Exploration: web-rendering (web-rendering-001)

Spike record: how should the web view of the estantería be produced and served?
Decision recorded in [ADR 0004](../adr/0004-static-render-shared-generator.md).

## Candidates

1. **Servidor en vivo (`estanteria serve`)** — `net/http` lee el ledger en cada
   request y renderiza HTML on the fly.
2. **Snapshot estático (`estanteria render`)** — un generador produce un
   `index.html` autocontenido desde el ledger; publicable en GitHub Pages sin
   backend.
3. **SPA con API JSON** — `serve` expone `/api/books` y un frontend JS renderiza
   en el cliente.

## Criteria

- Solo lectura garantizada (ADR 0003): nada de caminos de escritura web.
- Cero dependencias: stdlib únicamente; sin build step de frontend.
- El mismo HTML debe poder publicarse en Pages (issue #4) y verse en local.
- Mantenimiento: un solo generador, dos entradas (serve/render).

## Findings

- El snapshot estático (2) cumple los tres criterios de forma trivial: HTML
  autocontenido, sin runtime en el servidor de Pages (fuente:
  [GitHub Pages docs](https://docs.github.com/en/pages), consultado 2026-09-14).
- El servidor en vivo (1) es útil para *mirar el estado actual* sin regenerar —
  pero reutilizando el MISMO generador por request, el modo serve es solo un
  `http.HandlerFunc` alrededor de `render` (fuente:
  [net/http docs](https://pkg.go.dev/net/http), consultada 2026-09-14).
- La SPA (3) viola "sin build step de frontend" y agrega superficie (API JSON +
  JS embebido) para un tracker de un solo usuario; descartada.

## Recommendation

Un solo generador de HTML (función pura sobre `Shelf`, testeable con httptest) con
dos entradas: `estanteria serve` (re-genera por request, 127.0.0.1) y
`estanteria render --out` (snapshot para Pages). Formalizado en ADR 0004.

## Decision

[ADR 0004 — Render estático compartido: `serve` y `render` sobre un mismo generador](../adr/0004-static-render-shared-generator.md) — Accepted 2026-09-14.
