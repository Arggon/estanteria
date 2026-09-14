---
spec_id: backup-rotation-002
title: Rotación y retención de backups
status: implemented
created: 2026-09-14
---

# Spec: Rotación y retención de backups (backup-rotation-002)

## Purpose

Los backups acumulados en `~/.estanteria-backups/` crecen sin límite. Esta
spec agrega retención: `estanteria backup --keep N` poda, tras un backup
exitoso, los `ledger-*.json` más viejos dejando los `keep` más nuevos.

Invariantes:

- **La poda nunca toca el ledger principal**: solo opera dentro de `<dir>`.
- **La poda nunca elimina el backup más reciente** (ni ninguno de los `keep`
  más nuevos).
- **`--keep 0` = sin límite** (default; poda deshabilitada).
- Nunca se toca nada fuera de `<dir>` ni archivos que no matcheen el patrón
  `ledger-*.json` (p.ej. `otro.txt` sobrevive).

## Synopsis

```bash
estanteria backup [--dir RUTA] [--keep N]
```

- `--keep N` (int, default 0): tras el backup exitoso, elimina los
  `ledger-*.json` más viejos dejando los `N` más nuevos (orden por nombre,
  que es cronológico por el formato timestamp `ledger-YYYYMMDD-HHMMSS.json`).
- `N <= 0` → no hay poda. Si se podó algo, se imprime cuántos se eliminaron.
- `PruneBackups(dir string, keep int) (removed []string, err error)`:
  dir inexistente → no-op (sin error); errores de `os.Remove` se reportan.

## Acceptance

- [x] `PruneBackups` con 5 backups y keep=2 borra los 3 más viejos.
- [x] keep=0 (o negativo) no borra nada.
- [x] dir inexistente → no-op sin error (elegido y documentado).
- [x] Archivos que no matchean `ledger-*.json` nunca se borran.
- [x] El orden por nombre respeta el orden cronológico de los timestamps.
- [x] e2e: `backup --keep N` con varios backups pre-creados poda correcto.
- [x] Runbook y README actualizados con la política de retención.
