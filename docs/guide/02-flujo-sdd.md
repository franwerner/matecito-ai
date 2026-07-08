# 02 — El flujo SDD

[← 01 Overview](01-overview.md) · [Índice](README.md) · Siguiente: [03 — Las fases →](03-fases.md)

El SDD (Spec-Driven Development) es el **flujo del dominio development**: la capa de planificación estructurada para cambios sustanciales sobre un repo de código. Lleva un pedido en lenguaje natural hasta el código, por fases.

> **Qué es del núcleo y qué de development.** El *esqueleto* —base inmutable + add-ons opcionales, las lanes, el INTAKE GATE, el modelo configurable por agente— lo aporta el **núcleo** y es común a todos los dominios. SDD es cómo development **nombra y concreta** ese esqueleto: las fases `sdd-*`, el artefacto de alineación (`spec`), los EDRs y codegraph. Otro dominio (p. ej. design) reusa el mismo esqueleto con sus propios nombres de fase y herramientas.

## El pipeline

```
intake → [explore] → [propose] → spec → [design] → [tasks] → apply → verify → archive
```

- **Base inmutable** (siempre corre): `intake → spec → apply → verify → archive`.
- **Add-ons opcionales** (se activan según el cambio): `explore`, `propose`, `design`, `tasks`. El orquestador inserta cada add-on en su posición canónica.

No todo cambio recorre las nueve fases: un fix trivial va directo; un cambio grande activa todas.

## Las lanes

El orquestador recomienda una lane al estructurar el pedido (en `intake`), y el usuario confirma o ajusta:

| Lane | Qué corre | Cuándo |
|---|---|---|
| **direct** | sin SDD, implementación directa | cambio trivial, sin riesgo |
| **reduced** | base, 0 add-ons | default para cambios sustanciales sin trigger de escalada |
| **custom** | base + solo los add-ons que el cambio necesita | un trigger aislado (ej. una decisión de arquitectura → + `design`) |
| **full** | base + los 4 add-ons | cambio grande / arquitectura cruzando varios dominios |

Sesgo: la **lane mínima viable**. Se escala solo por un motivo concreto.

## El INTAKE GATE (control humano)

Después de `intake`, el orquestador **siempre** muestra el brief y espera **confirmar / ajustar / cancelar** antes de seguir — incluso en modo automático. Es el único punto de control obligatorio: ahí se confirma el alcance, la lane, y si el cambio amerita diagrama o verificación de UI.

## Modo de ejecución

Se elige una vez por sesión:

- **Interactive** (default): después de cada fase muestra el resumen y pregunta "¿Continuamos?".
- **Automatic**: corre las fases de corrido y muestra el resultado final. (El INTAKE GATE igual frena al inicio.)

## Cómo se pasa la información entre fases

Cada fase **lee** los artefactos de las fases previas y **escribe** el suyo. El medio es **Engram** (memoria persistente), no archivos sueltos. El orquestador pasa a cada sub-agente las **topic keys**, no el contenido; el sub-agente lee su artefacto directo de Engram.

Formato de topic key: `sdd/<change-name>/<artefacto>`.

| Fase | Lee (ideal full-lane) | Escribe |
|---|---|---|
| `intake` | el pedido crudo | `intake` |
| `explore` | intake | `explore` |
| `propose` | explore | `proposal` |
| `spec` | proposal (o intake si no hay) | `spec` |
| `design` | proposal + **EDRs** | `design` |
| `tasks` | spec + design | `tasks` |
| `apply` | tasks + spec + design + apply-progress | `apply-progress` |
| `verify` | spec + tasks + apply-progress + **EDRs tocados** | `verify-report` |
| `archive` | todos los artefactos | `archive-report` |

En lanes reducidas, cada fase lee el **upstream más cercano disponible** (ej. `spec` arranca del brief de intake si no hubo proposal).

> **Importante:** los **EDRs no viven en Engram** — viven solo en `.matecito-ai/edr/<dominio>/<slug>.md`. Los artefactos SDD pueden *referenciar* un EDR por `<dominio>/<slug>` (puntero de flujo), nunca su contenido. Ver [04](04-decisiones-edr.md).

## El rol del orquestador

El orquestador es el hilo principal de Claude (no un archivo): mantiene una conversación fina, **delega cada fase a un sub-agente** de contexto fresco, y sintetiza resultados. Los sub-agentes ejecutan; no orquestan ni lanzan otros sub-agentes. Las reglas que el orquestador sigue viven en `CLAUDE.md`.

Pasos condicionales que el orquestador evalúa entre fases (no son fases):
- **Review Workload Guard** (entre tasks y apply) — si el cambio puede pasar el presupuesto de review, propone PRs encadenados.
- **Decision-Gap Capture / mine gate** (entre verify y archive) — si está activo el auto-mine y hay huecos implementados, dispara la minería de EDRs. Ver [05](05-auto-mine.md).
