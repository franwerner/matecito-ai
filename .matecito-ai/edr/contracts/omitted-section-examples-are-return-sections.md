# EDR — Los ejemplos de sección legítimamente ausente son secciones de un retorno, nunca de un artefacto

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La lista de ejemplos de sección condicional legítimamente ausente, dentro de "Omitted → depends on
whether the section was unconditional" (`### Unresolved Decisions Guard`), usaba antes `## EDR
Conflicts` como primer ejemplo — una sección que sólo existe en el artefacto de DISEÑO, nunca en el
retorno de fase de `sdd-design`. Es exactamente la confusión artefacto/retorno que esta parte del
ecosistema existe para cerrar.

## Decisión
Los tres ejemplos que ilustran una sección legítimamente ausente son, los tres, secciones de un
RETORNO de fase — `### Early guard (EDRs)` de `sdd-intake`, el veredicto de UI, la tabla de evidencia
TDD — nunca una sección que sólo vive en un artefacto. `## EDR Conflicts` vive en el artefacto de
diseño, no en el retorno de `sdd-design`, así que usarla como ejemplo de "sección de retorno
legítimamente ausente" perpetuaba la confusión que la propia guía busca evitar.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Unresolved Decisions Guard
  (MANDATORY)`, subsección "Omitted → depends on whether the section was unconditional" — los tres
  ejemplos citados son secciones de un retorno de fase, nunca de un artefacto.

## Alternativas consideradas
Dejar `## EDR Conflicts` como ejemplo — descartada: es una sección del artefacto de diseño, no del
retorno de `sdd-design`, y usarla como ejemplo de "sección de retorno ausente" perpetúa la confusión
artefacto/retorno que esta misma guía distingue en otro lado.

## Consecuencias
Un lector que busca un ejemplo de sección condicional legítimamente ausente encuentra sólo ejemplos
reales de secciones de retorno — ninguno de artefacto.
