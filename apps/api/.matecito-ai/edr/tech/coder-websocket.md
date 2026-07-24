# coder/websocket

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** contracts
- **Date:** 2026-07-23

## Por qué la elegimos

Librería de WebSocket para la superficie WS-out (el broker hace fan-out de eventos a la UI). Su handshake valida Origin contra Host por default, lo que en mismo-origen (la UI se sirve desde el broker) cubre el caso sin configuración extra.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: contracts/api-contract. La validación Origin==Host por default es lo que sostiene la decisión de no aplicar CORS (ver la fila `cors` en el INDEX del dominio security).
