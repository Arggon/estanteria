---
type: task
status: todo
id: task-adopt-arggon
title: Adopt ArggonManager in this repo
parent: story-arggon-adoption
labels: []
created: "2026-09-14"
updated: "2026-09-14"
---
## Context

This repo is adopting ArggonManager over an existing documentation set: `arggon init` generated the governing docs (marked `<!-- arggon:generated ... -->`, with TODO placeholders), while any pre-existing docs were left untouched on disk. Your job as the executing agent: extract the valuable content from the adopter docs into the generated ones, archive what you replace, and report back on this task.

Current inventory (paths, sizes, managed vs adopter-owned, stack manifests):

```
arggon adopt --dry-run --json
```

## Checklist

- [ ] 1. Read the arggon-generated governing docs first: AGENTS.md, docs/convention.md, docs/engineering.md, docs/playbooks/ (if present). Follow them for the rest of this migration.
- [ ] 2. Sweep the existing repo docs (list them from the inventory above): extract the project description, conventions, workflows, and stack info. Extract, don't wholesale-copy — rewrite into the target doc's structure and drop duplicated or outdated material.
- [ ] 3. Complete the arggon-generated docs with the extracted content — fill the TODO placeholders: project description in AGENTS.md; CONTRIBUTING.md specifics (environment setup, build/test commands); ARCHITECTURE.md problem statement. The SECURITY.md contact is human input — leave it flagged for a human, never invent it.
- [ ] 4. Archive replaced originals to backup/<YYYY-MM-DD>/ preserving their relative paths (use today's date). Only docs you REPLACED get archived; never archive README.md — merge into it instead.
- [ ] 5. Detect the stack from the manifests (package.json / requirements.txt / go.mod / Cargo.toml / pom.xml); for each technology create a playbook (`arggon playbook new <tech>`), research current versions and best practices with dated sources, then record them with `arggon playbook refresh <tech> --version <v>`.
- [ ] 6. Baseline the sanctioned edits: run `arggon adopt --ack` so the generated docs you completed in step 3 become the new x-generated baseline (their checksums are refreshed and they stop reporting as modified). Hand edits made AFTER this ack still report modified — the protection stays intact.
- [ ] 7. Verify: `arggon validate` + `arggon spec validate` (if specs exist) + `arggon playbook status`.
- [ ] 8. Report: comment on this task (`arggon comment task-adopt-arggon`) listing the extracted content, archived files, and created playbooks; flip this task done when the human reviews.

### 2026-09-14 @Arggon
Adoption checklist executed (2026-09-14):

1. Generated governing docs read (AGENTS.md, convention.md, engineering.md).
2-3. Extraction: informal DECISIONS.md formalized into docs/adr/0001-go.md, 0002-json-file-storage.md, 0003-cli-first-web-readonly.md; AGENTS.md description filled; CONTRIBUTING.md got the Go toolchain + gates (go test/vet, was npm test); ARCHITECTURE.md completed with real layout/invariants; engineering.md testing expectations filled; SECURITY.md contact arggondev@gmail.com (owner-confirmed 2026-09-14) + re-confirm flag.
4. Archived: backup/2026-09-14/docs/DECISIONS.md (original informal log). README.md not archived (merged in place, per checklist).
5. Playbooks: docs/playbooks/go.md (v1.27.1) + docs/playbooks/golangci-lint.md (v2.13.2), researched 2026-09-14 with dated sources (go.dev release history, golangci-lint.run changelog).
6. arggon adopt --ack: 17 docs acked with sha256 checksums.
7. arggon validate ok:true, spec validate ok:true, playbook status staleCount:0.
Adopt hierarchy (fix #137 validated live): createdContainers [arggon-adoption, epic-arggon-adoption], storyCreated+taskCreated true.
