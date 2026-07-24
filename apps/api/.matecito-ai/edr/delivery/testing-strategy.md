# EDR — Estrategia de testing

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

El corazón del broker es la persistencia y el tail del transcript; la mayor parte del valor y del riesgo vive en esa interacción, no en lógica pura aislada. Por eso la capa de integración pesa más de lo típico. Como la persistencia es SQLite embebida y el tail es sobre archivos, se los puede ejercitar con recursos reales sin infraestructura externa, lo que hace innecesario mockearlos o levantar contenedores.

## Decisión

Adoptamos una pirámide clásica pragmática con la capa de integración más pesada de lo habitual, porque el store es el corazón: unit para la lógica pura (derivar estado del event-log, shaping del envelope, parseo del tail, validación de config) e integración para el store y el transporte; e2e mínimo u opcional. Usamos recursos reales, no mocks: SQLite real en tests (temporal o en memoria), sin mockear el store; el tail se testea contra archivos temporales reales; sin contenedores ni infraestructura externa; sin mockear interfaces internas del dominio. TDD no es obligatorio (sí recomendado para la derivación de estado y el parseo del tail incremental). La cobertura se mide de forma informativa, sin umbral formal ni gate. El framework es el de testing de la librería estándar (tests table-driven), sin librería de assertions externa.

## Reglas verificables

- **[manual]** el store se testea con SQLite real (temporal o en memoria), nunca mockeado.
- **[manual]** el tail se testea con archivos temporales reales.
- **[manual]** la suite no depende de contenedores ni de infraestructura externa (los tests corren solos).
- **[manual]** no se mockean las interfaces internas del dominio.

## Alternativas consideradas

- **Mockear el store y el tail:** descartado; ejercitar SQLite y archivos reales es barato y da confianza real sobre el corazón del sistema.
- **Contenedores / testcontainers:** descartados; la persistencia embebida no necesita infraestructura externa para testearse.
- **Umbral de cobertura como gate:** descartado; la cobertura se usa informativa, no como barrera.

## Consecuencias

- Los tests de integración dan confianza directa sobre el store y el tail, que es donde está el riesgo.
- La suite corre sola, sin dependencias de infraestructura.
- Trade-off: los tests de integración con recursos reales son más lentos que los unit puros.

## Relacionados

- `relacionado-con` → [ci-quality-gates.md](ci-quality-gates.md) — el gate de tests en CI se define ahí.
