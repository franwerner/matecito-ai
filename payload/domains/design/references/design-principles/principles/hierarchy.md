# Visual Hierarchy

## Intent

Arrange a composition so the eye is led through it in a deliberate order — most important element first, supporting elements after — instead of letting the viewer choose where to look. Hierarchy is the difference between a layout that is *read* and one that is merely *seen*.

## Why it matters (the perceptual reason)

Before a viewer consciously reads anything, the visual system triages the field by **salience** — size, weight, color, isolation, position. It decides "what matters here" in roughly 50 milliseconds, pre-attentively. If the design provides no salience cues, the brain still picks a focal point — just not the one you intended. Hierarchy is the act of taking control of that involuntary first pass so the message lands in the order you meant.

A design with no hierarchy forces the viewer to do the sorting work themselves, which is cognitively expensive. Most viewers won't pay that cost — they leave.

## The Problem

Every element competing for attention at the same level:

- Three headlines the same size — none reads as *the* headline.
- The CTA button the same weight as the footer links — the primary action disappears.
- A wall of body text with no entry point — nowhere for the eye to land.

When everything is emphasized, nothing is. "Flat" compositions read as overwhelming *or* boring, and both reactions end engagement.

## How to apply

Establish hierarchy through deliberate, layered cues — usually combining several:

- **Size** — bigger reads as more important. The clearest, bluntest tool.
- **Weight** — bold pulls forward; light recedes.
- **Color & contrast** — a saturated or high-contrast element jumps ahead of muted ones (see [contrast](contrast.md)).
- **Position** — top and left (in LTR reading cultures) and the optical center carry weight; the eye starts there.
- **White space** — isolating an element makes it dominant regardless of its size (see [white-space](white-space.md)).
- **Density & order** — a clear primary → secondary → tertiary scale (e.g. an H1/H2/body type scale) encodes hierarchy structurally, not by eyeballing each instance.

Decide the **single primary focal point first**, then rank everything else against it. There should be exactly one level-one element per view in most cases.

## Common violations to watch for

- **Too many focal points.** Two elements both screaming for first place cancel each other out.
- **Hierarchy by one cue alone.** Relying only on color fails for color-blind users and in grayscale; reinforce with size or weight ([accessibility](accessibility.md)).
- **Inverted importance.** Decorative or legal text given more visual weight than the actual message.
- **No tertiary level.** Only "huge" and "tiny" with nothing between — the eye jumps and the structure feels crude.
- **Flat type scale.** Body, captions, and headings nearly the same size, so the document has no skim-able structure ([typography](typography.md)).

## Connection to the guards

- **`visual-accessibility`** — checks that hierarchy survives without color (does size/weight/position still rank the content?) and that the reading order in the file matches the intended visual order, which screen readers depend on.
- **`brand-consistency`** — verifies the *same* hierarchy rules (the locked type scale, the standard CTA treatment) are applied across every piece, so importance is signaled the same way everywhere.

## Related principles

- **[Contrast](contrast.md)** — the primary mechanism for *creating* hierarchy.
- **[White Space](white-space.md)** — isolation as an emphasis tool.
- **[Typography](typography.md)** — the type scale is hierarchy encoded into the text system.
