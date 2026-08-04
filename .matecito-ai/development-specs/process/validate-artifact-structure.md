# Capability — Validar conformidad estructural de un artefacto escrito

- **Status:** Accepted
- **Date:** 2026-08-03
- **Components:** cli

## Propósito

Validar que un artefacto durable ya escrito (EDR o capability-spec) se conforma al contrato declarado de su tipo — secciones presentes y ordenadas, presencia condicionada a Status, valores de cabecera dentro de sus enums. El validador examina exactamente el alcance que se le nombra — un archivo (`--file`), el store completo de un tipo (`--store`), o los contratos y templates embarcados (`--self-check`) — modifica nada y reporta hallazgos.

## Actores

- Chequeos de coherencia en el flujo SDD (desarrollo-spec-validate)
- Validación manual ad-hoc de conformidad estructural
- Sistemas de CI/CD que detectan archivos no conformes ante hand-edits en git

## Precondiciones

- El artefacto a validar existe y es legible
- El contrato del tipo de artefacto existe

## Flujo principal

### Modo archivo (--file)

1. Leer el contrato para el tipo de artefacto
2. Leer el archivo a validar
3. Parsear el archivo (YAML header + markdown body)
4. Para cada campo de cabecera declarado:
   - Validar que existe si es requerido; registrar error si no
   - Validar que su valor está en el enum (si tiene enum); registrar error si no
5. Para cada sección declarada en el contrato:
   - Validar que está presente conforme a `emitted` y `Status`; registrar warning/nota si no
   - Validar que el orden es el esperado; registrar nota si está desordenada
6. Reportar todos los hallazgos (error/warning/nota) como un registro JSON único
7. Salir 0 si no hay errores; salir 1 si hay errores (warnings/notas no cierran el proceso)

### Modo store (--store)

1. Leer el contrato para el tipo de artefacto
2. Recorrer cada archivo del store
3. Para cada archivo, ejecutar el ciclo de validación del modo archivo (pasos 2-5)
4. Realizar chequeos entre archivos:
   - Validar que el índice de la carpeta está sincronizado (cada artefacto tiene una entrada; cada entrada apunta a un artefacto existente)
   - Validar que los links de secciones Referencias apuntan a artefactos existentes (no links colgados)
   - Validar que las carpetas dentro del store son canónicas (en la taxonomía declarada) o explícitamente permitidas
5. Acumular todos los hallazgos de todos los archivos
6. Reportar hallazgos agrupados por archivo, con resumen de conteos por tipo, en un registro JSON único
7. Salir 0 si no hay errores; salir 1 si hay errores; un archivo no conforme NO aborta el barrido

### Modo self-check (--self-check)

1. Leer los dos contratos (edr.yaml, capability.yaml)
2. Leer los dos templates (edr.md, capability.md)
3. Leer el archivo checks.yaml
4. Validar cada par contrato↔template:
   - Que cada `title:` declarado aparece como encabezado en el template
   - Que cada `##` encabezado del template está declarado en `sections[]`
   - Que cada valor en `statuses:` aparece en el enum de estados del template
5. Validar que cada fila de checks.yaml referencia secciones y carpetas que existen en sus contratos
6. Validar que los enlaces de ejemplo en ambos templates derivan su `../` depth de los `path:` declarados en los contratos
7. Reportar todos los hallazgos en un registro JSON único
8. Salir 0 si no hay drift; salir 1 si hay drift; salir 2 si un contrato, template o checks.yaml no es legible

## Ramas / flujos alternativos

- **Archivo no conforme (sección faltante, orden incorrecto, status ilegal)** → Registrar hallazgo(s); continuar checando; salir 0 si solo warning/nota, salir 1 si algún error
- **Tipo de artefacto no soportado** → Fallar a stderr; salida 2
- **Barrido de store sin permisos para la raíz (`--root`)** → Si un chequeo entre archivos requiere acceso al root (globs, configuración del repo) y `--root` no se proporcionó, salida 2 nombrando los chequeos bloqueados
- **Artefacto no conforme durante el barrido** → No aborta el barrido; continúa con el siguiente archivo; se reporta como hallazgo individual
- **Self-check: contrato, template o checks.yaml no legible** → Salida 2; stdout vacío; stderr nombra el archivo no legible

## Casos borde

