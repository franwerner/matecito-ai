---
name: rate-limiting
depth: light
domain: security
tipo: decisión
adr-output: rate-limiting
source: OWASP ASVS v4 §13.1 (Generic Web Service Security)
---

# Fase: Rate limiting

## Qué decide

Si se limita la cantidad de requests, con qué granularidad (IP / usuario / API key) y dónde se aplica el límite.

## Preguntas

### 1. Granularidad y punto de aplicación

> Sin límite, un cliente puede agotar los recursos del sistema o forzar credenciales por fuerza bruta. El punto de aplicación determina si la protección llega antes o después de tu código.

- **Por IP en el API gateway o reverse proxy (nginx, Traefik, Kong)** — *default para APIs públicas; bajo costo, sin lógica en la app.*
- Por usuario autenticado en middleware de la app.
- Por API key en middleware de la app.
- Mix: IP en el borde + usuario en la app para endpoints sensibles (login, reset de password).
- No aplica — la API es interna / sin acceso externo directo.
- No sé, recomendame.

## Notas de lógica (para el motor)

- Si el proyecto es `cli`, `libreria` o `script`, marcar esta fase como `Not Applicable`.
- Si se elige "Mix", pedir que el usuario liste los endpoints que necesitan límite más estricto.

## Tech a registrar

Si se elige una librería en la app (ej: `express-rate-limit`, `slowapi`, `throttler`), registrarla en el catálogo `tech/`.

## Qué materializar

ADR `rate-limiting` con: granularidad elegida, punto de aplicación, límites concretos si se definieron (requests/minuto por tier), y respuesta al cliente al superar el límite (429 + `Retry-After`).
