---
name: explain-concept
description: >
  Explains design concepts clearly and with a visual example, so the person learns along the way:
  visual hierarchy, contrast, grid and spacing (8pt), type scale, kerning/leading, color theory,
  Gestalt, balance, affordances, etc. USE THIS SKILL whenever the user asks "what is X", "why does
  X matter", "I don't get X", "explain X", "difference between X and Y" about a design topic, or
  when another skill touches a concept worth clarifying. It is the learning engine of the domain.
---

# Explain concept

Turns a design question into a clear, short mini-lesson with an example. This skill is the target
of the domain's cross-cutting mentor mode: when any phase or skill surfaces a concept the person
may not know, it derives here. The golden rule: **show, don't just describe.**

## How to explain

1. **A one-sentence definition**, in plain language (no needless jargon).
2. **Why it matters**: what problem it solves or what happens if it's ignored.
3. **A visual example** (the key step): see "How to show the example" by surface.
4. **How to apply it** in the person's real work (ideally tied to what they're doing).
5. **Optional**: a counterexample (good vs bad) if it helps fix the idea.

Keep it brief. It's a lesson in passing, not a chapter. If the person wants more, go deeper.

## How to show the example (by surface)

- **Claude.ai / Cowork**: generate an inline visual (diagram or illustration) that shows the
  concept — a type scale, two contrast chips with their ratio, an 8pt grid, etc.
- **Claude Code (terminal)**: no inline render; generate a short .svg or .html file the person can
  open, or a clear text/ASCII example if that's enough. Say where the file landed.

Adapt the example to the context: if the person is working on a concrete piece, use that piece as
the example instead of a generic one.

## Source of rationale

Ground explanations in the `design-principles` catalog at `~/.claude/references/design-principles/`
rather than improvising. When a principle has a canonical entry there, cite it so the lesson is
consistent with the rest of the domain.

## Rules

- Always concrete: no "use good contrast"; show a pair that passes and one that doesn't, with
  numbers.
- Sentence case, no unexplained jargon.
- Connect the concept to what the person is doing, so it doesn't stay abstract.
- Don't overrun: the strength is in the example, not the text.
