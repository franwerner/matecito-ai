# In-flow decision capture — development

<!-- matecito-ai: single definition of the mechanism this domain uses INSTEAD of the kernel's generic
     post-verify mine gate (the kernel's override clause lets a domain opt out that way — see
     `~/.claude/matecito-ai.md`, `### Decision-Gap Capture (mine gate)`). `design` keeps the kernel's
     mechanism unchanged; `development` declares this one. Every phase and script mentioned here reads
     THIS file for the mechanism's rules — they cite it, they do not restate it. -->

An architecture decision that surfaces while a `development` change is in flight is **proposed** by
`sdd-design` under its `### New Decisions` mailbox and **materialized** as an `Accepted` EDR by
`sdd-apply`, straight from the design artifact, in the same step that implements the code the decision
governs. The mailbox declares `gates: reported`: it never opens a confirmation gate, and there is no
forwarding channel, no ratification step, and no per-change ledger between propose and materialize.
There is no post-verify mining pass for `development` either — the kernel's generic mine gate never
fires here (see the override clause it carries).

With the ratification gate gone, nothing along this path stops to ask whether a given proposal is
actually a decision record rather than an implementation detail. The canonical concept definition in
`~/.claude/references/edr/README.md` — what IS and what is NOT a decision record — is the only filter
left: `sdd-design` applies it when it decides whether an item belongs under `### New Decisions` at
all, and it is what any reader of this mechanism checks a proposal against, not a rejection verdict
this path no longer carries.

## The proposal — one mailbox item, three tokens

A proposal is a single item under `sdd-design`'s `### New Decisions` return mailbox. It travels in two
halves, per the phase-return contract (`items.rationale` + `items.tokens`, rendered by the existing
engine — `render-return.js`/`validate-return.js`, unmodified in shape beyond the free-form-token fix
below):

```yaml
items:
  text: summary
  rationale: rationale
  tokens:
    - { name: blocking-test, field: blocking_test, values: [none, infra, contract, data-model], passing: [none] }
    - { name: record, field: record }          # free-form `<domain>/<slug>`; no `values`
    - { name: record-mode, field: record_mode, values: [create, modify] }   # closed set, no `passing:`
```

- **`summary`** — the point in dispute and the option taken, in one line.
- **`· rationale:`** — the alternatives considered and why this one, still one line (the full
  reasoning; the gate prints `summary`, `rationale` ships in the block for whoever asks).
- **`· blocking-test:`** — unchanged from before this change; see `sdd-design.md`, "The blocking-test
  token".
- **`· record: <domain>/<slug>`** — the identity the EDR occupies once materialized. Free-form: the
  engine accepts any present, non-null value for a token declared without `values` (see "The
  free-form-token fix" below). It is still **required** — an item with the token line missing fails
  `TOKEN-MISSING` at the Return Contract Check, same strict reading as any other omitted token.
- **`· record-mode: create | modify`** — whether the proposal creates a new record at `record:`, or
  edits an existing one in place. Closed value set, but declared with **no `passing:` key**: every
  value in the set is legal (`validate-return.js:240` resolves `passing = t.passing || legal`, so an
  absent `passing:` defaults to the full `values` list), so this is not a verdict with a failing
  outcome — only an absent token fails, `TOKEN-MISSING`, same strict reading as `record`. **It is a
  routing token, read verbatim by `sdd-apply` straight from the design artifact — not a verdict the
  orchestrator classifies.** It follows `record:`'s precedent, not the precedent of a token like
  `blocking-test`, which the orchestrator classifies to decide whether the flow stops; `record-mode`
  decides nothing about the flow — `sdd-apply` reads it to pick which of the two Materialization
  branches below to run, exactly as it already reads `record:` to pick the file. Stated explicitly so
  nobody reaches for the wrong precedent and builds an orchestrator classification table for it.

A proposal is **not** a record. It MUST NOT create, modify, or touch anything under `.matecito-ai/edr/`,
and it MUST NOT be persisted to Engram as a record — it only ever travels inside the phase's own return
mailbox (the gate half) and its artifact's `## New Decisions` prose (the body half, for whoever wants
the fuller writeup). Before `sdd-apply` runs, `.matecito-ai/edr/` is byte-identical to before the
proposal was made.

### The free-form-token fix

