---
playbook_id: go
version: v1.27.1
researched: 2026-09-14
status: current
---

# go playbook (go)

Technology playbook: the chosen version and the current best practices for
go. The version/best-practices research happened when this file was
created — cite dated sources (URL + access date) in every section so the next
reader can re-verify, and keep this file current via
`arggon playbook status`.

## Setup

- Versión fijada: **Go v1.27.1** (linux/amd64), estable vigente al 2026-09-14.
  Instalación local con mise: `mise use go@1.27.1` (fuente: [Go Release History](https://go.dev/doc/devel/release), consultado 2026-09-14 — 1.27.1 publicada el 2026-09-01).
- Soporte upstream: 1.27 y 1.26 reciben parches; **Go 1.25 llegó a EOL el 2026-08-19**
  (fuente: [endoflife.date/go](https://endoflife.date/go), consultado 2026-09-14). No arrancar branches sobre 1.25 o menor.
- Este repo no tiene dependencias externas: `go mod download` es no-op pero
  corre como `x-worktree.post-start` por uniformidad. Build: `go build -o estanteria .`.

## Conventions

- Solo stdlib para I/O básica (`encoding/json`, `net/http`, `flag`, `testing`).
  Si entra una dependencia nueva → es decisión transversal: ADR + exploración primero
  (metodología del skill arggon-cli, sección 5).
- Errores: `fmt.Errorf("contexto: %w", err)` (envolvimiento con `%w`), mensajes
  de cara al usuario en español consistente; nunca `panic` en rutas de CLI.
- Go 1.27 introduce **métodos genéricos** (type parameters propios en métodos),
  uno de los tres cambios notables de spec (fuentes: [Go 1.27 Release Notes](https://go.dev/doc/go1.27) y [InfoWorld, 2026-08](https://www.infoworld.com/article/4214999/go-1-27-brings-support-for-generic-methods.html), consultados 2026-09-14). Todavía no hay caso de uso en este repo; permitido pero no obligatorio.
- Estructura: dominio puro (`book.go`, `shelf.go`) sin I/O; toda escritura de
  archivo pasa por `storage.go` (ver `ARCHITECTURE.md` → Boundaries).

## Testing

- `go test ./...` — suite completa; `go vet ./...` como gate estático antes de push.
- Tests table-driven con subtests (`t.Run`) para toda regla de dominio: tabla
  de casos + `wantErr` como substring del mensaje (patrón dominante en este repo,
  ver `book_test.go`, `shelf_test.go`).
- E2E del CLI en proceso: `run(args)` devuelve código de salida (sin `os.Exit`)
  y los tests capturan stdout/stderr con `os.Pipe`; cada caso usa un ledger en
  `t.TempDir()` vía `t.Setenv("ESTANTERIA_FILE", ...)` (ver `main_test.go`).
- Invariantes de storage con tests dedicados: escritura atómica sin residuos
  `.tmp`, ledger corrupto falla ruidoso (`storage_test.go`).

## Security

- `gofix`/toolchain: la 1.27.1 trae fixes de cgo, compiler, runtime y `go fix`
  (fuente: [Go Release History](https://go.dev/doc/devel/release), 2026-09-01; consultado 2026-09-14) — mantener el toolchain al día en parches.
- El ledger vive en `~/.estanteria.json` con permisos `0600` (test lo verifica);
  la escritura atómica usa `os.CreateTemp` en el mismo directorio (sin carreras
  de /tmp compartido).
- Sin `cgo` ni red en el camino de escritura; la vista web bindea a 127.0.0.1.

## Upgrade policy

- Re-research al marcase stale (>90 días, `arggon playbook status`) o al salir
  un minor de Go (cadencia semestral: feb/ago).
- Cómo elegir próxima versión: último stable de [go.dev/doc/devel/release](https://go.dev/doc/devel/release);
  actualizar `mise use go@<v>`, `go.mod`, y este playbook con
  `arggon playbook refresh go --version <v>`.
- Secciones a actualizar en cada upgrade: Setup (versión + fecha), Conventions
  (cambios de lenguaje relevantes), Security (fixes del patch release).
