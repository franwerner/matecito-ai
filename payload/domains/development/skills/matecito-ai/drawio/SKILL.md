---
name: drawio
version: 1.14.0
description: Use when the user requests diagrams, flowcharts, architecture diagrams, ER diagrams, UML / sequence / class diagrams, network topology, ML/DL model figures (Transformer/CNN/LSTM), mind maps, or any visualization. Also use proactively when explaining systems with 3+ components, complex data flows, or relationships that benefit from visual representation. Best suited when the diagram needs custom styling, rich shape vocabulary, swimlanes, or precise geometry. Renders the diagram live (ephemeral preview) via the `mcp__drawio__*` MCP without writing any files to the working directory.
license: MIT
homepage: https://github.com/Agents365-ai/drawio-skill
metadata: {"hermes":{"tags":["drawio","diagram","flowchart","architecture","visualization","uml"],"category":"design","related_skills":["mermaid","excalidraw","plantuml"]},"author":"Agents365-ai","version":"1.14.0"}
---

# Draw.io Diagrams

**Ephemeral by default — zero `.drawio`/image artifacts in the working directory. Live preview via the `mcp__drawio__*` MCP. File export only on explicit user request, to a path outside the repo.**

## Overview

This skill builds `<mxGraphModel>` XML **in memory** and renders it live via the `mcp__drawio__*` MCP (ephemeral live preview). It **NEVER** writes `.drawio`/PNG files to the working directory. The skill **complements** the MCP: the skill supplies the rich vocabulary (shapes, icons, styles, layout) and the MCP supplies the live render. Export to a file happens **only** on explicit user request, via `mcp__drawio__export_diagram`, to a path **outside the repo**.

## When to use / when NOT to use

**Use this skill for:** polished, precise diagrams (architecture, network, strict UML, ERD), anything needing solid opaque fills, 10,000+ stock/branded shapes, swimlanes, or custom geometry.

**Do NOT use it — route elsewhere — for:**
- A casual hand-drawn / whiteboard look → **excalidraw** or **tldraw**.
- Diagrams-as-code that live in git / render in Markdown → **mermaid** (general) or **plantuml** (UML).
- Freeform infinite-canvas sketching or freehand strokes → **tldraw**.

## Bundled resources

When the workflow references one of these, read it on demand — none of them need to be in context up front.

| File | Read it when |
|---|---|
| `references/diagram-types.md` | The user names a specific diagram type (ERD, UML class, sequence, architecture, ML/DL, flowchart) |
| `references/shapes.md` + `scripts/shapesearch.py` | The diagram needs a **specific shape** — a cloud icon (AWS/Azure/GCP), Cisco/Kubernetes/network symbol, UML/BPMN/ER/electrical/P&ID element — or any time you'd otherwise guess a `style=` string. `shapesearch.py "<keywords>"` returns the exact official style for 10k+ shapes |
| `scripts/aiicons.py` | The diagram involves an **AI/LLM brand** (OpenAI, Claude, Gemini, Mistral, Llama, HuggingFace, Ollama, LangChain, …) — `aiicons.py "<brand>"` returns a draw.io `image` style for the brand logo (lobe-icons via CDN; `--embed` to inline). draw.io has no built-in AI logos. See `references/shapes.md` → "AI / LLM brand logos" |
| `references/style-presets.md` | The user asks to learn / save / list / set-default / delete a style preset, or you've resolved an active preset and need the application rules |
| `references/style-extraction.md` | You're inside the Learn flow and need the extraction procedure (called from `style-presets.md`) |
| `references/troubleshooting.md` | A rendering looks wrong or vision rejects an image |
| `references/autolayout.md` | The diagram is large or layout-heavy (dependency/call graph, code structure, >~15 nodes) and you want Graphviz to place nodes + route edges instead of hand-placing coordinates — describe the graph as JSON (by hand or from any external analyzer) and run `autolayout.py` |
| `scripts/validate.py` | You generated a `<mxGraphModel>` (especially via autolayout or for a large hand-placed diagram) and want a fast deterministic structural lint (dangling edges, dup/reserved ids, broken parents, overlaps) before the vision self-check |

