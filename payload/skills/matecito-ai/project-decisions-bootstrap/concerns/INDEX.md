# Catálogo de concerns — Índice raíz

Menú de fases (concerns) que la skill puede recorrer, **organizado por dominio**. El motor (`SKILL.md`) lee este índice para (a) entender el mapa de dominios y (b) armar la lista de fases relevantes según el tipo de proyecto. Recién después lee el archivo individual de cada fase que va a tratar, para no cargar al contexto lo que no aplica.

## Cómo lo usa el motor

1. Fase 0 detecta el tipo de proyecto → lo mapea a un token de abajo.
2. El motor recorre la matriz de aplicabilidad y arma dos grupos:
   - **Requerido** → se incluye por default (el usuario puede marcarlo `Not Applicable` con razón).
   - **Recomendado** → se sugiere; el usuario decide.
3. Muestra el set; el usuario elige qué definir ahora, y puede agregar una **fase custom**.
4. Por cada fase elegida, el motor lee `<dominio>/<slug>.md` y la trata con las reglas del motor.
5. Las relevantes NO elegidas → ADR `Not Applicable` / `Pending` + razón. Nunca hueco silencioso.

## Dominios canónicos (fijos)

La taxonomía de dominios es **cerrada y la impone el motor** — la misma para el catálogo interno y para la salida `.matecito-ai/adr/`, así todos los repos del equipo se ven igual. Cada dominio tiene su propio `INDEX.md` con el detalle de sus concerns y el **criterio de pertenencia** (cuándo un concern nuevo va ahí).

### Activos (con concerns)

| Dominio | Qué agrupa | Índice |
|---|---|---|
| `context` | Propósito, alcance, stack, equipo | [context/INDEX.md](context/INDEX.md) |
| `structure` | Patrón arquitectónico, capas, dependencias, organización de archivos | [structure/INDEX.md](structure/INDEX.md) |
| `runtime` | Errores, concurrencia, background jobs, resiliencia, caching | [runtime/INDEX.md](runtime/INDEX.md) |
| `data` | Modelado y acceso a datos persistentes | [data/INDEX.md](data/INDEX.md) |
| `observability` | Logs, métricas, trazas, health checks | [observability/INDEX.md](observability/INDEX.md) |
| `security` | Auth, autorización, validación, secretos, CORS, rate limiting | [security/INDEX.md](security/INDEX.md) |
| `contracts` | API / CLI / librería / eventos que el sistema expone | [contracts/INDEX.md](contracts/INDEX.md) |
| `delivery` | CI/CD, testing, DI, config, deployment, flags, docs | [delivery/INDEX.md](delivery/INDEX.md) |
| `frontend` | Implementación técnica de la UI (accesibilidad, etc.) | [frontend/INDEX.md](frontend/INDEX.md) |
| `quality` | NFRs medibles: performance, escalabilidad, i18n | [quality/INDEX.md](quality/INDEX.md) |

### Reservados (taxonomía completa, sin concerns todavía)

Casilleros válidos para que nada quede sin lugar. Se pueblan vía ratchet cuando un proyecto lo necesite.

| Dominio | Qué agruparía | Índice |
|---|---|---|
| `lifecycle` | Migraciones, backups, retención y borrado de datos | [lifecycle/INDEX.md](lifecycle/INDEX.md) |
| `integration` | Consumo de terceros, webhooks, mensajería, anti-corruption | [integration/INDEX.md](integration/INDEX.md) |
| `privacy` | Datos personales, consentimiento, minimización | [privacy/INDEX.md](privacy/INDEX.md) |
| `release` | Versionado de producto, changelogs, breaking changes | [release/INDEX.md](release/INDEX.md) |
| `domain-logic` | Reglas de negocio, invariantes, DDD | [domain-logic/INDEX.md](domain-logic/INDEX.md) |
| `compliance` | GDPR, HIPAA, PCI-DSS, auditoría | [compliance/INDEX.md](compliance/INDEX.md) |
| `ux-product` | Flujos de usuario, estados vacíos, copy de errores | [ux-product/INDEX.md](ux-product/INDEX.md) |

## Tokens de tipo de proyecto

`api-rest` · `api-graphql` · `cli` · `libreria` · `web-spa` · `web-ssr` · `microservicio` · `monolito-modular` · `script`

`todos` = cualquier tipo.

## Fase custom

Si el usuario tiene un tema fuera de este catálogo, el motor le hace las preguntas genéricas, determina a qué **dominio canónico** pertenece, crea `<dominio>/<slug>.md` con el formato estándar y suma la fila al índice de ese dominio + a la matriz de abajo. Antes de guardarlo pregunta: **¿reusable (queda en el catálogo) o solo para este proyecto (solo genera el ADR)?**

## Leyenda

- **Prof.** = profundidad: `deep` (cuestionario propio) · `light` (1-2 preguntas).
- **Tipo** = `decisión` (alternativas y trade-offs) · `convención` (acuerdo de estilo) · `política` (regla verificable).

---

## Matriz de aplicabilidad

Cada fila apunta a `<dominio>/<slug>.md`. La columna **Dominio** es la carpeta canónica (interna y de salida).

### context
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [context](context/context.md) | context | deep | todos (entrada, siempre primero) | — |

