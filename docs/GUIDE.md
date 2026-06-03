# Guía: agregar una nueva dependencia / integración

Esta guía describe el **flujo** (no el código) a seguir cada vez que se incorpora una
nueva dependencia o integración al ecosistema matecito-ai. El objetivo es que toda
integración futura siga los mismos pasos y no se saltee decisiones ni validaciones.

Ejemplo de referencia a lo largo de la guía: la integración de **ProofShot**
(CLI de verificación visual de UI). Donde diga "ProofShot", pensá en tu dependencia.

> Regla base: un cambio sustancial pasa por el **flujo SDD** (`intake → … → archive`),
> no por edición ad-hoc. Las decisiones se confirman con el usuario en los gates.

---

## Paso 0 — Clasificar la integración

Antes de tocar nada, definí dos ejes. Determinan el alcance y qué piezas tocás.

1. **¿Qué tipo de dependencia es?**
   - **MCP server** (ej: engram, codegraph, context7, drawio) → se registra como MCP
     (`claude mcp add`), entra en `permissions.allow`, y aparece en el Cross-check SDD ↔ MCP.
   - **CLI / binario suelto** (ej: ProofShot) → binary-only: **no** se registra como MCP,
     **no** va en `permissions.allow`, **no** entra al Cross-check. Tiene su sección propia
     en `verify`.

2. **¿Dónde impacta?**
   - **Comportamiento del flujo SDD** (cambia cómo se comportan los agentes `sdd-*`) →
     tocás los prompts en `payload/agents/` y los skills en `payload/skills/`.
   - **Solo el installer** (provisión del binario) → tocás `internal/setup/` (sync + install).
   - Puede ser **ambos** (ProofShot fue ambos: cambió `sdd-verify` y se sumó al installer).

3. **¿Cómo se instala?** Identificá el mecanismo (`npm install -g`, descarga de release,
   plugin de Claude Code, etc.) y **qué herramienta existente lo refleja** para copiar el
   patrón (ej: ProofShot es npm → se modeló como codegraph).

---

## Paso 1 — Verificar el contrato REAL de la dependencia

**No asumas.** Antes de escribir specs, confirmá contra la fuente/docs reales:

- Comandos y flags reales (incluido si hay flag para correr **no-interactivo**).
- Formato de salida que vas a consumir (¿qué expone? ¿por sesión o por ítem?).
- Mecanismo y nombre exacto del paquete / binario.
- Pasos post-install (¿descarga browsers? ¿prompts? ¿escribe archivos?).
- Concurrencia y estado (locks, directorios de salida, puertos).

> Gotcha ProofShot: se asumió selectores CSS y conteos de error por-scenario; el contrato
> real era refs de accessibility + agregados por-sesión. Verificarlo a tiempo evitó
> arrastrar supuestos equivocados al spec y al código.

---

## Paso 2 — Pasar por el flujo SDD

1. **Intake** (`sdd-intake`): estructura el pedido, clasifica, recomienda **lane**.
   - Bias: **lane mínima viable**. `reduced` por default para cambios sustanciales;
     `custom` (base + add-ons puntuales) si hay un trigger nombrado (decisión de diseño,
     varios dominios, superficie grande); `full` solo con trigger explícito.
   - **INTAKE GATE**: el alcance y la lane se confirman con el usuario. Siempre.
2. **Spec** (`sdd-spec`): requisitos + escenarios. Acá se formaliza todo contrato nuevo
   (schemas, formatos, comportamiento esperado).
3. **Design** (`sdd-design`, si la lane lo incluye): decisiones técnicas + diagrama si
   la complejidad estructural lo amerita.
4. **Tasks → Apply → Verify → Archive**: implementación en batches, validación contra
   spec, cierre.

> Las decisiones de scope se registran en Engram (`architecture/<tema>`) y los artefactos
> de cada fase bajo `sdd/<cambio>/<fase>`.

---

## Paso 3 — Installer (plan de `install` / `update`)

`install` y `update` comparten el **mismo motor** (`internal/setup/sync` → `Detect` →
`PlanSync` → `Sync`). Si la sumás bien, queda cubierta por **ambos** comandos automáticamente.

Pasos de flujo:

1. **Detección de versión**: agregar cómo se obtiene la última versión (ej: npm registry)
   y la instalada (`<bin> --version`), espejando una herramienta del mismo tipo.
2. **Inclusión en el plan**: que aparezca en `Detect()` como un componente más y en el
   dispatch de `Sync()`. Verificá que en `--dry-run` se liste (`<tool> — instalar/actualizar`).
