# Cómo crear un design concern (guía de autoría)

Guía canónica para agregar un concern al catálogo de diseño (el "ratchet"). Sirve a quien mantiene la skill —vos, el equipo, o Claude cuando el bootstrap detecta un tema fuera del catálogo—. Seguila para que todo concern nuevo nazca con el mismo formato y las mismas reglas de higiene que los existentes.

El valor de largo plazo de la skill es que **nunca se vuelva a olvidar un tema de diseño**. Un concern agregado una vez queda cubierto para todas las piezas futuras y se ofrece a las viejas vía modo update.

---

## Antes de escribir: 3 decisiones

1. **Surface canónica.** ¿A qué surface pertenece? Mirá el "criterio de pertenencia" en cada `concerns/<surface>/INDEX.md`. La taxonomía es **cerrada** (`foundation`, `components`, `layout`, `brand`, `accessibility`); no inventes una surface nueva. Si genuinamente no encaja en ninguna, es señal de que falta una surface en la taxonomía — eso es una decisión de catálogo aparte (tocar el motor + `concerns/INDEX.md`), no algo que se resuelve metiéndolo a la fuerza.
2. **Profundidad.** `deep` (cuestionario propio de 3-5 preguntas, para decisiones grandes con condicionales) o `light` (1-2 preguntas, para temas acotados). Referencia: `foundation/color-palette.md` es `deep`; `foundation/spacing-grid.md` es `light`.
3. **Type.** `decision` (alternativas y trade-offs reales) · `convention` (acuerdo de estilo, sin gran dilema) · `policy` (regla verificable, a menudo de accesibilidad).

---

## Template

Copiá esta estructura. Las secciones marcadas (opcional) se incluyen solo si aplican.

```markdown
---
name: <slug-en-kebab-case>
depth: <deep | light>
domain: <surface canónica>
type: <decision | convention | policy>
source: <taxonomía/estándar de origen: W3C Design Tokens, Atomic Design, WCAG 2.x, Material Design, etc.>
---

# Fase: <título legible>

## Qué decide

<1-3 líneas: qué decisión de diseño captura esta fase y por qué importa. Esta línea se usa como el "por qué importa" en la entrevista. Sin sermones.>

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default marcado, y siempre "no sé, recomendame".

### 1. <título de la pregunta>

> <una línea de por qué importa esta pregunta puntual>

- ***<opción A> — default para <condición>.*** <aclaración breve>
- **<opción B>** — <cuándo conviene>
- <opción C>
- No sé, recomendame.

### 2. <…> (las que hagan falta; light = 1-2, deep = 3-5)

## Notas de lógica (para el motor)   <!-- opcional: incluir si hay condicionales -->

- <defaults según tipo de pieza detectado en context>
- <preguntas condicionales: "si en la 1 eligió X, saltear la 2">
- <dependencias de otras fases: "requiere que <otra-fase> esté Accepted">

## Qué materializar

DDR `<name>` materializado según el template canónico [`templates/ddr.md`](~/.claude/references/ddr/templates/ddr.md). Especificá qué campos concretos y **verificables contra Figma** debe contener:
- **Reglas verificables:** nombrá las reglas como **valores concretos chequeables contra el archivo Figma** (hex exacto, ratio de contraste, escala modular, px, nombres de tokens) — NUNCA adjetivos vagos. Cada una con su mecanismo de verificación al inicio: `[tool: figma]` (chequeable leyendo styles/components del archivo), `[tool: contrast]` (cálculo de ratio WCAG) o `[manual]` (revisión humana).
- **Alcance** (solo concerns con locator espacial): la lista de tokens / set de componentes / frames —a nivel convención, patrones estables, no nodos concretos efímeros— que la decisión gobierna y que el validador usa como ancla de drift contra Figma.
- **Relacionados** (si aplica): vínculos tipados esperados con otros DDRs.
```

---

## Reglas de higiene (qué NO hacer)

Estas reglas protegen el patrón de la skill. Romperlas degrada el catálogo con el tiempo.

- ❌ **No metas instrucciones de implementación visual paso a paso.** El concern decide *qué* y *por qué*; el *cómo* (dibujar en Figma) lo hace el diseñador/agente después, leyendo el DDR. Sin recetas de "hacé un autolayout con gap 16 acá".
- ❌ **No pinees una herramienta como obligatoria.** Mal: "usá Tokens Studio". La herramienta la decide el usuario. Las herramientas se nombran solo **como ejemplos ilustrativos** de una opción (ej: "plugin de tokens (Tokens Studio, Figma Variables)").
- ❌ **No prescribas una paleta/escala concreta como dogma.** El concern ofrece opciones abstractas con ejemplos; el usuario elige y eso se registra en el DDR. El concern no dice "usá #2563EB" — dice "un primario de marca con sus tints/shades" como forma de la decisión.
- ❌ **No escribas opciones sin default ni sin "no sé, recomendame".** Toda pregunta ofrece un default razonado y la salida "no sé, recomendame" para quien no tiene opinión.
- ❌ **No uses lenguaje vago en "Qué materializar".** Mal: "buenos colores con buen contraste". Bien: reglas verificables con valores concretos chequeables contra Figma (ej: "texto normal ≥ 4.5:1 contra su fondo", "todo color de UI referencia un color style nombrado"), cada una con su mecanismo `[tool: figma]` / `[tool: contrast]` / `[manual]`. Lo que va al DDR tiene que ser chequeable contra el archivo.
- ❌ **No dupliques contenido de otra fase.** Si una decisión ya la captura otro concern, referencialo, no lo repitas.

---

## Checklist de autovalidación

Antes de guardar el concern, verificá:

- [ ] Frontmatter completo: `name`, `depth`, `domain`, `type`, `source`.
- [ ] `domain` es una de las surfaces canónicas y coincide con la carpeta donde va el archivo.
- [ ] Tiene las 3 secciones núcleo: `## Qué decide`, `## Preguntas`, `## Qué materializar`.
- [ ] Cada pregunta tiene línea de "por qué importa", default marcado, y "no sé, recomendame".
- [ ] `deep` tiene 3-5 preguntas; `light` tiene 1-2.
- [ ] **Cero recetas de implementación visual.**
- [ ] Las herramientas aparecen solo como ejemplos entre paréntesis, no como imposición.
- [ ] "Qué materializar" describe reglas verificables con valores concretos chequeables contra Figma (no adjetivos vagos), cada una con su mecanismo `[tool: figma]`/`[tool: contrast]`/`[manual]`; si el concern tiene locator, pide la sección `Alcance` con tokens/components/frames.
- [ ] Si hay condicionales o dependencias → están en "Notas de lógica".

---

## Pasos de integración

Una vez que el concern pasa el checklist:

1. **Colocá el archivo** en `concerns/<surface>/<slug>.md`.
2. **Actualizá el índice de la surface** (`concerns/<surface>/INDEX.md`): sumá una fila a la tabla "Concerns en esta surface" con `[<slug>](<slug>.md)`, profundidad, tipo y el "Qué decide" (resumido a ≤160 caracteres).
3. **Actualizá la matriz raíz** (`concerns/INDEX.md`): sumá la fila en la sección de la surface que corresponda, con la aplicabilidad por tipo de pieza (Requerido/Recomendado).
4. **Verificá los links:** que el link del nuevo concern en ambos índices resuelva a un archivo real.

Desde ese momento, todo bootstrap futuro considera el concern, y el modo update lo ofrece a piezas viejas.
