# EDR — Logging

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

El broker es un daemon local de larga vida; su diagnóstico se hace tailando su salida o, eventualmente, consumiéndola en JSON. Necesita logging estructurado sin sumar dependencias externas, legible cuando se lo mira a mano, y con correlación por identidad de dominio (proyecto, change, agente) en lugar de un request-id HTTP, que no captura la unidad de trabajo real. Además no debe filtrar contenido sensible.

## Decisión

Usamos el logger estructurado de la librería estándar (sin dependencia externa), con un handler de texto por default (legible al tailear el daemon) conmutable a JSON por configuración. Los niveles son debug/info/warn/error; el nivel mínimo es configurable, con default info. La correlación es por identidad de dominio, no por request-id HTTP: cada log lleva el proyecto, el change y el agente del contexto, adjuntados en el borde de transporte.

## Reglas verificables

- **[manual]** se usa el logger estructurado de la stdlib con formato configurable texto/JSON.
- **[manual]** nunca se loggean payloads de artefactos, contenido de transcript ni datos personales, en ningún nivel.
- **[manual]** todo error loggeado incluye su cadena de wrapping.
- **[manual]** los logs con contexto de dominio llevan el identificador de correlación (proyecto + change + agente) adjuntado en el borde.

## Alternativas consideradas

- **Librería de logging externa:** descartada; el logger estructurado de la stdlib cubre el caso sin sumar dependencias.
- **Correlación por request-id HTTP:** descartada; no representa la unidad de trabajo real, que es la identidad de dominio.

## Consecuencias

- Diagnóstico legible a mano por default y consumible en JSON cuando haga falta.
- La correlación por identidad de dominio permite seguir el trabajo de un change de punta a punta.
- Trade-off: adjuntar el contexto de dominio en el borde exige que ese contexto esté disponible ahí.

## Relacionados

- `relacionado-con` → [../runtime/error-handling.md](../runtime/error-handling.md) — la política de errores define qué se loggea y qué nunca se vuelca.
