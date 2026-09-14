---
type: story
status: in_progress
id: story-web-view
title: Vista web de solo lectura (serve)
assignee: Arggon
branch: feat/story-web-view
parent: epic-estanteria-app
labels: []
created: "2026-09-14"
updated: "2026-09-14"
claimed_at: "2026-09-14T13:31:06.504Z"
worktree_path: /home/arggon/Projects/estanteria-story-web-view
---
<!--
  Placement (v0): tasks/estanteria-app/epic-estanteria-app/story-web-view/story-web-view.md (story index; required).
  parent MUST be the epic id. Optional style prefixes (e.g. story-) are not type discriminators.
-->

# Vista web de solo lectura (serve)

## Context

<!-- Why this story exists. -->

## Acceptance

- [ ] estanteria serve [--addr] sirve HTML de solo lectura (por defecto 127.0.0.1)
- [ ] agrupado por status con ratings; ningún endpoint muta el ledger
- [ ] generador de HTML compartido con el render estático (ver ADR 0004)
- [ ] tests del handler con httptest; go test ./... && go vet ./... en verde
- [ ] Closes GitHub issue #3; task-issue-3 actualizado con referencia al PR

## Notes
