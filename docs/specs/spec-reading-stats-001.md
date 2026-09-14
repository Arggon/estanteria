---
spec_id: reading-stats-001
title: Stats de lectura
status: implemented
created: 2026-09-14
---

# Spec: Stats de lectura (reading-stats-001)

## Purpose

Responder "¿qué y cuánto leí?" sin abrir el ledger: totales por status,
promedio de ratings y páginas por mes (ventana de 12 meses). Comando de solo
lectura.

Invariantes:

- **Pure read — stats nunca muta el ledger**: `ComputeStats` es una función
  pura sobre `[]Book`; `estanteria stats` solo ejecuta `LoadShelf`, jamás
  `SaveShelf`.
- Libros sin `Pages > 0` o sin `Finished` no suman páginas, pero sí cuentan
  en totales y conteos por status.
- El promedio de ratings considera solo libros con `rating != 0`; sin datos,
  es 0. Redondeado a 2 decimales.

## Synopsis

```bash
estanteria stats
```

Salida (stdout): totales por status, promedio de ratings (`%.2f/5`) y una
línea por mes con páginas (`2006-01`), solo meses con datos > 0, en orden
ascendente.

## Acceptance

- [ ] `ComputeStats([]Book{}, now)` devuelve ceros y mapa vacío (o mapa sin meses con datos).
- [ ] Conteos por status correctos con libros mixtos (incluye sin pages/finished).
- [ ] Promedio = suma(ratings)/cantidad(ratings != 0), 2 decimales; 0 si no hay.
- [ ] `PagesPerMonth` acumula `Pages` en el mes `YYYY-MM` de `Finished`, solo dentro de la ventana de 12 meses que termina en el mes de `now`.
- [ ] `go test ./...` y `go vet ./...` en verde; tests table-driven de `ComputeStats` + e2e de `stats`.
- [ ] README (sección Uso) actualizado en el mismo PR.
