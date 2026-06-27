<p align="center">
  <img src="docs/visual/brand.png" alt="matecito-ai" width="640">
</p>

<p align="center">
  <em>Mientras la IA trabaja por vos, te tomás unos ricos mates.</em>
</p>

<p align="center">
  <a href="#instalación"><img alt="status" src="https://img.shields.io/badge/status-alpha-orange"></a>
</p>

---

**matecito-ai** es un ecosistema de trabajo asistido por IA, armado a medida sobre [Claude Code](https://claude.com/claude-code). No es una herramienta nueva: es la integración curada de varias piezas —propias y de terceros— en una **disciplina coherente**, donde cada decisión queda registrada y respetada a lo largo del tiempo y entre sesiones.

Su arquitectura es un **microkernel**: un **núcleo agnóstico** (la disciplina) más **un plugin por área** (development, design, …). El núcleo aporta el flujo estructurado, el punto de control humano, la memoria persistente y las decisiones capturadas; cada **dominio** aporta su vocabulario, sus fases, sus skills y sus herramientas. Instalás los dominios que te tocan y conviven bajo el mismo ecosistema.

La idea de fondo: que el agente **no reinvente las convenciones en cada sesión**. Las decisiones se capturan una vez, se respetan al ejecutar, y la memoria de trabajo persiste vía Engram. El humano decide; la IA ejecuta dentro de esas decisiones.

## Qué hace

Trabajar con agentes de IA tiene tres fugas recurrentes, y matecito-ai ataca las tres — el núcleo da el mecanismo, cada dominio lo concreta:

- **Amnesia entre sesiones** → **Engram** persiste la memoria de trabajo (descubrimientos, contexto, fixes) entre sesiones. _(núcleo, compartido por todos los dominios)_
- **Decisiones implícitas** → **decision records** capturan las decisiones una vez; el flujo las respeta y avisa si algo las contradice. Cada dominio define su propio tipo de record.
- **Exploración cara** → cada dominio explora por su **índice** propio, sin escanear a ciegas.

Sobre eso corre un **flujo estructurado** que lleva cada cambio de un pedido en lenguaje natural hasta el entregable, pasando por fases con un **punto de control humano al inicio**. Las fases las define cada dominio; el núcleo orquesta.

## Dominios

Un **dominio** es un área de trabajo: un plugin sobre el núcleo. Se eligen e instalan desde la TUI (_Configuración → Dominios_), conviven a la vez, y cada uno se **auto-configura** desde su contrato.

| Dominio | Para qué | Detalle |
| --- | --- | --- |
| **development** | Desarrollo de software con SDD sobre un repo de código. | [README](payload/domains/development/README.md) |
| **design** | Diseño visual full-spectrum (marca + UI/UX + prototipos + guías). No toca código. | [README](payload/domains/design/README.md) |

Agregar un dominio nuevo (marketing, video, contable, …) es **implementar el contrato de área**, sin tocar el núcleo — ver [cómo se arma un dominio](payload/domains/README.md).

## Componentes

### Núcleo (compartido por todos los dominios)

| Pieza | Rol |
| --- | --- |
| **Orquestador** | Coordina el flujo: delega cada fase, sintetiza, y mantiene el **gate humano** (siempre confirmás el alcance antes de seguir). |
| **Lane fork** | El flujo es una **base inmutable** + **add-ons opcionales**; el tamaño del cambio decide qué fases corren (`direct \| reduced \| full \| custom`). |
| **Engram** | Mecanismo de memoria persistente del orquestador: SQLite standalone con descubrimientos, contexto y fixes entre sesiones. Como todo, su instalación la **declara cada dominio** en su manifest (no es global). |
| **Contrato de área** | El `manifest.json` + el fragmento `CLAUDE.md` de cada dominio: cómo un plugin se enchufa al núcleo (vocabulario, fases, guards, `mcp`, `binaries`, config). |
| **Tier compartido** | Skills, agentes y referencias en `payload/shared/` se despliegan **siempre**, sin importar los dominios activos. Los hooks siempre activos usan `hook.SharedDomain = "shared"` (compiled-in, sin directorio en el payload). Catálogo de lo que entrega: [README del tier compartido](payload/shared/README.md). |

> **Nada se instala global.** Cada dominio declara en su `manifest.json` los `mcp` y `binaries` que usa; el ecosistema instala la **unión** de los dominios activos. Los **permisos** de Claude Code se infieren de esa declaración (`name → mcp__<name>__*`); el único permiso global fijo es `Skill`.

### Por dominio

Lo específico de cada área —su flujo, su tipo de decision record, su índice de exploración, sus fases, agentes, skills, guards, catálogos y MCP— vive en su propio README. Cada plugin se documenta a sí mismo: **[development](payload/domains/development/README.md)** · **[design](payload/domains/design/README.md)**.

## El flujo estructurado

Todo cambio sustancial pasa por un flujo de fases con un **punto de control humano al inicio**. El esqueleto es genérico; cada dominio nombra sus fases.

```
intake → … → verify → archive
   │
   │ estructura el pedido, pregunta lo que falta,
   │ y FRENA para confirmar el alcance (INTAKE GATE)
```

- **intake** es la fase de entrada: hace 2-4 preguntas, clasifica el pedido y produce un brief. El orquestador **siempre muestra ese brief y espera tu confirmación** antes de seguir.
- El flujo es una **base inmutable** (las fases obligatorias del dominio) más **add-ons opcionales**. No todo cambio recorre todas las fases: un fix trivial va directo; un cambio grande las activa todas (**lane fork**).
- Las **decisiones** que se toman quedan como **decision records** (cada dominio define su tipo), y las fases posteriores las **respetan y verifican**.
- Cada fase corre con su **modelo configurable por agente** y los **guards** del dominio (ver [Configuración](#configuración)).

El pipeline concreto de cada dominio —sus nombres de fase, su tipo de decision record y sus herramientas— está en su propio README: [development](payload/domains/development/README.md) · [design](payload/domains/design/README.md).

## Instalación

Requisitos:

- **[Claude Code](https://claude.com/claude-code)** instalado y autenticado
- **Node.js `≥ 18`** con `npm` y `npx` — lo necesitan los MCP/binarios que se invocan vía Node, según los **dominios activos**: en development, context7 (`npx -y @upstash/context7-mcp@latest`), CodeGraph y proofshot (`npm install -g`) y drawio (`claude mcp add … npx -y @next-ai-drawio/mcp-server@latest`, que abre un preview en el navegador en `localhost:6002`); en design, Figma (`claude mcp add`). Nada de esto se instala si su dominio no está activo.

Engram (binario de memoria, declarado por cada dominio) se descarga precompilado desde sus [GitHub Releases](https://github.com/Gentleman-Programming/engram/releases); no requiere Go.

### Instalación rápida (Linux y macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/franwerner/matecito-ai/master/scripts/install.sh | bash
```

El script detecta tu OS/arch, baja el binario apropiado desde la última release, y lo instala en `~/.local/bin/matecito-ai`.

Variables de entorno opcionales:

- `INSTALL_DIR=/usr/local/bin` para cambiar el destino (puede requerir `sudo`).
- `VERSION=v0.1.0` para pinear una versión específica en vez de `latest`.

### Descarga manual (todos los SO)

1. Andá a [Releases](https://github.com/franwerner/matecito-ai/releases) y descargá el asset que corresponda:
   - `matecito-ai_<version>_linux_amd64.tar.gz`
   - `matecito-ai_<version>_linux_arm64.tar.gz`
   - `matecito-ai_<version>_darwin_amd64.tar.gz` (macOS Intel)
   - `matecito-ai_<version>_darwin_arm64.tar.gz` (macOS Apple Silicon)
   - `matecito-ai_<version>_windows_amd64.zip`
   - `matecito-ai_<version>_windows_arm64.zip`
2. Extraé el archivo y mové el binario `matecito-ai` (o `matecito-ai.exe` en Windows) a una carpeta que esté en tu `PATH`.

### Build desde fuente

Requiere **Go `≥ 1.24`**:

```bash
go build -o matecito-ai ./cmd/matecito-ai
```

## Uso

El CLI verifica, inicia e instala las dependencias del ecosistema —la **unión** de los `mcp` y `binaries` declarados por los **dominios activos** (en development: engram, context7, codegraph, drawio + binarios engram/codegraph/proofshot; en design: engram, figma, canva)— sobre Claude Code, y deploya el núcleo + los dominios activos a `~/.claude/`. Nada se instala si su dominio no está activo. Una vez instalado, cada herramienta se usa con su propio binario; matecito-ai se ocupa del setup y la salud del entorno.

Sin subcomando, y en una terminal interactiva, abre una **TUI** desde donde ves el estado, instalás, **elegís qué dominios tener activos** (_Configuración → Dominios_) y configurás cada dominio. Con subcomando corre en modo directo.

```bash
# TUI interactiva (estado, instalación, selección y config de dominios)
matecito-ai

# Reportar estado del entorno (qué está instalado / registrado)
matecito-ai verify

# Instalar/actualizar lo que falte (binarios, MCPs, núcleo + dominios) con backup de la config
matecito-ai install --dry-run   # solo muestra el plan, no ejecuta nada
matecito-ai install             # aplica los cambios
matecito-ai install --yes       # sin confirmación interactiva (CI)
```

`install` es la única ruta de instalación y actualización: detecta qué falta o está desactualizado (matecito-ai, Engram, los MCP y binarios de los dominios activos, el núcleo y los fragmentos de cada dominio) y lo deja al día en un solo paso.

## Configuración

Los ajustes viven en archivos `config.json`, resueltos **por-proyecto → global → default**:

- `<repo>/.matecito-ai/config.json` — config específica del proyecto (no se versiona).
- `~/.matecito-ai/config.json` — config global, fallback cuando el proyecto no la define.

Se editan desde la TUI (`matecito-ai` → _Configuración_), organizada en **General** (compartido) + **una entrada por dominio activo**:

- **Compartido** — la **selección de dominios** (qué áreas tenés activas).
- **Por dominio** — cada dominio expone su propia config desde su contrato (`domainConfig.<dominio>`):
  - **Modelo por agente** — qué modelo usa cada fase del dominio. Sin valor configurado, cada agente usa su default curado.
  - **Guards** — p. ej. **Strict TDD** (test-first en `apply`/`verify`) en development.
  - **Auto-mine** (`flagDecisionGaps`) — opt-in, off por default. Detecta decisiones implementadas sin decision record durante el flujo; al cerrar, ofrece minarlas como records `Inferred` (siempre con tu confirmación). **Auto-mine ADR** en development, **Auto-mine DDR** en design.

La config de un dominio solo aparece si el dominio está **activo**.

## Documentación

- **Dominios:** [development](payload/domains/development/README.md) · [design](payload/domains/design/README.md) · [cómo armar un dominio](payload/domains/README.md)
- [Guía profunda del flujo SDD](docs/guide/README.md) — el dominio **development** de punta a punta: fases, herramientas y la capa de decisiones (bootstrap / validate / mine).
- [PRD](docs/PRD.md) — documento de producto del ecosistema.
- [Guía: agregar una dependencia](docs/workflow-dependecy.md) — cómo integrar una pieza nueva al ecosistema.

## Créditos

matecito-ai es una capa de integración sobre proyectos de terceros. El crédito del trabajo pesado es de ellos:

- **[gentle-ai](https://github.com/Gentleman-Programming/gentle-ai)** — el motor de Spec-Driven Development y TDD que matecito-ai forkea y adapta.
- **[engram](https://github.com/Gentleman-Programming/engram)** — memoria persistente entre sesiones.
- **[codegraph](https://github.com/colbymchenry/codegraph)** — grafo de conocimiento del código vía tree-sitter, expuesto como MCP.
- **[context7](https://github.com/upstash/context7)** — documentación de librerías en vivo, expuesta como MCP.
- **[next-ai-draw-io](https://github.com/DayuanJiang/next-ai-draw-io)** — generación de diagramas draw.io desde lenguaje natural, expuesta como MCP.
- **[mcp-debugger](https://github.com/debugmcp/mcp-debugger)** — debugging step-through headless sobre DAP (multi-lenguaje), expuesto como MCP.
- **[proofshot](https://github.com/AmElmo/proofshot)** — verificación visual de UI grabando sesiones de browser, integrada como CLI.
