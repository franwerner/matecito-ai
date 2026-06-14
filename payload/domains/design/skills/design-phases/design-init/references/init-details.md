# Design Init Details

## Design Capability Checklist

- Brand context: existing brand guide, `tokens.json` or design-token files, logo/asset folders, a prior design-system manifest.
- Figma: the `figma` MCP registered AND a file reachable (`get_file` succeeds).
- Canva: the `canva` MCP registered.
- Surface type: `landing | app-ui | brand-system | marketing`, inferred from the request and inspected files.
- DDR store: presence and content of `.matecito-ai/ddr/` (an `INDEX.md` or at least one DDR).

## Skill Loading (matecito-ai)

design-init builds no registry. Each design phase loads its own `SKILL.md`. Project conventions are read directly from their files: `.matecito-ai/ddr/` (design decisions), the brand guide, `CLAUDE.md`, and `config.yaml`.

## Engram Saves

```text
mem_save title/topic_key: design-init/{project}
type: architecture
content: detected design context markdown (includes the Design Capabilities block)
capture_prompt: false when available
```

design-init persists a single observation (`design-init/{project}`) that carries both the project context and the Design Capabilities block — the design equivalent of development's split `sdd-init/{project}` + `sdd/{project}/testing-capabilities`.

## Design Capabilities Format

```markdown
## Design Capabilities

**Detected**: {date}

### Surface

- Type: {landing | app-ui | brand-system | marketing}

### Brand Context

- Brand guide / design system: ✅ / ❌
- Path: `{path or —}`

### Connections

| Capability     | Available | Detail                                                |
| -------------- | --------- | ----------------------------------------------------- |
| figmaConnected | ✅ / ❌   | `figma` MCP reachable, a file connected / not found   |
| canvaConnected | ✅ / ❌   | `canva` MCP registered / not found                    |
| **available**  | ✅ / ❌   | figmaConnected ✅ (floor for reading/verifying work)  |

### DDR Store

- `.matecito-ai/ddr/`: present-with-content ✅ / absent-or-empty ❌

Detection notes:
- `figmaConnected`: the `figma` MCP must be registered AND a file reachable (`get_file` succeeds). If the MCP is not registered or no file is connected at init time, detected as ❌ even if a file exists in the user's Figma account. Document the limitation: Figma not connected at init → detected absent.
- `canvaConnected`: the `canva` MCP must be registered. Record ✅ / ❌.
- `available` = figmaConnected ✅; Figma is the read/verify index, so without it design-verify cannot run its guards against the real visual work.
- DDR store presence follows the DDR activation gate: active only when `.matecito-ai/ddr/` exists with content. design-init only records the gate; it never bootstraps DDRs.
```

## Output Templates

For each mode, include project, design context, persistence, the Design Capabilities table, the DDR-store presence, artifacts created/saved, limitations where relevant, and next steps. Engram mode must mention local/non-shareable limitations; none mode must recommend enabling persistence.
</content>
