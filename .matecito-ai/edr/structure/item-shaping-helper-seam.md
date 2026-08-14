# EDR — Un solo helper compartido produce el adorno de un ítem, en tabla o en lista

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Sólo las secciones de lista con ítems sabían imprimir los tokens y la razón de un ítem. Extender esa misma forma a una sección de tabla —cuyas filas ya tienen columnas propias y no admiten líneas de continuación— repetiría la lógica de adorno en cada renderizador que la necesitara.

## Decisión
Un único helper compartido devuelve las líneas de adorno de un ítem (tokens, campos, razón); cada renderizador conserva su propia primera línea (la fila de tabla, o el resumen de la lista) y delega el resto en el helper. Una fila de tabla no lleva líneas de continuación propias: su detalle de adorno se imprime aparte, como una lista debajo de la tabla y su pie, con cada bloque identificado por una clave que el ítem declara.

## Reglas verificables
- **[auto]** Todo renderizador que emite adorno de ítem (tokens, campos, razón) llama al mismo helper compartido, en vez de reimplementar la lógica.
- **[manual]** Una sección de tabla no imprime líneas de continuación dentro de la fila; el detalle de adorno va en un bloque aparte, debajo de la tabla, indexado por la clave del ítem.

## Alternativas consideradas
Duplicar la lógica de adorno en cada renderizador — descartada porque cada copia es un lugar donde el contrato puede desalinearse del original. Forzar líneas de continuación dentro de la fila de una tabla — descartada porque una fila de tabla ya tiene su propia forma columnar, y mezclarle líneas de continuación rompe esa forma.

## Consecuencias
Toda sección que declare ítems —tabla o lista— hereda el mismo comportamiento de adorno por construcción, sin importar cuántas se agreguen después — a costa de que una tabla con ítems ahora imprime dos bloques (la tabla y el detalle), no uno solo.

## Relacionados
- `relacionado-con` → [validate-side-needs-no-mirror.md](validate-side-needs-no-mirror.md) — por qué el lado de validación no necesita un helper espejo de este.
- `relacionado-con` → [../contracts/subverifier-item-shape-single-declaration.md](../contracts/subverifier-item-shape-single-declaration.md) — otra forma compartida declarada una sola vez, en otro punto del mismo pipeline.
