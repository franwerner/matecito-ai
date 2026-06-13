# Color

## Intent

Use a deliberate, limited, and meaningful palette. Color is the loudest signal in a composition — it carries emotion, encodes meaning, creates [contrast](contrast.md), and is the fastest-recognized brand asset. Because it's so powerful, undisciplined color does more damage than any other surface decision.

## Why it matters (the perceptual reason)

Color is processed **pre-attentively and emotionally**. A hue registers before shape or text, and it carries cultural and psychological associations (red = alert/passion, green = success/go, blue = trust/calm) that the viewer applies automatically. Color is also a **categorization tool**: the brain reads same-color items as belonging together (Gestalt similarity), so color encodes structure as much as it decorates.

But color perception is **relative and variable**: the same swatch looks different against different backgrounds, on different displays, and to different eyes — roughly 8% of men have a color-vision deficiency. This is why meaning must never rest on hue alone, and why contrast between color and background is a measurable, physiological requirement, not a preference.

## The Problem

Color chosen by mood, not by system:

- A dozen unrelated colors picked one at a time, with no relationships — visual chaos.
- Meaning carried only by color (red = error, green = ok) so color-blind users miss it entirely.
- Brand color on a CTA that fails contrast against its background — pretty and unreadable.
- Saturated colors competing everywhere, so nothing is emphasized and the eye has nowhere to rest.

## How to apply

- **Build a structured palette, not a list.** Define roles: a primary (brand), one or two secondaries, neutrals (the grays that do most of the actual work), and semantic colors (success/warning/error/info). Generate tints/shades as a consistent ramp, not by eye.
- **Limit the palette.** Fewer, well-related colors read as more intentional. Most strong systems run on a primary, a small accent set, and a robust neutral ramp.
- **Lean on neutrals.** The bulk of a UI is grays and near-blacks; reserve saturated color for emphasis and meaning so it retains its impact ([hierarchy](hierarchy.md)).
- **Verify contrast at definition time.** Every text/background and meaningful UI pairing must meet WCAG AA (4.5:1 body, 3:1 large/UI). Check it when you pick the color, not after.
- **Never encode meaning in hue alone.** Pair color with an icon, label, shape, or pattern so the meaning survives color-blindness and grayscale ([accessibility](accessibility.md)).
- **Use a perceptual color space** (e.g. OKLCH/HSL) for ramps so steps feel evenly spaced to the eye, not just numerically.

## Common violations to watch for

- **Palette sprawl.** Too many colors, or near-duplicates that drift over time ([repetition](repetition.md)).
- **Color-only meaning.** Status, links, or categories distinguished by hue with no secondary cue — a hard accessibility failure.
- **Failing contrast.** On-brand colors that drop below the required ratio against their background.
- **Over-saturation.** Everything vivid, so nothing stands out and the design feels tiring.
- **Ignoring context.** Not accounting for dark mode, varied displays, or print, where the same values behave differently.

## Connection to the guards

- **`visual-accessibility`** — measures color contrast ratios against WCAG AA/AAA for all text and meaningful UI, and checks that no information is conveyed by color alone. Color is one of this guard's central concerns alongside [contrast](contrast.md) and [typography](typography.md).
- **`brand-consistency`** — the palette is a locked system token set; any off-palette color in a produced piece is flagged. Color is often the most visible consistency violation.

## Related principles

- **[Contrast](contrast.md)** — color is a primary contrast dimension, and contrast ratio is color's accessibility constraint.
- **[Accessibility](accessibility.md)** — sets the contrast minimums and the "not by color alone" rule.
- **[Repetition](repetition.md)** — the palette is reused across the whole system.
- **[Hierarchy](hierarchy.md)** — reserved saturated color directs attention.
