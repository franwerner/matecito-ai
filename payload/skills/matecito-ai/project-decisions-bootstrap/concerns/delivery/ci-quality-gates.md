---
name: ci-quality-gates
depth: light
domain: delivery
tipo: política
adr-output: ci-quality-gates
source: práctica de CI/CD / checklists de production-readiness
---

# Fase: Quality gates de CI

## Qué decide

Qué checks corren automáticamente en CI y cuáles bloquean el merge. Incluye pre-commit hooks para fallas rápidas en local.

## Preguntas

### 1. Checks que bloquean el merge

> Cada check ausente es una categoría de problema que puede entrar silenciosamente a main.

Marcá los que querés que bloqueen el merge (podés elegir varios; el default recomendado está en *cursiva*):

- *Linter (análisis estático)* — ej: ESLint, Ruff, golangci-lint, Checkstyle.
- *Formateo automático* — ej: Prettier, Black, gofmt. Falla si hay diff.
- *Type-check* — ej: `tsc --noEmit`, `mypy`, `pyright`. Solo si el stack tiene tipos.
- *Tests* — falla si algún test falla (ver `testing-strategy` para el nivel de cobertura).
- *Cobertura mínima* — falla si baja de un umbral; requiere que esté definido en `testing-strategy`.
- *Enforcement de arquitectura* — solo si `arch-enforcement` está `Accepted`.
- Ninguno por ahora — CI existe pero no bloquea.
- No sé, recomendame.

### 2. Pre-commit hooks

> Correr los mismos checks en local antes del push reduce round-trips con CI.

- **Sí, con `pre-commit` framework** — *default; config en `.pre-commit-config.yaml`, mismos checks que CI.*
- Sí, scripts propios (Makefile / shell).
- No, solo CI.

## Tech a registrar

Si se usa `pre-commit` framework, registrarlo en `tech/`. Si se elige un linter o formatter específico que no estaba ya registrado, registrarlo también.

## Qué materializar

ADR `ci-quality-gates` con: lista de checks que bloquean el merge, herramientas para cada uno, si hay pre-commit y con qué herramienta, y la regla explícita "nada llega a main sin pasar X, Y, Z".