`validate-return.js`'s `checkItems` used to compute `legal = t.values || []` and then unconditionally
check `!legal.includes(value)` — so a token declared *without* `values` (meant to be free-form) failed
`TOKEN-ILLEGAL` for **every** value, because an empty `legal` array rejects everything. Reproduced by
execution before this change shipped: a hand-written return carrying `· record: contracts/some-slug`
against a token declared `{ name: record, field: record }` failed with `ERROR TOKEN-ILLEGAL: ...which
is not one of []`. `render-return.js`'s render-side check (`if (t.values && !t.values.includes(value))
fail(...)`) already skipped the legal-check when `values` was absent — only the validate side had the
bug. The fix: a token with no `values` is free-form (any present, non-null value passes; `TOKEN-MISSING`
still applies when the token is absent). A token that DOES declare `values` is byte-for-byte unchanged.
Covered by `payload/domains/development/dev-tests/validate-return-tokens.test.js`.

## Materialization — `sdd-apply` Step 4b, same step as the implementing task

For each entry under `## New Decisions` in the design artifact (`sdd/{change-name}/design`) whose
governing task this is, `sdd-apply` reads it straight from the artifact — the same one it already
reads as an ordinary upstream input — and materializes it in the **same work-unit step** that
implements the code the decision governs, never a separate pass before or after and never gated by a
confirmation step. Under a fan-out, only the single dispatch role the domain fragment names as the
writer materializes; an isolated run carries the unapplied INDEX rows in its Task Run Report, and the
consolidation run applies them once (see "The INDEX writer" below). The steps branch on the entry's
`record-mode` token.

### `record-mode: create`

The four steps as before this change — nothing here is new:

1. Build the EDR's `--data` JSON per `node ~/.claude/scripts/render-artifact.js --type edr --schema`
   (`status: Accepted`, `domain`/`slug` from the proposal's `record:` token, `title` from the
   proposal's summary, the rest of the body from the proposal's rationale plus whatever the
   implementing task itself establishes — `Reglas verificables`, `Alternativas consideradas`,
   `Consecuencias`). **Every item under `Reglas verificables` is `{ mechanism, rule }`, never a bare
   string** — `mechanism` is `auto` when a test, lint, schema or CI check enforces the rule, `manual`
   when nothing does. The renderer turns that into `- **[{mechanism}]** {rule}`, which is what
   `validate-artifact.js`'s `bullet-prefix` check (`edr.bullet-prefix.regla-sin-marca`) looks for; an
   item written as plain text, or with `mechanism` omitted, ships without the marker and is flagged.
   You decide `auto` vs `manual` from what the implementing task actually established — never default
   to one without checking whether the rule has a real enforcement mechanism in this change.
2. `node ~/.claude/scripts/render-artifact.js --type edr --data <file>` → the record's body. **Write
   it to `.matecito-ai/edr/<domain>/<slug>.md` yourself** — the script never writes to disk in either
   invocation (see "render-artifact.js never writes" below); this is `sdd-apply`'s own file write. A
   `create` naming a file that already exists is a failure, not an overwrite — see "Declaration versus
   reality" below.
3. `node ~/.claude/scripts/render-artifact.js --type edr --data <file> --index-entries` → the domain
   and root INDEX rows, as JSON, in a **second, separate call** (same `--data`, no writes either).
   Isolated Run Mode carries these rows in its Task Run Report under `### Decisions Materialized`
   instead of applying them — see "The INDEX writer" below.
4. Record the outcome as `materialized` — see "Recording the outcome" below.

### `record-mode: modify`

Edits the named record **in place** instead of rendering a new one. Steps 1-2 of the `create` path do
NOT run at all: `render-artifact.js` is a creator, and re-rendering means retyping every section this
change does not touch, which is the cost the ratified alternative (full re-render from `--data`, option
D) was rejected for.

1. Open `.matecito-ai/edr/<domain>/<slug>.md`. A `modify` naming a file that does not exist is a
   failure, not a first materialization — see "Declaration versus reality" below. **`modify` cannot
   bootstrap an absent store** — see "Bootstrapping" below.
2. Edit **only** the clauses the design's entry names, leaving every other byte identical. Rewrite
   the record's INDEX row (its "Consultá cuando…" / trigger cell) **only when the proposal states the
   record's trigger changed**; otherwise leave that row alone. No `render-artifact.js --index-entries`
   call runs, and no new INDEX row is added — the record already has its one row
   (`structure/root-index-cardinality-per-domain-type.md`); see "The INDEX writer" below.
