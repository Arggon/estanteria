---
type: story
status: done
id: story-pages-deploy
title: Deploy de la vista a GitHub Pages
assignee: Arggon
branch: feat/story-pages-deploy
parent: epic-estanteria-app
labels: []
created: "2026-09-14"
updated: "2026-09-14"
depends_on: [story-web-view]
---
<!--
  Placement (v0): tasks/estanteria-app/epic-estanteria-app/story-pages-deploy/story-pages-deploy.md (story index; required).
  parent MUST be the epic id. Optional style prefixes (e.g. story-) are not type discriminators.
-->

# Deploy de la vista a GitHub Pages

## Context

<!-- Why this story exists. -->

## Acceptance

- [x] estanteria render --out <file> genera la vista estática desde un ledger
- [x] workflow .github/workflows/pages.yml publica a Pages en push a main
- [x] sitio 100% de solo lectura; workflow commiteado en el PR
- [x] Closes GitHub issue #4; task-issue-4 actualizado con referencia al PR

## Notes
