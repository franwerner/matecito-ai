# Capability specs — `process`

Comportamiento reactivo/de fondo, disparado por un evento o por el sistema (no por un actor): file-watchers, jobs, reconciliación, indexado.

**Cuándo consultar este tipo:** antes de tocar cualquier proceso de fondo del daemon —el file-watch e indexado de records, su versionado o su proyección—, el motor que despliega el payload en el host destino, o un validador que el host dispara ante un evento, leé el spec para conocer su disparador, su flujo, sus casos borde y sus reglas.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `deploy-payload-to-host` | Mapea los componentes de los dominios activos y los compartidos al host destino, compone el archivo de instrucciones raíz, detecta clashes y decide backup por estado de archivo | Accepted | [`deploy-payload-to-host.md`](deploy-payload-to-host.md) |
| `dev-install-local-build` | Ciclo de desarrollo en un solo comando: compila el binario del árbol y lo deja operativo en el entorno local sin que la instalación lo pise con la release | Accepted | [`dev-install-local-build.md`](dev-install-local-build.md) |
| `render-durable-artifact` | Construye el cuerpo completo de un artefacto durable (EDR, capability-spec) a partir de datos y del contrato declarado del tipo, garantizando conformidad por construcción | Accepted | [`render-durable-artifact.md`](render-durable-artifact.md) |
| `validate-artifact-structure` | Valida que un artefacto escrito se conforma a su contrato (secciones presentes, orden, cabeceras válidas, sincronía entre archivos); valida los contratos y templates contra sí mismos (`--self-check`); examina el alcance nombrado (un archivo, el store completo, o los artefactos embarcados), reporta hallazgos, no modifica nada | Accepted | [`validate-artifact-structure.md`](validate-artifact-structure.md) |
| `validate-git-commit-message` | Bloquea el commit ante atribución de IA, avisa (sin bloquear) si el mensaje no sigue Conventional Commits, y falla abierto ante cualquier ambigüedad | Inferred | [`validate-git-commit-message.md`](validate-git-commit-message.md) |
| `isolate-change-workspace` | Sesiones concurrentes en un repo dejan de compartir un único árbol de trabajo. El aislamiento existe en un solo nivel: los workspaces por tarea que abre un batch de implementación sobre la rama de trabajo, siempre; no hay workspace de nivel de cambio en ninguna parte del pipeline | Accepted | [`isolate-change-workspace.md`](isolate-change-workspace.md) |
| `install-mcp-from-release-tarball` | Registra un servidor de integración cuya única distribución es un tarball de release en GitHub — sin entrada en un registro de paquetes — durante la corrida de instalación del ecosistema, idempotentemente y sin llevar el resto de la corrida hacia abajo | Accepted | [`install-mcp-from-release-tarball.md`](install-mcp-from-release-tarball.md) |
| `configure-record-search` | Configuración única por proyecto del camino de búsqueda semántica: una colección por store de registros, sus exclusiones, sus credenciales de embedding y el primer índice | Accepted | [`configure-record-search.md`](configure-record-search.md) |
| `binary-install-step-idempotency` | Garantiza que el sistema decide la presencia de un binario observándolo donde su paso lo deja, en ambas superficies (instalación y reporte de estado), para que una máquina no reporte contradicciones ni repita descargas de binarios sanos | Accepted | [`binary-install-step-idempotency.md`](binary-install-step-idempotency.md) |
| `materialize-change-delta-spec` | El pliegue del delta de comportamiento de un cambio al store durable lo hace la fase que implementa, una sola vez, en el despacho que cierra las tareas — y el pliegue es no destructivo: preserva lo que no menciona el delta | Accepted | [`materialize-change-delta-spec.md`](materialize-change-delta-spec.md) |
