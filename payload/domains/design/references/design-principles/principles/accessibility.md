# Accessibility

## Intent

Design so that everyone — including people with low vision, color-vision deficiencies, motor impairments, or cognitive differences — can perceive, understand, and use the work. Accessibility is not a feature added at the end; it is a constraint that shapes every other principle. A design that a meaningful fraction of its audience cannot read or operate has failed, however beautiful it looks.

## Why it matters (the perceptual and ethical reason)

Human perception varies far more than a designer's own eyes suggest. Roughly **1 in 12 men** has a color-vision deficiency; a large share of any audience has reduced contrast sensitivity (which worsens with age and in bright environments); many use assistive technology, larger text, or operate by keyboard or touch with limited precision. Designing only for ideal vision on a calibrated monitor silently excludes these people.

There's also a hard floor: legibility and contrast are **physiological**. Below a measurable contrast ratio, text is literally unresolvable for some eyes — no amount of taste overrides the math. And in many contexts (public sector, large products) accessibility is a **legal requirement** (WCAG, ADA, EN 301 549), not a nicety. Accessible design is also simply better design: it forces clarity, hierarchy, and restraint that benefit *all* users.

## The Problem

Designing for the designer's own eyes only:

- Light-gray text that looks refined on a Retina display and is invisible in sunlight or to a 60-year-old reader.
- Status conveyed only by red/green, so color-blind users can't tell error from success.
- Tap targets too small or too close, so anyone without fine motor control mis-taps.
- A logical-looking layout whose underlying order (the reading/DOM order) scrambles what a screen reader announces.
- No visible focus state, so keyboard users can't tell where they are.

## How to apply

Treat the **WCAG** success criteria as the baseline, targeting at least **AA**:

- **Contrast.** Text ≥ 4.5:1 (≥ 3:1 for large text ≥ 24px / 19px bold); meaningful UI and graphical elements ≥ 3:1 ([contrast](contrast.md), [color](color.md)).
- **Not by color alone.** Pair every color-coded meaning with an icon, label, shape, or text so it survives color-blindness and grayscale ([color](color.md)).
- **Text size & spacing.** Adequate base size, line height, and letter spacing; layouts that survive the user enlarging text up to 200% ([typography](typography.md), [white-space](white-space.md)).
- **Targets & spacing.** Touch targets ≥ ~44px with sufficient separation for motor-impaired users.
- **Reading order.** The visual order must match the logical/DOM order that assistive tech reads ([hierarchy](hierarchy.md), [alignment](alignment.md)).
- **Focus & state.** Clear, consistent visible focus, and distinguishable disabled/error/active states ([repetition](repetition.md)).
- **Test with real conditions.** Grayscale check, contrast checker, color-blindness simulation, and at increased zoom — early, not as a final audit.

## Common violations to watch for

- **Failing contrast** — the single most common accessibility defect, and the most measurable.
- **Color-only information** — status, links, categories, chart series distinguished by hue alone.
- **Missing/weak focus indicators** — keyboard users left lost.
- **Tiny or crowded targets** — interactive elements too small or too close.
- **Mismatched reading order** — visually correct but logically scrambled for screen readers.
- **Designing on ideal hardware only** — never checking how the work behaves for non-ideal vision or displays.

## Connection to the guards

- **`visual-accessibility`** — this principle *is* the guard's specification. `design-verify` runs `visual-accessibility` against the real colors, sizes, and structure in the Figma file: it measures contrast ratios against WCAG AA, checks minimum sizes and spacing, verifies meaning isn't color-only, and checks the reading order. Anything below AA is flagged. The other principles supply the rationale; this one supplies the pass/fail thresholds.
- **`brand-consistency`** — accessibility decisions (the accessible palette, the focus-state treatment, minimum sizes) become locked system tokens, so accessibility is repeated consistently rather than re-litigated per piece. A DDR that records "this is our accessible palette / focus style" is enforced here.

## Related principles

- **[Contrast](contrast.md)** & **[Color](color.md)** — supply the contrast-ratio and not-by-color-alone requirements.
- **[Typography](typography.md)** — supplies size and text-spacing minimums.
- **[Hierarchy](hierarchy.md)** & **[Alignment](alignment.md)** — provide the predictable reading order assistive tech depends on.
- **[White Space](white-space.md)** — target separation and text spacing.
- **[Repetition](repetition.md)** — consistent, learnable states across the system.
