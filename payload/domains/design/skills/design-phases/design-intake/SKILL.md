---
name: design-intake
description: >
  Method for the design INTAKE phase. Use as the FIRST step when a raw design request (a brand,
  screen, flyer, prototype, or visual change in natural language) must be turned into a structured
  brief-intake: ask targeted questions, classify the work, triage whether the full design flow is
  needed, and catch DDR conflicts or undecided brand questions before exploration begins.
---

# Design intake — method

The recipe the `design-intake` executor follows. This phase turns a vague chat request into a clear
brief-intake artifact that `design-explore` (or, in trivial cases, direct implementation) consumes.
It reads DDRs ONLY to catch early blockers; it does NOT inspect the Figma file or explore references
in depth (that is `design-explore`'s job).

## Reads / writes (design Phase Read/Write)

- **Reads:** the raw user request.
- **Writes:** `design/{change-name}/intake` (the brief-intake).

## Steps

1. Receive the raw user request (natural language from the chat — a brand, screen, asset, prototype).
2. Ask 2-4 targeted intake questions to lock down what is ambiguous (audience, medium, constraints, tone).
3. Classify the work: type (brand / ui / asset / prototype), surfaces touched, rough size.
4. Triage: does this warrant the full design flow, or is it trivial enough to go direct?
5. Early guard (DDR activation gate): if `.matecito-ai/ddr/` is absent or empty, DDRs are inactive —
   skip this step silently (`status: done`, no DDR mention in the brief). Only when it exists with
   content, check it for conflicts or undecided brand/visual questions this request raises.
6. Produce the structured brief-intake artifact and return it.

When DDRs are active and the request conflicts with an Accepted DDR → `status: blocked`; when it
raises an undecided brand question no DDR covers → `status: needs-decision`. When `.matecito-ai/ddr/`
is absent or empty, never emit `blocked` / `needs-decision` for DDR reasons and never mention DDRs —
treat such questions as ordinary design decisions for later phases (design-explore / design-system).

## Mentor mode

Explain the WHY behind each triage call in 1-2 lines — the design principle at stake (e.g. why a
flyer is `reduced` but a rebrand is `full`). Cite the `design-principles` catalog rather than
improvising; when a concept the person may not know surfaces, derive to `explain-concept`. Do not
paraphrase the full mentor rule (it lives in the domain CLAUDE.md).

Do NOT inspect the Figma file or explore references in depth, and do NOT design or produce — your job
is to turn a vague chat request into a clear brief, and to stop early on a DDR conflict or an
undecided brand question when DDRs are active.
