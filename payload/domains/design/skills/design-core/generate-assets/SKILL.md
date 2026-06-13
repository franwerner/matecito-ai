---
name: generate-assets
description: >
  Produces visual assets in the brand's style: posts, carousels, flyers, banners, stories,
  presentations, campaign images. USE THIS SKILL whenever the user asks "make a post", "build me
  a flyer", "I need a banner", "design this piece", "a carousel of X", "images for social/campaign",
  or wants to bring the brand down to concrete material. If the piece must respect the identity,
  this skill leans on `brand-guide`.
---

# Generate assets

Turns the brand and a concrete request into ready-to-use assets while keeping consistency.

## Before generating

1. Consult `brand-guide` for palette, typography, and rules. If there is no guide, ask for at
   least the base palette and typography, or derive them with `brand-from-references`.
2. Clarify the format and destination: network/platform, dimensions, number of pieces, included
   copy.

## How to produce (depending on available tool)

- **Canva (`canva` MCP)** — for editable, on-brand pieces with templates: create/edit
  designs, autofill brand templates, search the library, and export (PNG/PDF/MP4). Ideal for
  social and marketing material because it respects the brand kit and stays editable for the
  person. Registered by install via Canva's official hosted MCP (`mcp.canva.com/mcp`, OAuth per person).
- **Image-generation MCP** — for original images (backgrounds, illustrations, hero images) that
  don't come from a template. Pass the brand direction as a reference.
- **No tools connected** — deliver the design as a precise spec (layout, hierarchy, text, colors by
  hex, font and sizes) so the person can build it in their editor.

## Sets and consistency

- For several pieces (carousel, multiple sizes), keep grid, typography, and palette coherent
  across all of them. They should read as a family.
- Respect the logo's clear space and the minimum contrast (see the brand guide).

## Rules

- On-brand always: if a piece departs from the guide, flag it before delivering.
- Don't reproduce third-party brands, logos, or assets.
- Sentence case in the copy unless the brand says otherwise.
- Deliver in the requested format/dimensions; if unspecified, propose the platform standards.

## Mentor mode

Explain the why behind each decision in 1-2 lines — the design principle, not just the what.
Derive to the `explain-concept` skill for unfamiliar concepts. (Full rule in the domain CLAUDE.md.)
