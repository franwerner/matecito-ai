# EDR — La entrada del renderizador es agnóstica del productor, y lo derivable no se suministra

- **Status:** Accepted
- **Date:** 2026-08-06

## Contexto

El mismo artefacto durable lo puede producir una fase del flujo, una minería sobre material existente, o una persona escribiendo a mano. Los tres tienen contexto distinto y ninguno es más legítimo que los otros. Si cada uno alimenta al renderizador con una forma propia, el artefacto deja de tener un contrato y pasa a tener tres.

Hay además una clase de dato que el renderizador puede calcular a partir del resto — la ubicación de un artefacto a partir de sus coordenadas, la fila que le corresponde en un índice. Permitir que el productor lo suministre abre la posibilidad de que dos fuentes de verdad se contradigan, y de que la contradicción se materialice sin que nada la note.

## Decisión

La entrada del renderizador es **una sola, agnóstica de quién la produce**, y **lo derivable no se suministra: se calcula**. Un campo que el contrato marca como derivado no se acepta desde afuera.

## Reglas verificables

- **[auto]** El esquema de entrada declara explícitamente qué campos son derivados y no deben suministrarse.
- **[auto]** El renderizador aborta ante un campo que no sabe manejar, en vez de descartarlo; un dato que llega y no aparece en la salida es un error, nunca un silencio.
- **[manual]** No existe una forma de entrada por productor: el flujo, la minería y la escritura manual usan el mismo contrato.
- **[auto]** El esquema publicado enumera todos los campos que alguna invocación exige, incluidos los que sólo consume una salida secundaria.

## Alternativas consideradas

- **Una forma de entrada por productor.** Descartado: multiplica contratos que describen el mismo artefacto, y cada uno diverge con el primer cambio que se aplique de un solo lado.
- **Permitir suministrar los campos derivados.** Descartado: habilita que el valor suministrado y el calculado difieran, y el artefacto materializado no dice cuál ganó.
- **Descartar en silencio lo que el renderizador no entiende.** Descartado explícitamente: es la forma más barata de perder contenido, porque no deja rastro. Un campo desconocido aborta.

## Consecuencias

- Un productor nuevo funciona sin tocar el renderizador: si arma la entrada, renderiza.
- El contrato no puede expresar excepciones puntuales. Un caso que no encaja se resuelve extendiendo el contrato para todos, no con una vía lateral — más caro, y deliberadamente.
- Abortar ante lo desconocido hace ruidosa cualquier extensión no declarada. Es el costo de que nada se pierda callado.

## Relacionados

- `relacionado-con` → [../structure/contract-pair-in-templates.md](../structure/contract-pair-in-templates.md) — dónde vive el contrato que declara los derivados.
- `relacionado-con` → [../structure/root-index-cardinality-per-domain-type.md](../structure/root-index-cardinality-per-domain-type.md) — la fila de índice es uno de los valores que el renderizador calcula.