3. `node ~/.claude/scripts/validate-artifact.js --type edr --file <path>` — the structural check the
   `create` path gets for free from the renderer, run explicitly here since nothing rendered this time.
   Any finding is a failure.
4. Record the outcome as `modified` — see "Recording the outcome" below.

### Declaration versus reality, checked both directions

`modify` naming a file that does not exist, or `create` naming one that already exists, is a
**failure** — never a silent switch to the other path, and never an overwrite. `sdd-apply` never probes
the filesystem to *decide* which mode to use; it probes only to *contradict* a declaration that does not
match what is on disk.

### Recording the outcome

In `apply-progress`'s (or the Task Run Report's) `### Decisions Materialized` table — `record | task |
result` — `result` is `materialized` (created), `modified` (edited in place), or `failed: <reason>` on
any failure (invalid data, the renderer or validator refusing to run, an impossible write, or a
declaration-versus-reality mismatch). **A failed materialization does NOT mark the implementing task
complete**, does not leave a partial or malformed record on disk, and is named explicitly in the
return — the code already written is not reverted; the gap is what `sdd-verify`'s `decision-gaps` group
finds next.

### `render-artifact.js` never writes — confirmed by execution (closes the design's Open Question 1)

Both invocations of `render-artifact.js` (`--data` alone, and `--data --index-entries`) write **only
to stdout** — never to disk, in either call. Verified directly: running both against a fixture EDR
left `.matecito-ai/edr/` on disk completely untouched. There is no single-call form that both renders
the body and scaffolds an absent INDEX — it is always two separate invocations, and **materializing is
entirely the caller's job**: `sdd-apply` writes the `.md` body itself, and separately applies the
`--index-entries` JSON rows to the domain and root `INDEX.md` files, scaffolding either one from
`references/edr/templates/index-domain.md` / `index-root.md` when it does not exist yet — the script
does none of that for you.

## The INDEX writer — once, at batch close (not `single-writer-per-batch`)

Applying the `--index-entries` rows to `.matecito-ai/edr/<domain>/INDEX.md` and
`.matecito-ai/edr/INDEX.md` is governed by `structure/root-index-cardinality-per-domain-type.md`
(**one entry per record, no duplicates** — not by `contracts/single-writer-per-batch.md`, whose scope
is the `apply-progress` artifact and task state, and stretching it to the INDEX file would be
re-deciding it by analogy, not applying it). This split governs `create` rows only: a `modify` row
carries no `--index-entries` rows forward at all — the existing INDEX row is edited in place, per
"Materialization" above, so there is nothing for an isolated run to hand the consolidation run for it.

- **Isolated Run Mode**: writes only the record's `.md` body, inside its one commit, beside the code it
  governs (`contracts/one-commit-per-isolated-run.md`). It does NOT touch either INDEX file. It carries
  the rendered `--index-entries` rows in its Task Run Report's `### Decisions Materialized` section so
  the consolidation run has them.
- **Consolidation run**: applies every carried-forward `--index-entries` row **once**, after the
  cherry-pick loop, deduping the root INDEX row by `domain` (two isolated tasks materializing records in
  the same domain contribute one root row, not two) — `structure/consolidation-run-is-the-integrator.md`.
  Scaffolds an absent `INDEX.md` from the templates above on first write.
- **Serial mode**: does both — writes the record body and applies the INDEX rows — in the same step,
  since there is no isolation split to observe.

No registered record is ever lost or duplicated across a parallel batch: the body always lands (each
isolated task's own commit), and the INDEX row for it is applied exactly once, by the one writer the
kernel's continuity rule already names for this phase (`~/.claude/matecito-ai.md` →
"Apply-Progress Continuity").

## Bootstrapping — the mechanism runs even with no store

`.matecito-ai/edr/` absent does NOT disable this mechanism — it is the opposite of the presence-based
activation gate every EDR *reader* uses. The first materialization creates the domain folder, the
record, and both INDEX files (scaffolded from the templates) as a side effect of the ordinary
materialization step above. Every reader downstream of that point sees a store that now exists; the
gate they check keeps working exactly as documented, it simply now finds content.

