# EDR — La forma ratificada de un contrato viaja en el prompt de re-despacho, nunca se relee de un store

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Una fase que propone la forma de un contrato termina su corrida en el mismo turno en el que lo propuso; el gate que lo ratifica corre después, cuando esa fase ya no está viva. La forma ratificada (o ajustada por el usuario) tiene que llegar de alguna manera a la fase que la va a implementar, sin depender de que esa fase la vuelva a leer de algún store.

## Decisión
La forma ratificada viaja de vuelta en el prompt con el que se re-despacha a esa misma fase, identificada por el ancla y el resumen del ítem original. Ninguna fase lee la forma ratificada de un store — ni siquiera de su propio artefacto de progreso; el único canal es el prompt de despacho.

## Reglas verificables
- **[manual]** La fase que implementa un contrato ratificado recibe la forma completa (ancla, resumen, campos) en su propio prompt de despacho, nunca leyéndola de un store o artefacto propio.
- **[manual]** Una fase que llega al punto que un contrato gobierna sin la forma ratificada en su prompt devuelve un bloqueo nombrando la propuesta faltante, en vez de adivinarla o inferirla.

## Alternativas consideradas
Que la fase relea la forma ratificada de un store al momento de implementar — descartada porque introduce una segunda fuente de verdad (el store, además del prompt) que puede desalinearse de lo que el gate realmente ratificó, y porque el canal ya establecido para un caso equivalente (la resolución de una propuesta de decisión) usa el prompt, no un store.

## Consecuencias
La forma ratificada de un contrato es siempre exactamente lo que el prompt de re-despacho contiene, sin ambigüedad sobre cuál versión es la vigente — a costa de que quien construye ese prompt tiene que hacerlo con cuidado cada vez, y una fase que no la recibió no tiene ningún otro lugar donde buscarla antes de bloquear.

## Relacionados
- `relacionado-con` → [contract-item-anchors-once.md](contract-item-anchors-once.md) — la identidad que este reenvío usa para reconocer la propuesta.
- `relacionado-con` → [contract-proposal-has-no-persistence-slot.md](contract-proposal-has-no-persistence-slot.md) — por qué la propuesta original no declara dónde persistirá, reforzando que el prompt, no un store, es el canal.
