<!-- matecito-ai: DEVELOPMENT DOMAIN FRAGMENT.
     Appended to core/CLAUDE.md at deploy time. This fragment binds the
     kernel's generic nouns to development's concrete vocabulary and adds the
     code-specific rules that must NOT live in the kernel. -->

# matecito-ai — Development domain

## Domain vocabulary (binds the kernel's generic slots)

| Kernel slot | Development binding |
| --- | --- |
| Structured flow name | SDD (Spec-Driven Development) |
| Phase pipeline | `intake → explore → propose → spec → design → tasks → apply → verify → archive` |
| Mandatory base phases | `intake → spec → apply → verify → archive` |
| Optional add-on phases | `explore`, `propose`, `design`, `tasks` |
| Phase agents | `sdd-*` (`sdd-intake`, `sdd-explore`, …, `sdd-archive`) |
| Phase dispatch | the orchestrator dispatches the phase agent — there are no `/sdd-*` commands (see "How a phase runs") |
| Alignment artifact | `spec` |
| Decision record | `EDR`, stored in `.matecito-ai/edr/` |
| Decision-record concept reference | `~/.claude/references/edr/README.md` |
| Canonical catalog | `design-patterns` at `~/.claude/references/design-patterns/` (`Applied pattern: X` → `patterns/<x>.md`) |
| Decision-capture mechanism | in-flow — propose (per phase) → ratify once (lane gate) → materialize (`sdd-apply`); no flag, no post-verify mine gate — `~/.claude/references/decision-capture/in-flow-capture.md` |
| Decision-mining executor | `development-decisions-mine` — standalone brownfield scan (Mode A) only; its in-flow Mode B has no `development` caller |
| Decision-capture skill | `development-decisions-bootstrap` — standalone use only, no flow hook |
| Exploration index | CodeGraph (`mcp__codegraph__*`), active when `.codegraph/` exists |
| Library docs resolution (version, config, API, migration, debugging) | `resolve-library-docs` skill, backed by the `context7` MCP — the sole choke point before any version is fixed into an EDR (bootstrap's "Versión" question, mine's close-of-scan second pass) |
| Guards | `strict-tdd`, `review-workload` |
| Engram topic-key namespace | `sdd-init/{project}` · `sdd/{project}/testing-capabilities` · `sdd/{change-name}/{intake,explore,proposal,spec,design,tasks,apply-progress,verify-report,archive-report}` |<!-- matecito-ai: `state` removed — a declared format with no instructed producer and no consumer; recovery goes through the phase artifacts themselves. `testing-capabilities` added: it was missing from the namespace despite two phases reading it. -->
| Init topic key | `sdd-init/{project}` |

## Language
Code (variables, functions, classes, constants): English. Comments follow the `code-comments` skill (which also fixes their language as English).

## Contract & definition shapes — never inferred
When you are about to create or modify a **contract or definition** — a domain entity, a database model / migration / schema, an API request/response (DTO), a public/exported type, interface, or enum, an event payload, or a config schema — you MUST NOT infer its shape. Both **which properties it has** and **each property's type** are decisions the human owns; do not invent "obvious" fields or "most likely" types (an `id`, an `email`, a `createdAt`; `string` vs `uuid`; float vs integer minor units — these are decisions, not defaults).

