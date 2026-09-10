# Capability specs — `flow`

Operaciones de cara a un actor, con pasos, ramas y casos borde.

**Cuándo consultar este tipo:** antes de tocar cualquier tool MCP que dispare una operación del orquestador o de un agente de fase (enviar el artefacto de una fase, iniciar o continuar un change), la superficie de lectura que consume la UI (snapshot + push en vivo), o cualquier comando y pantalla del CLI que la persona dispara a mano (instalar, actualizar, configurar modelos, sincronizar) — leé el spec de la operación para conocer sus pasos, sus ramas y sus errores de cara al actor.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `ratify-gate-items` | Un gate abre con un índice único de items ratificables, luego presenta cada item uno a uno, y cada item nombra su fuente anclada | Accepted | [`ratify-gate-items.md`](ratify-gate-items.md) |
| `install-ecosystem` | Instala/actualiza en un comando lo que falte del ecosistema: plan combinado, dry-run, confirmación y ejecución continue-on-error | Accepted | [`install-ecosystem.md`](install-ecosystem.md) |
| `update-ecosystem` | Reconcilia binarios, payload y configuración del host; termina en error si algún componente falló | Accepted | [`update-ecosystem.md`](update-ecosystem.md) |
| `resume-self-replace-run` | La corrida relanzada tras el auto-reemplazo del ejecutable: excluye la acción propia, no pregunta, no imprime plan y sigue mostrando progreso | Inferred | [`resume-self-replace-run.md`](resume-self-replace-run.md) |
| `configure-agent-model` | La persona fija el modelo por agente de un dominio, por scope, con herencia explícita y sin persistir overrides espurios | Inferred | [`configure-agent-model.md`](configure-agent-model.md) |
| `discovery-runs-in-explore` | Discovery es formulada y resuelta dentro de la exploración: lee el código primero, luego pregunta, llevando las respuestas en el artefacto de exploración | Accepted | [`discovery-runs-in-explore.md`](discovery-runs-in-explore.md) |
| `sync-via-tui` | Sincronización desde la interfaz interactiva: plan, confirmación y progreso en vivo, sin re-ejecutarse nunca a sí misma; diferimiento del payload al arranque siguiente cuando hay auto-reemplazo | Accepted | [`sync-via-tui.md`](sync-via-tui.md) |
| `two-fixed-lanes` | Dos lanes fijos: `full` siempre por defecto, `direct` solo si el usuario lo pide explícitamente; sin fork, sin recomendación, sin confirmación | Accepted | [`two-fixed-lanes.md`](two-fixed-lanes.md) |
| `confirm-brief-before-dispatch` | Después que intake retorna el brief, un gate obligatorio lo ofrece para aceptar o corregir; nada se despacha hasta que el usuario responde | Accepted | [`confirm-brief-before-dispatch.md`](confirm-brief-before-dispatch.md) |
| `find-durable-records` | Localiza los registros duraderos —decisiones y comportamiento— que gobiernan un trabajo antes de escribir código, proponer un diseño o verificar contra ellos, a través de tres caminos acumulativos (índice, búsqueda literal, búsqueda semántica) con degradación por presencia | Accepted | [`find-durable-records.md`](find-durable-records.md) |
| `materialize-records-straight-through` | `sdd-apply` lee `## New Decisions` del artefacto de diseño y escribe el record `Accepted` directo, sin canal de forwarding, sin ledger por-cambio y sin gate de confirmación intermedio | Accepted | [`materialize-records-straight-through.md`](materialize-records-straight-through.md) |
| `claim-shared-branch-turn` | Una sesión reclama por llamada explícita el turno de la rama compartida, hace su unidad de trabajo y lo libera con el recibo que el reclamo le dio; asesor, sin interceptación, con reintento único y sin liberación automática | Accepted | [`claim-shared-branch-turn.md`](claim-shared-branch-turn.md) |
