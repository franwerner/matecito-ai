# EDR — El límite de sólo-discutir de la sesión lateral vive en prosa, sostenido por quien la mira

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
La sesión lateral tiene que quedarse en discutir — leer, razonar, devolver una conclusión — sin tocar el repo ni commitear. Ese límite se puede imponer con flags de lanzamiento (`--permission-mode`, `--disallowedTools`) o dejarlo sólo en prosa. En este ecosistema, una regla sólo-prosa normalmente no alcanza porque nadie está mirando la ejecución en el momento.

## Decisión
La sesión lateral se lanza en una terminal nueva sin modo de permisos y sin restricción de herramientas. El límite "sólo lee, razona y devuelve la conclusión; no toca archivos ni commitea" viaja únicamente como prosa, en la sección `## Return` del traspaso. Esto es defendible acá y no en cualquier otro lado del ecosistema porque la sesión lateral es interactiva y el usuario está sentado en ella — si empieza a escribir, lo ve.

## Reglas verificables
- **[manual]** El comando de lanzamiento de la sesión lateral no incluye `--permission-mode` ni ninguna lista de herramientas restringidas (`--allowedTools`/`--disallowedTools`).
- **[manual]** El límite "sólo discute" vive únicamente como prosa en la sección `## Return` del traspaso — no hay ningún mecanismo que lo haga cumplir mecánicamente.

## Alternativas consideradas
`--permission-mode plan` — descartado: agrega una rama de fallback para un CLI instalado que rechace el flag, y una pregunta abierta sobre si ese modo bloquearía la propia escritura a Engram que el mecanismo necesita. Una lista `--disallowedTools` enumerada — descartado por el mismo tipo de fragilidad, sin necesidad si el límite ya es defendible por prosa acá.

## Consecuencias
No hay rama de fallback para un flag no soportado, ni conflicto entre un modo de permisos y la escritura de la conclusión. El costo, aceptado explícitamente: si el usuario deja la sesión lateral sin mirarla, nada la frena de escribir o commitear.

## Relacionados
- `relacionado-con` → [side-discussion-pickup-on-consult.md](side-discussion-pickup-on-consult.md) — la otra decisión de este cambio ajustada en el mismo gate, sobre cómo el principal recoge la conclusión.
