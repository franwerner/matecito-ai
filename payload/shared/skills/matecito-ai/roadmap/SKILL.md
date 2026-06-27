---
name: roadmap
description: >
  Capa de planificación y continuidad por encima de los flujos SDD/design. Usá esta skill cuando el
  usuario quiera armar un roadmap, hacer un plan de implementación por fases, dividir un proyecto
  grande en etapas, retomar el roadmap de una sesión anterior, ver qué sigue en su roadmap, ejecutar
  /roadmap new, /roadmap continue, /roadmap next, o cuando describa una iniciativa con múltiples
  fases o pasos a completar en distintas sesiones.
---

# Roadmap — planificación y continuidad por fases

Guía al usuario para definir un roadmap multi-fase de forma conversacional, rastrea el avance
paso a paso, y emite un "next context prompt" al cerrar cada sesión para que la siguiente arranque
sin perder continuidad. Vive por encima del flujo SDD/design: organiza QUÉ hay que hacer; el flujo
ejecuta CÓMO hacerlo.

---

## Invocación

| Comando | Qué hace |
|---|---|
| `/roadmap new <titulo>` | Crea un roadmap nuevo con ese título. |
| `/roadmap continue` | Retoma el roadmap activo (resolución automática; ver regla abajo). |
| `/roadmap next` | Avanza al siguiente paso y, si aplica, propone el comando de flujo. |

Disparadores en lenguaje natural: "armar un roadmap", "plan de implementación por fases",
"dividir en fases", "retomar el roadmap", "qué sigue en mi roadmap", "roadmap", "etapas del
proyecto", "siguiente paso", o cualquier descripción de una iniciativa con múltiples fases.

---

## Layout de archivos (proyecto del usuario)

Todos los artefactos viven en el **proyecto del usuario**, no en el repo de matecito-ai:

```
<raiz-proyecto>/
└── .matecito-ai/
    └── roadmaps/
        └── <titulo>/          # un directorio por roadmap (slug del título)
            ├── INDEX.md       # objetivo, scope, dominio(s), metadata, rollup de progreso
            ├── STEP-0.md      # primer paso
            ├── STEP-1.md
            └── ...
```

`roadmaps/` es hermano de `adr/` y `ddr/`. Cada roadmap tiene su propia carpeta plana.

---

## Reglas de resolución de `/roadmap continue`

1. Si exactamente un roadmap tiene un step con `status: in-progress` → retomarlo directamente.
2. Si ningún roadmap tiene un step `in-progress` → retomar el step `pending` más antiguo del
   roadmap tocado más recientemente (el que tenga la fecha de modificación más reciente en
   su `INDEX.md`).
3. Si múltiples roadmaps tienen steps `in-progress` → listarlos y pedirle al usuario que elija:

   > Hay varios roadmaps en progreso:
   > - `roadmaps/feature-auth/` — STEP-2 en progreso
   > - `roadmaps/infra-migration/` — STEP-1 en progreso
   >
   > ¿Cuál querés retomar?

---

## `/roadmap new <titulo>` — flujo de creación

### Paso 1: Reconciliación inicial

Antes de hacer cualquier pregunta, leer el directorio `.matecito-ai/roadmaps/` del proyecto:

```bash
ls .matecito-ai/roadmaps/ 2>/dev/null || echo "(sin roadmaps)"
```

Si ya existe un roadmap con el mismo título (slug), avisar y preguntar si continuar el existente
o crear uno nuevo (con sufijo numérico para no pisar).

### Paso 2: Definición conversacional del roadmap

Hacer UNA pregunta abierta para arrancar:

> Contame brevemente de qué trata esta iniciativa: cuál es el objetivo principal, y si ya tenés
> idea de cuántas fases o etapas grandes tiene.

No pedir todo de una sola vez. A partir de la respuesta, inferir la estructura inicial y
confirmarla:

> Basándome en lo que contás, propongo estos pasos:
> - STEP-0: [nombre]
> - STEP-1: [nombre]
> - STEP-2: [nombre]
> …
>
> ¿Agregás, sacás o renombrás alguno?

### Paso 3: Definición de cada step (conversacional, profundidad proporcional)

Por cada step del set confirmado:

1. Mostrar el nombre del step.
2. Preguntar qué tareas o entregables concretos tiene ese step. **Profundidad proporcional:**
   - Step simple (configuración, doc, tarea única): una sola pregunta, checklist corto.
   - Step complejo (implementación, refactor grande, migración): profundizar: subtareas, criterios
     de completitud, dependencias.
