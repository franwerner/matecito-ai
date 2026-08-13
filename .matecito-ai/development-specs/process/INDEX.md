# Capability specs — `process`

Comportamiento reactivo/de fondo, disparado por un evento o por el sistema (no por un actor): file-watchers, jobs, reconciliación, indexado.

**Cuándo consultar este tipo:** antes de tocar cualquier proceso de fondo del daemon —el file-watch e indexado de records, su versionado o su proyección—, el motor que despliega el payload en el host destino, o un validador que el host dispara ante un evento, leé el spec para conocer su disparador, su flujo, sus casos borde y sus reglas.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `index-decision-records` | Indexa y versiona (copy-on-write lazy) los `.md` de EDR/spec del proyecto ante file-watch, sosteniendo el pin de la versión que aplicó cada evento | Accepted | [`index-decision-records.md`](index-decision-records.md) |
| `ingest-mechanical-events` | Ingesta automática de arranque/fin de sub-agentes (hooks del runtime) al event-log del change activo | Accepted | [`ingest-mechanical-events.md`](ingest-mechanical-events.md) |
| `capture-code-snapshots` | Fotos del contenido de cada archivo tocado por el apply, en los bordes de batch (antes sincrónico pre-edición, después al submit), por change | Accepted | [`capture-code-snapshots.md`](capture-code-snapshots.md) |
| `deploy-payload-to-host` | Mapea los componentes de los dominios activos y los compartidos al host destino, compone el archivo de instrucciones raíz, detecta clashes y decide backup por estado de archivo | Accepted | [`deploy-payload-to-host.md`](deploy-payload-to-host.md) |
| `dev-install-local-build` | Ciclo de desarrollo en un solo comando: compila el binario del árbol y lo deja operativo en el entorno local sin que la instalación lo pise con la release | Accepted | [`dev-install-local-build.md`](dev-install-local-build.md) |
| `render-durable-artifact` | Construye el cuerpo completo de un artefacto durable (EDR, capability-spec) a partir de datos y del contrato declarado del tipo, garantizando conformidad por construcción | Accepted | [`render-durable-artifact.md`](render-durable-artifact.md) |
| `validate-artifact-structure` | Valida que un artefacto escrito se conforma a su contrato (secciones presentes, orden, cabeceras válidas, sincronía entre archivos); valida los contratos y templates contra sí mismos (`--self-check`); examina el alcance nombrado (un archivo, el store completo, o los artefactos embarcados), reporta hallazgos, no modifica nada | Accepted | [`validate-artifact-structure.md`](validate-artifact-structure.md) |
| `validate-git-commit-message` | Bloquea el commit ante atribución de IA, avisa (sin bloquear) si el mensaje no sigue Conventional Commits, y falla abierto ante cualquier ambigüedad | Inferred | [`validate-git-commit-message.md`](validate-git-commit-message.md) |
| `isolate-change-workspace` | Sesiones concurrentes en un repo dejan de compartir un único árbol de trabajo. El aislamiento anida en dos niveles: un workspace de cambio que el orquestador abre e integra, y workspaces por tarea que ya abre un batch de implementación — ahora adentro del workspace del cambio en lugar de la rama original. Es opt-in; desactivado, se comporta exactamente como antes | Accepted | [`isolate-change-workspace.md`](isolate-change-workspace.md) |
