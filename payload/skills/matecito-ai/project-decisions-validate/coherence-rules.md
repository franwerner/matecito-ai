# Rúbrica de coherencia y completitud

Lista central de chequeos que aplica `project-decisions-validate`. Es **ratchet-able**: cuando aparece una contradicción nueva, se agrega acá y queda cubierta para siempre.

## Cómo la lee el validador

Cada chequeo tiene: **severidad** (CRITICAL / WARNING / SUGGESTION), una **condición** evaluada sobre los ADRs, el/los **dominio(s)** donde viven los ADRs involucrados (para localizar los archivos), y un **mensaje** (qué/por qué/sugerencia). El validador evalúa las condiciones contra `.claude/adr/<dominio>/` y reporta las que se cumplen.

## Mapa slug → dominio

Para localizar el archivo de cada ADR nombrado abajo. Un ADR vive en `.claude/adr/<dominio>/<slug>.md`. (El campo `Dominio:` del encabezado del ADR es la fuente de verdad si hay duda.)

| Slug | Dominio |
|---|---|
| context | context |
| architecture-style, layers-and-dependencies, inter-layer-communication, folder-structure | structure |
| error-handling, concurrency-async, background-jobs, resilience, caching | runtime |
| data-access, data-modeling | data |
| logging, metrics, tracing, health-checks | observability |
| auth, authorization, input-validation, rate-limiting, cors, secrets-management, dependency-scanning | security |
| api-contract, cli-contract, library-contract, event-contract | contracts |
| configuration, dependency-injection, testing-strategy, arch-enforcement, ci-quality-gates, deployment-topology, documentation, feature-flags | delivery |
| accessibility | frontend |
| nfr-performance, scalability, i18n-l10n | quality |

Dominios reservados (aparecen solo si el proyecto los pobló vía ratchet): `lifecycle`, `integration`, `privacy`, `release`, `domain-logic`, `compliance`, `ux-product`.

---

## Chequeos genéricos

### Completitud

- **[WARNING]** Una fase relevante para el tipo de proyecto no tiene ADR (ni siquiera `Not Applicable`). Hueco silencioso. *(Requiere la lista de fases relevantes o el catálogo `concerns/INDEX.md`; si no están disponibles, marcar como "no verificable".)*

### Higiene de status

- **[WARNING]** ADR `Accepted` sin sección "Decisión" con contenido concreto.
- **[WARNING]** ADR `Pending` o `Deferred` sin razón ni trigger.
- **[CRITICAL]** ADR `Superseded` sin link "Reemplazado por", o el ADR linkeado no existe. *(El link puede ser intra-dominio `<slug>.md` o cross-dominio `../<dominio>/<slug>.md`; verificá ambos.)*
- **[WARNING]** ADR `Not Applicable` sin razón.

### Verificabilidad

- **[WARNING]** `layers-and-dependencies` `Accepted` con reglas en prosa vaga en vez de globs/paths verificables.
- **[SUGGESTION]** Lenguaje vago ("tratá de no", "en lo posible", "idealmente", "evitar cuando se pueda") en las reglas de un ADR `Accepted`.

### Integridad de la taxonomía

- **[CRITICAL]** Existe una carpeta bajo `.claude/adr/` que no es un dominio canónico ni `tech/`. La taxonomía es cerrada; un dominio nuevo es decisión de catálogo, no de proyecto.
- **[WARNING]** El campo `Dominio:` del encabezado de un ADR no coincide con la carpeta en la que está el archivo. Mover el ADR a su carpeta correcta o corregir el campo.
- **[WARNING]** Un ADR está listado en el índice raíz (`.claude/adr/INDEX.md`) pero su dominio no tiene `INDEX.md`, o viceversa (índice de dominio sin entrada en el raíz). Índices desincronizados.
- **[SUGGESTION]** Un dominio tiene `INDEX.md` pero ningún ADR (carpeta de dominio vacía en la salida). Limpiar la carpeta o el índice.

### Coherencia del campo `tipo`

- **[SUGGESTION]** Un ADR marcado `tipo: convención` o `tipo: política` tiene una sección "Alternativas consideradas" sustanciosa → quizá es en realidad una `decisión`; revisar el tipo.
- **[SUGGESTION]** Un ADR marcado `tipo: decisión` y `Accepted` sin "Alternativas consideradas" ni "Consecuencias" → una decisión sin trade-offs documentados es sospechosa; o falta contenido o es en realidad una convención.
- **[WARNING]** Un ADR `tipo: política` `Accepted` sin "Reglas concretas" verificables → una política sin reglas accionables no se puede cumplir ni chequear.

---

## Contradicciones conocidas (combinaciones entre ADRs)

La columna **Dominio(s)** indica dónde viven los ADRs de la condición. "cross" = la contradicción cruza dominios.

