<!-- matecito-ai: template canónico del retorno de sdd-init.
     Caso aparte: esta fase NUNCA tuvo un bloque titulado. Su `Output Contract` re-enumeraba los
     campos del envelope por su cuenta —y encima incompleto: se olvidaba de `detailed_report` y de
     `skill_resolution`— que es justo el patrón que estos templates vienen a eliminar. Sin bloque
     titulado no hay nada que el Return Contract Check pueda matchear, así que el retorno de init era
     el único que nadie validaba. Este archivo define ese bloque; la skill lo referencia. -->

# Return template — `sdd-init`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.
- **What to detect and how to persist it**: `~/.claude/skills/sdd-init/references/init-details.md` —
  the detection checklist, the Engram payloads and the full Testing Capabilities format. That file
  owns what goes INTO the artifact; this one owns what comes BACK.

The block of this phase is `## SDD Initialized`. Section D.2 lists it as "the block its skill
defines" — this file is that definition, and the skill points here.

`sdd-init` runs BEFORE any flow phase (the orchestrator's Init Guard dispatches it when
`sdd-init/{project}` is not in Engram), so this return is what the orchestrator has when it decides
whether the project is ready for a flow at all.

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Project` | always | the orchestrator, as context for every later phase |
| `### Persistence` | always | the orchestrator: it caches the artifact-store mode from here |
| `### Testing Capabilities` | always | the orchestrator: it forwards the test command and the UI/debug availability |
| `### Artifacts Saved` | always | the orchestrator: what later phases can retrieve |
| `### Blocker` | only on `status: blocked` | the orchestrator: it puts the question to the user |
| `### Next Step` | always | the orchestrator, to route |

Titles are fixed. A capability that is absent is reported as absent (`❌`, `—`) — that is content,
not a reason to drop a row. Detection never guesses: an undetected runner is `❌`, never a plausible
default.

## `status: done` — the project is initialized

```markdown
## SDD Initialized

**Project**: {project}
**Mode**: engram | none

### Project
- **Stack**: {languages, frameworks, package manager}
- **Architecture**: {layout in one line — monorepo/apps, layers, entry points}
- **Conventions**: {what the project's own files fix — CLAUDE.md, config, lint/format rules}

### Persistence
- **Mode**: engram | none
- **Limitations**: {engram: memory is local and non-shareable across machines. none: nothing is
  persisted — every phase returns inline and the flow cannot resume across sessions; recommend
  enabling persistence.}

### Testing Capabilities
| Capability | Available | Detail |
|------------|-----------|--------|
| Test runner | ✅ / ❌ | `{command}` — {framework} |
| Unit / Integration / E2E | ✅ / ❌ each | {tool or —} |
| Coverage | ✅ / ❌ | `{command or —}` |
| Linter / Type checker / Formatter | ✅ / ❌ each | `{command or —}` |
| UI test (`uiTest.available`) | ✅ / ❌ | proofshot {✅/❌} + devServer {`command` / ❌} |
| Debugger toolchain | ✅ / ❌ | {adapter/toolchain found or —} |

{The full detail is in the persisted artifact — this table is the orchestrator's summary.}

### Artifacts Saved
- Engram `sdd-init/{project}` — {observation ID}
- Engram `sdd/{project}/testing-capabilities` — {observation ID}

{In `none` mode: "None — inline only, nothing persisted."}

### Next Step
Ready for intake — hand the orchestrator the request and it dispatches `sdd-intake` (or `sdd-explore`
first, to investigate).
```

## `status: blocked` — the project cannot be initialized as-is

Detection never guesses (`## Hard Rules`). When what you find does not resolve to one answer — no
recognizable project at this path, several incompatible manifests with no way to tell which one
governs, a test setup that contradicts itself — that is a question for the user, not a default for
you to pick.

Same block, with two differences. Everything you DID detect is still reported.

```markdown
## SDD Initialized

**Project**: {project or "not determinable — see Blocker"}
**Mode**: engram | none

### Project
{What IS settled. Anything undetermined is marked "not determinable — see Blocker", never filled in
with a plausible value.}

### Persistence
{As above — persistence mode usually resolves even when detection does not.}

### Testing Capabilities
{The rows you could establish. The rest: "not determinable — see Blocker".}

### Artifacts Saved
{Whatever you persisted before stopping, or "None — stopped before persisting."}

### Blocker
{The question, in one line, phrased so the user can answer it without reading the rest.}

**What you found**: {the ambiguous or contradictory evidence, concretely — which files, which
conflicting signals.}

**What unblocks it**: {the answer you need.}

### Next Step
None. Initialization cannot complete until the blocker is resolved.
```

`### Blocker` is where the question goes. Not `risks` — Section D.4 forbids routing a decision the
user owns through that field. And a capability you simply could not find is NOT a blocker: it is a
`❌` row.

## Artifact vs return — do not confuse them

The **artifacts** (Engram `sdd-init/{project}` and `sdd/{project}/testing-capabilities`) are the
full detected context and the complete Testing Capabilities document, in the format
`~/.claude/skills/sdd-init/references/init-details.md` fixes — later phases retrieve them directly.
The **return** is this block: the summary the orchestrator needs in order to route and to cache the
mode, the test command and the UI/debug availability. The orchestrator never reads the artifacts.
