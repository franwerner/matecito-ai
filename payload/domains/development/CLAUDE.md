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
| Phase commands | `/sdd-*` |
| Alignment artifact | `spec` |
| Decision record | `EDR`, stored in `.matecito-ai/edr/` |
| Decision-record concept reference | `~/.claude/references/edr/README.md` |
| Canonical catalog | `design-patterns` at `~/.claude/references/design-patterns/` (`Applied pattern: X` → `patterns/<x>.md`) |
| Decision-mining executor | `development-decisions-mine` |
| Decision-capture skill | `development-decisions-bootstrap` |
| Exploration index | CodeGraph (`mcp__codegraph__*`), active when `.codegraph/` exists |
| Guards | `strict-tdd`, `review-workload` |
| Engram topic-key namespace | `sdd-init/{project}` · `sdd/{change-name}/{intake,explore,proposal,spec,design,tasks,apply-progress,verify-report,archive-report,state}` |
| Init topic key | `sdd-init/{project}` |

## Language
Code (variables, functions, classes, constants): English. Comments follow the `code-comments` skill (which also fixes their language as English).

## Contract & definition shapes — never inferred
When you are about to create or modify a **contract or definition** — a domain entity, a database model / migration / schema, an API request/response (DTO), a public/exported type, interface, or enum, an event payload, or a config schema — you MUST NOT infer its shape. Both **which properties it has** and **each property's type** are decisions the human owns; do not invent "obvious" fields or "most likely" types (an `id`, an `email`, a `createdAt`; `string` vs `uuid`; float vs integer minor units — these are decisions, not defaults).

