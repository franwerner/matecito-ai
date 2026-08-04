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

- Última actualización: 2026-08-03 (archive `contract-template-self-check`) — Modificada `process/validate-artifact-structure` (alcance ampliado a validar contratos/templates contra sí mismos con `--self-check`, output pasado de texto a JSON)
- **Definir una capacidad nueva:** usá la skill `development-spec-bootstrap` (o escribí el spec desde `~/.claude/references/spec/templates/capability.md`); creá la carpeta del tipo si no existía y sumá la fila al `INDEX.md` de ese tipo (y a este índice raíz si el tipo es nuevo en el proyecto).
- **Actualizar comportamiento (cambio menor):** editá el spec. El historial lo lleva git.
- **Cambio de comportamiento vía flujo SDD:** no edites el spec a mano — el delta del cambio se mergea acá al archivar (`sdd-archive`).
- **Retirar una capacidad:** marcá el spec `Deprecated` con link a su reemplazo; no borres el archivo.
- **Validar coherencia entre specs:** usá la skill `development-spec-validate`.
