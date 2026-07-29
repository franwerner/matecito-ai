# EDR — Manejo de errores

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

La UI depende de un broker que puede fallar o caerse, y de render de árboles densos. Hace falta una estrategia que distinga lo inesperado (crashes de render) de lo esperado (errores de data), que mapee el error público del broker a algo user-facing y que haga visible el estado de la conexión.

## Decisión

Mix pragmático:
- excepciones para lo inesperado (crashes de render),
- errores esperados de data como error-state de TanStack Query,
- validación como valores.

Boundary: **React Error Boundaries en capas** — uno top-level en el app shell más un `errorComponent` por ruta de TanStack Router. Los errores de Query se muestran inline o se lanzan al boundary (`throwOnError`) según criticidad.

El contrato de error público del broker (`{error, code}`) se mapea a un tipo del cliente y a un mensaje user-facing; nunca se muestra crudo. Un tipo aparte modela el **estado de conexión WS** (caído / reconectando / stale), visible al usuario. Presentación: toasts (`sonner`) para transitorios, inline para data. Sin logging formal del cliente (ese concern quedó fuera; en dev, consola).

## Reglas verificables

- **[manual]** todo árbol de UI está cubierto por un error boundary (shell + por ruta); un crash de render no tumba la app.
- **[manual]** los errores esperados de data se manejan como estado de Query.
- **[manual]** el `{error, code}` del broker se mapea a un mensaje user-facing, nunca crudo.
- **[manual]** el estado de conexión WS es visible al usuario.
- **[manual]** nunca se vuelca a consola ni a la UI tokens ni payloads sensibles.

## Relacionados

- `relacionado-con` → [../frontend/data-fetching.md](../frontend/data-fetching.md) — los errores de data como estado de Query.
- `relacionado-con` → [resilience.md](resilience.md) — el estado de conexión y el resync tras la caída.