**`modify` cannot bootstrap.** A `modify` proposal against a domain/slug that does not exist is the
declaration-versus-reality failure above, not a first materialization — the first record for any given
`domain`/`slug` is always created, never edited into existence.

## `sdd-verify`'s `decision-gaps` group — always on, two checks, nothing else

Runs on **every** `development` change — no flag, and it does not depend on `.matecito-ai/edr/`
existing (a change that materialized nothing finds nothing and says nothing; a change that
materialized something is checked whether or not other, older EDRs exist). Source of the check list:
`apply-progress`'s `### Decisions Materialized` table (verify never reads the dispatch prompt — it was
not there; it reads what apply recorded).

For each row of that table:

- **`result` is not `materialized` and not `modified`** (a failed materialization) → CRITICAL, structure
  column names the failure reason, backing column is `—`. This is the "propuesta ratificada sin registro
  materializado" scenario.
- **`result` is `materialized`** → run exactly two checks, both structural/mechanical, never semantic:
  1. **Structure** — `node ~/.claude/scripts/validate-artifact.js --type edr --file
     .matecito-ai/edr/<domain>/<slug>.md`. Exit 0 → `OK`; any finding → CRITICAL, naming the file and
     the violated section.
  2. **Backing ("code that corresponds")** — a **join, never a search**: look up the implementing
     `task`'s row(s) in `apply-progress`'s `### Files Changed`. If that task's changed-file set
     contains **at least one path outside `.matecito-ai/edr/`**, backing is `OK`; if every path the
     task touched is under `.matecito-ai/edr/` (the record and the INDEX and nothing else), backing is
     CRITICAL — the decision has no implementation behind it in this change. No reading of the record's
     prose, no hunting for matching code elsewhere: the check is exactly this set difference.
- **`result` is `modified`** → the **structure** check runs unchanged (same command, same OK/CRITICAL
  reading, against the edited file). The **backing** check does **not** run: it is declared **not
  applicable**, and the cell reads `n/a — record edited in place, governance-only` — the row is neither
  `OK` nor `CRITICAL` on that column. Reason: the backing check's premise is "a decision with no
  implementation behind it **in this change**"; for an edit-in-place, the implementation is the payload
  edits landing in other tasks of the same change, and the check — a set difference over one task's own
  changed-file set — has no way to see them. Declaring it not-applicable is honest; passing it would be
  a false `OK`. Alternatives rejected: leaving the group unaware of `modified` skips the row silently,
  which is worse than either verdict (an unrecognized `result` value is an "absence is not a clean
  verdict" failure); running backing as written fails every in-place edit by construction (a
  governance-only task touches nothing outside `.matecito-ai/edr/`); forcing an unrelated file into the
  task to satisfy the set difference is gaming the check, not answering it. **Accepted cost, stated
  plainly and NOT mitigated**: a `modified` row whose edit really is unbacked — a governance change with
  nothing behind it anywhere in the change — now passes unnoticed. No check catches it, and none is
  proposed. The section's columns (`record | task | structure | backing`) are unchanged by this — the
  cells are free text, so no return-contract edit is needed.

Coherence BETWEEN records (whether two EDRs contradict each other) is explicitly out of scope for this
group — that is `development-decisions-validate`'s standing job, unrelated to what this change
materialized.

`## Decision Gaps` is emitted **only when the change materialized at least one record**
(`when: records_in_change`, read from whether `### Decisions Materialized` in `apply-progress` carries
any row) — a change that materialized nothing gets no section and no mention, exactly like the EDR
activation gate elsewhere. Columns: `record | task | structure | backing`.

## Confirming the second Open Question — `validate-return.js` names the finding, it does not crash

Before the free-form-token fix above, a return carrying an unrecognized-shape `record:` token did not
crash `validate-return.js` — it exited 1 with a **named finding** (`TOKEN-ILLEGAL`, quoting the item
and the illegal value), exactly like any other token violation. Reproduced by execution against a
hand-written fixture return. This closes the design's second Open Question: the validator degrades to a
reported finding, never to a run it cannot complete.

## Standalone paths are unaffected

`development-decisions-bootstrap` keeps its standalone role, callable any time a human wants to capture
or update a decision outside the flow; when both this mechanism and bootstrap could produce the same
record, the one this flow materializes is the one in effect.
