# ADR 0002 — JSON file como storage, escritura atómica

- **Status:** Accepted (2026-09-14)
- **Context:** el estado a persistir es la estantería de una persona (cientos de registros, no millones). Un servidor de base de datos es overhead puro a esta escala.
- **Decision:** un único archivo JSON (`~/.estanteria.json` o `$ESTANTERIA_FILE`) como fuente de verdad, escrito exclusivamente vía `SaveShelf`: temporario en el mismo directorio + `fsync` + `rename`. Ledger corrupto falla ruidoso, nunca se reinicia de cero.
- **Consequences:** estado legible, diffable con git y respaldable con `cp`; sin servidor que mantener. La escritura atómica elimina el riesgo de ledger a medias ante un crash. Si algún día hiciera falta más capacidad, la migración queda contenida en `storage.go`.
