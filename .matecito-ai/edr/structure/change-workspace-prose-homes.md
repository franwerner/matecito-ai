# EDR — La prosa del espacio de trabajo del cambio se reparte por especificidad, no se centraliza

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto
El mecanismo de aislamiento por cambio toca tres lugares distintos del payload: el kernel domain-agnóstico, la prosa específica del dominio de desarrollo, y la referencia del mecanismo de fan-out del batch paralelo. Meter todo en un único lugar obliga a elegir uno de los tres como hogar, y los otros dos terminan citándolo por un concepto que en realidad les pertenece en parte.

## Decisión
La política general —opt-in, momento de apertura, quién integra— vive en el kernel, en vocabulario neutral de dominio. Los detalles concretos —identidad del espacio de trabajo, los comandos del mecanismo elegido— viven en la guía del dominio de desarrollo, que ata el término neutral del kernel a ese mecanismo concreto. El anidamiento del batch de implementación paralelo sobre el espacio de trabajo del cambio vive en la referencia de fan-out del batch paralelo, donde ya vivía la nesting del batch mismo.

## Reglas verificables
- **[manual]** La política general del espacio de trabajo del cambio (opt-in, momento de apertura, quién integra) vive en `payload/core/CLAUDE.md`, en vocabulario neutral de dominio, sin nombrar git ni worktree.
- **[manual]** Los detalles de git (identidad, comandos) viven en `payload/domains/development/CLAUDE.md`, que ata el término neutral 'workspace' del kernel a un git worktree concreto.
- **[manual]** El anidamiento del batch de implementación paralelo sobre el espacio de trabajo del cambio vive en `parallel-batch.md`, junto al resto de la nesting del batch.

## Alternativas consideradas
Poner los comandos de git en el kernel. Descartado: el kernel declara mecanismos, no herramientas — cada dominio elige las suyas. Un archivo de referencia nuevo, dedicado sólo al espacio de trabajo del cambio. Descartado: fuera del alcance confirmado del cambio, y duplica contenido que ya tiene un hogar natural en los tres lugares existentes.

## Consecuencias
Tres archivos se tocan en vez de uno, pero ninguno de los tres termina citando un concepto que en realidad le pertenece a otro. Quien lee sólo el kernel entiende la política sin saber que el mecanismo concreto es git; quien lee el fragmento de desarrollo encuentra los comandos exactos sin tener que releer la política general.

## Relacionados
- `relacionado-con` → [change-level-worktree-isolation.md](change-level-worktree-isolation.md) — el mecanismo cuya prosa esta decisión reparte.
- `relacionado-con` → [consolidation-run-is-the-integrator.md](consolidation-run-is-the-integrator.md) — otra decisión de este mismo cambio que también fija quién ejecuta qué, en el mismo dominio.
