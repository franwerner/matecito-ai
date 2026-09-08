# EDR — El store de EDRs y de capability-specs es de raíz único tras el retiro del cockpit

- **Status:** Accepted
- **Date:** 2026-09-08

## Contexto
El repo pasa a describir un solo componente (`cli`). Antes de este cambio, las decisiones técnicas vivían por sub-app — un store de EDRs propio para cada uno de los dos sub-apps del cockpit (el broker/MCP y su UI), además del store raíz — y `CLAUDE.md` documentaba esa organización como la norma. Este cambio (`sdd/drop-apps-subtree`) borra los dos sub-apps y sus dos stores propios junto con el resto del cockpit, y `CLAUDE.md` necesita decir dónde viven las decisiones y el comportamiento del sistema en el repo que queda.

## Decisión
`.matecito-ai/edr/` pasa a ser el único store de decisiones, organizado por dominio, y gobierna todo el repo incluido `cli` — `map-records.js` ya lo resuelve así porque ningún componente está declarado `root: true` (verificado contra su salida: `componentsWithoutOwnStore = [cli]`). El store de capability-specs (`.matecito-ai/development-specs/`) se queda de raíz por la misma razón que ya tenía antes de este cambio — `rule/single-rooted-spec-store.md`: el comportamiento es intencionalmente multi-superficie y separarlo por componente escondería esa relación — y `CLAUDE.md` re-funda su afirmación de "por qué vive en el root" en esa regla en vez de en "el producto tiene dos mitades", que es la afirmación que este cambio falsifica.

## Reglas verificables
- **[manual]** La sección "Decisiones de ingeniería (EDRs)" de `CLAUDE.md` afirma que las decisiones viven en un solo store, `.matecito-ai/edr/`, organizado por dominio, que gobierna todo el repo incluido `cli`, y no nombra ninguna ruta bajo `apps/`.
- **[manual]** La sección "Comportamiento del sistema (capability-specs)" de `CLAUDE.md` funda su afirmación "vive en el root" en `rule/single-rooted-spec-store.md`, no en que el producto tenga dos mitades.

## Alternativas consideradas
Generalizar el framing vigente ("por sub-app, junto al código que gobiernan") a una forma abstracta en vez de afirmar un store único. Descartada: preserva una distinción que el repo ya no tiene (queda un solo componente) y reinvita un store por componente el día que se agregue uno nuevo. También se consideró declarar `cli` como `root: true` en `config.json` para hacer explícita la propiedad — descartada: es una edición fuera del alcance ratificado (el trim de `repo.components`) y no aporta nada con un solo componente.

## Consecuencias
Quien busque una decisión o un comportamiento lee `.matecito-ai/edr/INDEX.md` o `.matecito-ai/development-specs/INDEX.md` directamente, sin indirección por sub-app. Si en el futuro se agrega un segundo componente, esta afirmación (map-records.js resolviendo por ausencia de `root: true`) hay que revisarla junto con la declaración explícita en `config.json`.
