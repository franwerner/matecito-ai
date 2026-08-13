<!-- matecito-ai: template canónico del retorno de sdd-spec.
     Existe porque el formato vivía inline en la skill y estaba incompleto: declaraba UNA sola forma
     —la del camino feliz— y no decía nada de cómo se devuelve un `blocked`, que es precisamente lo
     que esta fase produce cuando la derivación de capabilities es ambigua. Y emitía
     `### Derived capabilities (unconfirmed)` sin marcarla como buzón Tier 1, que es lo único que
     hace disparar el gate del orquestador: sin ese rótulo, un ejecutor la trata como un comentario
     más y la omite cuando "no hay nada que decir".
     Este archivo es LA fuente del formato. La skill y el agente lo referencian; no lo copian. -->

# Return template — `sdd-spec`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Specs Written` | always | the orchestrator, as context |
| `### Coverage` | always | the orchestrator, as context |
| `### Blocker` | only on `status: blocked` | the orchestrator: it puts the question to the user |
| `### Derived capabilities (unconfirmed)` | always | Unresolved Decisions Guard — **Tier 1** |
| `### New Decisions` | conditional — only when the lane running has no `design` add-on active | Unresolved Decisions Guard — **Tier 1** |
| `### Next Step` | always | the orchestrator, to route |

Titles are fixed. This phase declares no accepted variants of them.

`### Derived capabilities (unconfirmed)` is the **only unconditional** Tier-1 mailbox of this phase.
It is what makes the Unresolved Decisions Guard fire: a derived capability mapping is a guess about
which durable capability-spec this change belongs to, and at archive that mapping decides which file
the delta gets merged into. Emit it **always** — with the `None — mapping was explicit.` sentinel when
the upstream proposal carried its own Capabilities section — because a missing section and an empty
one mean opposite things to the guard.

**Split into summary/rationale.** Each item declares two parts, `summary` and `rationale`, in the
`derived_capabilities` JSON: `summary` is what the gate prints, `rationale` is the full reasoning —
always emitted into this block, never printed by default. Both are non-empty, single-line strings; a
missing one, or one with an embedded newline, fails the render naming the item and the part, and
nothing reaches stdout. `summary`'s register is fixed once in Section D.3 of `sdd-phase-common.md` —
not restated here. `summary` also carries a **250-character cap**, enforced by `render-return.js`.

Every item also carries the `· anchor:` token, declared first so it prints directly under the
summary — the legitimate forms and the not-yet-written-target rule are fixed once in Section D.3 of
`sdd-phase-common.md`, not restated here.

<!-- matecito-ai: in-flow decision capture (development-specifics). Full mechanism:
     ~/.claude/references/decision-capture/in-flow-capture.md — this section only fixes the RETURN shape. -->
`### New Decisions` is this phase's **conditional** Tier-1 mailbox — the ratification gate for an
architecture decision this phase finds, in a lane where `sdd-design` will NOT also run (so no later
mailbox would ever ratify it). Emit it, with content or with the `None.` sentinel, **only** when
`decisions_gate_here` is `true` — read from the intake brief's `### Triage` line (`Lane: ... —
add-ons: [...]`): `true` iff `design` is NOT one of the add-ons listed. When `design` IS active, omit
this section entirely — `sdd-design`'s own `### New Decisions` is the single gate, and emitting a
second one here would ask the user to ratify the same decision twice. Same title as `sdd-design`'s
plain variant, byte-identical, because it is the same mailbox concept surfacing one lane earlier — the
Unresolved Decisions Guard applies the identical rule to it. Each item carries `summary` +
`· rationale:` (same split as above) plus three tokens, in this order: `· anchor:` (the concrete
source this decision is about — a `<repo-path>[:line]` or `<engram-key>`, per Section D.3 of
`sdd-phase-common.md`), `· blocking-test:` (identical meaning and values to `sdd-design`'s — see
`sdd-design.md`, "The blocking-test token") and `· record: <domain>/<slug>` — the EDR identity the
proposal would occupy if ratified, free-form (no closed value set), still required (an item missing
any of these tokens fails `TOKEN-MISSING`, the strict reading of an omission).

Only two statuses have a shape here: `done` and `blocked`. This phase's skill does not designate
`needs-input`, and a spec is written whole or not at all, so `partial` does not arise.

## `status: done` — the specs were written

```markdown
## Specs Created

**Change**: {change-name}

### Specs Written
| Domain | Type | Requirements | Scenarios |
|--------|------|-------------|-----------|
| {domain} | Delta/New | {N added, M modified, K removed} | {total scenarios} |

### Coverage
- Happy paths: {covered/missing}
- Edge cases: {covered/missing}
- Error states: {covered/missing}
- UI scenarios: {N written | none — ui-test: not-needed}

### Derived capabilities (unconfirmed)
{The capability mappings you derived because the upstream artifact carried no Capabilities section
— an intake brief in a `reduced` lane, or a proposal in the older format. One item per mapping:

- {the capability you derived, as `New` or `Modified`, and its name}
  · anchor: {the intake brief's Engram key, or the durable capability-spec path for a Modified capability}
  · rationale: {one line: what you derived it from — the brief's Affected Areas, its structured Request}

These are NOT contract until the main thread confirms them.
If the mapping came explicit from the proposal's Capabilities section: "None — mapping was explicit."}

### New Decisions
{Emitted ONLY when `decisions_gate_here` is true (no `design` add-on in this lane) — omit the section
entirely, do not even print "None.", when `design` is active. When emitted: the architectural choices
this phase found that pass the blocking test (see `sdd-design.md`, "The blocking-test token" — same
test, same values). One item per decision:

- {the choice}: {what you chose} — {alternatives weighed, and why this one}
  · anchor: {the concrete source this decision is about — a `<repo-path>[:line]` or `<engram-key>`}
  · blocking-test: none
  · record: {domain}/{slug}
  · rationale: {one line: the full reasoning}

If genuinely none: "None."}

### Next Step
Ready for design (sdd-design). If design already exists, ready for tasks (sdd-tasks).
```

