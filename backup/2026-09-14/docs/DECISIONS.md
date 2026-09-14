# Decisiones de diseño (informal)

Registro liviano de decisiones, en orden cronológico. Formato libre a
propósito: cuando el proyecto crezca, estas decisiones migran a ADRs formales.

## 1. Go como lenguaje (2026-09-14)

**Decisión:** todo el proyecto en Go, versión estable vigente (1.27.1).

**Por qué:** binario único y estático (cero runtime que instalar para usar el
CLI), stdlib con todo lo que este proyecto necesita (JSON, HTTP, testing), y
el modelo de concurrencia no se necesita pero el tipado fuerte sí — un ledger
personal se corrompe por bugs tontos, no por falta de frameworks. Alternativas
descartadas: Python (requiere interpreter + venv en cada máquina) y Rust
(excelente pero el costo de iteración no lo justifica para un CLI de ~500
líneas).

## 2. JSON file como storage (2026-09-14)

**Decisión:** el estado completo es un único archivo JSON, escrito
atómicamente (temp + rename en el mismo directorio).

**Por qué:** la estantería entera de una persona son cientos de libros, no
millones de filas — una base de datos es puro overhead. Un JSON es legible,
diffable con git, trivialmente respaldable y no tiene servidor que manterner.
La escritura atómica elimina el único riesgo real (ledger corrupto a medias si
muere el proceso escribiendo). La migración futura a SQLite, si algún día hiciera
falta, es un cambio contenido en `storage.go` — el resto del programa habla con
`Shelf`, no con el archivo.

## 3. CLI-first, web después (2026-09-14)

**Decisión:** la única forma *obligatoria* de interactuar es el CLI
(`add|list|status`). La vista web es un extra de solo lectura servida con
`net/http`, generada desde el ledger, sin capacidad de escritura.

**Por qué:** registrar un libro pasa mientras estoy leyendo o volviendo de la
biblioteca: la terminal está siempre a un keystroke, el navegador no. Una UI
web para escribir estado también arrastra autenticación, CSRF y validación de
formularios — tres fuentes de bugs para un beneficio nulo en un tracker de un
solo usuario. La vista web existe para *mirar* (la parte que sí es agradable en
un navegador), no para escribir.

## 4. Status canónicos en ASCII, display con acentos (2026-09-14)

**Decisión:** los status internos son `quiero-leer|leyendo|leido` (ASCII), el
display muestra `leído`. El parsing acepta ambas formas.

**Por qué:** los flags CLI con acentos son un infierno de teclados y shells;
perder la ortografía en el display sería un fastidio evitable. El folding de
acentos se aplica también a búsquedas: buscar `cronica` encuentra
`Crónica del pájaro...`.

## 5. ID = slug del título (2026-09-14)

**Decisión:** el id de cada libro es el slug de su título (minúsculas, sin
acentos, no-alfanuméricos → guión). Duplicado = error.

**Por qué:** evita inventar esquemas de ids que nadie va a recordar y hace que
el ledger sea auto-descriptivo (`el-aleph` se entiende sin lookup). El trade-off
asumido: dos libros con el mismo título no pueden coexistir — para un tracker
personal es un caso límite aceptable y el error es explícito.
