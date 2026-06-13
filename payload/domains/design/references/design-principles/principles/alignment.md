# Alignment

## Intent

Place every element on a shared, invisible line so the composition reads as intentional and ordered. Alignment connects elements that aren't physically near each other, and its absence is one of the fastest ways a layout starts to feel amateurish — even when nothing else is wrong.

## Why it matters (the perceptual reason)

The eye is exquisitely sensitive to **edges**. It detects a misaligned edge as easily as a misspelled word, and reads it as noise — a small, repeated "something is wrong" signal that accumulates into a sense of carelessness. Conversely, shared edges create **implied lines** (Gestalt continuity): the eye follows them, so aligned elements feel connected and the whole gains an invisible skeleton that organizes attention.

Alignment is also what makes a design feel *trustworthy*. Precision in the small things (every edge on a line) signals precision in the large things; sloppy edges signal the opposite, and viewers transfer that judgment to the content.

## The Problem

Elements placed by eye, each on its own line:

- A form where every field starts at a slightly different x-position — a ragged left edge that reads as broken.
- Centered text mixed with left-aligned text on the same card, with no rationale.
- Captions, headings, and images each indented differently — no underlying structure.
- "Close enough" placements that are off by a few pixels, just enough to register as wrong without the viewer knowing why.

## How to apply

- **Commit to a grid.** A column grid (and a baseline grid for type) gives every element a line to snap to. This is the structural form of alignment and the foundation of a coherent layout.
- **Prefer strong edges to centering.** Left-aligned (in LTR) text and elements share a crisp, scannable edge; centering creates a ragged left edge that's harder to follow for long content. Center sparingly, for short, symmetrical, display-level content.
- **Pick one alignment per group and hold it.** Don't mix left, center, and right within a single logical block without reason.
- **Align across groups, not just within them.** The real power is connecting a heading, an image, and a button three sections apart by a shared edge.
- **Use the tool's smart guides / layout grids** (e.g. Figma) so alignment is enforced, not eyeballed.

## Common violations to watch for

- **Pixel-off placements.** Elements *almost* aligned — worse than obviously misaligned because it looks accidental.
- **Centering everything.** Center alignment as a default produces ragged, hard-to-scan edges and a weak structure.
- **Too many alignment lines.** Every element on its own edge means no edges are shared — the grid dissolves.
- **Ignoring optical alignment.** Mathematically-aligned shapes that *look* off (e.g. a triangle "play" icon, or punctuation hanging) need optical adjustment, not metric alignment.

## Connection to the guards

- **`visual-accessibility`** — a consistent alignment grid supports predictable reading order and scannability, which benefits low-vision users and anyone relying on a logical, linear flow.
- **`brand-consistency`** — the grid and its alignment rules are part of the locked design system; a piece that abandons the shared grid breaks the brand's structural signature even if its colors and type are correct.

## Related principles

- **[Proximity & Gestalt](proximity-gestalt.md)** — alignment is Gestalt continuity; shared edges reinforce grouping.
- **[Repetition](repetition.md)** — the grid is a repeated structural decision applied across every page.
- **[White Space](white-space.md)** — margins and gutters are defined by the same grid that drives alignment.
