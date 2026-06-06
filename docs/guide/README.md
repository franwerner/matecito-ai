# Guía del flujo SDD de matecito-ai

Cómo funciona el ecosistema de punta a punta: las fases del SDD, cómo se conectan con las herramientas y skills, cómo se pasa la información entre fases, y la capa de decisiones (ADRs) con `bootstrap`, `validate` y `mine`.

## Mapa del ecosistema

```
        decisiones (ADRs)            memoria (Engram)         exploración (MCP)
        bootstrap · validate · mine  artifacts entre fases    codegraph · context7
                    │                        │                        │
                    └──────────────┬─────────┴────────────┬───────────┘
                                   │   FLUJO SDD           │
   intake → [explore] → [propose] → spec → [design] → [tasks] → apply → verify → archive
```

- **Flujo SDD**: lleva un pedido en lenguaje natural hasta el código, por fases, con un punto de control humano al inicio.
- **Capa de decisiones**: las decisiones de arquitectura se capturan como ADRs y el flujo las respeta.
- **Memoria**: cada fase deja su artefacto en Engram; la siguiente lo lee. Persiste entre sesiones.
- **Herramientas (MCP/CLI)**: codegraph, context7, drawio, proofshot — se enganchan en las fases donde aportan.

## Orden de lectura

1. [01 — Overview](01-overview.md) — qué problema resuelve y las tres piezas de base.
2. [02 — El flujo SDD](02-flujo-sdd.md) — pipeline, lanes, gate, y cómo pasa la información entre fases.
3. [03 — Las fases](03-fases.md) — qué hace cada fase, qué lee/escribe, qué herramienta usa.
4. [04 — Decisiones y ADRs](04-decisiones-adr.md) — qué es un ADR, concerns vs ADR, y la tríada bootstrap/validate/mine.
5. [05 — Auto-mine de ADRs](05-auto-mine.md) — detección de decisiones in-flow (`flagDecisionGaps`).
6. [06 — Herramientas](06-herramientas.md) — codegraph, context7, drawio, Engram, proofshot.
7. [07 — Configuración](07-configuracion.md) — `config.json`, scope, TUI.

## Principio rector

**El humano decide; la IA ejecuta dentro de esas decisiones.** El agente no reinventa las convenciones en cada sesión: se capturan una vez (ADRs), se respetan al implementar, y la memoria persiste (Engram).
