---
name: library-contract
depth: light
domain: contracts
tipo: decisión
adr-output: library-contract
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

ADR `library-contract` con: definición de la superficie pública (cómo se marca / dónde vive el index de exports), política de semver (qué cuenta como breaking change), ciclo de deprecación, y si existe changelog generado automáticamente (conventional commits + release tooling).