- **Artifact-pinned → execute.** If the shape is already fixed by an upstream artifact (spec, design, an EDR modeling policy, or the user's explicit request), that IS the mandate — implement it, do not re-ask. Only the un-pinned parts are open.
- **Unspecified → ask, per whole contract.** Propose the FULL contract (all fields + their types) as one reviewable unit — never field-by-field. With several unspecified contracts, default to one at a time (they often depend on each other); tell the user how many there are and offer "one-by-one or all-at-once" so they set the pace.
- **Where the answer lives.** The concrete shape (field names + types) belongs in the **code** (or the `design` artifact that materializes it) — NEVER copied into an EDR as a typed struct (that is a code calco; EDR reasoning stays conceptual). Only a **cross-cutting modeling policy** hidden in the answer ("identifiers are UUIDs", "money as integer minor units", "status as enums, not magic strings") may be captured as an EDR, expressed conceptually — and once captured it pins that part, so you stop re-asking it (per artifact-pinned). Offer to capture such a policy; never force it.
- **Scope.** Targets shapes that persist, cross a boundary, or are public. A transient internal struct used within a single function is execution detail, not a contract — no need to ask.

This is a specialization of the kernel's "Open question = blocked, not permission" for the high-stakes case of contracts, where inference is most tempting and most consequential (it propagates to DB, API, and tests).

## CodeGraph
Code exploration prefers CodeGraph when `.codegraph/` exists (structural questions); grep for literal text or non-indexed files. The SDD fork assumes the `mcp__codegraph__*` prefix — verify the tool names (`codegraph_search`, `codegraph_explore`, `codegraph_impact`, etc.) match the real MCP registration.

## Architecture diagrams (drawio)
Diagramming here is two complementary pieces: the **`drawio` skill** owns the *vocabulary* — how to build the diagram XML (shapes, branded/AI icons via `shapesearch`/`aiicons`, style presets, layout, diagram-type templates) — and the **`mcp__drawio__*` MCP** owns the *live render* — it renders the skill's `<mxGraphModel>` as an ephemeral preview; the skill itself never writes files. This rule is the **single source of truth for when to draw**. Diagrams are generated **on demand, never automatically**, and only when the change has structural complexity worth visualizing. **Diagram inference test — generate when** the change introduces or rewires ≥3-4 components with relationships, data flow crosses boundaries (layers/services/new modules), there is a non-trivial process with branches or states, or the task is to understand existing code spread across many files (CodeGraph can feed the graph) — **plus** capturing the shape of an architectural decision (EDR). **Do NOT generate** for a small fix, rename, config tweak, single-file/single-unit change, or linear logic — there prose or a snippet is clearer. **Model — offer-and-confirm, never unilateral.** **Decide vs generate (timing):** the structure does not exist yet at intake, so `sdd-intake` only *decides* — it sets `diagram: needed | not-needed` in the brief per this test, and the user confirms it at the **INTAKE GATE** (that gate IS the confirmation for the in-flow case; no re-ask later). **Generation is EPHEMERAL — always a live preview, NEVER a file in the project (zero `.drawio` artifacts in the repo).** The diagram is rendered by the **main thread** with a live preview (`mcp__drawio__*` — `start_session` reports the preview URL; the port is assigned dynamically, not fixed); nothing is exported or persisted. When the flag is `needed`, the main thread offers to render it live at the design step — the **headless `sdd-design` sub-agent does NOT generate or export diagrams**; it only notes that a live diagram is recommended. Same for a `direct` lane / outside the flow. Apply this same test before offering.

## Debugger MCP (mcp-debugger)
The debugger MCP (`mcp__debugger__*`, DAP step-through via `@debugmcp/mcp-debugger`) is **on-demand only** — it is NEVER invoked automatically. Its primary home is `sdd-apply`: when a runtime defect is encountered and the per-language debug toolchain is available (detected by `sdd-init` and cached in `sdd/{project}/testing-capabilities`), `sdd-apply` MAY diagnose the root cause AND apply a fix in the same context. In `sdd-verify`, the debugger is **diagnosis-only**: it MAY be used to understand why a test or scenario fails, but MUST NOT apply fixes there — any fix found belongs in a subsequent `sdd-apply` invocation. When the per-language debug toolchain is absent (`debugger.available = ❌` in testing-capabilities), both phases skip debugger usage silently — no error, no warning, no section. **For the full usage guide** — preflight (adapter vs. toolchain binary distinction), per-language install helper, and the debug loop — read the **`debugger` skill** (`payload/domains/development/skills/matecito-ai/debugger/SKILL.md`).

## SDD Flow

```
sdd-intake → sdd-explore → sdd-propose → sdd-spec → sdd-design → sdd-tasks → sdd-apply → sdd-verify → sdd-archive
                                              ^
                                           (design reads EDRs)
```

### Commands

- `/sdd-init` → initialize SDD context; detects stack, bootstraps persistence
- `/sdd-intake <request>` → structure a raw request into an Intake Brief (entry phase)
- `/sdd-explore <topic>` → investigate; reads codebase, compares approaches
- `/sdd-apply [change]` → implement tasks in batches
- `/sdd-verify [change]` → validate against specs + EDRs
- `/sdd-archive [change]` → close a change, persist final state in Engram
- `/sdd-onboard` → guided end-to-end walkthrough

Meta-commands (orchestrator handles them): `/sdd-new <change>`, `/sdd-continue [change]`, `/sdd-ff <name>`.

### SDD Phase Read/Write

| Phase | Reads | Writes |
| --- | --- | --- |
| `sdd-intake` | raw request | `intake` |
| `sdd-explore` | intake (brief) | `explore` |
| `sdd-propose` | exploration (optional) | `proposal` |
| `sdd-spec` | proposal (required) + **durable capability-spec** (for Modified Capabilities) | `spec` |
| `sdd-design` | proposal + **EDRs** + **durable capability-specs** (required) | `design` |
| `sdd-tasks` | spec + design + **durable capability-specs touched** (required) | `tasks` |
| `sdd-apply` | tasks + spec + design + apply-progress | `apply-progress` |
| `sdd-verify` | spec + tasks + apply-progress + **EDRs touched** + **capability-specs touched** | `verify-report` |
| `sdd-archive` | all artifacts | `archive-report` + **durable capability-specs (merge)** |

The "Reads" column lists the **full-lane** ideal. In `reduced`/`custom` lanes some upstream phases don't run, so each phase reads the **nearest available upstream**: `sdd-spec` falls back to the intake brief when there is no proposal; `sdd-apply` treats `spec` as the floor and skips `tasks`/`design` when absent. The **durable capability-specs** are read only when `.matecito-ai/development-specs/` exists; absent → skip silently (same presence-based gate as EDRs).

## Guards

### Strict TDD (resolution + forwarding)
Same precedence as model resolution — per-project `domainConfig.development.strictTdd` → global `domainConfig.development.strictTdd` → `false` (pre-M7 flat top-level `strictTdd` is auto-migrated into `domainConfig.development` on read). Resolve once per session, cache. If effective `strictTdd` is true, add to the `sdd-apply` / `sdd-verify` prompt: "STRICT TDD MODE IS ACTIVE. Test runner: {test_command}. Follow strict-tdd.md." The `{test_command}` comes from `sdd/{project}/testing-capabilities` in Engram.

### Review Workload Guard (MANDATORY)
After `sdd-tasks` and before `sdd-apply`, inspect `Review Workload Forecast`. If chained PRs recommended / 400-line budget risk High / decision needed → apply cached `delivery_strategy` (`ask-on-risk` default: STOP and ask chained PRs vs `size:exception`). Automatic mode does not override this guard.

## Decision-Gap Capture — development specifics
The kernel owns the generic mine gate. In development the mining executor is `development-decisions-mine`; confirmed candidates are materialized as `[Inferred]` `.md` EDRs and the `.matecito-ai/edr/INDEX.md` is updated **once at the end**; the EDRs live ONLY as `.md`, never recorded in Engram.

## Spec-Mine — development specifics
The kernel owns the generic Spec-Mine Trigger (brownfield, `flagSpecMine`-gated, Mode A only). In development the spec-mining executor is `development-spec-mine`; confirmed candidates are materialized as capability-specs with `Status: Inferred` under `.matecito-ai/development-specs/<type>/<capability>.md` (type ∈ `flow` | `rule` | `lifecycle` | `process`) and the `.matecito-ai/development-specs/INDEX.md` is updated **once at the end**; the specs live ONLY as `.md`, never recorded in Engram — same as EDRs and capability-specs generally.

**Asymmetry vs decision-mine (important):** an `Inferred` EDR is still enforced by `sdd-verify` (its EDR-compliance step does not filter by Status), but an `Inferred` capability-spec is **NOT** verified — `sdd-verify`'s durable-capability-spec check is scoped to `Status: Accepted`, so `Inferred` (like `Draft`) is skipped and is never a contract until a human ratifies it to `Accepted` (via `development-spec-bootstrap` update mode). This guardrail is what makes it safe to keep as-built-derived `Inferred` specs in the store: they are pending-ratification drafts, not the ratified intention.

## Capability-specs — development specifics
The system's **behavior** (the WHAT) is captured as durable **capability-specs**: files under `.matecito-ai/development-specs/<type>/<capability>.md` (type ∈ `flow` | `rule` | `lifecycle` | `process`), versioned in git and **never recorded in Engram** — exactly like EDRs. Concept and templates in `~/.claude/references/spec/README.md` and `~/.claude/references/spec/templates/`.

**Exception to engram-only (explicit).** The engram-only rule forbids file-based *proposal stores* of flow artifacts (like `openspec/`); it does NOT forbid durable repo knowledge. Capability-specs are durable knowledge that governs and verifies code — categorically the same as EDRs — so they live as files. The pipeline artifacts (`proposal`/`spec`/`design`/`tasks`/`verify-report`) stay in Engram; only the **accumulated behavior** is materialized to files. Never write pipeline artifacts to the filesystem.

**Who touches them:** `development-spec-bootstrap` authors them upfront (interview by capability, by type); `sdd-archive` merges each change's delta into them (scenario-anchored, non-destructive); `sdd-spec`/`sdd-design`/`sdd-tasks`/`sdd-verify` read them as the behavior contract; `development-spec-validate` checks coherence across them. Presence-based gate: absent store → every reader skips silently.
