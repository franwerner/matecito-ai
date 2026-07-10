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

## Dónde va cada nombre — el test del identificador volátil

La regla "el razonamiento va en conceptos" no prohíbe nombrar lo concreto; regula **dónde** va cada nombre. Un identificador (clase, método, columna, error interno, ruta, enum interno) tiene lugar en el EDR, pero en la **sección-ancla** (`## Alcance` con un glob estable, o `## Reglas verificables` donde el nombre es lo que se chequea), **nunca** en la prosa del razonamiento (Contexto / Decisión / Consecuencias / Alternativas). La diferencia es funcional: el ancla se chequea y avisa cuando el código se movió; el nombre suelto en la prosa solo se pudre callado con el primer rename.

**El test (3 preguntas).** Un nombre concreto puede quedar en el razonamiento solo si pasa las tres. Si falla alguna, es volátil → reubicalo al ancla.

1. **¿Lo nombra un consumidor externo?** Endpoint público, header del contrato, código de error expuesto, nombre de tecnología/librería. Si es puro interno, falla.
2. **¿Sobrevive a un rename interno sin que el EDR mienta?** Si renombrar la columna/método/carpeta mañana vuelve falsa una frase del razonamiento, falla.
3. **¿Es la decisión misma o solo cómo se implementó?** El nombre de la tecnología *es* la decisión (pasa). El método/columna/clase/ruta es el *cómo* (falla).

**Traducción — el mismo hecho, mal y bien:**

| Sección | ❌ Calco del código (en el razonamiento) | ✅ Concepto en el razonamiento + ancla aparte |
|---|---|---|
| Decisión | "las transiciones ocurren vía `markActive()` / `markDone(resultId)` / `markFailed(reason)`" | Razonamiento: "las transiciones de estado solo ocurren vía métodos de la entidad, nunca seteando el estado a mano." — y en `## Reglas verificables`: **[manual]** ningún repo/use case setea el estado directo. |
| Decisión | "tabla `orders` con `id`, `tenantId` FK, `ownerId` FK, `refId` nullable, `failureReason`" | Razonamiento: "se persiste la operación con audit trail: identidad, dueño, referencia opcional y motivo de fallo." — el modelo concreto vive en el código; si hace falta anclar, la regla verificable nombra la columna. |
| Decisión | "jerarquía `BaseError` → `ExpiredError`, `NotAllowedError` en `src/features/*/domain/errors/`" | Razonamiento: "errores de negocio como jerarquía propia; los del proveedor externo se traducen en el borde del adapter." — y en `## Alcance`: `src/**/domain/errors/**`. |
| Contexto | "el cambio `feature-x` (Fase 2 del roadmap); slice `modulo-saliente`" | "se agrega una operación remota de escritura que puede fallar transitoria o permanentemente." (el condicionante conceptual, no el nombre del ticket/slice) |

**Dos formas prohibidas** (mismo mal: el razonamiento calca la implementación, así que cada cambio de código obliga a parchear la prosa):

- **EDR-como-changelog.** Anotaciones de edición inline en la prosa — "(actualizada `<fecha>`)", "(renombrado tras el refactor de la capa X)". La evolución la lleva git; un cambio de fondo es un EDR nuevo + el viejo `Superseded`, no un parche entre paréntesis.
- **Anclar a nombres de planificación efímeros.** Slices, tickets, fases de roadmap, nombres de milestone en el razonamiento. Envejecen en semanas y no significan nada para quien lee el EDR después. El Contexto expresa el condicionante (qué cambió, por qué hace falta decidir), no cómo se llamaba el trabajo que lo trajo.

## Estados (un EDR no es estático)

Un EDR tiene ciclo de vida. Dos estados importan para el concepto:

- **Borrador / inferido** — se registró el **QUÉ** (qué se eligió, observable) pero el **PORQUÉ está vacío**: es *evidencia de que se tomó una decisión*, sin la razón ratificada. Es un candidato, no una decisión cerrada.
- **Aceptado** — una persona lo ratificó y completó el porqué (contexto, alternativas, consecuencias). Recién ahí es un EDR pleno.

La promoción de borrador a aceptado siempre la hace una persona; nunca es automática. Un borrador todavía no "decide" — solo deja constancia de que hay algo para decidir.
