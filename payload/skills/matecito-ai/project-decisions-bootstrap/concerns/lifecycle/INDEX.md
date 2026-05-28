# Dominio: `lifecycle`

Cómo se gestionan los datos a lo largo del tiempo: migraciones de esquema, backups y restore, retención, archivado y borrado de datos.

## Criterio de pertenencia

Un concern nuevo va en `lifecycle` si trata sobre el *ciclo de vida temporal* de los datos (versionado de esquema, cuánto se guardan, cómo se borran/archivan). Distinto de `data` (acceso en uso normal) y de `privacy` (derecho legal sobre datos personales).

## Concerns en este dominio

_Dominio reservado: todavía no tiene concerns. Es un casillero válido de la taxonomía, listo para poblarse vía ratchet cuando un proyecto lo necesite._

Para agregar el primer concern: creá `<slug>.md` con el formato estándar (ver `../runtime/error-handling.md` como referencia de fase `deep` o `../runtime/caching.md` para `light`), sumá la fila acá y la fila correspondiente en `../INDEX.md`.
