# EDR — El Return Contract Check corre un validador mecánico, nunca comparación a ojo

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Nada comprobaba que el retorno de una fase trajera lo que el template de esa fase exige. Si una fase se
comía una sección, el gate correspondiente no disparaba — en silencio, que es el modo de falla que más
cuesta detectar. Los cuatro chequeos de este guard vivían sólo como prosa, y el motor que los ejecuta
(`validate-return.js`) existía sin que ninguna regla lo invocara — su propio código ya se autodescribe
como el que hace que este check "stop being a human reading and comparing titles by eye", pero su nombre
no aparecía en ningún `.md` del dominio. Comparar títulos a ojo es exactamente el modo de falla que este
guard vino a cerrar, así que dejarlo en manos del ojo era el defecto reproduciéndose un nivel más arriba.

## Decisión
El Return Contract Check valida el retorno de toda fase contra el template canónico en
`~/.claude/references/phase-returns/<phase>/<phase>.md` (la especificación) corriendo
`node ~/.claude/scripts/validate-return.js --phase <phase> --status <status> --file <return.md>` primero,
siempre — nunca leyendo y comparando títulos a ojo. Un exit distinto de 0 se interpreta por código; un
exit 2 (el validador no pudo correr) se reporta explícitamente y nunca se degrada en silencio a lectura
manual.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Return Contract Check (MANDATORY)`,
  declara que la validación corre por `validate-return.js` antes de cualquier otra cosa, y que un exit 2
  se surface en vez de caer a lectura manual.
- **[manual]** La misma sección declara que los cuatro chequeos que siguen son la especificación de lo
  que el motor enforce — se leen para interpretar un finding, no para performar la comparación a mano.

## Alternativas consideradas
Dejar los cuatro chequeos como prosa sin motor que los corra, confiando en que quien despacha la fase los
compara a ojo — descartada: es exactamente el modo de falla (comparación manual de títulos) que el guard
existe para cerrar, reproducido un nivel más arriba en quien lo opera.

## Consecuencias
Todo retorno de fase pasa por el script antes de cualquier guard o dispatch. El costo es un paso mecánico
más antes de actuar sobre un retorno; el beneficio es que una sección perdida o un título re-nivelado
deja de ser invisible.
