---
spec_id: pages-deploy-002
title: Deploy de la vista estática a GitHub Pages
status: implemented
created: 2026-09-14
---

# Spec: Deploy de la vista estática a GitHub Pages (pages-deploy-002)

## Purpose

Publicar la vista web como sitio estático en GitHub Pages (issue #4), reusando
el generador de ADR 0004. Invariantes:

- El sitio publicado es 100% de solo lectura: un snapshot HTML autocontenido,
  sin API ni runtime.
- El workflow construye desde un ledger commiteado en `examples/`
  (`examples/ledger.json`), nunca del ledger personal de nadie.
- `render` reutiliza `RenderHTML` (mismo generador que `serve`); no hay
  segundo generador de HTML.

## Synopsis

```bash
estanteria render --out index.html [--file RUTA]
```

- `--out` (requerido): archivo de salida; escritura atómica (temp + rename),
  permisos 0644 (es contenido público). Nunca muta el ledger de origen.
- `--file`: ledger de entrada; default `ESTANTERIA_FILE` o `~/.estanteria.json`
  (igual que el resto de los comandos).
- Imprime la ruta escrita.

## Workflow

`.github/workflows/pages.yml`: en push a `main` → checkout → setup-go 1.27.1
→ `go build -o estanteria .` → `./estanteria render --out index.html --file
examples/ledger.json` → publicar con las actions oficiales
(`actions/configure-pages`, `actions/upload-pages-artifact`,
`actions/deploy-pages`) con `permissions: pages:write, id-token:write`,
environment `github-pages` y concurrency con cancel-in-progress.

## Acceptance

- [ ] `estanteria render --out X` escribe HTML válido que contiene los títulos del ledger y no deja `.tmp` residuales.
- [ ] `render` nunca modifica el ledger de entrada.
- [ ] El YAML del workflow parsea y usa solo actions oficiales de GitHub.
- [ ] `examples/ledger.json` parsea como ledger válido (statuses mixtos, ratings, finished).
