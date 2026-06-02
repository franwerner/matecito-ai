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
2. Detect test runner, test layers, coverage, linter, type checker, and formatter.
3. Initialize persistence for the resolved mode.
<!-- matecito-ai: paso de construir el registry removido -->
4. Persist testing capabilities and project context.
5. Return the structured initialization envelope.

## Output Contract

Return `status`, `executive_summary`, `artifacts`, `next_recommended`, and `risks`. Include project, stack, persistence mode, testing capability table, saved observation IDs/paths, and next `/sdd-explore` or `/sdd-new` step. <!-- matecito-ai: registry path removido del contrato -->-

## References

- [references/init-details.md](references/init-details.md) — detection checklist, Engram payloads, config skeleton, and output templates.
- `../_shared/engram-convention.md` — Engram artifact naming.
