# EDR — El footer de una tabla puede llevar dos totales en una sola línea formateada

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
`### Breakdown` en `sdd-tasks` pasa de `render: labeled-bullets` a `render: table`, pero la sección lleva DOS totales que hoy son dos bullets separados (`Total tasks`, `Work units`). `render-return.js` sólo admite un `footer` por sección de tabla, y agregarle una segunda entrada de footer es un cambio de script — prohibido por el alcance ratificado de este cambio. El total de tareas que el footer muestra lo escribe `sdd-tasks`, no lo suma de la columna `Tasks` de la tabla — lo cual tensiona con `contracts/data-contract-derived-and-producer-neutral` sin contradecirlo: el esquema no marca ningún campo como derivado, y el renderizador sigue abortando ante un campo desconocido.

## Decisión
El footer lleva los dos totales en una sola línea formateada, suministrados juntos como un solo objeto: `footer: { label: Totals, field: breakdown.totals, format: "{tasks} tasks · {work_units} work units" }`, que renderiza `**Totals**: 23 tasks · 4 work units`.

## Reglas verificables
- **[auto]** `sdd-tasks.yaml`'s `### Breakdown` section declares exactly one `footer`, whose `format` substitutes both `tasks` and `work_units` from the single object field `breakdown.totals`.
- **[manual]** The task total in the footer is authored by `sdd-tasks`, not derived from summing the table's own `Tasks` column — nothing currently detects a disagreement between the two.
- **[auto]** `render-return.js` and `validate-return.js` stay byte-identical for this change — no second `footer` key was taught to the renderer.

## Alternativas consideradas
(a) Move `Work units` into its own table row — rejected: the table is keyed per Phase, and a totals row is not a Phase. (b) Teach `render-return.js` a second `footer` entry — rejected: a script change, forbidden by this change's ratified scope. (c) Teach `derive()` a `sum:<field>.<column>` kind so the task total is computed from the rows — the alternative that best honors `contracts/data-contract-derived-and-producer-neutral`, set aside only because the ratified scope forbids script changes.

## Consecuencias
Both totals survive in one footer line, with no script change. The task total can silently disagree with the sum of the `Tasks` column — an accepted cost, stated plainly here, not mitigated. The tension against `contracts/data-contract-derived-and-producer-neutral` was raised at design time as `contested: contradicts-record` and never resolved unilaterally: the record's own verifiable rules all still hold (the schema marks no field derived, the renderer aborts on the unknown), so the record is not contradicted, only in tension.

## Relacionados
- `relacionado-con` → [data-contract-derived-and-producer-neutral.md](data-contract-derived-and-producer-neutral.md) — El campo en tensión: el total de tareas es suministrado, no derivado, por la razón que este EDR documenta.
