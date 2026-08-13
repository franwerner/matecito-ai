# Capability specs — `flow`

Operaciones de cara a un actor, con pasos, ramas y casos borde.

**Cuándo consultar este tipo:** antes de tocar cualquier tool MCP que dispare una operación del orquestador o de un agente de fase (enviar el artefacto de una fase, iniciar o continuar un change), la superficie de lectura que consume la UI (snapshot + push en vivo), o cualquier comando y pantalla del CLI que la persona dispara a mano (instalar, actualizar, configurar modelos, sincronizar) — leé el spec de la operación para conocer sus pasos, sus ramas y sus errores de cara al actor.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `submit-phase-artifact` | Recibe, valida y persiste el artefacto de una fase en el event-log del change | Accepted | [`submit-phase-artifact.md`](submit-phase-artifact.md) |
| `start-change` | Da de alta o continúa un change sobre la rama actual y registra el mapeo rama↔nombre | Accepted | [`start-change.md`](start-change.md) |
| `register-project` | Da de alta o linkea un proyecto por su `project-id` estable (`.id`); idempotente, arranca el watch y el índice | Accepted | [`register-project.md`](register-project.md) |
| `serve-change-state` | Sirve a la UI el estado de un change: snapshot del event-log + push en vivo por WS con resume por posición | Accepted | [`serve-change-state.md`](serve-change-state.md) |
| `serve-code-diff` | Sirve a la UI los diffs de código de un change, derivados on demand de las fotos: por batch, colapsado, contra el archivo real y cambios fuera del flujo | Accepted | [`serve-code-diff.md`](serve-code-diff.md) |
| `iteration-notes` | Notas de iteración creadas en la UI (ancladas a evento/archivo/líneas), señaladas al chat en el próximo mensaje y leídas/resueltas por el AI con confirmación del usuario | Accepted | [`iteration-notes.md`](iteration-notes.md) |
| `navigate-cockpit` | El usuario recorre las secciones del cockpit (rail de 5 destinos) y ve la identidad y el estado de la corrida en el encabezado | Accepted | [`navigate-cockpit.md`](navigate-cockpit.md) |
| `ratify-gate-items` | Un gate abre con un índice único de items ratificables, luego presenta cada item uno a uno, y cada item nombra su fuente anclada | Accepted | [`ratify-gate-items.md`](ratify-gate-items.md) |
| `materialize-mined-artifacts` | Post-gate de confirmación de minería, materializa candidatos confirmados a archivos; obtiene cuerpo e INDEX entries del renderer, nunca a mano | Accepted | [`materialize-mined-artifacts.md`](materialize-mined-artifacts.md) |
| `install-ecosystem` | Instala/actualiza en un comando lo que falte del ecosistema: plan combinado, dry-run, confirmación y ejecución continue-on-error | Accepted | [`install-ecosystem.md`](install-ecosystem.md) |
| `update-ecosystem` | Reconcilia binarios, payload y configuración del host; termina en error si algún componente falló | Accepted | [`update-ecosystem.md`](update-ecosystem.md) |
| `resume-self-replace-run` | La corrida relanzada tras el auto-reemplazo del ejecutable: excluye la acción propia, no pregunta, no imprime plan y sigue mostrando progreso | Inferred | [`resume-self-replace-run.md`](resume-self-replace-run.md) |
| `configure-agent-model` | La persona fija el modelo por agente de un dominio, por scope, con herencia explícita y sin persistir overrides espurios | Inferred | [`configure-agent-model.md`](configure-agent-model.md) |
| `sync-via-tui` | Sincronización desde la interfaz interactiva: plan, confirmación y progreso en vivo, sin re-ejecutarse nunca a sí misma; diferimiento del payload al arranque siguiente cuando hay auto-reemplazo | Accepted | [`sync-via-tui.md`](sync-via-tui.md) |
