---
name: design-init
description: "Trigger: design init, iniciar design. Initialize design context, design capabilities, and persistence."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "1.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are
> the ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to
> the dedicated `design-init` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Activation Contract

Run this phase when the orchestrator/user asks to initialize the design flow in a project. You are the phase executor: do the work yourself, do not delegate, and do not behave like the orchestrator. This is the setup/bootstrap phase — it sits OUTSIDE the `intake → archive` pipeline and runs once per project.

## Hard Rules

- Detect the real design context — connected Figma file, connected Canva, an existing brand guide / design system, the surface type, and whether `.matecito-ai/ddr/` exists; never guess.
- Always persist design capabilities as `design-init/{project}` in Engram, including the `designCapabilities` block (the equivalent of development's `uiTest`).
- Use `capture_prompt: false` for automated design/config saves when supported; omit it if the tool schema lacks it.

## Decision Gates

| Input | Action |
|---|---|
| `mode=engram` | Save context and capabilities to Engram only. |
| `mode=none` | Return detected context only; write no design artifacts. |

## Execution Steps

1. Inspect project files (brand guide, `tokens.json`/design-token files, `README`, asset folders, any existing design-system manifest) and summarize the design context and conventions.
2. Detect the design capabilities. Also detect connection capability:
   a. Check if a Figma file is connected (the `figma` MCP is available and a file is reachable). Record ✅ or ❌. Note: figma MCP not registered or no file connected at init time → detected as ❌.
   b. Check if Canva is connected (the `canva` MCP is available). Record ✅ or ❌.
   c. Detect the surface type (`landing | app-ui | brand-system | marketing`) from the request and the inspected files; record it.
   d. Detect whether a prior brand guide / design system exists. Record resolved path or ❌.
   e. Derive `designCapabilities.available` = figmaConnected ✅ (the floor for reading/verifying real visual work).
3. Initialize persistence for the resolved mode.
4. Persist design capabilities and project context, including the `designCapabilities` block (figmaConnected, canvaConnected, surface, brandGuide, available) per the `### Design Capabilities` section in `references/init-details.md`.
5. Detect whether `.matecito-ai/ddr/` exists with content (DDR activation gate). Record it; do NOT bootstrap DDRs here.
6. Return the structured initialization envelope.

## Output Contract

Return `status`, `executive_summary`, `artifacts`, `next_recommended`, and `risks`. Include project, design context, persistence mode, Design Capabilities table, DDR-store presence, saved observation IDs/paths, and next `/design-intake` step.

## References

- [references/init-details.md](references/init-details.md) — detection checklist, Engram payloads, and output templates.
- `../_shared/engram-convention.md` — Engram artifact naming.
