---
name: resilience
depth: light
domain: runtime
tipo: decisión
adr-output: resilience
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

ADR `resilience` con: política elegida, valores concretos si se definieron (timeout en ms, max retries, backoff base, umbral del circuit breaker), y librería si aplica. Las reglas deben ser verificables ("timeout de 5s en todas las llamadas HTTP salientes").
