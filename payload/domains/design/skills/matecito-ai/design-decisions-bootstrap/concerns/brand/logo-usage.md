---
name: logo-usage
depth: light
domain: brand
type: policy
source: práctica clásica de brand guidelines (logo usage / clear space) · guías de identidad
---

# Fase: Uso del logo

## Qué decide

Las reglas de uso del logo: qué variantes son válidas, el área de protección (clear space), el tamaño mínimo y los usos prohibidos. Sin reglas explícitas, el logo termina deformado, sobre fondos ilegibles o en tamaños donde se rompe.

## Preguntas

Una o dos, según haga falta.

### 1. Variantes y reglas de uso

> Definir las variantes válidas y cuándo usar cada una evita que se elija a ojo y se rompa la identidad.

- **Set de variantes (full color / monocromo / negativo) + regla de cuándo cada una** — *default para `brand-system`.*
- Una sola variante con sus reglas de fondo — para marcas simples.
- Logo externo provisto sin reglas propias — *Pending; documentar la fuente.*
- No sé, recomendame.

### 2. Protección y mínimos

> Clear space y tamaño mínimo son las reglas que mantienen el logo legible y respirado. **Solo si la marca tiene logo propio.**

- **Clear space en múltiplos de una medida del logo (ej: ½ x la altura) + tamaño mínimo en px** — *default; verificable.*
- Solo tamaño mínimo — parcial.

## Notas de lógica (para el motor)

- **Default según tipo de pieza:** Requerido para `brand-system`. Para `landing`/`marketing-asset` que solo *usan* un logo existente, suele bastar referenciar el DDR del brand-system; si no existe, capturarlo acá.
- Si el logo es externo sin guía propia, materializá con `Status: Pending` y la fuente.

## Qué materializar

DDR `logo-usage` materializado según el template `~/.claude/references/ddr/templates/ddr.md`. La **Decisión** captura: las variantes válidas del logo, la regla de cuándo usar cada una, el clear space (en múltiplos concretos), el tamaño mínimo en px, y los usos prohibidos enumerados.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[tool: figma]** cada variante del logo existe como un componente nombrado (`Logo/FullColor`, `Logo/Mono`, `Logo/Negative`); no hay logos pegados como imagen suelta sin variante.
- **[tool: figma]** ninguna instancia del logo está por debajo del tamaño mínimo en px decidido.
- **[manual]** el clear space alrededor del logo respeta el múltiplo declarado; ningún elemento invade esa zona.
- **[manual]** no aparece ninguno de los usos prohibidos enumerados (logo deformado, recoloreado fuera de la paleta, sobre fondo de bajo contraste).

**Alcance:** los componentes de logo nombrados (`Logo/*`) — el ancla de drift contra Figma.

Si el logo es externo sin guía propia, el DDR va con `Status: Pending` y la fuente; en ese caso no lleva Reglas verificables propias.
