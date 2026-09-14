---
spec_id: ledger-backup-001
title: Backup fechado del ledger
status: implemented
created: 2026-09-14
---

# Spec: Backup fechado del ledger (ledger-backup-001)

## Purpose

Copias de seguridad fechadas del ledger, para poder recuperar el estado antes
de una edición a mano masiva o un accidente. Invariantes:

- **El backup nunca corrompe el ledger principal**: la operación es de
  solo-lectura sobre la fuente; ninguna ruta de `backup` escribe el ledger.
- **Copia atómica**: el backup se escribe con temp + fsync + rename; un corte
  de luz nunca deja un backup a medias con nombre final.
- **Sin backup cruzado entre ledgers distintos**: se hace `LoadShelf` de la
  fuente antes de escribir nada; un archivo que no parsea no se backuppea.

## Synopsis

```bash
estanteria backup [--dir RUTA]
```

- Fuente: `$ESTANTERIA_FILE` o `~/.estanteria.json` (`LedgerPath()`).
- Destino: `<dir>/ledger-YYYYMMDD-HHMMSS.json` (timestamp **UTC** de `now`);
  colisión en el mismo segundo → sufijo `-1`, `-2`, …
- `<dir>` default: `~/.estanteria-backups`; se crea con `MkdirAll` si falta.
- Escritura atómica: temp en `<dir>` + fsync + rename, permisos `0600`.
- Si el ledger no existe o no parsea → error y **no se crea ningún archivo**.
- Salida: imprime la ruta del backup creado. La rotación/`--keep` NO va aquí
  (story-backup-rotation, con dependencia a esta).

## Acceptance

- [ ] `estanteria backup` crea `<dir>/ledger-YYYYMMDD-HHMMSS.json` parseable e idéntico en contenido al ledger.
- [ ] Colisión de timestamp en el mismo segundo produce sufijo `-1`.
- [ ] Ledger inexistente o corrupto → error, sin dejar archivos (ni `.tmp` residuales).
- [ ] El comando `backup` nunca escribe el ledger principal.
- [ ] Runbook `docs/runbooks/backup-restore.md` + README actualizados.
