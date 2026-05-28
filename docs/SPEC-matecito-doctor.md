# matecito-doctor — Especificación (CLI de setup de matecito-ai)

> Spec para construir con Claude Code. Define requisitos, alcance y comportamiento. No es código — es el "qué" y el "por qué".

## 1. Propósito

CLI en Go, binario único, multiplataforma, que **verifica, inicia e instala** las dependencias del ecosistema matecito-ai: Engram, CodeGraph y los MCP servers (incluido context7), y valida que la configuración MCP coincida con lo que el fork del SDD espera.

Nombre propuesto del binario: `matecito` (o `matecito-doctor`).

## 2. Requisitos del usuario (prerequisites)

El CLI debe **detectar y reportar** la presencia de estos prerequisites. No los instala (son responsabilidad del usuario / gestor de paquetes del SO), pero avisa si faltan:

| Requisito | Para qué | Cómo se detecta |
|-----------|----------|-----------------|
| **Claude Code** (CLI `claude`) | Host de los MCP y agentes | `claude --version` en PATH |
| **Node.js 18+** | CodeGraph y context7 corren vía npm/npx | `node --version` ≥ 18 |
| **npm / npx** | Instalar CodeGraph, ejecutar context7 | `npm --version`, `npx --version` |
| **Homebrew** (macOS/Linux) | Método recomendado para instalar Engram | `brew --version` (solo si se usará ese método) |
| **Go 1.18+** *(opcional)* | Solo si se instala Engram vía `go install` | `go version` |
| **git** | Sync de Engram, operaciones de repo | `git --version` |

Plataformas objetivo: **Linux, macOS, Windows** (como gentle-ai). En Windows considerar rutas y el formato de `~/.claude.json`.

## 3. Componentes que gestiona

### 3.1 Engram (memoria de sesión)
- **Repo:** github.com/Gentleman-Programming/engram — Go binary, SQLite+FTS5, MCP server. Sin Node/Python/Docker.
- **Instalación (recomendada):** `brew install gentleman-programming/tap/engram`
- **Instalación (alternativas):** binario de releases de GitHub; `go install` desde el repo. (Dejar la vía brew como default y releases como fallback.)
- **Registro como MCP en Claude Code:** `claude plugin marketplace add Gentleman-Programming/engram && claude plugin install engram`
  - Alternativa genérica MCP: comando `engram`, args `["mcp"]`, transport stdio.
- **Verificación:** binario `engram` en PATH (`engram version`); DB en `~/.engram/engram.db`; MCP registrado.
- **Nota:** Engram expone 19 tools MCP con prefijo observado en el SDD `mcp__plugin_engram_engram__*` (ej: `mem_search`, `mem_save`, `mem_get_observation`, `mem_update`). El CLI debe validar que este prefijo coincida (ver §5).

### 3.2 CodeGraph (grafo de código)
- **Repo:** github.com/colbymchenry/codegraph — Node/TypeScript, MCP server, 100% local.
- **Instalación:** `npm install -g @colbymchenry/codegraph` (o el instalador interactivo `npx @colbymchenry/codegraph`).
- **Registro como MCP en `~/.claude.json`:**
  ```json
  { "mcpServers": { "codegraph": { "type": "stdio", "command": "codegraph", "args": ["serve", "--mcp"] } } }
  ```
- **Init por proyecto:** `codegraph init -i` (crea `.codegraph/` en el repo).
- **Verificación:** binario `codegraph` en PATH; entrada en `~/.claude.json`; existencia de `.codegraph/` en el proyecto actual (`codegraph status`).
- **Tools MCP esperadas (prefijo `mcp__codegraph__*`):** `codegraph_search`, `codegraph_explore`, `codegraph_context`, `codegraph_callers`, `codegraph_callees`, `codegraph_impact`, `codegraph_node`, `codegraph_status`, `codegraph_files`.

### 3.3 context7 (docs de librerías)
- **Paquete:** `@upstash/context7-mcp` — corre vía npx, transport stdio (o HTTP remoto).
- **Instalación/registro (recomendado, vía Claude Code):**
  `claude mcp add --scope user context7 -- npx -y @upstash/context7-mcp@latest`
  (Si requiere API key: `... --api-key YOUR_KEY`, o vía env `CONTEXT7_API_KEY`.)
- **Registro manual en `~/.claude.json`:**
  ```json
  { "mcpServers": { "context7": { "command": "npx", "args": ["-y", "@upstash/context7-mcp@latest"] } } }
  ```
- **Verificación:** entrada `context7` en `~/.claude.json`; resolución de `npx`.
- **Tools MCP esperadas (prefijo a CONFIRMAR):** el fork del SDD asume `mcp__context7__resolve_library_id` y `mcp__context7__query`. El nombre real depende del registro; el CLI debe reportar el prefijo real para que el usuario lo ajuste en el SDD (ver §5). Nombres documentados de las tools: `resolve-library-id` y `query-docs`.

## 4. Comandos (subcomandos del CLI)

