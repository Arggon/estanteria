<!-- runbook: backup-restore -->

# Runbook: backup / restore del ledger

Backup y restauración del ledger `~/.estanteria.json` (o `$ESTANTERIA_FILE`).
Backup: `estanteria backup` (story-ledger-backup), con retención opcional vía
`--keep N` (story-backup-rotation, ver sección **Retención** abajo).

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

### 2. Retención (`--keep N`)

```bash
estanteria backup --keep 7
```

Tras un backup exitoso, `--keep N` elimina los `ledger-*.json` más viejos del
directorio destino dejando los `N` más nuevos (orden por nombre, que es
cronológico por el formato del timestamp). Garantías de la poda:

- Nunca toca el ledger principal ni nada fuera del directorio de backups.
- Nunca elimina un archivo que no matchee `ledger-*.json` (p.ej. notas
  propias dentro del directorio sobreviven).
- Nunca elimina ninguno de los `N` backups más nuevos (el más reciente está
  siempre a salvo).
- `--keep 0` (default) = **sin límite**: no hay poda.

Política recomendada: si corrés `backup` con frecuencia (pre-edición,
pre-restore), usá `--keep 30`; el default ilimitado solo es razonable si el
uso es esporádico. Si la poda falla a mitad de camino, `backup` devuelve
error nonzero: los eliminados ya removidos se listan y el resto puede
podiarse re-corriendo `backup --keep N`.

### 3. Verificar que un backup sirve

```bash
# A: con la propia herramienta
ESTANTERIA_FILE=~/.estanteria-backups/ledger-20260914-134530.json estanteria list

# B: solo parseo JSON
python3 -c "import json; json.load(open('.../ledger-20260914-134530.json'))"
```

Si `list` muestra los libros esperados, el backup es utilizable.

### 4. Restaurar

```bash
estanteria backup --dir /tmp/pre-restore   # cinturón extra antes de pisar
cp ~/.estanteria-backups/ledger-20260914-134530.json ~/.estanteria.json
chmod 600 ~/.estanteria.json
ESTANTERIA_FILE=~/.estanteria.json estanteria list   # verificación final
```

Importante: el `cp` pisa el ledger actual sin copia de seguridad suya — por
eso el backup previo del paso 3. No editar el archivo durante el restore.

### 5. Si el ledger principal está corrupto y no hay backup

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
