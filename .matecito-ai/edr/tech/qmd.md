# EDR — qmd is recorded in the root tech/ catalogue as a repo-level technology

- **Status:** Accepted
- **Date:** 2026-09-01

## Contexto
`.matecito-ai/edr/tech/INDEX.md` declares itself a living register of technologies chosen at repo level and instructs to consult it before installing anything new. Its own scope line says repo-level technologies are the ones that affect all sub-apps; qmd is installed on the *user's machine* by the payload, not compiled into or run by any sub-app of this repo — so the record's own instruction and its own scope statement pull in opposite directions for this case. The catalogue's two existing entries, husky and lint-staged, are both root build tooling — evidence for the narrower "affects all sub-apps" reading, and exactly why this point was raised rather than silently absorbed.

## Decisión
qmd gets its own record in the root `tech/` catalogue, following the catalogue's own instruction to consult and record before installing something new, rather than being left out on the "doesn't affect a sub-app" reading.

## Reglas verificables
- **[manual]** .matecito-ai/edr/tech/qmd.md exists and .matecito-ai/edr/tech/INDEX.md names it.

## Alternativas consideradas
Leaving qmd uncatalogued, on the reading that a tool the payload installs on a user's machine is not "this repository's own technology" and therefore falls outside the catalogue's declared scope. Not chosen — ratified in the opposite direction.

## Consecuencias
The tech/ catalogue's declared scope line ("repo-level technologies are ones that affect all sub-apps") is now in tension with its own contents, since qmd affects none of this repo's sub-apps — nobody has resolved that tension. A future reader consulting tech/INDEX.md before installing something will find qmd there despite the scope note appearing to exclude it.
