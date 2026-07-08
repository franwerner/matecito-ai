# Dominio: `domain-logic`

Decisiones de **cómo se modela el negocio**: el paradigma de modelado elegido, la estructura de las entidades y sus relaciones, dónde viven las reglas e invariantes, y la terminología del dominio. NO las reglas de negocio en sí (eso son capability-specs) — acá va la *decisión de modelado*, no el comportamiento. (DDD —agregados, límites de consistencia, lenguaje ubicuo— es un ejemplo de vocabulario, no un requisito.)

## Criterio de pertenencia

Un concern nuevo va en `domain-logic` si trata sobre *una decisión de cómo se modela el dominio* independiente de la tecnología — el paradigma, la forma de las entidades, dónde se ubican reglas e invariantes. Distinto de `structure` (cómo se organiza el código que las implementa) y de los capability-specs (cuáles son las reglas y qué hace el sistema).

## Concerns en este dominio

_Dominio reservado: todavía no tiene concerns. Es un casillero válido de la taxonomía, listo para poblarse vía ratchet cuando un proyecto lo necesite._

Para agregar el primer concern: creá `<slug>.md` con el formato estándar (ver `../runtime/error-handling.md` como referencia de fase `deep` o `../runtime/caching.md` para `light`), sumá la fila acá y la fila correspondiente en `../INDEX.md`.
