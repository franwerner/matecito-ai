<p align="center">
  <img src="docs/visual/brand.png" alt="matecito-ai" width="640">
</p>

<p align="center">
  <em>Mientras la IA trabaja por vos, te tomás unos ricos mates.</em>
</p>

<p align="center">
  <a href="#instalación"><img alt="status" src="https://img.shields.io/badge/status-alpha-orange"></a>
  <img alt="go" src="https://img.shields.io/badge/go-1.22-blue">
  <img alt="license" src="https://img.shields.io/badge/license-private-lightgrey">
</p>

---

**matecito-ai** es un ecosistema de desarrollo asistido por IA, armado a medida sobre [Claude Code](https://claude.com/claude-code). No es una herramienta nueva: es la integración curada de varias piezas —propias y de terceros— en un flujo coherente, donde cada decisión de arquitectura queda registrada y respetada a lo largo del tiempo y entre sesiones.

La idea de fondo: que el agente **no reinvente las convenciones del proyecto en cada sesión**. Las decisiones se capturan una vez (como ADRs), se respetan al implementar, y la memoria de trabajo persiste vía Engram. El humano decide; la IA ejecuta dentro de esas decisiones.

## Problema que resuelve

Trabajar con agentes de IA sobre un proyecto tiene tres fugas recurrentes:

1. **Amnesia entre sesiones.** Cada sesión nueva arranca sin saber qué se decidió antes — re-sugiere librerías ya descartadas, reinventa convenciones, contradice decisiones previas.
2. **Decisiones implícitas.** La arquitectura vive en la cabeza del autor, no escrita. El agente no tiene cómo respetarla.
3. **Exploración cara.** El agente gasta tokens y tool calls escaneando el codebase archivo por archivo para entender estructura que podría consultarse de forma indexada.

matecito-ai ataca las tres: **ADRs** para las decisiones, **Engram** para la memoria de sesión, **codegraph** para la exploración eficiente.

## Componentes

| Capa        | Pieza                          | Rol                                                                                  |
|-------------|--------------------------------|--------------------------------------------------------------------------------------|
| **Skills**  | `project-decisions-bootstrap`  | Entrevista por fases que captura decisiones de ingeniería y las materializa como ADRs por dominio. |
| **Skills**  | `project-decisions-validate`   | Validador consultivo: coherencia, completitud y verificabilidad de los ADRs.         |
| **Skills**  | `issue-brief`                  | Puente entre un issue y los ADRs que le aplican; arma briefing de restricciones.     |
| **Skills**  | `SDD` *(fork del Gentleman)*   | Workflow de fases: explore → propose → spec → design → tasks → apply → verify → archive. |
| **MCP**     | `codegraph`                    | Grafo de código pre-indexado (tree-sitter + SQLite) para explorar por estructura.    |
| **MCP**     | `context7`                     | Documentación de librerías al día, contra APIs no alucinadas.                        |
| **Agentes** | Sub-agentes del SDD            | Uno por fase, con contexto propio. Forkeados y modificados.                          |
| **Engram**  | Memoria persistente            | SQLite standalone con descubrimientos, contexto y fixes entre sesiones.              |

## Instalación

Requisitos:
- Go `1.22+`
- [Claude Code](https://claude.com/claude-code) instalado y autenticado

Build local:

```bash
go build -o matecito-ai ./cmd/matecito-ai
```

## Uso

El CLI verifica, inicia e instala las dependencias del ecosistema (Engram, CodeGraph, context7) sobre Claude Code, y deploya el fork del SDD a `~/.claude/`.

```bash
# Reportar estado del entorno
matecito-ai verify

# Instalar todo lo que falte (prereqs detectados auto)
matecito-ai install --dry-run
matecito-ai install

# Inicializar el ecosistema en un proyecto
matecito-ai init
```

## Flujo típico

1. **Setup del proyecto** (una vez): `project-decisions-bootstrap` captura las decisiones → genera `.claude/adr/` + `CLAUDE.md`.
2. **Al implementar un issue:** se describe el issue → el flujo SDD lo lleva por sus fases. `design` y `apply` leen los ADRs vigentes; `explore` usa codegraph; `apply` usa context7. Si el issue choca con un ADR o destapa una decisión nueva, se frena y se captura vía bootstrap antes de codear.
3. **Al cerrar:** `verify` chequea que el cambio respete los ADRs que tocó. Engram guarda lo aprendido.
4. **Mantenimiento del catálogo:** concerns nuevos se agregan vía `CONCERN-TEMPLATE.md`; coherencia entre ADRs se revisa con `project-decisions-validate`.

## Decisiones de diseño

| # | Decisión | Razón |
|---|----------|-------|
| 1 | **ADRs = decisiones de arquitectura; Engram = memoria de sesión** | Evitar solapamiento. ADR es decisión deliberada y verificable; Engram es lo que el agente aprendió trabajando. |
| 2 | **Fork directo del SDD, no inyección** | Máxima personalización. Se acepta el costo de mantenimiento a cambio de control total y coherencia con la forma propia de trabajar. |
| 3 | **Engram-only en el SDD (sin openspec/hybrid)** | El proyecto no usa persistencia basada en archivos del SDD; Engram cubre la memoria. |
| 4 | **Sin mecanismo de inyección (registry/resolver removidos)** | Las convenciones se leen de los archivos del proyecto (ADRs, CLAUDE.md), no de un registry intermedio. |
| 5 | **codegraph-first en exploración, grep como fallback** | codegraph para estructura/relaciones (eficiente); grep para texto literal o archivos no indexados. |
| 6 | **`design` y `apply` leen los ADRs; `verify` los chequea** | Las decisiones se respetan en el momento de diseñar e implementar; verify confirma que el código no las viole. |
| 7 | **No construir un issue-implementer propio** | El SDD ya es el implementador disciplinado. Construir otro sería duplicar infraestructura madura. |

## Fuera de alcance

- No es un instalador de agentes ni un producto distribuible — es un entorno personal.
- No reconstruye memoria persistente ni el orquestador de fases desde cero (se integran Engram y el SDD).
- No incluye los componentes de gentle-ai que no se usan (persona, GGA, skills de frameworks específicos).
- El enforcement de que el código respete los ADRs se delega a herramientas determinísticas (linters vía `arch-enforcement`) + `verify`.

## Riesgos conocidos

- **Frescura de los ADRs.** Si quedan desactualizados respecto al código, el agente trabaja sobre decisiones viejas con confianza. Mitigación: actualizar la decisión vía bootstrap update.
- **Frescura del índice de codegraph.** Un grafo desactualizado en el que el agente confía ciegamente es peor que no tenerlo. Mitigación: `codegraph status` antes de sesiones importantes.
- **Divergencia del upstream.** El fork se alejó del SDD original; portar mejoras del Gentleman será trabajo manual. Mitigación: `vendor-original/` + marcadores `matecito-ai`.
- **Madurez de las dependencias.** gentle-ai/SDD está en `v0.1.x` (APIs cambiarán).

## Documentación

- [PRD completo](docs/PRD.md)

## Estado

Alpha. Uso personal del autor. APIs y convenciones pueden cambiar sin previo aviso.
