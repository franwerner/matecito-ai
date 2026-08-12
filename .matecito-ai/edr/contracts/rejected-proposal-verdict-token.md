# EDR — Token de veredicto obligatorio para una propuesta rechazada: design-conflict, ambos valores passing

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
sdd-apply ya bloquea (vía la guardia de contenido) cuando la implementación de una propuesta rechazada difiere de lo que describe el enfoque del diseño para el mismo punto — pero esa evaluación quedaba sin rastro verificable en el retorno: nada obligaba a declarar que la corrida efectivamente chequeó cada propuesta rechazada forwardeada, ni distinguía "chequeé y no hay conflicto" de "no llegué a chequear". Sin un token, el orquestador no puede auditar mecánicamente si la guardia corrió.

## Decisión
Cada propuesta rechazada cuya tarea gobernante esta corrida alcanzó emite un ítem con dos tokens: `record: <domain>/<slug>` (la identidad, sin `values`, formato libre, ya usado por el mecanismo de captura in-flow) y `design-conflict: none|conflicts` (el veredicto). Los dos valores de `design-conflict` están declarados `passing` — a diferencia de `mandate`/`verify-checks`, no hay un valor "ilegal en esta sección": un valor no-passing haría fallar el render exactamente cuando la corrida intenta reportar un conflicto honesto, lo que volvería imposible declararlo.

## Reglas verificables
- **[auto]** `render-return.js` exige ambos tokens por ítem de `rejected_proposals` y falla el render si falta cualquiera de los dos (comportamiento heredado del motor, confirmado corriendo `--data` incompleto).
- **[auto]** Los dos valores de `design-conflict` (`none`, `conflicts`) están declarados `passing` en `sdd-apply.yaml`, así que ningún valor legal del token puede hacer fallar el render.

## Alternativas consideradas
Un solo token `verdict` con `passing: [none]`: descartado — un valor no-passing haría fallar el render, volviendo imposible reportar honestamente un conflicto (misma razón por la que `verify-checks` mantiene ambos valores como passing). Una columna `Record` en vez de un token `record` libre: descartado — el spec pide un ítem por `record:`, y el token libre es la forma que ya usa el mailbox existente de decisiones materializadas.

## Consecuencias
Cada tarea gobernada por una propuesta rechazada que esta corrida alcanza queda auditada mecánicamente: el orquestador puede verificar que existe un ítem por cada `record:` forwardeado. Un veredicto `conflicts` es reportable sin que el render lo rechace, preservando la posibilidad de un `status: blocked` honesto.

## Relacionados
- `relacionado-con` → [verdict-in-its-own-conditional-section.md](verdict-in-its-own-conditional-section.md) — el hogar de este token — sección propia, no un mailbox existente.
