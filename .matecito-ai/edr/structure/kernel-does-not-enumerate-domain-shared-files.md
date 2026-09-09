# EDR — El kernel no enumera por nombre los archivos compartidos de un dominio

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
`### State and Conventions` nombraba a mano `persistence-contract`, y ese archivo fue borrado — nada lo leía, y su contenido era una copia paralela de la sección de persistencia del protocolo de fases. El kernel no tiene por qué enumerar la lista de archivos compartidos de un dominio: es una lista que se desactualiza cada vez que un dominio agrega o saca uno, exactamente como acababa de pasar.

## Decisión
El kernel declara sólo que las convenciones compartidas se despliegan como skills y que cada dominio declara cuáles, sin nombrar ningún archivo específico por nombre. El dominio es dueño de esa lista.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### State and Conventions`, no nombra ningún archivo compartido de dominio por su nombre de archivo.

## Alternativas consideradas
Actualizar la mención a la lista vigente de archivos compartidos — descartada: vuelve a quedar desactualizada la próxima vez que un dominio agregue o saque un archivo, el mismo modo de falla que produjo la mención rota original.

## Consecuencias
El kernel deja de tener una lista propia que mantener sincronizada con lo que cada dominio decide tener como archivo compartido.
