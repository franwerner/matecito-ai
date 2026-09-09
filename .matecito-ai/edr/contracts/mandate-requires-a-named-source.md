# EDR — Un mandato necesita una fuente nombrada; un pedido crudo no alcanza por su forma gramatical

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La regla del hard-stop por desviación definía el mandato como "los artefactos acordados + el alcance confirmado", pero no daba ningún test para decidir si algo se había acordado. Una prueba funcional expuso la consecuencia: un pedido imperativo de una sola línea ("Agregá un comando X que devuelva JSON") se leyó como un mandato ya confirmado, y el agente "ejecutó libremente dentro de él" — inventó un contrato JSON público, creó dos paquetes y refactorizó un tercer archivo, sin mostrar nunca la bifurcación de Lane. Su propio relato: "traté el mensaje como una asignación de implementación ya acotada por vos". La bifurcación de Lane nunca dispara porque llega a una pregunta que el agente cree ya contestada.

## Decisión
Un mandato existe sólo cuando tiene una de dos fuentes: un artefacto de flujo confirmado (intake brief, spec, design, tasks), o una confirmación explícita del usuario en la conversación actual. Un pedido crudo — por más detallado o imperativo que sea — NO es un mandato: es el insumo que tiene que llegar a una de esas dos fuentes, no una fuente en sí mismo. La forma gramatical no otorga autoridad: un imperativo ("Agregá X", "Necesito que hagas X") recibe el mismo tratamiento que una pregunta ("¿podés agregar X?").

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### Deviation hard-stop (anchored to the mandate)`, declara el test explícito: un mandato existe sólo si tiene una fuente de las dos nombradas, y un pedido crudo no es una fuente en sí mismo.
- **[manual]** La misma sección declara que la forma gramatical del pedido (imperativo, pregunta, urgencia percibida) no cambia el tratamiento.

## Alternativas consideradas
Restablecer la regla original ("ejecutá dentro del mandato acordado") sin agregar el test de fuente — descartada: es exactamente lo que ya estaba escrito y lo que la prueba funcional mostró que fallaba en la práctica.

## Consecuencias
Un pedido imperativo sin una de las dos fuentes confirmadas frena en la bifurcación de Lane en vez de ejecutarse como si ya estuviera acotado. El costo es un paso más de verificación antes de cualquier trabajo que no venga de un artefacto de flujo ya confirmado.
