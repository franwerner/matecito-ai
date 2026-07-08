<!-- Canonical template: índice RAÍZ de capability-specs (`.matecito-ai/development-specs/INDEX.md`). Enruta por tipo; el detalle de cada capability está en el INDEX de su tipo. Consumido por development-spec-bootstrap (materialización) y actualizado por sdd-archive al consolidar un delta. -->

# Capability specs — Índice raíz

El comportamiento del sistema, capturado por **capacidad** y organizado por **tipo**. Cada capability-spec dice *qué hace* el sistema; el *por qué* de cada elección técnica vive en `../edr/`, y el *cómo* literal en el código.

## Cómo usar este índice

1. Identificá qué tipo de comportamiento vas a tocar.
2. Encontrá el tipo abajo y abrí su `INDEX.md`.
3. Leé los capability-specs relevantes **antes** de escribir código o tests.
4. Si hay contradicción entre tu plan y un spec `Accepted`: pará y preguntale al usuario.

## Tipos de este proyecto

(Solo se listan los tipos con al menos un spec-archivo.)

| Tipo | Qué agrupa | Índice |
|---|---|---|
| `flow` | Operaciones de cara a un actor, con pasos y ramas (connect, send, receive) | [flow/INDEX.md](flow/INDEX.md) |
| `rule` | Reglas de negocio transversales, sin flujo (ventanas, dedup, políticas) | [rule/INDEX.md](rule/INDEX.md) |
| `lifecycle` | Máquinas de estado de una entidad (su ciclo de vida y transiciones) | [lifecycle/INDEX.md](lifecycle/INDEX.md) |
| `process` | Comportamiento reactivo/de fondo, disparado por evento (reconciliación, jobs) | [process/INDEX.md](process/INDEX.md) |

**Leyenda de status:** `Draft` = en escritura, no es fuente de verdad todavía · `Accepted` = ratificado, el código se valida contra él · `Deprecated` = capacidad retirada/reemplazada (se conserva por trazabilidad).

## Tipos sin uso en este proyecto

(Tipos sin ningún spec-archivo — no tienen carpeta. Se listan acá para dejar constancia de que se consideraron.)

| Tipo | Razón |
|---|---|
| `<tipo>` | <por qué este proyecto no tiene comportamiento de este tipo todavía> |

## Estado y mantenimiento

- Última actualización: <YYYY-MM-DD>
- **Definir una capacidad nueva:** usá la skill `development-spec-bootstrap` (o escribí el spec desde `~/.claude/references/spec/templates/capability.md`); creá la carpeta del tipo si no existía y sumá la fila al `INDEX.md` de ese tipo (y a este índice raíz si el tipo es nuevo en el proyecto).
- **Actualizar comportamiento (cambio menor):** editá el spec. El historial lo lleva git.
- **Cambio de comportamiento vía flujo SDD:** no edites el spec a mano — el delta del cambio se mergea acá al archivar (`sdd-archive`).
- **Retirar una capacidad:** marcá el spec `Deprecated` con link a su reemplazo; no borres el archivo.
- **Validar coherencia entre specs:** usá la skill `development-spec-validate`.
