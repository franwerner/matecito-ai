# Dominios de área

Cada subdirectorio acá es un **dominio de área** (development, design, …) — un plugin
sobre el kernel agnóstico de dominio en `../core/`. Un dominio entrega:

```
domains/<id>/
├── manifest.json     # contrato legible por máquina (Go: deploy, MCP, gating de checks)
├── CLAUDE.md         # fragmento para el agente, se anexa a core/CLAUDE.md en el deploy
├── agents/           # agentes de fase (opcional)        → ~/.claude/agents/
├── references/       # catálogos consultables (opcional) → ~/.claude/references/
└── skills/<group>/   # skills (opcional)                 → ~/.claude/skills/
```

`agents/`, `references/` y `skills/<group>/` se aplanan dentro de los árboles
compartidos `~/.claude/...` en el momento del deploy (la capa `<group>` bajo
`skills/` es un nivel organizativo y se descarta).

## Convención de nombres de skills (sin prefijo automático)

Como todo se aplana dentro del `~/.claude/skills/` compartido, **el nombre de
carpeta de una skill debe ser globalmente único entre todos los dominios
instalados** — y ese nombre de carpeta ES el comando de invocación en Claude Code
(el frontmatter `name:` del SKILL.md es solo una etiqueta de display). **No hay
prefijo de dominio automático.**

Nombrá las carpetas de tus skills para que no choquen con otros dominios; prefijá
con el id del dominio ante cualquier riesgo (p. ej. `design-audit` en vez de
`audit`). Si dos dominios activos exponen la misma carpeta de skill, el deploy
**falla rápido** con un error de colisión consciente del dominio (ver
`deploy.clashError`) — nada se sobrescribe en silencio.

## Shared tier

`payload/shared/` es un nivel de deploy **transversal**: sus componentes se despliegan **incondicionalmente**, sin importar qué dominios tenga activos el usuario.

```
payload/shared/
├── skills/<group>/   # skills cross-domain  → ~/.claude/skills/
├── agents/           # agentes cross-domain → ~/.claude/agents/
└── references/       # catálogos           → ~/.claude/references/
```

Las mismas reglas de aplanamiento que los dominios aplican: `skills/<group>/<skill>/` se aplana a `~/.claude/skills/<skill>/`; la capa `<group>` es solo organizativa. La convención de nombres únicos de skills sigue vigente: si un componente shared y uno de dominio apuntan al mismo destino, el deploy **falla con un error de colisión** (el guard `seen[Target]→Source` de `Plan()` no tiene excepciones para el tier shared).

**No existe `payload/shared/hooks/`.** Los hooks siempre activos se registran como compiled-in usando `hook.SharedDomain = "shared"` — cualquier hook con `Domain == "shared"` es incluido por `ForDomains` sin importar el set de dominios activos, sin necesitar un directorio en el payload.

El catálogo de QUÉ entrega este tier vive en [`../shared/README.md`](../shared/README.md).

## Campos de manifest.json

| Campo | Significado |
| --- | --- |
| `id` / `label` | id del dominio (coincide con la carpeta) / etiqueta humana |
| `workspace` | dónde vive el estado del proyecto: `repository` \| `folder` |
| `alignmentArtifact` | el término del dominio para el spec/brief |
| `decisionRecord` | `{ term, dir }` — tipo de decision-record y ruta del store |
| `canonicalCatalog` | el catálogo que citan los decision records |
| `phases` | el pipeline de fases del dominio |
| `guards` | gates de verificación que corre el flujo |
| `explorationTool` | nombre del índice de exploración (opcional) |
| `mcp` | servidores MCP que el dominio registra (mapeados a pasos de instalación por nombre) |
| `binaries` | binarios/CLI que el dominio necesita instalados (p. ej. `engram`, `codegraph`, `proofshot`) |

**Nada es global.** Un MCP o binario se instala **solo** si algún dominio activo lo declara acá — no hay set "base" que se instale siempre. La instalación monta la **unión** de lo que declaran los dominios activos. Los **permisos** de Claude Code (`permissions.allow`) **no se declaran**: se infieren de `mcp` vía un map en Go (`name → mcp__<name>__*`, con override para casos como engram, que es un plugin); el único permiso global fijo es `Skill`.
