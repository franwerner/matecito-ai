# EDR — La continuidad de apply-progress se ata al rol que escribe, no a la fase entera

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La regla de continuidad de apply-progress asumía un solo rol de despacho por batch de apply — leer y fusionar, un solo escritor. Un fragmento de dominio puede declarar más de un rol para la misma fase (por ejemplo, un rol aislado que nunca persiste, junto a un rol de consolidación que sí — ver el fan-out de `sdd-apply` en development). Sin esta aclaración, el fragmento de dominio tendría que contradecir la regla del kernel en vez de especializarla.

## Decisión
Cuando el fragmento del dominio declara más de un rol de despacho para la misma fase, la regla de continuidad ata SÓLO al rol que el fragmento nombra como el escritor. Un rol que el fragmento dice que nunca persiste tampoco lee `apply-progress` — el mismo principio de escritor único, aplicado tanto a la lectura como a la escritura.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `#### Apply-Progress Continuity (MANDATORY)`, ata la regla de continuidad sólo al rol que el fragmento de dominio nombra como escritor, cuando declara más de un rol de despacho.

## Alternativas consideradas
Dejar la regla del kernel asumiendo un solo rol y forzar al fragmento de dominio a contradecirla explícitamente para su caso — descartada: una regla del kernel que un dominio tiene que contradecir para funcionar correctamente es una regla mal alcanzada, no una excepción legítima.

## Consecuencias
El kernel queda domain-agnostic sobre cuántos roles de despacho existen por fase, sin perder la garantía de escritor único que la regla original perseguía.
