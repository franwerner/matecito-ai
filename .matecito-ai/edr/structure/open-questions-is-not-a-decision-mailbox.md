# EDR — `Open Questions` es puramente informativo; lo que fija una decisión va a `New Decisions`

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
`New Decisions` y `Open Questions` se solapaban: las dos recibían decisiones pendientes, el ejecutor
terminaba duplicando contenido entre ambas, y el usuario confirmaba lo mismo dos veces — la fatiga de
confirmación que este guard existe para evitar.

## Decisión
`Open Questions` deja de ser un buzón de decisiones. Carga lo que NO fija una decisión; cualquier cosa
que sí fija una queda en `New Decisions`, el buzón de decisiones propio de esa fase (`gates: reported`,
así que nada ahí gatea). `sdd-apply` lee `New Decisions` directo del artefacto de diseño.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Unresolved Decisions Guard (MANDATORY)`, declara que `Open Questions` es informativo y que lo que fija una decisión pertenece a `New Decisions`.

## Alternativas consideradas
Mantener ambos buzones aceptando decisiones pendientes, confiando en que el ejecutor elija uno —
descartada: es el estado que produjo la duplicación de contenido y la doble confirmación.

## Consecuencias
Un item que fija una decisión tiene un solo destino posible; `Open Questions` deja de necesitar
revisarse dos veces por el mismo contenido.
