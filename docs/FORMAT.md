# FORMAT.md — esquema del estado

## Ubicación

El estado completo vive en un único archivo JSON:

- `$ESTANTERIA_FILE` si está definida la variable de entorno, o
- `~/.estanteria.json` por defecto.

Escritura **atómica**: se escribe un temporario (`.<nombre>-*.tmp`, permisos
0600) en el mismo directorio, se le hace `fsync`, y se hace `rename` sobre el
archivo final. Los lectores nunca ven estado parcial. Archivos `*.tmp`
residuales no deberían existir tras un save exitoso.

## Documento raíz

```json
{
  "books": [ ... ]
}
```

Archivo inexistente == estantería vacía (no es error). `books` siempre está
presente en la salida (nunca `null`).

## Book

| Campo       | Tipo            | Presencia   | Descripción |
| ----------- | --------------- | ----------- | ----------- |
| `id`        | string          | siempre     | slug del título: minúsculas, sin acentos (`ñ`→`n`), no-alfanuméricos → guión. Único. |
| `title`     | string          | siempre     | título original, trimado |
| `author`    | string          | opcional    | omitido si vacío |
| `pages`     | int             | opcional    | omitido si 0; nunca negativo |
| `status`    | string          | siempre     | `quiero-leer` \| `leyendo` \| `leido` (ASCII canónico; display `leído`) |
| `rating`    | int             | opcional    | 1–5; **solo** tiene sentido con `status: leido`. Dejar `leido` vuelve a `leyendo` lo limpia. |
| `added`     | RFC 3339 UTC    | siempre     | momento del `add` |
| `finished`  | RFC 3339 UTC    | opcional    | se fija al pasar a `leido`; se limpia al salir de `leido` |

### Reglas de transición de status

- `quiero-leer → leyendo → leido` es el flujo normal; los saltos hacia atrás
  están permitidos.
- `rating` solo se acepta en la transición **hacia** `leido` (con `--set leido
  --rating N`); sin `--set`, un `--rating` se ignora — *(bug conocido, ver
  issue "rating sin --set se ignora silenciosamente")*.
- Pasar a `leido` fija `finished`; salir de `leido` limpia `rating` y
  `finished`.

## Invariantes

- Nunca se escribe un ledger que no parseó: un JSON corrupto **falla con
  error** en vez de arrancar de cero (protege los datos).
- `add` con id duplicado falla y **no muta** el shelf en memoria.
- Toda escritura pasa por `SaveShelf` (atómica); ninguna otra ruta escribe el
  archivo.
