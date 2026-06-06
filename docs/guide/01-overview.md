# 01 — Overview

[← Índice](README.md) · Siguiente: [02 — El flujo SDD →](02-flujo-sdd.md)

## El problema

Trabajar con un agente de IA sobre un proyecto tiene tres fugas recurrentes:

| Fuga | Síntoma | Pieza que la ataca |
|---|---|---|
| **Amnesia entre sesiones** | el agente olvida lo descubierto/decidido | **Engram** — memoria persistente |
| **Decisiones implícitas** | las convenciones viven en la cabeza del autor; el agente las reinventa o las viola | **ADRs** — decisiones capturadas y respetadas |
| **Exploración cara** | el agente escanea archivo por archivo para entender el código | **codegraph** — grafo pre-indexado |

## La idea de fondo

Que el agente **no reinvente las convenciones del proyecto en cada sesión**:

1. Las decisiones se capturan **una vez** (como ADRs).
2. El flujo SDD las **respeta** al implementar y **avisa** si algo las contradice.
3. La memoria de trabajo **persiste** entre sesiones (Engram).

El humano decide; la IA ejecuta dentro de esas decisiones.

## Las tres capas

- **Flujo SDD** — el andamio que lleva un pedido de lenguaje natural a código, por fases, con control humano. Ver [02](02-flujo-sdd.md) y [03](03-fases.md).
- **Capa de decisiones** — ADRs + las skills `bootstrap` (capturar), `validate` (chequear) y `mine` (descubrir desde código). Ver [04](04-decisiones-adr.md).
- **Herramientas** — codegraph, context7, drawio, Engram, proofshot, enganchadas donde aportan. Ver [06](06-herramientas.md).

## Qué NO es

No es una herramienta nueva: es la **integración curada** de varias piezas —propias y de terceros— en un flujo coherente. El crédito del trabajo pesado es de los proyectos que integra (ver [README raíz](../../README.md#créditos)).
