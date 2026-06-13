---
name: brand-from-references
description: >
  Defines a brand or visual identity from reference images: the person uploads screenshots,
  photos, designs they like, or links to sites/Dribbble, and Claude distills the palette,
  typography, shape language, and style descriptors, then captures them as reusable tokens.
  USE THIS SKILL whenever the user wants to "define my brand", "build an identity", "pull the
  palette from these images", "take inspiration from these designs", "create a moodboard",
  "what style is this", or uploads images/refs asking to derive a visual direction. Also when
  they want to generate new designs "in this style".
---

# Brand from references

Turns a set of visual references into a concrete, reusable brand direction. Claude analyzes
images with vision directly: no external tool is needed to extract them; tools only help to
*gather* references and to *apply* the style afterward.

## Input

One or more of:
- Images uploaded to the chat (screenshots, photos, pieces they like).
- Links to live sites, Dribbble, Behance, etc.
- A text brief ("I want something minimal, warm, with a serious feel").

If none of these is present, ask for at least 3-5 references before continuing: with a single
one, the direction comes out thin.

## Flow

1. **Gather references (optional).** If a live site is provided, read it to extract real tokens;
   otherwise work with what the user uploaded.
2. **Read each reference** and note: palette (with approximate hex), typography (serif/sans,
   weight, contrast, personality), spacing/density, shapes (radii, borders, shadows), and the
   "mood" in 3-5 adjectives.
3. **Distill the common direction.** Don't copy one reference: find the pattern they share and
   propose ONE coherent direction. If the references contradict each other (e.g. one minimal and
   one maximalist), show 2 directions and let the person choose.
4. **Pour it into the brand guide.** Produce the output in the sections used by the `brand-guide`
   skill (palette with usage, typography, shapes, tone), with semantic names. If a guide already
   exists, integrate there instead of creating something parallel.
5. **Show and validate.** Present the palette and the type samples so the person confirms or
   adjusts before locking the brand.

## Applying the style

Once the brand is defined:
- Lock it in the guide with `brand-guide`.
- Produce pieces in that style with `generate-assets` (Canva / image MCP).
- Explore more directions from it with `explore-variations`.

## Rules

- The goal is to **take inspiration and distill an own brand**, not to clone a third party's
  identity. Do not reproduce logos, trademarks, or specific assets of another company.
- Sentence case, and semantic, consistent token names.
- Always validate the direction with the person before locking it: the brand is their decision.
- Give concrete hex and names, no "a nice blue". If you estimate a color from an image, say so.

## Mentor mode

Explain the why behind each decision in 1-2 lines — the design principle, not just the what.
Derive to the `explain-concept` skill for unfamiliar concepts. (Full rule in the domain CLAUDE.md;
cite the `design-principles` catalog at `~/.claude/references/design-principles/` for rationale.)
