---
name: documentation
depth: light
domain: delivery
type: convention
source: arc42 §1 / práctica de engineering documentation
---

# Fase: Documentación

## Qué decide

Qué se documenta, dónde vive cada tipo de doc, y qué convenciones se siguen. El objetivo es que la documentación sea mantenible, no exhaustiva.

## Preguntas

### 1. Qué se documenta y dónde

> Sin una convención explícita, la documentación se dispersa o directamente no existe. El objetivo no es documentar todo — es que lo que importa tenga un lugar definido.

Elegí qué tipos de doc aplican al proyecto (podés combinar):

- **`README.md` en raíz** — *default para todos; describe qué es el proyecto, cómo levantarlo, cómo correr tests, cómo contribuir.*
- **ADRs en `.matecito-ai/adr/`** — decisiones arquitectónicas; ya cubierto por esta skill.
- **Docs de API** — *recomendado si hay API pública o consumida por terceros*; formato: OpenAPI / Swagger, GraphQL schema, doc generada desde código.
- **Docs de módulos / librerías** — docstrings / JSDoc con generador (Sphinx, TypeDoc, Javadoc, etc.); solo si es una librería pública o con múltiples consumidores.
- **Runbooks / docs operacionales** — cómo deployar, cómo hacer rollback, cómo resolver alertas comunes.
- Solo README, nada más por ahora.
- No sé, recomendame.

### 2. Dónde vive la documentación de API

> **Solo si eligió "Docs de API" arriba.** Definir esto evita que cada dev elija un formato distinto.

- **Generada desde código** (decoradores, anotaciones, comentarios) — fuente de verdad en el código; siempre sincronizada.
- Archivo estático versionado (`openapi.yaml` / `schema.graphql` commitado) — explícito, revisable en PRs.
- Plataforma externa (Confluence, Notion, Postman) — fácil de editar, difícil de mantener sincronizada.

## Qué materializar

ADR `documentation` materializado según `~/.claude/references/adr/templates/adr.md`. Debe contener:

- **Contexto** y **Decisión**: qué tipos de doc se mantienen (README, ADRs, docs de API, docs de módulos, runbooks), dónde vive cada una, y el formato de API docs si aplica (generada desde código / archivo estático versionado / plataforma externa).
- **Reglas verificables**: las convenciones de mantenimiento como aserciones con su mecanismo al inicio. Ej: `[tool: <CI doc check>]` `openapi.yaml` se regenera y el PR falla si queda desincronizado del código; `[manual]` cualquier cambio de interfaz pública actualiza el doc de API en el mismo PR; `[manual]` el `README.md` en raíz describe qué es el proyecto, cómo levantarlo y cómo correr tests. Conservá el detalle de qué tipo de doc vive dónde.
- **Relacionados** (opcional): vinculá con `ci-quality-gates` si la verificación de docs corre como gate.
