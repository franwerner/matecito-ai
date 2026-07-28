# Capability specs — `flow`

Operaciones de cara a un actor, con pasos, ramas y casos borde.

**Cuándo consultar este tipo:** antes de tocar cualquier tool MCP que dispare una operación del orquestador o de un agente de fase (enviar el artefacto de una fase, iniciar o continuar un change), o la superficie de lectura que consume la UI (snapshot + push en vivo) — leé el spec de la operación para conocer sus pasos, sus ramas y sus errores de cara al actor.

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