## `status: blocked` — the derivation is ambiguous, or a contract shape is not yours to pin

Same block, with three differences. Everything else is emitted as usual — a blocked spec still
carries whatever specs you did write.

```markdown
## Specs Created

**Change**: {change-name}

### Specs Written
{The same table, for the specs that DID get written. If the blocker stopped you before any spec
landed: "None — blocked before any spec was written." Whether what you wrote was persisted is
stated in the envelope's `artifacts` field, not here.}

### Coverage
{Of what you did write, same bullets — including `UI scenarios`. Or "not assessable — depends on the
resolution of the blocker."}

### Blocker
{The question, in one line, phrased so the user can answer it without reading the rest.}

**Options weighed**: {each one, with what it costs and what it buys. When the blocker is an
ambiguous capability derivation, each option IS one possible reading — name the capability it maps
to and the evidence in the upstream artifact that supports it. List every reading you found
plausible; the point of blocking is that you did not narrow them down.}

**Why this is not mine to settle**: {for an ambiguous derivation: the mapping is not a labelling
choice — it decides which durable capability-spec the delta merges into at archive, so picking the
wrong reading silently rewrites the behavior of a capability nobody asked to touch, and naming a
capability that does not exist invents behavior. For a contract or definition shape: name it as
what it is under "Contract & definition shapes — never inferred" (entity, DB model/migration/schema,
DTO, public type/interface/enum, event payload, config schema) and say whether it is unspecified or
pinned by something that conflicts or does not cover this case.}

**What unblocks it**: {the answer you need. For a contract shape, propose the FULL contract — every
field and its type — as ONE reviewable unit, never field by field.}

### Derived capabilities (unconfirmed)
{Any mapping you derived that holds regardless of how the blocker resolves, same shape as `done` —
each item still needs both `summary` and `· rationale:`. If the blocker IS the derivation:
"None — the derivation is the blocker; see above."}

### Next Step
None. This phase cannot continue until the blocker is resolved.
```

`### Blocker` is where the possible readings go. Not `risks` — Section D.4 forbids routing a
decision the user owns through that field — and not `### Derived capabilities (unconfirmed)`, which
is for a mapping you actually settled on and are flagging as unconfirmed. A blocked derivation was
never settled: filing the readings there would hand downstream a contract you do not have.

## Artifact vs return — do not confuse them

The spec **artifact** (persisted to Engram, format in the skill) is the delta itself: ADDED /
MODIFIED / REMOVED requirements with their Given/When/Then scenarios. The **return** is this block:
what the orchestrator needs in order to route and to gate. The orchestrator never reads the artifact.

<!-- matecito-ai: `ui-scenarios` es el caso más tentador de confundir las dos cosas. El bloque es
     INSUMO de sdd-verify, que lee el artefacto — pegarlo en el retorno lo dejaría donde nadie lo
     ejecuta, y además no hay decisión pendiente que gatear: el flag `ui-test` ya lo confirmó el
     usuario en el INTAKE GATE. Por eso el retorno lleva sólo el conteo, y dentro de `### Coverage`
     (sección que ya se emite siempre) en vez de abrir una sección nueva que sería un buzón fantasma. -->
The `ui-scenarios` block written in the skill's Step 4b follows that same split: it belongs to the
**artifact**, because `sdd-verify` executes it from there. The return carries only the `UI scenarios`
count on the `### Coverage` line — deliberately not a section of its own. It is not a Tier-1 mailbox
and opens no gate: the `ui-test` flag it derives from was already confirmed by the user at the INTAKE
GATE, so there is nothing left for the orchestrator to decide. When the brief says `not-needed`, the
bullet reads `none — ui-test: not-needed`; the bullet itself is never dropped, so a missing count is
a defect and not "no UI work".

<!-- matecito-ai: a third form used to live here — `not authored — ui-test: needed, see Blocker` — for the
     case where the flag said `needed` and the block could not be authored, because the old schema demanded
     a `url` and role+name targets that neither the brief nor a not-yet-existing UI could supply. That case
     no longer exists: this phase now authors only the BEHAVIORAL half (domain language, no routes, no
     locators), and the executable counterpart is `sdd-apply`'s job — it knows the real targets because it
     wrote them. Two forms again, and this time they cover reality. -->
Two forms, and they are exhaustive: a count, or `none — ui-test: not-needed`. **There is no "could not
author it" case** — the behavioral half needs no information that does not exist yet, so the flag saying
`needed` always resolves to a count. Do not misreport either way: `none — ui-test: not-needed` when the
flag said `needed` lies about the flag, and `0 written` lies about the block, and both make
`sdd-verify`'s UI gate close **silently** so nobody learns the verification will not run.

Note the asymmetry with `sdd-design`, where `New Decisions` exists in both the artifact and the
return. Here the Tier-1 mailbox lives **only in the return** — a delta spec has no section for it,
and no other copy of it exists anywhere. Dropping it does not degrade the gate; it deletes the
information.
