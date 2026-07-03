<!-- Canonical template: ADR individual (`.matecito-ai/adr/<dominio>/<slug>.md`). Consumido por la fase de Materialización de SKILL.md. -->

# ADR — <título>

- **Status:** <Inferred | Accepted | Pending | Deferred | Superseded>
- **Type:** <decision | convention | policy>
- **Date:** <YYYY-MM-DD>
- **Applied pattern:** <Opcional. Solo si la decisión mapea a un patrón del catálogo canónico en `~/.claude/references/design-patterns/`. Formato: `<Nombre> — <1 línea de por qué este patrón>`. Ej: `Repository — necesitamos swap SQLite↔Postgres en tests sin tocar dominio`. Si no aplica, omitir la línea completa.>

<!-- Header mínimo a propósito: el dominio y el slug ya están en la ruta del archivo; el autor y las fechas de edición los lleva git. `Date` es la fecha de la decisión vigente, no la de la última edición. -->

## Contexto

<por qué hace falta esta decisión, qué condicionantes hay del proyecto/stack/equipo/alcance>

## Decisión

<lo decidido, en imperativo. Ej: "Usamos JWT con refresh tokens y rotación; access token de 15min, refresh de 7d.">

<!-- Vocabulario (aplica a Contexto/Decisión/Consecuencias/Alternativas): conceptos, patrones y límites — NUNCA identificadores internos volátiles (clase, método, columna, error interno, ruta de archivo). Si te sale escribir uno, reubicalo: límite estable → glob en `## Alcance`; aserción chequeable → `## Reglas verificables` (ahí sí podés nombrar la clase, es el ancla). Excepción: nombre de tecnología/librería y contrato público (endpoint público, código de error expuesto). El ejemplo de arriba es conceptual a propósito. Ver `~/.claude/references/adr/README.md` → "No es el cómo". -->

<!-- Si Status es Pending o Deferred, REEMPLAZAR "Decisión" por:

## Razón de omisión / aplazamiento

**Status:** <Pending | Deferred>

<1-2 líneas con el motivo, honesto y concreto.
- Pending: indicá el trigger esperado ("cuando llegue X").
- Deferred: fecha o condición de revisión.>
(Los `Not Applicable` no usan este template — viven como fila en el INDEX del dominio.)
-->

<!-- Si Status es Superseded, agregar:
## Reemplazado por
[<slug-del-nuevo>.md](<slug-del-nuevo>.md) — <1 línea de por qué cambió la decisión>
(Si el ADR nuevo está en otro dominio, usar ruta relativa: [../<otro-dominio>/<slug>.md](../<otro-dominio>/<slug>.md))
-->

<!-- Si Status es Inferred (ADR minado por development-decisions-mine desde el código, NO decidido aún por un humano), agregar esta sección. Es TRANSITORIA: al promoverse a Accepted vía bootstrap, se elimina (git conserva la traza). El humano completa Contexto/Decisión/Consecuencias; mine NUNCA infiere el porqué — solo registra el qué observado.

## Evidencia (inferida)

- **kind:** <estructural | config | patrón | ausencia>
- **observado:** <lo que se vio en el código — el QUÉ, sin el porqué>
- **prevalencia:** <ej: 40/42 handlers. Omitir si el kind no aplica (config/ausencia).>

El locator estructural NO va en esta sección: para kind `estructural`/`patrón` se llena `## Alcance` (los globs), que el validador ya usa como ancla de drift. Para `config` la evidencia es la entrada del manifest. Para `ausencia` no hay glob. Así la sección Evidencia queda acotada a la metadata de la inferencia.
-->

## Alcance

<!-- Opcional. Incluir SOLO en decisiones espaciales/estructurales (típicamente los dominios structure / folder-structure / layers-and-dependencies). Omitir la sección entera si la decisión no gobierna ubicación de código (ej. una política de rate limiting). -->

Globs **a nivel convención** —patrones estructurales estables, no archivos concretos— que esta decisión gobierna. El validador chequea que sigan matcheando algo; si dejan de matchear, es drift a resolver (el código se movió o la decisión quedó obsoleta).

- `<glob>` — <qué representa este patrón>
- Ej: `src/**/*.routes.ts` — handlers HTTP finos, uno por slice.

## Reglas verificables

<!-- Solo si Accepted. Cada regla es una aserción chequeable con valores concretos, no un adjetivo vago. Marcá el mecanismo de verificación al inicio de cada una:
- [tool: <herramienta/comando>] → lo enforced una herramienta (linter, dependency-cruiser, type-check, test).
- [manual] → hoy solo se verifica en revisión humana, no hay check automático. -->

- **[tool: <herramienta>]** <regla concreta con valores>. Ej: **[tool: dependency-cruiser]** ningún import desde `domain/**` hacia `infra/**`.
- **[manual]** <regla que hoy solo se chequea en review>.

## Alternativas consideradas

<si Accepted, listá alternativas evaluadas con por qué no se eligieron. Si no se decidió, omitir.>

## Consecuencias

<si Accepted, positivas y trade-offs. Si no se decidió, omitir.>

## Relacionados

<!-- Opcional. Links tipados a otros ADRs. Tipos: `depende-de`, `refina`, `relacionado-con`. Para reemplazos usar la sección "Reemplazado por" de arriba, no esta. Mismo dominio: ruta corta `<slug>.md`; otro dominio: ruta relativa `../<dominio>/<slug>.md`. Omitir la sección si no hay vínculos. -->

- `depende-de` → [<slug>.md](<slug>.md) — <1 línea>
- `relacionado-con` → [../<dominio>/<slug>.md](../<dominio>/<slug>.md) — <1 línea>

<!--
Notas del contrato (no van en el ADR generado):
- No hay sección `Historial`. El historial de ediciones lo lleva git; la evolución de decisiones se ve en la cadena de `Superseded`.
- Header en inglés (`Status`, `Type`, `Date`, `Applied pattern`); nombres de sección y prosa en español.
- Self-check antes de dar por escrito el ADR: releé Contexto/Decisión/Consecuencias/Alternativas y por cada nombre de clase/método/columna/archivo/error interno, convertilo en un glob (`## Alcance`) o una regla (`## Reglas verificables`), o reformulá la frase en términos de concepto. Excepción: tecnología/librería y contrato público.
-->