3. **Instalación**: la función de install espeja el patrón del tipo (npm → como codegraph).
   - Si hay **post-step** (ej: `proofshot install`): debe correr **no-interactivo** y
     **no colgar** (stdin cerrado / timeout). Si falla, **no abortar** el resto del plan:
     seguir y mostrar instrucción manual (política "se continúa ante errores de binarios").
4. **MCP (solo si aplica)**: registrar el MCP y sumar el patrón a `permissions.allow`.
   Para un CLI suelto, **saltear** este punto.

---

## Paso 4 — Capability detection + `verify`

1. **Check propio**: agregar `internal/checks/<tool>/` (detección PATH-based + versión).
2. **Sección en `verify`**: cablear el check en **las dos superficies** —
   `internal/cli/verify.go` y `internal/tui/screens/verify/verify.go`— para que tenga su
   **sección propia** en el reporte. (Mantener CLI y TUI en paridad.)
3. **Si los agentes SDD dependen de la herramienta**: `sdd-init` detecta la capability y la
   cachea en `sdd/<project>/testing-capabilities`; las fases que la usan **saltean en
   silencio** si está ausente (nunca rompen).

> Nota: un CLI suelto NO va al Cross-check SDD ↔ MCP (eso valida solo MCPs). Su estado
> de instalación se ve en su sección propia.

---

## Paso 5 — PATH y distribución

- **PATH**: si la herramienta instala su binario en un dir que puede no estar en el PATH
  (ej: `~/.npm-global/bin`), asegurate de que el installer lo agregue al shell rc de forma
  **idempotente**. No asumas que el PATH ya está configurado.
- **Caveat de distribución (importante)**: el `install`/`update` se **auto-actualiza
  bajando la release publicada**. Mientras los cambios no estén en una release, un build
  local NO se propaga vía self-update; hay que instalar el binario a mano en el entorno de
  prueba, y un `install` posterior lo pisará con la release. Documentar/avisar esto.

---

## Paso 6 — Testing

1. **Automático**: `go build ./...`, `go vet ./...`, `go test ./...`. Agregar/extender tests
   espejando los existentes del componente del mismo tipo.
2. **End-to-end en un entorno limpio**: probar en un usuario de test (`tester`) — instalar,
   confirmar que:
   - la herramienta queda **instalada y en el PATH**,
   - `verify` muestra su **sección propia** correctamente,
   - el payload/agentes que la usan quedaron deployados.
3. Si el cambio afecta el TUI, validar la salida/estado de la terminal al salir.

---

## Paso 7 — Archive

Cerrar el ciclo SDD (`sdd-archive`): persistir el estado final, dejar follow-ups anotados.
Los cambios quedan en el working tree; el commit/PR los maneja el usuario (atómicos por cambio).

---

## Checklist rápido

- [ ] Clasificada: ¿MCP o CLI? ¿toca agentes SDD, installer, o ambos? ¿qué herramienta refleja?
- [ ] Contrato real verificado (flags, salida, paquete, post-install, no-interactividad).
- [ ] Pasó por SDD con lane confirmada en el INTAKE GATE.
- [ ] Detección de versión + inclusión en el plan (`Detect` + dispatch). Aparece en `--dry-run`.
- [ ] Install espeja el patrón del tipo; post-step no-interactivo y con fallo aislado.
- [ ] MCP registrado + permisos (solo si es MCP); salteado si es CLI suelto.
- [ ] Check propio + sección en `verify` (CLI **y** TUI en paridad).
- [ ] Capability en `sdd-init`/`testing-capabilities` si los agentes la usan (silent-skip si falta).
- [ ] PATH asegurado idempotente; caveat de self-update/release documentado.
- [ ] `go build/vet/test` en verde + prueba end-to-end en usuario limpio.
- [ ] Ciclo SDD archivado.

---

## Gotchas aprendidos (ProofShot)

- **Verificá el contrato externo antes de specear**: supuestos equivocados (selectores,
  formato de errores) se arrastran al código si no se chequean a tiempo.
- **`install` y `update` comparten motor**: sumando bien al `sync`, ambos quedan cubiertos.
- **Post-step que cuelga**: correr no-interactivo (stdin cerrado/timeout) y aislar el fallo.
- **CLI ≠ MCP**: no forzar un CLI suelto dentro del Cross-check ni de `permissions.allow`;
  su lugar es una sección propia en `verify`.
- **Self-update pisa builds locales**: el binario de prueba se reemplaza por la release en el
  próximo `install` hasta que publiques.
- **Bug de PATH "user-local"**: si el prefix npm ya era user-local, el installer saltaba el
  `EnsurePathInShell` → el bin quedaba fuera del PATH. Asegurar el PATH **siempre**, no solo
  al reconfigurar.
