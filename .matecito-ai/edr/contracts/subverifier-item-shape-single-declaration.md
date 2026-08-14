# EDR — La forma de un hallazgo se declara una sola vez, no en cada prompt de sub-verificador

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
La fase de verificación se despacha como varias instancias concurrentes, una por grupo de chequeos, y cada una devuelve sus propios hallazgos para consolidar en un solo reporte. Si la forma de un hallazgo (los campos que lleva, el ancla, el resumen) se declarara en el prompt de cada sub-verificador por separado, varias copias del mismo contrato tendrían que mantenerse alineadas a mano.

## Decisión
La forma de un hallazgo se declara una sola vez, en el documento que describe el sobre de reporte que todo sub-verificador ya lee antes de correr — no en cada prompt de despacho por separado.

## Reglas verificables
- **[auto]** Ningún prompt de despacho de un sub-verificador redeclara la forma de un hallazgo; todos la heredan del documento único del sobre de reporte.
- **[manual]** Un cambio a la forma de un hallazgo se edita en un solo lugar y alcanza a todos los sub-verificadores sin tocar sus prompts.

## Alternativas consideradas
Declarar la forma en cada uno de los prompts de despacho por separado — descartada porque varias copias del mismo contrato son varios lugares donde puede desalinearse, y ese desvío no se nota hasta que un sub-verificador devuelve una forma distinta a las demás.

## Consecuencias
Agregar o cambiar un campo del hallazgo es una sola edición que alcanza a todos los sub-verificadores por construcción — a costa de que el documento del sobre de reporte se vuelve una dependencia que todo sub-verificador tiene que leer antes de empezar.

## Relacionados
- `relacionado-con` → [../structure/item-shaping-helper-seam.md](../structure/item-shaping-helper-seam.md) — misma lógica de un solo lugar declarando una forma compartida, en otro punto del mismo pipeline.
