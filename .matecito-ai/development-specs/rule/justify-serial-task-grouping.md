# Capability — Justify serial task grouping in phase breakdown

- **Status:** Accepted
- **Date:** 2026-08-31
- **Components:** cli

## Propósito

Dejar los tasks de una Fase sin marca de ejecución concurrente solía ser libre y silencioso — el camino serial, exactamente el que paga el crecimiento, era el default que nadie tenía que defender. La fase de descomposición ahora adeuda un veredicto explicado por Fase: reportado, nunca gatillando.

## Actores

- **Fase de tareas**: redacta un `### Parallelization Verdict` con un entry por Fase; cada entry declara el veredicto (marcado o serial) y su razonamiento
- **Orquestador**: imprime el summary de cada entry en el resumen between-phase; nunca como un gate que detiene
- **Validador**: rechaza entries incompletas y detecta drift entre contrato y template

## Precondiciones

- El retorno de `sdd-tasks` declara una sección `### Parallelization Verdict` con tabla, `gates: reported`
- Cada entry en la sección lleva tres partes: `summary` (máximo 250 caracteres), `anchor` y `rationale`
- Hay exactamente un entry por Fase de `### Breakdown` — ni más ni menos
- La sección no declara token `contested` (contested es inerte en una sección `reported`)

## Flujo principal

1. La fase de tareas evalúa la independencia de cada Fase según su criterio de elegibilidad actual (conservador; por defecto serial)
2. La fase redacta un entry `### Parallelization Verdict` por Fase con su veredicto y razonamiento
3. Si la Fase fue marcada para concurrencia, el `summary` lo dice y la columna `Group` lleva el id del grupo
4. Si la Fase se mantiene serial, el `summary` lo dice con el razonamiento por el cual sus tareas no son independientes
5. La fase se renduriza con exactamente un entry por Fase de `### Breakdown`, llevando las tres partes + token `anchor`
6. El orquestador imprime solo `summary` + `anchor` en el resumen entre-fases; la `rationale` completa vive en el bloque persistido del `detailed_report`
7. **Nunca abre un gate** — la sección `gates: reported` significa que los entries siempre surfacean, nunca frenan el flujo

## Ramas / flujos alternativos

- **Fase serializada que contiene un par genuinamente independiente** → su entry serial nombra ese par explícitamente en lugar de generalizar
- **Sección vacía (blocked return, no Phases)** → se emite el sentinel `None.`
- **Sección sin items** → comportamiento sin cambios — el sentinel se emite

## Casos borde

- **Una Fase con todas las tareas marcadas** → entry nombra el grupo y declara "parallelized"
- **Una Fase totalmente serial con razonamiento simple** → entry nombra una restricción que aplica a todas las tareas
- **Una Fase con dos subtareas genuinamente independientes en un contexto serial general** → entry nombra esas dos tareas explícitamente por id, declarando por qué el resto está acoplado
- **Una sección que reporta entradas pero ninguna gatilla** → no pasa por el gate, solo aparece en el resumen, como cualquier sección `reported`

## Reglas de negocio

- El veredicto es obligatorio — cada Fase debe tener un entry `### Parallelization Verdict`
- El criterio de independencia NO cambia — el que existe hoy se cita sin cambios
- El default NO se invierte — la Fase comienza serial; las marks se aplican solo cuando se establece independencia explícita
- La sección NUNCA gatilla — se renderiza como `reported`, lo que significa que los entries siempre surfacean al usuario, nunca frenan
- El `anchor` es siempre suministrado por la fase de tareas, referenciando la Fase en la descomposición que dio origen a esta entrada
- Un entry cuya `summary` es "serial" pero que nombra un par independiente es una descripción de caso-borde legal, no una contradicción
- La rationale PUEDE ser larga — no está capped como el summary (250 caracteres vs. ningún límite en rationale)
- El `summary` DEBE ser una línea no-vacía; la ausencia o vaciado causa fallo de renderización
- El token `anchor` DEBE estar presente; la ausencia causa fallo de renderización
- **No hay token `contested`** — un veredicto contested sería inerte en `reported`, así que la sección lleva `anchor` solo, igual que `## Decision Gaps`

## Entidades y estados

- **Veredicto** — la declaración de una Fase: parallelized (con id de grupo) o serial (con razonamiento)
- **Entrada** — un item en la tabla `### Parallelization Verdict`, una por Fase; lleva `phase`, `verdict`, `group` (o vacío si serial), `summary`, `anchor`, `rationale`
- **Independencia** — criterio de elegibilidad existente aplicado a la Fase; conservador; por defecto serial sin aplicación de marks
- **Par independiente** → dos tareas dentro de una Fase que podrían ejecutarse concurrentemente pero se mantienen en contexto serial más amplio — debe nombrarse explícitamente en la rationale si es relevante