- **Archivo editado a mano con secciones adicionales no declaradas** → Registrar hallazgo `UNKNOWN-SECTION`; continuar; no fallar por solo eso
- **Archivo con orden de secciones alterado** → Registrar nota `SECTION-ORDER`; no considerar falta de conformidad total
- **EDR con Status ilegal** → Registrar error `STATUS-ILLEGAL`; salida 1
- **Capability-spec con `Status: Deprecated` sin línea `**Reemplazado por:**`** → Registrar warning o error según severidad declarada en contrato
- **Validación conforme, sin hallazgos** → Registro JSON con `findings: []`, `counts: {error: 0, warning: 0, nota: 0}`; salida 0
- **Validación con solo hallazgos warning/nota** → Registro JSON con hallazgos listados en `findings[]`, `counts.error: 0`; salida 0
- **Validación con hallazgos error** → Registro JSON con hallazgos listados en `findings[]`, `counts.error > 0`; salida 1

## Reglas de negocio

- El validador examina exactamente el alcance que se le nombra — un archivo (`--file`), el store completo de un tipo (`--store`), o los contratos y templates embarcados (`--self-check`) — nunca escribe y nunca lo amplía por su cuenta
- Tres niveles de severidad en hallazgos: `error` (rompe conformidad, salida 1), `warning` (se reporta, no rompe, herencia futura hacia error), `nota` (informativo; se reporta, no rompe)
- La severidad de `SECTION-MISSING` depende de la declaración de la sección en el contrato:
  - Sección con `emitted: always` → ausencia es `error`
  - Sección con `emitted: on-status <status>` cuyo status aplica → ausencia es `warning`
  - Sección con `when-present` → nunca es hallazgo por ausencia (no es aplicable)
- La conformidad se representa en el contrato de salida: un registro JSON único con `findings[]` (vacío si conforme) y `counts.error` (cero si conforme)
- El contrato es la fuente de verdad; un archivo que el renderer produjo siempre conforma (por construcción)

## Entidades y estados

- **Artefacto escrito** — Un archivo `.md` con cabecera YAML. Estados: conforme → no conforme (detectado); no escrito → escrito (el validador lee, no escribe)

## Errores de cara al actor

El validador emite un único registro JSON en stdout, con este esquema:

```json
{
  "tool": "validate-artifact",
  "type": "<type-of-artifact>",
  "mode": "<file | store | self-check>",
  "root": "<project-root-or-null>",
  "target": "<path-or-store-dir-or-null>",
  "scanned": <count>,
  "findings": [
    {
      "severity": "<error | warning | nota>",
      "display": "<descriptive-title>",
      "code": "<ERROR-CODE>",
      "file": "<path-to-artifact-or-reference>",
      "line": <line-number-or-null>,
      "related": [<related-paths>],
      "message": "<detailed-message>"
    }
  ],
  "counts": {
    "error": <count>,
    "warning": <count>,
    "nota": <count>
  },
  "skipped": [<paths-skipped-if-any>]
}
```

**Interpretación de la salida:**
- Conformidad: derivada de `counts.error === 0` (no hay campo `conformant`)
- Hallazgos reportados: en el array `findings[]` agrupados por tipo en su order de aparición
- Exit 0: `counts.error === 0` (sin errores, warnings/notas no previenen salida 0)
- Exit 1: `counts.error > 0` (hay al menos un error)
- Exit 2: no legible (stdout vacío, stderr contiene razón) — usado cuando un contrato, template o checks.yaml no se puede leer

## Escenarios

### Scenario: Archivo conforme

- **GIVEN** un archivo que el renderer produjo
- **WHEN** se valida contra su contrato con `validate-artifact.js --type <type> --file <path>`
- **THEN** se emite un registro JSON con `findings: []`, `counts.error: 0`; salida 0

### Scenario: Desviaciones reportadas individualmente

- **GIVEN** un archivo que falta una sección declarada, tiene orden incorrecto, o un valor de cabecera ilegal
- **WHEN** se valida
- **THEN** cada desviación es un hallazgo distinto en `findings[]`, nombrando sección o campo; salida 0 si son warnings/notas, salida 1 si hay errores

### Scenario: Solo el archivo nombrado es examinado

- **GIVEN** un store con muchos artefactos
- **WHEN** se valida un solo path con `validate-artifact.js --type <type> --file <path>`
- **THEN** solo ese archivo es leído; `scanned: 1`; nada es escrito; no hay barrido de directorio

### Scenario: Status ilegal rompe conformidad

- **GIVEN** un archivo con `Status:` a un valor fuera del enum (e.g., `Rejected`)
- **WHEN** se valida
- **THEN** se reporta un hallazgo `ERROR STATUS-ILLEGAL` en `findings[]`; `counts.error: 1`; salida 1

### Scenario: Archivo conforme con solo hallazgos no-error