| # | Severidad | Dominio(s) | Condición | Mensaje |
|---|---|---|---|---|
| 1 | CRITICAL | structure | `architecture-style` ∈ {Clean, Hexagonal} **y** `inter-layer-communication` permite que las entidades de dominio crucen los bordes | Clean/Hexagonal exige DTOs en los bordes; entidades crudas cruzando rompen el aislamiento del dominio. Definir mapeo a DTOs. |
| 2 | CRITICAL | structure | `architecture-style` ∈ {Clean, Hexagonal} **y** `inter-layer-communication` = "no usamos interfaces" | La inversión de dependencias —núcleo de Clean/Hexagonal— requiere interfaces en los bordes I/O. Sin interfaces no hay desacople real. |
| 3 | CRITICAL | security | `authorization` `Accepted` **y** `auth` = `Not Applicable` | No se puede autorizar sin autenticar. Si hay modelo de permisos, tiene que haber identidad. |
| 4 | WARNING | structure + delivery (cross) | `layers-and-dependencies` `Accepted` (hay capas) **y** `arch-enforcement` ∈ {`Not Applicable`, "sin enforcement"} | Reglas de capas sin enforcement automatizado se degradan en el primer sprint apurado. Considerar import-linter/dependency-cruiser/ArchUnit. |
| 5 | WARNING | runtime + tech (cross) | `error-handling` = "Result/Either" **y** stack Python/JS/TS sin librería de Result en `tech/` | Result types sin librería en ese stack es incómodo de sostener. Confirmar la elección o registrar la librería (returns, neverthrow). |
| 6 | WARNING | data + structure (cross) | `data-access` usa patrón Repository **y** `architecture-style` = "sin patrón formal" | Repository sin arquitectura en capas suele ser inconsistente u over-engineering. Revisar si hace falta. |
| 7 | WARNING | security + context (cross) | `auth` = `Not Applicable` **y** tipo de proyecto ∈ {api-rest, api-graphql, web-ssr, web-spa} | API/web marcando auth como N/A: confirmar que es realmente interno o sin usuarios. |
| 8 | WARNING | delivery + tech (cross) | `testing-strategy` con TDD obligatorio o cobertura mínima **y** sin test framework en `tech/` | Se exige TDD/cobertura pero no hay framework de test registrado. Registrar el framework. |
| 9 | WARNING | security + structure (cross) | `input-validation` "solo en el borde" **y** `inter-layer-communication` "validación solo en dominio" | La validación quedó con un hueco o duplicada. Definir cuál es la fuente de verdad (o defensa en profundidad explícita). |
| 10 | SUGGESTION | observability + delivery (cross) | `metrics` / `tracing` / `health-checks` `Accepted` **y** `deployment-topology` = `Not Applicable` | Observabilidad definida sin saber dónde/cómo corre el sistema. Conviene definir la topología. |
| 11 | SUGGESTION | runtime + quality/delivery (cross) | `caching` con cache distribuido **y** `scalability`/`deployment-topology` indica una sola instancia | Un cache distribuido aporta poco con una sola instancia; quizá alcanza in-memory. |

### Contradicciones que involucran dominios reservados

Aplican solo si el proyecto pobló esos dominios vía ratchet. Se anticipan acá para que el día que aparezcan, el validador ya las atrape.

| # | Severidad | Dominio(s) | Condición | Mensaje |
|---|---|---|---|---|
| 12 | CRITICAL | privacy + lifecycle (cross) | `privacy` define un período de retención de datos personales **y** `lifecycle` define una política de borrado/retención incompatible (ej: borrado inmediato vs retención de 5 años) | Retención y borrado se contradicen. La política de datos personales y el ciclo de vida de datos tienen que dar el mismo número. |
| 13 | CRITICAL | compliance + security (cross) | `compliance` exige cifrado en reposo / control de acceso específico (ej: PCI-DSS, HIPAA) **y** `security`/`secrets-management` no lo provee | El requisito regulatorio no está soportado por la decisión de seguridad. Sin esto, el proyecto no cumple. |
| 14 | WARNING | privacy + data (cross) | `privacy` exige minimización de datos **y** `data-modeling` define un modelo que captura datos personales no justificados | Minimización vs modelo que sobre-captura. Revisar qué PII se guarda realmente y por qué. |
| 15 | WARNING | compliance + lifecycle (cross) | `compliance` exige retención mínima por norma (ej: registros financieros 7 años) **y** `lifecycle` define una retención menor | El borrado viola la obligación de retención regulatoria. Alinear el período al mínimo legal. |
| 16 | WARNING | integration + runtime (cross) | `integration` define consumo de APIs de terceros **y** `resilience` = `Not Applicable` o sin política de timeout/retry para externos | Consumir terceros sin política de resiliencia: un proveedor lento/caído cuelga el sistema. Definir timeouts, retries y fallback. |
| 17 | WARNING | release + contracts (cross) | `release` define versionado semántico con breaking-change policy **y** `api-contract`/`library-contract` no define estrategia de versionado de la interfaz | La política de versionado del producto no está respaldada por una estrategia de versionado del contrato. Alinear ambas. |
| 18 | SUGGESTION | domain-logic + structure (cross) | `domain-logic` define agregados/invariantes ricos (DDD) **y** `architecture-style` = {Layered N-tier simple, sin patrón formal} | Lógica de dominio rica sobre una arquitectura sin aislamiento del dominio tiende a "anemic domain model". Considerar Clean/Hexagonal. |
| 19 | SUGGESTION | ux-product + frontend (cross) | `ux-product` define estados de error/carga de cara al usuario **y** `frontend`/`accessibility` no contempla cómo se comunican accesiblemente | Los estados de UX definidos deberían cumplir accesibilidad (foco, anuncios a screen reader). Cerrar el gap. |

---

## Ratchet

Cuando encontrás una contradicción o un chequeo útil que no está acá:

- Genérico → agregalo como bullet en la sección "Chequeos genéricos" (en la subsección que corresponda).
- Combinación entre ADRs → agregá una fila en "Contradicciones conocidas" (o en la subtabla de dominios reservados si involucra uno).

Siempre con severidad, dominio(s) y mensaje (qué/por qué/sugerencia). Si la combinación cruza dominios, marcala "cross" en la columna Dominio(s). Así el validador la atrapa de ahí en más.
