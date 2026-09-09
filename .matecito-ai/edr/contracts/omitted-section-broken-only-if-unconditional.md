# EDR — Una sección ausente sólo es un retorno roto cuando esa sección era incondicional

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La regla anterior trataba como retorno roto TODA sección ausente, y hay reglas vigentes que ordenan
omitir secciones enteras de forma legítima y condicional (conflictos de EDR con el store inactivo, el
veredicto de UI cuando no aplica, la evidencia de TDD fuera de estricto): las tres se marcaban como
error. Además exigía "pedile que re-emita", una obligación sin mecanismo — no existe comando ni status
de re-emisión, y re-despachar re-ejecuta la fase entera (en `sdd-apply` un re-despacho está definido
como batch de continuación, no como re-emisión).

## Decisión
Sólo las secciones declaradas **incondicionales** están rotas cuando faltan. Una sección que el propio
skill de la fase declara **condicional**, y cuya condición no se cumple, está legítimamente ausente — sin
gate, sin mención. Cuando falta una sección incondicional, el guard no exige "re-emisión" (que no
existe): trata la omisión como contenido gateable sin resolver y abre el mismo gate que abriría para
contenido real.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Unresolved Decisions Guard (MANDATORY)`, distingue secciones incondicionales (rotas si faltan) de condicionales (legítimamente ausentes) y no exige una "re-emisión" inexistente.

## Alternativas consideradas
Seguir tratando cualquier sección ausente como retorno roto — descartada: marca como error tres casos
legítimos de omisión ya vigentes en el ecosistema. Exigir una "re-emisión" del retorno — descartada: no
existe ese comando ni ese status; lo que existe es re-despachar la fase entera.

## Consecuencias
Una sección ausente legítimamente (condición no cumplida) no genera ruido; una sección incondicional
ausente sí abre gate, con las opciones reales que el ecosistema ofrece (tratar como vacío, re-despachar,
ajustar).