3. Confirmar el checklist resultante antes de pasar al siguiente.
4. Si el step claramente mapea a un flujo de desarrollo o diseño (ver "Handoff a flujos"), anotarlo
   en ese momento en la metadata del step.

**Regla**: no sobre-especificar steps simples. Una línea de contexto + 3-4 ítems de checklist es
suficiente para un step sencillo.

### Paso 4: Materializar

Crear los archivos usando `INDEX-TEMPLATE.md` y `STEP-TEMPLATE.md` como forma canónica:

1. Crear el directorio: `.matecito-ai/roadmaps/<slug-del-titulo>/`.
2. Escribir `INDEX.md` con el rollup inicial (todos los steps en `pending`).
3. Escribir un `STEP-N.md` por cada step, con `status: pending` y el checklist definido.
4. Reportar los archivos creados.

---

## `/roadmap continue` — flujo de reanudación

### Paso 1: Reconciliación

Leer todos los roadmaps del proyecto y sus steps. Aplicar la regla de resolución de más arriba.

### Paso 2: Reconciliar estado

Al cargar un step `in-progress`, verificar:

- ¿Los ítems del checklist principal reflejan el estado real? (el usuario puede haber marcado
  manualmente)
- ¿El `status:` del header es consistente con el checklist?
  - Si el checklist está todo ticked → proponer cambiar a `done`.
  - Si hay ítems sin ticked → mantener `in-progress`.
- ¿Los `## Pendientes` del step previo (si aplica) están todos resueltos? Si no:
  - Los Pendientes abiertos del step anterior deben aparecer en `## Pendientes` del step actual.
  - Avisar al usuario: "El step anterior tenía N pendientes sin resolver: [lista]. Los moví a
    este step."

### Paso 3: Mostrar estado actual

> Roadmap: **\<titulo\>** — `.matecito-ai/roadmaps/<titulo>/`
> Progreso: X/N steps completados
>
> Step actual: **STEP-K — [nombre]** (in-progress)
> Tareas pendientes en este step:
> - [ ] tarea pendiente 1
> - [ ] tarea pendiente 2
>
> Pendientes del step:
> - [ ] pendiente abierto A
>
> ¿Continuamos con este step o querés hacer algo diferente?

### Paso 4: Trabajo en el step

Asistir al usuario en completar las tareas del checklist. A medida que se completan:

- Marcar `- [x]` en el `STEP-N.md`.
- Si se resuelve un Pendiente, marcarlo `- [x]` también.
- Si aparece un loose end que no es parte del checklist principal y no se decide/resuelve ahora,
  agregarlo como `- [ ]` en `## Pendientes` del step actual.

Actualizar el `STEP-N.md` y el `INDEX.md` (rollup) al final de la sesión.

---

## `/roadmap next` — avanzar al siguiente step

### Paso 1: Verificar completitud del step actual

Un step es `done` SOLO cuando su checklist principal está completamente ticked.

- Si quedan ítems del checklist principales sin ticked → NO avanzar. Mostrar los pendientes y
  preguntar si el usuario quiere marcarlos como completados o dejarlos para después.
- Si el checklist principal está completo pero hay `## Pendientes` abiertos → el step SÍ puede
  marcarse `done`, pero los Pendientes deben carried forward:
  1. Marcar el step actual como `done`.
  2. Copiar los Pendientes abiertos al `## Pendientes` del siguiente step.
  3. Actualizar `INDEX.md` (rollup), incluyendo `pendientes abiertos: N`.
  4. Avisar al usuario: "Step [K] marcado como done. Los N pendientes abiertos se trasladaron
     al Step [K+1]."

### Paso 2: Activar el siguiente step

Marcar el siguiente `STEP-N.md` como `status: in-progress`. Actualizar `INDEX.md`.

### Paso 3: Handoff a flujos (cuando aplica)

Si el step activo mapea a trabajo de desarrollo o diseño, proponer el comando de flujo con el
scope pre-llenado. **El handoff PROPONE; nunca ejecuta autónomamente.** La propuesta pasa por
el INTAKE GATE del flujo destino — el usuario confirma, ajusta o cancela.

**Detección:** el step mapea a un flujo cuando su objetivo es implementar código, un feature, un
refactor, una migración de datos, o producir assets visuales / un sistema de diseño.

**Formato de propuesta:**

