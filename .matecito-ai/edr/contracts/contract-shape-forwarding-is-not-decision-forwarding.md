# EDR — El canal de forwarding explícito es exclusivo de la forma de contrato; las decisiones ya no lo usan

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
El párrafo que introduce "Forwarding a ratified contract shape to the proposing phase", dentro de
`### Unresolved Decisions Guard`, distingue dos tipos de item que ese guard puede hacer disparar: una
decisión bajo `## New Decisions` y una forma de contrato bajo `### Contract Shapes Proposed`. Sin
explicar la diferencia, un lector podía asumir que ambos tipos viajan por el mismo canal de forwarding
de vuelta a la fase que los propuso — y no es así: sólo uno de los dos lo necesita.

## Decisión
`sdd-apply` lee `## New Decisions` directamente del artefacto de diseño (ver "In-Flow Decision
Capture"); no hay ningún gate intermedio que reescriba el item antes de que `sdd-apply` lo consuma, así
que no hace falta ningún canal de forwarding para las decisiones. La forma de contrato es distinta: se
ratifica en un gate que corre DESPUÉS de que la fase que la propuso ya retornó, y sólo esa misma fase
puede escribir la forma ratificada en su propio artefacto o código — el forwarding explícito, en las
instrucciones del re-despacho, es la única manera de devolvérsela.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Unresolved Decisions Guard
  (MANDATORY)`, subsección "Forwarding a ratified contract shape to the proposing phase" — el
  forwarding explícito aplica sólo a la forma de contrato, nunca a un item de `## New Decisions`.

## Alternativas consideradas
Mantener el mismo canal de forwarding para decisiones y formas de contrato — descartada: las decisiones
no pasan por ningún gate intermedio que las reescriba antes de que `sdd-apply` las lea, así que
forwardearlas sería un paso sin destinatario real.

## Consecuencias
Un lector de la sección de forwarding no necesita inferir por qué sólo la forma de contrato aparece
ahí: la exclusividad queda fijada en este record, no en un comentario del payload.
