---
name: tone-of-voice
depth: light
domain: brand
type: convention
source: práctica clásica de brand voice / content design · guías de marca
---

# Fase: Tono de voz

## Qué decide

El tono y los atributos de voz del copy (cómo "habla" la marca) y los do/don't que lo mantienen consistente en todas las piezas. Sin esto, cada texto suena distinto y la identidad se diluye.

## Preguntas

Una o dos, según haga falta.

### 1. Atributos de voz

> Definir 3-4 atributos concretos hace el tono accionable y revisable; "amigable pero profesional" sin más es vago.

- **3-4 atributos con su opuesto (`cercano no informal`, `claro no simplista`)** — *default; accionable y verificable.*
- Una frase de posicionamiento de voz — más liviano, menos chequeable.
- Sin definición formal — *solo si el copy lo provee otro equipo.*
- No sé, recomendame.

### 2. Reglas de do/don't

> Ejemplos concretos de qué decir y qué no es lo que hace el tono replicable. **Solo si se definieron atributos.**

- **Tabla do/don't con ejemplos reales del producto** — *default.*
- Solo lineamientos generales (sin ejemplos) — menos útil.

## Notas de lógica (para el motor)

- **Default según tipo de pieza:** `brand-system`/`marketing-asset` → este concern es Requerido. `app-ui` → suele delegarse a microcopy y puede ser `Recomendado`.
- Si el copy lo provee otro equipo/externo, materializá con `Status: Pending` y razón.

## Qué materializar

DDR `tone-of-voice` materializado según el template `~/.claude/references/ddr/templates/ddr.md`. La **Decisión** captura: los atributos de voz con su opuesto, y la tabla de do/don't con ejemplos concretos del producto.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** el copy de cada pieza es coherente con los atributos de voz declarados; un atributo violado (ej: texto informal cuando se decidió "no informal") es un hallazgo.
- **[manual]** los textos no caen en los "don't" enumerados (ej: jerga prohibida, mayúsculas de grito, exclamaciones múltiples).
- **[manual]** los CTAs y mensajes de error/vacío siguen el tono decidido, no solo los titulares.

Si el copy lo provee otro equipo, el DDR va con `Status: Pending` y la razón concreta; en ese caso las reglas quedan como referencia para ese equipo.
