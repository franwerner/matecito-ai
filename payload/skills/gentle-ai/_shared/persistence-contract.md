# Persistence Contract (shared across all SDD skills)

<!-- matecito-ai: modos openspec e hybrid REMOVIDOS. Este ecosistema persiste solo en Engram (memoria de sesión). Modos válidos: engram | none. -->

## Mode Resolution

The orchestrator passes `artifact_store.mode` with one of: `engram | none`.

Default: if Engram is available → `engram`. Otherwise → `none`.

(matecito-ai is Engram-only. The `openspec` and `hybrid` file-based modes were removed; do NOT create `openspec/` directories or write artifact files to the repo.)

## Mode Roles

- **`engram`**: Working memory between sessions. Upserts overwrite — no iteration history. Local only.
- **`none`**: Ephemeral. Results returned inline; lost when the conversation ends.

### `engram` mode limitation

Engram uses `topic_key`-based upserts. Re-running a phase for the same change **overwrites** the previous version — no revision history is kept (git history of the actual code is your audit trail). The archive phase saves a summary report, not a full artifact folder.

## Behavior Per Mode

| Mode | Read from | Write to | Project files |
|------|-----------|----------|---------------|
| `engram` | Engram | Engram | Never |
| `none` | Orchestrator prompt context | Nowhere | Never |

## State Persistence (Orchestrator)

The orchestrator persists DAG state after each phase transition to enable SDD recovery after compaction.

| Mode | Persist State | Recover State |
|------|--------------|---------------|
| `engram` | `mem_save(topic_key: "sdd/{change-name}/state", capture_prompt: false*)` | `mem_search("sdd/*/state")` → `mem_get_observation(id)` |
| `none` | Not possible — warn user | Not possible |

*For state automated artifacts, set `capture_prompt: false` when the Engram tool schema supports it; if an older schema rejects or does not expose the field, omit it rather than failing.

## Common Rules

- `none` → do NOT create or modify any project files; return results inline only
- `engram` → do NOT write any project files; persist to Engram and return observation IDs
- NEVER create `openspec/` or any artifact files in the repo (file-based modes were removed)
- If unsure which mode to use, default to `none`

## Sub-Agent Context Rules

Sub-agents launch with a fresh context and NO access to the orchestrator's instructions or memory protocol.

Who reads, who writes:
- Non-SDD (general task): orchestrator searches engram, passes summary in prompt; sub-agent saves discoveries via `mem_save`
- SDD (phase with dependencies): sub-agent reads artifacts directly from Engram; sub-agent saves its artifact
- SDD (phase without dependencies, e.g. explore): nobody reads; sub-agent saves its artifact

Why this split:
- Orchestrator reads for non-SDD: it knows what context is relevant; sub-agents doing their own searches waste tokens
- Sub-agents read for SDD: SDD artifacts are large; inlining them in the orchestrator prompt would consume the context window
- Sub-agents always write: they have the complete detail on what happened

## Orchestrator Prompt Instructions for Sub-Agents

Non-SDD:
```
PERSISTENCE (MANDATORY):
If you make important discoveries, decisions, or fix bugs, you MUST save them to engram before returning:
  mem_save(title: "{short description}", type: "{decision|bugfix|discovery|pattern}",
           project: "{project}", content: "{What, Why, Where, Learned}")
Do NOT return without saving what you learned.
```

SDD (with dependencies):
```
Artifact store mode: {engram|none}
Read these artifacts before starting (search returns truncated previews):
  mem_search(query: "sdd/{change-name}/{type}", project: "{project}") → get ID
  mem_get_observation(id: {id}) → full content (REQUIRED)

PERSISTENCE (MANDATORY — do NOT skip):
After completing your work, you MUST call:
  mem_save(
    title: "sdd/{change-name}/{artifact-type}",
    topic_key: "sdd/{change-name}/{artifact-type}",
    type: "architecture",
    project: "{project}",
    capture_prompt: false,
    content: "{your full artifact markdown}"
  )
If you return without calling mem_save, the next phase CANNOT find your artifact and the pipeline BREAKS.
```

SDD (no dependencies):
```
Artifact store mode: {engram|none}

PERSISTENCE (MANDATORY — do NOT skip):
After completing your work, you MUST call:
  mem_save(
    title: "sdd/{change-name}/{artifact-type}",
    topic_key: "sdd/{change-name}/{artifact-type}",
    type: "architecture",
    project: "{project}",
    capture_prompt: false,
    content: "{your full artifact markdown}"
  )
```

For SDD artifacts, `capture_prompt: false` is explicit and mandatory when the Engram tool schema supports it. Engram v1.15.3 defaults `capture_prompt` to true for normal human/proactive saves, but automated pipeline artifacts must not. If an older schema rejects or does not expose `capture_prompt`, omit it rather than failing.

## Skill Loading

<!-- matecito-ai: el registry/inyección fue removido. Las fases cargan su propia skill directamente y leen las convenciones del proyecto desde sus archivos (.claude/adr/, CLAUDE.md, config.yaml). -->

Each SDD phase loads its own `SKILL.md` (plus any module it explicitly references). Project conventions are read directly from the project's files — `.claude/adr/` for architecture decisions, `CLAUDE.md`, and `config.yaml`. There is no skill registry or pre-injected `Project Standards` block in matecito-ai. See Section A of `sdd-phase-common.md`.

## Detail Level

The orchestrator may pass `detail_level`: `concise | standard | deep`. This controls output verbosity but does NOT affect what gets persisted — always persist the full artifact.