### structure
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [architecture-style](structure/architecture-style.md) | structure | deep | todos salvo `script` | `script` |
| [layers-and-dependencies](structure/layers-and-dependencies.md) | structure | deep | los que eligieron un patrón en architecture-style | — |
| [inter-layer-communication](structure/inter-layer-communication.md) | structure | deep | proyectos con capas (Clean / Layered / Hexagonal) | `monolito-modular` |
| [folder-structure](structure/folder-structure.md) | structure | light | todos | — |

### runtime
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [error-handling](runtime/error-handling.md) | runtime | deep | todos | — |
| [concurrency-async](runtime/concurrency-async.md) | runtime | light | `microservicio`, `api-rest`, `api-graphql` | `web-ssr`, `cli` |
| [background-jobs](runtime/background-jobs.md) | runtime | light | — | `api-rest`, `api-graphql`, `microservicio`, `web-ssr` |
| [resilience](runtime/resilience.md) | runtime | light | `microservicio` | `api-rest`, `api-graphql` |
| [caching](runtime/caching.md) | runtime | light | — | `api-rest`, `web-ssr`, `microservicio` |

### data
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [data-access](data/data-access.md) | data | deep | proyectos con persistencia | — |
| [data-modeling](data/data-modeling.md) | data | light | proyectos con persistencia | — |

### observability
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [logging](observability/logging.md) | observability | light | todos | — |
| [metrics](observability/metrics.md) | observability | light | `microservicio` | `api-rest`, `api-graphql`, `web-ssr` |
| [tracing](observability/tracing.md) | observability | light | `microservicio` | `api-rest` distribuidas |
| [health-checks](observability/health-checks.md) | observability | light | `microservicio` | `api-rest`, `web-ssr` |

### security
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [auth](security/auth.md) | security | deep | proyectos con usuarios / expuestos a red | — |
| [authorization](security/authorization.md) | security | light | donde aplica auth | — |
| [input-validation](security/input-validation.md) | security | light | `api-rest`, `api-graphql`, `web-ssr`, `web-spa` | — |
| [rate-limiting](security/rate-limiting.md) | security | light | APIs públicas | `microservicio` |
| [cors](security/cors.md) | security | light | APIs consumidas por browser | — |
| [secrets-management](security/secrets-management.md) | security | light | proyectos con secretos | — |
| [dependency-scanning](security/dependency-scanning.md) | security | light | — | todos los productivos |

### contracts
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [api-contract](contracts/api-contract.md) | contracts | light | `api-rest`, `api-graphql`, `microservicio` | — |
| [cli-contract](contracts/cli-contract.md) | contracts | light | `cli` | — |
| [library-contract](contracts/library-contract.md) | contracts | light | `libreria` | — |
| [event-contract](contracts/event-contract.md) | contracts | light | sistemas event-driven / `microservicio` que publica eventos | `api-rest` con webhooks |

### delivery
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [configuration](delivery/configuration.md) | delivery | light | todos | — |
| [dependency-injection](delivery/dependency-injection.md) | delivery | deep | donde architecture-style usa DI | — |
| [testing-strategy](delivery/testing-strategy.md) | delivery | deep | todos | — |
| [arch-enforcement](delivery/arch-enforcement.md) | delivery | light | donde hay capas / reglas de dependencia definidas | — |
| [ci-quality-gates](delivery/ci-quality-gates.md) | delivery | light | — | todos los productivos |
| [deployment-topology](delivery/deployment-topology.md) | delivery | light | `microservicio`, `api-rest`, `web-ssr`, `web-spa` | — |
| [documentation](delivery/documentation.md) | delivery | light | — | todos los productivos |
| [feature-flags](delivery/feature-flags.md) | delivery | light | — | `microservicio`, `api-rest`, `web-ssr`, `web-spa` |

### frontend
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [accessibility](frontend/accessibility.md) | frontend | light | `web-spa`, `web-ssr` | — |

### quality
| Fase | Dominio | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [nfr-performance](quality/nfr-performance.md) | quality | light | — | servicios / APIs |
| [scalability](quality/scalability.md) | quality | light | `microservicio` | APIs de alta escala |
| [i18n-l10n](quality/i18n-l10n.md) | quality | light | — | apps user-facing |

---

## Mantenimiento (ratchet)

- **Agregar una fase:** decidí a qué **dominio canónico** pertenece (mirá el criterio de pertenencia en `<dominio>/INDEX.md`). Creá `<dominio>/<slug>.md` con el formato estándar (ver `runtime/error-handling.md` para `deep`, `runtime/caching.md` para `light`), sumá la fila al `<dominio>/INDEX.md` y a la matriz de arriba.
- **Poblar un dominio reservado:** el primer concern convierte el dominio de reservado a activo. Movelo de la tabla "Reservados" a "Activos" en este índice.
- **No crear dominios nuevos por proyecto:** la taxonomía es fija. Si de verdad falta un dominio, es una decisión de catálogo (agregarlo acá y en el motor), no algo improvisado en un repo.
- **Origen del catálogo:** sembrado de ISO/IEC 25010, 12-factor, arc42, OWASP ASVS y checklists de production-readiness. Nace casi completo y solo crece.
