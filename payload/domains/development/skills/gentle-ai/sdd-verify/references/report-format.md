<!-- matecito-ai: este archivo llevaba la plantilla completa del reporte, y era una SEGUNDA descripción
     del mismo bloque que el retorno de la fase — con dos huecos que la otra descripción sí tenía
     (`## Decision Gaps`, `## UI Verdict`). Dos plantillas de la misma cosa se desincronizan siempre.
     La FORMA se fue al template de retorno, que es lo que el orquestador valida; acá queda lo que
     este archivo sí es dueño y el template no repite: qué SIGNIFICA cada celda. -->

# SDD Verify Report Format

## Where the shape lives

The report's shape — its sections, their titles and order, which ones are conditional, and what
changes per status — is defined once, in `~/.claude/references/phase-returns/sdd-verify/sdd-verify.md`. That is
the file the orchestrator validates your return against, so it is the one you follow when writing the
report. This file is not a second template and does not restate it.

That template is also where the two **conditional** sections have their place: `## UI Verdict` (only
when the UI gate passed) and `## Decision Gaps` (only when `flagDecisionGaps` is on). Neither takes a
`None` sentinel — when its condition does not hold, the section simply does not exist.

## Compliance Statuses

The vocabulary of the `Result` column in the Spec Compliance Matrix:

- ✅ `COMPLIANT`: covering test exists and passed.
- ❌ `FAILING`: covering test exists but failed.
- ❌ `UNTESTED`: no covering test found.
- ⚠️ `PARTIAL`: test passes but covers only part of the scenario.

A scenario is `COMPLIANT` only when a covering test **passed at runtime**. Static inspection alone
never yields `COMPLIANT` — that evidence belongs in the Correctness (Static Evidence) table, which is
not a compliance verdict.

<!-- matecito-ai: acá vive el SIGNIFICADO de las celdas, y las de coherencia de diseño no lo tenían.
     Sin esto, la columna `Verifiable per apply` del template es un dato suelto: se ve, pero nadie dice
     de dónde sale ni qué severidad implica. -->
## Design Coherence

The vocabulary of the `### Coherence (Design)` table:

- `Followed?` — ✅ `Yes` when the changed code honors the design decision; ❌ `Deviated` when it does
  not. Deviation is a fact about the code, established here; it is not taken on faith from apply.
- `Verifiable per apply` — what `sdd-apply` declared about that deviation in its `### Deviations from
  Design` block: whether **this** phase would check it against the design, i.e. whether it touches a
  spec requirement or a task `criteria:`. ✅ `Declared verifiable` / ❌ `Declared not verifiable` /
  ➖ `Not declared`. Only apply has the context to make that call; you read it, you do not re-derive it.
- `Severity` — CRITICAL when the deviation is declared verifiable **or** was not declared at all;
  WARNING only when apply explicitly declared it not verifiable.

Why CRITICAL is the default: `sdd-apply` must return `blocked` when a design gap affects what it is
about to write. So a deviation that survives to verification is one of two things — execution detail
(not this phase's business, and apply is the one who says so), or a decision taken without asking,
which is a hard-stop violation. Absent a declaration, assume the second: the strict reading is the
safe one, and only CRITICAL blocks archive.

A deviation that also breaks a spec needs no special case here — the Spec Compliance Matrix catches it
on its own as CRITICAL. That is why this table carries no "unless it breaks a spec" escape: breaking a
spec makes a deviation worse, never milder.

<!-- matecito-ai: the EDR check is a hard rule and a mandatory step, and its only home in the report was
     `### Issues Found` — where it appears only if there WAS a violation, so the all-clear case left no
     trace and "no violation" read the same as "never checked". It now has `### Coherence (EDRs)`, and
     the meaning of its cells belongs here, next to the design-coherence vocabulary, not in the template. -->
## EDR Coherence

The vocabulary of the `### Coherence (EDRs)` table:

- `EDR` — the record's path relative to `.matecito-ai/edr/`, e.g. `security/auth.md`. Cite the file, not
  a paraphrase of its title: the reader has to be able to open it.
- `Followed?` — ✅ `Yes` when the changed code honors that EDR's concrete rules; ❌ `Violated` when it
  does not. This is a fact you established by inspecting the shipped code, scoped to **this change** —
  not a re-audit of the whole catalog, and not taken on faith from the design's `EDR Alignment`.
- `Severity` — CRITICAL for every ❌ `Violated` row, with no gradations: an `EDR-VIOLATION` blocks
  archive. A ✅ row carries `—`.
- `Notes` — for a violation, the concrete rule the code does not honor (from the EDR's `Reglas
  verificables`), so the fix is actionable without re-reading the record. Empty on a ✅ row.

Which EDRs get a row: those the design's `EDR Alignment` section lists, or — when there is no design
artifact — the EDRs of the domains this change touched. `None.` when the change touched no domain that
has one. The section is emitted whenever the decision store is active, **including the all-clear case**:
its reason for existing is that "checked, nothing found" and "not checked" must not look the same.

An `Inferred` EDR is checked like any other: unlike capability-specs, the EDR compliance step does not
filter by `Status`.

## Capability-Spec Coherence

The vocabulary of the `### Coherence (Capability-Specs)` table. Same reason for existing as the EDR one
— the check is a hard rule, and without a table the all-clear case leaves no trace:

- `Capability-spec` — the path relative to `.matecito-ai/development-specs/`, e.g. `flow/export-report.md`.
- `Status` — the spec's own `Status`, read at the load point. It is what decides whether the row is a
  verdict or a skip, so it is a column and not a footnote.
- `Followed?` — ✅ `Yes` when the code honors the **accumulated** intended behavior that spec describes,
  not merely the delta this change introduced; ❌ `Diverged` when it does not; ➖ `Not checked` for a
  spec whose `Status` is `Inferred` or `Draft`.
- `Severity` — CRITICAL for every ❌ `Diverged` row (`SPEC-VIOLATION`, blocks archive). `—` for ✅ and
  for ➖.
- `Notes` — for a divergence, the accumulated behavior the code does not honor. For a ➖ row, why it was
  skipped (non-ratified draft).

**The asymmetry with EDRs is deliberate and load-bearing.** An `Inferred` EDR is enforced; an `Inferred`
capability-spec is not — it is a pending-ratification draft, never a contract, and can never produce a
`SPEC-VIOLATION` until a human promotes it to `Accepted`. That guardrail is what makes it safe to keep
as-built-mined specs in the store. Its cost is that a reader cannot tell "checked and clean" from
"skipped as a draft" unless the report says which — hence the `Status` column and the ➖ row.

## Command Evidence

The Build & Tests Execution section reports commands you actually ran, in this run:

- **Build**: the exact command, plus the output that justifies the ✅/❌ — the failing lines when it
  failed, a short tail when it passed.
- **Tests**: the exact command, the passed/failed/skipped counts as the runner reported them, and the
  failure details verbatim for each failure.
- **Coverage**: the measured percentage against the project's threshold. `➖ Not available` when the
  project has no coverage tooling — never an estimate.

Never report a command you did not run, and never paraphrase output. A command that could not be run
at all is not evidence of anything: say so, and see the `blocked` case in the return template.

## Strict TDD

When Strict TDD is active, the report also carries the sections of `## Report Template Extension`
from `../strict-tdd-verify.md` (TDD compliance, test layer distribution, changed-file coverage,
assertion quality, quality metrics). Their shape lives in that module; their position in the report is
fixed by the return template. Outside Strict TDD they are absent — a legitimate conditional absence,
not a broken report.
