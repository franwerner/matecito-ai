---
name: testing-strategy
depth: deep
domain: delivery
type: decision
source: práctica clásica de pirámide de tests · arc42 §8 (conceptos transversales)
---

# Fase: Estrategia de testing

## Qué decide

La pirámide objetivo de tests, política de mocks vs reales, si TDD es obligatorio, y la cobertura mínima si se mide. Define qué tan costoso es cambiar el sistema y qué nivel de confianza da la suite.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. Pirámide de tests objetivo

> La proporción define dónde está la mayoría del costo de mantenimiento y qué tan rápido corre la suite.

- ***Pirámide clásica — default recomendado:*** mayoría unit (70%), integración (20%), e2e/acceptance (10%).
- **Pirámide invertida** — mayoría integración/e2e. Más confianza, suite más lenta y frágil.
- **Solo integración** — sin unit tests del dominio. Ok para CRUD sin lógica de negocio real.
- **Sin estrategia formal** — tests donde surge la necesidad.
- No sé, recomendame.

### 2. Mocks vs fakes vs reales

> Usar mocks en exceso produce tests que pasan pero el sistema falla en producción. Usar reales hace la suite lenta o no reproducible.

- ***Mocks para dependencias externas al proceso (DB, HTTP, filesystem); reales para lógica de dominio — default recomendado.***
- **Fakes in-memory** (repositorios en memoria, DB en memoria) — más realistas que mocks, más rápidos que reales.
- **Testcontainers para DB real** — containers efímeros de DB en tests de integración. Alta confianza, requiere Docker.
- **Mocks para todo** — rápido de escribir, frágil con el tiempo.
- No sé, recomendame.

### 3. TDD obligatorio

> TDD no es solo una práctica de testing — cambia cómo se diseña la API pública de cada componente.

- **Sí, obligatorio** — todo código nuevo va precedido de test rojo.
- ***No obligatorio, pero sí recomendado para lógica de dominio y casos de uso — default razonable.***
- **No, tests se escriben después** — válido para prototipos o cuando la interfaz cambia mucho.
- No sé, recomendame.

### 4. Cobertura mínima

> Una cobertura mínima sin contexto es una métrica engañosa. Solo tiene valor si se acuerda qué medir y qué no.

- **Sin umbral formal** — la calidad importa más que el porcentaje.
- ***Umbral razonable: 80% líneas en `domain/` y `application/`; sin umbral en `infrastructure/` — default si se quiere medir algo.***
- **Umbral global** (ej: 80% de líneas en todo el proyecto).
- No sé, recomendame.

## Tech a registrar

Framework de tests (`pytest.md`, `vitest.md`, `jest.md`, `junit.md`, `rspec.md`, `go-test.md`), librería de mocking si es separada del framework (`unittest.mock` es built-in de Python y no necesita registro; `mockito.md` para Java; `testify.md` para Go si se usa), `testcontainers.md` si se eligió para integración.

## Qué materializar

ADR `testing-strategy` materializado según `../../templates/adr.md`. Debe contener:

- **Contexto** y **Decisión**: proporciones de la pirámide objetivo (unit / integración / e2e, ej: 70/20/10), política de mocks vs fakes vs reales con criterio claro (qué se mockea y qué no y por qué), si TDD es obligatorio o recomendado y para qué tipo de código, y el umbral de cobertura si se acordó (sobre qué capas y con qué métrica).
- **Reglas verificables**: cada política como aserción con su mecanismo al inicio. Ej: `[manual]` nunca mockear clases internas del dominio; `[manual]` siempre mockear llamadas HTTP salientes y acceso a DB en unit tests; `[tool: <coverage tool>]` cobertura mínima 80% de líneas en `domain/**` y `application/**`, sin umbral en `infrastructure/**`; `[manual]` todo código nuevo de dominio/casos de uso va precedido de test rojo si TDD es obligatorio. Nombrá el framework de tests y la herramienta de cobertura concretos en los `[tool: ...]`. Conservá los porcentajes de la pirámide y el umbral por capa.
- **Relacionados** (opcional): vinculá con `ci-quality-gates` (donde tests y cobertura corren como gate).
