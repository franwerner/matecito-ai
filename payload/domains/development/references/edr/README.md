# Qué es (y qué no es) un EDR

Referencia canónica del **concepto** de EDR (Engineering Decision Record). Es la fuente de verdad de la *idea*; cualquier skill o agente que trabaje con EDRs apunta acá en vez de redefinirla. La *estructura/plantilla* concreta se define por separado; esto define qué cuenta como EDR y qué no.

> **Por qué "Engineering" y no "Architecture".** El término estándar en la industria es "ADR" (*Architecture* Decision Record), pero acá el artefacto cubre **más que arquitectura**: decisiones, convenciones y políticas de ingeniería que gobiernan el código (desde el patrón de capas hasta "usá enums, no magic strings" o un límite de rate limiting). "Engineering" refleja ese alcance real y evita la lectura estrecha de "solo arquitectura". Lo que une a los tres sabores **no es la arquitectura** — es que son *elecciones deliberadas que gobiernan el código y se pueden chequear*. (El paralelo en el dominio de diseño es **DDR** — Design Decision Record.)

## Qué ES un EDR

Un EDR captura una **decisión** de ingeniería: una elección deliberada entre alternativas, con su razón. Tres rasgos lo definen:

1. **Decide** algo — elige una opción frente a otras. No describe ni ejecuta: elige.
2. **Tiene un porqué** — el contexto y el trade-off que justifican la elección.
3. **Perdura como restricción** — gobierna el código futuro; quien escribe código después la respeta, y se puede chequear si el código se desvió de ella.

Cubre tres sabores, todos "decisiones" en sentido amplio:
- **decisión** — un trade-off real (p. ej. qué motor de persistencia).
- **convención** — un acuerdo de estilo o estructura (p. ej. cómo se nombran los módulos).
- **política** — una regla verificable (p. ej. límites de rate limiting).

Responde: *"qué decidimos, por qué, y qué gobierna de acá en adelante"*.

## Qué NO es un EDR

- **No es una unidad de trabajo.** Una tarea o paso de implementación se ejecuta y se termina; no decide nada. La mayoría del trabajo es mecánico y no merece EDR. Forzar un EDR por cada cosa hecha es ruido.
- **No es una verificación.** Un criterio de aceptación o un test confirma que algo funciona ("input → resultado"). No tiene porqué ni alternativas, y no perdura como restricción. Verifica; no decide.
- **No es una señal ni un recordatorio.** Un "falta decidir X", un TODO o una nota *apuntan* a una decisión ausente — son el dedo señalando, no la decisión.
- **No es código incidental.** Un patrón que aparece pocas veces o sin intención clara no es una decisión, es ruido. Sin evidencia de elección deliberada, no hay EDR.
- **No es el "cómo".** El detalle de implementación vive en el código y sus comentarios. El EDR captura el *porqué* de una elección, no el paso a paso. En concreto, las secciones de **razonamiento** (Contexto, Decisión, Consecuencias, Alternativas) se escriben en términos de **conceptos, patrones y límites** — nunca nombrando **identificadores internos volátiles**: clases, métodos, columnas de base de datos, errores internos ni rutas de archivo concretas. Un identificador así en el razonamiento vuelve al EDR un calco del código que se pudre con el primer rename. Si hace falta anclar a algo concreto, va a la sección de anclaje/enforcement de la plantilla (un glob estable, o una regla verificable —que sí puede nombrar la clase, porque es el ancla que se chequea—), no al razonamiento. Excepción: el nombre de una tecnología/librería (que ES la decisión) y el contrato público de cara al consumidor (endpoints públicos, códigos de error expuestos).
- **No es un porqué adivinado.** Inferir la razón de una decisión sin que conste es inventar. El porqué lo aporta una persona.

## Estados (un EDR no es estático)

Un EDR tiene ciclo de vida. Dos estados importan para el concepto:

- **Borrador / inferido** — se registró el **QUÉ** (qué se eligió, observable) pero el **PORQUÉ está vacío**: es *evidencia de que se tomó una decisión*, sin la razón ratificada. Es un candidato, no una decisión cerrada.
- **Aceptado** — una persona lo ratificó y completó el porqué (contexto, alternativas, consecuencias). Recién ahí es un EDR pleno.

La promoción de borrador a aceptado siempre la hace una persona; nunca es automática. Un borrador todavía no "decide" — solo deja constancia de que hay algo para decidir.
