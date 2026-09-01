# EDR — El aislamiento por espacio de trabajo tiene dos niveles anidados

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto

Varias sesiones concurrentes sobre el mismo repositorio, cada una con un cambio distinto, comparten un
único directorio de trabajo: se pisan los archivos entre ellas y el estado sucio de una contamina la
base de las otras. El aislamiento que ya existe cubre un solo nivel —las corridas de un batch de
implementación paralelo, dentro de una fase— y está atado al ciclo de vida de esa fase: se abre al
despachar el batch y se libera al integrarlo. Un cambio, en cambio, abarca el pipeline entero, desde
que se confirma su alcance hasta que se archiva, así que nada de lo existente lo cubre.

## Decisión

El aislamiento se anida en dos niveles. El orquestador abre un **espacio de trabajo aislado por
cambio** y coordina el ciclo adentro; el batch de implementación paralelo abre los suyos **sobre ese
espacio**, no sobre la rama original; al cerrar el ciclo, el orquestador integra el espacio del cambio
a la rama original.

Se activa solo cuando el pedido del usuario lo pide **explícitamente**, nunca por default; ningún gate
lo confirma como ítem propio. El momento de apertura depende de si el trabajo pasa por el flujo: cuando
sí, el espacio se abre una vez que se confirma en el Brief Confirmation Gate un brief con el aislamiento
activo, antes de despachar la fase siguiente; cuando no —trabajo directo o edición ad-hoc, que nunca
produce un brief—, el momento sigue siendo el mismo de siempre: antes del primer archivo que se escribe.

## Reglas verificables

- **[manual]** Con el aislamiento por cambio activo, las corridas aisladas de un batch paralelo parten del espacio de trabajo del cambio, nunca de la rama original.
- **[manual]** La apertura del espacio de trabajo del cambio ocurre una vez que se confirma en el Brief Confirmation Gate un brief con el aislamiento activo, antes de despachar la fase siguiente, cuando el trabajo pasa por el flujo; y antes del primer archivo escrito, cuando no.
- **[manual]** El aislamiento por cambio se activa solo cuando el pedido lo pide explícitamente; ningún camino lo activa por default ni lo asume activo sin haberlo visto elegir.

## Alternativas consideradas

- **Un solo nivel, reusando el aislamiento del batch para el cambio entero.** Descartado: ese aislamiento está atado al ciclo de vida de una fase y se libera al integrarla, mientras que un cambio abarca el pipeline completo — estirarlo obligaría a sostener un espacio de trabajo a través de fases que hoy no lo conocen.
- **Aislar todo cambio sustancial por default.** Descartado: agrega manejo de ramas y cambio de directorio a quien trabaja en una sola sesión, que es el caso mayoritario, para resolver un problema que solo aparece con sesiones concurrentes.
- **Cerrar el ciclo dejando la rama lista y que la integración la haga una persona a mano.** Descartado: el cierre del ciclo es justamente el momento en que el trabajo quedó verificado; postergar la integración a un paso manual deja el aislamiento sin cierre y el trabajo fuera de la rama por tiempo indefinido.

## Consecuencias

- Sesiones concurrentes dejan de compartir directorio de trabajo, que es el problema que motiva la decisión: el estado sucio de una ya no es la base de otra.
- La base de una corrida aislada deja de ser un valor fijo y pasa a depender de si el nivel de cambio está activo — lo que obliga a resolverla explícitamente en cada despacho en vez de asumirla.
- El orquestador ejecuta una mutación irreversible del repositorio al cerrar el ciclo, lo que acota una regla previa que se leía como prohibición general de que el orquestador integre.
- Cuando la opción se activa, cada cambio cuesta un espacio de trabajo propio y su rama; con la opción apagada el comportamiento es el de antes de esta decisión.

## Relacionados

- `relacionado-con` → [consolidation-run-is-the-integrator.md](consolidation-run-is-the-integrator.md) — qué ejecutor le corresponde a cada uno de los dos niveles.
- `relacionado-con` → [../contracts/worktree-base-handshake.md](../contracts/worktree-base-handshake.md) — cómo se verifica la base cuando el nivel de cambio la desplaza.
- `relacionado-con` → [dispatch-batch-bound-integration.md](dispatch-batch-bound-integration.md) — el batch cuyo destino de integración pasa a ser el espacio del cambio.