### `matecito verify` (default, solo lectura — riesgo cero)
Reporta el estado completo, sin modificar nada:
- Prerequisites (§2): presencia y versión de cada uno; marcar faltantes.
- Engram: instalado? versión? DB existe? MCP registrado?
- CodeGraph: instalado? versión? MCP registrado? `.codegraph/` en el cwd?
- context7: MCP registrado?
- **Cross-check con el SDD (§5).**
- Salida tipo checklist con ✓ / ✗ / ⚠ y, para cada faltante, el comando exacto para resolverlo.
- Exit code 0 si todo OK, ≠0 si falta algo crítico (para uso en scripts).

### `matecito init`
Inicializa lo que es por-proyecto en el cwd:
- Si CodeGraph está instalado pero falta `.codegraph/` → corre `codegraph init -i`.
- Reporta qué inicializó. No toca config global.

### `matecito install`
Instala/registra lo que falte. **Acciones de riesgo medio — requieren confirmación y backup:**
- Antes de tocar `~/.claude.json`: **backup** a `~/.claude.json.bak.<timestamp>`.
- **Merge, no overwrite:** agregar entradas MCP sin pisar las existentes. Si una entrada ya existe, preguntar antes de reemplazar.
- Confirmación explícita (`--yes` para saltarla en modo no interactivo).
- `--dry-run`: mostrar qué haría sin ejecutar (como gentle-ai).
- Instala en orden: Engram (brew/releases) → CodeGraph (npm) → registra MCP (engram plugin, codegraph, context7).
- Reporta cada paso y su resultado; si uno falla, no continúa con los dependientes y deja el backup intacto.

### `matecito doctor`
`verify` + diagnóstico accionable: por cada problema, la causa probable y el comando para arreglarlo. (Inspirado en `engram doctor` / `codegraph status`.)

## 5. Cross-check con el fork del SDD (requisito clave)

El CLI debe **leer los agentes del SDD forkeado** (en `~/.claude/agents/sdd-*.md` o donde estén instalados) y extraer los nombres de tools MCP declarados en su frontmatter (`tools:`). Luego comparar contra los MCP realmente registrados en `~/.claude.json`:

- Si el SDD declara `mcp__codegraph__codegraph_impact` pero el server real expone otro prefijo → **⚠ MISMATCH**, reportar ambos nombres y sugerir corregir el frontmatter del agente o el registro del MCP.
- Igual para context7 y engram.
- Esto cierra el pendiente conocido de matecito-ai: "verificar que los nombres de tools MCP coincidan con los que el SDD espera".

Salida: una tabla `Tool esperada por el SDD | ¿Registrada? | Prefijo real | Estado`.

## 6. Requisitos no funcionales

- **Binario único, sin dependencias de runtime** (Go estático), multiplataforma (Linux/macOS/Windows; amd64 + arm64).
- **Idempotente:** correr `install` dos veces no duplica entradas ni rompe nada.
- **Seguro por defecto:** `verify` es la acción default; nada se instala/modifica sin `install` + confirmación.
- **Backup siempre** antes de escribir en `~/.claude.json`.
- **Salida legible** con códigos de color y exit codes usables en scripts.
- **Detección de SO/arch** para elegir el método de instalación correcto (brew en mac/linux, etc.).

## 7. Fuera de alcance (v1)

- No instala los prerequisites (Claude Code, Node, Go, brew) — solo los detecta y reporta.
- No gestiona Engram Cloud (solo el modo local).
- No configura los agentes del SDD (eso es el fork, ya hecho) — solo verifica coincidencia de nombres.
- No soporta otros agentes además de Claude Code en v1 (OpenCode, Cursor, etc. quedan para después).

## 8. Criterios de aceptación

- `matecito verify` en una máquina limpia reporta correctamente qué falta y cómo instalarlo, con exit code ≠0.
- `matecito verify` en una máquina configurada reporta todo ✓ con exit code 0.
- `matecito install --dry-run` muestra el plan sin tocar nada.
- `matecito install` hace backup de `~/.claude.json`, mergea sin pisar entradas previas, y deja Engram + CodeGraph + context7 registrados.
- `matecito init` crea `.codegraph/` en un proyecto que no lo tenía.
- El cross-check (§5) detecta un mismatch deliberado entre el frontmatter del SDD y el registro real.
- Corre en Linux, macOS y Windows.

## 9. Notas / riesgos

- **Nombres de tools MCP:** los prefijos (`mcp__codegraph__*`, `mcp__context7__*`, `mcp__plugin_engram_engram__*`) deben confirmarse contra la instalación real — son la fuente principal de mismatch. El §5 existe justamente para detectarlo.
- **`~/.claude.json` es config global:** un error de escritura puede romper TODOS los MCP del usuario, no solo los de matecito. El backup + merge son obligatorios, no opcionales.
- **context7 puede requerir API key** según el plan; el CLI debe detectar la ausencia y avisar, no asumir que funciona sin ella.
- **Madurez:** Engram (v1.15.x, activo) y CodeGraph (sin releases formales) están en evolución; los comandos de instalación pueden cambiar. El CLI debería tolerar fallos de instalación con mensajes claros, no asumir éxito.
