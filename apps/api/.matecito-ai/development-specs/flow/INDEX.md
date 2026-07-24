# Capability specs — `flow`

Operaciones de cara a un actor, con pasos, ramas y casos borde.

**Cuándo consultar este tipo:** antes de tocar cualquier tool MCP que dispare una operación del orquestador o de un agente de fase (enviar el artefacto de una fase, iniciar o continuar un change) — leé el spec de la operación para conocer sus pasos, sus ramas y sus errores de cara al actor.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `submit-phase-artifact` | Recibe, valida y persiste el artefacto de una fase en el event-log del change | Accepted | [`submit-phase-artifact.md`](submit-phase-artifact.md) |
| `start-change` | Da de alta o continúa un change sobre la rama actual y registra el mapeo rama↔nombre | Accepted | [`start-change.md`](start-change.md) |
| `register-project` | Da de alta o linkea un proyecto por su `project-id` estable (`.id`); idempotente, arranca el watch y el índice | Accepted | [`register-project.md`](register-project.md) |
