---
spec_id: web-view-001
title: Vista web de solo lectura
status: implemented
created: 2026-09-14
---

# Spec: Vista web de solo lectura (web-view-001)

## Purpose

Ver la estantería en el navegador sin escribir estado (ADR 0003), reusando un
solo generador HTML que la story de `render`/Pages (issue #4) también usará
(ADR 0004).

Invariantes:

- **Pure read** — ningún endpoint muta el ledger; la única escritura del
  ledger sigue siendo `SaveShelf` desde el CLI.
- **Bind 127.0.0.1 por defecto** — el servidor no se expone a la red salvo
  que el usuario pase `--addr` explícitamente.
- **Un solo generador HTML compartido** — `RenderHTML` es una función pura
  `([]Book, time.Time) -> []byte`; `serve` solo la envuelve en `http.Handler`.

## Synopsis

```bash
estanteria serve [--addr 127.0.0.1:8080]
```

- Carga el ledger (`ESTANTERIA_FILE` o `~/.estanteria.json`) en cada request.
- Solo `GET`; cualquier otro método responde `405 Method Not Allowed`.
- Respuesta `text/html; charset=utf-8`, HTML autocontenido (CSS inline, sin
  JS obligatorio, escapado con `html.EscapeString`).

## Acceptance

- [x] GET sobre el handler devuelve 200 con Content-Type `text/html; charset=utf-8` y el título de un libro escapado.
- [x] POST/PUT/DELETE devuelven 405; el ledger no cambia.
- [x] Ledger vacío devuelve 200 con una marca de estantería vacía.
- [x] Un título `<script>` aparece escapado en el HTML (`&lt;script&gt;`).
- [x] `go test ./...` y `go vet ./...` en verde.
