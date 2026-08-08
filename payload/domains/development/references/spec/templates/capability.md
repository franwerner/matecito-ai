<!-- Canonical template: capability-spec (`.matecito-ai/development-specs/<type>/<capability>.md`, type ∈ flow | rule | lifecycle | process). Consumido por la fase de Materialización de development-spec-bootstrap y por el merge de sdd-archive. --><!-- matecito-ai: this path omitted the `<type>/` segment the whole rest of the ecosystem requires (the domain fragment, sdd-archive, sdd-verify, spec-mine) — and the file's own line-8 comment contradicted it. -->

# Capability — <nombre en lenguaje de dominio>

- **Status:** <Inferred | Draft | Accepted | Deprecated>
- **Date:** <YYYY-MM-DD>
- **Components:** <componente>, <componente> <!-- solo si el proyecto declaró `repo.components`; omitir la línea entera si el eje no está activo -->

<!-- Header mínimo a propósito: el TIPO y el slug ya están en la ruta del archivo (`<type>/<capability>.md`) — no se repiten como propiedad. El historial lo lleva git. `Date` es la fecha del comportamiento vigente, no la de la última edición. Si Status es Deprecated, agregar una línea `**Reemplazado por:** [<capability>.md](<capability>.md)`. `Components` es multivaluado y cada valor debe estar en el set declarado en `repo.components` (ver `~/.claude/references/repo-components/README.md` → "Declaración" y "Gate presence-based"); gate presence-based: sin ese eje declarado, no se escribe esta línea. -->

<!-- VOCABULARIO (aplica a TODAS las secciones): se escribe en el idioma del dominio + contrato público (endpoints públicos, códigos de error expuestos). NUNCA identificadores internos volátiles: clases, métodos, columnas de base de datos, errores internos, rutas de archivo. El "cómo" es del código; el "por qué" es del EDR. Ver `~/.claude/references/spec/README.md`. -->

## Propósito

<1-2 líneas: qué logra esta capacidad y para quién. El resultado de negocio, no el mecanismo.>

## Actores

<quién participa: el consumidor/usuario que la dispara, sistemas externos, procesos de fondo. Uno por línea con una frase de rol.>

## Precondiciones

<qué debe ser cierto antes de que el flujo pueda ocurrir. Omitir si no hay ninguna relevante.>

## Flujo principal

<el happy path, en pasos numerados. Cada paso es una acción observable, en lenguaje de dominio.>

1. <paso>
2. <paso>

## Ramas / flujos alternativos

<cada bifurcación del happy path: la condición que la dispara → el resultado. Omitir la sección si no hay ramas.>

- **<condición>** → <qué hace el sistema>.

## Casos borde

<situaciones límite y qué hace el sistema en cada una. Es lo que más se pierde si no se escribe. Omitir si no aplica.>

- **<situación límite>** → <comportamiento>.

## Reglas de negocio

<invariantes explícitas que gobiernan la capacidad, con valores concretos. No "se valida la ventana" sino "la ventana es de 24h; la abre un mensaje entrante del consumidor; al vencer no se puede enviar". Omitir si no hay.>

- <regla con valores concretos>.

## Entidades y estados

<las entidades de dominio que la capacidad toca, por su semántica (no su tabla), y las transiciones de estado relevantes. Omitir si no aplica.>

- **<entidad>** — <qué representa>. Estados: <estado> → <estado> (<qué dispara la transición>).

## Errores de cara al actor

<qué recibe el actor en cada falla — el contrato de error observable, no la excepción interna. Omitir si no aplica.>

- **<situación de falla>** → <qué se responde al actor>.

## Escenarios

<lo que vuelve al spec verificable. Cada regla/flujo/borde importante tiene al menos un escenario Given/When/Then. Son la aserción testeable del comportamiento.>

### Scenario: <nombre>

<!-- La línea `verification:` es OPCIONAL y va primero, antes del GIVEN. Su ausencia significa
     `in-scope`, así que un escenario sin ella no es un escenario sin marcar. Los valores y su
     semántica los define `~/.claude/references/spec/README.md` → "El token `verification:` por
     escenario"; este template la muestra, no la define. -->
- **verification:** `deferred → <cambio>`
- **GIVEN** <estado inicial>
- **WHEN** <acción>
- **THEN** <resultado observable esperado>

## Referencias

<!-- Opcional. Links a otros capability-specs de este mismo store que se relacionan con este comportamiento — nunca a un EDR, un PRD/proposal ni ningún artefacto fuera de `.matecito-ai/development-specs/`: el store es cerrado, y esa relación se explica en prosa, sin materializarla como link. Omitir si no hay. -->

- **<capability-spec relacionado>** → [`../<type>/<capability>.md`](../<type>/<capability>.md) — <qué relación tiene con este comportamiento>.

<!--
Notas del contrato (no van en el capability-spec generado):
- No hay sección `Historial`. El historial de ediciones lo lleva git.
- Nombres de sección y prosa en español; header (`Status`, `Date`) en inglés.
- Self-check antes de dar por escrito el spec: releé todas las secciones y por cada nombre de clase/método/columna/archivo/error interno, reformulá la frase en lenguaje de dominio o contrato público. El store es cerrado: ninguna sección — ni siquiera «Referencias» — linkea ni nombra por título/slug un EDR, un PRD/proposal ni ningún artefacto fuera de `.matecito-ai/development-specs/`; si hace falta anclar a una decisión técnica, la relación queda en prosa, conceptual, sin identificar el artefacto.
- Toda sección marcada "Omitir si no aplica" se borra entera cuando no aplica — no se deja vacía ni con "N/A".
-->