## Prerequisites

The render is provided by the `mcp__drawio__*` MCP, already available in the `development` domain — **no draw.io desktop CLI required**. Graphviz (`dot`) remains optional, needed **only** for `scripts/autolayout.py`.

## Workflow

Before starting the workflow, assess whether the user's request is specific enough. If key details are missing, ask 1-3 focused questions:
- **Diagram type** — which preset? (ERD, UML, Sequence, Architecture, ML/DL, Flowchart, or general)
- **Scope/fidelity** — how many components? Any specific technologies or labels?

Skip clarification if the request already specifies these details or is clearly simple (e.g., "draw a flowchart of X").

**Step 0 — Resolve active preset.** Determine which (if any) user-defined style preset applies to this generation.

- Scan the user's message for a phrase that clearly names a style preset: "use my `<name>` style", "with my `<name>` style", "in `<name>` mode", "in the style of `<name>`". A bare `with <name>` does **not** count — "draw a diagram with redis" names a component, not a style. If a clear match is found → active preset = `<name>`.
- Else, check `~/.drawio-skill/styles/` for any file with `"default": true`. If found → active preset = that one.
- Else → no preset active; fall through to the built-in color/shape/edge conventions for the rest of the workflow.

Load the preset JSON from `~/.drawio-skill/styles/<name>.json`, falling back to `<this-skill-dir>/styles/built-in/<name>.json`. If the named preset exists in neither location, tell the user the name is unknown, list the available presets (user dir + built-in), and stop — do **not** silently fall back to defaults.

When a preset loads successfully, mention it in the first line of the reply: *"Using preset `<name>` (confidence: `<level>`)."* See the **Applying a preset** subsection below for how the preset changes color/shape/edge/font decisions.

