---
name: code-conventions
depth: deep
domain: structure
type: convention
source: prácticas idiomáticas por lenguaje · Clean Code · guías de estilo oficiales (PEP 8, Effective Go, Airbnb, Rustfmt/Clippy)
---

# Fase: Convenciones de código

## Qué decide

Las convenciones **estilísticas** de cómo se escribe el código a nivel micro: cómo se representan conjuntos de valores, la ausencia, la inmutabilidad, los tipos, los literales, el casing de nombres, la forma de las funciones, la comparación y la iteración. Son idiomáticas del lenguaje: la *pregunta* es agnóstica, la *respuesta* (y su default) depende del stack.

**Frontera con otras fases (no re-preguntar acá):**
- El **sufijo de rol** de un archivo (`.controller`, `.use-case`) y su ubicación → `folder-structure`. Acá va solo el **casing** de nombres, no el sufijo.
- Las **reglas de qué capa importa a qué** → `layers-and-dependencies`. Acá va solo el **estilo/orden** de imports, no la semántica de capas.
- **throw vs Result/Either** → `runtime/error-handling`. No se decide acá.

## Preguntas

Una por turno, pero como casi todas son lintables con un default claro del lenguaje, el motor puede **proponer el default del lenguaje para todas juntas y confirmar en bloque**, abriendo pregunta puntual solo donde el usuario quiera desviarse. Para cada una: incluí "no sé, recomendame".

### 1. Conjuntos cerrados de valores

> Un set fijo de valores comparado con strings sueltos es frágil: un typo no lo atrapa el compilador.

- ***Tipo dedicado (enum nativo / union de literales / typed constants) — default.*** Nunca comparar contra magic strings.
- Strings/números sueltos con constantes nombradas — solo si el lenguaje no tiene mejor herramienta.

### 2. Representación de la ausencia

> Cómo se modela "no hay valor" define cuántos null-checks y cuántos NPE vas a tener.

- Según el lenguaje: `Option`/`Maybe` (Rust, FP), `Optional` (Java), nullable explícito (`T | null` en TS con strict), `None` (Python). Evitar sentinels (`-1`, `""`) para representar ausencia.

### 3. Inmutabilidad por defecto

> Estado mutable compartido es la fuente #1 de bugs difíciles.

- ***Inmutable por defecto (`const`/`readonly`/`final`/`val`), mutabilidad explícita y justificada — default.***
- Mutable por defecto — solo si el lenguaje/perf lo exige.

### 4. Estrictez de tipos

> Los escape hatches vacían de valor al type system.

- ***Prohibir `any`/`interface{}`/`dynamic` salvo excepción documentada; tipos explícitos en bordes públicos, inferencia adentro — default.***
- Tipado laxo — solo prototipos/scripts.

### 5. Literales mágicos

> Un `86400` o un `"active"` desperdigado es ilegible y multi-fuente.

- ***Números/strings mágicos → constante o config nombrada — default.*** (Se solapa con #1 para sets de valores.)

### 6. Casing y naming de símbolos

> Sin regla, la búsqueda por nombre y la lectura se degradan. (Solo el **casing**; el **sufijo de rol** es de `folder-structure`.)

- Casing por clase de símbolo según el lenguaje (archivos, tipos/clases, funciones, variables, constantes).
- Booleanos con prefijo `is`/`has`/`can`; evitar abreviaturas no estándar.
- Confirmá el set del lenguaje detectado o ajustá.

### 7. Forma de función

> Anidamiento profundo y firmas anchas disparan la complejidad.

- ***Guard clauses / early return sobre anidamiento; límite de parámetros (≈3-4) → options object/struct por encima; evitar boolean params (→ enum/options) — default.***

### 8. Igualdad y comparación

> La coerción implícita mete bugs sutiles.

- ***Comparación estricta sin coerción (`===` en JS/TS, evitar `==`); explicitar valor vs identidad — default.***

### 9. Iteración y colecciones

> Loops imperativos con mutación son más propensos a error que transformaciones declarativas.

- ***`map`/`filter`/`reduce` u equivalentes sobre loops manuales cuando aplica; no mutar la colección mientras se itera — default.*** No dogmático: un loop es válido cuando es más claro.

### 10. Estilo de imports

> Orden y forma de imports afectan legibilidad y refuerzan el encapsulamiento. (Solo estilo; las **reglas de capa** son de `layers-and-dependencies`.)

- Orden/agrupado de imports; absolute vs relative; barrel files sí/no; prohibir deep imports a internos de otro módulo (importar por su índice público).

## Notas de lógica (para el motor)

- **Default por lenguaje (Fase 0):** proponé el idiom del stack como default de cada pregunta y pedí confirmación. Ej — TS: enum/union, `T | null` con strict, `readonly`/`const`, no `any`, `===`, kebab archivos + PascalCase tipos + camelCase funciones. Go: typed constants + `iota`, valores cero/`ok`-idiom, no `interface{}` gratuito, comparación por valor, MixedCaps. Python: `Enum`, `None`/`Optional`, dataclasses frozen, type hints + mypy strict, snake_case + PascalCase clases. Rust: `enum`+match, `Option`, inmutable por defecto (ya lo es), sin `unwrap` gratuito, snake_case + CamelCase tipos.
- **Marcá el enforcement de cada regla:** la mayoría son lintables → en el EDR van como `[tool: eslint / golangci-lint / ruff / clippy]` con la regla concreta (ej. `no-magic-numbers`, `@typescript-eslint/no-explicit-any`, `eqeqeq`). Las no chequeables por linter (ej. "options object sobre boolean param") van `[manual]`.
- **Batch:** no hagas 10 turnos si el usuario acepta los defaults del lenguaje — proponé el set completo, confirmá en bloque, y profundizá solo donde quiera desviarse.
- Si el proyecto ya eligió linter/formatter en `tech`, alineá estas reglas con esa herramienta (no propongas una convención que el linter elegido no pueda enforzar sin config extra; avisá si hace falta).

## Qué materializar

EDR `code-conventions` materializado según `~/.claude/references/edr/templates/edr.md`. Debe contener:

- **Contexto** y **Decisión**: las convenciones elegidas por cada punto tratado, en términos de concepto (no nombres de clases del proyecto). Enunciá cada una como una regla ("los conjuntos cerrados se modelan con enum/union; la ausencia con `T | null`; …").
- **Reglas verificables**: cada convención como aserción chequeable con su mecanismo al inicio — `[tool: <linter> <regla>]` para las lintables, `[manual]` para el resto. Es la parte más importante: el registro real de estas convenciones es la regla de linter; el EDR las agrupa y justifica.
- **Relacionados**: `relacionado-con` → `folder-structure` (naming: sufijos allá, casing acá), `layers-and-dependencies` (imports: reglas de capa allá, estilo acá), `runtime/error-handling` (representación de errores allá).
