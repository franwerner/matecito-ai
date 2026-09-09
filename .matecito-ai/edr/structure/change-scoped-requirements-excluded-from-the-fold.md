# EDR — El pliegue del delta de comportamiento usa un manifiesto cerrado en el diseño, no un scan de encabezados del spec

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La instrucción del Step 5b de sdd-apply es "para cada capacidad que toca, plegar el delta". El spec artifact de un cambio puede cargar más cuerpos de requisitos que capacidades reales — este cambio trae tres, y sólo dos son capacidad-spec. Escanear el artefacto por encabezados `## ADDED Requirements` encuentra los tres por igual, y plegar el tercero (un bloque de requisitos change-scoped: chequeo de tope de palabras, regla de supervivencia de un claim, las cinco duplicaciones nombradas, el barrido de clases de comentario, el alcance de una reubicación) dentro de `.matecito-ai/development-specs/` es exactamente la falla concreta que este cambio diseñó para evitar.

## Decisión
La lista de documentos de comportamiento que el paso que implementa puede actualizar al cierre de un cambio se escribe entera en el diseño, como un manifiesto de pliegue cerrado que nombra las dos capacidades y su ruta destino (fila YES) más cualquier cuerpo de requisitos que NO se pliega, sin ruta destino (fila NO) — restablecido como el criterio de aceptación de la tarea que cierra el cambio. El Step 5b de sdd-apply lee ese manifiesto en vez de escanear los encabezados `## ADDED Requirements` del artefacto de spec.

## Reglas verificables
- **[manual]** El diseño de un cambio cuyo spec incluye requisitos change-scoped no plegados declara un manifiesto cerrado con las filas YES (capacidad + ruta destino) y NO (sin ruta destino).
- **[manual]** La tarea que cierra el cambio restablece ese manifiesto como su propio criterio de aceptación.
- **[manual]** El Step 5b de sdd-apply lee el manifiesto del diseño en vez de escanear los encabezados `## ADDED Requirements` del artefacto de spec para decidir qué capacidad-spec tocar.

## Alternativas consideradas
Confiar en el propio encabezado "NOT folded" del artefacto de spec — descartada: es el estado que ya fallaba, porque una fase headless que escanea encabezados de requisito encuentra las tres secciones por igual, sin distinguir cuál es change-scoped. Agregar un token `fold-scope:` al contrato del artefacto de spec — descartada: edita los contratos de payload de `sdd-spec` y `sdd-apply`, explícitamente fuera de alcance de este cambio.

## Consecuencias
sdd-apply lee dos canales independientes que ya lee de todos modos (tasks + spec + design), sin requerir un cambio de contrato de payload. El manifiesto queda cerrado: un requisito change-scoped que no figura en la lista YES no se pliega, sea cual sea su forma o cuán completo esté escrito.
