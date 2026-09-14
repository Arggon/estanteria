# ADR 0004 — Render estático compartido: `serve` y `render` sobre un mismo generador

- **Status:** Accepted (2026-09-14)
- **Context:** la vista web (issue #3) y el deploy a Pages (issue #4) necesitan el mismo HTML. Exploración en `docs/explorations/exploration-web-rendering-001.md` (2026-09-14): servidor en vivo vs snapshot estático vs SPA.
- **Decision:** un único generador de HTML (función pura sobre `Shelf`, testeable con httptest) con dos entradas: `estanteria serve` lo sirve por request con `net/http` bindeando a 127.0.0.1, y `estanteria render --out <file>` lo escribe como snapshot autocontenido para GitHub Pages. Sin API JSON ni frontend con build step.
- **Consequences:** el deploy de Pages no tiene runtime (solo archivos); `serve` siempre muestra el estado actual del ledger sin regenerar a mano; agregar cualquier vista futura (stats detalladas, filtros) se hace una sola vez en el generador. Trade-off: el snapshot puede quedar desactualizado respecto del ledger local — aceptable porque es un mirror de solo lectura.
