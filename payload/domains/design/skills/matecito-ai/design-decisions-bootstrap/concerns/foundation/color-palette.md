---
name: color-palette
depth: deep
domain: foundation
type: decision
source: Material Design (color system) · W3C Design Tokens (color)
---

# Fase: Paleta de color

## Qué decide

La paleta de color del sistema: qué roles existen (primario, secundario, semánticos, neutros), sus valores concretos, y cómo se materializan como color styles nombrados en Figma. Es la decisión de foundation más visible: define identidad y, junto con `contrast-target`, accesibilidad.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. Estructura de roles

> Los roles definen para qué sirve cada color, no solo qué color es. Sin roles, cada pieza inventa sus colores.

- **Primario + neutros + semánticos (success/warning/error/info)** — *default para `app-ui` y `brand-system`.*
- Primario + secundario + neutros + semánticos — cuando la marca necesita un acento secundario fuerte.
- Solo marca + neutros — *default para `marketing-asset` y `landing` simples.*
- No sé, recomendame.

### 2. Generación de tints y shades

> Cada color de rol suele necesitar variaciones (hover, disabled, fondos). Definir cómo se generan evita un arcoíris improvisado.

- **Escala fija por color (ej: 50–900, 10 pasos)** — *default; predecible y tokenizable, estilo Material.*
- Solo base + hover/active (3 pasos) — suficiente para piezas chicas.
- Ad-hoc por necesidad — *no recomendado para un sistema.*
- No sé, recomendame.

### 3. Materialización como styles nombrados

> Un color suelto repetido a mano no es una decisión chequeable; un color style nombrado sí. Esto es lo que hace la paleta verificable contra Figma.

- **Todo color de la paleta es un color style nombrado (`Primary/500`, `Neutral/100`)** — *default; condición para que `mine` y `verify` chequeen drift.*
- Solo los roles principales como styles, neutros sueltos — parcial.
- Sin styles, colores aplicados a mano — *no recomendado; rompe la verificabilidad.*
- No sé, recomendame.

### 4. Modo oscuro

> Si la pieza necesita dark mode, la paleta debe preverlo desde el inicio (tokens semánticos vs. valores crudos).

- No, solo light — *default si no hay requerimiento.*
- Sí, light + dark vía tokens semánticos (un token apunta a distinto valor por modo).
- Sí, pero más adelante (DDR `Pending`).

## Notas de lógica (para el motor)

- **Default según tipo de pieza:** `marketing-asset` y `landing` simples → proponé "solo marca + neutros" en la pregunta 1. `app-ui`/`brand-system` → proponé roles completos con semánticos.
- **Pregunta 4 condicional:** si el tipo de pieza no contempla dark mode y el usuario no lo pide, materializá "solo light" sin abrir la pregunta.
- Depende de `contrast-target`: los pares de color texto/fondo deben cumplir el ratio decidido ahí.

## Qué materializar

DDR `color-palette` materializado según el template `~/.claude/references/ddr/templates/ddr.md`. La **Decisión** captura: los roles elegidos, los valores hex concretos de cada color de rol, la escala de tints/shades, y si hay dark mode.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[tool: figma]** cada color de la paleta existe como un color style nombrado en el archivo (`Primary/500`, `Neutral/100`, `Error/500`); no hay colores de rol aplicados como valor suelto.
- **[tool: figma]** el hex de cada color style coincide exactamente con el valor enumerado en la Decisión (ej: `Primary/500` = `#2563EB`).
- **[tool: figma]** todo fill/stroke de un elemento de UI referencia un color style de la paleta, no un hex inline fuera de la lista.
- **[tool: contrast]** los pares texto/fondo definidos cumplen el ratio mínimo decidido en `contrast-target` (delegado a ese DDR).
- **[manual]** si hay dark mode: cada token semántico resuelve a un valor válido en ambos modos.

**Alcance:** la lista de color styles nombrados que componen la paleta (`Primary/*`, `Neutral/*`, `Error/*`, …) — el ancla que `mine`/`verify` usan para detectar drift contra Figma.
