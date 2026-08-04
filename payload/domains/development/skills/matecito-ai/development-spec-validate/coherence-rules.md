# Rúbrica de coherencia y completitud — capability-specs

Lista central de chequeos que aplica `development-spec-validate`. Es **ratchet-able**: cuando aparece una contradicción nueva, se agrega acá y queda cubierta para siempre.

Retiene solo la mitad **semántica** de la rúbrica: contradicción de significado, lenguaje vago, identificadores volátiles, clasificación por tipo/dominio y completitud-por-status. La mitad **mecánica** (secciones esqueleto, sincronía de índices, taxonomía de carpetas, links colgados, forma GIVEN/WHEN/THEN, reuso por un solo referrer, el eje `components`) vive en `~/.claude/references/artifact-checks/checks.yaml` y la evalúa `~/.claude/scripts/validate-artifact.js --type capability-spec --store .matecito-ai/development-specs`; esta skill invoca ese script e ingiere su JSON en vez de re-derivar esos chequeos — ver `SKILL.md`.

## Cómo la lee el validador

Cada chequeo tiene: **severidad** (CRITICAL / WARNING / SUGGESTION), una **condición** evaluada sobre los specs, y un **mensaje** (qué/por qué/sugerencia). El validador evalúa las condiciones contra `.matecito-ai/development-specs/<type>/` y reporta las que se cumplen. El tipo de cada spec sale de su carpeta; el status, del header.

---

## Completitud

- **[NOTA — Draft]** Specs con `Status: Draft` NO cierran el comportamiento: no reportes secciones esqueleto ni escenarios faltantes como defecto (esperados en Draft). Sí aplican los chequeos de coherencia contra los `Accepted`.
- **[NOTA — Inferred]** Specs con `Status: Inferred` se tratan como `Draft` para completitud: no reportes secciones esqueleto ni escenarios faltantes como defecto (es un borrador no-confiable minado del código as-built por `development-spec-mine`, se espera que le falten hasta la ratificación humana). Sí aplican los chequeos de coherencia contra los `Accepted`, pero con severidad capada — ver "Coherencia entre capabilities".

## Verificabilidad

- **[WARNING]** Una regla de negocio, rama o caso borde enunciado en prosa que no tiene ningún escenario que lo cubra → afirmación no testeable.
- **[SUGGESTION]** Lenguaje vago ("debería", "en lo posible", "idealmente", "normalmente", "evitar cuando se pueda") en una regla o comportamiento de un spec `Accepted` → el comportamiento tiene que ser determinista.

## Coherencia entre capabilities (el núcleo del validador)

- **[CRITICAL]** Dos specs describen el **mismo comportamiento de forma contradictoria** (ej: un `flow` dice que ante X el sistema responde A, y otro spec dice que ante X responde B).
- **[CRITICAL]** Una regla (`rule`) **prohíbe** lo que un `flow`/`process` **hace** (o al revés) → la regla y el flujo se contradicen.
- **[CRITICAL]** Un escenario de una capability asume un comportamiento que el escenario de otra capability contradice.
- **[EXCEPCIÓN — Inferred]** Si alguna de las dos capabilities en contradicción tiene `Status: Inferred`, la severidad de las tres reglas CRITICAL de arriba se **capa en 🟡 WARNING** ("posible drift as-built vs intención"), NUNCA en 🔴 CRITICAL: un `Inferred` es un borrador no-confiable minado del código, no una afirmación validada por una persona — escalarlo a CRITICAL penalizaría al humano por algo que la máquina minó, no que alguien afirmó.
- **[WARNING]** Dos specs describen la **misma capability** con distinto nombre (duplicado) → consolidar en uno.
- **[WARNING]** Un `flow`/`process` referencia una **entidad o estado** que ningún spec `lifecycle` (ni ninguna sección "Entidades y estados") define → estado colgado.
- **[WARNING]** Un spec referencia una **transición de estado** que el `lifecycle` de esa entidad no contempla.

## Referencias

- **[SUGGESTION]** Un comportamiento claramente gobernado por una decisión técnica (ej: una política de reintentos, un formato de error) no linkea ningún EDR → puede faltar el EDR o la referencia. (Solo sugerencia: no todo comportamiento tiene un EDR.)

## Vocabulario (separación qué-hace vs cómo)

- **[WARNING]** Un spec nombra **identificadores internos volátiles** (clases, métodos, columnas de base de datos, errores internos, rutas de archivo) en cualquier sección → el spec es el *qué hace*, en idioma de dominio + contrato público; el *cómo* es del código y el *por qué* del EDR. Excepción: nombre de tecnología/librería y contrato público (endpoints públicos, códigos de error expuestos). Ver `~/.claude/references/spec/README.md` → "No es el cómo".
- **[SUGGESTION]** Un spec incluye justificación/argumentación de por qué se eligió un enfoque técnico → eso es un EDR; el spec especifica, no argumenta.

## Higiene de status

- **[SUGGESTION]** Un spec `Draft` de larga data referenciado por el flujo como fuente de verdad → conviene completarlo a `Accepted` o quitar la dependencia.
- **[SUGGESTION]** Un spec `Inferred` de larga data sin ratificar → conviene revisarlo con `development-spec-bootstrap` (modo update, caso "Ratificar un Inferred") antes de que el código diverja más del candidato minado.

## Integridad de la taxonomía

- **[WARNING]** El contenido de un spec no corresponde a su carpeta-tipo según la **regla de clasificación** (`~/.claude/references/spec/README.md` → «Cómo clasificar el tipo»): ej. un archivo en `rule/` que en realidad describe un flujo disparado por un actor con pasos y ramas → moverlo al tipo correcto.
