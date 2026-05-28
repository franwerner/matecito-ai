# Dominio: `integration`

Cómo el sistema se integra con servicios externos: clientes de APIs de terceros, webhooks entrantes, mensajería/colas, patrones anti-corruption, política de reintentos a terceros.

## Criterio de pertenencia

Un concern nuevo va en `integration` si trata sobre *consumir o intercambiar* con sistemas externos. Distinto de `contracts` (lo que este sistema expone) y de `runtime/resilience` (que cubre el comportamiento ante fallos en general).

## Concerns en este dominio

_Dominio reservado: todavía no tiene concerns. Es un casillero válido de la taxonomía, listo para poblarse vía ratchet cuando un proyecto lo necesite._

Para agregar el primer concern: creá `<slug>.md` con el formato estándar (ver `../runtime/error-handling.md` como referencia de fase `deep` o `../runtime/caching.md` para `light`), sumá la fila acá y la fila correspondiente en `../INDEX.md`.
