# Rúbrica de coherencia y completitud (diseño)

Lista central de chequeos que aplica `design-decisions-validate`. Es **ratchet-able**: cuando aparece una contradicción nueva, se agrega acá y queda cubierta para siempre.

## Cómo la lee el validador

Cada chequeo tiene: **severidad** (CRITICAL / WARNING / SUGGESTION), una **condición** evaluada sobre los DDRs, la/las **surface(s)** donde viven los DDRs involucrados (para localizar los archivos), y un **mensaje** (qué/por qué/sugerencia). El validador evalúa las condiciones contra `.matecito-ai/ddr/<surface>/` y reporta las que se cumplen.

## Mapa slug → surface

Para localizar el archivo de cada DDR nombrado abajo. Un DDR vive en `.matecito-ai/ddr/<surface>/<slug>.md`. (La ruta del archivo es la fuente de verdad si hay duda.)

| Slug | Surface |
|---|---|
| color-palette, type-scale, spacing-grid, design-tokens | foundation |
| component-library | components |
| grid-breakpoints | layout |
| tone-of-voice, logo-usage | brand |
| contrast-target | accessibility |

Los DDRs custom project-local viven igual bajo una surface canónica (`.matecito-ai/ddr/<surface>/<slug>.md`); su slug no está en esta tabla pero la surface se lee de la ruta.

---

## Chequeos genéricos

### Completitud

- **[WARNING]** Una fase relevante para el tipo de pieza no tiene DDR (ni siquiera `Not Applicable` como fila en el INDEX de su surface). Hueco silencioso. *(Requiere la lista de fases relevantes o el catálogo `concerns/INDEX.md`; si no están disponibles, marcar como "no verificable".)*
- **[NOTA — Inferred]** DDRs con `Status: Inferred` NO cierran la preocupación: no los contés como decisión tomada en el conteo de completitud; no reportes como defecto las secciones Contexto/Decisión/Consecuencias/Alternativas/Reglas verificables vacías ni el porqué vacío (esperados en Inferred); sí verificá que `## Alcance` siga matcheando el archivo Figma (ancla de drift, solo si tenés acceso al archivo). Un Inferred convive con un WARNING de completitud si la preocupación sigue sin una decisión `Accepted`.

### Higiene de status

- **[WARNING]** DDR `Accepted` sin sección "Decisión" con contenido concreto (valores chequeables contra Figma).
- **[WARNING]** DDR `Pending` o `Deferred` sin razón ni trigger/condición de revisión.
- **[CRITICAL]** DDR `Superseded` sin link "Reemplazado por", o el DDR linkeado no existe. *(El link puede ser intra-surface `<slug>.md` o cross-surface `../<surface>/<slug>.md`; verificá ambos.)*
- **[WARNING]** Una fila `Not Applicable` en el INDEX de una surface sin razón.

### Verificabilidad

> Esta es la columna vertebral del sistema de diseño: una `## Reglas verificables` solo vale si se puede chequear contra Figma. Cada regla debe ser una aserción con **valores concretos** (hex, ratio de contraste, escala/pasos, px, nombres de tokens/styles/componentes) y su **mecanismo** al inicio (`[tool: figma]`, `[tool: contrast]`, `[manual]`).

- **[WARNING]** DDR `Accepted` con sección `## Reglas verificables` vacía o ausente → una decisión de diseño sin reglas chequeables no se puede verificar contra Figma.
- **[WARNING]** Una regla bajo `## Reglas verificables` formulada como **adjetivo vago** en vez de un valor concreto chequeable (ej: "colores armoniosos", "buen contraste", "tipografía legible", "espaciado consistente") → no es chequeable contra Figma. Debe ser un valor: hex (`#2563EB`), ratio (`≥ 4.5:1`), escala (`50–900, 10 pasos`), px (`8px base`), o nombre de token/style/componente (`Primary/500`).
- **[SUGGESTION]** Una regla bajo `## Reglas verificables` sin marca de mecanismo al inicio (`[tool: figma]` / `[tool: contrast]` / `[manual]`) → no queda claro cómo se chequea; agregar el mecanismo.
- **[SUGGESTION]** Lenguaje vago de obligatoriedad ("tratá de", "en lo posible", "idealmente", "evitar cuando se pueda", "preferiblemente") en las reglas de un DDR `Accepted` → ablanda la regla hasta volverla no-chequeable.

### Integridad de la taxonomía

- **[CRITICAL]** Existe una carpeta bajo `.matecito-ai/ddr/` que no es una surface canónica (`foundation` · `components` · `layout` · `brand` · `accessibility`). La taxonomía es cerrada; una surface nueva es decisión de catálogo, no de proyecto.
- **[WARNING]** Un DDR está listado en el índice raíz (`.matecito-ai/ddr/INDEX.md`) pero su surface no tiene `INDEX.md`, o viceversa (índice de surface con un DDR que no figura en el raíz). Índices desincronizados.
- **[SUGGESTION]** Una surface tiene `INDEX.md` pero ningún DDR-archivo (carpeta de surface vacía en la salida). Limpiar la carpeta o el índice, o listar la surface como "sin uso" en el raíz.

### Coherencia del campo `Type`

