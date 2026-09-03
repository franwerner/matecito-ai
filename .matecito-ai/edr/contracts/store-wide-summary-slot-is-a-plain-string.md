# EDR — El resumen de alcance store es un string plano, no un objeto con formato

- **Status:** Accepted
- **Date:** 2026-09-03

## Contexto
La línea de alcance store — evidencia de que un chequeo corrió sobre el store entero, no un resumen de las filas de la tabla — necesita un lugar propio en el contrato de retorno de `### Coherence (Capability-Specs)`. Su hermana `### Spec Compliance Matrix` ya tiene un footer parecido (`compliance_summary`), pero ese footer agrega las filas de arriba; éste no: expresa tres estados distintos (corrió con hallazgos, corrió limpio, no corrió), y uno de ellos — "no corrió" — no tiene un valor numérico honesto que representarlo sin generar la misma confusión que el cambio busca cerrar.

## Decisión
El lugar es `field: spec_store_structure`, `label: "Store structure (pre-existing)"`, sin `format`: un string plano. El renderer garantiza que el slot existe y es obligatorio; el contenido interno de ese string queda fijado en prosa (la skill de `sdd-verify`, paso 6e), no en la forma del contrato.

## Reglas verificables
- **[auto]** El schema de `sdd-verify` publica `spec_store_structure: string ("Store structure (pre-existing)")` bajo `### Coherence (Capability-Specs)`.
- **[manual]** El valor de `spec_store_structure` es siempre uno de tres estados en prosa: corrió con hallazgos pre-existentes, corrió limpio, o no corrió — nunca un objeto con conteos por severidad.
- **[manual]** Los hallazgos que este cambio SÍ tocó nunca viajan en este slot: van a `### Issues Found`. El slot sólo carga lo pre-existente.

## Alternativas consideradas
`field: coherence_specs_summary`, imitando el nombre de la hermana `compliance_summary`: descartada, nombra el slot como un resumen de las filas, que no es lo que es — es evidencia de que corrió un chequeo distinto sobre todo el store. `format:` con conteos de severidad, imitando el objeto de dos valores de `compliance_summary`: descartada, no puede expresar "no corrió" sin imprimir ceros, que es exactamente la confusión que este cambio cierra.

## Consecuencias
El nombre se alinea con las otras dos claves ya en juego para esta sección: el booleano de gate `spec_store_active` y las filas `coherence_specs`. El valor viene de la corrida de una herramienta externa (`validate-artifact.js`), así que lo suministra su productor en vez de derivarlo el renderer — y `schema()` ya publica esa obligación en su rama de footer existente.

## Relacionados
- `relacionado-con` → [../structure/footer-survives-the-sentinel.md](../structure/footer-survives-the-sentinel.md) — el mismo footer necesita sobrevivir al caso `None.` de la tabla — dos decisiones sobre el mismo slot
