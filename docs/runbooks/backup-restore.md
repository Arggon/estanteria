<!-- runbook: backup-restore -->

# Runbook: backup / restore del ledger

Backup y restauración del ledger `~/.estanteria.json` (o `$ESTANTERIA_FILE`).
Backup: `estanteria backup` (story-ledger-backup). La rotación/pruning de
backups viejos es otra story (story-backup-rotation) y todavía no existe.

## Trigger — cuándo correrlo

- **Antes de cualquier edición masiva a mano del ledger** (reordenar, merge
  de ids duplicados, migración de esquema a mano).
- Antes de restaurar por encima del ledger actual (el paso de restore lo
  pisa).
- Antes de operaciones experimentales contra el ledger (scripts de una vez,
  `jq` con salida redirigida, etc.).
- No hay trigger automático: es manual y barato (milisegundos, unos KB).

## Diagnosis — cuándo restaurar

- `estanteria list` falla con `parseando ...` → el ledger está corrupto. Por
  diseño (ver `docs/FORMAT.md`) falla ruidoso: **nunca** se arranca de cero,
  los datos siguen ahí pero inlegibles para la app.
- Perdiste libros por una edición a mano defectuosa.
- Confirmá el daño: `cat ~/.estanteria.json | python3 -m json.tool` (o
  `jq .`) → error de parseo = ledger corrupto; JSON válido con libros
  faltantes = edición defectuosa.

## Mitigation — verificación y restore

### 1. Crear (o confirmar que existe) un backup

```bash
estanteria backup
# imprime la ruta, p.ej.: /home/tu/.estanteria-backups/ledger-20260914-134530.json
```

Default: `~/.estanteria-backups/`; con `--dir RUTA` elegís otro destino.
El backup se crea atómicamente (temp + fsync + rename, permisos 0600) y el
ledger principal **nunca** es escrito por `backup`.

### 2. Verificar que un backup sirve

```bash
# A: con la propia herramienta
ESTANTERIA_FILE=~/.estanteria-backups/ledger-20260914-134530.json estanteria list

# B: solo parseo JSON
python3 -c "import json; json.load(open('.../ledger-20260914-134530.json'))"
```

Si `list` muestra los libros esperados, el backup es utilizable.

### 3. Restaurar

```bash
estanteria backup --dir /tmp/pre-restore   # cinturón extra antes de pisar
cp ~/.estanteria-backups/ledger-20260914-134530.json ~/.estanteria.json
chmod 600 ~/.estanteria.json
ESTANTERIA_FILE=~/.estanteria.json estanteria list   # verificación final
```

Importante: el `cp` pisa el ledger actual sin copia de seguridad suya — por
eso el backup previo del paso 3. No editar el archivo durante el restore.

### 4. Si el ledger principal está corrupto y no hay backup

El diseño es fail-loud (`docs/FORMAT.md`): la app se niega a operar sobre un
ledger corrupto para no pisar datos. Mirá si el archivo tiene un prefijo
reconocible (p.ej. un `.` backup de editor, `~/.estanteria.json.orig`) o si
está versionado con `git` (caso de uso documentado en el README). No
reconstruyas "a ojo" salvo que el contenido sea trivial.

## Escalation

- Sin backup utilizable y datos irrecuperables localmente → abrí un issue en
  github.com/Arggon/estanteria con: salida exacta del error, `ls -la` de
  `~/.estanteria-backups/`, y qué comandos corriste. No borres el ledger
  corrupto: puede recuperarse parcialmente a mano.

## Rollback

- El restore es reversible: el ledger pisado queda como `git`-versionable o,
  mínimamente, duplicá antes: `cp ~/.estanteria.json ~/.estanteria.json.pre-restore`.
- `backup` es de solo-lectura sobre el ledger: para "deshacer" un backup,
  simplemente borrá el archivo en `~/.estanteria-backups/` (no afecta al
  ledger principal).
