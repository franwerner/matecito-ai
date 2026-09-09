# EDR — Los flags de decisión del brief se declaran donde el orquestador los lee, no sólo en el agente de intake

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Con las dos lanes fijas, ningún gate ratifica ya los flags de decisión del brief (`diagram`, `ui-test`,
`components`, `worktree-isolation`): `sdd-intake` los decide y simplemente los reporta. Pero la
instrucción de que el orquestador tiene que surfacear cada uno en una línea de aviso, y de que ningún
gate los confirma como ítem propio, sigue haciendo falta en un lugar que el orquestador realmente lee —
y `agents/sdd-intake.md` y su skill son leídos por el ejecutor de intake, nunca por el orquestador.

## Decisión
La tabla de flags y la instrucción de reporte viven en el fragmento del dominio
(`payload/domains/development/CLAUDE.md`), que el orquestador sí lee, no sólo en el agente de intake que
los decide. Esto no cambió con el pase a lanes fijas: la tabla ya vivía ahí antes de ese cambio.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Brief decision flags (decided by
  intake, confirmed with the brief)`, declara la tabla de flags y la instrucción de reporte en un
  documento que el orquestador lee, no en `agents/sdd-intake.md`.

## Alternativas consideradas
Dejar la tabla sólo en `agents/sdd-intake.md` y su skill, ya que es intake quien decide los flags —
descartada: el orquestador es quien tiene que reportar cada flag en una línea de aviso, y no lee ese
archivo.

## Consecuencias
Un cambio a un flag de decisión del brief tiene que actualizar la tabla acá, no en el agente de intake,
para que el orquestador vea el cambio.
