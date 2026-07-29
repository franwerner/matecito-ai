# EDR — Documentación

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

Proyecto local single-dev: la documentación necesaria es la mínima para levantar y buildear la UI, más las decisiones ya registradas. La doc de API no se escribe a mano porque el broker la genera.

## Decisión

`README.md` en `apps/ui` (qué es, cómo levantarlo, cómo buildear) más los EDRs de este store. La **doc de API es el OpenAPI generado del broker** (generado desde código, siempre en sync). Sin runbooks (es local).

## Reglas verificables

- **[manual]** el README describe qué es el proyecto, cómo levantarlo y cómo buildear.
- **[manual]** la doc de API es el OpenAPI generado del broker.

## Relacionados

- `relacionado-con` → [ci-quality-gates.md](ci-quality-gates.md) — el sync de los tipos generados desde ese OpenAPI.
