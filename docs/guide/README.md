# Guía del dominio development (flujo SDD)

Cómo funciona el dominio **development** de matecito-ai de punta a punta: las fases del SDD, cómo se conectan con las herramientas y skills, cómo se pasa la información entre fases, y la capa de decisiones (EDRs) con `bootstrap`, `validate` y `mine`.

> **Dónde encaja.** matecito-ai es un **microkernel**: un núcleo agnóstico (orquestación, flujo estructurado, gate humano, memoria) más **un plugin por área** (`development`, `design`, …). Esta guía profundiza el dominio **development** — desarrollo de software con SDD sobre un repo de código. El núcleo y el resto de los dominios están en el [README raíz](../../README.md) y en el [contrato de área](../../payload/domains/README.md). Lo que sigue es el *binding concreto* de development; el esqueleto (fases base + add-ons, INTAKE GATE, modelo por agente) lo aporta el núcleo, los nombres y herramientas los aporta este dominio.

## Mapa del dominio

```
        decisiones (EDRs)            memoria (Engram)         exploración
        bootstrap · validate · mine  artifacts entre fases    codegraph · context7
         (development)                (núcleo, compartido)     (development · núcleo)
                    │                        │                        │
                    └──────────────┬─────────┴────────────┬───────────┘
                                   │   FLUJO SDD           │
   intake → [explore] → [propose] → spec → [design] → [tasks] → apply → verify → archive
```

- **Flujo SDD**: el pipeline de development — lleva un pedido en lenguaje natural hasta el código, por fases, con un punto de control humano al inicio. El esqueleto (base inmutable + add-ons, gate) es del núcleo; los nombres de fase (`sdd-*`) son de development.
- **Capa de decisiones**: las decisiones de arquitectura se capturan como **EDRs** (el decision record de development; en otros dominios cambia el tipo — DDR en design) y el flujo las respeta.
- **Memoria (núcleo)**: cada fase deja su artefacto en Engram; la siguiente lo lee. Persiste entre sesiones. Compartida por todos los dominios.
- **Herramientas**: las declara cada dominio en su manifest (nada global). development usa Engram (memoria del núcleo), context7, codegraph, drawio y proofshot — se enganchan en las fases donde aportan.

## Orden de lectura

1. [01 — Overview](01-overview.md) — qué problema resuelve y las tres piezas de base.
2. [02 — El flujo SDD](02-flujo-sdd.md) — pipeline, lanes, gate, y cómo pasa la información entre fases.
3. [03 — Las fases](03-fases.md) — qué hace cada fase, qué lee/escribe, qué herramienta usa.
4. [04 — Decisiones y EDRs](04-decisiones-edr.md) — qué es un EDR, concerns vs EDR, y la tríada bootstrap/validate/mine.
5. [05 — Auto-mine de EDRs](05-auto-mine.md) — detección de decisiones in-flow (`flagDecisionGaps`).
6. [06 — Herramientas](06-herramientas.md) — codegraph, context7, drawio, Engram, proofshot.
7. [07 — Configuración](07-configuracion.md) — `config.json`, scope, TUI.

## Principio rector

**El humano decide; la IA ejecuta dentro de esas decisiones.** El agente no reinventa las convenciones en cada sesión: se capturan una vez (EDRs en development), se respetan al implementar, y la memoria persiste (Engram). Es el principio del núcleo; este dominio lo concreta para desarrollo de software.
