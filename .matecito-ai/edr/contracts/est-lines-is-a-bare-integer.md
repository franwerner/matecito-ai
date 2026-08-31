# EDR — `Est. lines` es un entero simple por Phase, nunca un rango

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
`### Breakdown`, en el retorno de `sdd-tasks`, gana una columna `Est. lines` por Phase. Ese número alimenta al futuro guard `Apply Dispatch Budget` (materializado aparte, en un commit posterior de este mismo cambio), que lo acumula aritméticamente para decidir dónde cortar un despacho serial de `sdd-apply`. `### Review Workload Forecast` ya tiene una línea agregada, `Estimated changed lines`, que tolera "estimate or range" porque la lee una persona.

## Decisión
`Est. lines` es un único número entero no negativo por Phase — sin `~`, sin `200-300`, sin unidades — sobre la lista de tareas completa de esa Phase. Es una divergencia deliberada respecto de la línea agregada del forecast, que conserva su libertad de rango.

## Reglas verificables
- **[auto]** `sdd-tasks.yaml` declara `est_lines` como una columna de tabla simple, sin envoltorio derivado ni de formato — una fila sin ese valor hace fallar el render nombrando la fila y la clave faltante.
- **[manual]** La SKILL de `sdd-tasks` establece la regla "entero simple, nunca un rango" y el motivo declarado de la divergencia respecto de la línea agregada del forecast.
- **[manual]** El guard Apply Dispatch Budget degrada en silencio — un solo despacho, sin error — cuando el `Est. lines` de cualquier Phase está ausente o no se puede parsear, en vez de intentar interpretar un rango.

## Alternativas consideradas
Aceptar un rango y tomar su cota superior — descartado: agrega una regla de parseo a un guard de cinco pasos y sesga hacia cortar de más. Aceptar un rango y dejar que el guard degrade — descartado: apaga el corte justo donde la estimación es menos certera, que es el caso donde más se lo necesita.

## Consecuencias
La regla de degradación del guard ("ausente o no parseable → un despacho, sin error") se mantiene simple: un rango degrada exactamente igual que un valor ausente, los dos son no-ops silenciosos. `sdd-tasks` ahora escribe dos tipos de estimación distintos para la misma magnitud — una tolerante a rangos (la agregada) y una estricta (la de Est. lines por Phase) — y no debe confundirlas.
