---
name: development-spec-validate
description: Validador de coherencia, completitud y verificabilidad de los capability-specs (comportamiento) de un proyecto, organizados por tipo en .matecito-ai/development-specs/<type>/. Usá esta skill cuando el usuario pida "validar los specs", "revisar el comportamiento", "chequear coherencia entre capabilities", "¿mis specs se contradicen?", después de editar specs a mano o de correr development-spec-bootstrap. Lee `.matecito-ai/development-specs/` (recursivo por tipo) y reporta hallazgos con severidad. NO modifica nada — es consultiva.
---

# Development Spec Validate

Lee los capability-specs producidos por `development-spec-bootstrap` (o materializados por `sdd-archive`, o editados a mano) y los chequea contra una rúbrica: **completitud**, **coherencia entre capabilities**, **verificabilidad**, **referencias** e **integridad de la taxonomía**. Reporta hallazgos con severidad. No modifica archivos.

Los specs están organizados por tipo (`.matecito-ai/development-specs/<type>/<capability>.md`, `type` ∈ `flow` | `rule` | `lifecycle` | `process`), con un índice raíz (`.matecito-ai/development-specs/INDEX.md`) y un índice por tipo (`.matecito-ai/development-specs/<type>/INDEX.md`).

## Por qué contexto fresco

Esta validación SOLO sirve si es adversarial: leé únicamente lo que está ESCRITO en los specs. No asumas la intención del autor ni el comportamiento "obvio". Lo que no está en el archivo, no existe. Por eso esta skill corre con contexto limpio (standalone, o lanzada por el bootstrap como sub-agente).

## Concepto y tipos

El concepto canónico —qué ES/NO es un capability-spec, los tipos y sus secciones esqueleto— vive en `~/.claude/references/spec/README.md`. La skill lo **aplica**; no lo redefine.

Taxonomía de tipos **cerrada**: `flow` · `rule` · `lifecycle` · `process`. Cualquier carpeta bajo `.matecito-ai/development-specs/` que no sea uno de estos tipos es un hallazgo de integridad de taxonomía.

Secciones **esqueleto** por tipo (para el chequeo de completitud):

| `type` | Secciones esqueleto (deben tener contenido en un spec `Accepted`) |
|---|---|
| `flow` | Flujo principal · Errores de cara al actor · Escenarios |
| `rule` | Reglas de negocio · Escenarios |
| `lifecycle` | Entidades y estados · Escenarios |
| `process` | Flujo principal · Reglas de negocio · Escenarios |

## Pre-flight

Leé `.matecito-ai/development-specs/INDEX.md`. Si no existe, no hay nada que validar → sugerí correr `development-spec-bootstrap` y frená.

## Proceso

1. **Inventariá la estructura.** `find .matecito-ai/development-specs -name '*.md'`. Identificá el índice raíz, los índices de tipo (`<type>/INDEX.md`) y los specs (`<type>/<capability>.md`).
2. **Leé cada spec.** El tipo sale de la carpeta; el status del header.
   - **Specs con `Status: Draft`** NO son fuente de verdad cerrada: no reportes como defecto las secciones esqueleto o los escenarios faltantes (se espera que falten); sí verificá coherencia contra los `Accepted` (un Draft que ya contradice a un Accepted es un hallazgo).
   - **`Deprecated`:** verificá que tenga link a su reemplazo (si aplica) y que el reemplazo exista.
3. **Leé `coherence-rules.md`** (en esta misma skill) y aplicá cada chequeo.
4. **Emití el reporte** agrupado por tipo y, dentro de cada tipo, por severidad; cerrá con hallazgos cross-capability.

## Resolución de archivos

Un spec vive en `.matecito-ai/development-specs/<type>/<capability>.md`. Las contradicciones **entre capabilities** (que pueden estar en tipos distintos) requieren abrir varios archivos — recorré recursivo, no con glob plano.

## Formato del reporte

Agrupado **por tipo**, y dentro de cada tipo **por severidad**. Cerrá con una sección de hallazgos **cross-capability** (contradicciones que involucran más de un spec) y un veredicto final.

```
## Tipo: flow
🔴 CRITICAL — <qué> · specs: <cuáles> · <por qué> · <sugerencia>
🟡 WARNING  — ...
🔵 SUGGESTION — ...
(si un nivel no tiene hallazgos en el tipo, omitilo)

## Cross-capability
🔴/🟡/🔵 — hallazgos que cruzan capabilities (ej: una regla que contradice un flujo)

## Veredicto
<una línea: ej "1 CRITICAL, 2 WARNING — resolver el CRITICAL antes de codear">
```

Leyenda de severidad:

```
🔴 CRITICAL — contradicen el comportamiento entre sí; hay que resolverlas.
🟡 WARNING  — inconsistencias, huecos de verificabilidad o riesgo de pudrición.
🔵 SUGGESTION — mejoras de claridad o robustez.
```

Si un tipo no tiene ningún hallazgo, no lo listes. Si NO hay hallazgos en ningún lado, decilo explícitamente y dá un veredicto verde.

## Después del reporte

- **No modifiques specs.** Si el usuario quiere resolver un hallazgo, derivá a `development-spec-bootstrap` en modo update para el spec afectado.
- **Ratchet:** si detectaste una contradicción real que NO está en `coherence-rules.md`, ofrecé agregarla a la rúbrica (con su severidad y mensaje qué/por qué/sugerencia).

## Anti-patterns

- ❌ Inferir comportamiento no escrito para "salvar" una contradicción → si no está en el spec, es un hallazgo.
- ❌ Modificar o arreglar specs vos mismo → solo reportás; el usuario resuelve vía update.
- ❌ Reportar todo como CRITICAL → reservá CRITICAL para lo que contradice comportamiento; WARNING/SUGGESTION para el resto.
- ❌ Buscar specs con glob plano (`.matecito-ai/development-specs/*.md`) → están en subcarpetas de tipo; recorré recursivo.
- ❌ Ignorar carpetas no canónicas o índices desincronizados → son hallazgos de integridad, repórtalos.
- ❌ Tratar un `Draft` como si fuera `Accepted` para el chequeo de completitud → un Draft puede tener huecos por diseño.
