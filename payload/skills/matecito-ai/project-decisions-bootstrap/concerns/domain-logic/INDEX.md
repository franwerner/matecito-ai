# Dominio: `domain-logic`

Cómo se modela el negocio en sí: reglas de negocio, invariantes, lenguaje ubicuo, agregados y límites de consistencia (DDD).

## Criterio de pertenencia

Un concern nuevo va en `domain-logic` si trata sobre *las reglas del negocio* independientes de la tecnología. Distinto de `structure` (cómo se organiza el código que las implementa).

## Concerns en este dominio

_Dominio reservado: todavía no tiene concerns. Es un casillero válido de la taxonomía, listo para poblarse vía ratchet cuando un proyecto lo necesite._

Para agregar el primer concern: creá `<slug>.md` con el formato estándar (ver `../runtime/error-handling.md` como referencia de fase `deep` o `../runtime/caching.md` para `light`), sumá la fila acá y la fila correspondiente en `../INDEX.md`.
