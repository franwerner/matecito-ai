---
name: figma-hygiene
description: >
  Reviews the health/organization of a Figma file (via the figma MCP): layer names, auto-layout
  usage, component structure and reuse, detached instances, and unused styles. USE THIS SKILL
  whenever the user says "tidy up the file", "this is a mess", "layer naming", "should I use
  components/auto-layout?", "clean up the Figma", or when a file grew and became hard to maintain.
---

# Figma hygiene

Keeps the file tidy and maintainable, especially as it grows (whole apps, many frames). It does
not change the design: it improves how it's built underneath.

## Requirements

- `figma` MCP connected (read-only) to read the file structure.

## What it reviews

1. **Naming**: layers with default names ("Frame 127", "Rectangle 8"); lack of a consistent
   convention; unnamed groups.
2. **Auto-layout**: hand-positioned elements that should use auto-layout to be responsive and
   easy to edit.
3. **Components**: repeated pieces that should be a component; detached instances; variants worth
   grouping.
4. **Styles / variables**: loose colors and text that should be styles/variables; styles defined
   but unused.
5. **Page structure**: page/section organization so another person can find their way.

## Output

A prioritized list of maintainability improvements, each with the problem, where it is, and the
concrete step to fix it. Distinguish "quick and high impact" from "nice to have".

## Rules

- Don't touch design decisions (that's for other skills); focus on construction.
- Prioritize what saves the most time later (auto-layout and components usually win).
- Report on what is actually in the file, not assumptions.

## Mentor mode

Explain the why behind each finding in 1-2 lines — why auto-layout helps, what a component
prevents, etc. — not just the what. Derive to the `explain-concept` skill for unfamiliar
concepts. (Full rule in the domain CLAUDE.md.)
