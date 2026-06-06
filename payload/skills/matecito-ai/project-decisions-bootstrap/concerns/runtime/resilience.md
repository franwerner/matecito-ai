---
name: resilience
depth: light
domain: runtime
type: decision
source: SRE (site reliability engineering)
---

# Fase: Resiliencia — timeouts, retries y circuit breakers

## Qué decide

Cómo el sistema se comporta cuando una dependencia externa (DB, API tercera, cola) es lenta o falla. Sin política explícita, el default es "esperar indefinidamente y propagar el error".

## Preguntas

Una o dos, según haga falta.

### 1. Política de resiliencia en llamadas externas

> Una llamada sin timeout puede bloquear un thread/goroutine indefinidamente y agotar el pool. Retries sin backoff amplifican la carga sobre un servicio caído.

- **Timeout + retry con backoff exponencial** — *default recomendado para cualquier llamada externa.*
- Solo timeout, sin retry — simple; aceptable si la operación no es idempotente o el cliente reintenta.
- Timeout + retry + circuit breaker — para servicios críticos donde la falla en cascada es el riesgo principal.
- Sin política formal por ahora — *solo si el proyecto no tiene dependencias externas en este momento.*
- No sé, recomendame.

### 2. Implementación del circuit breaker

> El circuit breaker evita llamar a un servicio que ya se sabe caído. **Solo si en la 1 eligió circuit breaker.**

- Librería dedicada (resilience4j, opossum, pybreaker, go-resilience) — *default recomendado.*
- Implementación custom — solo si la librería no existe para el stack o hay restricciones de dependencias.

## Notas de lógica (para el motor)

- Si eligió "Sin política formal por ahora", materializá con `Status: Pending` y razón ("proyecto sin dependencias externas todavía; revisar cuando se integre la primera").
- Si eligió circuit breaker, hacer pregunta 2.

## Tech a registrar

Si se elige una librería de resiliencia concreta, registrarla en el catálogo `tech/`.

## Qué materializar

ADR `resilience` materializado según el template `~/.claude/references/adr/templates/adr.md`. La **Decisión** captura: política elegida (solo timeout / timeout+retry con backoff / timeout+retry+circuit breaker / sin política), los valores concretos si se definieron (timeout en ms, max retries, backoff base, umbral del circuit breaker) y la librería de resiliencia si aplica.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** toda llamada externa saliente (DB, API tercera, cola) tiene un timeout explícito; ninguna espera indefinidamente. Ej: timeout de 5s en todas las llamadas HTTP salientes.
- **[manual]** los retries usan backoff exponencial con un tope de reintentos definido; no hay retries fijos sin backoff que amplifiquen la carga sobre un servicio caído.
- **[manual]** si se eligió circuit breaker: las llamadas a servicios críticos pasan por el breaker con el umbral definido, no lo bypassean.

Si se eligió "Sin política formal por ahora", el ADR va con `Status: Pending` y la razón concreta ("proyecto sin dependencias externas todavía; revisar cuando se integre la primera"); en ese caso no lleva Reglas verificables.

**Relacionados:** vincular con `background-jobs` si los reintentos de jobs comparten la misma política de backoff.
