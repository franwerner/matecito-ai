# EDR — Ruteo

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

El cockpit es read-first: un índice de changes y el detalle de un change (canvas + inspector + timeline). Buena parte del estado de vista debe ser remontable por URL (deep-link), y el snapshot inicial conviene prefetchearlo por ruta.

## Decisión

**TanStack Router file-based** (directorio de rutas más el plugin de Vite `@tanstack/router-plugin`); las rutas son finas e importan de las features. Rutas read-first: índice de changes → detalle de un change (canvas + inspector + timeline). El estado remontable (lente/filtro activo, nodo seleccionado, posición del scrubber) va en **search params tipados**. Los loaders de ruta prefetchean el snapshot vía Query (`ensureQueryData`); Query sigue siendo la cache.

## Alcance

- `apps/ui/src/routes/**` — rutas file-based finas que importan de las features.

## Reglas verificables

- **[manual]** las rutas son file-based y finas, e importan de las features.
- **[manual]** el estado deep-linkable va en search params tipados.
- **[manual]** los loaders prefetchean el snapshot vía Query.
