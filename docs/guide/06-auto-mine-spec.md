# 06 — Auto-mine de specs (`flagSpecMine`)

[← 05 Auto-mine de EDRs](05-auto-mine.md) · [Índice](README.md) · Siguiente: [07 — Herramientas →](07-herramientas.md)

> Esta página describe el auto-mine de **capability-specs**, el contrato observable del dominio **development** (ver [03](03-fases.md#spec-base)). El mecanismo de fondo (Mode A brownfield, gate obligatorio, `Inferred` como borrador no-confiable) es el mismo motor que el auto-mine de EDRs ([05](05-auto-mine.md)); esta página cubre dónde diverge para specs.

`mine` (spec) tiene **un solo modo**:

- **Mode A — scan brownfield**: escanea un repo (existente o con código sin specs) y propone capability-specs `Inferred` a partir del comportamiento as-built (rutas, validaciones, máquinas de estado, handlers de eventos, con tests como oráculo de confianza).

No hay Mode B in-flow. Esto es una **asimetría deliberada** frente al auto-mine de EDRs: una decisión puede quedar implementada sin EDR (nadie la capturó a propósito), pero el comportamiento de un cambio **siempre** queda capturado — lo autora la fase `spec`, que es **base obligatoria** del flujo SDD, no un add-on opt-in. No existe entonces un "hueco de spec" que detectar dentro del flujo; Mode A sirve para repos brownfield cuyo código es anterior al flujo SDD, o cuyas specs quedaron desactualizadas frente al código.

## El gate es la INTENCIÓN, no la presencia de specs

`flagSpecMine` es **opt-in, off por default**, hermano directo de `flagDecisionGaps` (mismo tipo, misma precedencia por-dominio). Con el flag off: silencio total. Con el flag on, el orquestador **ofrece** Mode A (en una línea, declinable en una palabra) cuando detecta un repo con código pero con `.matecito-ai/development-specs/` ausente o escaso — al iniciar sesión o después de `sdd-init`. Escasez = el store no existe, o sus specs cubren solo una fracción chica del código con comportamiento (muchas rutas / máquinas de estado / validaciones / handlers sin capability-spec).

Hay **dos confirmaciones distintas**, no las confundas: (1) la **oferta-a-escanear** — antes de cualquier scan; declinarla significa que no se escanea nada; y (2) el **gate de materialización** — después del scan, sobre los candidatos reales.

## El flujo

```
trigger     → orquestador detecta store ausente/escaso (flagSpecMine on, post-sdd-init o inicio de sesión)
oferta      → ofrece minar en una línea (1ª confirmación); declinable en una palabra — sin aceptar NO escanea
executor    → (solo si se acepta) escanea el repo (Mode A: rutas, validaciones, estados, eventos + tests como oráculo)
gate        → thread principal, obligatorio (2ª confirmación, materialización): confirmar/editar/saltar candidatos — nada se escribe sin confirmar
materialize → .matecito-ai/development-specs/<type>/<capability>.md   Status: Inferred
sdd-verify  → IGNORA el spec mientras sea Inferred (no es contrato, no puede fallar verify)
ratificar   → development-spec-bootstrap (modo update): revisa/corrige, agrega escenarios → Status: Accepted
sdd-verify  → a partir de ahí lo exige (SPEC-VIOLATION si el código diverge)
```

No hay hooks en `tasks`/`verify` ni boundary dispatch para specs — a diferencia del auto-mine de EDRs, que sí engancha ahí (ver [05](05-auto-mine.md#el-flujo-fase-por-fase)). El único disparador de spec-mine es su propio trigger post-`sdd-init`.

## `Inferred` es un borrador, no la verdad

Igual que un EDR `Inferred`, una capability-spec `Inferred` es un **borrador no-confiable**: describe el código **as-built**, bugs incluidos si los hay — no la intención. Por eso el gate es obligatorio (nada se materializa sin confirmar) y por eso, mientras esté `Inferred`, no es un contrato: `sdd-verify` la saltea.

## Distinto de los EDR Inferred: qué chequea `sdd-verify`

Acá diverge del auto-mine de EDRs, y vale la pena marcarlo porque es contraintuitivo:

- **EDR `Inferred`**: el chequeo de cumplimiento de EDRs en `verify` (6b) no filtra por `Status` — si un EDR (incluido uno `Inferred`) está en `.matecito-ai/edr/` y el cambio lo toca, `verify` **sí** confirma que el código lo honra.
- **Capability-spec `Inferred`**: el chequeo `SPEC-VIOLATION` (6d) está explícitamente **scoped a `Status: Accepted`** — un spec `Inferred` (como uno `Draft`) se saltea siempre, nunca puede disparar `SPEC-VIOLATION`.

La razón: una capability-spec `Inferred` calca el código as-built con precisión (rutas, validaciones, estados) — exigirla como contrato inmediatamente congelaría bugs existentes como si fueran comportamiento intencional. Un EDR `Inferred` es distinto: captura una decisión ya tomada (el QUÉ), y esa regla concreta vale como restricción aunque falte el porqué.

## Después: ratificar

Lo que mine materializa son **borradores `Inferred`** (comportamiento observado, sin intención confirmada). Para promoverlos a `Accepted`, corrés **`development-spec-bootstrap` en modo update**: revisás el contenido minado contra la intención real (¿es a propósito o es un bug?), sacás cualquier identificador interno que haya quedado, completás escenarios faltantes, y pasás `Status → Accepted`. Recién ahí `sdd-verify` empieza a exigirlo.

## Quién hace qué

- **executor (mine spec)**: mode-agnóstico, `scope → candidates[]`. No lee el flag, no escribe nada — igual que el miner de EDRs.
- **orquestador**: el único que dispara el executor (su propio trigger, no vía tasks/verify), presenta el gate, materializa lo confirmado.
- **`sdd-verify`**: audita `Accepted`, ignora `Inferred`.
- **`development-spec-bootstrap`**: ratifica `Inferred → Accepted` vía entrevista humana.
