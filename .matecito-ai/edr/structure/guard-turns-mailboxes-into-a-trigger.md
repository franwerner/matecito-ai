# EDR — El Unresolved Decisions Guard existe para convertir los buzones de cada fase en un disparador, no un depósito

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Antes de este guard, los buzones que cada fase retorna (Open Questions, New Decisions, Deviations,
Derived capabilities, risks) existían como campos declarados sin que nada los consumiera: eran el lugar
barato donde una fase depositaba lo no resuelto y el flujo seguía igual. El guard nació específicamente
para cerrar eso — convertir esos buzones en un disparador real que frena el flujo cuando corresponde.

## Decisión
El Unresolved Decisions Guard inspecciona el envelope de retorno de cada fase después de que retorna y
antes de despachar la siguiente, y dispara un gate por item (no por sección) según dos triggers
cerrados. Los buzones dejan de ser un depósito pasivo: un item que dispara frena el flujo y espera al
usuario.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Unresolved Decisions Guard (MANDATORY)`, abre declarando que inspecciona el envelope de retorno después de cada fase y antes de despachar la siguiente, y que un gate dispara por item.

## Alternativas consideradas
Dejar los buzones como campos meramente informativos sin ningún mecanismo que los lea — descartada: es
exactamente el estado que motivó crear este guard, donde lo no resuelto se acumulaba sin que nadie lo
viera.

## Consecuencias
Toda fase que declara un item en uno de estos buzones sabe que puede disparar un gate, no que está
escribiendo a un registro que nadie lee.
