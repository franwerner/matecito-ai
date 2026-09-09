# EDR — Una sección con solo el sentinel `None` cuenta como vacía, nunca como contenido que abre un gate

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Las secciones gateables se emiten SIEMPRE, así que sin esta regla toda sección tiene "contenido" — la
cadena `None…` — y el gate se abriría en cada fase, exactamente la fatiga de confirmación que el guard
existe para evitar.

## Decisión
Una sección gateable cuyo cuerpo es solo un sentinel vacío — cualquier línea que empieza con `None`, con
o sin explicación final — cuenta como EMPTY, no como contenido. El sentinel también se reconoce en
español (`Ninguna`, `Ninguno`, `Nada`), porque el ecosistema conversa en español y los ejecutores derivan
hacia ahí.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Unresolved Decisions Guard (MANDATORY)`, declara que una sección con solo el sentinel `None` (o su equivalente en español) cuenta como EMPTY y no abre gate.

## Alternativas consideradas
Tratar cualquier sección emitida, incluida la del sentinel vacío, como contenido gateable — descartada:
abre un gate en cada fase sobre nada, la fatiga de confirmación que el guard entero existe para evitar.

## Consecuencias
Sin contenido disparador, no hay gate: la siguiente fase se despacha sin mencionar este guard.