- **GIVEN** un archivo con sección desordenada (nota) o warning menor
- **WHEN** se valida
- **THEN** hallazgos se reportan en `findings[]`; `counts.error: 0` pero `counts.warning > 0` o `counts.nota > 0`; salida 0

### Scenario: Línea de cierre precisa en hallazgos

- **GIVEN** un archivo conforme
- **WHEN** se valida
- **THEN** `findings: []` si cero hallazgos; `findings` listado si los hay; nunca un campo booleano de conformidad

### Scenario: Barrido del store completo

- **GIVEN** un store con varios artefactos, alguno no conforme
- **WHEN** se valida el store entero con `validate-artifact.js --type <type> --store <dir> --root <raíz>`
- **THEN** se examina cada artefacto; `scanned` cuenta el total; cada hallazgo lleva su `file`; uno no conforme no aborta el barrido; todos los hallazgos se acumulan en `findings[]`; salida 1 si hay errores, salida 0 si solo warnings/notas

### Scenario: Chequeos entre archivos

- **GIVEN** un índice desincronizado (un artefacto sin entrada en INDEX.md o una entrada que apunta a un archivo inexistente), un link colgado en la sección Referencias, o una carpeta no canónica dentro del store
- **WHEN** se barre el store
- **THEN** cada anomalía es un hallazgo propio en `findings[]` que nombra los archivos o carpetas involucrados; esos hallazgos se reportan además de los por-archivo

### Scenario: Shipped contracts agree

- **GIVEN** los contratos y templates tal como los deja este cambio
- **WHEN** `--self-check` se ejecuta
- **THEN** no hay hallazgos (`findings: []`, `counts.error: 0`); salida 0

### Scenario: An input cannot be read

- **GIVEN** un contrato, template o checks.yaml no legible
- **WHEN** `--self-check` se ejecuta
- **THEN** salida 2; stdout vacío; stderr nombra el archivo no legible

### Scenario: Declared section missing from its template

- **GIVEN** una sección cuyo título está en `contract.sections[]` pero no aparece como encabezado `##` en el template correspondiente
- **WHEN** `--self-check` se ejecuta
- **THEN** un hallazgo `error DRIFT-SECTION` nombrando la sección y su contrato; salida 1

### Scenario: Template heading the contract does not declare

- **GIVEN** un encabezado `##` que aparece en el template pero no está declarado en `contract.sections[]`
- **WHEN** `--self-check` se ejecuta
- **THEN** un hallazgo `error DRIFT-SECTION` nombrando el encabezado; salida 1

### Scenario: A row references a renamed section

- **GIVEN** una fila de checks.yaml cuyo `section:` es un título que su contrato no declara
- **WHEN** `--self-check` se ejecuta
- **THEN** un hallazgo `error DRIFT-CHECKS-SECTION` nombrando el código de fila y el título sin concordancia; salida 1

### Scenario: A row references a folder outside the taxonomy

- **GIVEN** una fila nombrando una carpeta ausente de las `folders` y `extra_folders` declaradas de su store
- **WHEN** `--self-check` se ejecuta
- **THEN** un hallazgo `error DRIFT-CHECKS-FOLDER` nombrando el código de fila y la carpeta desconocida; salida 1

### Scenario: Cross-store link one level short

- **GIVEN** un capability.md que linkea a `../edr/<domain>/<slug>.md`, mientras los `path:` declarados ponen un spec dos niveles bajo `.matecito-ai/`
- **WHEN** `--self-check` se ejecuta
- **THEN** un hallazgo `error DRIFT-LINK-DEPTH` nombrando el link y la profundidad esperada; salida 1

### Scenario: Same-store link

- **GIVEN** un edr.md cuyo "## Relacionados" linkea a `../<domain>/<slug>.md`
- **WHEN** `--self-check` se ejecuta
- **THEN** la profundidad coincide con lo que edr.yaml declara; sin hallazgo

### Scenario: Self-check reads no project artifact

- **GIVEN** un repo cuyos stores están ausentes o son no-conformes
- **WHEN** `--self-check` se ejecuta
- **THEN** solo los contratos, templates y checks.yaml se leen; no se barre ningún store; `scanned` cuenta solo las referencias (2: edr y capability)

## Referencias

- **EDR** → [`../../edr/contracts/three-level-checker-severity.md`](../../edr/contracts/three-level-checker-severity.md) — Decisión de error/warning/nota y salida 1 solo en error
- **EDR** → [`../../edr/structure/two-scripts-render-and-validate.md`](../../edr/structure/two-scripts-render-and-validate.md) — Decisión de script independiente
