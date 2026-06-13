# 04 — Decisiones y ADRs

[← 03 Las fases](03-fases.md) · [Índice](README.md) · Siguiente: [05 — Auto-mine de ADRs →](05-auto-mine.md)

La capa de decisiones es lo que hace que el agente respete las convenciones del proyecto en vez de reinventarlas. Se apoya en **ADRs** y en tres skills: `bootstrap`, `validate` y `mine`.

> **Decision record por dominio.** "Capturar las decisiones una vez y respetarlas" es un mecanismo del **núcleo**; cada dominio define su *tipo* de decision record. En **development** es el **ADR** (Architecture Decision Record), que es lo que describe esta página. En otros dominios cambia el nombre y el catálogo (p. ej. **DDR** —Design Decision Record— en design), pero el ciclo de vida y la tríada productor/consumidor son los mismos.

## Qué es (y qué no es) un ADR

Un **ADR** captura una **decisión** de ingeniería: una elección deliberada entre alternativas, con su razón. Tres rasgos lo definen:

1. **Decide** algo (elige una opción frente a otras),
2. **tiene un porqué** (contexto + trade-off),
3. **perdura como restricción** (el código futuro la respeta; se chequea su drift).

Un ADR **no es**: una tarea, un criterio de verificación, una señal/TODO, código incidental, ni un "cómo" de implementación. Y el **porqué nunca se adivina** — lo aporta una persona.

> La definición canónica y agnóstica de flujo vive en [`references/adr/README.md`](../../payload/references/adr/README.md); la **estructura** del archivo está en `references/adr/templates/`. Cualquier skill la consume; no la redefine.

**Estados** (un ADR no es estático): `Inferred` (borrador: el QUÉ, sin el porqué) → `Accepted` (ratificado por una persona, con el porqué); más `Pending`/`Deferred` (aplazados) y `Superseded` (reemplazado).

Dónde viven: `.matecito-ai/adr/<dominio>/<slug>.md`, agrupados por dominio, con índices (`INDEX.md` raíz + por dominio). **Nunca en Engram.**

## Concerns vs ADR

Se confunden fácil. La regla corta: **un concern es la PREGUNTA; un ADR es la RESPUESTA.**

| | Concern | ADR |
|---|---|---|
| Qué es | un tema que suele requerir decisión (la pregunta) | la decisión capturada (la respuesta) |
| Dónde vive | catálogo de `bootstrap` (`concerns/`) | el proyecto (`.matecito-ai/adr/`) |
| Alcance | compartido, todos los repos | este repo |
| Naturaleza | guía: qué preguntar + defaults | registro: qué se decidió + por qué |

El puente: un concern, cuando un proyecto lo trata, **produce un ADR**. No es 1:1 — un concern puede quedar `Not Applicable` (sin ADR), y un ADR puede no tener concern (decisión fuera del catálogo → custom local + flag de catalog-gap).

## La tríada

### bootstrap — producir decisiones (Accepted)

Entrevista por fases. Recorre los concerns relevantes al proyecto, **propone** opciones con default y un porqué según el stack (recommend-and-confirm: vos aceptás, corregís u override), y materializa **ADRs `Accepted`** (con el porqué que confirmaste). Es el **productor** de decisiones. También tiene **modo update**: resuelve `Pending`/`Deferred`, y **ratifica `Inferred` → `Accepted`** (ver abajo).

### validate — chequear coherencia

Validador **consultivo** (no modifica nada): lee `.matecito-ai/adr/` y reporta coherencia entre ADRs, completitud y verificabilidad, con severidad. Trata un `Inferred` como **no-decidido** (no cuenta para completitud; no marca su porqué vacío como defecto). Los hallazgos los resuelve el humano vía bootstrap modo update.

### mine — descubrir desde código (Inferred)

Mina decisiones implícitas en el código de un repo existente y las **propone** como ADRs `Inferred` (borradores). **No es productor de decisiones: descubre y propone.** Observa el **QUÉ** (visible en el código) y deja el **PORQUÉ vacío** — no lo infiere, porque no hay humano en el loop al descubrir. Ver [05](05-auto-mine.md) para su funcionamiento en el flujo.

## El ciclo de vida de una decisión

```
mine (desde código)  →  Inferred (.md: QUÉ + evidencia, sin porqué)
                              │  acto humano deliberado, cuando querés ratificar
bootstrap modo update  →  entrevista el PORQUÉ → completa secciones
                          → descarta `## Evidencia (inferida)` → Accepted
validate  →  chequea coherencia/completitud (consultivo, en cualquier momento)
```

- **mine propone**, **bootstrap produce/ratifica**, **validate chequea**.
- La promoción `Inferred → Accepted` la hace **siempre una persona** vía bootstrap; nunca es automática. Sin la entrevista del porqué, un Inferred no se convierte en decisión.

## Productores vs consumidores (y sus gates)

- **Consumidores** de ADRs (los leen para chequear): `intake` (guardia), `design` (alinea), `verify` (cumplimiento). Gateados por **presencia**: si no hay `.matecito-ai/adr/`, no hay contra qué chequear → silencio.
- **Productores/proponentes**: `bootstrap` y `mine`. Gateados por **intención** (lo invocás / el flag), no por presencia — pueden **bootstrapear los primeros ADRs desde cero**.

> El **ADR activation gate** (presencia de `.matecito-ai/adr/` con contenido) decide si los **consumidores** actúan. No condiciona a los productores.
