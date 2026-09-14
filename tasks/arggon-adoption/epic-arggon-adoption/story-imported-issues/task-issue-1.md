---
type: task
status: todo
id: task-issue-1
title: "issue #1: Stats de lectura (páginas/mes, ratings promedio)"
parent: story-imported-issues
labels: [enhancement]
created: "2026-09-14"
updated: "2026-09-14"
issue: 1
---
## Propuesta

Comando `estanteria stats` que resuma la actividad de lectura:

- páginas leídas por mes (a partir de `finished` y `pages`)
- promedio de ratings de los libros terminados
- total de libros por status

## Criterios de aceptación

- [ ] `estanteria stats` muestra páginas/mes de los últimos 12 meses
- [ ] muestra el promedio de ratings de los libros con rating
- [ ] libros sin `pages` o sin `finished` no rompen el cálculo
- [ ] tests table-driven del cálculo
> imported from issue #1

### 2026-09-14 @Arggon
Implementado y mergeado en PR #8 (story-reading-stats). Spec+plan: spec-reading-stats-001.
