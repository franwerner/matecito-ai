# Contrast

## Intent

Make elements that are different *look* clearly different — in size, weight, color, shape, or space. Contrast is what separates the important from the incidental and what makes a composition legible at a glance. Weak contrast reads as muddy and timid; strong, intentional contrast reads as confident and clear.

## Why it matters (the perceptual reason)

The human visual system is a **difference detector**. The retina and early visual cortex respond to edges, luminance steps, and changes — not to absolute values. Attention is drawn to wherever the field changes most. Two elements that are *almost* the same create a worse experience than two that are obviously the same or obviously different: the near-match makes the brain ask "are these meant to match or not?" and that unresolved question is friction.

Contrast also underwrites legibility at the physiological level: text is readable only when its luminance differs enough from the background for the eye to resolve the letterforms. This is not subjective — it's measurable as a contrast ratio.

## The Problem

Timid, low-differentiation design:

- A heading set two points larger than body text — the difference is invisible, so it doesn't function as a heading.
- Light gray text on a white background — "elegant" until a real user in sunlight can't read it.
- A primary and secondary button styled almost identically — the viewer can't tell which action is recommended.
- A palette of five nearly-identical blues — nothing in the design has a clear job.

## How to apply

Contrast must be **decisive**: if two things differ, make the difference unmistakable.

- **Size contrast** — don't go 16px → 18px; go 16px → 32px. Use a type scale with real jumps ([typography](typography.md)).
- **Weight contrast** — pair a true bold with a regular, not two mediums.
- **Color/luminance contrast** — ensure text-to-background ratios meet WCAG AA (4.5:1 body, 3:1 large text); use it to make CTAs pop ([color](color.md)).
- **Shape contrast** — a pill button among rectangular cards; a circular avatar in a grid of squares.
- **Space contrast** — a generous gap against tight clusters makes the isolated element dominant ([white-space](white-space.md)).

Contrast is the *mechanism* by which you build [hierarchy](hierarchy.md): decide what matters most, then make it contrast hardest with everything else.

## Common violations to watch for

- **The "almost" trap.** Sizes, colors, or weights close enough to look like a mistake rather than a decision.
- **Contrast that fails measurement.** A pairing that looks fine on a designer's calibrated monitor but drops below 4.5:1 — a guard failure, not an opinion.
- **Contrast everywhere.** If every element high-contrasts with its neighbor, you've recreated a no-hierarchy flat field. Contrast must be *selective* to mean anything.
- **Color-only contrast.** Distinguishing items by hue alone fails for ~8% of men (color-blindness); add a second cue ([accessibility](accessibility.md)).

## Connection to the guards

- **`visual-accessibility`** — directly measures contrast ratios against WCAG AA/AAA for every text-on-background and meaningful UI pairing in the Figma file; flags anything below threshold. This is the most frequently triggered guard check, and contrast is its core metric.
- **`brand-consistency`** — confirms contrast is applied the *same* way across pieces (e.g. the primary CTA always uses the brand's high-contrast treatment, never a washed-out variant).

## Related principles

- **[Hierarchy](hierarchy.md)** — contrast is the tool; hierarchy is the goal.
- **[Color](color.md)** and **[Typography](typography.md)** — where contrast is operationalized and measured.
- **[Accessibility](accessibility.md)** — sets the non-negotiable minimum contrast levels.
