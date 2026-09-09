# EDR — El fragmento de dominio se carga por un acto que siempre sucede, no por clasificar dentro del flujo

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
El disparador de carga del fragmento de dominio decía "apenas determinás a qué dominio pertenece un pedido sustantivo — a más tardar cuando intake clasifica". Eso encadenaba la carga a un acto (clasificar) que sólo pasa DENTRO del flujo: un agente que va directo al código nunca clasifica, nunca determina un dominio, y nunca carga el fragmento. Una prueba funcional lo confirmó — el agente tenía el kernel y el CLAUDE.md del proyecto, nunca leyó el fragmento, y lo dijo: "never ran the domain resolution step". El costo no es una regla perdida entre muchas: "Contract & definition shapes — never inferred" vive SÓLO en el fragmento, así que la regla más fuerte contra inventar un contrato era inalcanzable justo en el lane donde ningún guard de fase está mirando tampoco.

## Decisión
El disparador de carga es ahora un acto que siempre sucede, sin importar el lane: antes de crear o modificar el primer archivo del material de un dominio — código, tests, config, assets de diseño. Esto dispara en TODO lane, incluido `direct`, y no depende de haber corrido intake, clasificado nada, o entrado al flujo. El segundo disparador (cuando intake clasifica el pedido) sigue existiendo para el trabajo que sí pasa por el flujo, pero ya no es el único.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### Domain resolution & on-demand loading`, declara dos disparadores: antes de crear/modificar el primer archivo del material de un dominio (todo lane, incluido `direct`), y cuando intake clasifica el pedido (flujo).
- **[manual]** El primer disparador no depende de haber clasificado nada ni de haber entrado al flujo.

## Alternativas consideradas
Dejar el disparador atado a la clasificación y confiar en que el lane `direct` es raro — descartada: es precisamente el lane donde el fragmento hace más falta, porque es el camino más corto al código y ningún guard de fase lo cubre.

## Consecuencias
El fragmento se carga en más momentos que antes, incluido cada arranque de trabajo directo/ad-hoc — el costo es una lectura más frecuente de un archivo que antes se saltaba silenciosamente fuera del flujo.
