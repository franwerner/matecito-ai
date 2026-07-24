# Catálogo de tecnologías

Registro vivo de las tecnologías concretas elegidas. Cada entrada apunta a un mini-EDR con el "por qué" y alternativas descartadas.

**Para Claude:** consultá esta tabla antes de sugerir agregar una nueva dependencia. Si lo que vas a agregar pisa con algo ya elegido, **no lo agregues sin preguntar**.

## Por categoría

### Lenguaje y runtime
| Tech | Versión | Por qué (resumen) |
|---|---|---|
| [Go](go.md) | sin pinear | Concurrencia idiomática, binario único autocontenido y cross-compilación trivial. |

### Framework principal
| Tech | Versión | Por qué |
|---|---|---|
| [huma](huma.md) | sin pinear | Genera OpenAPI 3.1 + validación desde structs sobre net/http; el contrato Go-first que consume la UI. |

### Base de datos
| Tech | Versión | Por qué |
|---|---|---|
| [modernc.org/sqlite](modernc-sqlite.md) | sin pinear | Driver SQLite en Go puro (sin cgo): cross-compilación trivial. |

### ORM / Acceso a datos
| Tech | Versión | Por qué |
|---|---|---|
| [sqlc](sqlc.md) | sin pinear | Queries type-safe generadas desde SQL a mano, sin ORM. |

### Testing
| Tech | Versión | Por qué |
|---|---|---|

### Logging
| Tech | Versión | Por qué |
|---|---|---|

### Configuración / Secretos
| Tech | Versión | Por qué |
|---|---|---|

### Auth
| Tech | Versión | Por qué |
|---|---|---|

### Inyección de dependencias
| Tech | Versión | Por qué |
|---|---|---|

### Otros
| Tech | Versión | Por qué |
|---|---|---|
| [golangci-lint](golangci-lint.md) | sin pinear | Enforcer del estilo por encima de gofmt + go vet. |
| [golang.org/x/sync](golang-x-sync.md) | sin pinear | Grupo de goroutines con cancelación propagada para el lifecycle. |
| [goose](goose.md) | sin pinear | Migraciones versionadas embebidas, aplicadas al arranque. |
| [coder/websocket](coder-websocket.md) | sin pinear | WebSocket para el WS-out hacia la UI; valida Origin==Host por default. |
| [SDK Go oficial de MCP](mcp-go-sdk.md) | sin pinear | SDK oficial de MCP para Go; soporta streamable HTTP para la superficie MCP cara a Claude. |

## Stdlib usada (sin entrada propia)

Estas piezas de la librería estándar de Go se usan pero no meritan un mini-EDR (no hubo elección entre alternativas): `slog` (logging estructurado, ver observability/logging), `database/sql` (ejecución de queries, ver data/data-access), `testing` (framework de tests, ver delivery/testing-strategy), `flag` y `os` (config por flags/env, ver delivery/configuration), `embed` (migraciones y bundle de la UI embebidos, ver data/data-access y delivery/deployment-topology), `net/http` (ServeMux del HTTP-in, sobre el que corre huma, ver contracts/api-contract).

## Mantenimiento

- **Agregar tech:** crear `<nombre>.md`, sumar fila en la categoría.
- **Reemplazar tech:** marcar el viejo `Superseded`, crear el nuevo, sacar del INDEX el viejo (o moverlo a "Históricas").
- **Actualizar versión:** editar el archivo, anotar en Notas si hay breaking changes.
