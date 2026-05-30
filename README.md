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

**matecito-ai** es un ecosistema de desarrollo asistido por IA, armado a medida sobre [Claude Code](https://claude.com/claude-code). No es una herramienta nueva: es la integración curada de varias piezas —propias y de terceros— en un flujo coherente, donde cada decisión de arquitectura queda registrada y respetada a lo largo del tiempo y entre sesiones.

La idea de fondo: que el agente **no reinvente las convenciones del proyecto en cada sesión**. Las decisiones se capturan una vez (como ADRs), se respetan al implementar, y la memoria de trabajo persiste vía Engram. El humano decide; la IA ejecuta dentro de esas decisiones.

## Qué hace

Trabajar con agentes de IA sobre un proyecto tiene tres fugas recurrentes, y matecito-ai ataca las tres:

- **Amnesia entre sesiones** → **Engram** persiste la memoria de trabajo (descubrimientos, contexto, fixes) entre sesiones.
- **Decisiones implícitas** → **ADRs** capturan las decisiones de arquitectura una vez; el flujo las respeta y avisa si algo las contradice.
- **Exploración cara** → **codegraph** indexa el código para explorarlo por estructura, sin escanear archivo por archivo.

Sobre eso corre un **flujo de desarrollo guiado (SDD)** que lleva cada cambio de un pedido en lenguaje natural hasta el código, pasando por fases con un punto de control humano al inicio.

## Componentes

| Capa        | Pieza                          | Rol                                                                                  |
|-------------|--------------------------------|--------------------------------------------------------------------------------------|
| **Skills**     | `project-decisions-bootstrap`  | Entrevista por fases que captura decisiones de ingeniería y las materializa como ADRs por dominio. |
| **Skills**     | `project-decisions-validate`   | Validador consultivo: coherencia, completitud y verificabilidad de los ADRs.         |
| **Skills**     | `SDD` *(fork del Gentleman)*   | Flujo de fases: intake → explore → propose → spec → design → tasks → apply → verify → archive. |
| **Referencia** | `design-patterns`              | Catálogo canónico de patrones de diseño consultable. Los ADRs lo citan por nombre; `sdd-design` lo respeta cuando un ADR declara `Patrón aplicado`. |
| **MCP**        | `codegraph`                    | Grafo de código pre-indexado (tree-sitter + SQLite) para explorar por estructura.    |
| **MCP**        | `context7`                     | Documentación de librerías al día, contra APIs no alucinadas.                        |
| **Agentes**    | Sub-agentes del SDD            | Uno por fase, con contexto propio. Forkeados y modificados.                          |
| **Engram**     | Memoria persistente            | SQLite standalone con descubrimientos, contexto y fixes entre sesiones.              |

## El flujo SDD

```
intake → explore → propose → spec → design → tasks → apply → verify → archive
   │                                   │                          │
   │ estructura el pedido,             │ lee los ADRs vigentes    │ chequea que el código
   │ pregunta lo que falta,            │ y respeta los Accepted   │ respete los ADRs que tocó
   │ y FRENA para confirmar alcance    │                          │
```

- **intake** es la fase de entrada: hace 2-4 preguntas para estructurar el pedido, lo clasifica, y produce un brief. El orquestador **siempre muestra ese brief y espera tu confirmación** antes de seguir.
- **design** y **apply** leen los ADRs vigentes; **explore** usa codegraph; **apply** usa context7.
- Cuando un ADR declara `Patrón aplicado: X`, **design** consulta el catálogo `design-patterns` y respeta la definición canónica del patrón.
- **verify** confirma que el cambio no viole los ADRs que tocó.

## Instalación

Requisitos:
- **[Claude Code](https://claude.com/claude-code)** instalado y autenticado
- **Node.js `≥ 18`** con `npm` y `npx` — CodeGraph se instala una vez con `npm install -g`; context7 se invoca runtime en cada sesión con `npx -y @upstash/context7-mcp@latest`
- **git `≥ 2.23`** — el workflow SDD asume historial git para versionar ADRs y commits; la skill de git usa `git restore` (introducido en 2.23)

Engram se descarga como binario precompilado desde sus [GitHub Releases](https://github.com/Gentleman-Programming/engram/releases); no requiere Go.

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

Requiere **Go `≥ 1.22`**:

```bash
go build -o matecito-ai ./cmd/matecito-ai
```

## Uso

El CLI verifica, inicia e instala las dependencias del ecosistema (Engram, codegraph, context7) sobre Claude Code, y deploya el fork del SDD a `~/.claude/`. Una vez instalado, cada herramienta se usa con su propio binario; matecito-ai se ocupa del setup y la salud del entorno.

```bash
# Reportar estado del entorno (qué está instalado / registrado)
matecito-ai verify


# Instalar lo que falte (con backup de la config y confirmación)
matecito-ai install --dry-run
matecito-ai install

# Actualizar todo el ecosistema a sus últimas versiones (matecito-ai, Engram, CodeGraph)
matecito-ai update
```

## Documentación

- [PRD](docs/PRD.md) — documento de producto del ecosistema.

## Créditos

matecito-ai es una capa de integración sobre proyectos de terceros. El crédito del trabajo pesado es de ellos:

- **[gentle-ai](https://github.com/Gentleman-Programming/gentle-ai)** — el motor de Spec-Driven Development y TDD que matecito-ai forkea y adapta.
- **[engram](https://github.com/Gentleman-Programming/engram)** — memoria persistente entre sesiones.
- **[codegraph](https://github.com/colbymchenry/codegraph)** — grafo de conocimiento del código vía tree-sitter, expuesto como MCP.
- **[context7](https://github.com/upstash/context7)** — documentación de librerías en vivo, expuesta como MCP.
