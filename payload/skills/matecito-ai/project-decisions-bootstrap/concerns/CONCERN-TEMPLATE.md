# Cómo crear un concern (guía de autoría)

Guía canónica para agregar un concern al catálogo (el "ratchet"). Sirve a quien mantiene la skill —vos, el equipo, o Claude cuando el bootstrap detecta un tema fuera del catálogo—. Seguila para que todo concern nuevo nazca con el mismo formato y las mismas reglas de higiene que los 39 existentes.

El valor de largo plazo de la skill es que **nunca se vuelva a olvidar un tema**. Un concern agregado una vez queda cubierto para todos los proyectos futuros y se ofrece a los viejos vía modo update.

---

## Antes de escribir: 3 decisiones

1. **Dominio canónico.** ¿A qué dominio pertenece? Mirá el "criterio de pertenencia" en cada `concerns/<dominio>/INDEX.md`. La taxonomía es **cerrada** (10 activos + 7 reservados); no inventes un dominio nuevo. Si genuinamente no encaja en ninguno, es señal de que falta un dominio en la taxonomía — eso es una decisión de catálogo aparte (tocar el motor + `concerns/INDEX.md`), no algo que se resuelve metiéndolo a la fuerza.
2. **Profundidad.** `deep` (cuestionario propio de 3-5 preguntas, para decisiones grandes con condicionales) o `light` (1-2 preguntas, para temas acotados). Referencia: `runtime/error-handling.md` es `deep`; `runtime/caching.md` es `light`.
3. **Type.** `decision` (alternativas y trade-offs reales) · `convention` (acuerdo de estilo, sin gran dilema) · `policy` (regla verificable, a menudo de seguridad/operación).

---

## Template

Copiá esta estructura. Las secciones marcadas (opcional) se incluyen solo si aplican.

```markdown
---
name: <slug-en-kebab-case>
depth: <deep | light>
domain: <dominio canónico>
type: <decision | convention | policy>
source: <taxonomía/estándar de origen: ISO/IEC 25010, 12-factor, arc42 §N, OWASP ASVS §N, SRE, etc.>
---

# Fase: <título legible>

## Qué decide

<1-3 líneas: qué decisión captura esta fase y por qué importa. Esta línea se usa como el "por qué importa" en la entrevista. Sin sermones.>

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

- <defaults según stack/tipo de proyecto detectado en context>
- <preguntas condicionales: "si en la 1 eligió X, saltear la 2">
- <dependencias de otras fases: "requiere que <otra-fase> esté Accepted">

## Tech a registrar   <!-- opcional: incluir solo si la fase puede involucrar una tecnología concreta -->

<qué mini-ADR de tech crear si se elige una herramienta específica. Ej: "la librería de Result si se usa (returns, neverthrow)">

## Qué materializar

ADR `<name>` materializado según el template canónico [`templates/adr.md`](~/.claude/references/adr/templates/adr.md). Especificá qué campos concretos y verificables debe contener:
- **Reglas verificables:** nombrá las reglas como valores concretos (no adjetivos vagos), cada una con su mecanismo de verificación al inicio: `[tool: <herramienta>]` o `[manual]`.
- **Alcance** (solo concerns espaciales/estructurales): los globs a nivel convención —patrones estables, no archivos concretos— que la decisión gobierna.
- **Relacionados** (si aplica): vínculos tipados esperados con otros ADRs.
```

---

## Reglas de higiene (qué NO hacer)

Estas reglas protegen el patrón de la skill. Romperlas degrada el catálogo con el tiempo.

- ❌ **No metas código de implementación.** El concern decide *qué* y *por qué*; el *cómo* (el código) lo escribe el agente después, en el repo, leyendo el ADR. Sin snippets de Python/JS/etc. mostrando cómo implementar la decisión.
  - **Única excepción:** una fase cuyo *output sea un artefacto de configuración* (como `arch-enforcement`, que produce config de linter). Aun así, el concern describe *qué* artefacto generar y con qué reglas — **no trae la config hardcodeada**. El agente la escribe traduciendo las reglas reales del proyecto.
