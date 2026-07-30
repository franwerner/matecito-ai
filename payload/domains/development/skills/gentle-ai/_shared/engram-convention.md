# Engram Artifact Convention (reference documentation)

NOTE: Critical engram calls (`mem_search`, `mem_save`, `mem_get_observation`) are inlined directly in each skill's SKILL.md. This document is supplementary reference — sub-agents do NOT need to read it to function.

## Naming Rules

ALL SDD artifacts persisted to Engram MUST follow this deterministic naming:

```
title:     sdd/{change-name}/{artifact-type}
topic_key: sdd/{change-name}/{artifact-type}
type:      architecture
project:   {detected or current project name}
scope:     project
capture_prompt: false
```

Set `capture_prompt: false` when the Engram tool schema supports it; if an older schema rejects or does not expose the field, omit it rather than failing.

### Artifact Types

| Artifact Type | Produced By | Description |
|---------------|-------------|-------------|
<!-- matecito-ai: `intake` was missing despite being the ENTRY phase, and despite its brief being read by
     sdd-spec (the `ui-test` flag), sdd-design (the `diagram` flag) and sdd-verify. `testing-capabilities`
     was missing too, and sdd-apply and sdd-verify consult it. A "canonical" table with holes is not
     canonical. -->
| `intake` | sdd-intake | Intake Brief — entry artifact of the flow (Pass 2 only; Pass 1 persists nothing) |
| `explore` | sdd-explore | Exploration analysis |
| `proposal` | sdd-propose | Change proposal |
| `spec` | sdd-spec | Delta specifications (all domains concatenated) |
| `design` | sdd-design | Technical design |
| `tasks` | sdd-tasks | Task breakdown |
| `apply-progress` | sdd-apply | Implementation progress (one per batch) |
| `verify-report` | sdd-verify | Verification report |
| `archive-report` | sdd-archive | Archive closure with lineage |
<!-- matecito-ai: the `state` artifact is gone. It had a declared format (here and in
     persistence-contract.md), NO instructed producer — the row said "orchestrator", but nothing in the
     kernel or the domain fragment ever told the orchestrator to write it — and NO consumer at all. The
     inverse of every other defect in this sweep: not a reader without an emitter, a format with neither
     end. Post-compaction recovery already works through the phase artifacts themselves (deterministic
     `topic_key` + the kernel's Recovery Rule), so `state` would have been a SECOND source of truth about
     which phase ran and which tasks are done — duplicated information that drifts, which is the failure
     mode this whole round has been undoing. Implementing it would have added the problem, not fixed it. -->

### Project-scoped artifacts (not under `sdd/{change-name}/`)

Two artifacts belong to the PROJECT, not to a change, so they do not take a `{change-name}` segment.
They are produced once by `sdd-init` and read by later phases of every change.

| Topic key | Produced By | Description |
|-----------|-------------|-------------|
| `sdd-init/{project}` | sdd-init | Detected project context: stack, architecture, conventions |
| `sdd/{project}/testing-capabilities` | sdd-init | Test runner, layers, coverage, quality tools, `uiTest.*`, `debugger.*` — format and canonical keys in `~/.claude/skills/sdd-init/references/init-details.md` |


## Recovery Protocol (2 steps)

```
Step 1: mem_search(query: "sdd/{change-name}/{artifact-type}", project: "{project}") → truncated preview + ID
Step 2: mem_get_observation(id: {observation-id}) → complete content
```

When retrieving multiple artifacts, group all searches first, then all retrievals:

```
STEP A — SEARCH (get IDs only):
  mem_search(query: "sdd/{change-name}/proposal", ...) → save ID
  mem_search(query: "sdd/{change-name}/spec", ...) → save ID
  mem_search(query: "sdd/{change-name}/design", ...) → save ID

STEP B — RETRIEVE FULL CONTENT (mandatory):
  mem_get_observation(id: {proposal_id})
  mem_get_observation(id: {spec_id})
  mem_get_observation(id: {design_id})
```

Loading project context:
```
mem_search(query: "sdd-init/{project}", project: "{project}") → get ID
mem_get_observation(id) → full project context
```

## Writing Artifacts

Standard write:
```
mem_save(
  title: "sdd/{change-name}/{artifact-type}",
  topic_key: "sdd/{change-name}/{artifact-type}",
  type: "architecture",
  project: "{project}",
  capture_prompt: false,
  content: "{full markdown content}"
)
```

Concrete example — saving a proposal for `add-dark-mode`:
```
mem_save(
  title: "sdd/add-dark-mode/proposal",
  topic_key: "sdd/add-dark-mode/proposal",
  type: "architecture",
  project: "my-app",
  capture_prompt: false,
  content: "## Proposal\n\nAdd dark mode toggle..."
)
```

`capture_prompt: false` is REQUIRED for SDD artifacts when the Engram tool schema supports it. Engram v1.15.3 captures user prompts by default for human/proactive saves, but SDD artifacts are automated pipeline outputs. Do not infer this from `type` because both SDD artifacts and human architecture decisions use `architecture`. If an older schema rejects or does not expose `capture_prompt`, omit it rather than failing.

Update existing artifact (when you have the observation ID):
```
mem_update(id: {observation-id}, content: "{updated full content}")
```

Use `mem_update` when you have the exact ID. Use `mem_save` with same `topic_key` for upserts.

### Browsing All Artifacts for a Change

```
mem_search(query: "sdd/{change-name}/", project: "{project}")
→ Returns all artifacts for that change
```

## Project Name Resolution (engram v1.11.0+)

Engram auto-detects the project name from the git remote at MCP startup. The `--project` flag and `ENGRAM_PROJECT` env var can override detection. All project names are normalized to lowercase and trimmed.

If the agent saves a memory under a project name that doesn't match existing observations, engram warns about potential name drift. Use `mem_merge_projects` (MCP tool) or `engram projects consolidate` (CLI) to merge variants.

## Upsert Behavior

Same `topic_key` + `project` + `scope` → UPDATE (overwrite), not INSERT. Previous content is lost — `revision_count` increments but old content is NOT saved. This is by design — engram is working memory, not an audit trail. The git history of your code is the audit trail. <!-- matecito-ai: openspec/hybrid removidos -->

## Why This Convention

- Deterministic titles → recovery works by exact match
- `topic_key` → enables upserts without duplicates
- `sdd/` prefix → namespaces all SDD artifacts
- Two-step recovery → search previews are always truncated; `mem_get_observation` is the only way to get full content
- Lineage → archive-report includes all observation IDs for complete traceability
