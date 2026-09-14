# estanteria

Tracker de lectura personal. Un ledger JSON plano + un CLI mínimo en Go, sin
dependencias: tu estantería te tiene que sobrevivir a las modas de bases de
datos.

## Demo online (GitHub Pages)

Cada push a `main` dispara [`.github/workflows/pages.yml`](.github/workflows/pages.yml):
compila el binario, renderiza el ledger de ejemplo
[`examples/ledger.json`](examples/ledger.json) con `estanteria render --out
index.html` y lo publica en GitHub Pages. El resultado es un sitio estático
100% de solo lectura (mismo generador que `serve`, ADR 0004) — una demo, no
tu ledger personal.

## Stack

- **Go 1.27.1** (estable al 2026-09-14; 1.27 salió el 19/8/2026 y 1.27.1 el 1/9/2026)
- Solo stdlib: `encoding/json`, `net/http` (vista web mínima), `flag`, `testing`
- Tests table-driven con `go test ./...`

## Motivación

Quiero llevar registro de lo que leo (quiero-leer / leyendo / leído + rating al
terminar) sin fricción: un comando desde la terminal, un archivo JSON que
puedo leer, versionar y respaldar con `git`. Nada de servidores, cuentas ni
interfaces web obligatorias — la vista web es un extra de solo lectura, no la
casalinga del estado.

## Uso

```bash
estanteria add "El Aleph" --autor "Jorge Luis Borges" --paginas 210
estanteria status "el al" --set leyendo
estanteria status "el aleph" --set leido --rating 5
estanteria list                    # toda la estantería
estanteria list --status leido     # solo terminados (con sus ratings)
estanteria backup                  # copia fechada en ~/.estanteria-backups/
estanteria backup --dir /tmp/bk    # copia fechada en otro directorio
estanteria stats                   # totales, rating promedio y páginas/mes (solo lectura)
estanteria serve                   # vista web de SOLO LECTURA en http://127.0.0.1:8080
```

`serve` levanta una vista HTML agrupada por status (quiero-leer / leyendo /
leído) con ratings y estadísticas: es **solo lectura** (ningún endpoint muta
el ledger, métodos que no sean GET responden 405) y bindea a `127.0.0.1` por
defecto (`--addr` para cambiarlo). La página se genera en cada request desde
el ledger, así que siempre muestra el estado actual.

El libro se busca por id exacto o por prefijo de título, sin distinguir
mayúsculas ni acentos. El rating (1-5) solo aplica al terminar un libro.

El estado vive en `$ESTANTERIA_FILE` o `~/.estanteria.json`. La escritura es
atómica (archivo temporario + rename), así que un corte de luz nunca deja un
ledger corrupto a medias.

`estanteria backup` crea una copia fechada del ledger
(`<dir>/ledger-YYYYMMDD-HHMMSS.json`, timestamp UTC, permisos 0600, escritura
atómica) sin tocar jamás el ledger principal. Es la operación recomendada
antes de editar el ledger a mano; la verificación y el restore están
documentados en [docs/runbooks/backup-restore.md](docs/runbooks/backup-restore.md).

## Esquema del estado

Ver [docs/FORMAT.md](docs/FORMAT.md). Resumen: `{"books": [...]}` con cada
libro con `id`, `title`, `status` (`quiero-leer|leyendo|leido`), `rating`
opcional y timestamps.

## Desarrollo

```bash
go test ./...      # suite completa
go vet ./...       # análisis estático
go build -o estanteria .
```

## Decisiones

Las decisiones de diseño están formalizadas como ADRs en
[docs/adr/](docs/adr/): por qué Go, por qué un JSON file y por qué CLI-first.
