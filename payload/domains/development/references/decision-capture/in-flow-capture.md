# In-flow decision capture — development

<!-- matecito-ai: single definition of the mechanism this domain uses INSTEAD of the kernel's generic
     post-verify mine gate (the kernel's override clause lets a domain opt out that way — see
     `~/.claude/matecito-ai.md`, `### Decision-Gap Capture (mine gate)`). `design` keeps the kernel's
     mechanism unchanged; `development` declares this one. Every phase and script mentioned here reads
     THIS file for the mechanism's rules — they cite it, they do not restate it. -->

An architecture decision that surfaces while a `development` change is in flight is **proposed** by
the phase that finds it, **ratified once** at the lane's gate, and **materialized** as an `Accepted`
EDR in the same `sdd-apply` step that implements the code the decision governs. There is no post-verify
mining pass for `development` — the kernel's generic mine gate never fires here (see the override
clause it carries).

## The proposal — one mailbox item, two tokens

A proposal is a single item under a phase's `### New Decisions` return mailbox (`sdd-spec`'s or
`sdd-design`'s — see "The ratification gate" below for which one, per lane). It travels in two
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
```

- **`summary`** — the point in dispute and the option taken, in one line.
- **`· rationale:`** — the alternatives considered and why this one, still one line (the full
  reasoning; the gate prints `summary`, `rationale` ships in the block for whoever asks).
- **`· blocking-test:`** — unchanged from before this change; see `sdd-design.md`, "The blocking-test
  token".
- **`· record: <domain>/<slug>`** — the identity the EDR would occupy if ratified. Free-form: the
  engine accepts any present, non-null value for a token declared without `values` (see "The
  free-form-token fix" below). It is still **required** — an item with the token line missing fails
  `TOKEN-MISSING` at the Return Contract Check, same strict reading as any other omitted token.

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

## The ratification gate — exactly once, per lane

| Lane | Gate |
| --- | --- |
| `full`, and `custom` with the `design` add-on | `sdd-design`'s `### New Decisions` (unchanged from before this change) |
| `reduced`, and `custom` without the `design` add-on | `sdd-spec`'s `### New Decisions` (new, conditional — see below) |
| `direct` | none — the flow does not run, so there is no proposal and no record |

No later phase re-asks a proposal the gate already ratified, and `sdd-apply` never opens a second
confirmation for it — the ratified text reaches it verbatim through the orchestrator's dispatch prompt
(the same channel already used for `delivery_strategy`, strict-TDD, and the apply-progress continuity
note). An adjustment the user makes AT the gate wins for free: what is in the dispatch prompt IS what
was ratified. Automatic mode does not skip this gate — same as every other Tier-1 mailbox.

### `sdd-spec`'s `### New Decisions` — conditional, same title as `sdd-design`'s

Emitted **only** when the lane running has no `design` add-on active — read from the intake brief's
`### Triage` line (`Lane: ... — add-ons: [...]`, a line `sdd-spec` already reads for other purposes),
never re-derived any other way. When `design` IS in the lane's add-ons, `sdd-spec` emits nothing here:
`sdd-design`'s mailbox is the one and only gate, so a second one would ask the user to ratify the same
decision twice. The title is byte-identical to `sdd-design`'s plain variant (`### New Decisions`) —
same guard rule, same items shape, same tokens — because it is the same mailbox concept surfacing at a
different point in the pipeline, not a new kind of thing.

## Materialization — `sdd-apply` Step 4b, same step as the implementing task

For each ratified proposal forwarded in the dispatch prompt, `sdd-apply` materializes it in the **same
work-unit step** that implements the code the decision governs — never a separate pass before or after.

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
   invocation (see "render-artifact.js never writes" below); this is `sdd-apply`'s own file write.
3. `node ~/.claude/scripts/render-artifact.js --type edr --data <file> --index-entries` → the domain
   and root INDEX rows, as JSON, in a **second, separate call** (same `--data`, no writes either).
   Isolated Run Mode carries these rows in its Task Run Report under `### Decisions Materialized`
   instead of applying them — see "The INDEX writer" below.
4. Record the outcome in `apply-progress`'s (or the Task Run Report's) `### Decisions Materialized`
   table — `record | task | result` — `result` is `materialized` on success, or `failed: <reason>` on
   any failure (invalid data, the renderer refusing to run, an impossible write). **A failed
   materialization does NOT mark the implementing task complete**, does not leave a partial or
   malformed record on disk, and is named explicitly in the return — the code already written is not
   reverted; the gap is what `sdd-verify`'s `decision-gaps` group finds next.

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
re-deciding it by analogy, not applying it):

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

## `sdd-verify`'s `decision-gaps` group — always on, two checks, nothing else

Runs on **every** `development` change — no flag, and it does not depend on `.matecito-ai/edr/`
existing (a change that materialized nothing finds nothing and says nothing; a change that
materialized something is checked whether or not other, older EDRs exist). Source of the check list:
`apply-progress`'s `### Decisions Materialized` table (verify never reads the dispatch prompt — it was
not there; it reads what apply recorded).

For each row of that table:

- **`result` is not `materialized`** (a failed materialization) → CRITICAL, structure column names the
  failure reason, backing column is `—`. This is the "propuesta ratificada sin registro materializado"
  scenario.
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

`development-decisions-mine` keeps its Mode A brownfield scan, invoked directly by the user — unrelated
to this mechanism, and unchanged. Its Mode B (in-flow, driven by the kernel's post-verify mine gate)
has no `development` caller anymore: nothing in this domain's flow dispatches it. `development-decisions-
bootstrap` keeps its standalone role, callable any time a human wants to capture or update a decision
outside the flow; when both this mechanism and bootstrap could produce the same record, the one this
flow materializes is the one in effect.
