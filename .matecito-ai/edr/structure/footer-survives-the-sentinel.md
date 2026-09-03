# EDR — El footer de una tabla sobrevive a su sentinel, por opt-in

- **Status:** Accepted
- **Date:** 2026-09-03

## Contexto
`renderTable` devuelve `'None.'` en cuanto la lista de filas está vacía y la sección declara `sentinel`, antes de leer `section.footer` en ningún punto. Una sección cuyo footer registra evidencia de un chequeo DISTINTO al de las filas — no un agregado de ellas — pierde esa línea exactamente cuando la tabla queda vacía, que es el caso en que "limpio" y "nunca corrió" más necesitan distinguirse. El mismo camino también evita el fallo por campo faltante que un footer obligatorio exige.

## Decisión
Un footer puede declarar `on_sentinel: true` para sobrevivir al retorno temprano `'None.'` de una tabla vacía; sin ese flag, el retorno sigue siendo el `'None.'` desnudo de siempre. El flag es por footer, nunca global: sólo lo declara la sección cuyo footer registra un chequeo de alcance store, no un resumen de filas.

## Reglas verificables
- **[auto]** `renderTable` con `rows: []`, `sentinel: true` y un footer con `on_sentinel: true` (más el campo del footer presente en los datos) devuelve `'None.'` seguido de la línea del footer.
- **[auto]** La misma llamada sin `on_sentinel` sigue devolviendo el `'None.'` desnudo, sin ninguna línea de footer.
- **[auto]** La misma llamada con `on_sentinel: true` y el campo del footer ausente en los datos falla nombrando el campo faltante.
- **[auto]** Los tres footers preexistentes (`compliance_summary`, `breakdown.totals`, `error_gate`) siguen renderizando byte a byte igual que antes de este cambio — ninguno declara `on_sentinel`.

## Alternativas consideradas
Aceptar que una tabla vacía no lleve footer: descartada, viola dos requisitos de la nueva capacidad (la línea se emite siempre que la sección se emite; un fragmento sin la clave falla el render). Sacarle `sentinel: true` a la sección: descartada, una tabla vacía pasaría a renderizar un header pelado en vez de `None.`, reformando una sección fuera del alcance ratificado de este cambio. Que todo footer de tabla sobreviva a su sentinel: descartada, vuelve obligatorios en el caso vacío a `compliance_summary` y `breakdown.totals`, dos secciones que hoy no cargan nada cuando no hay filas — un cambio de exigencia que este cambio no pidió.

## Consecuencias
Unas ~6 líneas de lógica de render más, y una clave más en el vocabulario del contrato (`footer.on_sentinel`). Los tres footers preexistentes quedan byte-idénticos, cubierto por prueba de regresión. La línea de alcance store sobrevive al caso `None.`, que es exactamente donde "limpio" y "nunca corrió" se confundían.

## Relacionados
- `relacionado-con` → [../contracts/data-contract-derived-and-producer-neutral.md](../contracts/data-contract-derived-and-producer-neutral.md) — el footer que este EDR habilita a sobrevivir es el mismo cuya forma fija `contracts/store-wide-summary-slot-is-a-plain-string.md`
