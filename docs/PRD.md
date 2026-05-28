# matecito-ai — PRD

> *Mientras la IA trabaja por vos, te tomás unos ricos mates.*

Documento de producto (liviano) del ecosistema de desarrollo con IA **matecito-ai**.

---

## Qué es

matecito-ai es un **ecosistema de desarrollo asistido por IA**, armado a medida sobre Claude Code. No es una herramienta nueva: es la integración curada de varias piezas —propias y de terceros— en un flujo coherente, donde cada decisión de arquitectura queda registrada y respetada a lo largo del tiempo y entre sesiones.

La idea de fondo: que el agente **no reinvente las convenciones del proyecto en cada sesión**. Las decisiones se capturan una vez (como ADRs), se respetan al implementar, y la memoria de trabajo persiste vía Engram. El humano decide; la IA ejecuta dentro de esas decisiones.

---

## Problema que resuelve

Trabajar con agentes de IA sobre un proyecto tiene tres fugas recurrentes:

1. **Amnesia entre sesiones.** Cada sesión nueva arranca sin saber qué se decidió antes — re-sugiere librerías ya descartadas, reinventa convenciones, contradice decisiones previas.
2. **Decisiones implícitas.** La arquitectura vive en la cabeza del autor, no escrita. El agente no tiene cómo respetarla.
3. **Exploración cara.** El agente gasta tokens y tool calls escaneando el codebase archivo por archivo para entender estructura que podría consultarse de forma indexada.

matecito-ai ataca las tres: ADRs para las decisiones, Engram para la memoria de sesión, codegraph para la exploración eficiente.

---

## Usuario

Un solo desarrollador (el dueño del ecosistema) que usa IA intensivamente, trabaja sobre repos propios recurrentes, y prioriza un entorno 100% adaptado a su forma de trabajar por encima de la facilidad de mantenimiento. Acepta dar mantenimiento al fork a cambio de control total.

---

## Componentes

```
matecito-ai/
├── SKILLS    project-decisions (ADRs) + issue-brief + SDD (Gentleman, fork)
├── MCP       codegraph + context7
├── AGENTES   sub-agentes del SDD (fork modificado)
└── ENGRAM    memoria de sesión (standalone)
```

### SKILLS

- **project-decisions-bootstrap** — entrevista por fases que captura decisiones de ingeniería (arquitectura, convenciones, políticas) y las materializa como ADRs por dominio en `.claude/adr/<dominio>/`. Catálogo de 39 concerns en 17 dominios canónicos. Incluye `CONCERN-TEMPLATE.md` como guía de autoría para expandir el catálogo (ratchet).
- **project-decisions-validate** — validador consultivo: chequea coherencia entre ADRs, completitud, verificabilidad e integridad de taxonomía. No modifica nada.
- **issue-brief** — puente entre un issue descrito y los ADRs que le aplican; arma un briefing de restricciones, frena ante choques y deriva a bootstrap ante decisiones nuevas. *(Estado: definido; integración con el flujo SDD pendiente de decidir — ver Pendientes.)*
- **SDD (fork del Gentleman)** — workflow de fases (explore → propose → spec → design → tasks → apply → verify → archive) con sub-agentes reales. Modificado para integrarse con ADRs, codegraph y context7.

### MCP

- **codegraph** — grafo de código pre-indexado (tree-sitter + SQLite). Exploración por estructura/relaciones en vez de grep. Eficiente en tokens y tool calls.
- **context7** — documentación de librerías al día, para no implementar contra APIs desactualizadas o alucinadas.

### AGENTES

Los sub-agentes del SDD (uno por fase), con contexto propio. Forkeados y modificados (ver Decisiones).

### ENGRAM

Memoria de sesión persistente en SQLite, instalada standalone (no vía gentle-ai, para no acoplar a su ciclo de releases). Guarda descubrimientos, contexto y fixes entre sesiones.

---

## Decisiones de diseño

