---
name: cli-contract
depth: light
domain: contracts
tipo: decisión
adr-output: cli-contract
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

ADR `cli-contract` con: librería de parsing elegida, tabla de exit codes (0 = éxito, 1 = error genérico, y cualquier código semántico adicional), política de stdout vs stderr (datos vs diagnóstico/errores), y formato de output estándar.
