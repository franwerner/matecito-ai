---
name: library-contract
depth: light
domain: contracts
type: decision
source: arc42 §8 (cross-cutting concepts) · semver.org
---

# Fase: Contrato de librería

## Qué decide

Qué expone la librería como superficie pública, cómo se versiona, y cuál es la política de backward compatibility y deprecación.

## Preguntas

### 1. Superficie pública y semver

> Sin una definición explícita de la API pública, cualquier cambio interno puede romper a los consumidores. Semver sin disciplina pierde su significado como contrato.

- **API pública declarada explícitamente (exports/index) + semver estricto** — *default recomendado: major = breaking, minor = additive, patch = fix.*
- Todo lo que no es `_private` o `internal` es público — solo para librerías internas de un monorepo.
- No sé, recomendame.

### 2. Política de deprecación

> Los consumidores necesitan tiempo para migrar. Eliminar sin deprecación previa es un breaking change encubierto.

- **Deprecation warning en minor + eliminación en el siguiente major** — *default recomendado.*
- Sin política formal — solo para librerías internas con consumidores conocidos y bajo control.

## Qué materializar

ADR `library-contract` materializado según `~/.claude/references/adr/templates/adr.md`. Debe contener:

- **Contexto**: por qué sin una API pública declarada cualquier cambio interno puede romper consumidores, y por qué semver sin disciplina pierde su significado como contrato.
- **Decisión**: definición de la superficie pública (cómo se marca / dónde vive el index de exports), política de semver (qué cuenta como breaking change: major = breaking, minor = additive, patch = fix), ciclo de deprecación, y si existe changelog generado automáticamente (conventional commits + release tooling).
- **Reglas verificables** (cada una con su mecanismo):
  - `[manual]` solo lo declarado en la superficie pública (exports/index) se considera contrato; lo demás es interno y puede cambiar sin major.
  - `[manual]` un breaking change incrementa la versión major; lo additive, minor; el fix, patch.
  - `[manual]` toda eliminación pasa primero por un deprecation warning en un minor antes de removerse en el siguiente major.
  - Si hay tooling: `[tool: <release tooling>]` el changelog se genera automáticamente a partir de conventional commits.
- **Alternativas consideradas**: tratar todo lo no `_private`/`internal` como público y la ausencia de política de deprecación, con su trade-off para librerías internas vs públicas.
- **Consecuencias**: disciplina requerida en cada release y previsibilidad que ganan los consumidores.
