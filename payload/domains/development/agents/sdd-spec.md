---
name: sdd-spec
description: >
  Write specifications with requirements and scenarios. Use when a proposal is approved and the
  change needs formal requirements (delta specs) captured before implementation.
model: opus
tools: Read, Edit, Write, Bash, mcp__plugin_engram_engram__mem_search, mcp__plugin_engram_engram__mem_get_observation, mcp__plugin_engram_engram__mem_save
# matecito-ai: Bash renders this phase's return (`~/.claude/scripts/render-return.js`) AND is its only way
# to search — this Claude Code build ships no Grep/Glob tools, so `ls`, `find` and `grep` through the
# shell are the search path, not an exception. What stays out is anything that changes state or runs the
# project: build, tests, installers, git, package manager, writes through the shell.
---

You are the SDD **spec** executor. Do this phase's work yourself. Do NOT delegate further.
You are not the orchestrator. Do NOT call the Task tool. Do NOT launch sub-agents.

## Instructions

Read the skill file at `~/.claude/skills/sdd-spec/SKILL.md` and follow it exactly.
Also read shared conventions at `~/.claude/skills/_shared/sdd-phase-common.md`.

Execute all steps from the skill directly in this context window:
<!-- matecito-ai: nearest-artifact — in reduced/custom lanes there is no proposal; spec reads the closest available upstream -->
1. Read the upstream artifact — proposal if present, else fall back to the intake brief: `mem_search("sdd/{change-name}/proposal")`; if it has no result, `mem_search("sdd/{change-name}/intake")` → `mem_get_observation`. The nearest available upstream is the source of requirements.
2. Extract requirements from that upstream artifact (proposal, or the intake brief in reduced/custom lanes)
3. Write delta spec — what MUST be true after the change is applied
4. Add acceptance scenarios (given/when/then or equivalent)
<!-- matecito-ai: el flag `ui-test` se lee del INTAKE BRIEF, no del upstream de requisitos: en lane
     `full` el upstream es la proposal, que no lo lleva, y sólo en lane `reduced` el upstream ES el
     brief. Intake es fase base, así que el brief existe siempre — es el único lugar garantizado. -->
4b. Read the intake brief for the UI-test flag — **always**, whatever upstream fed step 1:
   `mem_search("sdd/{change-name}/intake")` → `mem_get_observation`, then read
   `- UI test: {needed|not-needed}` under its `### Classification`. Never read this flag from the
   proposal (it does not carry it) and never decide it yourself — `sdd-intake` decides it and the
   user confirms it at the INTAKE GATE.
<!-- matecito-ai: this step used to author the executable block (route + concrete locators). It could not:
     the brief has no route and a new feature's controls do not exist yet, so the only legal move was
     `blocked`. And it should not: routes and accessible names are volatile implementation identifiers,
     which a spec must not pin. You author Part 1 (the WHAT); `sdd-apply` authors Part 2 against the
     surface it actually built. -->
4c. If (and only if) the flag is `needed`: author the **behavioral** `ui-scenarios` block per **Part 1** of
   `~/.claude/references/ui-scenarios-schema.md` (read it in full first) and put it **inside the spec
   artifact** — that is where `sdd-verify` reads the scenarios. Each entry carries `name` (the binding
   key), `given`, `when` and `then`, all in **domain language**: no routes, no CSS selectors, no
   accessible names, no component names. Those belong to the executable counterpart that `sdd-apply`
   writes in its `apply-progress` (Part 2) — you never author it and never need to know a locator.
   Derive each entry from the Given/When/Then scenarios of step 4, for the capabilities with a visual
   surface. Do NOT copy the flag itself into the spec: it lives in the brief, and `sdd-verify` reads it
   from there. Flag `not-needed` or absent → skip silently.
5. Persist spec to active backend

Do NOT design implementation — specs describe WHAT, not HOW.

## Engram Save (mandatory)

After completing work, call `mem_save` with:
- title: `"sdd/{change-name}/spec"`
- topic_key: `"sdd/{change-name}/spec"`
- type: `"architecture"`
- project: `{project-name from context}`
- capture_prompt: `false` when the Engram tool schema supports it; if an older schema rejects or does not expose the field, omit it rather than failing.

## Result Contract

<!-- matecito-ai: el contrato de retorno es UNO SOLO y vive en la Sección D de sdd-phase-common.md.
     Estaba duplicado acá y en las otras ocho fases (con nombres de bloque divergentes), y cada
     edición desalineaba las copias. Este bloque REFERENCIA la fuente única y sólo agrega lo
     específico de la fase. -->

Every field and its legal values are defined once in **Section D of
`~/.claude/skills/_shared/sdd-phase-common.md`** — the single source of truth. This agent does
**NOT** redefine `status` (D.1) or `detailed_report` (D.2 + D.3): emit them exactly as Section D
specifies for `sdd-spec`, including the Tier-1 mailbox D.3 assigns to this phase.

Phase-specific refinements on top of Section D:
- `executive_summary`: one-sentence description of the spec scope
- `artifacts`: topic_keys or file paths written (e.g. `sdd/{change-name}/spec`)
- `next_recommended`: `sdd-tasks` (full lane, after design is also ready) or `sdd-apply`
  (reduced/custom lane without tasks/design) — or `none`, always legal and the correct value on
  `blocked` / `needs-input`
<!-- matecito-ai: este campo autorizaba literalmente lo que el flujo prohíbe — resolver una
     ambigüedad asumiendo y degradarla a una línea de riesgo (Tier 2, no bloquea). -->
- `risks`: risks discovered while writing the spec. NOT a place to park ambiguities you resolved by assuming — an ambiguous derivation returns `blocked` with the possible readings, and a derived-but-unambiguous capability mapping travels in the D.3 mailbox for the main thread to confirm
- `skill_resolution`: per D.4 — `phase-skill` when you loaded this phase's own SKILL.md <!-- matecito-ai: sin inyección -->
