<!-- matecito-ai: template canónico del retorno de sdd-verify.
     Existe porque el formato vivía en tres lugares que no coincidían: el `Output Contract` de la
     skill (una lista en prosa), la plantilla de `references/report-format.md`, y el contrato del
     ejecutor. Dos secciones CONDICIONALES quedaron sin lugar en ninguna de las tres —
     `## Decision Gaps`, que alimenta el mine gate post-verify, y `## UI Verdict`, que bloquea el
     archive — así que un ejecutor podía emitir un retorno "completo" y dejar un gate sin disparar
     nunca. Este archivo es LA fuente del formato: la skill y el agente lo referencian, no lo copian. -->

# Return template — `sdd-verify`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.
- **Meaning of the cells**: `~/.claude/skills/sdd-verify/references/report-format.md` — the compliance
  status vocabulary (`COMPLIANT` / `FAILING` / `UNTESTED` / `PARTIAL`) and what counts as execution
  evidence. That file owns what a cell *means*; this one owns what the report *looks like*. They
  never restate each other.

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Completeness` | always | the orchestrator, to see what apply left open |
| `### Build & Tests Execution` | always | the orchestrator: the runtime evidence behind the verdict |
| `### Spec Compliance Matrix` | always | the orchestrator: scenario-by-scenario contract check |
| `### Correctness (Static Evidence)` | always | the orchestrator, as context |
<!-- matecito-ai: la fila creció de "el diseño quedó stale" a "acá se ve por qué un desvío es CRITICAL o
     WARNING". La declaración de `sdd-apply` no tenía dónde verse, y sin verse la severidad quedaba sin
     respaldo en el reporte. No se inventa sección: entra como columnas de la tabla que ya existía. -->
| `### Coherence (Design)` | always | the orchestrator: a deviation here means the design artifact is stale, and the row shows why it landed as CRITICAL or WARNING |
<!-- matecito-ai: the EDR compliance check is a hard rule of the skill and a mandatory step (6b), and it
     had NO section: its only home was `### Issues Found`, where the finding shows up only if there WAS a
     violation. With nothing found, no trace remained that the check ran at all — "no violation" and "I
     never checked" looked identical. Exactly the hole this template says it came to close for
     `## Decision Gaps` and `## UI Verdict`, and asymmetric with `### Coherence (Design)`, which does
     have its own table for a rule of the same rank. -->
| `### Coherence (EDRs)` | only when the decision store is active | the orchestrator: proof the EDR check ran, and which EDR a CRITICAL cites |
<!-- matecito-ai: same hole, one check over: the capability-spec compliance step (6d) is a hard rule with
     no section either. Fixed together with the EDR one on purpose — fixing this defect one instance at a
     time is a documented failure mode of this ecosystem. -->
| `### Coherence (Capability-Specs)` | only when the capability-spec store is active | the orchestrator: proof the check ran, which spec a CRITICAL cites, and which specs were skipped as drafts |
| *Strict TDD extension* | only in Strict TDD mode | the orchestrator, as evidence |
| `## UI Verdict` | only when the UI gate passed | the orchestrator: a FAIL is CRITICAL and blocks archive |
| `## Decision Gaps` | only when `flagDecisionGaps` is on | the kernel's post-verify **mine gate** |
| `### Issues Found` | always | the orchestrator: what goes back to `sdd-apply` |
| `### Verdict` | always | the orchestrator, to route: archive or back to apply |

Titles are fixed.

**Two sections are `##`, not `###` — that is deliberate, do not "fix" it.** `## Decision Gaps` is
matched at that level by the mine gate (`~/.claude/skills/_shared/sdd-phase-common.md`, Section D.3)
and `## UI Verdict` at that level by the UI step. Levelling them down to match their neighbours is a
gate that never fires.

**The conditional sections are conditional in the strict sense** of the Return Contract Check: when
their condition does not hold they are simply ABSENT, and that is a valid return. They carry no
`None` sentinel — unlike the unconditional sections, whose absence is a broken return. Concretely:
with `flagDecisionGaps` off, never emit `## Decision Gaps`, not even empty; with the UI gate not
passed, never mention the UI step at all; outside Strict TDD, no TDD evidence; with the decision store
absent or empty, no `### Coherence (EDRs)` — naming EDRs at all when the store is inactive contradicts
the activation gate's total-silence rule.

The two `Coherence` sections below `### Coherence (Design)` are conditional on their **store**, not on
their findings. With the store active each is emitted **always**, including the all-clear case: that is
the whole point of them existing — a check whose only trace is `### Issues Found` is indistinguishable
from a check that never ran. A row per applicable record, `None.` when the change touched nothing that
has one. The two stores gate independently: one can be active while the other is not.