1. **Ensure the preview session** — call `mcp__drawio__start_session` once to open the live preview in the browser. `start_session` returns the actual preview URL — the port is assigned **dynamically** (commonly `6002`, or the next free one such as `6003`), so use whatever it reports rather than assuming a fixed port. No binary resolution is needed; the MCP provides the render.
2. **Plan** — identify shapes, relationships, layout (LR or TB), group by tier/layer. Start margins at `x=40, y=40` and space shapes per the **Layout tips** table below (gaps scale with node count). The MCP's `800×600` is the **single-page target for simple diagrams (≤5 nodes)** — keep those within it. Medium/large diagrams *will* exceed it, and that's fine: the live preview scrolls and zooms. **Prioritize legibility (the table's gaps) over fitting one page** — never compress a large diagram to force it into `800×600`.
3. **Generate** — build the `<mxGraphModel>` XML **in memory** (do not write to the working dir). Hand-place coordinates for small/styled diagrams. **For large or layout-heavy diagrams (dependency/call graphs, code structure, >~15 nodes), don't hand-place** — describe the graph as JSON (write it by hand, or generate it with any external analyzer that emits nodes/edges — see `references/autolayout.md`) and run `python3 <this-skill-dir>/scripts/autolayout.py graph.json -o /tmp/<name>.drawio` to compute node positions + orthogonal edge routing via Graphviz. **`autolayout.py` writes its `.drawio` output to a system temp path (e.g. `/tmp`), never to the working dir** — you then **read the XML back from the temp file** to pass it to the MCP. Run `python3 <this-skill-dir>/scripts/validate.py /tmp/<name>.drawio` against that temp file for a fast structural lint (dangling edges, dup ids, overlaps) before rendering. **IMPORTANT:** the MCP consumes a bare `<mxGraphModel>` **without** the `<mxfile><diagram>` wrapper; if a tool produced the full `<mxfile>`, extract only the inner `<mxGraphModel>...</mxGraphModel>` before passing it to the MCP.
4. **Render** — pass the `<mxGraphModel>` XML to `mcp__drawio__create_new_diagram(xml)`; the diagram appears live in the preview opened by `start_session`. There is no preview PNG export, no `-e` flag, no PNG repair. (Note: `create_new_diagram` **destroys** any prior diagram — use it only for the first render or an explicit "start over".)
5. **Self-check** — optional, if the model has vision. Call `mcp__drawio__get_diagram` to inspect the structure and/or look at the live preview, then apply the checks in the table below. Apply any fixes via `mcp__drawio__edit_diagram`.

| Check | What to look for | Fix action |
|-------|-----------------|-----------------|
| Overlapping shapes | Two or more shapes stacked on top of each other | Shift shapes apart by ≥200px |
| Clipped labels | Text cut off at shape boundaries | Increase shape width/height to fit label |
| Missing connections | Arrows that don't visually connect to shapes | Verify `source`/`target` ids match existing cells |
| Off-canvas shapes | Shapes at negative coordinates, or stranded far from the main cluster | Move to positive coordinates near the cluster (exceeding `800×600` is fine for medium/large diagrams) |
| Edge-shape overlap | An edge/arrow visually crosses through an unrelated shape | Add waypoints (`<Array as="points">`) to route around the shape, or increase spacing between shapes |
| Stacked edges | Multiple edges overlap each other on the same path | Distribute entry/exit points across the shape perimeter (use different exitX/entryX values) |

- Max **2 self-check rounds** — if issues remain after 2 fixes, show the user anyway.

6. **Review loop** — show the live preview to the user and collect feedback. Apply **minimal** edits via `mcp__drawio__get_diagram` (always read the current state first — it includes the user's manual edits) followed by `mcp__drawio__edit_diagram` (add/update/delete operations keyed by `cell_id`, each carrying a complete `mxCell` in `new_xml`).

**Targeted edit rules** — for each type of feedback, apply the minimal MCP operation:

| User request | MCP edit operation |
|-------------|----------------|
| Change color of X | `update` the cell whose `value` matches X — new `mxCell` with `fillColor`/`strokeColor` adjusted in `style` |
| Add a new node | `add` a new `mxCell` vertex with the next available `cell_id`, positioned near related nodes |
| Remove a node | `delete` the `mxCell` vertex and `delete` any edges with matching `source`/`target` |
| Move shape X | `update` the matching cell — new `mxCell` with adjusted `x`/`y` in `mxGeometry` |
| Resize shape X | `update` the matching cell — new `mxCell` with adjusted `width`/`height` in `mxGeometry` |
| Add arrow from A to B | `add` a new `mxCell` edge with `source`/`target` matching A and B ids |
| Change label text | `update` the matching cell — new `mxCell` with the new `value` |
| Change layout direction | **Full regeneration** — rebuild the `<mxGraphModel>` and re-render with `create_new_diagram` |

**Rules:**
- For single-element changes: `update`/`add`/`delete` individual cells — preserves layout tuning from prior iterations.
- For layout-wide changes (e.g., swap LR↔TB, "start over"): regenerate the full `<mxGraphModel>` and call `create_new_diagram` again.
- Always call `mcp__drawio__get_diagram` before editing so you operate on the current state.
- Loop continues until the user says approved / done / LGTM.

7. **Final (no export by default)** — the deliverable is the **live preview**; nothing is written to disk. **Only** if the user explicitly asks for a file, export via `mcp__drawio__export_diagram(path, format)` to a path **outside the repo** that the user specifies.

## Style Presets

A **style preset** is a named JSON file capturing a user's visual preferences (palette, shapes, font, edges). When active, it fully replaces the built-in color/shape conventions in this skill.

**Lookup order** when SKILL.md's Step 0 resolves a preset name:
1. `~/.drawio-skill/styles/<name>.json` — user presets (survive `git pull`)
2. `<this-skill-dir>/styles/built-in/<name>.json` — shipped built-ins (`default`, `corporate`, `handdrawn`)

Always lowercase the user-provided name before any file operation — the schema enforces lowercase.

**For everything else — Learn flow (extracting a preset from a file), management ops (list/default/delete/rename), application rules (color lookup, shape keywords, edges, fonts, extras, interaction with diagram-type presets), and validation — read `references/style-presets.md`.** It's only needed when the user invokes those flows or when an active preset must be applied to the current generation.

## Draw.io XML Structure

### File skeleton

The MCP consumes the inner `<mxGraphModel>` **directly** — pass it **without** the `<mxfile><diagram>` wrapper:

```xml
<mxGraphModel>
  <root>
    <mxCell id="0" />
    <mxCell id="1" parent="0" />
    <!-- user shapes start at id="2" -->
  </root>
</mxGraphModel>
```

The full `<mxfile>` wrapper below only applies if you **export to a file on explicit request**; it is not what the MCP expects for live rendering:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<mxfile host="drawio" version="26.0.0">
  <diagram name="Page-1">
    <mxGraphModel>
      <root>
        <mxCell id="0" />
        <mxCell id="1" parent="0" />
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>
```

**Rules:**
- `id="0"` and `id="1"` are required root cells — never omit them
- User shapes start at `id="2"` and increment sequentially
- All shapes have `parent="1"` (unless inside a container — then use container's id)
- All text uses `html=1` in style for proper rendering
- **Never use `--` inside XML comments** — it's illegal per XML spec and causes parse errors
- Escape special characters in attribute values: `&amp;`, `&lt;`, `&gt;`, `&quot;`
- **Multi-line text in labels:** use `&#xa;` for line breaks inside `value` attributes (not literal `\n`). Example: `value="Line 1&#xa;Line 2"`

### Shape types (vertex)

| Style keyword | Use for |
|--------------|---------|
| `rounded=0` | plain rectangle (default) |
| `rounded=1` | rounded rectangle — services, modules |
| `ellipse;` | circles/ovals — start/end, databases |
| `rhombus;` | diamond — decision points |
| `shape=mxgraph.aws4.resourceIcon;` | AWS icons |
| `shape=cylinder3;` | cylinder — databases |
| `swimlane;` | group/container with title bar |

For **vendor/branded icons** (AWS/Azure/GCP/Cisco/Kubernetes) and any non-trivial shape, don't guess the `shape=mxgraph.*` name — a wrong name renders as a blank box. Run `python3 <this-skill-dir>/scripts/shapesearch.py "<keywords>"` to get the exact official style + size, or see `references/shapes.md` for the hand-writable cheatsheet. For **AI/LLM brand logos** (OpenAI, Claude, Gemini, …), which draw.io has none of, use `python3 <this-skill-dir>/scripts/aiicons.py "<brand>"`.

### Required properties

```xml
<!-- Rectangle / rounded box -->
<mxCell id="2" value="Label" style="rounded=1;whiteSpace=wrap;html=1;fillColor=#dae8fc;strokeColor=#6c8ebf;" vertex="1" parent="1">
  <mxGeometry x="100" y="100" width="160" height="60" as="geometry" />
</mxCell>

<!-- Cylinder (database) -->
<mxCell id="3" value="DB" style="shape=cylinder3;whiteSpace=wrap;html=1;fillColor=#f5f5f5;strokeColor=#666666;fontColor=#333333;" vertex="1" parent="1">
  <mxGeometry x="350" y="100" width="120" height="80" as="geometry" />
</mxCell>

<!-- Diamond (decision) -->
<mxCell id="4" value="Check?" style="rhombus;whiteSpace=wrap;html=1;fillColor=#fff2cc;strokeColor=#d6b656;" vertex="1" parent="1">
  <mxGeometry x="100" y="220" width="160" height="80" as="geometry" />
</mxCell>
```

### Containers and groups

For architecture diagrams with nested elements, use draw.io's parent-child containment — do **not** just place shapes on top of larger shapes.

| Type | Style | When to use |
|------|-------|-------------|
| **Group** (invisible) | `group;pointerEvents=0;` | No visual border needed, container has no connections |
| **Swimlane** (titled) | `swimlane;startSize=30;` | Container needs a visible title bar, or container itself has connections |
| **Custom container** | Add `container=1;pointerEvents=0;` to any shape | Any shape acting as a container without its own connections |

**Key rules:**
- Add `pointerEvents=0;` to container styles that should not capture connections between children
- Children set `parent="containerId"` and use coordinates **relative to the container**

```xml
<!-- Swimlane container -->
<mxCell id="svc1" value="User Service" style="swimlane;startSize=30;fillColor=#dae8fc;strokeColor=#6c8ebf;" vertex="1" parent="1">
  <mxGeometry x="100" y="100" width="300" height="200" as="geometry"/>
</mxCell>
<!-- Child inside container — coordinates relative to parent -->
<mxCell id="api1" value="REST API" style="rounded=1;whiteSpace=wrap;html=1;" vertex="1" parent="svc1">
  <mxGeometry x="20" y="40" width="120" height="60" as="geometry"/>
</mxCell>
<mxCell id="db1" value="Database" style="shape=cylinder3;whiteSpace=wrap;html=1;" vertex="1" parent="svc1">
  <mxGeometry x="160" y="40" width="120" height="60" as="geometry"/>
</mxCell>
```

### Connector (edge)

**CRITICAL:** Every edge `mxCell` must contain a `<mxGeometry relative="1" as="geometry" />` child element. Self-closing edge cells (`<mxCell ... edge="1" ... />`) are **invalid** and will not render. Always use the expanded form.

```xml
<!-- Directed arrow — always include rounded, orthogonalLoop, jettySize for clean routing -->
<mxCell id="10" value="" style="edgeStyle=orthogonalEdgeStyle;rounded=1;orthogonalLoop=1;jettySize=auto;html=1;" edge="1" parent="1" source="2" target="3">
  <mxGeometry relative="1" as="geometry" />
</mxCell>

<!-- Arrow with label + explicit entry/exit points to control direction -->
<mxCell id="11" value="HTTP/REST" style="edgeStyle=orthogonalEdgeStyle;rounded=1;orthogonalLoop=1;jettySize=auto;html=1;exitX=0.5;exitY=1;exitDx=0;exitDy=0;entryX=0.5;entryY=0;entryDx=0;entryDy=0;" edge="1" parent="1" source="2" target="4">
  <mxGeometry relative="1" as="geometry" />
</mxCell>

<!-- Arrow with waypoints — use when edge must route around other shapes -->
<mxCell id="12" value="" style="edgeStyle=orthogonalEdgeStyle;rounded=1;orthogonalLoop=1;jettySize=auto;html=1;" edge="1" parent="1" source="3" target="5">
  <mxGeometry relative="1" as="geometry">
    <Array as="points">
      <mxPoint x="500" y="50" />
    </Array>
  </mxGeometry>
</mxCell>
```

**Edge style rules:**
- **Animated connectors:** add `flowAnimation=1;` to any edge style to show a moving dot animation along the arrow — ideal for data-flow and pipeline diagrams. Example: `style="edgeStyle=orthogonalEdgeStyle;flowAnimation=1;rounded=1;..."`
- **Always** include `rounded=1;orthogonalLoop=1;jettySize=auto` — these enable smart routing that avoids overlaps
- Pin `exitX/exitY/entryX/entryY` on every edge when a node has 2+ connections — distributes lines across the shape perimeter
- Add `<Array as="points">` waypoints when an edge must detour around an intermediate shape
- **Leave room for arrowheads:** the final straight segment between the last bend and the target shape must be ≥20px long. If too short, the arrowhead overlaps the bend and looks broken. Fix by increasing node spacing or adding explicit waypoints

### Distributing connections on a shape

When multiple edges connect to the same shape, assign different entry/exit points to prevent stacking:

| Position | exitX/entryX | exitY/entryY | Use when |
|----------|-------------|-------------|----------|
| Top center | 0.5 | 0 | connecting to node above |
| Top-left | 0.25 | 0 | 2nd connection from top |
| Top-right | 0.75 | 0 | 3rd connection from top |
| Right center | 1 | 0.5 | connecting to node on right |
| Bottom center | 0.5 | 1 | connecting to node below |
| Left center | 0 | 0.5 | connecting to node on left |

**Rule:** if a shape has N connections on one side, space them evenly (e.g., 3 connections on bottom → exitX = 0.25, 0.5, 0.75)

### Color palette (fillColor / strokeColor)

*Used only when no preset is active (see "Applying a preset" above).*

| Color name | fillColor | strokeColor | Use for |
|-----------|-----------|-------------|---------|
| Blue | `#dae8fc` | `#6c8ebf` | services, clients |
| Green | `#d5e8d4` | `#82b366` | success, databases |
| Yellow | `#fff2cc` | `#d6b656` | queues, decisions |
| Orange | `#ffe6cc` | `#d79b00` | gateways, APIs |
| Red/Pink | `#f8cecc` | `#b85450` | errors, alerts |
| Grey | `#f5f5f5` | `#666666` | external/neutral |
| Purple | `#e1d5e7` | `#9673a6` | security, auth |

### Layout tips

**Spacing — scale with complexity:**

| Diagram complexity | Nodes | Horizontal gap | Vertical gap |
|-------------------|-------|----------------|--------------|
| Simple | ≤5 | 200px | 150px |
| Medium | 6–10 | 280px | 200px |
| Complex | >10 | 350px | 250px |

These gaps take priority over the `800×600` single-page target: only a Simple diagram is expected to fit it. Medium/Complex layouts grow past the viewport by design — the live preview scrolls and zooms, so keep the gaps and let the canvas extend.

**Routing corridors:** between shape rows/columns, leave an extra ~80px empty corridor where edges can route without crossing shapes. Never place a shape in a gap that edges need to traverse.

**Grid alignment:** snap all `x`, `y`, `width`, `height` values to **multiples of 10** — this ensures shapes align cleanly on draw.io's default grid and makes manual editing easier.

**General rules:**
- Plan a grid before assigning x/y coordinates — sketch node positions on paper/mentally first
- Group related nodes in the same horizontal or vertical band
- Use `swimlane` cells for logical grouping with visible borders
- Place heavily-connected "hub" nodes centrally so edges radiate outward instead of crossing
- To force straight vertical connections, pin entry/exit points explicitly on edges:
  `exitX=0.5;exitY=1;exitDx=0;exitDy=0;entryX=0.5;entryY=0;entryDx=0;entryDy=0`
- Always center-align a child node under its parent (same center x) to avoid diagonal routing
- **Event bus pattern**: place Kafka/bus nodes in the **center of the service row**, not below — services on either side can reach it with short horizontal arrows (`exitX=1` left side, `exitX=0` right side), eliminating all line crossings
- Horizontal connections (`exitX=1` or `exitX=0`) never cross vertical nodes in the same row; use them for peer-to-peer and publish connections

**Avoiding edge-shape overlap:**
- Before finalizing coordinates, trace each edge path mentally — if it must cross an unrelated shape, either move the shape or add waypoints
- For tree/hierarchical layouts: assign nodes to layers (rows), connect only between adjacent layers to minimize crossings
- For star/hub layouts: place the hub center, satellites around it — edges stay short and radial
- When an edge must span multiple rows/columns, route it along the outer corridor, not through the middle of the diagram

## Export (only on explicit user request)

By default **nothing is exported** — the deliverable is the **live preview**. If, and only if, the user explicitly asks for a file, use `mcp__drawio__export_diagram(path, format)` with `format` of `png`, `svg`, or `drawio`, writing to a path **outside the repo** that the user specifies. There is no draw.io CLI, no PNG repair, no browser fallback, and no PATH/WSL2 handling — the MCP owns the render and the export.

## Common Mistakes

When something looks wrong (layout broken, edges misroute, vision rejects an image), see `references/troubleshooting.md` for a row-by-row mistake → fix table.

## Diagram Type Presets

When the user requests a specific diagram type, read `references/diagram-types.md` for the matching preset (shapes, edges, layout direction). Pick by user phrasing:

| User says | Section in `references/diagram-types.md` |
|---|---|
| "ER diagram", "schema diagram", "data model" | ERD |
| "UML class diagram", "class diagram" | UML Class |
| "sequence diagram", "interaction diagram", "lifeline" | Sequence |
| "architecture", "system diagram", "service diagram" | Architecture |
| "neural network", "model architecture", "ML diagram", "deep learning" | ML / Deep Learning Model |
| "flowchart", "decision tree", "process flow" | Flowchart |

The diagram-type preset sets **structural** style keywords. If a user style preset is also active (see `## Style Presets`), keep the structural keywords and layer color/font/edge/extras on top — read `references/style-presets.md` → "Interaction with diagram-type presets" for the merge rules.
