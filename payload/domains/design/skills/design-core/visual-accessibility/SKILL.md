---
name: visual-accessibility
description: >
  Checks the visual accessibility of a design (via the figma MCP or on a screenshot): WCAG color
  contrast on the real text/background pairs, minimum text and touch-target sizes, and reading
  order/hierarchy. USE THIS SKILL whenever the user says "is it accessible?", "check the
  contrast", "does it meet WCAG", "is it readable?", "button sizes", "accessibility", or wants to
  validate that a design can be used comfortably by everyone.
---

# Visual accessibility

Validates that the design is legible and usable for everyone. This is accessibility on the design
side (contrast, sizes, hierarchy), not code.

## Input

- Preferred: the Figma file via the `figma` MCP (read-only), to read the real colors and sizes.
- Alternative: a screenshot uploaded to the chat (contrast is estimated visually, less precise).

## What it checks

1. **Color contrast (WCAG)** on each real text/background pair:
   - body text: target AA = 4.5:1 (AAA = 7:1);
   - large text (>=24px or >=19px bold) and UI elements: AA = 3:1.
   Report the ratio and whether it passes, per pair.
2. **Sizes**: text too small for its role; small touch targets (comfortable reference ~44x44px)
   that are hard to tap.
3. **Hierarchy and reading order**: whether the visual path is clear and focus lands where it
   should.
4. **Don't rely on color alone**: states (error, success) distinguished only by color.

## Output

A per-item list with: what was checked, the measured value (e.g. ratio 3.1:1), pass / fail, and
the concrete fix (e.g. "darken the text to #5F5E5A to reach 4.6:1").

## Rules

- Give real numbers when reading from the file; if estimating from a screenshot, say so.
- Prioritize: what blocks reading first, the details after.
- Don't reduce accessibility to "contrast passes": also look at size, hierarchy, and color-only.

## Mentor mode

Explain the why behind each finding in 1-2 lines — why 4.5:1, what happens with low vision, etc. —
not just the what. Derive to the `explain-concept` skill for unfamiliar concepts. (Full rule in
the domain CLAUDE.md; cite the `design-principles` catalog at
`~/.claude/references/design-principles/` for rationale.)
