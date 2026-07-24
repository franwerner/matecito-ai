# EDR — Modelado de datos

- **Status:** Pending
- **Type:** decision
- **Date:** 2026-07-23

## Razón de omisión / aplazamiento

**Status:** Pending

El schema concreto de tablas es un contrato y se define cuando se construya el store, no antes. Trigger esperado: cuando aterrice la construcción del store.

Cuando se defina, incluirá una tabla de proyectos (identidad, ruta, nombre, alta) y el scope por `project_id` en todas las tablas (eventos, changes, agentes, artefactos, referencias de código, índice de decisiones), en línea con el modelo de instancia única global. Acá también se fijan la clave concreta de idempotencia de las escrituras y el shape del envelope `{type, payload}` que consume la UI.