The *Strict TDD extension* is not one section but the block of sections defined in
`~/.claude/skills/sdd-verify/strict-tdd-verify.md` (`## Report Template Extension`: `### TDD
Compliance`, `### Test Layer Distribution`, `### Changed File Coverage`, `### Assertion Quality`,
`### Quality Metrics`). Its shape lives there and is not duplicated here; its **position** is fixed
here — after the last `Coherence` section that exists: `### Coherence (Capability-Specs)`, falling back
to `### Coherence (EDRs)` and then to `### Coherence (Design)` as their stores turn out to be inactive.

## `status: done` — verification ran and produced a verdict

A `FAIL` verdict is still `status: done`: the phase did its job, the change did not pass. Do not
return `blocked` or `partial` for a failing verdict — `blocked` is for a verification that could not
be performed at all (below).

~~~markdown
## Verification Report

**Change**: {change-name}
**Version**: {spec version or N/A}
**Mode**: {Strict TDD | Standard}

### Completeness
| Metric | Value |
|--------|-------|
| Tasks total | {N} |
| Tasks complete | {N} |
| Tasks incomplete | {N} |

### Build & Tests Execution
**Build**: ✅ Passed / ❌ Failed
```text
{build command and relevant output}
```

**Tests**: ✅ {N} passed / ❌ {N} failed / ⚠️ {N} skipped
```text
{test command and failure details}
```

**Coverage**: {N}% / threshold: {N}% → ✅ Above / ⚠️ Below / ➖ Not available

### Spec Compliance Matrix
| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| {REQ-01} | {Scenario} | `{file} > {test}` | ✅ COMPLIANT |
| {REQ-02} | {Scenario} | (none found) | ❌ UNTESTED |

**Compliance summary**: {N}/{total} scenarios compliant

### Correctness (Static Evidence)
| Requirement | Status | Notes |
|------------|--------|-------|
| {Req name} | ✅ Implemented | {brief note} |

### Coherence (Design)
| Decision | Followed? | Verifiable per apply | Severity | Notes |
|----------|-----------|----------------------|----------|-------|
| {Decision} | ✅ Yes | — | — | |
| {Decision} | ❌ Deviated | ✅ Declared verifiable / ❌ Declared not verifiable / ➖ Not declared | CRITICAL / WARNING | {the deviation, and the apply declaration it came from} |

### Coherence (EDRs)

{CONDITIONAL — emit ONLY when the decision store is active (`.matecito-ai/edr/` exists and has
content). Absent or empty store → omit the whole section, and do not mention EDRs anywhere in the
report. With the store active, emit it ALWAYS, including the all-clear case: a check whose only trace
is `### Issues Found` is indistinguishable from a check that never ran. One row per EDR applicable to
this change — those the design's `EDR Alignment` lists, or the EDRs of the domains this change touched.
`None.` when the change touched no domain that has one.}

| EDR | Followed? | Severity | Notes |
|-----|-----------|----------|-------|
| `<domain>/<slug>.md` | ✅ Yes | — | |
| `<domain>/<slug>.md` | ❌ Violated | CRITICAL | {the concrete rule the code does not honor} |

### Coherence (Capability-Specs)

{CONDITIONAL — emit ONLY when the capability-spec store is active (`.matecito-ai/development-specs/`
exists and has content). Absent or empty store → omit the whole section. With the store active, emit it
ALWAYS, all-clear case included. One row per capability this change touches that has a durable spec.
`None.` when it touches none. The check is scoped to `Status: Accepted`; an `Inferred` or `Draft` spec
still gets its row, marked skipped — a spec that exists but is not yet a contract is information, not an
absence.}

| Capability-spec | Status | Followed? | Severity | Notes |
|-----------------|--------|-----------|----------|-------|
| `<type>/<capability>.md` | Accepted | ✅ Yes | — | |
| `<type>/<capability>.md` | Accepted | ❌ Diverged | CRITICAL | {the accumulated behavior the code does not honor} |
| `<type>/<capability>.md` | Inferred | ➖ Not checked | — | non-ratified draft — never a contract |

{Strict TDD only: insert here the sections of `## Report Template Extension` from
`~/.claude/skills/sdd-verify/strict-tdd-verify.md`. Outside Strict TDD, nothing goes here.}

## UI Verdict

{CONDITIONAL — emit ONLY when the UI gate passed: `ui-test == needed` AND `uiTest.available = ✅`.
Otherwise omit the whole section, silently.}

<!-- matecito-ai: the condition used to include "AND static validation passed", which made the section
     disappear exactly when something was wrong with the inputs — a failed coverage or static check had
     nowhere to be reported and the silence looked identical to "no UI work was needed". The gate (flag +
     availability) decides whether the section exists; everything that went wrong AFTER the gate is
     reported INSIDE it. Hence the Counterpart column: `UNTESTED` is now expressible. -->
