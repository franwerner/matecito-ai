---
name: cli-contract
depth: light
domain: contracts
type: decision
source: arc42 §8 (cross-cutting concepts) · POSIX utility conventions
---

# Fase: Contrato de CLI

## Qué decide

Cómo se parsean los argumentos, qué exit codes se usan, qué va a stdout vs stderr, y el formato de output.

## Preguntas

### 1. Parsing y salida

> Un CLI sin exit codes correctos no puede integrarse en scripts o pipelines. La mezcla de datos y errores en stdout rompe el parseo downstream.

- **Librería de parsing estándar del ecosistema + exit 0 ok / exit 1 error** — *default recomendado; ej: `cobra` (Go), `click` (Python), `commander` (Node), `clap` (Rust).*
- Parsing manual con `argv` — aceptable solo para scripts de un solo comando sin subcomandos.
- No sé, recomendame.

### 2. Formato de output

> Si el CLI es consumido por humanos y por scripts, el formato de salida define su composabilidad.

- **Texto plano legible para humanos** — *default para herramientas de uso directo.*
- JSON con `--json` flag — para CLIs que alimentan pipelines o dashboards.
- Mix: texto por default, JSON con flag.

## Tech a registrar

Si se elige una librería de parsing, registrarla en el catálogo `tech/`.

## Qué materializar

ADR `cli-contract` materializado según `../../templates/adr.md`. Debe contener:

- **Contexto**: por qué un CLI sin exit codes correctos no se integra en scripts o pipelines, y por qué mezclar datos y errores en stdout rompe el parseo downstream.
- **Decisión**: librería de parsing elegida (ej. cobra, click, commander, clap) o parsing manual, tabla de exit codes (0 = éxito, 1 = error genérico, y cualquier código semántico adicional), política de stdout vs stderr (datos vs diagnóstico/errores), y formato de output estándar (texto plano, JSON con `--json`, o mix).
- **Reglas verificables** (cada una con su mecanismo):
  - `[tool: test]` un comando exitoso retorna exit `0`; un error retorna el código de la tabla decidida.
  - `[tool: test]` los datos van a stdout y el diagnóstico/errores a stderr, sin mezclarlos.
  - `[manual]` con el flag de formato máquina (ej. `--json`), la salida es JSON parseable; sin él, texto plano legible.
- **Alternativas consideradas**: parsing manual con `argv` y los formatos de output evaluados, con su trade-off de composabilidad.
- **Consecuencias**: dependencia de la librería de parsing y disciplina para mantener la tabla de exit codes estable como contrato.