- **Artifact-pinned → execute, unless it conflicts.** If the shape is already fixed by an upstream artifact (spec, design, an EDR modeling policy, or the user's explicit request) **and it is coherent**, that IS the mandate — implement it, do not re-ask. Only the un-pinned parts are open. But check coherence before executing: STOP and open the discussion when
  - **the user's request contradicts an Accepted EDR** — do not silently pick a winner. Show both sides (what the user asked for; what the EDR fixes and why) and discuss it. The outcome may be adjusting the request, or updating the EDR via `development-decisions-bootstrap` (update mode); or
  - **the artifact that pins the shape is internally inconsistent, or does not cover this case** — an EDR was written at a point in time and may not have foreseen this. Say so: do not stretch it by analogy, and do not follow it to the letter knowing it breaks something.

  Silently executing a mandate you can see is inconsistent is NOT respecting the decision — it is hiding the problem. The goal is a back-and-forth until you land on a middle ground, not the artifact winning by default nor the request winning by default. As a headless executor you cannot hold that discussion yourself: return `blocked` with the conflict stated and the concrete options; the main thread carries the conversation.
- **Unspecified → ask, per whole contract.** Propose the FULL contract (all fields + their types) as one reviewable unit — never field-by-field. With several unspecified contracts, default to one at a time (they often depend on each other); tell the user how many there are and offer "one-by-one or all-at-once" so they set the pace.
- **Where the answer lives.** The concrete shape (field names + types) belongs in the **code** (or the `design` artifact that materializes it) — NEVER copied into an EDR as a typed struct (that is a code calco; EDR reasoning stays conceptual). Only a **cross-cutting modeling policy** hidden in the answer ("identifiers are UUIDs", "money as integer minor units", "status as enums, not magic strings") may be captured as an EDR, expressed conceptually — and once captured it pins that part, so you stop re-asking it (per artifact-pinned). Offer to capture such a policy; never force it.
- **Scope.** Targets shapes that persist, cross a boundary, or are public. A transient internal struct used within a single function is execution detail, not a contract — no need to ask.

This is a specialization of the kernel's "Open question = blocked, not permission" for the high-stakes case of contracts, where inference is most tempting and most consequential (it propagates to DB, API, and tests).

## CodeGraph
Code exploration prefers CodeGraph when `.codegraph/` exists (structural questions); grep for literal text or non-indexed files. The SDD fork assumes the `mcp__codegraph__*` prefix. Reference the server by capability — never hardcode individual tool names (they drift between server versions); resolve the actual registered tool names at use time.

## Architecture diagrams (drawio)
Diagramming here is two complementary pieces: the **`drawio` skill** owns the *vocabulary* — how to build the diagram XML (shapes, branded/AI icons via `shapesearch`/`aiicons`, style presets, layout, diagram-type templates) — and the **`mcp__drawio__*` MCP** owns the *live render* — it renders the skill's `<mxGraphModel>` as an ephemeral preview; the skill itself never writes files. This rule is the **single source of truth for when to draw**. Diagrams are generated **on demand, never automatically**, and only when the change has structural complexity worth visualizing. **Diagram inference test — generate when** the change introduces or rewires ≥3-4 components with relationships, data flow crosses boundaries (layers/services/new modules), there is a non-trivial process with branches or states, or the task is to understand existing code spread across many files (CodeGraph can feed the graph) — **plus** capturing the shape of an architectural decision (EDR). **Do NOT generate** for a small fix, rename, config tweak, single-file/single-unit change, or linear logic — there prose or a snippet is clearer. **Model — offer-and-confirm, never unilateral.** **Decide vs generate (timing):** the structure does not exist yet at intake, so `sdd-intake` only *decides* — it sets `diagram: needed | not-needed` in the brief per this test, and the user confirms it at the **INTAKE GATE** (that gate IS the confirmation for the in-flow case; no re-ask later). **Generation is EPHEMERAL — always a live preview, NEVER a file in the project (zero `.drawio` artifacts in the repo).** The diagram is rendered by the **main thread** with a live preview (`mcp__drawio__*` — `start_session` reports the preview URL; the port is assigned dynamically, not fixed); nothing is exported or persisted. When the flag is `needed`, the main thread offers to render it live at the design step — the **headless `sdd-design` sub-agent does NOT generate or export diagrams**; it only notes that a live diagram is recommended. Same for a `direct` lane / outside the flow. Apply this same test before offering.

## Debugger MCP (mcp-debugger)
The debugger MCP (`mcp__debugger__*`, DAP step-through via `@debugmcp/mcp-debugger`) is **on-demand only** — it is NEVER invoked automatically. Its primary home is `sdd-apply`: when a runtime defect is encountered and the per-language debug toolchain is available (detected by `sdd-init` and cached in `sdd/{project}/testing-capabilities`), `sdd-apply` MAY diagnose the root cause AND apply a fix in the same context. In `sdd-verify`, the debugger is **diagnosis-only**: it MAY be used to understand why a test or scenario fails, but MUST NOT apply fixes there — any fix found belongs in a subsequent `sdd-apply` invocation. When the per-language debug toolchain is absent (`debugger.available = ❌` in testing-capabilities), both phases skip debugger usage silently — no error, no warning, no section. **For the full usage guide** — preflight (adapter vs. toolchain binary distinction), per-language install helper, and the debug loop — read the **`debugger` skill** (`~/.claude/skills/debugger/SKILL.md`).

## SDD Flow

```
sdd-intake → sdd-explore → sdd-propose → sdd-spec → sdd-design → sdd-tasks → sdd-apply → sdd-verify → sdd-archive
                                              ^
                                           (design reads EDRs)
```

<!-- matecito-ai: esta sección listaba `/sdd-init`, `/sdd-intake`, … como comandos de usuario, y no
     existen: no hay `~/.claude/commands/` en el deploy y las 11 skills `sdd-*` llevan
     `disable-model-invocation: true` junto con `user-invocable: false`, que cierran las dos puertas.
     Un usuario que tipeaba `/sdd-verify` no encontraba nada, y la lista además omitía cuatro fases
     (propose, spec, design, tasks) siendo `spec` obligatoria. Lo que se documenta ahora es la vía que
     realmente existe. -->
### How a phase runs

There are no `/sdd-*` slash-commands: the phase skills are not user-invocable, and the model cannot
self-invoke them either. A phase runs by **the orchestrator dispatching its agent** (`sdd-intake`,
`sdd-spec`, …), and that agent reads its own method from `~/.claude/skills/<phase>/SKILL.md`. Most
phases dispatch as exactly one agent, one call, one return; two named cases fan out into more than one
dispatch instead — see "Phase fan-out" below.

So you ask for the work in plain language ("arreglá X", "agregá Y"); the orchestrator resolves the lane
at the INTAKE GATE and drives the pipeline from there. Naming a phase ("corré el verify", "seguimos con
apply") is a request to dispatch that agent, not a command the harness resolves.

<!-- matecito-ai: phase fan-out. `sdd-verify` used to run as ~9 serial steps in one agent (~30 min
     wall-clock); that account is unchanged below. `sdd-apply` gained a second, independent fan-out
     case for a batch with independence-marked tasks — the mechanism, dispatch shape and consolidation
     step differ enough from verify's that each keeps its own definition file; only the naming and the
     exclusivity clause are shared prose. Neither partition is restated here — verify's lives in
     `subverifier-groups.md`, apply's in `parallel-batch.md` — so the skill, the agent and this
     fragment cite one definition each instead of drifting apart. -->
### Phase fan-out (the two declared cases)

Two phases fan out instead of running as a single dispatch. Both are named here on purpose: naming
only one would read as licence to add a third without its own change — it is not.

**`sdd-verify`.** Instead of one instance running every check in series, the orchestrator dispatches,
**in a single message**, one `sdd-verify` instance per group whose gate is active — so their runs
overlap in time instead of chaining. Each instance receives a `group` (one of `execution`,
`correctness`, `design-coherence`, `edr-coherence`, `spec-coherence`, `ui`, `decision-gaps`) and returns
only that group's fragment; the orchestrator merges every fragment into the single `## Verification
Report` and persists it once. Full partition — which group owns which check, the gate that turns each
on, the Sub-Report envelope, and the mechanical merge —
`~/.claude/references/phase-returns/sdd-verify/subverifier-groups.md`.

**`sdd-apply`, for each `· parallel-group: <id>` with two or more tasks.** Eligibility is evaluated
**per group**, never over the whole tasks artifact — two groups of two are two eligible rounds, not one
eligible batch of four. Before dispatching an eligible round, the orchestrator runs the
**Uncommitted-Work Gate** (see "Uncommitted-Work Gate" under Guards, below); a dirty tree no longer
degrades a round to serial by itself — the gate's three outcomes decide that. For each eligible group's
round, the orchestrator dispatches, in one
message, one `sdd-apply` **isolated run** per task of that group — each repositioning its worktree onto
the round's **immediate container's** local `HEAD` before writing anything (the change workspace's HEAD
when change-level isolation is active, the working branch's local HEAD when it is not — never the
harness's starting point, which can be `origin/<branch>` and lag behind unpushed local commits), each in
its own worktree, each closing with exactly one commit on its own branch, none of them writing
`apply-progress` or marking a task. Once the round returns, the orchestrator dispatches a further
`sdd-apply` **consolidation run** — no isolation, a second mode of the same agent, never the orchestrator
and never a new agent — which cherry-picks every commit onto that same container one at a time (the
change workspace when isolation is active, the working branch when not), in ascending task-id order, and
is the round's only writer. Groups never mix in the same round; rounds run one after another, in
ascending order of each group's lowest task id. A group with no marks, or fewer than two members, runs the unchanged
single-dispatch path. Before forming any round, the orchestrator validates the tasks artifact's marks by
script — see "Parallel-mark validation" under Guards, below. Full mechanism — eligibility, the
uncommitted-work gate, repositioning onto the local base, the base handshake, the commit convention, the
Task Run Report, integration and conflict handling — `~/.claude/references/phase-returns/sdd-apply/parallel-batch.md`.

Each case has its own shape and its own definition file; read the one for the phase you are dispatching
or consolidating — this fragment does not keep a second copy of either. The Executor boundary still
holds for every dispatched instance in both cases: none of them dispatches anything, including each
other. **This exception is scoped to exactly these two phases** — no other phase in the pipeline fans
out, and this section is not an invitation to add a third without its own change.

### SDD Phase Read/Write

| Phase | Reads | Writes |
| --- | --- | --- |
| `sdd-intake` | raw request | `intake` |
| `sdd-explore` | intake (brief) | `explore` |
| `sdd-propose` | exploration (optional) | `proposal` |
<!-- matecito-ai: spec pasó a leer el intake brief SIEMPRE, no sólo como upstream de fallback: es el
     único lugar que lleva el flag `ui-test`, y la proposal no lo transporta. Sin esa lectura, la
     producción de `ui-scenarios` funcionaría en lane `reduced` y fallaría en `full`. -->
| `sdd-spec` | proposal (required) + **intake brief (always, for the `ui-test` flag)** + **durable capability-spec** (for Modified Capabilities) | `spec` (incl. the **behavioral** `ui-scenarios` when `ui-test: needed` — domain language, no routes or locators) |
| `sdd-design` | proposal + **intake brief (always, for the `diagram` flag)** + **EDRs** + **durable capability-specs** (required) | `design` |
| `sdd-tasks` | spec + design + **durable capability-specs touched** (required) | `tasks` |
| `sdd-apply` | tasks + spec (incl. the behavioral `ui-scenarios`) + design + apply-progress + **ratified decision proposals, forwarded verbatim in the dispatch prompt** (never re-read from Engram or an artifact — see `in-flow-capture.md`) | `apply-progress` (incl. `### UI Scenario Counterparts` — the **executable** half, with the real routes and locators it built — and `### Decisions Materialized`, when it materialized at least one ratified proposal) |
<!-- matecito-ai: parallel-batch note — the row above is the full-lane, single-dispatch ideal; a batch
     with independence-marked tasks splits it across the two fan-out roles instead of changing the row. -->
| ↳ isolated run (parallel batch) | tasks (its one task) + spec + design — never `apply-progress` | nothing (single-writer rule — returns a Task Run Report, see `parallel-batch.md`) |
| ↳ consolidation run (parallel batch) | the batch's N Task Run Reports + prior `apply-progress` | `apply-progress` (as above, plus `### Integration Log`) |
| `sdd-verify` | spec (incl. the behavioral `ui-scenarios`) + design + tasks + apply-progress (incl. the **counterparts**, paired by `name`) + **intake brief (always, for the `ui-test` flag)** + **EDRs touched** + **capability-specs touched** | `verify-report` |
| `sdd-archive` | all artifacts | `archive-report` + **durable capability-specs (merge)** |

The "Reads" column lists the **full-lane** ideal. In `reduced`/`custom` lanes some upstream phases don't run, so each phase reads the **nearest available upstream**: `sdd-spec` falls back to the intake brief when there is no proposal; `sdd-apply` treats `spec` as the floor and skips `tasks`/`design` when absent. The **durable capability-specs** are read only when `.matecito-ai/development-specs/` exists; absent → skip silently (same presence-based gate as EDRs).

<!-- matecito-ai: the kernel's INTAKE GATE surfaces the lane plus "the decision flags the domain declares",
     and stays domain-agnostic on purpose — this block is that declaration. Before it existed, the
     instruction to surface these two lived only in `agents/sdd-intake.md` and its skill, both read by the
     intake executor and neither read by the orchestrator that has to act on it. The flags were visible
     inside the brief but never raised as something to ratify, and since neither is re-asked later, the
     phases downstream acted on values nobody confirmed. -->
### Brief decision flags (confirmed at the INTAKE GATE)

`sdd-intake` decides these on the user's behalf and writes them into the brief under
`### Classification`. The orchestrator surfaces each at the INTAKE GATE, by name, with its value and
its one-line reason. **That gate is their only confirmation** — no phase re-asks, and each reader
treats an absent or unconfirmed flag as `not-needed` and closes **silently**, so a flag that slips
through the gate is a check nobody notices was skipped.

| Flag | Line in the brief | Decided by | Read by | What it drives |
| --- | --- | --- | --- | --- |
| `diagram` | `- Diagram: {needed\|not-needed}` | `sdd-intake` per the diagram inference test in `## Architecture diagrams (drawio)` above | `sdd-design` | Whether a **drawio** architecture diagram is warranted. `sdd-design` only NOTES the recommendation in its `executive_summary`; the main thread renders it live via `mcp__drawio__*`. Nothing is ever written to the repo. |
| `ui-test` | `- UI test: {needed\|not-needed}` | `sdd-intake`, by keyword inference over the request (`browser`, `page`, `form`, `screen`, `visual`, `click`, `render`), overridable explicitly in the request | `sdd-spec`, `sdd-verify` | Whether UI verification via **proofshot** is warranted. `sdd-spec` authors the `ui-scenarios` block only when this is `needed`; `sdd-verify` runs the ProofShot session only when this is `needed` AND `uiTest.available = ✅`. |
| `components` | `- Components: {name[, name...] \| unassigned}` | `sdd-intake`, by mapping the request's scope against `repo.components[].paths` | ninguno | Nothing — it is metadata for the person confirming the gate, not a phase input. **Presence-based, unlike the two above**: with no `repo.components` declared for the project the field does not exist and is never mentioned; declared, it is multivalued and always emitted (`unassigned` when no `paths` match — never omitted to mean "no match"). |
| `isolation` | `- Isolation: {active\|inactive}` | `sdd-intake`, recommended together with the lane per `structure/change-isolation-activation-flag.md` | the orchestrator (kernel's "Change Workspace (opt-in)") | Whether the orchestrator opens a dedicated change workspace for this change. For `direct`/ad-hoc work — which never reaches the INTAKE GATE — the choice is confirmed at the lane fork itself instead; a fork never surfaced means inactive. |

None of these is executed by intake: it decides, others (or, for `components`, no one) act. Adjusting
one at the gate updates the brief like any other correction.

<!-- matecito-ai: git mechanics for the kernel's domain-neutral "Change Workspace (opt-in)" policy —
     `structure/change-workspace-prose-homes.md` fixes this split: the kernel keeps the policy in neutral
     wording, this fragment holds the concrete git commands. -->
## Change Workspace — git mechanics

Binds the kernel's `### Change Workspace (opt-in)` policy to this domain's concrete mechanism: a git
worktree, on its own branch, added from the main repo.

**Identity.** Branch `matecito-ai/<change-name>`; directory
`<repo>/.matecito-ai/workspaces/<change-name>` (`structure/change-workspace-identity.md`). Opened with
`git worktree add <dir> -b matecito-ai/<change-name> <original-branch>`. At open time, before any phase
writes into it, confirm `.matecito-ai/workspaces/` is listed in the repo's `.gitignore` — add the line if
it is not, so the Uncommitted-Work Gate never reads the workspace directory itself as dirty work in the
main repo.

**Forwarding.** While the workspace is open, every phase dispatch prompt carries a `workspace: <absolute
path>` line — readers included (`contracts/workspace-forwarded-in-dispatch.md`). A phase resolves repo
paths under that path and runs git with `git -C <workspace> ...` rather than assuming the session's own
working directory; a phase that runs a project command (a test runner, a linter, a build) runs it with
that same path as the command's working directory, for the identical reason.

**Nesting.** A parallel implementation batch dispatched while the workspace is open nests on it exactly
as "Phase fan-out" above describes: the round's `base` is `git -C <workspace> rev-parse HEAD`, and the
consolidation run's cherry-picks land on the workspace's branch, never on the original branch. Full
mechanics: `~/.claude/references/phase-returns/sdd-apply/parallel-batch.md`.

**Integration (orchestrator, once, at cycle close).** From the main repo: `git merge --no-ff
matecito-ai/<change-name>`. On failure, the orchestrator does not abort right away — it first runs `git
rebase <original-branch>` **inside the change workspace**. A clean rebase means retrying the merge (now a
fast-forward). If the rebase also conflicts, the orchestrator runs `git rebase --abort` then `git merge
--abort`, and reports to the user what happened — which file, in which commit, which side each version
comes from — with a recommendation suited to the case (resolve by hand in the workspace · leave it open ·
another way out). Nothing is forced at any point; the workspace stays intact on every failure path
(`structure/change-level-integration-act.md`).

**Cleanup (after a clean integration only).** `git worktree unlock <dir>` → `git worktree remove <dir>` →
`git branch -D matecito-ai/<change-name>` — never `remove -f -f`; a failure at any step is recorded, not
forced, same as the batch-level cleanup pass in `parallel-batch.md`
(`structure/change-workspace-cleanup.md`). After a failed integration, the workspace and its branch are
both kept, untouched, for inspection.

## Guards

### Strict TDD (resolution + forwarding)
Same precedence as model resolution — per-project `domainConfig.development.strictTdd` → global `domainConfig.development.strictTdd` → `false` (pre-M7 flat top-level `strictTdd` is auto-migrated into `domainConfig.development` on read). Resolve once per session, cache. If effective `strictTdd` is true, add to the `sdd-apply` / `sdd-verify` prompt: "STRICT TDD MODE IS ACTIVE. Test runner: {test_command}. Follow strict-tdd.md." The `{test_command}` comes from `sdd/{project}/testing-capabilities` in Engram.

<!-- matecito-ai: `sdd-intake` tenía orden de "ask 2-4 questions" corriendo headless, sin canal con el usuario: se las autocontestaba y el brief salía con respuestas inventadas que el flujo trataba como mandato confirmado. El discovery pasa a dos pasadas. -->
### Discovery Gate (MANDATORY)
`sdd-intake` runs headless and CANNOT ask the user anything. It formulates the discovery form and returns `status: needs-input` with the questions. That return is neither an error nor a blocker — it is the normal first pass.

`needs-input` is a legal envelope status, not an error. The legal status values are enumerated once in the canonical contract (`_shared/sdd-phase-common.md`, Section D.1) — read them there; this gate does not re-declare them.

When intake returns `needs-input`: put its questions to the user yourself (you own the channel), then **re-dispatch `sdd-intake`** with the raw request plus the answers verbatim so it produces the brief on Pass 2. **An empty question list still requires you to go to the user** — intake returns its one-line reading of the request and the user confirms or corrects it; never treat "no questions" as licence to skip straight to the brief. Never answer them on the user's behalf, never trim the list down to the ones you find interesting, and never skip ahead to another phase — no brief exists yet, so there is nothing downstream to run. If the user leaves one open, hand it back as open instead of resolving it for them.

The questions are walked through the shared presentation in `~/.claude/references/gate-presentation.md`
— one index when there are two or more, the fixed item template either way — under the Discovery
Gate's own declared exception: no anchor, since the request has not yet produced anything to point at.
This gate states no index or bulk-action wording of its own.

This gate sits BEFORE the INTAKE GATE and does not replace it: discovery answers first, then the brief, then confirm / adjust / cancel over that brief. **Automatic mode does NOT skip this gate** — Automatic only skips the between-phase "¿Continuamos?" checkpoint, never a question the user has to answer.

<!-- matecito-ai: nada comprobaba que un retorno trajera lo que debía traer. Si una fase se comía una
     sección, el gate correspondiente no disparaba — en silencio, que es el modo de falla que más
     costó detectar. Los templates de `~/.claude/references/phase-returns/` son la especificación
     contra la que se valida; este guard es quien la aplica. -->
### Return Contract Check (MANDATORY)
Before acting on ANY phase return — before the Unresolved Decisions Guard, before routing, before dispatching anything — validate it against that phase's template at `~/.claude/references/phase-returns/<phase>/<phase>.md`. That file is the canonical shape of the return, and matching is **literal**: a title that differs in wording, casing or heading level is a section you will not find, and a gate that will not fire.

<!-- matecito-ai: los cuatro chequeos vivían sólo como prosa, y el motor que los ejecuta existía sin que
     ninguna regla lo invocara: `validate-return.js` se autodescribe como el que hace que este check
     "stop being a human reading and comparing titles by eye", y su nombre no aparecía en ningún `.md`
     del dominio. Comparar títulos a ojo es exactamente el modo de falla que este guard vino a cerrar,
     así que dejarlo en manos del ojo era el defecto reproduciéndose un nivel más arriba. -->
**Run the mechanical engine first.** Write the phase's return block to a temp file and validate it:

```
node ~/.claude/scripts/validate-return.js --phase <phase> --status <status> --file <return.md>
```

Exit 0 → clean, move on to the guards. Exit 1 → it names the violated check with a code
(`SECTION-MISSING`, `SECTION-RELEVELLED`, `SUMMARY-CONTRADICTS-BODY`, `TOKEN-WRONG-MAILBOX`, …); handle
it per "When something fails" below. **Exit 2 means the validator could not run at all — surface that,
never swallow it and never fall back to reading by eye silently**: a check that quietly stops running is
indistinguishable from a check that passes, which is the failure this guard exists to prevent.

The four checks below are the **specification** of what the engine enforces and what each finding means
— read them to interpret a finding, not to perform the comparison by hand.

Four checks:
1. **Every unconditional section is present.** The template marks which ones are always emitted. A missing unconditional section is NOT "nothing to report" — the phase is required to emit it with a `None…` sentinel instead. A section the template marks conditional is legitimately absent when its condition does not hold.
2. **The titles match the template**, including the accepted variants it declares.
3. **The shape matches the status**: a `blocked` return carries the section the template designates for the blocker; a `needs-input` return carries its questions; and so on.
<!-- matecito-ai: a return carrying `### New Decisions: None.` in the body and `Summary: "Key Decisions: 2
     documented"` passed all three checks and dispatched silently: this check looks at presence, titles
     and status, and the Unresolved Decisions Guard reads the sentinel and stays quiet. Nobody put the
     two claims side by side. It goes here and not in the guard because this check runs BEFORE it, and
     the guard is precisely the one that goes mute. And it is deliberately MECHANICAL: "is the summary
     faithful to the body?" would be interpretation, the one thing this check does not do. -->
4. **The `Summary` does not contradict the body.** A section whose body is only the empty sentinel is **declared empty**. If the envelope's `Summary` states a non-zero count or asserts the presence of content for that same section, that is a contradiction — not a wording nuance. The reverse counts too: a section carrying real rows whose `Summary` declares it at zero. **No other disagreement between summary and body fires this check**: you are comparing two explicit claims about the same section, never judging whether prose is a faithful rendering of a table.

**When something fails**, do NOT fill it in yourself and do NOT silently continue — that is the whole point. Surface it: name the phase, the missing, malformed or contradicted section, and let the user choose — proceed treating it as empty, re-run that phase, or adjust. On a check-4 contradiction, quote the two claims side by side and add one factual line, not a verdict: the body is what the gates and the downstream phases read, and the `Summary` is read by nobody else. Which of the two is right is not yours to settle. Re-running re-executes the phase in full (for `sdd-apply`, a re-dispatch is a continuation batch, not a re-emission): say so when you offer it.

**When everything checks out, say nothing** and move on to the guards. This check only speaks when something is wrong.

<!-- matecito-ai: sdd-tasks-parallel-group-contract — the tasks artifact's `· parallel-group:` marks
     were "as emitted" with no mechanical check anywhere, one level below the Return Contract Check.
     Same principle, one artifact down: a script gates it, never eyeballing. -->
### Parallel-Mark Validation (MANDATORY)
Before the orchestrator forms any `sdd-apply` batch from a tasks artifact's `· parallel-group:` marks
(see "Phase fan-out" → `sdd-apply` above), validate it mechanically:

```
node ~/.claude/scripts/validate-parallel-marks.js --file <tasks artifact>
```

Exit 0 → clean, no output, form the batch. Exit 1 → one line per finding, naming the offending task;
that task alone degrades to serial (unmarked) — the finding never halts the phase or the rest of its
group. Exit 2 → the check could not run at all: surface that, never record the artifact as clean, and
form no batch from the marks — every task takes the serial path. Full contract — the closed list of
three malformed shapes, eligibility per group — `~/.claude/references/phase-returns/sdd-apply/parallel-batch.md`.

<!-- matecito-ai: uncommitted-work-gate — short guard that points, same shape as Parallel-Mark
     Validation above; the full mechanism (when it runs, the component mapping and its fallback, the
     three outcomes, the orphan case, the notice) lives once in parallel-batch.md, never duplicated
     here. -->
### Uncommitted-Work Gate (MANDATORY)
For each eligible `sdd-apply` round that will use worktree isolation, immediately after
`validate-parallel-marks.js` and **before** the orchestrator reads `HEAD` as that round's `base`,
inspect the round's **immediate container's** uncommitted changes — the change workspace's tree when
change-level isolation is active, the main repo's when it is not
(`contracts/uncommitted-gate-follows-the-container.md`). Clean tree, or dirty with no relevant
intersection → silent, nothing to do. Dirty and relevant → present exactly three outcomes (commit first
· continue anyway · work on that same container's branch without a worktree) and dispatch nothing until
the user picks one, not even in Automatic mode; picking "continue anyway" leaves a trace the
consolidation run records. Those three outcomes are presented through the shared walkthrough in
`~/.claude/references/gate-presentation.md`, anchored to the dirty paths `git status --porcelain`
already printed — this gate states no index or bulk-action wording of its own. A serial dispatch and
the consolidation run never trigger this gate; it does
not re-ask within the same phase run when the dirty set hasn't changed. Full mechanism — legal inputs,
the component mapping (and its fallback when the project declares no `repo.components`), the orphan
case, the notice shape — `~/.claude/references/phase-returns/sdd-apply/parallel-batch.md` →
"Uncommitted-Work Gate".

<!-- matecito-ai: los buzones de cada fase (Open Questions, New Decisions, Deviations, Derived capabilities, risks) existían sin que nadie los consumiera: eran el lugar barato donde depositar lo no resuelto y seguir. Este guard los convierte en disparador. -->
### Unresolved Decisions Guard (MANDATORY)
After EVERY phase returns and before dispatching the next one, inspect the return envelope for unresolved-decision mailboxes. **Two tiers — do not conflate them.**

<!-- matecito-ai: la lista de buzones vivía duplicada acá y en el contrato canónico, y cada edición
     desalineaba una de las dos copias. Este guard ya NO mantiene su propia lista: referencia la tabla. -->
**Tier 1 — pending decision → STOP and ask.** The Tier-1 mailboxes are **exactly** the sections marked Tier 1 in the canonical mailbox table (`_shared/sdd-phase-common.md`, **Section D.3**) — read the list there. This guard does NOT keep a parallel copy of it. One matching note only: `sdd-design`'s `New Decisions` is titled `New Decisions (not yet in EDRs)` when the EDR store is active and plain otherwise — **emitted either way**, the activation gate does not suppress it. Content in any Tier-1 section means the phase produced something the user has not agreed to. Present it and wait before the next dispatch. **Automatic mode does NOT skip this gate** (same pattern as the INTAKE GATE and the mine gate). When a Tier-1 section's contract splits its items into `summary`/`rationale` (the phase-return contract's `items.rationale`), present each item's `summary` at the gate — that is the section's own **declared** presentation, not you judging what is brief — and reproduce an item's `rationale` verbatim, from the block already in context, only when the user asks for it.

<!-- matecito-ai: `sdd-design`'s blocking test was pure self-assessment — the executor ran it in its head
     and published only the verdict. A decision whose OWN text said "this needs a queue and a worker the
     project does not have today" (axis 1, literally) arrived under `New Decisions` with `status: done`,
     and this guard had no mechanical way to notice: noticing meant reading the decision's prose and
     re-running the test, which is the interpretation this guard exists not to do. Same fix as
     `verify-checks:` for design deviations, one level up: the phase declares, the guard classifies. -->
**Reading the `blocking-test` token (`sdd-design` only).** Each item under `### New Decisions` carries
a `· blocking-test: none | infra | contract | data-model` line. It declares whether the alternatives
differ in any of the blocking test's three axes, and you classify on the token **alone** — never by
reading the decision and re-running the test yourself:

| Token | What it asserts | What you do |
| --- | --- | --- |
| `none` | the test ran and came back negative | ordinary Tier 1 — present it with the rest of the batch |
| an axis named (`infra` / `contract` / `data-model`) | the item is in the wrong mailbox: a differing axis makes the decision `blocked`, not a Tier-1 note | stop and surface it as you would a `blocked` return, quoting the token |
| absent, or hedged | the test did not run, or the answer is being withheld | Tier 1 under the strict reading — same default an undeclared deviation gets in `sdd-verify` |

The token is the only evidence the test ran at all. Do not accept a decision's prose as a substitute
for it, and do not fill one in on the phase's behalf.

<!-- matecito-ai: in-flow decision capture (development-specifics). Full mechanism, the ratification
     gate per lane, the materialization contract: in-flow-capture.md. This is the orchestrator-side
     half — WHO forwards the resolution and WHEN — that neither sdd-spec/sdd-design (who propose) nor
     sdd-apply (who materializes) can instruct on their own, since it is the launch-prompt construction
     step between them. Written for its two readers: the orchestrator, who builds the prompt (first two
     paragraphs), and the `sdd-apply` executor, who reads this fragment as part of its mandatory load
     protocol (`_shared/sdd-phase-common.md`, Section A) and acts on the last two paragraphs. -->
**Forwarding a proposal's resolution to `sdd-apply`.** Every item that reached this gate under `New
Decisions` — `sdd-spec`'s conditional mailbox or `sdd-design`'s — carries a `record: <domain>/<slug>`
token, and stays in the design's own `## New Decisions` prose whatever the gate decided: the item's
mere presence there is NOT evidence of ratification or rejection, and `sdd-apply` MUST NOT read it as
either. Forward each item's resolution explicitly, in the launch prompt of the `sdd-apply` dispatch
that implements the task governing it — never re-written into an Engram key, never left for `sdd-apply`
to re-derive:

- **Ratified** (confirmed or adjusted at this gate) — its ratified text (the adjusted summary/rationale
  if the user corrected it at this gate, not the originally-proposed one) verbatim, plus its `record:`
  token, marked ratified. Unchanged from before this rule.
- **Rejected** — its `record:` token and `summary` alone, marked rejected. The full text already
  travels in the design's `## New Decisions`; the token and summary are enough for `sdd-apply` to know
  which item and what it was about.

The resolution never travels as a token on the mailbox item itself, and none should be added: `New
Decisions` has no `resolution:` field, because `sdd-spec`/`sdd-design` write the item at propose time,
before the gate has run — the phase authoring the item cannot fill a field for an outcome that does not
exist yet. The gate happens after the item is written, and the orchestrator, at the moment it forwards,
is the only participant who ever learns that outcome. That is why this instruction lives here, in the
guard the orchestrator reads to build the dispatch prompt, and not as a field on an item a different,
earlier phase already authored.

No proposal reached the gate this change → nothing to forward, and the prompt makes no mention of this
mechanism. Every proposal that did reach the gate travels with a named resolution — if `sdd-apply`
reaches a task governed by an item the prompt does not mention, it MUST NOT treat it as ratified or
rejected on its own guess; it returns `status: blocked` naming the item and the missing resolution.

Naming a rejection is only how `sdd-apply` learns the resolution — never an instruction to record it.
A rejected item MUST NOT produce an EDR, an INDEX row, or a `### Decisions Materialized` row, and
`sdd-verify`'s `decision-gaps` group MUST NOT check it; `sdd-apply` implements the governed task per
what the design's approach describes. When that approach matches what the rejected proposal proposed,
there is nothing further to do. When they describe **different** implementations for the same point,
the design is internally inconsistent, and `sdd-apply` applies the domain's existing rule for exactly
that case — **"Contract & definition shapes — never inferred"**, above, the "artifact that pins the
shape is internally inconsistent" clause and its close ("return `blocked` with the conflict stated and
the concrete options") — showing both versions (the design's approach, the rejected proposal) rather
than choosing one or stretching either by analogy. This is not an `### Unmandated Forks` item: that
mailbox is for a point NO artifact fixes, and here the design fixes it twice, incompatibly.

This is the one and only channel `sdd-apply` reads a proposal's resolution from — see
`~/.claude/references/decision-capture/in-flow-capture.md`.

**Ratification ledger (re-emergence support).** At gate close — in the same step that builds the
`sdd-apply` dispatch prompt — the orchestrator writes one row per ratified item to
`sdd/{change-name}/ratified-decisions`: `record` (the item's `<domain>/<slug>` token), `ratified_summary`
(the adjusted text when the user corrected it at this gate, the offered text otherwise), `anchor` (per
`~/.claude/references/gate-presentation.md`'s anchor rule), and `gate` (which gate ratified it). Only a
LATER gate of the same change reads this key — to recognize a decision that re-emerges (matched on
`record`, exact string — see `~/.claude/references/gate-presentation.md` → "Re-emergence") and offer its
short form instead of walking it again. **`sdd-apply` MUST NOT read `sdd/{change-name}/ratified-decisions`**:
its one and only channel for a proposal's resolution stays the dispatch prompt, per the paragraph
above — a second, independently-read channel is exactly the drift `in-flow-capture.md`'s single-channel
rule exists to prevent. Not pruned: the per-change key is its own cleanup, same as `apply-progress`.

<!-- matecito-ai: return-side half of the forwarding paragraph above. The send-side half says WHAT gets
     forwarded and WHEN; this half says how the orchestrator reads what came back — the mandatory
     `design-conflict` verdict `sdd-apply` now declares for each rejection it checked
     (`~/.claude/references/phase-returns/sdd-apply/sdd-apply.md`, `### Rejected Proposals Checked`).
     Same pattern as "Reading the `blocking-test` token" below: classify on the token alone, never by
     re-deriving the verdict yourself. -->
**Reading the `design-conflict` token.** Every item under `sdd-apply`'s conditional
`### Rejected Proposals Checked` — present whenever the dispatch prompt forwarded ≥1 rejected proposal
whose task this run reached — carries a `· record: <domain>/<slug>` line and a `· design-conflict:
none | conflicts` line. Classify on the token **alone**, mirroring "Reading the `blocking-test` token"
below:

| Token | Asserts | Orchestrator |
|---|---|---|
| `none` | the design's approach and the rejected proposal describe the same implementation for that point | nothing — proceed, no gate, no mention |
| `conflicts` | the design fixes the point twice, incompatibly | the return MUST be `blocked`; present both versions (the design's approach, the rejected proposal) and the options. A `conflicts` item on a non-`blocked` return is a contract violation — stop and surface it, do not route it through the ordinary Unresolved Decisions Guard flow |
| no item for a `record:` this run forwarded as rejected, whose governed task the run reached (not left untouched in `### Remaining Tasks`) | no pronouncement | stop; name the record left without one; the return is not treated as correct |

**The expected set is yours, not the return's.** You are the one who forwarded each rejection in the
dispatch prompt, so you hold the expected set of `record:` identities that owe a verdict — nothing in
the return declares it, and nothing needs to: compare it, literally, against the `· record:` lines the
run actually returned. A forwarded rejection with no matching item is the third row above, not silence
to read as `none`. This check stays entirely on your side — `validate-return.js` and `render-return.js`
are unchanged by this mechanism; they only prove that whatever item IS present carries both tokens
legally, never that every rejection you forwarded got one.

<!-- matecito-ai: `New Decisions` y `Open Questions` se solapaban — las dos recibían decisiones
     pendientes, el ejecutor terminaba duplicando contenido y el usuario confirmaba lo mismo dos
     veces (fatiga de confirmación, el fallo que este guard existe para evitar). Desde ahora el
     único buzón Tier 1 de `sdd-design` es `New Decisions`. -->
**Tier 2 — accomplished fact or information → surface, do not block.** The Tier-2 mailboxes are **exactly** the sections marked Tier 2 in that same canonical table (`_shared/sdd-phase-common.md`, **Section D.3**) — read the list there. This guard does NOT keep a parallel copy of it either: enumerating Tier 1 by reference and Tier 2 inline was the same duplication, one tier later. Plus `risks`, which is an envelope field (D.4) and therefore not in that table, and is Tier 2 for every phase. These report what already happened or what is merely worth knowing. `Open Questions` is **informative**: it carries what does NOT fix a decision — anything that does fix one belongs in `New Decisions`, the single Tier-1 mailbox of that phase. Show them verbatim in the between-phase summary — never compress them to "no news" — and call out a deviation explicitly when it touches something `sdd-verify` will check against the design, because the design artifact is now stale. Where a section's contract splits its items into `summary`/`rationale`, "verbatim" binds to the `summary` — the part the contract declares printable; the `rationale` still ships in the block and is reproduced verbatim on request, not printed by default.

**One gate per phase, walked per the shared presentation.** Collect every Tier-1 item the phase returned and walk it through `~/.claude/references/gate-presentation.md`: one index over the phase's Tier-1 items, then each one shown in turn — item by item is the default here, not a batch-first summary — with "confirm the rest" the only bulk shortcut, offered before the first item and at any item while walking. This guard states no batching mechanic of its own. Gate fatigue is still the failure mode this exists to avoid — a gate the user clicks through without reading is worse than no gate — and the index up front plus "confirm the rest" mid-walk are what covers it now. This matters most in a repo with no `.matecito-ai/edr/`, where every architectural choice lands under `New Decisions`.

<!-- matecito-ai: las secciones se emiten SIEMPRE, así que sin esta regla toda sección tiene "contenido"
     (la cadena `None…`) y el gate se abriría en cada fase — la gate fatigue que este guard dice evitar. -->
**Empty → silent.** A Tier-1 section whose body is only an empty sentinel — any line starting with `None`, with or without a trailing explanation (`None.`, `None — mapping was explicit.`, `None — every task links to spec or design.`) — counts as EMPTY, not as content. **The sentinel is also recognized in Spanish** — a line starting with `Ninguna`, `Ninguno` or `Nada`, same rule, same trailing-explanation tolerance. Phase bodies are written in English and the canonical sentinel is `None`, but this ecosystem converses in Spanish and executors drift into it; a `Ninguna.` read as content opens a gate over nothing, and a gate the user clicks through without reading is the failure mode this guard is trying to avoid. No Tier-1 content means no gate: dispatch the next phase without mentioning this guard at all.

<!-- matecito-ai: esta regla trataba como retorno roto TODA sección ausente, y hay reglas vigentes que
     ordenan omitir secciones enteras de forma legítima y condicional (conflictos de EDR con el store
     inactivo, el veredicto de UI cuando no aplica, la evidencia de TDD fuera de estricto): las tres
     se marcaban como error. Además exigía "pedile que re-emita", una obligación sin mecanismo — no
     existe comando ni status de re-emisión, y re-despachar re-ejecuta la fase entera (en `sdd-apply`
     un re-despacho está definido como batch de continuación, no como re-emisión). -->
<!-- matecito-ai: esta frase decía "las enumeradas en D.3", y era cierta hasta que D.3 sumó una fila
     CONDICIONAL (`## Decision Gaps`, sólo cuando el cambio materializó al menos un decision record).
     Leída al pie, convertía la ausencia normal de esa sección en retorno roto → un gate espurio en cada
     verify. El corte es la columna "Emitted", no la pertenencia a la tabla. -->
**Omitted → depends on whether the section was unconditional.** Only sections declared **unconditional** are broken when missing — in the mailbox table of `_shared/sdd-phase-common.md`, **Section D.3**, those are the rows marked `always`; that table also carries conditional rows, which follow the rule below. The phase's own return template (`~/.claude/references/phase-returns/<phase>/<phase>.md`) marks the same distinction for every section it declares, mailbox or not. A section that the phase's own skill declares **conditional**, and whose condition does not hold, is **legitimately absent** — not a broken return, no gate, no mention (e.g. `sdd-intake`'s `### Early guard (EDRs)` when the EDR store is inactive, the UI verdict when the UI check does not apply, the TDD evidence table outside Strict TDD). <!-- matecito-ai: the first example used to be `## EDR Conflicts`, which exists only in the design ARTIFACT and never in a return — precisely the artifact/return confusion these rules exist to close. The other two really are return sections. -->

When an **unconditional** section is missing, do NOT assume there was nothing and do NOT silently dispatch the next phase — but do not demand a "re-emission" either: no such command or status exists, and re-dispatching re-runs the whole phase (for `sdd-apply` a re-dispatch is defined as a **continuation batch**, not a re-emission). Handle it with what exists: treat the omission as unresolved Tier-1 content and open the same gate you would open for real content, naming the phase and the missing section, and let the user pick — proceed as if empty / re-run that phase / adjust. You own the channel; the decision is theirs, not a repair you improvise.

### Review Workload Guard (MANDATORY)
After `sdd-tasks` and before `sdd-apply`, inspect `Review Workload Forecast`. If chained PRs recommended / 400-line budget risk High / decision needed → apply cached `delivery_strategy` (`ask-on-risk` default: STOP and ask chained PRs vs `size:exception`). Automatic mode does not override this guard.

The decision this guard raises is presented through the shared walkthrough in
`~/.claude/references/gate-presentation.md`, anchored to `sdd/{change-name}/tasks` — no index or
bulk-action wording of its own.

### Validator Findings (presentation)
`development-decisions-validate` and `development-spec-validate` are consultative — they only report,
never write. Their findings are walked through the shared presentation in
`~/.claude/references/gate-presentation.md`: one index when there are two or more findings, then each
one shown through the fixed item template, anchored to the EDR or capability-spec file the finding is
about. Neither skill states an index or bulk-action wording of its own.

<!-- matecito-ai: development declares its OWN decision-capture mechanism, per the kernel's override
     clause (`~/.claude/matecito-ai.md` → "Decision-Gap Capture (mine gate)"). This is why: propose ·
     ratify once · materialize in apply, instead of the kernel's generic post-verify mine. Full
     mechanism — the proposal shape, the ratification gate per lane, the materialization contract, the
     INDEX writer, and `sdd-verify`'s two checks — lives once in the reference below; this guard only
     states that it is MANDATORY and points at it. `design`, which declares its own post-verify
     mechanism, is unaffected. -->
### In-Flow Decision Capture (MANDATORY)
`development` does NOT use the kernel's generic Decision-Gap Capture (mine gate) — this domain declares
its own mechanism instead (see the kernel's override clause). Every phase that reaches an architecture
decision proposes it in its own return, the lane's gate ratifies it exactly once, and `sdd-apply`
materializes it as an `Accepted` EDR in the same step that implements the governing code — no
post-verify mining pass. Full mechanism: `~/.claude/references/decision-capture/in-flow-capture.md`.
Read it before touching any of: `sdd-spec`'s or `sdd-design`'s `### New Decisions` mailbox,
`sdd-apply`'s materialization step, or `sdd-verify`'s `decision-gaps` group.

## Spec-Mine — development specifics
The kernel owns the generic Spec-Mine Trigger (brownfield, `flagSpecMine`-gated, Mode A only). In development the spec-mining executor is `development-spec-mine`; confirmed candidates are materialized as capability-specs with `Status: Inferred` under `.matecito-ai/development-specs/<type>/<capability>.md` (type ∈ `flow` | `rule` | `lifecycle` | `process`) and the `.matecito-ai/development-specs/INDEX.md` is updated **once at the end**; the specs live ONLY as `.md`, never recorded in Engram — same as EDRs and capability-specs generally.

**Asymmetry vs decision-mine (important):** an `Inferred` EDR is still enforced by `sdd-verify` (its EDR-compliance step does not filter by Status), but an `Inferred` capability-spec is **NOT** verified — `sdd-verify`'s durable-capability-spec check is scoped to `Status: Accepted`, so `Inferred` (like `Draft`) is skipped and is never a contract until a human ratifies it to `Accepted` (via `development-spec-bootstrap` update mode). This guardrail is what makes it safe to keep as-built-derived `Inferred` specs in the store: they are pending-ratification drafts, not the ratified intention.

## Components axis (repo-level)
A **component** is a surface the product's consumer recognizes (`api`, `cli`, `ui`) — not every package or internal folder. The set is declared **per-project only**, in `.matecito-ai/config.json`'s top-level `repo.components` (never inherited from a global config, never folded into `domainConfig`). **Gate:** no `repo.components` declared → the axis does not exist for any consumer — no header line anywhere, no validation finding. Concept, declaration shape, and gate: `~/.claude/references/repo-components/README.md`.

## Capability-specs — development specifics
The system's **behavior** (the WHAT) is captured as durable **capability-specs**: files under `.matecito-ai/development-specs/<type>/<capability>.md` (type ∈ `flow` | `rule` | `lifecycle` | `process`), versioned in git and **never recorded in Engram** — exactly like EDRs. Concept and templates in `~/.claude/references/spec/README.md` and `~/.claude/references/spec/templates/`.

**Exception to engram-only (explicit).** The engram-only rule forbids file-based *proposal stores* of flow artifacts (like `openspec/`); it does NOT forbid durable repo knowledge. Capability-specs are durable knowledge that governs and verifies code — categorically the same as EDRs — so they live as files. The pipeline artifacts (`proposal`/`spec`/`design`/`tasks`/`verify-report`) stay in Engram; only the **accumulated behavior** is materialized to files. Never write pipeline artifacts to the filesystem.

**Who touches them:** `development-spec-bootstrap` authors them upfront (interview by capability, by type); `sdd-archive` merges each change's delta into them (scenario-anchored, non-destructive); `sdd-spec`/`sdd-design`/`sdd-tasks`/`sdd-verify` read them as the behavior contract; `development-spec-validate` checks coherence across them. Presence-based gate: absent store → every reader skips silently.

**Components axis (opt-in, presence-based) — the per-capability projection.** A capability-spec MAY carry a `- **Components:** api, ui` header line naming which surface(s) of the repo the described behavior touches — this line is the axis's **whole surface** inside the store: nothing else changes, both `INDEX.md` levels are untouched. Why the store never splits by app even though this line names surfaces: `~/.claude/references/spec/README.md` → "Por qué el store nunca se parte por app". `development-spec-validate` only reports (value outside the set, a spec missing the line, or the store found outside the repo root); it never writes and never auto-consolidates.
