<!-- matecito-ai: template canónico del retorno de sdd-archive.
     Existe porque el bloque vivía inline en la skill y sólo contemplaba el caso feliz: las dos
     situaciones en que esta fase TIENE que parar —un verify-report con CRITICAL, y un merge que
     sería destructivo— estaban escritas como reglas ("avisá al orquestador") sin ninguna forma de
     retorno que las transportara. Un ejecutor que paraba tenía que inventar la sección.
     Este archivo es LA fuente del formato. La skill y el agente lo referencian; no lo copian. -->

# Return template — `sdd-archive`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Archive Report (Engram)` | always | the orchestrator: the audit trail |
| `### Source of Truth Updated` | always | the orchestrator, as context |
| `### Blocker` | only on `status: blocked` | the orchestrator: it puts the question to the user |
| `### SDD Cycle Complete` | always | the orchestrator: the closing line it relays |

Titles are fixed. The unconditional sections are emitted even when there is nothing to report — with
a `None…` sentinel that says WHY (`None — `none` mode, no durable store.`), never omitted. In `none`
mode nothing is persisted and no store is merged, and that is exactly what the sentinels must state.

## `status: done` — the change is closed

```markdown
## Change Archived

**Change**: {change-name}
**Archived to**: Engram archive report (engram) | inline (none)

### Archive Report (Engram)
- proposal, spec, design, tasks, verify-report observation IDs recorded

{In `none` mode: "None — inline only, nothing persisted."}

### Source of Truth Updated
{One line stating what closing this change leaves behind — the archive report is now the record of
this change; the durable capability-specs, if this change touched any, were already folded by
`sdd-apply`'s Step 5b before verification ran, so this phase reports on them, it does not update them.}

### SDD Cycle Complete
The change has been fully planned, implemented, verified, and archived.
Ready for the next change.
```

Never list EDRs here — of any status, `Inferred` included. They live only as `.md` under
`.matecito-ai/edr/`, never in the archive report and never in Engram.

## `status: blocked` — the change must not be closed yet

One cause, a stop this phase owns and not yours to resolve:

1. **The verification report carries CRITICAL issues.** Archiving a change with CRITICALs is
   forbidden. Do not persist the archive report.

Same block, with the differences below. Everything you did complete is still reported — and, per
Section D.4, anything already persisted is listed in `artifacts` whatever the status.

```markdown
## Change Archived

**Change**: {change-name}
**Archived to**: not archived — see Blocker.

### Archive Report (Engram)
Not persisted — the report is written only once the change is actually archived.

### Source of Truth Updated
Not updated — this phase never writes to the durable capability-specs; the change's delta, if it
touched any, was already folded by `sdd-apply` before verification ran.

### Blocker
{The question, in one line, phrased so the user can answer it without reading the rest.}

**Cause**: CRITICAL issues in the verification report

**What is at stake**: {list the CRITICALs verbatim from the verify-report.}

**What unblocks it**: {go back to `sdd-apply` to fix the CRITICALs, then re-dispatch `sdd-verify`.}

### SDD Cycle Complete
Not complete. The change stays open until the blocker is resolved.
```

The title stays `## Change Archived` even here: Section D.2 fixes one block per phase, and the
orchestrator matches it literally. The BODY says the truth — `**Archived to**: not archived` and a
`### SDD Cycle Complete` that states the cycle did not close.

`### Blocker` is where the question goes. Not `risks` — Section D.4 forbids routing a decision the
user owns through that field.

## Artifact vs return — do not confuse them

One thing is persisted by this phase, and it is not this block. The **archive report** (Engram
`sdd/{change-name}/archive-report`) is the audit trail with every observation ID. This phase writes no
file: the durable capability-specs, when a change touches any, are `sdd-apply`'s output (Step 5b), not
this phase's. The **return** is this block: what the orchestrator needs in order to close the cycle or
to stop it. The orchestrator never reads the archive report.
