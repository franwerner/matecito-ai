---
name: sdd-init
description: "Trigger: sdd init, iniciar sdd. Initialize SDD context, testing capabilities, and persistence."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "3.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are
> the ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to
> the dedicated `sdd-init` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Activation Contract

Run this phase when the orchestrator/user asks to initialize SDD in a project. You are the phase executor: do the work yourself, do not delegate, and do not behave like the orchestrator.

## Hard Rules

- Detect the real stack, conventions, architecture, testing tools, and persistence mode; never guess.
- Always persist testing capabilities as `sdd/{project}/testing-capabilities` in Engram.
<!-- matecito-ai: registry removido — sdd-init ya no construye .atl/skill-registry.md -->
- Use `capture_prompt: false` for automated SDD/config saves when supported; omit it if the tool schema lacks it.

## Decision Gates

| Input | Action |
|---|---|
| `mode=engram` | Save context and capabilities to Engram only. |
| `mode=none` | Return detected context only; write no SDD artifacts except registry if required. |

## Execution Steps

1. Inspect project files (`package.json`, `go.mod`, `pyproject.toml`, CI, lint/test config) and summarize stack/conventions.
2. Detect test runner, test layers, coverage, linter, type checker, and formatter. Also detect UI test capability:
   a. Check if `proofshot` is on PATH (equivalent to `exec.LookPath("proofshot")`). Record ✅ or ❌. Note: proofshot installed but not on PATH at init time → detected as ❌.
   b. Detect dev-server command from `package.json` scripts (`dev`, `start`, `serve` in priority order) or framework config (`vite.config.*`, `next.config.*`). Record resolved command or ❌.
   c. Derive `uiTest.available` = proofshot ✅ AND devServer ✅.
<!-- matecito-ai: el fragmento del dominio afirma que este toolchain lo detecta init y lo cachea acá,
     y `sdd-apply` / `sdd-verify` consultan `debugger.available` para decidir si pueden abrir una
     sesión. Nadie lo detectaba: el campo nunca se escribía, ambas fases lo leían ausente y saltaban
     el debugger EN SILENCIO. Toda la integración estaba muerta y parecía deliberado. -->
   d. Detect debug capability for the project's primary language (the one you resolved in step 1). This is **procedural, not a lookup table** — same criterion as the `debugger` skill: identify that language's standard debug toolchain binary and check whether it is actually installed (e.g. `dlv` for Go, `debugpy` for Python, the inspector for Node). Record the binary you looked for and ✅ / ❌.
   e. Derive `debugger.available` = the toolchain binary for the primary language is present ✅. Absent → ❌, which is a normal outcome, not a blocker: `sdd-apply` and `sdd-verify` simply skip debugger usage, and the `debugger` skill offers the install command when someone actually needs a session. Do NOT install anything here — init detects, it does not provision.
3. Initialize persistence for the resolved mode.
<!-- matecito-ai: paso de construir el registry removido -->
4. Persist project context per `## Project Context Format` in `references/init-details.md` (`### Stack`, `### Architecture`, `### Conventions` — an axis you could not establish is `— not detected`, never a plausible value), and testing capabilities emitting the canonical keys **literally** as that same file fixes them — `test_runner.command` / `test_runner.framework`, the `### UI Test` block (`uiTest.proofshot`, `uiTest.devServer.command`, `uiTest.available`) and the `### Debugger` block (`debugger.language`, `debugger.toolchain`, `debugger.available`). Downstream phases look them up by those exact names; a key rewritten as a prose label is a key they will not find.
5. Return the structured initialization envelope, with the `## SDD Initialized` block from `~/.claude/references/phase-returns/sdd-init/sdd-init.md`.

## Output Contract

<!-- matecito-ai: esta sección re-enumeraba los campos del envelope por su cuenta — y se olvidaba de `detailed_report` y de `skill_resolution`, así que init era la única fase sin bloque titulado y por lo tanto la única cuyo retorno nadie podía validar. Los campos se heredan de la Sección D; la forma del bloque vive en el template. -->
Return the Section D envelope from `~/.claude/skills/_shared/sdd-phase-common.md` — do not re-declare its fields here — carrying the `## SDD Initialized` block **exactly as `~/.claude/references/phase-returns/sdd-init/sdd-init.md` defines it**. Follow it literally: the orchestrator validates your return against that same file, matching section titles literally.

What the return must CARRY (the template fixes how it looks): project and stack, architecture and conventions detected, persistence mode with its limitations, the testing-capability table (runner, layers, coverage, quality tools, `uiTest.available`, debugger toolchain), the saved observation IDs, and the next step (intake, or exploration first when the area is unclear). Detection never guesses: an undetected capability is `❌`, and evidence that does not resolve to one answer is the template's `blocked` block, not a plausible default.

## References

- `~/.claude/references/phase-returns/sdd-init/sdd-init.md` — **the** shape of the block you return.
- [references/init-details.md](references/init-details.md) — detection checklist, Engram payloads, and the full Testing Capabilities format (what goes into the artifact).
- `~/.claude/skills/_shared/engram-convention.md` — Engram artifact naming.