- **[SUGGESTION]** Un DDR marcado `Type: convention` o `Type: policy` tiene una sección "Alternativas consideradas" sustanciosa → quizá es en realidad una `decision`; revisar el type.
- **[SUGGESTION]** Un DDR marcado `Type: decision` y `Accepted` sin "Alternativas consideradas" ni "Consecuencias" → una decisión sin trade-offs documentados es sospechosa; o falta contenido o es en realidad una convention.
- **[WARNING]** Un DDR `Type: policy` `Accepted` sin "Reglas verificables" accionables → una política sin reglas chequeables contra Figma no se puede cumplir ni verificar.

### Trazabilidad a Figma (sección `Alcance`)

- **[WARNING]** Un DDR con sección `## Alcance` cuyo locator (styles / components / frames nombrados) no existe en el archivo Figma conectado → drift: el sistema cambió o la decisión quedó obsoleta. *(Requiere acceso al archivo Figma vía el MCP figma; si no está disponible, marcar como "no verificable" y delegar la detección de drift a `design-decisions-mine`.)*
- **[SUGGESTION]** Un DDR de foundation/components/layout `Accepted` que gobierna styles/componentes/frames concretos pero sin sección `## Alcance` → una decisión que ancla en Figma sin locator no es verificable contra drift; considerar agregarlo.

---

## Contradicciones conocidas (combinaciones entre DDRs)

La columna **Surface(s)** indica dónde viven los DDRs de la condición. "cross" = la contradicción cruza surfaces.

| # | Severidad | Surface(s) | Condición | Mensaje |
|---|---|---|---|---|
| 1 | CRITICAL | foundation | Dos DDRs `Accepted` de foundation definen **paletas de color incompatibles** (ej: `color-palette` enumera `Primary/500` = `#2563EB` y otro DDR/regla usa `Primary/500` = `#1D4ED8`) | El mismo token de rol no puede tener dos hex. Definir una sola fuente de verdad para cada color style. |
| 2 | CRITICAL | foundation | Dos DDRs `Accepted` definen **escalas tipográficas o de espaciado distintas** para el mismo eje (ej: `type-scale` con 8 pasos vs otra regla con 6; `spacing-grid` base 8px vs 4px) | Dos escalas para el mismo eje rompen la consistencia y el tokenizado. Unificar la escala (pasos y base). |
| 3 | CRITICAL | accessibility + foundation (cross) | `contrast-target` `Accepted` exige un ratio mínimo (ej: ≥ 4.5:1) **y** un par texto/fondo definido en `color-palette` no lo cumple con los hex enumerados | La paleta viola el objetivo de contraste declarado. Ajustar los hex del par o el rol semántico para cumplir el ratio. |
| 4 | WARNING | foundation | `color-palette` `Accepted` materializa colores como **styles nombrados** (`Primary/500`) **y** otra regla aplica colores de rol como **hex inline suelto** fuera de la lista de styles | Color de rol aplicado a mano rompe la verificabilidad contra Figma. Todo color de rol debe referenciar un color style. |
| 5 | WARNING | foundation + components (cross) | `design-tokens` / `color-palette` / `type-scale` `Accepted` (hay tokens) **y** `component-library` define componentes con **valores crudos** en vez de referenciar esos tokens | Componentes con valores hardcodeados se desincronizan de la foundation. Los componentes deben consumir los tokens/styles. |
| 6 | WARNING | components | `component-library` `Accepted` define un set de variantes/estados **y** falta el estado de error/disabled/focus que las reglas o el tipo de pieza requieren | Set de componentes incompleto: los estados faltantes se improvisan por pieza. Completar las variantes/estados. |
| 7 | WARNING | layout + foundation (cross) | `grid-breakpoints` `Accepted` define una grilla/espaciado **y** su base no es múltiplo de la base de `spacing-grid` (ej: grid de 12 col con gutter 10px vs base 8px) | La grilla y el espaciado base no encajan; los componentes no alinean a la grilla. Alinear el gutter/margen a la base de espaciado. |
| 8 | WARNING | accessibility + components (cross) | `contrast-target` / objetivos de accesibilidad `Accepted` **y** `component-library` define estados (hover/disabled) cuyos colores no se chequearon contra el ratio | Estados de componente con contraste no verificado: disabled/hover suelen romper accesibilidad. Chequear el ratio de cada estado. |
| 9 | WARNING | foundation | `color-palette` define **dark mode** vía tokens semánticos **y** no todos los tokens semánticos resuelven a un valor válido en ambos modos | Token semántico sin valor en un modo deja un agujero en dark mode. Completar el valor del token para cada modo. |
| 10 | SUGGESTION | brand + foundation (cross) | `tone-of-voice` / `logo-usage` `Accepted` define una expresión de marca (ej: "minimalista, sobria") **y** `color-palette` propone una paleta saturada/recargada que choca con esa expresión | La paleta no acompaña la expresión de marca declarada. Revisar si la foundation refleja la identidad. |
| 11 | SUGGESTION | brand | `logo-usage` `Accepted` define un área de protección / tamaño mínimo en px **y** no hay regla verificable con el valor concreto (solo prosa) | Una regla de uso de logo sin px concretos no es chequeable contra Figma. Enumerar el área de protección y el tamaño mínimo en px. |

---

## Ratchet

Cuando encontrás una contradicción o un chequeo útil que no está acá:

- Genérico → agregalo como bullet en la sección "Chequeos genéricos" (en la subsección que corresponda).
- Combinación entre DDRs → agregá una fila en "Contradicciones conocidas".

Siempre con severidad, surface(s) y mensaje (qué/por qué/sugerencia). Si la combinación cruza surfaces, marcala "cross" en la columna Surface(s). Así el validador la atrapa de ahí en más.