## Errores de cara al actor

- **Entry faltando la rationale**: renderizador falla nombrando la Fase y el campo; exit 1, stdout vacío
- **Entry faltando anchor**: renderizador falla nombrando la Fase y el campo; exit 1, stdout vacío
- **Summary sobre 250 caracteres**: renderizador falla; exit 1, stdout vacío
- **Número de entries ≠ número de Fases en Breakdown**: validación falla
- **Token `contested` declarado**: contrato-template drift — `--self-check` reporta DRIFT

## Escenarios

### Scenario: Una Fase cuyas tareas fueron marcadas

- **GIVEN** una Fase cuyas tareas todas llevan `· parallel-group: infra`
- **WHEN** se renderiza el retorno
- **THEN** su entry declara la Fase como marcada y la columna `Group` lleva `infra`

### Scenario: Una Fase mantenida serial declara su razonamiento

- **GIVEN** una Fase cuyas tareas fueron evaluadas por independencia y dejadas sin marcar
- **WHEN** se renderiza el retorno
- **THEN** su entry declara el veredicto como serial y su `summary` da el razonamiento por el cual los tasks no son independientes

### Scenario: Exactamente un entry por Fase, siempre

- **GIVEN** una descomposición de cinco Fases
- **WHEN** se renderiza el retorno
- **THEN** `### Parallelization Verdict` lleva exactamente cinco entries, clave-ados por Fase, ninguno omitido

### Scenario: Un par independiente dentro de una Fase serial se nombra

- **GIVEN** una Fase mantenida serial que además contiene un par genuinamente independiente de tasks
- **WHEN** su veredicto se redacta
- **THEN** el razonamiento nombra ese par explícitamente en lugar de generalizar sobre la Fase completa

### Scenario: La sección nunca gatilla

- **GIVEN** un retorno cuyo `### Parallelization Verdict` lleva entries para cada Fase
- **WHEN** el Unresolved Decisions Guard inspecciona el retorno
- **THEN** ningún gate se abre sobre esta sección y ningún turno del usuario se agrega
- **AND** cada entry surfacea en el resumen between-phase

### Scenario: Nada que reportar

- **GIVEN** un retorno `blocked` sin Phases descompuestas
- **WHEN** se renderiza
- **THEN** `### Parallelization Verdict` emite el sentinel `None.` en lugar de ser drapeado

### Scenario: Contrato y template deben coincidir

- **GIVEN** `sdd-tasks.yaml` declarando la nueva sección con `items.rationale`
- **WHEN** `node ~/.claude/scripts/validate-return.js --phase sdd-tasks --self-check` corre contra un `sdd-tasks.md` que NO documenta la sección, o lo hace sin una línea `· rationale:` bajo ese heading
- **THEN** reporta DRIFT y exit 1
- **AND** con ambos archivos alineados exit 0

### Scenario: La fila en la tabla canonical mailbox existe y es legible por el guard

- **GIVEN** la tabla mailbox en la Sección D.3 de `_shared/sdd-phase-common.md` después de este cambio
- **WHEN** el Unresolved Decisions Guard lee las candidatas de sección
- **THEN** encuentra la nueva fila con `reported`/`always`, y trata sus entries como surfaceando, nunca gatillando

### Scenario: La regla de firing-trigger no cambia

- **GIVEN** `.matecito-ai/development-specs/rule/gate-firing-triggers.md`
- **WHEN** este cambio aterriza
- **THEN** está sin modificar — una nueva sección `reported` es una instancia de esa regla, no un cambio a ella

### Scenario: La obligación de veredicto está en el SKILL

- **GIVEN** `sdd-tasks`'s SKILL después de este cambio
- **WHEN** su sección parallel-group se lee
- **THEN** declara que cada Fase adeuda un veredicto, y el criterio de independencia existente se cita sin cambios

### Scenario: El default se mantiene serial

- **GIVEN** una Fase cuyas tareas no fueron evaluadas como independientes
- **WHEN** la descomposición se redacta
- **THEN** los tasks se mantienen sin marcar y la Fase lleva un razonamiento explicado — las marks no se aplican por defecto para evitar escribir un razonamiento

## Referencias

- **Criterio de independencia** → El existente, citado en `sdd-tasks`'s SKILL, sin cambios
- **Formato sección `reported`** → `sdd-verify`'s `## Decision Gaps` — precedent; `anchor` solo, no `contested`
- **Mecanismo split** → `mailbox-item-summary-rationale-split.md` — tres partes, summary capped, anchor libre-form
- **Enumeración canonical** → Sección D.3 de `_shared/sdd-phase-common.md` — `sdd-tasks` | `### Parallelization Verdict` | reported | always
