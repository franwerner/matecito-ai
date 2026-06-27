---
step: N
title: "Título del step"
status: pending
# status: pending | in-progress | done
---

# STEP-N — Título del step

> **Objetivo:** Una oración que describe qué se logra al completar este step.

## Checklist

- [ ] Tarea o entregable 1
- [ ] Tarea o entregable 2
- [ ] Tarea o entregable 3

<!-- Agregar o quitar ítems según la complejidad real del step.
     Un step simple puede tener 2-4 ítems; uno complejo puede tener más.
     El step es "done" SOLO cuando todos los ítems están ticked. -->

## Pendientes

<!-- Loose ends que aparecieron MIENTRAS se trabajó este step y no están en el checklist
     principal. Sin timestamps. Se marcan con - [x] cuando se resuelven.
     Los que quedan abiertos al marcar el step "done" se trasladan al step siguiente. -->

- [ ] (ejemplo) Decidir entre opción A y opción B antes de avanzar
- [ ] (ejemplo) Revisar X que quedó sin definir

<!-- Si no hay pendientes, dejar la sección vacía (no eliminar el encabezado). -->

## Next context prompt

<!-- Esta sección es para el usuario: copiar y pegar en una nueva sesión para retomar
     sin perder contexto. Se actualiza al cerrar la sesión en este step.
     Solo es "viva" para el step `in-progress` (el que vas a retomar); en steps `done` o
     `pending` es opcional y puede omitirse o quedar desactualizada. -->

**Roadmap:** [título del roadmap]
**Carpeta:** `.matecito-ai/roadmaps/<titulo>/`
**Step completado:** STEP-[N-1] — [nombre del step anterior y resumen de una línea de lo logrado]
**Step actual:** STEP-N — [título de este step]
**Objetivo del step:** [una oración]
**Tareas pendientes en este step:**
- [ ] [tarea sin completar 1]
- [ ] [tarea sin completar 2]
**Pendientes abiertos:** [N] (ver `## Pendientes` arriba)
**Archivo a leer primero:** `.matecito-ai/roadmaps/<titulo>/STEP-N.md`
