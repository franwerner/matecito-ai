# EDR — El kernel afirma el invariante de discovery, no el mecanismo que lo resuelve

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La regla del invariante de discovery decía antes "asks the discovery form" — un CÓMO. Aplicado a una fase headless, ese cómo se traducía en que el agente se contestara su propio formulario, porque no tiene canal para preguntarle al usuario. El slot y el invariante (que el formulario se resuelva CON el usuario) son del kernel; el mecanismo concreto de cómo se resuelve es del dominio.

## Decisión
El kernel afirma sólo el invariante — el formulario de discovery se resuelve con el usuario antes de despachar la fase que fija el alcance del cambio — y deja explícitamente el CÓMO (y qué fase lo posee) como decisión del fragmento de dominio. `development`, por ejemplo, corre un ciclo de dos pasadas `needs-input` a través de su Discovery Gate.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `## Structured Flow` → "Discovery invariant", afirma que el formulario se resuelve CON el usuario, sin especificar el mecanismo concreto de cómo se resuelve.
- **[manual]** La misma sección declara explícitamente que el CÓMO y qué fase lo posee es decisión del fragmento de dominio.

## Alternativas consideradas
Mantener "asks the discovery form" como instrucción genérica del kernel — descartada: una fase headless no tiene canal para preguntar, así que la instrucción sólo puede cumplirse contestándose a sí misma, exactamente lo que el invariante existe para prohibir.

## Consecuencias
El kernel queda más corto y más correcto para cualquier dominio cuyo mecanismo de discovery no sea un ciclo needs-input de dos pasadas.
