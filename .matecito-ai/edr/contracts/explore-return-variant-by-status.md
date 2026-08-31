# EDR — El retorno de sdd-explore expresa sus dos pasadas como variants por status, no como secciones on-status

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
sdd-explore pasó a tener dos pasadas con contenido estructuralmente distinto: un Discovery Form (Pass 1, status needs-input) y una Exploration (Pass 2, status done/blocked). render-return.js y validate-return.js resuelven el contrato de una fase por status. Declarar las secciones exclusivas de Pass 1 (Approaches, Recommendation, Ready for Proposal) como emitted: on-status dentro de un único bloque las deja estructuralmente expresables en un retorno Pass-1, aunque ese contenido no puede existir todavía — depende de código que recién se lee en Pass 2.

## Decisión
sdd-explore.yaml declara un `variants:` de nivel superior, con una entrada por grupo de statuses — `[needs-input]` para el Discovery Form y `[done, blocked]` para la Exploration — igual que ya hace sdd-intake.yaml. `resolveVariant` (render-return.js, validate-return.js) ya es genérico y no hardcodea el nombre de ninguna fase, así que el mecanismo se reutiliza sin tocar ningún script.

## Reglas verificables
- **[auto]** `render-return.js --phase sdd-explore --data <payload Pass-1>` no puede emitir `### Approaches`, `### Recommendation` ni `### Ready for Proposal` — la variante `[needs-input]` no las declara en sus `sections:`, así que ningún dato puede producirlas.
- **[auto]** `validate-return.js --self-check --phase sdd-explore` falla si alguna de las dos variantes (bloques o secciones declaradas) diverge de lo que `sdd-explore.md` muestra en su template.

## Alternativas consideradas
Mantener un único bloque marcando las secciones exclusivas de Pass 1 como `emitted: on-status, statuses: [needs-input]` — descartado porque deja Approaches/Recommendation/Ready for Proposal expresables en un retorno que todavía no leyó código, exactamente lo que `variants` hace estructuralmente imposible.

## Consecuencias
El contrato de sdd-explore queda simétrico al que sdd-intake tenía antes de su propio colapso a un solo bloque: dos bloques distintos, cada uno con su propio conjunto de secciones. Dos comentarios en render-return.js y validate-return.js que nombraban a sdd-intake como la única fase con variants ahora describen dos fases; son comentarios, no comportamiento, y no se tocan si eso implicara editar el script.
