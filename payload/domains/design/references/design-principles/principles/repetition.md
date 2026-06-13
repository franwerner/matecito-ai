# Repetition

## Intent

Reuse visual decisions — colors, type styles, spacing, button shapes, icon style, layout patterns — consistently across a design and across a whole system of pieces. Repetition is what turns a collection of individual screens or pages into one recognizable, coherent thing. It is the principle most directly responsible for "brand."

## Why it matters (the perceptual reason)

Repetition exploits **similarity** (a Gestalt law) and **familiarity** at scale. When a viewer sees the same treatment a second time, the brain recognizes it instead of re-parsing it — recognition is faster and cheaper than fresh perception. Across an interface, this means a user *learns* the system once (this is what a button looks like, this is how a heading reads) and then navigates the rest effortlessly. The reused vocabulary becomes invisible infrastructure.

Repetition also builds **trust**. Consistency reads as competence and intention; inconsistency reads as either error or as a signal that "this isn't the same thing," forcing the viewer to re-evaluate. A brand that looks the same everywhere feels reliable; one that drifts feels improvised.

## The Problem

Re-deciding the same thing every time:

- Three slightly different blues across three screens because each was picked independently.
- Buttons with 6px corners on one page and 10px on another.
- Five heading sizes that don't belong to any scale, each chosen by eye.
- The same concept (a "warning") styled differently in two places, so users don't recognize it as the same concept.

Each individual choice may be fine; together they tell the viewer "no system here," and the cognitive cost of relearning every screen falls back on them.

## How to apply

- **Build a system, then draw from it.** Lock tokens — a color palette, a type scale, a spacing scale, corner radii, elevation/shadow levels, icon style — and reuse them. This is exactly what the `design-system` phase produces.
- **Componentize.** Reusable components (in Figma, real components/variants) make repetition structural rather than disciplined: change once, propagate everywhere.
- **Repeat patterns, not just atoms.** Reuse layout patterns (card structure, section rhythm, list item shape), not only individual colors and fonts.
- **Make exceptions deliberate and rare.** Repetition needs occasional contrast to avoid monotony — but a break should be a *decision* (captured in a DDR), not drift.

## Common violations to watch for

- **Drift.** Near-duplicate values (#2B6CB0 vs #2C6DB2, 16px vs 15px) that accumulate as the work grows — the most common and most corrosive violation.
- **Reinventing components.** Building a new button instead of reusing the system's, so variants multiply.
- **Inconsistent meaning.** The same visual treatment used for different concepts, or the same concept shown differently, breaking the learned mapping.
- **Monotony (the over-application).** Repeating so rigidly that nothing stands out — repetition with no [contrast](contrast.md) erases [hierarchy](hierarchy.md).

## Connection to the guards

- **`brand-consistency`** — this is repetition's home guard. `design-verify` checks every produced piece against the locked system and the accepted DDRs, flagging any element that uses an off-system color, type style, spacing, or component. Repetition is the property being verified.
- **`visual-accessibility`** — consistent, repeated treatment of UI states (focus, error, disabled) means assistive-tech users and low-vision users can rely on a learned, predictable vocabulary instead of decoding each instance.

## Related principles

- **[Proximity & Gestalt](proximity-gestalt.md)** — repetition is the similarity law applied system-wide.
- **[Alignment](alignment.md)** — the grid is a repeated structural decision.
- **[Color](color.md)** & **[Typography](typography.md)** — the palette and type scale are the repeated tokens.
- **[Contrast](contrast.md)** — the deliberate exception that keeps repetition from becoming monotony.
