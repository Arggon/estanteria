---
type: story
status: in_progress
id: story-reading-stats
title: "Stats de lectura: páginas/mes y ratings promedio"
assignee: Arggon
branch: feat/story-reading-stats
parent: epic-estanteria-app
labels: []
created: "2026-09-14"
updated: "2026-09-14"
claimed_at: "2026-09-14T13:31:02.868Z"
worktree_path: /home/arggon/Projects/estanteria-story-reading-stats
---
<!--
  Placement (v0): tasks/estanteria-app/epic-estanteria-app/story-reading-stats/story-reading-stats.md (story index; required).
  parent MUST be the epic id. Optional style prefixes (e.g. story-) are not type discriminators.
-->

# Stats de lectura: páginas/mes y ratings promedio

## Context

<!-- Why this story exists. -->

## Acceptance

- [ ] estanteria stats muestra páginas leídas por mes (últimos 12 meses, desde finished+pages)
- [ ] muestra promedio de ratings y conteos por status
- [ ] libros sin pages o sin finished no rompen el cálculo
- [ ] tests table-driven del cálculo; go test ./... && go vet ./... en verde
- [ ] Closes GitHub issue #1; task-issue-1 actualizado con referencia al PR

## Notes
