# Capability specs — `rule`

Reglas de negocio transversales, sin flujo (scoping, políticas, invariantes).

**Cuándo consultar este tipo:** antes de tocar cómo un evento se asocia a un (proyecto, change), cualquier lógica que dependa de la identidad por request (proyecto, change, sesión), la resiliencia de cualquier productor que ingesta contenido al broker, o cualquier resolución transversal del CLI: qué dominios están activos, de dónde sale un valor de configuración, y qué se registra o se auto-aprueba en el host destino.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `event-scoping` | Asocia cada evento (de una tool MCP o de un hook) a un (proyecto, change) | Accepted | [`event-scoping.md`](event-scoping.md) |
| `ingestion-spool` | Resiliencia transversal de la ingesta: fire-and-forget, spool local completo ante falla y reconciliación idempotente en orden | Accepted | [`ingestion-spool.md`](ingestion-spool.md) |
| `domain-activation-shim` | Qué dominios se consideran activos: conjunto vacío = todos los del payload; conjunto explícito filtrado contra los descubiertos | Inferred | [`domain-activation-shim.md`](domain-activation-shim.md) |
| `agent-model-resolution-precedence` | Precedencia proyecto → global → default al resolver el modelo de un agente; la ausencia no se sustituye | Inferred | [`agent-model-resolution-precedence.md`](agent-model-resolution-precedence.md) |
| `strict-tdd-resolution-precedence` | Precedencia proyecto → global → false del guard de TDD estricto, con clave ausente distinta de false | Inferred | [`strict-tdd-resolution-precedence.md`](strict-tdd-resolution-precedence.md) |
| `typed-flag-resolution-precedence` | Flags tipados por dominio en tri-estado (ausente ≠ false), aislados entre dominios, con la misma precedencia | Inferred | [`typed-flag-resolution-precedence.md`](typed-flag-resolution-precedence.md) |
| `legacy-config-migration` | Plegado de las claves heredadas al dominio por defecto sin pisar lo presente, y migración de una sola vez del archivo de modelos heredado | Inferred | [`legacy-config-migration.md`](legacy-config-migration.md) |
| `model-value-validation` | Se valida el valor del modelo contra el conjunto que declara el host; los nombres de agente no se restringen (se descubren del payload) | Inferred | [`model-value-validation.md`](model-value-validation.md) |
| `mcp-permission-auto-approval` | Los patrones de auto-aprobación se derivan de las integraciones de los dominios activos, por convención con override, más un patrón fijo | Inferred | [`mcp-permission-auto-approval.md`](mcp-permission-auto-approval.md) |
| `hook-reconciliation-by-matecito-id` | Reconciliación por marcador de identidad propio: reemplaza lo propio obsoleto, preserva lo del usuario, limpia grupos vacíos, idempotente | Inferred | [`hook-reconciliation-by-matecito-id.md`](hook-reconciliation-by-matecito-id.md) |
| `hook-registry-domain-filtering` | Un hook del dominio compartido entra siempre; uno de dominio concreto solo si su dominio está activo | Inferred | [`hook-registry-domain-filtering.md`](hook-registry-domain-filtering.md) |
| `repo-component-declaration` | Declaración del set de componentes (superficies) del repo en el config de proyecto; no se hereda del global | Accepted | [`repo-component-declaration.md`](repo-component-declaration.md) |
| `capability-spec-components-axis` | Línea `Components:` en el header de cada spec para registrar qué superficies implementan el comportamiento | Accepted | [`capability-spec-components-axis.md`](capability-spec-components-axis.md) |
| `component-projection-independence` | Independencia de proyecciones: cada consumidor del set de componentes proyecta con su propia granularidad, vida y ratificación; comparten set+gate, nunca el valor | Accepted | [`component-projection-independence.md`](component-projection-independence.md) |
| `component-inference-ratification` | Propuesta de componentes por bootstrap (del set y per-spec), por mine (longest-prefix), y por cambio (intake), con ratificación explícita requerida en cada vía | Accepted | [`component-inference-ratification.md`](component-inference-ratification.md) |
| `contract-shape-proposal` | Proponer una forma de contrato no especificada como un item compuesto único (una línea + fields tipados); cómo se ratifica y vuelve a la fase que la pidió | Accepted | [`contract-shape-proposal.md`](contract-shape-proposal.md) |
| `single-rooted-spec-store` | El store de specs vive siempre en el root; antes de partir por app se ofrece el eje `components` como alternativa | Accepted | [`single-rooted-spec-store.md`](single-rooted-spec-store.md) |
| `rubric-mechanical-semantic-split` | La mitad mecánica de las rúbricas vive como filas en un contrato de datos; la skill retiene solo el criterio semántico irreductible | Accepted | [`rubric-mechanical-semantic-split.md`](rubric-mechanical-semantic-split.md) |
| `mailbox-item-summary-rationale-split` | Items en secciones buzón llevan un split opt-in: `summary` se imprime en el gate, `rationale` se emite siempre en el mismo bloque; el split es declarado por la sección, sea cual sea su forma de renderización | Accepted | [`mailbox-item-summary-rationale-split.md`](mailbox-item-summary-rationale-split.md) |
| `decision-re-emergence-reconfirmation` | Una decisión ratificada en un gate que resurge en un gate posterior se reconoce en una sola pregunta de sí/no; diferentes contenidos o records nuevos se presentan en forma completa | Accepted | [`decision-re-emergence-reconfirmation.md`](decision-re-emergence-reconfirmation.md) |
| `deviation-consult-before-applying` | Consultar un fork sin mandato antes de aplicarlo: si una edición abre un punto no-fijado con más de una resolución válida, parar y preguntar; si no hay alternativa válida, absorber y nombrar la restricción | Accepted | [`deviation-consult-before-applying.md`](deviation-consult-before-applying.md) |
| `scenario-verification-scope-token` | Token `verification:` por escenario para asociar cada escenario a su dueño de verificación (in-scope, deferred → cambio, standing → dueño) | Accepted | [`scenario-verification-scope-token.md`](scenario-verification-scope-token.md) |
| `intake-brief-components-line` | Línea `Components` en el Intake Brief: clasificación por-cambio de superficies del repo tocadas, con gate presence-based y ratificación en la INTAKE GATE | Accepted | [`intake-brief-components-line.md`](intake-brief-components-line.md) |
