---
type: story
status: in_progress
id: story-backup-rotation
title: Rotación y retención de backups
assignee: Arggon
branch: feat/story-backup-rotation
parent: epic-estanteria-ops
labels: []
created: "2026-09-14"
updated: "2026-09-14"
claimed_at: "2026-09-14T13:43:11.382Z"
depends_on: [story-ledger-backup]
worktree_path: /home/arggon/Projects/estanteria-story-backup-rotation
---
<!--
  Placement (v0): tasks/estanteria-ops/epic-estanteria-ops/story-backup-rotation/story-backup-rotation.md (story index; required).
  parent MUST be the epic id. Optional style prefixes (e.g. story-) are not type discriminators.
-->

# Rotación y retención de backups

## Context

<!-- Why this story exists. -->

## Acceptance

- [ ] estanteria backup --keep N poda los backups más viejos (retención)
- [ ] la poda nunca toca el ledger principal ni el backup más reciente
- [ ] tests table-driven de la rotación; runbook actualizado
- [ ] cierra el alcance del issue #5 (rotación); task-issue-5 comentado

## Notes