| Scenario | Counterpart | Covers | STATE | Failure Reason |
|----------|-------------|--------|-------|----------------|
| {behavioral scenario `name`} | ✅ Found | ✅ full | ✅ PASS / ❌ FAIL | {reason or —} |
| {behavioral scenario `name`} | ✅ Found | ⚠️ partial | ✅ PASS | {what the counterpart leaves unverified — WARNING even on PASS} |
| {behavioral scenario `name`} | ❌ Missing | — | ❌ UNTESTED | no counterpart in apply-progress — CRITICAL |

{One row per **behavioral** scenario declared by the spec — never one per counterpart, so a scenario
apply never implemented still gets its row. Orphan counterparts (a `name` matching no behavioral
scenario) are listed under WARNING in `### Issues Found`, not here.

When the whole `### UI Scenario Counterparts` section is absent while the spec declares scenarios:
every row is `❌ Missing` / `❌ UNTESTED`, and one CRITICAL states the UI check could not run at all.}

**Error gate**: consoleErrorCount {N} / serverErrorCount {N} → ✅ PASS / ❌ FAIL
**Artifacts**: `proofshot-artifacts/{outputDir}/`

## Decision Gaps

{CONDITIONAL — emit ONLY when `flagDecisionGaps` is on. Its presence does NOT depend on the EDR
store existing: with no store, every decision-touching task is a gap. One row per `· edr:
<domain>/<slug>` ref in the tasks artifact whose file does not exist. With the flag off, omit the
whole section — never emit it empty.}

| domain/slug | task | implemented? |
|-------------|------|--------------|
| {domain}/{slug} | {task title} | yes / no |

### Issues Found
**CRITICAL**: {list or None}
**WARNING**: {list or None}
**SUGGESTION**: {list or None}

### Verdict
{PASS / PASS WITH WARNINGS / FAIL}
{one-line reason}
~~~

The verdict must already account for everything above it, including the conditional sections: a UI
`STATE` FAIL or a failing error gate is CRITICAL and therefore `FAIL`. A `❌ Violated` row in
`### Coherence (EDRs)` is CRITICAL `EDR-VIOLATION` and therefore `FAIL`, and a `❌ Diverged` row in
`### Coherence (Capability-Specs)` is CRITICAL `SPEC-VIOLATION` and therefore `FAIL` — those sections
are evidence, not softer verdicts of their own. A `➖ Not checked` row influences nothing: skipping a
draft is the rule, not a finding. A design deviation marked CRITICAL in `### Coherence (Design)` is
`FAIL` too — only a deviation `sdd-apply` declared not
verifiable against the design is a WARNING, and it caps the verdict at `PASS WITH WARNINGS`.
`## Decision Gaps` is the one section that does **not** influence the verdict — it is input for a
gate that runs after this phase, never a finding.

## `status: blocked` — verification could not be performed

Not "the change failed" — that is `done` + `FAIL`. This is the case where you cannot judge at all:
a required upstream artifact is missing, or no execution evidence can be obtained (no runner, the
environment cannot build) so every scenario would be `UNTESTED` for a reason that is not the
change's fault.

Same block, with these differences. Everything you DID manage to check is still emitted.

~~~markdown
## Verification Report

**Change**: {change-name}
**Version**: {spec version or N/A}
**Mode**: {Strict TDD | Standard}

### Blocker
{What stopped the verification, in one line, phrased so the user can act on it without reading the
rest: the artifact that is missing, or the evidence that cannot be produced and why.}

**What unblocks it**: {re-running an upstream phase, a missing test command, an environment fix —
whichever it is, name it concretely.}

{Then every unconditional section as above, filled with what you could establish and with the rest
marked as not determinable — never invented. The conditional sections follow their own conditions,
unchanged.}

### Verdict
FAIL — not verifiable. {one-line reason}
~~~

`### Blocker` is where this goes — not `risks` (Section D.4 forbids routing a decision or a stop the
user owns through that field) and not `### Issues Found`, which is for findings about the change.

## Artifact vs return — do not confuse them

For this phase the persisted **artifact** (`sdd/{change-name}/verify-report`) and the **return** are
the same document: the block above IS the verify-report. That is why this template, not the skill,
is where its shape is fixed.

The consequence still bites the same way as elsewhere: a section that exists only in the persisted
copy is invisible to every gate. `## Decision Gaps` was exactly that defect — the skill said "add it
to the verify-report", the return contract said it travels inside `detailed_report`, and nothing
guaranteed both were the same text. Emit the block once, complete, and persist that.
