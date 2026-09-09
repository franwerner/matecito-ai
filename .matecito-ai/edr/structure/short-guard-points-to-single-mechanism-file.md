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

## Decisión
Ambos guards declaran que son MANDATORY, resumen su disparo mínimo, y apuntan a `parallel-batch.md` como
la única fuente del mecanismo completo — el fragmento no mantiene una segunda copia de ninguno de los
dos.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Parallel-Mark Validation (MANDATORY)`, valida las marcas del artefacto de tareas por script, nunca a ojo, y apunta a `parallel-batch.md` para el mecanismo completo.
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Uncommitted-Work Gate (MANDATORY)`, apunta a `parallel-batch.md` para el mecanismo completo sin restatearlo en el fragmento.

## Alternativas consideradas
Documentar el mecanismo completo de cada guard en el propio fragmento, además de en `parallel-batch.md`
— descartada: produce dos copias que divergen en la próxima edición, el mismo defecto que otras
relocaciones de este cambio ya corrigieron.

## Consecuencias
Editar el mecanismo de cualquiera de los dos guards se hace una sola vez, en `parallel-batch.md`; el
fragmento hereda el cambio por referencia.
