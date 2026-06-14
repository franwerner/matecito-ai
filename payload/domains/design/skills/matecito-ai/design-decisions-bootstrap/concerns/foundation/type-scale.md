---
name: type-scale
depth: deep
domain: foundation
type: decision
source: Material Design (typography) · escala modular tipográfica clásica · W3C Design Tokens (typography)
---

# Fase: Escala tipográfica

## Qué decide

La escala tipográfica del sistema: las familias de fuente, el ratio modular que genera los tamaños, los pesos disponibles, y los text styles nombrados que materializan cada nivel (`Heading/H1`, `Body/Base`). Define la jerarquía y la legibilidad de todo texto.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. Familias de fuente

> Las familias definen carácter y carga de render. Demasiadas fragmentan la identidad y pesan.

- **Una familia para todo (display + body)** — *default; máxima consistencia y mínimo peso.*
- Dos familias (una display/headings, una body) — cuando la marca quiere contraste tipográfico.
- Tres o más — *raro; solo brand-systems con razón explícita.*
- No sé, recomendame.

### 2. Ratio de la escala modular

> El ratio define el salto entre niveles. Uno consistente hace que la jerarquía se sienta intencional; tamaños sueltos la rompen.

- **1.250 (Major Third)** — *default; jerarquía clara sin saltos excesivos.*
- 1.125 (Major Second) — escala densa, para UI con mucha información.
- 1.333 (Perfect Fourth) — escala expresiva, para `landing`/`marketing-asset`.
- No sé, recomendame.

### 3. Materialización como text styles nombrados

> Un tamaño suelto repetido a mano no es chequeable; un text style nombrado sí. Es lo que hace la escala verificable contra Figma.

- **Cada nivel es un text style nombrado (`Heading/H1`…`Body/Base`, `Caption`)** — *default; condición para chequear drift.*
- Solo headings como styles, body suelto — parcial.
- Sin text styles — *no recomendado; rompe la verificabilidad.*
- No sé, recomendame.

### 4. Pesos disponibles

> Limitar los pesos mantiene la jerarquía legible y el peso de carga bajo.

- **Regular + Semibold/Bold (2 pesos)** — *default.*
- Regular + Medium + Bold (3 pesos).
- Solo Regular — para piezas mínimas.

## Notas de lógica (para el motor)

- **Default según tipo de pieza:** `landing`/`marketing-asset` → proponé ratio expresivo (1.333) en la pregunta 2. `app-ui` → proponé escala densa o Major Third.
- Los tamaños de cuerpo deben respetar el mínimo legible decidido en `contrast-target`/`accessibility` (cuerpo ≥ 16px en web salvo razón).
- Depende de `design-tokens` si la escala se expresa como Figma Variables además de text styles.

## Qué materializar

DDR `type-scale` materializado según el template `~/.claude/references/ddr/templates/ddr.md`. La **Decisión** captura: las familias elegidas, el ratio modular, la lista de niveles con su tamaño/line-height/peso concreto, y los pesos disponibles.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[tool: figma]** cada nivel de la escala existe como un text style nombrado (`Heading/H1`, `Body/Base`, `Caption`); no hay tamaños de texto aplicados como valor suelto.
- **[tool: figma]** el tamaño en px de cada text style coincide con la escala modular declarada (ej: base 16, ratio 1.250 → H3 = 25, H2 = 31, H1 = 39).
- **[tool: figma]** la familia de cada text style es una de las declaradas en la Decisión; no aparecen familias fuera de la lista.
- **[tool: figma]** el peso de cada text style está entre los pesos decididos; no hay pesos fuera de la lista.
- **[manual]** el tamaño de cuerpo cumple el mínimo legible (ej: ≥ 16px en web).

**Alcance:** la lista de text styles nombrados que componen la escala (`Heading/*`, `Body/*`, `Caption`) — el ancla que `mine`/`verify` usan para detectar drift contra Figma.
