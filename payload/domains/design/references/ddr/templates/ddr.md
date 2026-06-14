<!-- Canonical template: DDR individual (`.matecito-ai/ddr/<surface>/<slug>.md`). Consumido por la fase de Materialización de design-decisions-bootstrap/SKILL.md y por design-decisions-mine (status Inferred). -->

# DDR — <título>

- **Status:** <Inferred | Accepted | Pending | Deferred | Superseded>
- **Type:** <decision | convention | policy>
- **Date:** <YYYY-MM-DD>
- **Applied principle:** <Opcional. Solo si la decisión mapea a un principle del catálogo canónico en `~/.claude/references/design-principles/`. Formato: `<Nombre> — <1 línea de por qué este principle>`. Ej: `Contrast — la jerarquía de la paleta se apoya en diferencia de luminancia, no de tono`. Si no aplica, omitir la línea completa.>

<!-- Header mínimo a propósito: la surface y el slug ya están en la ruta del archivo; el autor y las fechas de edición los lleva git. `Date` es la fecha de la decisión vigente, no la de la última edición. -->

## Contexto

<por qué hace falta esta decisión, qué condicionantes hay de la pieza/marca/sistema/alcance. En DDRs `Inferred` queda vacío — mine NUNCA infiere el porqué.>

## Decisión

<lo decidido, en imperativo, con valores concretos chequeables contra Figma. Ej: "La paleta usa `Primary/500` = `#2563EB`, neutros `Neutral/50–900` en escala de 10 pasos, semánticos success/warning/error/info; todos materializados como color styles nombrados.". En DDRs `Inferred` queda vacío.>

<!-- Si Status es Pending o Deferred, REEMPLAZAR "Decisión" por:

## Razón de omisión / aplazamiento

**Status:** <Pending | Deferred>

<1-2 líneas con el motivo, honesto y concreto.
- Pending: indicá el trigger esperado ("cuando llegue X").
- Deferred: fecha o condición de revisión.>
(Los `Not Applicable` no usan este template — viven como fila en el INDEX de la surface.)
-->

<!-- Si Status es Superseded, agregar:
## Reemplazado por
[<slug-del-nuevo>.md](<slug-del-nuevo>.md) — <1 línea de por qué cambió la decisión>
(Si el DDR nuevo está en otra surface, usar ruta relativa: [../<otra-surface>/<slug>.md](../<otra-surface>/<slug>.md))
-->

## Consecuencias

<si Accepted, positivas y trade-offs. Si no se decidió, omitir.>

## Alternativas consideradas

<si Accepted, listá alternativas evaluadas con por qué no se eligieron. Si no se decidió, omitir.>

## Reglas verificables

<!-- Solo si Accepted. Cada regla es una aserción chequeable con VALORES CONCRETOS contra Figma (hex, ratio, escala, px, nombres de tokens) — no un adjetivo vago. Marcá el mecanismo de verificación al inicio de cada una:
- [tool: figma] → chequeable leyendo el archivo Figma vía el MCP figma (styles/components/variables nombrados, hex, px, escala).
- [tool: contrast] → chequeable con un cálculo de ratio de contraste (WCAG).
- [manual] → hoy solo se verifica en revisión visual humana, no hay check automático. -->

- **[tool: figma]** <regla concreta con valores>. Ej: **[tool: figma]** el color style `Primary/500` tiene hex exactamente `#2563EB`.
- **[tool: contrast]** <regla de ratio>. Ej: **[tool: contrast]** texto sobre `Primary/500` cumple ratio ≥ 4.5:1.
- **[manual]** <regla que hoy solo se chequea en review visual>.

## Alcance

<!-- Opcional. Incluir SOLO cuando la decisión gobierna un locator verificable en Figma (típicamente surfaces `foundation` / `components` / `layout`). Omitir la sección entera si la decisión no ancla en estilos/componentes/frames concretos (ej. una política de tono de voz). -->

Locator **a nivel sistema** —styles, sets de componentes o frames estables, no nodos efímeros— que esta decisión gobierna. El validador (y `design-decisions-mine` en re-run) chequea que el locator siga existiendo en el archivo Figma; si desaparece, es drift a resolver (el sistema cambió o la decisión quedó obsoleta).

- `<locator>` — <qué representa>
- Ej: `Primary/*`, `Neutral/*`, `Error/*` — color styles que componen la paleta.

## Evidencia (inferida)

<!-- TRANSITORIA: solo en DDRs `Inferred` (drafteados por design-decisions-mine desde el archivo Figma, NO decididos aún por un humano). Al ratificarse a `Accepted` vía design-decisions-bootstrap modo update, esta sección se ELIMINA (git conserva la traza). Mine NUNCA infiere el porqué — solo registra el qué observado. -->

- **kind:** <token | component | pattern | absence>
- **observado:** <lo que se vio en el archivo Figma — el QUÉ, sin el porqué>
- **prevalencia:** <ej: usado en 12/14 frames. Omitir si el kind no aplica (absence).>

<!-- El locator estructural NO va en esta sección: para kind `token`/`component`/`pattern` se llena `## Alcance` (la lista de tokens / frames / set de componentes), que el validador y mine usan como ancla de drift. Para `absence` no hay locator. Así la sección Evidencia queda acotada a la metadata de la inferencia. -->

## Relacionados

<!-- Opcional. Links tipados a otros DDRs. Tipos: `depende-de`, `refina`, `relacionado-con`. Para reemplazos usar la sección "Reemplazado por" de arriba, no esta. Misma surface: ruta corta `<slug>.md`; otra surface: ruta relativa `../<surface>/<slug>.md`. Omitir la sección si no hay vínculos. -->

- `depende-de` → [<slug>.md](<slug>.md) — <1 línea>
- `relacionado-con` → [../<surface>/<slug>.md](../<surface>/<slug>.md) — <1 línea>

<!--
Notas del contrato (no van en el DDR generado):
- No hay sección `Historial`. El historial de ediciones lo lleva git; la evolución de decisiones se ve en la cadena de `Superseded`.
- Header en inglés (`Status`, `Type`, `Date`, `Applied principle`); nombres de sección y prosa en español.
- `Inferred` llena solo header + `## Evidencia (inferida)` + `## Alcance` (cuando el kind tiene locator). `## Contexto`, `## Decisión`, `## Consecuencias`, `## Alternativas consideradas`, `## Reglas verificables` quedan vacías hasta la ratificación.
-->
