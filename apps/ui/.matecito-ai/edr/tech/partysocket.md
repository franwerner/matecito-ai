# partysocket

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** frontend
- **Date:** 2026-07-23

## Por qué la elegimos

Cliente WebSocket único hacia el broker: auto-reconnect con backoff, buffer mientras está caído y heartbeat de fábrica. Bridgea los eventos `{type, payload}` a la cache de Query y al store Zustand.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: frontend/data-fetching (cliente WS único) y runtime/resilience (reconnect + backoff + buffer + heartbeat).
