---
spec_id: issue-search-001
title: Búsqueda por título/autor
status: implemented
created: 2026-09-14
---

# Spec: Búsqueda por título/autor (issue-search-001)

## Purpose

Encontrar libros en el ledger sin conocer el id exacto: `estanteria search <q>`
lista todo libro cuyo título o autor contenga la consulta, sin distinguir
mayúsculas ni acentos (mismo folding que `Find`: `foldAccents` + `ToLower`).

Invariantes:

- **Pure read**: `search` nunca muta el ledger — carga con `LoadShelf` y solo
  imprime. Ningún camino de código de search llama a `SaveShelf`.
- Matching por *contains* (no solo prefijo), sobre título O autor.
- Resultados determinísticos: ordenados por título.
- Query vacía o solo whitespace → sin resultados (no error).

## Synopsis

```bash
estanteria search <q>    # imprime filas igual que `list` + "total: N"
```

Sin flags. Salida: misma fila de `list` (id, título, autor, status, rating);
sin matches imprime "(sin resultados para \"q\")" y `total: 0`.

## Acceptance

- [ ] Match por substring de título y de autor, case/accent-insensitive en ambos sentidos (query con acentos, datos con acentos).
- [ ] Sin match → salida vacía de filas + `total: 0`, exit 0.
- [ ] Query vacía/whitespace → `total: 0`, exit 0.
- [ ] Múltiples matches ordenados por título.
- [ ] `search` nunca modifica el archivo del ledger (verificado por test e2e).
- [ ] `go test ./...` y `go vet ./...` en verde.
