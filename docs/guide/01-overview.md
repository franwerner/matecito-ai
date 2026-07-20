# 01 — Overview

[← Índice](README.md) · Siguiente: [02 — El flujo SDD →](02-flujo-sdd.md)

> **Esta guía cubre el dominio `development`.** matecito-ai es un **microkernel**: un núcleo agnóstico más un plugin por área (`development`, `design`, …). El núcleo aporta el mecanismo (flujo estructurado, gate humano, memoria persistente); cada dominio lo concreta con su vocabulario, sus fases y sus herramientas. Lo que sigue describe cómo development concreta ese mecanismo para desarrollo de software con SDD. Visión general del ecosistema: [README raíz](../../README.md).

## El problema

Trabajar con un agente de IA sobre un proyecto tiene tres fugas recurrentes. El **núcleo** da el mecanismo para atacarlas; el dominio **development** las concreta con su pieza:

| Fuga | Síntoma | Pieza que la ataca (en development) |
|---|---|---|
| **Amnesia entre sesiones** | el agente olvida lo descubierto/decidido | **Engram** — memoria persistente _(núcleo, compartida por todos los dominios)_ |
| **Decisiones implícitas** | las convenciones viven en la cabeza del autor; el agente las reinventa o las viola | **EDRs** — decisiones capturadas y respetadas _(el decision record de development; en design es DDR)_ |
| **Exploración cara** | el agente escanea archivo por archivo para entender el código | **codegraph** — grafo pre-indexado _(índice de exploración de development; en design es Figma)_ |

## La idea de fondo

Que el agente **no reinvente las convenciones del proyecto en cada sesión**:

1. Las decisiones se capturan **una vez** (como EDRs).
2. El flujo SDD las **respeta** al implementar y **avisa** si algo las contradice.
3. La memoria de trabajo **persiste** entre sesiones (Engram).

El humano decide; la IA ejecuta dentro de esas decisiones.

## Las tres capas (en development)

- **Flujo SDD** — el andamio que lleva un pedido de lenguaje natural a código, por fases, con control humano. El esqueleto (base + add-ons, INTAKE GATE) es del núcleo; las fases `sdd-*` son de development. Ver [02](02-flujo-sdd.md) y [03](03-fases.md).
- **Capa de decisiones** — EDRs + las skills `bootstrap` (capturar), `validate` (chequear) y `mine` (descubrir desde código). El concepto de decision record es del núcleo; el tipo concreto (EDR) es de development. Ver [04](04-decisiones-edr.md).
- **Herramientas** — las declara cada dominio en su manifest (nada global): development usa Engram (memoria del núcleo), context7, codegraph, drawio y proofshot, enganchadas donde aportan. Ver [07](07-herramientas.md).

## Dónde encaja development

development es **uno** de los dominios del ecosistema; convive con otros (`design`, …) bajo el mismo núcleo. Agregar un dominio nuevo es **implementar el contrato de área** —`manifest.json` + un fragmento `CLAUDE.md`— sin tocar el núcleo (ver [cómo se arma un dominio](../../payload/domains/README.md)). Por eso esta guía habla de SDD, EDR y codegraph: son el *binding* de development, no el ecosistema entero.

## Qué NO es

No es una herramienta nueva: es la **integración curada** de varias piezas —propias y de terceros— en un flujo coherente. El crédito del trabajo pesado es de los proyectos que integra (ver [README raíz](../../README.md#créditos)).
