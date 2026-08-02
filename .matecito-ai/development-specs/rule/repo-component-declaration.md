# Capability — Declaración del set de componentes del repo

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Permitir que un repositorio con múltiples superficies (e.g., CLI, API, UI) declare sus componentes en la configuración, de modo que cada capability-spec pueda indicar qué superficies implementan ese comportamiento.

## Reglas de negocio

- El autor del proyecto declara el bloque `repo` en el config **de proyecto**, al mismo nivel que `domains` y `domainConfig`; una declaración global no se hereda ni se mezcla — es la única excepción explícita a la precedencia global → proyecto. Requiere que el repositorio tenga un archivo `.matecito-ai/config.json`.
- El bloque `repo` contiene una lista `components`; cada componente declara `name` (la superficie que el consumidor reconoce) y `paths` (uno o más directorios donde vive). El parser de configuración carga esta lista y resuelve el set de componentes del proyecto.
- Los nombres de componentes **dentro del set deben ser únicos** — no hay dos componentes con el mismo `name`. Si se declaran duplicados, el sistema rechaza el config y reporta qué nombres se repiten.
- El set se escribe **una sola vez** y se ratifica explícitamente por el usuario; no se aplica automáticamente.
- **Gate presence-based:** sin `repo.components` declarado, el eje no existe en ningún lado del sistema — sin líneas en headers, sin findings en validación, sin propuestas en minería.

## Entidades y estados

- **Componente** — una superficie reconocida por el consumidor (e.g., `cli`, `api`, `ui`). Cada componente tiene un nombre único y una lista de directorios (puede ser más de uno, ej. `cli` puede estar en `cmd/` e `internal/`).

## Escenarios

### Scenario: el global no aporta componentes

- **GIVEN** un config global con `repo.components` y uno de proyecto sin él
- **WHEN** se carga la configuración del repo
- **THEN** el set queda vacío — el global no se hereda

### Scenario: el bloque `repo` no se pliega

- **GIVEN** un config con `repo` y un flag heredado, ambos al tope
- **WHEN** se carga
- **THEN** el flag queda bajo el dominio por defecto y `repo` sigue al tope, intacto — no participa en la normalización que pliega otros bloques

### Scenario: sin declaración, el eje no existe

- **GIVEN** un repo sin `repo.components`
- **WHEN** se autoriza, valida o mina specs
- **THEN** todo se comporta como antes del eje: sin línea en el header y sin findings del eje

## Referencias

- **Concepto** → `~/.claude/references/spec/README.md` — el capability-spec como fuente de verdad del comportamiento.
