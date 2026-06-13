---
name: git
description: Define formato (Conventional Commits), tipos válidos, scope por aggregate, idioma (inglés), reglas duras de atribución (sin Claude, sin Co-Authored-By), atomicidad, y la regla #1 — la IA NUNCA modifica código al commitear, ni siquiera para forzar atomicidad. Solo el usuario, con autorización EXPLÍCITA, autoriza tocar código. Mantenelo simple — esto es la guía completa.
---

# Spec-graph Git — Convenciones de commits

## ⛔ REGLA #0 — La IA NO modifica código al commitear

**Esta regla está por encima de TODAS las demás. Si entra en conflicto con cualquier otra cosa (atomicidad, formato, "mejor práctica"), ESTA gana.**

Al commitear, la IA **NUNCA** debe:

- Editar archivos para "limpiar" antes del commit.
- Reformatear, renombrar variables, reordenar imports, borrar líneas en blanco.
- "Arreglar typos al pasar", aunque sean obvios.
- Mover código entre archivos para que la atomicidad cierre mejor.
- Dividir un cambio existente partiéndolo con `git add -p` y re-editando el resto.
- Generar archivos nuevos (incluido `.gitignore`, configs, docs) sin pedido explícito.
- Hacer `git restore`, `git checkout --`, `git stash`, `git reset --hard`, o cualquier operación destructiva sobre el working tree.

**La IA commitea LO QUE EL USUARIO YA ESCRIBIÓ, tal cual está.** Punto.

### Qué hacer cuando NO se puede atomizar

Caso típico: el working tree tiene cambios mezclados (feat + fix + chore en los mismos archivos) y no se pueden separar limpiamente sin editar código.

**Flujo obligatorio:**

1. **STOP.** No hacer `git add`, no hacer `git commit`, no editar nada.
2. **Reportar la inconsistencia al usuario**: qué cambios ves mezclados, en qué archivos, y por qué no se pueden separar con `git add -p` puro.
3. **Proponer opciones** sin ejecutarlas:
   - "Puedo commitear todo junto como `feat(...)` aunque incluya un fix — ¿lo autorizás?"
   - "Para separarlo necesitaría editar X archivo y mover Y línea — ¿me autorizás a tocar código?"
   - "Podés vos hacer el split y yo commiteo después."
4. **Esperar autorización EXPLÍCITA del usuario.** "Dale", "sí", "hacelo", "autorizado" cuentan. Silencio, ambigüedad, o un "ok" sin contexto NO cuentan.
5. **Solo entonces** ejecutar lo autorizado, y SOLO lo autorizado.

**Autorización explícita ≠ autorización amplia.** Si te autoriza a "mover esa línea", no podés además renombrar la variable. Una autorización = una acción concreta.

### Por qué esta regla existe

- El usuario es el dueño del código. La IA es un ejecutor de commits, no un editor encubierto.
- Cambios "limpios al pasar" rompen `git blame`, ensucian diffs y meten cambios no revisados al historial.
- Si la IA edita por iniciativa propia, el commit miente sobre lo que hizo el humano.
- La atomicidad es un objetivo, no un mandato que justifique tocar código.

## Formato

```
<type>(<scope>): <subject>

[body opcional]
```

Conventional Commits estándar. Sin extensiones custom.

## Types

| Type | Cuándo |
|------|--------|
| `feat` | Feature nueva (código de producto). |
| `fix` | Bug fix (código de producto). |
| `refactor` | Cambio que no altera comportamiento ni agrega feature ni fixea bug. |
| `docs` | Solo documentación (README, SKILL.md, comentarios masivos, anotaciones swag). |
| `test` | Solo agregar o ajustar tests. |
| `chore` | Mantenimiento que no entra en otra categoría (deps, config, gitignore, scripts). |
| `build` | Cambios en el sistema de build o dependencias compiladas. |
| `perf` | Optimización de performance sin cambio funcional. |
| `style` | Formato, whitespace, semicolons. NO cambios de lógica. |
| `revert` | Reversión de un commit anterior. |

Si dudás entre dos types, prevalece el más específico (`fix` > `refactor`, `feat` > `chore`).

## Scope

Nombre del aggregate o área tocada, en lowercase singular:

- `feat(tag): expose public API`
- `fix(project): correct repository tag id type`
- `docs(api): update CLAUDE.md with open-api skill`
- `chore(deps): bump gorm to v1.25.5`
- `refactor(shared): extract transactioner to its own package`

