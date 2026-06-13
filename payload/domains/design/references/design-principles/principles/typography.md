# Typography

## Intent

Choose, size, and arrange type so text is read effortlessly and its structure is legible at a glance. Typography is the bulk of most interfaces and documents — it is where the design either supports reading or quietly sabotages it. Good typography is invisible; the reader notices only that they read comfortably.

## Why it matters (the perceptual reason)

Reading is **pattern recognition of word shapes**, not letter-by-letter decoding. Anything that disrupts those shapes — poor letterforms, cramped or loose spacing, low contrast, over-long lines — forces the eye to slow down and the brain to spend effort that should have gone to the content. Legibility (can I distinguish the characters?) and readability (can I read sustained text comfortably?) are measurable ergonomic properties, not matters of taste.

A **type scale** also encodes [hierarchy](hierarchy.md) structurally: distinct, related sizes let the reader skim and parse a document's structure pre-attentively. Typography is therefore both a legibility system and a hierarchy system at once.

## The Problem

Type treated as decoration rather than a reading system:

- A body font chosen for personality at the cost of legibility, fatiguing the reader.
- Headings only two points larger than body — no skim-able structure ([contrast](contrast.md)).
- Line length running 120 characters across — the eye loses its place returning to the next line.
- Line height set tight (1.0–1.2) on body text, so lines crowd and reading slows.
- Five fonts and a dozen ad-hoc sizes, with no scale and no consistency ([repetition](repetition.md)).

## How to apply

- **Define a type scale.** A small set of related sizes (e.g. a modular scale) for display / heading levels / body / caption, with clear jumps between them — this is hierarchy made structural.
- **Limit typefaces.** One or two families is plenty (often a display + a text face). More fonts fragment the system.
- **Optimize body for reading.** Aim for ~45–75 characters per line, line height around 1.4–1.6 for body, and a comfortable base size (≈16px / 1rem on web, larger for long-form).
- **Tune the micro-spacing.** Slightly tighten tracking on large display type; leave body tracking near default; use line height as the dominant rhythm control ([white-space](white-space.md)).
- **Establish weight and style roles.** Decide what bold, italic, caps, and color mean, and apply them consistently rather than for decoration.
- **Maintain contrast.** Text must meet WCAG AA against its background (4.5:1 body, 3:1 large); never pursue "subtle" gray text below the threshold ([color](color.md)).

## Common violations to watch for

- **Flat scale.** Headings, body, and captions nearly the same size — no hierarchy, nothing to skim.
- **Lines too long or too short.** Over ~75 or under ~45 characters both hurt readability.
- **Tight leading.** Body line height below ~1.4 crowds the text.
- **Too many fonts/sizes.** A grab-bag of typefaces and one-off sizes that no scale governs.
- **Light gray body text.** Low-contrast "elegant" text that fails the contrast minimum and excludes low-vision readers.
- **All caps for long text.** Capitals destroy word-shape recognition and slow reading; reserve for short labels.

## Connection to the guards

- **`visual-accessibility`** — checks minimum text sizes, text-to-background contrast, line-height/letter-spacing (WCAG text-spacing criteria), and that the visual reading order matches the document order screen readers follow. Typography is a primary surface this guard inspects.
- **`brand-consistency`** — the typeface set and type scale are locked system tokens; off-scale sizes or unauthorized fonts in a produced piece are flagged. Typographic drift is a common consistency failure.

## Related principles

- **[Hierarchy](hierarchy.md)** — the type scale is hierarchy encoded into text.
- **[Contrast](contrast.md)** — size/weight contrast in type, and text contrast ratio.
- **[White Space](white-space.md)** — line height and tracking are micro white space.
- **[Accessibility](accessibility.md)** — sets the size and contrast minimums for text.
- **[Repetition](repetition.md)** — the type scale is a reused system decision.
