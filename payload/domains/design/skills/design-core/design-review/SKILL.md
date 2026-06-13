---
name: design-review
description: >
  Reviews a design (a Figma frame, a screenshot, or a coded component) with UX, hierarchy,
  consistency, accessibility, and copy criteria, and returns prioritized findings. USE THIS SKILL
  whenever the user says "review this design", "what do you think of this screen", "critique this
  frame", "is this good UX", "give me feedback on this piece", or shares a frame/screenshot/
  component asking for a structured design review.
---

# Design review

Performs a structured design review of whatever the user indicates: a Figma frame (via the
`figma` MCP, read-only), an attached screenshot, or an already-coded component.

## Steps

1. Identify what is being reviewed and bring in the context (Figma screenshot, file, or code).
2. Evaluate these dimensions, one by one:
   - **Visual hierarchy**: what is seen first, contrast, typographic weight, focus.
   - **Consistency**: uses the brand system tokens; flag any off-system value. Lean on
     `consistency-audit` for a deep check against the brand guide and DDRs.
   - **Usability**: clarity of actions, affordances, states (empty, loading, error).
   - **Accessibility (WCAG AA)**: color contrast, target sizes, reading order, focus. Lean on
     `visual-accessibility` for measured ratios.
   - **UI copy**: clear microcopy, sentence case, useful error messages.
3. Return the findings **prioritized** (blocking / important / minor), each with the "what" and
   the concrete "how to fix it". No "could be improved" without the action.
4. Don't rewrite the whole design: point out what matters and let the person decide.

If the team has consistency or accessibility skills installed, lean on them (`consistency-audit`,
`visual-accessibility`).

## Mentor mode

Explain the why behind each finding in 1-2 lines — the design principle behind it (contrast,
hierarchy, rhythm, etc.), not just the what. Derive to the `explain-concept` skill for unfamiliar
concepts. (Full rule in the domain CLAUDE.md.)