Si el cambio es transversal y no toca un aggregate específico, omitir scope:

- `chore: update gitignore`
- `docs: add architecture overview to README`

## Subject

- **Inglés** siempre.
- **Imperativo**: `add`, `fix`, `remove`, `update` — no `added`, `fixing`, `updates`.
- **Lowercase**: `add tag entity`, no `Add Tag Entity`.
- **Sin punto final**: `expose public API`, no `expose public API.`
- **Max ~70 caracteres**. Si necesitás más, falta en el body.
- **Qué hace, no cómo**: `add tag entity` es bueno; `add Tag struct with ID Name Color fields` es ruido.

## Body

**Solo cuándo el WHY no es obvio del subject**. El body explica motivación, no implementación (el diff ya muestra el qué).

Casos típicos donde body sí aporta:

- Decisión arquitectónica con alternativas evaluadas.
- Bug fix donde el root cause es no obvio.
- Cambio que rompe expectativa razonable de un futuro lector.

Casos donde body NO aporta:

- Repetir el subject con más palabras.
- Listar archivos cambiados (el diff los muestra).
- Narrar paso a paso el cambio.

Formato: línea en blanco después del subject, body en prosa o bullets cortos.

```
fix(project): correct projectID type in SaveProjectTag

The interface was using `int` for projectID while DeleteProjectTag
used `uuid.UUID`. Reconciled to UUID since project IDs are UUIDs
across the rest of the codebase.
```

## Reglas duras

1. **NO `Co-Authored-By`** ni ninguna línea de atribución.
2. **NO mencionar a Claude** ni a ninguna IA en el commit (subject ni body ni footer).
3. **NO `Generated with Claude Code`** ni variantes.
4. **Conventional Commits siempre**. Sin commits estilo "wip", "stuff", "trying things".
5. **NO modificar código al commitear** — ver [Regla #0](#-regla-0--la-ia-no-modifica-código-al-commitear). Solo con autorización EXPLÍCITA del usuario.

## Atomicidad

**Un commit = un cambio lógico.**

Si en una sesión hiciste:

- Agregar Tag aggregate.
- Arreglar bug en project repository.
- Bumpear gorm.

Son **3 commits separados**, no 1:

```
feat(tag): add tag entity and repository
fix(project): correct projectID type in SaveProjectTag
chore(deps): bump gorm to v1.25.5
```

Razones:

- `git revert` puede deshacer un cambio sin tocar los otros.
- `git log` cuenta una historia coherente.
- Code review puede aprobar/rechazar por separado.
- `git bisect` apunta al commit exacto que rompió algo.

> **Atomizar NO autoriza a la IA a tocar código.** Si los cambios no se separan con `git add -p` puro, aplica el flujo de [Regla #0](#-regla-0--la-ia-no-modifica-código-al-commitear): STOP, reportar, esperar autorización explícita.

## Reglas operativas

- **No commit con código que no compila** (excepto WIP explícito en branch propio que no se va a mergear).
- **No commit mezclando cambios no relacionados**. Si te das cuenta tarde, `git reset` y separá en commits atómicos.
- **No `git commit --amend` después de pushear**. Reescribe historial compartido.
- **Si te equivocaste antes de pushear**, `git commit --amend` o `git reset` está OK.

## Anti-patterns

| ❌ Mal | ✅ Bien |
|--------|--------|
| `update stuff` | `feat(tag): add public API for cross-context reads` |
| `fix bug` | `fix(project): handle nil tag IDs in SaveProjectTag` |
| `feat: add a lot of things` | (separar en commits atómicos) |
| `Added new tag entity.` (pasado, mayúscula, punto final) | `add tag entity` |
| Commit con `Co-Authored-By: Claude` | Commit sin atribución de IA |
| `WIP` mergeado a main | (squash o reescribir antes de mergear) |

## Workflow rápido (para la IA)

1. **Leer el estado actual** (`git status`, `git diff`) — sin editar nada.
2. **Si los cambios son atomizables con `git add` / `git add -p` puros** → agrupar por unidad lógica y commitear cada grupo por separado.
3. **Si NO son atomizables sin editar código** → STOP. Reportar al usuario, proponer opciones, esperar autorización EXPLÍCITA. Ver [Regla #0](#-regla-0--la-ia-no-modifica-código-al-commitear).
4. Escribir mensaje siguiendo el formato.
5. `git commit` (sin `--no-verify` salvo pedido explícito del usuario).
6. `git push` solo si el usuario lo pidió.
