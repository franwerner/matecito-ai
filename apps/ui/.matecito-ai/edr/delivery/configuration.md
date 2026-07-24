# EDR — Configuración

- **Status:** Accepted
- **Type:** convention
- **Date:** 2026-07-23

## Contexto

La UI se sirve embebida en el binario del broker (same-origin), así que en producción no hay URL de broker que configurar. La superficie de configuración es casi solo de desarrollo, y conviene validarla al arranque para fallar claro.

## Decisión

**`import.meta.env` de Vite** con variables prefijadas `VITE_`. Superficie mínima: en producción los endpoints son **relativos** (same-origin, el broker sirve la UI), así que no hay URL de broker que configurar; la configuración es casi solo de dev (URL del broker / target del proxy de Vite). **Validación al startup con Zod**: se parsea `import.meta.env` al boot y se aborta con mensaje claro si falta o está mal. El `.env` local va en `.gitignore`; el `.env.example` se commitea.

## Reglas verificables

- **[tool: zod]** la configuración se valida contra un schema Zod al arranque y aborta claro si falta o está mal.
- **[manual]** en producción los endpoints son relativos (same-origin).
- **[tool: gitignore]** el `.env` local está ignorado; el `.env.example` se commitea.

## Relacionados

- `relacionado-con` → [deployment-topology.md](deployment-topology.md) — el same-origin es lo que hace relativa la config en prod.
- `relacionado-con` → [../security/input-validation.md](../security/input-validation.md) — Zod también valida la config al boot.
