# Troubleshooting — Common Mistakes

Read this when something looks wrong in the live preview (rendering, layout, edges) or vision rejects an image. Most rows have a one-line fix. The render is provided by the `mcp__drawio__*` MCP — there is no CLI to fail, no binary to locate, and no PNG export in the default flow.

| Mistake | Fix |
|---------|-----|
| Missing `id="0"` and `id="1"` root cells | Always include both at the top of `<root>` |
| Shapes not connected | `source` and `target` on edge must match existing shape `id` values |
| Self-closing edge `mxCell` (`<mxCell ... edge="1" />`) | Use the expanded form with `<mxGeometry relative="1" as="geometry" />` child — self-closing edges won't render |
| `--` inside XML comments | Illegal per XML spec — use single hyphens or rephrase |
| Special characters in `value` | Use XML entities: `&amp;` `&lt;` `&gt;` `&quot;` |
| Literal `\n` in label text | Use `&#xa;` for line breaks in `value` attributes |
| Shape renders as a blank box | The `style=` string is malformed or names a non-existent `shape=mxgraph.*` library — run `shapesearch.py "<keywords>"` for the exact official style instead of guessing |
| Invalid XML passed to the MCP (diagram doesn't render) | The XML wasn't a well-formed `<mxGraphModel>` — run `scripts/validate.py` on it, or check that you stripped the `<mxfile><diagram>` wrapper before calling `create_new_diagram` |
| Passed the full `<mxfile>` wrapper to `create_new_diagram` | The MCP consumes a bare `<mxGraphModel>` — extract only the inner `<mxGraphModel>...</mxGraphModel>` |
| Overlapping shapes | Scale spacing with complexity (200–350px); leave routing corridors |
| Edges crossing through shapes | Add waypoints, distribute entry/exit points, or increase spacing |
| Arrowhead overlaps bend | Final edge segment before target must be ≥20px — increase spacing or add waypoints |
| Iteration loop never ends | After 5 rounds, export the `.drawio` on the user's request (`mcp__drawio__export_diagram` to a path outside the repo) so they can fine-tune it in draw.io desktop |
| Edit didn't apply / hit a stale cell | Always call `mcp__drawio__get_diagram` before `edit_diagram` — it returns the current state including the user's manual edits |
| Vision returns "Unable to resize image — dimensions exceed the 2576x2576px limit" | The diagram is too large/dense for Claude's vision API. Inspect the structure via `mcp__drawio__get_diagram` (XML) instead of an image, or split the diagram across logical sub-views and render them one at a time. |
