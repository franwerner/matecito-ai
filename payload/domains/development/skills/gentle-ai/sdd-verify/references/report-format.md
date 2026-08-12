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
when the UI gate passed) and `## Decision Gaps` (only when this change materialized at least one
decision record — no flag; see `~/.claude/references/decision-capture/in-flow-capture.md`). Neither
takes a `None` sentinel — when its condition does not hold, the section simply does not exist.

## Compliance Statuses

The vocabulary of the `Result` column in the Spec Compliance Matrix:

- ✅ `COMPLIANT`: covering test exists and passed.
- ❌ `FAILING`: covering test exists but failed.
- ❌ `UNTESTED`: no covering test found.
- ⚠️ `PARTIAL`: test passes but covers only part of the scenario.
- ➖ `OUT-OF-SCOPE — {token}`: the scenario carries a `verification: deferred → <change>` or
  `standing → <owner>` token (canonical definition: `~/.claude/references/spec/README.md`). Inline the
  token verbatim in the cell — e.g. `➖ OUT-OF-SCOPE — deferred → auth-redesign` — so the owner of that
  verification stays visible without opening another artifact.

A scenario is `COMPLIANT` only when a covering test **passed at runtime**. Static inspection alone
never yields `COMPLIANT` — that evidence belongs in the Correctness (Static Evidence) table, which is
not a compliance verdict.

<!-- matecito-ai: verification-scope-token. This paragraph owns the denominator rule; the token's own
     form and semantics stay in `~/.claude/references/spec/README.md` — not restated here. -->
`OUT-OF-SCOPE` is never a substitute for a real result: a `deferred`/`standing` scenario whose test ran
and failed is `FAILING`, not `OUT-OF-SCOPE` — the token exempts missing coverage, never a negative
outcome. The **Compliance summary** (`{N}/{total} scenarios compliant`) counts only `in-scope`
scenarios in `{total}`: an `OUT-OF-SCOPE` row neither inflates nor dilutes it, and stays visible in the
matrix with its token. A scenario with no `verification:` line is `in-scope` and classifies with the
other four values, unchanged.

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

<!-- matecito-ai: el resultado estructural (paso 6e) no lleva tabla propia a propósito. Es por store, no
     por capability, así que una fila por spec repetiría el mismo dato; y lo que este reporte necesita
     no es el detalle —el JSON del motor ya lo tiene— sino la constancia de que el chequeo corrió, que
     es justamente lo que faltaba cuando un merge mal formado pasó tres verificaciones seguidas. -->
Below the table, one line records the outcome of the **structural** validation (Execution Step 6e), so
that "clean" and "never ran" stop looking the same: the command that ran, and the pre-existing finding
counts by severity. Findings on a capability-spec **this change touched** never live here — they are
this change's own and go under `### Issues Found` with the mapped severity. This line carries only what
the store already owed before the change, and it never pushes the verdict.

<!-- matecito-ai: verification-scope-token. This paragraph is the 6d-scoped counterpart of the OUT-OF-SCOPE
     paragraph above; the token's own form and semantics stay in `~/.claude/references/spec/README.md`. -->
A scenario inside an `Accepted` spec that carries a `verification: deferred → <change>` or
`standing → <owner>` token (canonical definition: `~/.claude/references/spec/README.md`) does not, by
itself, make the row `❌ Diverged` — its behavior is intentionally absent (`deferred`) or verified by
exercising its owner's flow (`standing`), not by this change. Cite the token verbatim in `Notes`; the
row still reads `✅ Yes` unless some other scenario in the same spec genuinely diverges, or that same
token-carrying scenario's behavior was directly exercised and failed — the token exempts missing
coverage, never a negative result.

**The asymmetry with EDRs is deliberate and load-bearing.** An `Inferred` EDR is enforced; an `Inferred`
capability-spec is not — it is a pending-ratification draft, never a contract, and can never produce a
`SPEC-VIOLATION` until a human promotes it to `Accepted`. That guardrail is what makes it safe to keep
as-built-mined specs in the store. Its cost is that a reader cannot tell "checked and clean" from
"skipped as a draft" unless the report says which — hence the `Status` column and the ➖ row.

<!-- matecito-ai: in-flow decision capture (development-specifics). Full mechanism, the ratification
     gate, and the materialization contract: in-flow-capture.md — this section owns only the cell
     vocabulary of `## Decision Gaps`. -->
## Decision-Gaps Coherence

The vocabulary of the `## Decision Gaps` table. Runs on every `development` change — no flag; the
section itself is emitted only when `apply-progress`'s `### Decisions Materialized` carries at least
one row:

- `Record` — the EDR's path relative to `.matecito-ai/edr/`, e.g. `contracts/some-slug.md`, taken from
  the materialized record's `record` identity, whether or not the file successfully exists.
- `Task` — the implementing task's id, from the same `### Decisions Materialized` row.
- `Structure` — ✅ `OK` when `validate-artifact.js --type edr --file .matecito-ai/edr/<domain>/<slug>.md`
  exits 0; CRITICAL, naming the finding, on any violation. When the materializing row's own `result`
  was `failed: <reason>` (apply could not write the record at all), skip straight to CRITICAL naming
  that reason — there is no file to run the structural check against.
- `Backing` — "code that corresponds" is a **mechanical join, never a search**: look up the implementing
  task's row(s) in `apply-progress`'s `### Files Changed`. ✅ `OK` when that task's changed-file set
  contains at least one path outside `.matecito-ai/edr/`; CRITICAL when every path it touched is under
  `.matecito-ai/edr/` (only the record and the INDEX — no governing code in this change). `—` when
  `Structure` already failed on a `result: failed` row (there is no record to check backing for).

Coherence **between** records — whether two EDRs contradict each other — is explicitly out of scope for
this group; that check belongs to `development-decisions-validate`, run standalone, unrelated to what
this change materialized. Any CRITICAL in `Structure` or `Backing` is `FAIL`, exactly like the other
coherence sections — see the return template's verdict rule.

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
