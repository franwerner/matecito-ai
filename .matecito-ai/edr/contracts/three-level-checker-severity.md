# EDR — Tres niveles de severidad, y sólo el más alto rompe

- **Status:** Accepted
- **Date:** 2026-08-06

## Contexto

Un validador que sólo distingue "pasa" de "falla" obliga a elegir, por cada chequeo, entre dos opciones malas: romper el trabajo por algo cosmético, o no reportarlo. La segunda gana casi siempre, y el resultado es un validador que sólo mira lo grave mientras el resto del deterioro avanza sin registro.

Los stores durables acumulan además defectos que no son de quien está trabajando hoy: material heredado, artefactos de cambios anteriores, deuda que nadie introdujo en esta tanda. Un validador binario los convierte a todos en bloqueo, y entonces se deja de correr.

## Decisión

El validador clasifica cada hallazgo en **tres niveles** y **sólo el más alto produce salida no exitosa**. Los otros dos se reportan y no frenan nada.

## Reglas verificables

- **[auto]** El validador sale con código de error únicamente cuando encontró al menos un hallazgo del nivel más alto.
- **[auto]** Todo hallazgo lleva su nivel explícito en la salida; ninguno se emite sin clasificar.
- **[manual]** Un hallazgo sobre material que el trabajo actual no tocó se reporta como contexto y nunca empuja el veredicto de ese trabajo.
- **[manual]** Que el validador no pueda correr en absoluto es un resultado distinto de que corra limpio, y se surfacea como tal — nunca se lo trata como éxito ni se degrada en silencio a una revisión manual.

## Alternativas consideradas

- **Dos niveles.** Descartado: empuja todo lo cosmético a ser invisible o bloqueante, que es exactamente el dilema del que este esquema sale.
- **Niveles configurables por proyecto.** Descartado: agrega superficie de configuración sin demanda real, y hace que el mismo hallazgo signifique cosas distintas según el repo — perdiendo la propiedad que hace útil al validador, que es que su salida se lea igual en todos lados.
- **Que todo hallazgo bloquee y se silencien los indeseados con excepciones.** Descartado: mueve la decisión de severidad a una lista de excepciones que nadie revisa, y que crece.

## Consecuencias

- El store puede acumular hallazgos de nivel medio sin frenar el trabajo, lo que permite correr el validador sobre material heredado desde el primer día.
- El riesgo real es que esos niveles se normalicen y nadie los mire. Lo mitiga que el trabajo en curso sí responda por lo suyo: los hallazgos sobre lo que uno tocó son propios, con la severidad que les corresponde.
- La distinción entre "corrió limpio" y "no pudo correr" pasa a ser parte del contrato, no un detalle de implementación: un chequeo que deja de correr en silencio es indistinguible de uno que pasa.

## Relacionados

- `relacionado-con` → [../structure/two-scripts-render-and-validate.md](../structure/two-scripts-render-and-validate.md) — el ejecutable separado que emite estos niveles.
