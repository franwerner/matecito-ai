# EDR — NFR de performance

- **Status:** Pending
- **Date:** 2026-07-23

## Razón de omisión / aplazamiento

**Status:** Pending

No hay objetivos numéricos para un cliente local (no aplica un P99). Postura cualitativa registrada: el riesgo es el canvas y el feed con muchos nodos, y las vistas de diff sobre archivos grandes; la mitigación ya elegida es virtualización de listas/feed (TanStack Virtual), culling de viewport de ReactFlow, y virtualización del viewer de diffs. Trigger: fijar números formales cuando el sistema escale a muchos nodos/eventos.
