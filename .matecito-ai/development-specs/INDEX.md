# Capability specs — Índice raíz

El comportamiento del sistema, capturado por **capacidad** y organizado por **tipo**. Cada capability-spec dice *qué hace* el sistema; el *por qué* de cada elección técnica vive en `../edr/`, y el *cómo* literal en el código.

**Alcance:** el store cubre las dos mitades del producto — el **cockpit** (el broker/MCP de `apps/api` y la UI de `apps/ui`) y el **CLI** de matecito-ai (instalación, actualización, despliegue del payload en el host destino, configuración por dominio, hooks y chequeo del entorno).

## Cómo usar este índice

1. Identificá qué tipo de comportamiento vas a tocar.
2. Encontrá el tipo abajo y abrí su `INDEX.md`.
3. Leé los capability-specs relevantes **antes** de escribir código o tests.
4. Si hay contradicción entre tu plan y un spec `Accepted`: pará y preguntale al usuario.

## Tipos de este proyecto

(Solo se listan los tipos con al menos un spec-archivo.)

| Tipo | Qué agrupa | Índice |
|---|---|---|
| `flow` | Operaciones de cara a un actor, con pasos y ramas (submit de fase, iniciar change, instalar, configurar) | [flow/INDEX.md](flow/INDEX.md) |
| `rule` | Reglas de negocio transversales, sin flujo (scoping de eventos, activación de dominios, precedencias de configuración, políticas) | [rule/INDEX.md](rule/INDEX.md) |
| `lifecycle` | Máquinas de estado de una entidad (su ciclo de vida y transiciones) | [lifecycle/INDEX.md](lifecycle/INDEX.md) |
| `process` | Comportamiento reactivo/de fondo, disparado por evento o por el sistema (file-watch, indexado, versionado, despliegue, validadores del host) | [process/INDEX.md](process/INDEX.md) |

**Leyenda de status:** `Inferred` = borrador no-confiable minado del código as-built, pendiente de ratificación humana · `Draft` = en escritura, no es fuente de verdad todavía · `Accepted` = ratificado, el código se valida contra él · `Deprecated` = capacidad retirada/reemplazada (se conserva por trazabilidad).

## Estado y mantenimiento

- Última actualización: 2026-08-14 (archive `sdd/contract-shape-gate`) — Creada `rule/contract-shape-proposal` (contrato como item compuesto tipado, ratificación via instrucciones de re-despacho, test de alcance angostado); Extendida `flow/ratify-gate-items` (cuarta forma compound-item, excepción de ritmo para contratos múltiples); Modificada `rule/mailbox-item-summary-rationale-split` (items pueden llevar lista de campos tipados bajo summary, 16 secciones en lugar de 12, caps per-field de 160 caracteres)
- Anterior: 2026-08-14 (archive `sdd/gate-coverage-gaps`) — Extendida `flow/ratify-gate-items` (seis momentos de orquestador, regla de conteo 0/1/≥2, excepción de anchor en Discovery Gate, risks informativos); Creada `rule/decision-re-emergence-reconfirmation` (decisiones ratificadas que resurgen, ledger per-change, reconfirmación corta); Modificada `rule/mailbox-item-summary-rationale-split` (split declarado por sección, formas tabla/labeled-list igual que listas, 12 secciones en lugar de 9)
- Anterior: 2026-08-13 (archive `sdd/sequential-decision-gate`) — Creada `flow/ratify-gate-items` (gate uno-a-uno con evidencia anclada); Modificada `rule/mailbox-item-summary-rationale-split` (3 partes: summary+anchor+rationale, tope al summary); Modificada `flow/materialize-mined-artifacts` (candidatos por template compartido)
- 2026-08-13 (archive `sdd/worktree-isolation`) — Creada `process/isolate-change-workspace` (aislamiento anidado de workspace: cambio en nivel propio, tarea dentro del cambio)
- **Definir una capacidad nueva:** usá la skill `development-spec-bootstrap` (o escribí el spec desde `~/.claude/references/spec/templates/capability.md`); creá la carpeta del tipo si no existía y sumá la fila al `INDEX.md` de ese tipo (y a este índice raíz si el tipo es nuevo en el proyecto).
- **Actualizar comportamiento (cambio menor):** editá el spec. El historial lo lleva git.
- **Cambio de comportamiento vía flujo SDD:** no edites el spec a mano — el delta del cambio se mergea acá al archivar (`sdd-archive`).
- **Retirar una capacidad:** marcá el spec `Deprecated` con link a su reemplazo; no borres el archivo.
- **Validar coherencia entre specs:** usá la skill `development-spec-validate`.
