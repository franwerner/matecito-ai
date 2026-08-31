# EDR — El contrato de retorno de sdd-propose se suma al sweep de vocabulario de lane

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
two-lanes-fixed-flow reemplaza el modelo de cuatro lanes (direct/reduced/full/custom, con propose como add-on opcional) por dos lanes fijos donde toda fase corre siempre, y colapsa el modo Interactive/Automatic a una garantía única sin checkpoint entre fases. El Scope ratificado del cambio no incluía `payload/domains/development/references/phase-returns/sdd-propose/sdd-propose.md`, pero el archivo describía dos cosas que el cambio retira: la línea 32 decía "Propose es un add-on opcional", y las líneas 194-195 justificaban por qué `### Scope and approach (unconfirmed)` es la única superficie de revisión contrastando el modo Interactive (donde el checkpoint entre fases muestra el artefacto) con el modo Automatic (donde nada lo muestra). Ambas afirman un mecanismo que, tras el cambio, no existe.

## Decisión
El archivo entra al sweep aunque no estuviera en el Scope original — es el mismo caso que `structure/gating-vocabulary-sweep-scope.md` ya ratificó: un archivo de lectura en runtime que, sin editar, queda apuntando a algo que ya no resuelve. La línea 32 se reescribe para no llamar a propose "opcional". Las líneas 194-195 no se borran: con un solo modo activo, la razón por la que este mailbox es la única superficie de revisión no desaparece — se vuelve más fuerte, porque ahora es la única superficie de revisión sin excepción de modo — así que se reescriben incondicionalmente en vez de eliminarse.

## Alcance
- `payload/domains/development/references/phase-returns/sdd-propose/sdd-propose.md` — Reescrito fuera del Scope original, per gating-vocabulary-sweep-scope.md.

## Reglas verificables
- **[manual]** `sdd-propose.md` no describe a `sdd-propose` como un add-on opcional del pipeline.
- **[manual]** El párrafo que justificaba `### Scope and approach (unconfirmed)` como única superficie de revisión ya no contrasta modo Interactive vs Automatic; afirma la garantía sin condicionarla a un modo.

## Alternativas consideradas
Dejar el archivo fuera del Scope, como residuo aceptado (el mismo tratamiento que `structure/retired-vocabulary-in-record-stores.md` dio a otros tres archivos). Descartada: a diferencia de esos tres, acá el sustantivo que queda stale no es incidental — la oración entera describe un mecanismo (add-on opcional, dos modos) que deja de existir, no una etiqueta dentro de una regla que sigue siendo correcta. Es el caso que `gating-vocabulary-sweep-scope.md` ya distinguió del residuo aceptable.

## Consecuencias
El sweep final del cambio no encuentra vocabulario de lane retirado en `sdd-propose.md`. El archivo pasa a formar parte del changeset de la tarea 3.16 (Slice C), sumándose a los dieciséis ya ratificados en el Scope original.

## Relacionados
- `edr` → [structure/gating-vocabulary-sweep-scope.md](structure/gating-vocabulary-sweep-scope.md) — Fija el criterio de cuándo un archivo fuera del Scope original entra igual al sweep.
- `edr` → [structure/retired-vocabulary-in-record-stores.md](structure/retired-vocabulary-in-record-stores.md) — Fija el criterio contrario — cuándo el sustantivo stale queda como residuo aceptado — que este caso distingue.
