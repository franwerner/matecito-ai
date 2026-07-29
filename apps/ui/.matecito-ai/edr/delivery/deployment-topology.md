# EDR — Topología de despliegue

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

El broker es un daemon local en Go que ya sirve la UI. No hay ni se quiere un host separado para el frontend; la SPA es stateless del lado cliente (su estado es efímero o derivado del broker).

## Decisión

La UI se distribuye como **build estático** (Vite → `apps/ui/dist`) **embebido en el binario del broker** vía `embed.FS`. Sin host separado (nada de CDN / Vercel / Netlify): el broker sirve los assets. El **build de la UI corre antes del build de Go** en la release (hook de goreleaser). En desarrollo: Vite dev server con **proxy al broker** (same-origin desde el browser). La SPA es stateless del lado cliente y se sirve desde la instancia única global del broker.

## Reglas verificables

- **[manual]** la UI se distribuye como build estático embebido en el binario del broker (embed.FS), no como host separado.
- **[manual]** el build de la UI corre antes del build de Go en la release.
- **[manual]** en desarrollo la UI corre en Vite dev server con proxy al broker.

## Relacionados

- `relacionado-con` → [configuration.md](configuration.md) — el same-origin hace relativa la config en prod.
- `relacionado-con` → [../frontend/styling.md](../frontend/styling.md) — el build estático que se embebe incluye los assets de estilo.
