# EDR — El veredicto de una propuesta rechazada vive en su propia sección condicional, no en un mailbox existente

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
El veredicto de una propuesta rechazada no encaja en ninguno de los dos mailboxes existentes de sdd-apply: `### Unmandated Forks` es para un punto que NINGÚN artefacto fija (acá el diseño SÍ lo fija, dos veces si hay conflicto); `### Mandated Departures` es para lo que la corrida SÍ aplicó (un veredicto `conflicts` es exactamente lo que la corrida NO aplicó). Forzarlo en cualquiera de los dos sería una clasificación falsa que además arrastraría los tokens `mandate:`/`verify-checks:`, que no significan nada para este chequeo.

## Decisión
Sección propia, `### Rejected Proposals Checked`, `emitted: conditional`, gatillada por `has_forwarded_rejections` — un hecho sobre el PROMPT de despacho, no derivable de la lista (una lista vacía con el gate en `true` es "todo lo forwardeado se chequeó y no hubo conflicto"; el gate en `false` es "no se forwardeó nada"). Declara `items.rationale` (split summary/rationale) igual que sus dos hermanas. No lleva fila de Tier en la Sección D.3 de `sdd-phase-common.md`: no es un mailbox que el guard clasifique en Tier 1/2 — su único consumidor es el chequeo mecánico `conflicts → status: blocked` que vive en la guardia de contenido.

## Reglas verificables
- **[auto]** `sdd-apply.yaml` declara la sección `emitted: conditional`, `when: has_forwarded_rejections`, sin gate derivable — `render-return.js` falla el render si el campo no se provee explícitamente.
- **[manual]** Ningún guard del orquestador clasifica esta sección como Tier 1 o Tier 2; su lectura es la comparación literal `conflicts → status: blocked` de la guardia de contenido, no una decisión que el usuario deba ratificar en el gate de mailboxes.

## Alternativas consideradas
Meterlo en `### Unmandated Forks`: descartado — ese mailbox es para un punto que NINGÚN artefacto fija; acá el diseño lo fija (a veces dos veces, incompatiblemente). Meterlo en `### Mandated Departures`: descartado — ese mailbox reporta lo que SÍ se aplicó; un veredicto `conflicts` es justamente lo que no se aplicó. `emitted: always` con un sentinel: descartado — el spec lee un sentinel y una ausencia de la misma forma, así que `always` no compra nada sobre `conditional`.

## Consecuencias
Un batch que no forwardeó ninguna propuesta rechazada no gana ninguna sección nueva ni ninguna obligación nueva. El guard de Unresolved Decisions no necesita una fila de Tier nueva: la clasificación de esta sección vive en la guardia de contenido, no en el guard genérico de mailboxes.

## Relacionados
- `relacionado-con` → [rejected-proposal-verdict-token.md](rejected-proposal-verdict-token.md) — los dos tokens que esta sección declara por ítem.
