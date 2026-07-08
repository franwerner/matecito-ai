---
name: rate-limiting
depth: light
domain: security
type: decision
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

EDR `rate-limiting` materializado según `~/.claude/references/edr/templates/edr.md`. Esta es una decisión de tipo `policy`; sus reglas deben quedar especialmente accionables. Debe contener:

- **Contexto**: por qué sin límite un cliente puede agotar recursos o forzar credenciales por fuerza bruta, y cómo el punto de aplicación determina si la protección llega antes o después del código de la app.
- **Decisión**: granularidad elegida (IP, usuario autenticado, API key, o mix), punto de aplicación (API gateway/reverse proxy vs middleware de la app), y los límites concretos si se definieron.
- **Reglas verificables** (cada una con su mecanismo):
  - `[tool: test]` superado el límite decidido (ej. requests/minuto por tier), el sistema responde `429` con header `Retry-After`.
  - `[manual]` el límite se aplica con la granularidad elegida en el punto de aplicación decidido.
  - Para "Mix": `[manual]` los endpoints sensibles listados (login, reset de password) tienen el límite más estricto aplicado.
- **Alternativas consideradas**: las otras granularidades y puntos de aplicación evaluados y por qué no se eligieron.
- **Consecuencias**: lógica añadida en la app vs costo en el borde, y comportamiento esperado del cliente ante un `429`.
