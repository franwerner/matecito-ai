# 08 — Configuración

[← 07 Herramientas](07-herramientas.md) · [Índice](README.md)

Los ajustes viven en archivos `config.json` y se editan desde la **TUI** (`matecito-ai`). Desde la migración a multi-dominio, la config está **organizada por dominio**: lo compartido arriba, y lo de cada área bajo `domainConfig.<dominio>`.

## Dónde vive y precedencia

Resolución **por-proyecto → global → default**:

- `<repo>/.matecito-ai/config.json` — config específica del proyecto.
- `~/.matecito-ai/config.json` — config global, fallback cuando el proyecto no la define.

```json
{
  "domains": ["development", "design"],
  "domainConfig": {
    "development": {
      "models": { "sdd-design": "opus", "sdd-apply": "sonnet" },
      "strictTdd": false,
      "flagDecisionGaps": false
    },
    "design": {
      "models": { "design-system": "opus" },
      "flagDecisionGaps": false
    }
  }
}
```

> **Compatibilidad:** un `config.json` plano de antes de multi-dominio (con `models` / `strictTdd` / `flagDecisionGaps` en la raíz) se **migra automáticamente** a `domainConfig.development` al leerlo —no hace falta tocar nada— y se persiste en el formato nuevo en el próximo guardado.

## Qué se configura

### Compartido (cross-dominio)

- **Dominios activos** (`domains`) — qué áreas tenés instaladas/activas. Vacío = todos los presentes en el payload (shim de compatibilidad). Gobierna el deploy, los MCP que se registran y qué config aparece.

### Por dominio (`domainConfig.<dominio>`)

- **Modelo por agente** (`models`) — qué modelo usa cada fase **de ese dominio** (`sdd-*` en development, `design-*` en design). Sin valor configurado, cada agente usa su default curado (no el modelo de la conversación).
- **Guards** — los controles propios del dominio. En development: **Strict TDD** (`strictTdd`) — si está activo, `apply` y `verify` siguen el ciclo test-first.
- **Auto-mine** (`flagDecisionGaps`) — opt-in, off por default. Detecta decisiones implementadas sin decision record durante el flujo; al cerrar, ofrece minarlas como `Inferred`. **Auto-mine EDR** en development, **Auto-mine DDR** en design. Ver [05](05-auto-mine.md).

**Scope** — si los cambios desde la TUI aplican al proyecto actual o al global.

## Quién resuelve la config (no los ejecutores)

Clave de diseño: **los sub-agentes ejecutores no leen `config.json`.** El **orquestador** resuelve `model` / `strictTdd` / `flagDecisionGaps` (por la precedencia de arriba, una vez por sesión, cacheado) **antes** de despachar cada fase, leyendo el path **por-dominio**, y le pasa el valor ya resuelto al sub-agente. Es el mismo patrón para los tres: una sola fuente de resolución, ejecutores que reciben el valor.

- `models[<agente>]` → se lee de `domainConfig[<dominio del agente>].models[<agente>]` y se pasa como el modelo de ese sub-agente; si no hay valor, se omite y aplica el default del agente. (El dominio del agente sale de quién lo trae: `sdd-*` → `development`, `design-*` → `design`.)
- `strictTdd` → `domainConfig.development.strictTdd`; se inyecta en el prompt de `apply`/`verify` si está activo.
- `flagDecisionGaps` → `domainConfig[<dominio>].flagDecisionGaps`; habilita los hooks de gap en `tasks`/`verify` y el mine gate del boundary.

## La TUI

Sin subcomando y en terminal interactiva, `matecito-ai` abre una TUI con el estado del entorno, la instalación, y la **Configuración**, organizada en **General** + **una entrada por dominio activo**:

- **General** — la **selección de dominios** (qué áreas tenés activas).
- **Por dominio** — entrás a un dominio y ves su config, renderizada desde su contrato (`manifest.json`): en **development**, _Models per agent · Strict TDD · Auto-mine EDR_; en **design**, _Models per agent · Auto-mine DDR_. La config de un dominio **solo aparece si está activo**.

Los toggles son scope-aware (editan el config global o el del proyecto según el scope activo).

```bash
matecito-ai            # TUI
matecito-ai verify     # estado del entorno
matecito-ai install    # instalar/actualizar lo que falte
```
