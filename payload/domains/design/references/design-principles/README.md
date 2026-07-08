# Design Principles — Reference catalog

> **This is a consultable catalog, not a skill.** The `design-*` phase agents and design skills reference it by principle name. When a DDR (Design Decision Record) declares `Applied principle: X`, the canonical definition lives in [`principles/<x>.md`](principles/).

## Overview

Design principles are the perceptual and cognitive laws that govern how people read, parse, and trust a visual composition. They are not stylistic preferences — they are constraints rooted in how human vision and attention actually work. A layout that respects them feels effortless; one that violates them feels "off" even to viewers who can't say why.

**Core Philosophy:** Principles are not a checklist to decorate a design after the fact — they are the reasoning you apply *while* deciding. Use a principle when it solves a real perception problem (a viewer can't find the primary action, can't tell two groups apart, can't read the body text), not to justify a choice already made on taste alone.

This catalog is the **single source of truth for design rationale** in the design domain. Per the domain's **mentor mode** rule, every phase and skill explains the *why* behind a decision by citing a principle here rather than improvising the reasoning. The `design-system` phase consults it before locking the visual system; `design-verify` and its guards (`visual-accessibility`, `brand-consistency`) check produced work against it.

## How this catalog is used

- **Skills cite it.** `brand-guide`, `explore-variations`, `consistency-audit`, `design-review`, and the `explain-concept` mentor skill reference a principle by name instead of restating the rationale. The canonical text lives here once.
- **DDRs reference it.** A Design Decision Record captures `Applied principle: X` the same way a development EDR captures `Applied pattern: X`. Before implementing the decision, consult `principles/<x>.md` to know the principle's contract; if you deviate, justify it in the DDR.
- **The design phases consult it before deciding.** `design-system` locks the palette, type scale, grid, and spacing *against* these principles. `design-verify` flags any produced piece that breaks one.
- **The guards enforce a subset of it.** `visual-accessibility` operationalizes [`accessibility`](principles/accessibility.md), [`contrast`](principles/contrast.md), [`color`](principles/color.md), and [`typography`](principles/typography.md) into WCAG-measurable checks. `brand-consistency` leans on [`repetition`](principles/repetition.md) to detect pieces that drift from the system.

## Foundational idea

Every principle below answers one of two questions the viewer's brain asks automatically:

| The viewer asks… | Principles that answer it |
|------------------|---------------------------|
| **"What do I look at first, and in what order?"** | [Hierarchy](principles/hierarchy.md), [Contrast](principles/contrast.md), [White Space](principles/white-space.md) |
| **"What belongs together, and what is separate?"** | [Proximity & Gestalt](principles/proximity-gestalt.md), [Alignment](principles/alignment.md), [Repetition](principles/repetition.md) |
| **"Can I actually read and use this?"** | [Color](principles/color.md), [Typography](principles/typography.md), [Accessibility](principles/accessibility.md) |

Hierarchy is the goal; the rest are the tools that produce it without the viewer noticing the machinery.

## Index

### Tier 1: Structural principles (decide these first)

| Principle | One-line | Perceptual reason | Reference |
|-----------|----------|-------------------|-----------|
| **Hierarchy** | Guide the eye through content in order of importance | The brain triages by salience before it reads | [hierarchy.md](principles/hierarchy.md) |
| **Contrast** | Make different things look clearly different | Difference is what draws and holds attention | [contrast.md](principles/contrast.md) |
| **Proximity & Gestalt** | Group related elements; the mind sees the whole | We perceive grouped objects as one unit | [proximity-gestalt.md](principles/proximity-gestalt.md) |
| **Alignment** | Place elements on shared invisible lines | The eye follows edges; misalignment reads as noise | [alignment.md](principles/alignment.md) |

### Tier 2: System principles (lock these in the design system)

| Principle | One-line | Perceptual reason | Reference |
|-----------|----------|-------------------|-----------|
| **Repetition** | Reuse visual decisions consistently | Repetition builds recognition and trust | [repetition.md](principles/repetition.md) |
| **White Space** | Let elements breathe; emptiness is active | Space separates, emphasizes, and reduces load | [white-space.md](principles/white-space.md) |
| **Color** | Use a deliberate, limited, meaningful palette | Color carries meaning, emotion, and contrast | [color.md](principles/color.md) |
| **Typography** | Choose and scale type for reading, not decoration | Legibility and rhythm govern whether text is read | [typography.md](principles/typography.md) |

### Tier 3: Inclusive principle (constrains all of the above)

| Principle | One-line | Perceptual reason | Reference |
|-----------|----------|-------------------|-----------|
| **Accessibility** | Design so everyone can perceive and use it | A design only works if its audience can read it | [accessibility.md](principles/accessibility.md) |

## How the principles interact

Principles are not independent knobs — they trade against and reinforce each other:

- **Contrast serves Hierarchy.** You create hierarchy *by means of* contrast (size, weight, color, space). Hierarchy is the intent; contrast is the mechanism.
- **Proximity and White Space are the same gesture from two sides.** Pulling related items together (proximity) is what creates the gaps (white space) that separate groups.
- **Alignment and Repetition produce consistency.** Shared edges plus reused decisions are what make a multi-page system read as one brand — the basis of the `brand-consistency` guard.
- **Color and Typography are where Contrast and Accessibility collide.** A palette can satisfy brand and still fail a contrast ratio; the `visual-accessibility` guard exists precisely at that intersection.

## Common mistakes (across all principles)

| Mistake | Symptom | Fix |
|---------|---------|-----|
| **Everything emphasized** | Nothing stands out; flat, overwhelming | Establish one clear focal point (Hierarchy + Contrast) |
| **Decorative variation** | Fonts/colors/spacings differ for no reason | Reuse a small set of decisions (Repetition) |
| **Cramming** | No room to breathe; groups blur together | Add White Space; apply Proximity deliberately |
| **Taste over measurement** | "Looks fine to me" but fails for real users | Verify against Accessibility, Contrast ratios |
| **Style before structure** | Polishing color/type before hierarchy is set | Resolve structure (Tier 1) before surface (Tier 2) |

## Consultation checklist

Before locking a visual decision in `design-system` or shipping a piece in `design-produce`:

- [ ] Is the primary focal point unambiguous? (Hierarchy)
- [ ] Does each group read as a unit, clearly separate from others? (Proximity, White Space)
- [ ] Is every element aligned to a shared line or grid? (Alignment)
- [ ] Are type, color, and spacing decisions reused, not reinvented? (Repetition)
- [ ] Do text and UI meet contrast and size minimums for the real audience? (Contrast, Color, Typography, Accessibility)
- [ ] If a DDR drove this, does the result honor its `Applied principle`? (or is the deviation justified)
