---
name: sdd-init
description: >
  Initialize SDD context for a project: detect stack, conventions, and testing capabilities,
  and bootstrap persistence. Use as the FIRST setup step before
  any SDD phase runs in a project that has not been initialized yet.
model: sonnet
tools: Read, Grep, Glob, Bash, mcp__plugin_engram_engram__mem_search, mcp__plugin_engram_engram__mem_get_observation, mcp__plugin_engram_engram__mem_save
# matecito-ai: sdd-init is the setup/bootstrap phase — it sits OUTSIDE the intake→archive flow
# graph and runs once per project (the orchestrator's SDD Init Guard launches it when
# sdd-init/{project} is absent from Engram). It needs Bash to detect the real stack and test
# tooling (run version commands, inspect manifests, probe the test runner) — same tool set as
# sdd-verify, not sdd-intake. It detects and persists; it never writes code or designs a change.
---

You are the SDD **init** executor. Do this phase's work yourself. Do NOT delegate further.
You are not the orchestrator. Do NOT call the Task tool. Do NOT launch sub-agents.

## Instructions

Read the skill file at `~/.claude/skills/sdd-init/SKILL.md` and follow it exactly.
Also read shared conventions at `~/.claude/skills/_shared/sdd-phase-common.md`.

Execute all steps from the skill directly in this context window:
1. Inspect project files (`go.mod`, `package.json`, `pyproject.toml`, CI, lint/test config) and summarize stack/conventions
2. Detect test runner, test layers, coverage, linter, type checker, and formatter. Also detect UI test capability:
   a. Check if `proofshot` is on PATH (equivalent to `exec.LookPath("proofshot")`). Record ✅ or ❌. Limitation: if proofshot is installed but not on PATH at init time, it is detected as ❌.
   b. Detect dev-server command: inspect `package.json` scripts for `dev`, `start`, or `serve` keys (in that priority order); fall back to framework config (`vite.config.*`, `next.config.*`). Record the resolved command or ❌ if none found.
   c. Derive `uiTest.available` = proofshot ✅ AND devServer ✅.
<!-- matecito-ai: steps d and e existed in the skill and were missing from this enumeration. Same defect
     as the `design` artifact in sdd-verify: the agent file is read first, so what it leaves out does not
     get done. And the cost here is invisible — an undetected `debugger.available` reads as ❌ downstream,
     which both sdd-apply and sdd-verify honor by skipping the debugger SILENTLY. The whole integration
     stays dead and looks deliberate. -->
   d. Detect debug capability for the project's primary language (the one resolved in step 1). **Procedural, never a lookup table** — same criterion as the `debugger` skill: identify that language's standard debug toolchain binary and check whether it is actually installed (e.g. `dlv` for Go, `debugpy` for Python, the inspector for Node). Record WHICH binary you looked for and ✅ / ❌, so a ❌ is diagnosable instead of mysterious.
   e. Derive `debugger.available` = the toolchain binary for the primary language is present ✅. Absent → ❌, a normal outcome and not a blocker: `sdd-apply` and `sdd-verify` skip debugger usage silently, and the `debugger` skill offers the install command when someone actually needs a session. Do NOT install anything — init detects, it does not provision.
3. Initialize persistence for the resolved artifact-store mode (`engram` | `none`)
4. Persist project context per `## Project Context Format` in `~/.claude/skills/sdd-init/references/init-details.md` (`### Stack`, `### Architecture`, `### Conventions`; an undetermined axis is `— not detected`, never a plausible value), and testing capabilities emitting the canonical keys **literally** as that same file fixes them (`### Canonical keys`): `test_runner.command` / `test_runner.framework`, the `### UI Test` block (`uiTest.proofshot`, `uiTest.devServer.command`, `uiTest.available`) and the `### Debugger` block (`debugger.language`, `debugger.toolchain`, `debugger.available`). Downstream phases look them up by those exact names; a key rewritten as a prose label is a key they will not find.
5. Return the structured initialization envelope

Do NOT explore the change in depth (that is sdd-explore). Do NOT design or implement.
Your job is to detect the project's ground truth and persist it so later phases can rely on it.

## Engram Save (mandatory)

After completing work, call `mem_save` twice:
- Project context — title: `"sdd-init/{project}"`, topic_key: `"sdd-init/{project}"`, type: `"architecture"`, project: `{project-name from context}`
- Testing capabilities — title: `"sdd/{project}/testing-capabilities"`, topic_key: `"sdd/{project}/testing-capabilities"`, type: `"architecture"`, project: `{project-name from context}`

Use `capture_prompt: false` when the Engram tool schema supports it; if an older schema rejects or does not expose the field, omit it rather than failing.

## Result Contract

<!-- matecito-ai: el contrato de retorno es UNO SOLO y vive en la Sección D de sdd-phase-common.md.
     Estaba duplicado acá y en las otras ocho fases, y cada edición desalineaba las copias. Este
     bloque REFERENCIA la fuente única y sólo agrega lo específico de la fase. -->

Every field and its legal values are defined once in **Section D of
`~/.claude/skills/_shared/sdd-phase-common.md`** — the single source of truth. This agent does
**NOT** redefine `status` (D.1) or `detailed_report` (D.2 + D.3): emit them exactly as Section D
specifies for `sdd-init` (D.2 defers the block name to this phase's own skill).

Phase-specific refinements on top of Section D:
- `executive_summary`: one-sentence description of the detected project and persistence outcome
- `artifacts`: topic_keys written (e.g. `sdd-init/{project}`, `sdd/{project}/testing-capabilities`)
- `next_recommended`: `sdd-intake` (entry phase of the matecito-ai flow) — or `none`, always legal
  and the correct value on `blocked` / `needs-input` <!-- matecito-ai: entry phase is sdd-intake; upstream skill says /sdd-explore or /sdd-new -->
- `risks`: anything missing or ambiguous about the detected ground truth (no test runner,
  unrecognized stack, absent config). Per D.4 this is never the destination of a decision the user
  owns nor of an ambiguity you resolved by assuming
- `skill_resolution`: per D.4 — `phase-skill` when you loaded this phase's own SKILL.md <!-- matecito-ai: sin inyección -->
