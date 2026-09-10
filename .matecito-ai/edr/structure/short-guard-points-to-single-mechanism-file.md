# EDR — Un guard corto declara MANDATORY y apunta; el mecanismo completo vive una sola vez en el archivo referenciado

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
`### Parallel-Mark Validation` y `### Uncommitted-Work Gate` son dos guards cortos con la misma forma:
las marcas del artefacto de tareas vivían "como se emiten", sin ningún chequeo mecánico, un nivel por
debajo del Return Contract Check — mismo principio, un artefacto más abajo: un script gatea, nunca a
ojo. El mecanismo completo de cada uno (cuándo corre, el mapeo de componentes y su fallback, las tres
salidas, el caso huérfano, la forma del aviso) vive una sola vez en `parallel-batch.md`, nunca duplicado
en el fragmento.

**Extensión (port-turn-mechanism):** `contracts/uncommitted-gate-follows-the-container` se retiró en vez
de reescribirse — su única decisión era la condición sobre cuál contenedor inspeccionaba el gate, y con
el workspace de nivel de cambio retirado esa condición ya no existe. La oración que sobrevive de ese
registro —el gate inspecciona el repositorio principal, siempre, y su salida de degradación corre sobre
la rama de trabajo— no gana un registro propio: se plegó directamente en la sección `### Uncommitted-Work
Gate` de `parallel-batch.md`, exactamente donde este registro ya fija que vive el mecanismo completo. El
guard corto del fragmento (`payload/domains/development/CLAUDE.md`) no cita ningún registro para esa
afirmación — sigue apuntando sólo a `parallel-batch.md`, sin restatearla.

## Decisión
Ambos guards declaran que son MANDATORY, resumen su disparo mínimo, y apuntan a `parallel-batch.md` como
la única fuente del mecanismo completo — el fragmento no mantiene una segunda copia de ninguno de los
dos. El enunciado de contenedor único del gate de trabajo sin commitear entra en esta misma regla: vive
en `parallel-batch.md`, sin un registro EDR detrás que lo respalde, porque ya no hay una condición que
decidir.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Parallel-Mark Validation (MANDATORY)`, valida las marcas del artefacto de tareas por script, nunca a ojo, y apunta a `parallel-batch.md` para el mecanismo completo.
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Uncommitted-Work Gate (MANDATORY)`, apunta a `parallel-batch.md` para el mecanismo completo sin restatearlo en el fragmento, y no cita ningún registro EDR para la afirmación de contenedor único.
- **[manual]** `parallel-batch.md`, sección `## Uncommitted-Work Gate`, enuncia que el gate inspecciona el repositorio principal, siempre, y que su salida de degradación corre sobre la rama de trabajo, sin ninguna forma condicional sobre cuál contenedor inspecciona.

## Alternativas consideradas
Documentar el mecanismo completo de cada guard en el propio fragmento, además de en `parallel-batch.md`
— descartada: produce dos copias que divergen en la próxima edición, el mismo defecto que otras
relocaciones de este cambio ya corrigieron.

Reescribir `contracts/uncommitted-gate-follows-the-container` en vez de retirarlo, para conservar el
rastro de por qué se preguntó cuál era el contenedor — descartada en `sdd-design` (contra la
recomendación de la proposal): ese rastro sobrevive en el historial de git y en el propio artefacto de
diseño, y ninguno de los dos es razón para mantener un registro Accepted cuya decisión ya es el default.

## Consecuencias
Editar el mecanismo de cualquiera de los dos guards se hace una sola vez, en `parallel-batch.md`; el
fragmento hereda el cambio por referencia. Lo mismo vale ahora para el enunciado de contenedor único del
gate: se edita una sola vez, en `parallel-batch.md`, y ningún registro EDR necesita tocarse cuando cambie
su redacción.
