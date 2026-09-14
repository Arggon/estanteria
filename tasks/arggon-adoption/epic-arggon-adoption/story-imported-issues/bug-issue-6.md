---
type: bug
status: in_progress
id: bug-issue-6
title: "issue #6: status <libro> --rating N sin --set ignora el rating en silencio"
assignee: Arggon
branch: fix/bug-issue-6
parent: story-imported-issues
labels: [bug]
created: "2026-09-14"
updated: "2026-09-14"
claimed_at: "2026-09-14T13:49:49.055Z"
issue: 6
worktree_path: /home/arggon/Projects/estanteria-bug-issue-6
---
## Bug

`estanteria status <libro> --rating N` (sin `--set`) sale con código 0, imprime el libro como si nada, y el rating **no se aplica ni se avisa**. El usuario cree que calificó el libro; el ledger no cambia.

## Reproducción (verificada en v0.1.0, Go 1.27.1)

```
$ estanteria status el-aleph
El Aleph — leído (5/5) — terminado: 2026-09-14
$ estanteria status el-aleph --rating 1
El Aleph — leído (5/5) — terminado: 2026-09-14     # exit 0
$ estanteria status el-aleph
El Aleph — leído (5/5) — terminado: 2026-09-14     # rating sigue en 5
```

## Causa

En `cmdStatus` (main.go), si `--set` está vacío se salta `SetStatus` por completo y `--rating` se descarta sin validación ni aviso.

## Comportamiento esperado

Error explícito: el rating solo aplica con `--set leido` (o aplicar el rating si el libro ya está `leido` — a decidir en el fix, con test).

## Criterios de aceptación

- [ ] `status <libro> --rating N` sin `--set` falla con mensaje claro (o aplica el rating a un libro ya leído, decidido y documentado en FORMAT.md)
- [ ] test table-driven del caso
> imported from issue #6
