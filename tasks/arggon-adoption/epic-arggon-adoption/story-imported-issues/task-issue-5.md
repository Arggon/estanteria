---
type: task
status: todo
id: task-issue-5
title: "issue #5: Backup automático del ledger"
parent: story-imported-issues
labels: [enhancement]
created: "2026-09-14"
updated: "2026-09-14"
issue: 5
---
## Propuesta

El ledger es el activo irremplazable del proyecto. Mecanismo de backup: copia fechada del ledger (local) con retención, pensada para correr antes de cada escritura o como subcomando.

## Criterios de aceptación

- [ ] copia de backup fechada del ledger
- [ ] retención configurable (N backups, los más viejos se podan)
- [ ] el backup nunca corrompe el ledger principal (escritura atómica también para el backup)
> imported from issue #5
