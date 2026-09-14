---
type: task
status: todo
id: task-issue-3
title: "issue #3: Vista web de la estantería (solo lectura, net/http stdlib)"
parent: story-imported-issues
labels: [enhancement]
created: "2026-09-14"
updated: "2026-09-14"
issue: 3
---
## Propuesta

`estanteria serve` que sirva una vista HTML de solo lectura en localhost con net/http (stdlib): estantería agrupada por status, ratings y stats básicos.

## Criterios de aceptación

- [ ] `estanteria serve [--addr 127.0.0.1:8080]` sirve HTML legible
- [ ] solo lectura: ningún endpoint muta el ledger
- [ ] bind por defecto a 127.0.0.1 (no exponer a la red)
- [ ] test del handler con httptest
> imported from issue #3
