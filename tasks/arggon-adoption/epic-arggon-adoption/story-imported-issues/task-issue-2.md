---
type: task
status: in_progress
id: task-issue-2
title: "issue #2: Búsqueda por título/autor"
assignee: Arggon
parent: story-imported-issues
labels: [enhancement]
created: "2026-09-14"
updated: "2026-09-14"
claimed_at: "2026-09-14T13:42:21.941Z"
issue: 2
---
## Propuesta

Hoy `status <libro>` busca por prefijo de título. Falta un `estanteria search <query>` que busque en título Y autor, listando todos los matches (no solo el mejor).

## Criterios de aceptación

- [ ] `estanteria search <q>` matchea título y autor, sin distinguir mayúsculas ni acentos
- [ ] lista TODOS los matches con su status
- [ ] reutiliza el folding de acentos existente
- [ ] tests table-driven
> imported from issue #2

### 2026-09-14 @Arggon
Implementado y mergeado en PR #9 (story-issue-search, via claim race — fix #134 validado). Spec+plan: spec-issue-search-001.