- ❌ **No pinees versiones de tecnologías.** Mal: "usá FastAPI 0.115". La versión la decide el usuario y vive en el catálogo `tech/`, no hardcodeada en el concern. Las tecnologías se nombran solo **como ejemplos ilustrativos** de una opción (ej: "proveedor externo (Auth0, Keycloak, Cognito)").
- ❌ **No prescribas una tecnología.** El concern ofrece opciones abstractas con ejemplos; el usuario elige y eso se registra en `tech/`. El concern no dice "usá Redis" — dice "cache distribuido (Redis/Memcached)" como una opción entre otras.
- ❌ **No escribas opciones sin default ni sin "no sé, recomendame".** Toda pregunta ofrece un default razonado y la salida "no sé, recomendame" para quien no tiene opinión.
- ❌ **No uses lenguaje vago en "Qué materializar".** Mal: "manejá bien los errores". Bien: reglas verificables con valores concretos (ej: "access token 15min, refresh 7d con rotación"), cada una con su mecanismo `[tool: <herramienta>]` o `[manual]`. Lo que va al ADR tiene que ser chequeable.
- ❌ **No dupliques contenido de otra fase.** Si una decisión ya la captura otro concern, referencialo, no lo repitas.

---

## Checklist de autovalidación

Antes de guardar el concern, verificá:

- [ ] Frontmatter completo: `name`, `depth`, `domain`, `type`, `source`.
- [ ] `domain` es uno de los canónicos y coincide con la carpeta donde va el archivo.
- [ ] Tiene las 3 secciones núcleo: `## Qué decide`, `## Preguntas`, `## Qué materializar`.
- [ ] Cada pregunta tiene línea de "por qué importa", default marcado, y "no sé, recomendame".
- [ ] `deep` tiene 3-5 preguntas; `light` tiene 1-2.
- [ ] **Cero bloques de código** (salvo que la fase sea de artefacto-output, y aun así sin config hardcodeada).
- [ ] **Cero versiones pineadas** de tecnologías.
- [ ] Las tecnologías aparecen solo como ejemplos entre paréntesis, no como imposición.
- [ ] "Qué materializar" describe reglas verificables con valores concretos (no adjetivos vagos), cada una con su mecanismo `[tool]`/`[manual]`; si el concern es estructural, pide la sección `Alcance` con globs.
- [ ] Si hay condicionales o dependencias → están en "Notas de lógica".
- [ ] Si puede involucrar una tecnología → está la sección "Tech a registrar".

---

## Pasos de integración

Una vez que el concern pasa el checklist:

1. **Colocá el archivo** en `concerns/<dominio>/<slug>.md`.
2. **Actualizá el índice del dominio** (`concerns/<dominio>/INDEX.md`): sumá una fila a la tabla "Concerns en este dominio" con `[<slug>](<slug>.md)`, profundidad, tipo y el "Qué decide" (resumido a ≤160 caracteres).
3. **Actualizá la matriz raíz** (`concerns/INDEX.md`): sumá la fila en la sección del dominio que corresponda, con la aplicabilidad por tipo de proyecto (Requerido/Recomendado).
4. **Si el dominio era reservado** (estaba sin concerns): pasalo de la tabla "Reservados" a "Activos" en `concerns/INDEX.md`, y quitá el `.gitkeep` de la carpeta. El `INDEX.md` del dominio cambia su nota de "reservado" por la tabla de concerns.
5. **Verificá los links:** que el link del nuevo concern en ambos índices resuelva a un archivo real.

Desde ese momento, todo bootstrap futuro considera el concern, y el modo update lo ofrece a proyectos viejos.

---

## Si la contradicción también es nueva

Si al agregar el concern detectás una combinación que puede contradecir a otro (ej: el concern nuevo choca con `auth` en ciertos casos), agregá también la regla a la rúbrica del validador (`project-decisions-validate/coherence-rules.md`), con severidad, dominio(s) y mensaje qué/por qué/sugerencia. Así el validador la atrapa de ahí en más.
