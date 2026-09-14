# ADR 0003 — CLI-first; la web es solo lectura

- **Status:** Accepted (2026-09-14)
- **Context:** registrar un libro ocurre mientras se lee o se vuelve de la biblioteca; la terminal está siempre a mano, el navegador no.
- **Decision:** la única forma obligatoria de escribir estado es el CLI (`add|list|status`). La vista web (`net/http` stdlib) es un extra opcional de solo lectura, sin endpoints de escritura.
- **Consequences:** cero autenticación/CSRF/formularios para el camino de escritura (tres fuentes de bugs evitadas en un tracker de un solo usuario). La web existe para mirar (ratings, stats), no para escribir.
