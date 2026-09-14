---
type: task
status: in_progress
id: task-issue-4
title: "issue #4: Deploy automático de la vista web (GitHub Actions → Pages)"
assignee: Arggon
parent: story-imported-issues
labels: [enhancement]
created: "2026-09-14"
updated: "2026-09-14"
claimed_at: "2026-09-14T13:47:11.475Z"
issue: 4
---
## Propuesta

Workflow de GitHub Actions que genere la vista web como HTML estático y lo publique en GitHub Pages en cada push a main. El ledger commiteado es la fuente; el deploy es un snapshot de solo lectura.

## Criterios de aceptación

- [ ] workflow que construye la vista estática desde el ledger
- [ ] publica a Pages en push a main
- [ ] el sitio estático es 100% de solo lectura
> imported from issue #4
