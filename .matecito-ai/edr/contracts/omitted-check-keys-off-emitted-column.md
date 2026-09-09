# EDR — El chequeo de sección omitida se define por la columna `Emitted`, no por pertenecer a la tabla D.3

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La frase anterior decía "las [secciones] enumeradas en D.3", y era cierta hasta que la Sección D.3 sumó
una fila CONDICIONAL (`## Decision Gaps`, sólo cuando el cambio materializó al menos un decision record).
Leída al pie, esa frase convertía la ausencia normal de esa sección condicional en un retorno roto — un
gate espurio en cada verify.

## Decisión
El corte que decide si una sección ausente es un retorno roto es la columna `Emitted` de la Sección D.3
— si esa fila está marcada `always` — nunca la mera pertenencia de la sección a la tabla.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Unresolved Decisions Guard (MANDATORY)`, ata "sección incondicional" a la columna `Emitted` de la Sección D.3, no a la pertenencia a la tabla.

## Alternativas consideradas
Seguir leyendo "pertenece a la tabla D.3" como sinónimo de "incondicional" — descartada: deja de ser
cierto en cuanto la tabla gana una fila condicional, y produce un gate espurio en cada verify sobre una
ausencia legítima.

## Consecuencias
Agregar una fila condicional nueva a la Sección D.3 no rompe este chequeo — el chequeo ya lee la columna
que distingue el caso, no la mera presencia en la tabla.
