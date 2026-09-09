# EDR — El alcance de la excepción de emisión se nombra, en vez de quedar posicional

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
El párrafo que garantiza que el contenido mandado por otra regla se emite entero decía que "la prosa alrededor de ese material" queda fuera del presupuesto de brevedad. "Alrededor" es una posición, no un alcance: se lee como permiso para un preámbulo propio, justo lo que la misma sección prohíbe tres párrafos antes en su lista `Never`. La excepción terminaba autorizando por una preposición lo que la regla general prohibía por nombre.

## Decisión
El párrafo nombra exactamente qué puede cubrir la prosa propia del asistente cuando una forma cerrada aplica: la única línea que nombra qué se está mostrando, y la respuesta a una pregunta que el usuario hizo sobre ese contenido. Se nombra también, explícitamente, lo que NO queda cubierto: preámbulo, oración de encuadre, restablecimiento de por qué importa, o cierre recapitulativo. El resto del párrafo — que el contenido mandado se emite entero, y el ejemplo del split summary/rationale — queda sin cambios.

## Reglas verificables
- **[manual]** El párrafo de excepción de emisión en `payload/core/CLAUDE.md`, sección `### Explaining — fixed forms, not free prose`, nombra las dos cosas que la prosa propia puede cubrir (la línea que nombra qué se muestra; la respuesta a una pregunta del usuario) en vez de decir "la prosa alrededor de ese material".
- **[manual]** El mismo párrafo nombra explícitamente lo que queda fuera: preámbulo, encuadre, restablecimiento de por qué importa, cierre recapitulativo.
- **[manual]** El contenido que otra regla manda emitir sigue emitiéndose entero — incluido el split summary/rationale — sin que el alcance recortado de esta excepción lo alcance.

## Alternativas consideradas
Dejar la oración de alcance y confiar en que la lista `Never` la acota — descartada: el trabajo del párrafo es decir que el presupuesto de brevedad no anula la emisión, y una excepción sin límite explícito se lee como "nada acá está presupuestado". Reformular a "prosa que precede a ese material" — descartada: sigue siendo posicional, el mismo defecto una preposición más adelante.

## Consecuencias
La excepción deja de poder leerse como licencia para un preámbulo propio. El costo es que el párrafo nombra explícitamente cuatro formas prohibidas que antes quedaban implícitas por la lista `Never` de más arriba — una repetición acotada, no una regla nueva.

## Relacionados
- `relacionado-con` → [between-phase-report-act-form.md](between-phase-report-act-form.md) — el cuarto acto que esta excepción ya no puede usar como puerta trasera para un preámbulo propio.
