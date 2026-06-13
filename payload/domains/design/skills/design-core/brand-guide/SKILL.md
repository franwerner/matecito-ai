---
name: brand-guide
description: >
  Builds and maintains the team's brand guide as a design DELIVERABLE: palette with usage,
  typographic system, logo usage (do / don't), iconography, photography/illustration, and
  examples. It is the visual source of truth the other skills consult. USE THIS SKILL whenever
  the user mentions "brand guide", "brand manual", "branding", "visual identity", "palette",
  "style guide", wants to document or export the brand, or asks for something to be "on brand".
  Also when another design task needs to know a color, typeface, or brand rule: consult it
  before inventing anything.
---

# Brand guide

Keeps the team's visual identity in one place, in design language (not code), so everything
designed comes out consistent and can be delivered as a document.

## Source of truth

The guide lives as a project document/file (editable markdown + PDF export when needed). That
document is the truth: pieces align to it. Brand decisions captured in it are recorded as Design
Decision Records (DDRs) under `.matecito-ai/ddr/` — the palette, type system, and logo rules are
each a decision that later phases respect and verify.

If the guide does not exist yet, offer to build it (from `brand-from-references` or a brief).

## What it contains

1. **Palette** — each color with its role (primary, secondary, background, text, states) and hex.
   State valid uses and combinations, not just a swatch board.
2. **Typography** — families, scale (display/title/body/caption), weights, line-height, and
   recommended pairing.
3. **Logo** — versions, clear space, minimum size, and a clear "what NOT to do" section (distort,
   recolor, on low-contrast backgrounds).
4. **Shapes & iconography** — radii, stroke weight, icon style.
5. **Photography / illustration** — visual direction, treatment, examples.
6. **Visual tone** — 3-5 adjectives summarizing the personality.

## Flow

1. If a guide already exists, read it and work on top of it. If not, build it section by section,
   validating with the person.
2. When another skill or the user asks for a color/typeface/rule, answer FROM the guide.
3. If a piece appears with a value outside the guide, flag it: either it becomes a new rule (a new
   DDR) or it gets corrected.
4. Brand decisions are captured as DDRs in `.matecito-ai/ddr/`. When a decision is made or
   changed, record/update the DDR so `consistency-audit` and the verify phase can check against it.
5. To deliver: offer to export the guide as a document (use the PDF/DOCX creation skill).

## Rules

- No "a nice blue": always hex and a semantic name.
- Brand changes = edit the guide (and its DDR) first, then propagate to the pieces.
- Sentence case and consistent names throughout the guide.
- The brand is the person's decision: propose, don't impose.

## Mentor mode

Explain the why behind each decision in 1-2 lines — the design principle, not just the what.
Derive to the `explain-concept` skill for unfamiliar concepts. (Full rule in the domain CLAUDE.md.)
