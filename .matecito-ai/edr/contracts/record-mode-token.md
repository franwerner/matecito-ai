# EDR — Una proposal declara `record-mode: create | modify`; sdd-apply la lee como token de ruteo, no como veredicto

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
El mecanismo de captura in-flow sólo sabía crear un EDR nuevo. No existía forma de que un cambio in-flow editara un EDR Accepted ya existente, y esto dejó sin camino formal ediciones como las de Slice F sobre structure/change-isolation-activation-flag.md y structure/change-level-worktree-isolation.md.

**Nota (port-turn-mechanism):** ambos registros nombrados arriba fueron retirados en ese cambio — el
aislamiento por espacio de trabajo que cada uno gobernaba dejó de existir. La mención sigue siendo válida
como ejemplo histórico de por qué `modify` hacía falta; no es una cita resoluble y no se repara.

## Decisión
Una proposal declara `· record-mode: create | modify` junto a `· record:` en el ítem de `### New Decisions` de sdd-design. Closed value set pero sin `passing:` — cualquier valor declarado es legal; sólo la ausencia falla `TOKEN-MISSING`. `sdd-apply` Step 4b ramifica sobre el token: `create` son los cuatro pasos de siempre (render, escribir, aplicar filas de INDEX); `modify` abre el archivo nombrado, edita sólo las cláusulas que la proposal ratificada nombra, corre `validate-artifact.js`, y no llama a `render-artifact.js` ni agrega fila de INDEX. Una declaración que no coincide con la realidad del disco (`modify` sobre un archivo inexistente, o `create` sobre uno que ya existe) es una falla — nunca un cambio silencioso de rama, nunca un overwrite. El token es de ruteo, leído verbatim por `sdd-apply` — no un veredicto que el orquestador clasifica: sigue el precedente de `record:`, no el de `structure/verdict-classified-by-the-orchestrator.md`.

## Reglas verificables
- **[auto]** `sdd-design.yaml` declara `record-mode` como token con `values: [create, modify]` y sin `passing:` en los `items.tokens` de `### New Decisions`; `validate-return-tokens.test.js` cubre la presencia y el conjunto de valores.
- **[auto]** `validate-return.js` rechaza un ítem de `### New Decisions` sin `· record-mode:` con `TOKEN-MISSING`, y un valor fuera de `[create, modify]` con `TOKEN-ILLEGAL`.
- **[manual]** `sdd-apply` Step 4b ramifica la Materialización sobre `record-mode`: `create` renderiza + escribe + aplica las filas de INDEX; `modify` abre el archivo, edita sólo las cláusulas nombradas, y corre `validate-artifact.js` sin renderizar y sin agregar fila de INDEX.
- **[manual]** `modify` sobre un archivo inexistente, o `create` sobre uno que ya existe, es una falla explícita — nunca un cambio silencioso de modo ni un overwrite.

## Alternativas consideradas
B — supersede-and-replace: contradice la regla ya ratificada de editar en el lugar sin borrar, y necesita un status que la leyenda de EDR no tiene. C — modo update de development-decisions-bootstrap fuera del flujo: cuesta la corrida desatendida que este cambio busca, y parte un mismo cambio entre dos mecanismos sin nada que los una; rechazada explícitamente por el costo a la corrida desatendida. D — re-render completo desde `--data`: cualquier contenido no retipeado desaparece en silencio. Las tres se presentaron al usuario con su costo antes de elegir A.

## Consecuencias
`modify` no puede bootstrapear un store ausente — un `modify` contra un `domain`/`slug` inexistente es la falla de declaración-vs-realidad, nunca una primera materialización. El límite "sólo esa cláusula cambió" es legible pero no chequeable por máquina: `validate-artifact.js` prueba que el registro sigue siendo estructuralmente válido, nunca que el cambio tocó únicamente lo que la proposal decía. Costo aceptado, no mitigado.
