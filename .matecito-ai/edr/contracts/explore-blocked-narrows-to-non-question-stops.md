# EDR — El status blocked de sdd-explore se restringe a lo que una segunda pasada con respuestas no puede resolver

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
sdd-explore corre headless y, antes de este cambio, usaba `blocked` tanto para un pedido demasiado ambiguo como para un obstáculo real (falta de acceso, una contradicción entre entradas ya ratificadas). El contrato de retorno ya distingue needs-input (una pregunta que un re-despacho con respuestas resuelve) de blocked (lo que no).

## Decisión
Una pregunta sobre el pedido en sí — el caso "demasiado vago para explorar" — es ahora una pregunta de discovery y se devuelve como `needs-input`. `blocked` se restringe a un obstáculo que un re-despacho CON respuestas no puede resolver: falta de acceso, o una contradicción entre las entradas ya ratificadas.

## Reglas verificables
- **[manual]** sdd-explore.md designa `needs-input` para la ambigüedad del pedido y reserva `blocked` para un obstáculo que un re-despacho con respuestas no resuelve — enforced por lectura del template en cada despacho, no por un script.
- **[manual]** Ningún script puede clasificar si un obstáculo dado es ambigüedad del pedido o falta de acceso/contradicción — es criterio del ejecutor, chequeado sólo por revisión.

## Alternativas consideradas
Mantener ambas rutas y dejar que el ejecutor elija caso por caso — descartado porque dos statuses para la misma situación es cómo una fase termina usando el que frena más duro, en vez del que corresponde.

## Consecuencias
Un pedido ambiguo ya no cierra la exploración: vuelve como una pregunta de discovery que el usuario puede responder, en vez de un blocker que exige una intervención distinta.
