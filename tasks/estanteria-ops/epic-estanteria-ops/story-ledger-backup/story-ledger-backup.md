---
type: story
status: done
id: story-ledger-backup
title: Backup fechado del ledger
assignee: Arggon
branch: feat/story-ledger-backup
parent: epic-estanteria-ops
labels: []
created: "2026-09-14"
updated: "2026-09-14"
---
<!--
  Placement (v0): tasks/estanteria-ops/epic-estanteria-ops/story-ledger-backup/story-ledger-backup.md (story index; required).
  parent MUST be the epic id. Optional style prefixes (e.g. story-) are not type discriminators.
-->

# Backup fechado del ledger

## Context

<!-- Why this story exists. -->

## Acceptance

- [x] estanteria backup crea copia fechada del ledger (escritura atómica)
- [x] runbook de backup/restore en docs/runbooks/
- [x] tests table-driven; go test ./... && go vet ./... en verde
- [x] Closes GitHub issue #5; task-issue-5 actualizado con referencia al PR

## Notes
