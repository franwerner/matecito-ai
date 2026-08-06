# EDR — El índice raíz enruta por agrupador, no por artefacto

- **Status:** Accepted
- **Date:** 2026-08-06

## Contexto

Los stores durables tienen dos niveles de índice: uno raíz que enruta hacia un agrupador (el dominio de una decisión, el tipo de una capacidad) y uno por agrupador que lista sus artefactos. Cuando se materializa una tanda de varios artefactos, cada uno querría insertar su fila en ambos — y en el raíz, varios artefactos del mismo agrupador producirían la misma fila repetida.

El renderizador conoce **un artefacto por invocación**. No sabe si es el primero de su agrupador o el quinto.

## Decisión

El índice raíz lleva **una fila por agrupador**, no por artefacto, y la deduplicación es responsabilidad de quien invoca al renderizador, no del renderizador. El renderizador emite la fila que corresponde a su artefacto; el caller decide si ya la escribió.

## Alcance

- `.matecito-ai/*/INDEX.md` — los índices raíz de los stores durables, que enrutan por agrupador.
- `.matecito-ai/*/*/INDEX.md` — los índices por agrupador, que sí llevan una fila por artefacto.

## Reglas verificables

- **[auto]** El índice raíz de un store no repite un agrupador en más de una fila.
- **[manual]** El renderizador no recibe ni consulta el lote completo de una materialización: cada invocación conoce sólo su artefacto.
- **[manual]** El índice del store se actualiza una sola vez al cerrar la tanda, nunca una vez por artefacto.

## Alternativas consideradas

- **Que el renderizador deduplique.** Descartado: lo obligaría a conocer el lote completo, y dejaría de ser una función de un artefacto para pasar a ser una de una tanda. Esa dependencia se paga en todos los llamadores, incluidos los que materializan uno solo.
- **Un índice raíz con una fila por artefacto.** Descartado: convierte el índice de enrutamiento en un listado plano que duplica lo que ya dice el índice del agrupador, y crece hasta dejar de servir para lo único que existe — decidir dónde mirar.

## Consecuencias

- El caller debe agrupar antes de escribir. Es trabajo real, y es donde aparece el error si alguien materializa artefacto por artefacto sin agrupar.
- A cambio, el renderizador se mantiene como una función de un solo artefacto, invocable desde cualquier productor sin contexto extra.
- El índice raíz se mantiene corto y legible: su tamaño crece con la cantidad de agrupadores, no con la de artefactos.

## Relacionados

- `relacionado-con` → [contract-pair-in-templates.md](contract-pair-in-templates.md) — el contrato declara la forma de ambas filas.
- `relacionado-con` → [../contracts/data-contract-derived-and-producer-neutral.md](../contracts/data-contract-derived-and-producer-neutral.md) — por qué la fila la calcula el renderizador y no la suministra el productor.
