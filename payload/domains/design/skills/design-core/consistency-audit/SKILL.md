---
name: consistency-audit
description: >
  Reviews a Figma file (via the figma MCP) and detects everything that drifted out of the system:
  colors outside the brand palette, multiple grays that should be one, arbitrary spacings,
  inconsistent typography, and duplicated components that should be the same. USE THIS SKILL
  whenever the user says "check consistency", "is everything on brand?", "audit this file", "is
  anything off-system?", "clean up inconsistencies", or shares a Figma file/frame to check it
  respects the brand guide.
---

# Consistency audit

Compares what is in the Figma file against the brand guide and flags each deviation. It is the
skill that most relies on the MCP: it needs to READ the real file, not a screenshot.

## Requirements

- `figma` MCP connected and authenticated (read-only).
- A reference `brand-guide` (palette, typography, spacing, components) and its DDRs under
  `.matecito-ai/ddr/`. The accepted brand DDRs are the contract the audit checks against. If no
  guide/DDRs exist, flag it: without a reference there is no "system" to audit against; offer to
  build it first.

## Flow

1. **Read the file** with the MCP tools: styles, variables, and the nodes of the indicated frame
   or page. If the scope is unclear, ask whether it's a frame, a page, or the whole file.
2. **Compare against the guide and the accepted DDRs** on these dimensions:
   - **Color**: hex not in the palette; near-identical variants (e.g. #2B2B2B vs #2C2C2C) that
     should be unified; color used outside its role.
   - **Typography**: families/weights/sizes outside the defined scale.
   - **Spacing**: values that break the grid (e.g. stray 7px, 13px, 19px).
   - **Components**: detached instances, or copies that should be a single component.
3. **Return prioritized findings**: blocking (breaks the brand) / important / minor. Each with
   where it is, which rule (or DDR) it violates, and how to unify it.
4. **Propose the fix**, not just the diagnosis: which token/style each case should map to.

## Rules

- Don't invent values: report the ones you actually read from the file.
- If a deviation should really be a new rule (not an error), suggest adding it to the guide and
  capturing it as a new DDR.
- Be specific: "3 different grays in the buttons" is useful; "there are inconsistencies" is not.
- A piece that contradicts an accepted DDR is a blocking finding.

## Mentor mode

Explain the why behind each finding in 1-2 lines — the design principle, not just the what.
Derive to the `explain-concept` skill for unfamiliar concepts. (Full rule in the domain CLAUDE.md;
cite the `design-principles` catalog at `~/.claude/references/design-principles/` for rationale.)