| # | Decisión | Razón |
|---|----------|-------|
| 1 | **ADRs = decisiones de arquitectura; Engram = memoria de sesión** | Evitar solapamiento. ADR es decisión deliberada y verificable; Engram es lo que el agente aprendió trabajando. |
| 2 | **Fork directo del SDD, no inyección** | Máxima personalización. Se acepta el costo de mantenimiento a cambio de control total y coherencia con la forma propia de trabajar. |
| 3 | **Engram-only en el SDD (sin openspec/hybrid)** | El proyecto no usa persistencia basada en archivos del SDD; Engram cubre la memoria. Se purgaron los modos openspec/hybrid. |
| 4 | **Sin mecanismo de inyección (registry/resolver removidos)** | Coherente con el fork directo: las convenciones se leen de los archivos del proyecto (ADRs, CLAUDE.md), no de un registry intermedio. |
| 5 | **codegraph-first en exploración, grep como fallback** | codegraph para estructura/relaciones (eficiente); grep para texto literal, archivos no indexados, o cuando codegraph no resuelve. |
| 6 | **design y apply leen los ADRs; verify los chequea** | Las decisiones se respetan en el momento de diseñar e implementar; verify confirma que el código del cambio no viole los ADRs que tocó. |
| 7 | **No construir un issue-implementer propio** | El SDD ya es el implementador disciplinado. Construir otro sería duplicar infraestructura madura. |
| 8 | **Integrar infraestructura cara, construir solo la diferenciación** | ADRs e issue-brief son propios (diferenciación); SDD, Engram, codegraph y context7 se integran (resueltos). |

---

## Fuera de alcance

- No es un instalador de agentes ni un producto distribuible — es un entorno personal.
- No reconstruye memoria persistente ni el orquestador de fases desde cero (se integran Engram y el SDD).
- No incluye los componentes de gentle-ai que no se usan (persona, GGA, skills de frameworks específicos).
- El enforcement de que el código respete los ADRs se delega a herramientas determinísticas (linters vía `arch-enforcement`) + verify, no a un agente vigilante.

---

## Cómo se usa (flujo típico)

1. **Setup del proyecto** (una vez): `project-decisions-bootstrap` captura las decisiones → genera `.claude/adr/` + `CLAUDE.md`.
2. **Al implementar un issue:** se describe el issue → el flujo SDD lo lleva por sus fases. design y apply leen los ADRs vigentes; explore usa codegraph; apply usa context7. Si el issue choca con un ADR o destapa una decisión nueva, se frena y se captura vía bootstrap antes de codear.
3. **Al cerrar:** verify chequea que el cambio respete los ADRs que tocó. Engram guarda lo aprendido.
4. **Mantenimiento del catálogo:** concerns nuevos se agregan vía `CONCERN-TEMPLATE.md`; coherencia entre ADRs se revisa con `project-decisions-validate`.

---

## Pendientes / decisiones abiertas

- [ ] **Verificar nombres reales de las tools MCP.** El fork del SDD asume `mcp__codegraph__*` y `mcp__context7__*`. Confirmar contra la instalación real (`claude mcp list`) y ajustar el frontmatter de los agentes si difieren.
- [ ] **Probar el fork end-to-end.** Correr el flujo SDD completo sobre un proyecto chico para validar que cada fase persiste en Engram y que las modificaciones (ADRs, codegraph, context7) funcionan en la práctica. Es la única prueba real tras la cirugía de purga.
- [ ] **Decidir el rol final de `issue-brief`.** Como design ya lee los ADRs, parte de lo que haría issue-brief quedó absorbido. Definir si entra como pieza separada o si el flujo SDD ya cubre el puente issue→ADRs.
- [ ] **Confirmar política codegraph en explore.** Validar con una corrida real que `sdd-explore` efectivamente prefiere codegraph sobre grep cuando `.codegraph/` existe.
- [ ] **Mantenimiento del fork.** Establecer el hábito de: ante un update del Gentleman, hacer diff contra `vendor-original/` y reaplicar los bloques `matecito-ai`.

---

## Riesgos conocidos

- **Frescura de los ADRs.** Si los ADRs quedan desactualizados respecto al código, el agente trabaja sobre decisiones viejas con confianza. Mitigación: actualizar la decisión (vía bootstrap update) cuando cambia, no solo el código.
- **Frescura del índice de codegraph.** Un grafo desactualizado en el que el agente confía ciegamente es peor que no tenerlo. Mitigación: el auto-sync de codegraph; ante dudas, `codegraph status` antes de sesiones importantes.
- **Divergencia del upstream.** El fork (registry + openspec purgados, fases modificadas) se alejó bastante del SDD original. Portar mejoras del Gentleman será trabajo manual. Mitigación: `vendor-original/` + marcadores `matecito-ai`.
- **Madurez de las dependencias.** gentle-ai/SDD está en v0.1.x ("APIs will change"). Los updates pueden requerir re-trabajo del fork.
