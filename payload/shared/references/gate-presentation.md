# Gate presentation — one walkthrough, one template

<!-- matecito-ai: single statement of how a gate puts ratifiable items to the user — the index, the
     item-by-item walkthrough, the bulk shortcut, and the fixed slots an item is shown through. Every
     gate that ratifies items cites this file; none of them restates any of it. Lives in `shared/`
     rather than a domain fragment because two of its three citing gates are kernel-general, not
     `development`-specific — see "Where this governs" below. -->

A gate hands the user everything a return carries at once today, and each item arrives without naming
what it is about — so the user can neither judge one item on its own nor ask for its detail without
scrolling back through the whole batch. This file is the single, shared definition of how a gate avoids
that: one index of what a return brought, then the items walked one at a time, each one shown through a
template with no room for narrative.

## The walkthrough

### One index, never accumulated

A gate opens with exactly **one index**, covering everything ratifiable that the return (or the mining
run) in front of it carries — whatever section each item came from. The index states the total count
and groups the items by the section that produced them. It carries nothing already decided at an
earlier gate: each gate's index is its own. The flow never grows a running list across phases — the
next return that opens a gate gets its own index, starting from zero.

Nothing ratifiable in the return → no index, no walkthrough, no mention of this mechanism at all. The
gate is silent, and the next step proceeds.

### Item by item is the default

After the index, items are presented **one at a time, in index order**. Each one needs an explicit
outcome — confirmed, adjusted, or rejected — before the next one is shown. Nothing downstream is
dispatched while an item is still undecided; an item still awaiting an outcome is what holds the flow,
not a batch summary. This is the default in every execution mode, including the unattended one: the
unattended mode skips the between-phase checkpoint, never this walkthrough.

When an item is adjusted rather than confirmed as offered, the **adjusted** content is what counts as
ratified downstream — never the originally offered text.

### "Confirm the rest" — available before the first item, and at any item

The one bulk shortcut this walkthrough offers confirms every item **not yet decided**. It is available
in exactly two moments: **before the first item is shown** (over the index itself) and **at any item
while walking** (from the second item onward, same offer). Invoking it closes the walkthrough
immediately: items already decided keep the outcome they were given: nothing already confirmed,
adjusted or rejected is revisited; everything still undecided is recorded as confirmed.

No other bulk action exists. A gate does not invent "accept all high-confidence", "per-type",
"per-domain" or any other grouping of its own — "confirm the rest" is the only shortcut, and it is
offered in the same two moments everywhere this walkthrough runs.

### Asking for detail does not answer the item

An item is presented with its summary and its anchor only (see "The fixed item template" below).
Asking for its full reasoning, or for the material its anchor points at, is always available and never
counts as an outcome:

- **The reasoning** — the item's `rationale`, already carried in the return's block — is reproduced
  verbatim on request.
- **The anchored source**, when it points at code — a `<repo-path>[:line]` anchor — is retrieved and
  shown on request, without pulling the whole file into the presenting agent's context: extract the
  relevant range to a scratch file (`sed -n 'A,Bp' <file> > <scratch>/snippet`) and show that; only the
  path and the line count travel into the conversation, never the file's full content. When the anchor
  is an `<engram-key>` instead, retrieve that observation the ordinary way.

Either way, after answering, the walkthrough resumes at the **same** item — still awaiting a decision.

### Resume

Once every item has an outcome — whether reached one by one or closed early via "confirm the rest" —
the walkthrough is done and the gate proceeds exactly as it does today: dispatching the next step,
routing on the decisions made, or stopping on whatever the gate's own rules say a rejection triggers.

## The fixed item template

Every gate that ratifies items presents them through this **one** template. Its slots are fixed and
nothing outside them reaches the presentation — not the item's full reasoning, not commentary around
it, not a gate's own framing:

```
Decision {n} of {total} — {anchor or section label}
{Anchor slot — the concrete source, per the anchor criterion in
 `~/.claude/skills/_shared/sdd-phase-common.md`, Section D.3}
{Summary slot — the item's `summary`, as authored, one line}
Ratify it? (yes / no / see detail)
```

Illustrative, filled in:

```
Decision 1 of 2 — structure/event-vocabulary-home
File: apps/api/.matecito-ai/edr/structure/event-vocabulary-home.md
What's being decided: the event vocabulary lives in the feature, not in shared/
Ratify it? (yes / no / see detail)
```

Three slots, always the same three: the item's **summary**, its **anchor**, and the **actions**
available on it (confirm / adjust / reject / see detail, plus "confirm the rest" at the moments named
above). No gate states a presentation of its own, and no gate states its own wording for the actions —
the words above are the shared wording, not a per-gate choice.

### Tier-2 items stay out of the walkthrough

An item that is informative rather than ratifiable — a Tier-2 mailbox item, per the domain fragment's
tier definitions — never enters the index and never blocks anything. It keeps appearing in the
between-phase summary exactly as it does today, and it still carries its anchor like any other item:
staying out of the walkthrough is about not gating on it, not about withholding where it comes from.

## Where this governs

This walkthrough and this template govern exactly three gates:

- **The pending-decisions gate** — the domain's guard that ratifies a phase's pending decisions before
  the next phase dispatches (in `development`, the Unresolved Decisions Guard's Tier-1 mailboxes).
- **The confirmed-scope gate** — the INTAKE GATE. Its ratifiable items are the recommended lane and
  each of the brief's decision flags; each one's anchor slot is filled from fields the brief already
  carries (its Engram key, the `### Classification` block, the flag's own label) — the brief is not
  required to declare a new formal section for this.
- **Both mining confirmation gates** — the gate that precedes materializing decision-mine candidates,
  and the gate that precedes materializing spec-mine candidates. Candidates are indexed and walked like
  any other item, anchored to the source they were mined from.

A gate outside this list states its own presentation as it always has; this file does not reach for
one it was not asked to govern.
