# EDR — El brief de sdd-intake conserva cuatro secciones y pierde discovery, triage y el tamaño

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
sdd-intake pasa de dos pasadas (formular preguntas, luego producir el brief) a una sola pasada sin ciclo de discovery, sin recomendación de lane y sin estimar tamaño — el modelo de dos lanes fijos no deja nada que triagear, y el ciclo de discovery se mudó a sdd-explore.

## Decisión
El bloque `[done, needs-decision, blocked]` de sdd-intake.yaml conserva `### Request (structured)`, `### Classification`, `### Early guard (EDRs)` (condicional, sin cambios) y `### Next`. Pierde `### Discovery answers` (se mudó a sdd-explore), `### Triage` (no hay lane que recomendar) y el bullet `Size` dentro de Classification. Al quedar una sola variante, `variants:` colapsa a un `block:` + `sections:` de nivel superior, y `needs-input` deja la lista de statuses de esta fase.

## Reglas verificables
- **[auto]** `render-return.js --phase sdd-intake --data <payload con status needs-input>` falla con `status is "needs-input", not one of [...]` — needs-input ya no es un status legal de esta fase.
- **[auto]** `validate-return.js --self-check --phase sdd-intake` falla si `### Discovery answers`, `### Triage` o un bullet `Size` reaparecen en el template sin estar declarados en sdd-intake.yaml, o si se declaran sin aparecer en el template.

## Alternativas consideradas
Mantener `### Triage` como una línea fija `Lane: full` — descartado porque un campo con un único valor posible se lee como un mecanismo vivo para quien lo encuentre después, exactamente el tipo de falla que este cambio completo viene a eliminar.

## Consecuencias
El brief queda más corto y sin ningún campo que ya no varía. sdd-intake deja de ser la única fase con un contrato `variants:` por status; ese patrón ahora vive en sdd-explore.
