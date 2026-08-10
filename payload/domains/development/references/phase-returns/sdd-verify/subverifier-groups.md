<!-- matecito-ai: sole home of the sdd-verify fan-out partition — which group owns which check, what a
     sub-verifier returns, and how the orchestrator merges N sub-reports into the one
     `## Verification Report` block. No `.yaml` pair: nothing renders or validates this file
     mechanically, unlike `sdd-verify.yaml` (which stays the shape of the CONSOLIDATED block and is
     unchanged by this partition). Field names below are not invented — they are what
     `render-return.js --schema --phase sdd-verify` publishes; only the three sub-report envelope
     labels (`Group` / `Result` / `Why incomplete`) are new, and they exist to be read, not rendered. -->

# `sdd-verify` sub-verifier groups

Verification runs as a fan-out: the orchestrator dispatches, **in one message**, one sub-verifier per
group whose gate is active, and consolidates their fragments into a single `## Verification Report`.
This is the **sole** exception to one-agent-per-phase dispatch in this pipeline — it applies to
`sdd-verify` only, and this file is not a template for another phase to copy.

## Groups (the partition)

Seven groups, fixed and named. Every check the phase performs belongs to exactly one — none is lost,
none runs twice. A group's **gate** decides whether it is dispatched at all; the orchestrator resolves
every gate (see "Who resolves what" below) — no sub-verifier re-evaluates its own.

| Group | Gate (orchestrator resolves) | SKILL steps it executes | Data keys it returns |
|---|---|---|---|
| `execution` | always | 3, 4, 5 (runtime half), 5b, 7, 8, Strict-TDD extension | `version`, `completeness`, `build`, `tests`, `coverage`, `spec_compliance`, `compliance_summary` |
| `correctness` | always | 5 (static half) | `correctness` |
| `design-coherence` | always | 6 | `coherence_design` |
| `edr-coherence` | `.matecito-ai/edr/` non-empty | 6b | `coherence_edrs` |
| `spec-coherence` | `.matecito-ai/development-specs/` non-empty | 6d, 6e | `coherence_specs` |
| `ui` | `ui-test: needed` ∧ `uiTest.available` ✅ | 3b-bis–3f | `ui_verdict`, `error_gate` |
| `decision-gaps` | `flagDecisionGaps` resolved `true` | 6c | `decision_gaps` |

The step numbers are the ones `~/.claude/skills/sdd-verify/SKILL.md` already declares — this partition
does not renumber or re-specify any of them. A sub-verifier runs **only** the steps its row lists;
every other step is not its business, and it does not run it, check it, or comment on it.

## Who resolves what

**The orchestrator owns:** `status`, `change`, `mode`, `blocker.*`, the four gate booleans
(`edr_store_active`, `spec_store_active`, `ui_gate_passed`, `flag_decision_gaps`), `verdict`, and the
concatenation of `issues.critical` / `issues.warning` / `issues.suggestion` across every dispatched
group. **Each sub-verifier owns exactly** the data keys its row above lists — nothing else. Every
schema key has exactly one owner, so no check is lost and none runs twice.

`execution` and `correctness` and `design-coherence` are **always** dispatched: their gate is `always`,
so N is never below three. The other four (rows 4-7) turn on and off with the gate their row names —
**that is what "N is stable in definition, variable in execution" means.** A group whose gate is off
MUST NOT be dispatched, and its absence MUST be exactly as silent as the skip its gate already
produced before this partition existed: no placeholder, no mention, no row.

**Invariant.** A group whose gate can turn off (rows 4-7) owns *only* conditional sections of the
consolidated report. Gate off → the orchestrator never dispatches it → the boolean is `false` → the
section is absent, with no mention anywhere. Gate on but the group does not complete → the boolean
stays `true` → the section IS present, filled with the placeholder from "Merge" below, and a CRITICAL
names the incomplete group. These two outcomes must stay legible apart without opening any other
artifact — see the spec's "gate apagado y grupo incompleto no se confunden" scenario.

## Sub-verifier return (`detailed_report`)

A sub-verifier is **not** returning a phase result — `sdd-verify`'s Return Contract Check runs once,
against the orchestrator's consolidated block, never against a sub-report. What a sub-verifier hands
back is this envelope:

~~~markdown
## Sub-Report: <group>
**Group**: <group>
**Result**: complete | incomplete
**Why incomplete**: <one line> | —

```json
{ …exactly the data keys of its row above…,
  "issues": { "critical": [], "warning": [], "suggestion": [] } }
```
~~~

Rules for this envelope:

- **No key outside its row.** The consolidation step (below) ignores unknown keys — a stray one is
  content silently lost, not an error anyone sees.
- **No envelope metadata inside the JSON.** `Group`, `Result` and `Why incomplete` are markdown lines
  above the fence, not JSON fields; the fence carries only the data keys and `issues`.
- `artifacts`: `none` — a sub-verifier never persists anything of its own. There is exactly one
  `verify-report` per run, written once by the orchestrator after the merge.
- A sub-verifier **never re-evaluates its own gate** — being dispatched IS the instruction that the
  gate is already active; it does not re-check `.matecito-ai/edr/` presence, the UI flag, or anything
  else that decided whether it runs at all.
- A sub-verifier **never dispatches** anything. The Executor boundary (`_shared/sdd-phase-common.md`)
  holds exactly as it does for every other phase agent: if it needs evidence another group owns, it
  reports that in `Why incomplete` — it does not launch a sub-agent to get it.

## Merge (mechanical, orchestrator-only)

1. **Union the fragments.** Take the JSON object of every dispatched sub-verifier and merge their data
   keys into one object; concatenate `issues.critical` / `issues.warning` / `issues.suggestion` across
   all of them.
2. **Detect incomplete groups.** A group is incomplete when its returned key set does not match its
   row above, or when it returned `blocked`, or when it returned nothing at all (a sub-verifier that
   never delivers its Sub-Report is treated exactly like `blocked` — never a silent omission).
3. **Placeholder for what an incomplete group owed.** For each data key an incomplete group `g` owned:
   a table becomes one row, every cell `⛔ Not run`, last column `group <g> incomplete — <reason>`; a
   scalar becomes that same string; `compliance_summary` becomes `{compliant: "—", total: "—"}`. Then
   push `<g> did not complete: <reason>` onto `issues.critical`.
4. **Compute the verdict.** Any CRITICAL anywhere in the unioned `issues.critical` → `FAIL`; else any
   WARNING → `PASS WITH WARNINGS`; else `PASS`. The CRITICAL step 3 injects for an incomplete group is
   what forbids `PASS` for that case — there is no second rule for it.
5. **Compute `status`.** `blocked` if and only if `execution` is incomplete (no execution evidence at
   all is exactly the template's `blocked` case); otherwise `done`.
6. **Render, validate, persist once.** Assemble the full data object (the union plus the orchestrator-
   owned fields), run `render-return.js --phase sdd-verify`, then `validate-return.js --phase sdd-verify`,
   then a single `mem_save` of `sdd/{change-name}/verify-report`. No sub-verifier persists anything —
   this is the only write.

## Scope

This partition exists for `sdd-verify` alone. It is not a reusable "fan-out" mechanism offered to any
other phase — every other phase in the pipeline still dispatches a single agent, exactly as before this
file existed. Generalizing it is a later change, not something this file, or any prose near it,
proposes.
