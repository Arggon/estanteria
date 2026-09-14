---
playbook_id: golangci-lint
version: v2.13.2
researched: 2026-09-14
status: current
---

# golangci lint playbook (golangci-lint)

Technology playbook: the chosen version and the current best practices for
golangci-lint. The version/best-practices research happened when this file was
created — cite dated sources (URL + access date) in every section so the next
reader can re-verify, and keep this file current via
`arggon playbook status`.

## Setup

- Versión vigente: **golangci-lint v2.13.2**, publicada el 2026-08-28 (fuentes:
  [changelog oficial](https://golangci-lint.run/docs/product/changelog/) y
  [GitHub releases](https://github.com/golangci/golangci-lint/releases), consultados 2026-09-14).
- La serie v1 está cerrada: la última release v1 fue v1.64.8 (marzo 2025); todo
  lo nuevo es v2 (misma fuente).
- Instalación (binario oficial, NO `go install` — el proyecto lo desaconseja por
  builds no reproducibles): ver [guía de instalación](https://golangci-lint.run/docs/welcome/install/), consultada 2026-09-14.
- Corrida local: `golangci-lint run`; en CI como job aparte del `go vet`.

## Conventions

- Config v2: `.golangci.yml` en la raíz; habilitar por linter explícito en vez
  de confiar en defaults silenciosos (los defaults cambian entre minors).
- Linters base recomendados para este repo: `errcheck`, `staticcheck`, `govet`,
  `ineffassign`, `unused` — suficientes para un CLI stdlib-only sin ruido.
- `staticcheck` se actualiza con cada release de golangci-lint (v2.13.2 lo lleva
  a 0.8.0, [changelog](https://golangci-lint.run/docs/product/changelog/), 2026-08-28).
- Los hallazgos no se silencian con `//nolint` sin comentario explicativo; si un
  linter molesta estructuralmente, se discute en el ADR/PR correspondiente.

## Testing

- El lint NO reemplaza `go test ./...`: primero tests en verde, después lint.
- Gate recomendado en cada PR: `go vet ./... && golangci-lint run` — ambos
  deben salir limpios antes de push (se agrega a CI cuando el repo tenga workflow).
- Para probar que una regla nueva de lint vale la pena: correrla sobre el árbol
  completo y evaluar el ratio señal/ruido antes de fijarla en `.golangci.yml`.

## Security

- v2.13.2 es un release de fixes: menor entropía de caché + bugs de linters
  (`iface` 1.5.1, `staticcheck` 0.8.0) — fuente:
  [changelog oficial](https://golangci-lint.run/docs/product/changelog/), 2026-08-28, consultado 2026-09-14.
- No hay advisories conocidos afectando v2.13.2 al 2026-09-14.
- Pinneado: fijar la versión en CI (imagen/tag exacto, no `latest`) para builds
  reproducibles.

## Upgrade policy

- Re-research al marcarse stale (>90 días, `arggon playbook status`) o ante un
  minor nuevo de la serie v2 (ritmo actual: varios por trimestre).
- Cómo elegir próxima versión: último tag de
  [GitHub releases](https://github.com/golangci/golangci-lint/releases); leer el
  changelog de linters nuevos/deprecados antes de subir.
- Tras re-research: `arggon playbook refresh golangci-lint --version <v>` y
  actualizar Setup/Security con la fecha de la release.
