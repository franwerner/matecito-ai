# Proximity & Gestalt

## Intent

Control how the viewer groups elements into meaningful units. Things placed close together are read as related; things placed apart are read as separate. The broader Gestalt laws describe how the mind assembles individual marks into wholes — and a designer who knows them composes *with* perception instead of against it.

## Why it matters (the perceptual reason)

The brain does not see a screen as a list of independent shapes — it organizes them automatically into groups, before conscious thought, according to the **Gestalt principles** (early-20th-century perceptual psychology). The strongest of these for layout is **proximity**: spatial nearness is the cheapest, most powerful grouping signal the visual system has. When you group related items by distance, you offload the "what goes with what" question from the viewer's working memory onto their visual system, which answers it for free.

Ignore this and the brain *still* groups — just incorrectly. A label too far from its field, or too close to the next field, gets attached to the wrong thing, and the viewer makes errors they'll blame on themselves.

## The Gestalt laws you'll use most

- **Proximity** — near = related. The default grouping tool.
- **Similarity** — same color/shape/size = same category, even when apart.
- **Common region** — a shared container (card, box, background) groups its contents.
- **Continuity** — the eye follows lines and aligned edges as connected paths ([alignment](alignment.md)).
- **Closure** — the mind completes implied shapes; gaps read as wholes.
- **Figure/ground** — the eye separates a subject from its background; ambiguous figure/ground is disorienting.

## The Problem

Grouping that fights the content's structure:

- A form where labels sit equidistant between two fields — each label could belong to either.
- A pricing table where the gap between rows equals the gap between columns — the grid loses its read.
- Related navigation links spaced as far apart as unrelated ones — no group reads as a group.
- Caption text floating closer to the *next* image than to the one it describes.

## How to apply

- **Make intra-group spacing smaller than inter-group spacing.** This single rule (relationships expressed as distance ratios, not absolute gaps) does most of the work. Tie the gaps to a spacing scale ([white-space](white-space.md)).
- **Use one grouping cue per relationship,** reinforced if needed — proximity for the primary grouping, similarity or a shared container as backup.
- **Don't add boxes when space will do.** Borders and cards are heavier than whitespace; reach for common-region containers only when proximity alone can't carry the grouping.
- **Respect figure/ground** — keep enough separation between a subject and its background that the eye doesn't have to work to find the figure.

## Common violations to watch for

- **Equal spacing everywhere.** When every gap is identical, proximity conveys nothing and the layout flattens into an undifferentiated grid.
- **Orphaned labels/captions.** Text closer to the wrong element than the right one.
- **Box overload.** Solving every grouping with another bordered container, producing visual clutter (often a [white-space](white-space.md) and [repetition](repetition.md) problem too).
- **Ambiguous figure/ground.** Busy backgrounds that compete with foreground content.

## Connection to the guards

- **`visual-accessibility`** — grouping by proximity (not by color alone) helps users with low vision and screen-reader users, since a logical visual grouping should match the document/DOM order. Captions and labels must be programmatically associated with what they describe, which proximity should mirror.
- **`brand-consistency`** — the spacing ratios that express grouping (the spacing scale) are part of the locked system; pieces that invent ad-hoc gaps drift from the brand's rhythm.

## Related principles

- **[White Space](white-space.md)** — the medium through which proximity operates; grouping *is* spacing.
- **[Alignment](alignment.md)** — Gestalt continuity; shared edges reinforce grouping.
- **[Repetition](repetition.md)** — Gestalt similarity applied across a whole system.
