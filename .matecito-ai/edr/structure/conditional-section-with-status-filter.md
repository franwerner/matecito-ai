# EDR — Una sección condicional puede filtrar además por el status del retorno

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
La sección que propone un contrato tiene que aparecer en dos momentos distintos según la fase: cuando la fase tiene una propuesta que hacer, y sólo si además esa fase terminó en un status bloqueado. Antes de esta decisión, una sección condicional sólo sabía prenderse o apagarse por un hecho declarado por la fase — no sabía filtrar además por el status del retorno.

## Decisión
Ambos casos montan sobre la misma forma ya existente de sección condicional, extendida con un filtro opcional de status. La sección se emite cuando el hecho que la fase declara es verdadero, y —si el contrato declara el filtro— sólo además cuando el status del retorno está entre los permitidos.

## Reglas verificables
- **[auto]** Una sección condicional con filtro de status sólo se emite cuando su hecho declarado es verdadero y el status del retorno está entre los declarados; sin filtro, se comporta como cualquier sección condicional existente.
- **[manual]** El filtro de status es opcional: una sección condicional sin él sigue funcionando exactamente como antes de esta decisión.

## Alternativas consideradas
Una forma de sección nueva y separada para este caso — descartada porque hubiera duplicado la lógica de sección condicional que ya existe, en vez de extenderla con un filtro adicional sobre algo que ya declara cada retorno: su status.

## Consecuencias
El mecanismo de sección condicional gana una capacidad reutilizable —filtrar además por status— sin que las secciones condicionales existentes, que no la usan, cambien de comportamiento.

## Relacionados
- `relacionado-con` → [../contracts/contract-proposal-has-no-persistence-slot.md](../contracts/contract-proposal-has-no-persistence-slot.md) — la sección que usa este filtro por primera vez.