Para desarrollo:
> El Step [K] — "[nombre]" es trabajo de desarrollo. Propongo iniciar el flujo SDD con este scope:
>
> `/sdd-new [titulo-roadmap] step [K]: [objetivo del step]`
>
> Scope pre-llenado para intake:
> - Objetivo: [objetivo del step]
> - Tareas: [lista de ítems del checklist principal]
> - Contexto del roadmap: `.matecito-ai/roadmaps/<titulo>/STEP-K.md`
>
> El flujo va a pedir confirmación antes de arrancar (INTAKE GATE). ¿Iniciamos?

Para diseño:
> El Step [K] — "[nombre]" es trabajo de diseño. Propongo iniciar el flujo design con este scope:
>
> `/design-new [titulo-roadmap] step [K]: [objetivo del step]`
>
> Scope pre-llenado para intake:
> - Objetivo: [objetivo del step]
> - Entregables: [lista de ítems del checklist principal]
> - Contexto del roadmap: `.matecito-ai/roadmaps/<titulo>/STEP-K.md`
>
> El flujo va a pedir confirmación antes de arrancar (INTAKE GATE). ¿Iniciamos?

Si el step NO mapea a un flujo (es investigación, definición, doc, decisión, reunion):
> El Step [K] — "[nombre]" no requiere un flujo SDD/design. Trabajamos directamente.
> ¿Arrancamos?

---

## Reconciliación de estado en cada sesión

Al cargar cualquier roadmap (no solo en `continue`), verificar:

1. El `status:` del header de cada step es consistente con su checklist:
   - Checklist 100% ticked y `status: in-progress` → proponer cambiar a `done`.
   - Checklist vacío o sin ticks y `status: done` → inconsistencia; alertar.
2. El rollup en `INDEX.md` refleja el estado real de los steps.
3. Los Pendientes del step anterior marcados como `- [x]` se eliminan del carried-forward (ya
   resueltos).

Si hay inconsistencias, reportarlas y proponer la corrección antes de continuar.

---

## Actualización del INDEX.md (rollup)

Actualizar el rollup de `INDEX.md` cada vez que cambie el `status:` de un step o la cuenta de
pendientes abiertos. El rollup tiene formato machine-readable para que la skill lo lea sin
ambigüedad:

```
<!-- rollup
total: N
done: X
in-progress: Y
pending: Z
pendientes-abiertos: P
-->
```

Los `pendientes-abiertos` son la suma de todos los `- [ ]` en secciones `## Pendientes` de
steps que no están en `done`.

---

## Reglas de Pendientes

- Son loose ends que aparecieron MIENTRAS se trabajaba un step y NO están en el checklist
  principal: preguntas sin responder, decisiones aplazadas, tareas que surgieron fuera de scope,
  descubrimientos que necesitan seguimiento.
- Se escriben como `- [ ]` en `## Pendientes`. Sin timestamps. Tickables cuando se resuelven.
- Cuando un step se marca `done`, sus Pendientes abiertos se TRASLADAN al step siguiente (no se
  pierden, no se eliminan en silencio).
- Se cuentan en el rollup de `INDEX.md` como `pendientes-abiertos: N`.
- Un step puede marcarse `done` con Pendientes abiertos, siempre que el checklist principal esté
  completo. Los Pendientes abiertos son deuda visible, no bloqueantes de `done`.

---

## Handoff a flujos — invariantes

- `/roadmap next` NUNCA ejecuta `/sdd-new` o `/design-new` autónomamente.
- El scope se pre-llena como PROPUESTA de input para el intake del flujo destino.
- El usuario confirma, ajusta o cancela en el INTAKE GATE del flujo.
- El roadmap no tiene autoridad sobre el flujo — solo provee contexto de partida.

---

## Anti-patterns

- No sobre-especificar steps simples. Si el step es "actualizar dependencias", un checklist de
  3 ítems es suficiente.
- No crear Pendientes con timestamps — son solo `- [ ]` con texto.
- No eliminar Pendientes abiertos al cerrar un step — trasladarlos al siguiente.
- No ejecutar el flujo SDD/design sin pasar por el INTAKE GATE.
- No avanzar al siguiente step si el checklist principal tiene ítems sin ticked (a menos que el
  usuario lo indique explícitamente y se documente la razón).
- No sobrescribir un `INDEX.md` existente sin antes leerlo y preservar el contenido anterior.
- No asumir que el estado en memoria coincide con el estado en disco — siempre leer los archivos
  al iniciar.
