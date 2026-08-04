# Chequeos mecánicos de artefactos — `checks.yaml`

Contrato de datos que `~/.claude/scripts/validate-artifact.js` evalúa contra un store completo
(`--store <dir>`, junto con `--root <dir>` cuando el chequeo lo necesita). Lo consumen
`development-decisions-validate` (store `edr`) y `development-spec-validate` (store
`capability-spec`) — ninguna de las dos re-deriva estos chequeos en prosa; invocan el script e
ingieren su JSON. Su rúbrica (`coherence-rules.md`) retiene solo la mitad **semántica**: contradicción
de significado, lenguaje vago, identificadores volátiles, anotaciones de edición inline,
completitud-por-concern y `Not Applicable` sin razón — nada de eso es decidible por inspección
estructural, así que no está acá.

## Qué es un `kind`

Un **kind** es un chequeo mecánico que el motor de `validate-artifact.js` ya sabe evaluar —
presencia/vacío de sección, sincronía índice↔archivo, un link que resuelve o no, una carpeta dentro o
fuera de la taxonomía. Cada **fila** bajo `checks:` es una instancia concreta de un kind: a qué store
apunta, con qué severidad, con qué mensaje. Un kind es código; una fila es dato.

## Los 19 kinds

| Kind | Qué chequea |
|---|---|
| `section-non-empty` | Una sección que el status del artefacto exige tiene contenido, no solo el encabezado. |
| `section-required-when` | Una sección opcional que, bajo una condición (status + dominio/tipo), debería estar presente. |
| `section-set-symmetry` | Un par de secciones que deberían aparecer juntas o ninguna — una con contenido y la otra vacía es asimetría. |
| `section-minimum` | Un `Accepted` que solo llena las dos secciones mínimas, sin ninguna de las que anclan o registran trade-offs. |
| `bullet-prefix` | Un ítem de lista que debería abrir con una de un set fijo de marcas (`[auto]`/`[manual]`). |
| `legacy-status` | Un valor de header deprecado que ya no es parte del enum vigente. |
| `scenario-shape` | Un escenario que no tiene los tres componentes que lo hacen verificable. |
| `skeleton-sections` | Las secciones "esqueleto" de un tipo de capability-spec, cada una con su propio set. |
| `header-in-set` | Un valor de header multivaluado fuera del set canónico declarado en la config del proyecto. |
| `header-required-when-source` | Un header que debería estar presente porque el eje que lo gobierna está declarado. |
| `glob-drift` | Un glob de `## Alcance` que ya no matchea ningún archivo real del repo. |
| `folder-taxonomy` | Una carpeta dentro del store que no es ni un dominio/tipo canónico ni una carpeta legal fuera de la taxonomía (`extra_folders`). |
| `index-sync` | Un artefacto y su fila de índice — a nivel archivo↔carpeta o carpeta↔raíz — donde uno existe sin el otro. |
| `location-by-slug` | Un artefacto cuyo slug corresponde, según el catálogo, a un dominio distinto de la carpeta donde vive. |
| `dangling-link` | Un link markdown que no resuelve a nada — ni en el store propio ni en el store que declara como destino. |
| `orphan-folder` | Una carpeta de dominio/tipo con `INDEX.md` pero sin ningún artefacto real adentro. |
| `inbound-reference-count` | Una capability compartible referenciada por un solo consumidor — candidata a vivir adentro de ese consumidor en vez de aparte. |
| `deprecated-referenced` | Un artefacto `Deprecated` todavía linkeado como si estuviera vigente. |
| `store-uniqueness` | Más de una carpeta con el nombre del store en el repo — el store se partió por accidente. |

Conjunto **cerrado**: entra un kind si, y solo si, tiene al menos una fila real en el `checks.yaml`
inicial. `stale-status` quedó afuera a propósito — ambos candidatos que lo habrían alimentado (`Draft` /
`Inferred` "de larga data") necesitan un umbral de vigencia que nadie fijó, y cada uno tiene una mitad
semántica que sigue viviendo en la rúbrica.

## Severidad de `SECTION-MISSING` (chequeo de forma por archivo, no un `kind`)

Antes de llegar al motor de `checks.yaml`, `validate-artifact.js` valida cada artefacto contra el
contrato de forma (`edr.yaml` / `capability.yaml`) que también usa `render-artifact.js` para
producirlo — ese es el chequeo `SECTION-MISSING` (y su contraparte `SECTION-UNEXPECTED`), no una fila
de `checks:`. Su severidad **no es uniforme**: depende de qué significa `emitted:` para esa sección.

- **`emitted: always`** y ausente → **`error`**. Es incondicional — su ausencia es un defecto
  estructural real, sea cual sea el status.
- **`emitted: on-status`** cuyo `status` aplica, y ausente → **`warning`**. Ese declarativo se escribió
  para `render-artifact.js` ("emitila cuando el status sea X"), una regla de **producción** — no una
  obligación de conformidad sobre un archivo escrito a mano. La rúbrica original que este store migró
  nunca trató esa ausencia como error: sus chequeos más cercanos (`section-set-symmetry`,
  `section-minimum`) son estrictamente más débiles y migraron a `warning`/`nota`.
- **`emitted: when-present`** → nunca es hallazgo por ausencia (sin cambios).

`SECTION-UNEXPECTED` (una sección presente cuando su status no la admite) sigue siendo **`error`**
siempre — eso sí es un defecto real, no una omisión de producción. Los enums del header
(`STATUS-ILLEGAL`, `HEADER-MISSING`, `HEADER-UNEXPECTED`) tampoco cambian: siguen `error`.

**Si una sección concreta merece `error` aun siendo `on-status`**, no se hardcodea una excepción en el
script — se declara como fila explícita en `checks:` con la severidad que corresponda (mismo patrón que
cualquier otro `kind`). Este contrato es el único lugar donde una excepción de severidad se justifica.

## `extra_folders`

Carpetas **legales** dentro de un store pero **fuera** de su taxonomía cerrada — hoy el único caso real
es `tech/` en el store de EDRs (mini-EDRs por tecnología concreta, sin dominio). `folder-taxonomy` nunca
las marca como drift, y ningún otro chequeo las trata como si fueran uno de los dominios/tipos
canónicos. El nombre por sí solo no lo dice — de ahí esta aclaración.

## El ratchet

Agregar una **fila** de un kind que ya existe no toca código: es una línea de dato en `checks:`, y ese
es precisamente el punto de migrar la mitad mecánica — el costo de sumar un chequeo nuevo de un kind
conocido es el mismo que agregarlo a la rúbrica en prosa, sin perder la evaluación automática. El
conjunto de **kinds** crece solo cuando aparece un chequeo mecánico que ningún kind existente cubre —
eso sí es trabajo de motor, en `validate-artifact.js`, no una fila.
