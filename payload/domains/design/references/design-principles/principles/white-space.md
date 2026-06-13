# White Space

## Intent

Treat empty space as an active design element, not leftover background. White space (also "negative space" — it needn't be white) separates groups, emphasizes focal points, sets rhythm, and gives the eye room to rest. The space *between* things is as much a design decision as the things themselves.

## Why it matters (the perceptual reason)

Empty space is how the visual system **delimits** content. Proximity grouping only works because gaps exist to define group boundaries — space is the medium that proximity speaks through. Space around an element also creates **isolation**, and isolation is one of the strongest emphasis cues there is: a single item in a sea of space dominates a composition regardless of its size.

There's a cognitive-load dimension too. Dense, crammed layouts force the eye to fight through clutter, raising the effort of every glance; generous spacing lowers that effort and reads as calm, premium, and confident. Cramped reads as cheap and anxious — not because of taste, but because the brain is working harder.

## The Problem

Treating empty space as wasted space:

- Filling every pixel "to use the room" — producing a dense, exhausting wall with no entry point.
- Uniform tight margins everywhere, so nothing is grouped and nothing breathes.
- A headline jammed against the edge and the body text, with no room to function as a focal point.
- Equal padding inside and between cards, so the card boundaries blur (a [proximity](proximity-gestalt.md) failure caused by space).

## How to apply

- **Use a spacing scale.** Define spacing tokens (e.g. 4 / 8 / 16 / 24 / 32 / 48) and compose with them, so gaps are intentional and repeatable ([repetition](repetition.md)). Avoid arbitrary one-off values.
- **Distinguish macro and micro white space.** Macro = margins, gutters, gaps between sections (structure and grouping). Micro = line height, letter spacing, padding inside a button (legibility and comfort). Both matter.
- **Let space do the grouping** before reaching for borders and boxes — more space between groups than within them ([proximity](proximity-gestalt.md)).
- **Isolate to emphasize.** Give the primary element room; surrounding emptiness promotes it in the [hierarchy](hierarchy.md) without enlarging it.
- **Don't fear the void.** Resist the urge to fill — restraint is the point. Generous margins are a feature, not a gap to be plugged.

## Common violations to watch for

- **Horror vacui** ("fear of empty space") — the compulsion to fill every area, the single most common white-space failure.
- **Insufficient line height.** Body text set too tight (under ~1.4) is physically harder to read ([typography](typography.md)).
- **Uniform spacing.** Identical gaps everywhere defeat proximity and flatten structure.
- **Inconsistent padding.** Ad-hoc internal spacing that drifts from the scale, breaking rhythm ([repetition](repetition.md)).
- **Edge-crowding.** Content touching frame edges, with no breathing margin.

## Connection to the guards

- **`visual-accessibility`** — adequate spacing supports readability and touch-target separation (minimum tap sizes and gaps), and line-height/letter-spacing minimums are part of WCAG text-spacing criteria. Cramped layouts fail real low-vision and motor-impaired users.
- **`brand-consistency`** — the spacing scale is a locked system token; pieces with off-scale gaps drift from the brand's rhythm and trip the consistency audit.

## Related principles

- **[Proximity & Gestalt](proximity-gestalt.md)** — white space *is* the medium proximity works in.
- **[Hierarchy](hierarchy.md)** — isolation by space is a primary emphasis tool.
- **[Typography](typography.md)** — line height and letter spacing are micro white space.
- **[Repetition](repetition.md)** — the spacing scale is a reused system decision.
