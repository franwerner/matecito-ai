# 07 — Configuración

[← 06 Herramientas](06-herramientas.md) · [Índice](README.md)

Los ajustes del flujo SDD viven en archivos `config.json` y se editan desde la **TUI** (`matecito-ai`).

## Dónde vive y precedencia

Resolución **por-proyecto → global → default**:

- `<repo>/.matecito-ai/config.json` — config específica del proyecto.
- `~/.matecito-ai/config.json` — config global, fallback cuando el proyecto no la define.

```json
{
  "models": { "sdd-design": "opus", "sdd-apply": "sonnet" },
  "strictTdd": false,
  "flagDecisionGaps": false
}
```

## Qué se configura

- **Modelo por agente** (`models`) — qué modelo usa cada fase del SDD (`sdd-intake`, `sdd-spec`, `sdd-design`, …). Sin valor configurado, cada agente usa su default curado (no el modelo de la conversación).
- **Strict TDD** (`strictTdd`) — si está activo, `apply` y `verify` siguen el ciclo test-first.
- **Auto-mine ADR** (`flagDecisionGaps`) — opt-in, off por default. Activa la detección de decisiones implementadas sin ADR durante el flujo; al cerrar, ofrece minarlas como `Inferred`. Ver [05](05-auto-mine.md).
- **Scope** — si los cambios de configuración desde la TUI aplican al proyecto actual o al global.

## Quién resuelve la config (no los ejecutores)

Clave de diseño: **los sub-agentes ejecutores no leen `config.json`.** El **orquestador** resuelve `model` / `strictTdd` / `flagDecisionGaps` (por la precedencia de arriba, una vez por sesión, cacheado) **antes** de despachar cada fase, y le pasa el valor ya resuelto al sub-agente. Es el mismo patrón para los tres: una sola fuente de resolución, ejecutores que reciben el valor.

- `strictTdd` → se inyecta en el prompt de `apply`/`verify` si está activo.
- `flagDecisionGaps` → habilita los hooks de gap en `tasks`/`verify` y el mine gate del boundary.
- `models[<agente>]` → se pasa como el modelo de ese sub-agente; si no hay valor, se omite y aplica el default del agente.

## La TUI

Sin subcomando y en terminal interactiva, `matecito-ai` abre una TUI con el estado del entorno, la instalación, y la **Configuración** (Modelos por agente · Strict TDD · Auto-mine ADR · Scope). Los toggles son scope-aware (editan el config global o el del proyecto según el scope activo).

```bash
matecito-ai            # TUI
matecito-ai verify     # estado del entorno
matecito-ai install    # instalar/actualizar lo que falte
```
