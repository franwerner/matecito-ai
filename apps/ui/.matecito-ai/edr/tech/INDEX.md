# Catálogo de tecnologías

Registro vivo de las tecnologías concretas elegidas. Cada entrada apunta a un mini-EDR con el "por qué" y alternativas descartadas.

**Para Claude:** consultá esta tabla antes de sugerir agregar una nueva dependencia. Si lo que vas a agregar pisa con algo ya elegido, **no lo agregues sin preguntar**.

## Por categoría

### Lenguaje y runtime
| Tech | Versión | Por qué (resumen) |
|---|---|---|
| [TypeScript](typescript.md) | sin pinear | Tipado estricto de toda la UI; habilita los tipos generados desde el OpenAPI del broker. |

### Framework principal
| Tech | Versión | Por qué |
|---|---|---|
| [React](react.md) | sin pinear | Base de la SPA del cockpit (canvas, inspector, timeline). |

### Base de datos
| Tech | Versión | Por qué |
|---|---|---|

### ORM / Acceso a datos
| Tech | Versión | Por qué |
|---|---|---|

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
| [Vite](vite.md) | sin pinear | Build tool + dev server; produce el build estático embebido y el proxy de dev. |
| [pnpm](pnpm.md) | sin pinear | Gestor de paquetes; su `audit` es gate de vulnerabilidades en CI. |
| [ReactFlow](reactflow.md) | sin pinear | Motor del canvas; nodos/aristas derivados de la proyección del event-log. |
| [TanStack Query](tanstack-query.md) | sin pinear | Cache del server/remote state (snapshot HTTP + eventos WS). |
| [TanStack Router](tanstack-router.md) | sin pinear | Router file-based con search params tipados y loaders que prefetchean. |
| [TanStack Table](tanstack-table.md) | sin pinear | Modelo headless de tablas para vistas tabulares. |
| [TanStack Virtual](tanstack-virtual.md) | sin pinear | Virtualización de listas/feed (mitigación de performance). |
| [Zustand](zustand.md) | sin pinear | Store del estado efímero de UI. |
| [Tailwind CSS](tailwindcss.md) | sin pinear | Motor de estilos utility-first + theming por CSS variables. |
| [shadcn/ui](shadcn-ui.md) | sin pinear | Componentes Radix copiados y editables; base del component base. |
| [Radix UI](radix-ui.md) | sin pinear | Primitives accesibles (foco, ARIA, teclado) bajo shadcn. |
| [sonner](sonner.md) | sin pinear | Toasts para errores transitorios. |
| [partysocket](partysocket.md) | sin pinear | Cliente WS único: reconnect + backoff + buffer + heartbeat. |
| [Zod](zod.md) | sin pinear | Validación por schema en runtime (entrada del broker + config). |
| [Kubb](kubb.md) | sin pinear | Genera tipos TS + schemas Zod desde el OpenAPI del broker. |
| [ESLint](eslint.md) | sin pinear | Enforcer de las convenciones de código; gate de merge. |
| [Prettier](prettier.md) | sin pinear | Formateador; su diff bloquea el merge. |
| [eslint-plugin-jsx-a11y](eslint-plugin-jsx-a11y.md) | sin pinear | Chequeo de accesibilidad en dev-time; falla el build. |
| [husky](husky.md) | sin pinear | Hooks de git para los checks en pre-commit. |
| [lint-staged](lint-staged.md) | sin pinear | Corre los checks solo sobre los archivos staged. |

## Nota

**Dependabot** no es un paquete sino un servicio de GitHub (PRs automáticos de seguridad para las deps npm): se documenta en security/dependency-scanning, no como entrada de este catálogo.

## Mantenimiento

- **Agregar tech:** crear `<nombre>.md`, sumar fila en la categoría.
- **Reemplazar tech:** marcar el viejo `Superseded`, crear el nuevo, sacar del INDEX el viejo (o moverlo a "Históricas").
- **Actualizar versión:** editar el archivo, anotar en Notas si hay breaking changes.
