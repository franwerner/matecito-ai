# EDR — Resiliencia

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-27

## Contexto

El broker es un daemon local único del que dependen tres canales con modos de falla distintos: los hooks sincrónicos del runtime del coder (que meten latencia en el camino del trabajo y no pueden colgar la sesión), los productores de ingesta fire-and-forget (cuya conducta ante broker caído —spool local completo y reconciliación— ya está fijada por el comportamiento especificado del sistema), y el WebSocket de lectura de la UI (cuya reanudación sin huecos por posición del log también está fijada). Lo que queda por decidir acá son las políticas concretas: presupuestos de tiempo, reintentos y límites.

## Decisión

- **Presupuesto de los hooks sincrónicos: 500 ms en total.** El daemon responde en localhost en milisegundos; agotar el presupuesto significa broker caído o colgado. Al vencer: el hook de captura pre-edición cae al spool local (el cliente lee el archivo él mismo) y libera la edición; la señal de notas del hook de prompt se saltea en silencio, sin spool — es solo una señal, reaparece en el próximo mensaje.
- **Reconexión del WebSocket de la UI: backoff exponencial con jitter — 1 s inicial, tope 30 s, reintentos infinitos.** La desconexión típica es un reinicio del daemon de segundos: el primer reintento rápido casi siempre alcanza; el tope evita martillar en caídas largas; nunca se rinde (la UI indica el estado de reconexión y, al volver, reanuda por posición del log).
- **Spool local: sin proceso residente del lado cliente y sin límites.** El flush ocurre únicamente en el próximo contacto exitoso (lo pendiente antes que lo nuevo); no hay timers ni daemons del lado del hook. Sin tope de tamaño ni expiración: la garantía es que nada se pierde, y un límite reintroduciría la pérdida que el spool existe para evitar.
- **Techo de recepción de cabeceras en el borde HTTP: 5 s.** Existe para reclamar conexiones de clientes que mueren a mitad de request —el caso realista en un daemon local—, no como defensa ante un atacante: el borde escucha en localhost. El valor es holgado a propósito; en localhost un cliente sano no se acerca. Los presupuestos de escritura e inactividad quedan **sin fijar hasta que exista el canal de lectura por WebSocket**: son los que pueden cortar conexiones largas legítimas, y elegirlos sin ese canal a la vista sería decidir a ciegas.

## Alcance

- `apps/api/internal/**` — hub de suscripciones y superficie de ingesta.
- El cliente mínimo de hooks que despliega el installer.

## Reglas verificables

- **[manual]** los hooks sincrónicos operan con presupuesto total de 500 ms; al vencerlo, la captura pre-edición spoolea localmente y la señal de notas se saltea en silencio — en ningún caso la sesión del coder queda colgada o con error visible.
- **[manual]** la reconexión WS usa backoff exponencial con jitter (1 s inicial, tope 30 s) y no tiene límite de reintentos.
- **[manual]** el cliente de hooks no tiene procesos residentes ni timers: el flush del spool solo ocurre en el próximo contacto exitoso, entregando lo pendiente antes que lo nuevo.
- **[manual]** el spool local no tiene tope de tamaño ni expiración; ningún ítem se descarta por antigüedad o volumen.
- **[tool: gosec]** el servidor del borde HTTP declara un techo de recepción de cabeceras; construirlo sin ese techo rompe el gate de lint.

## Alternativas consideradas

- **Presupuesto de 200 ms:** cae a spool más rápido ante un broker lento, a costa de spoolear de más en hiccups transitorios. **De 1 s:** casi nunca spoolea de más, pero una edición puede sentirse trabada.
- **Reintentos WS finitos con reconexión manual:** más control, más fricción; innecesario contra un daemon local.
- **Intervalo fijo de reconexión:** o martilla o tarda de más; el backoff con tope cubre ambos extremos.
- **Flusher con timer del lado cliente:** vaciaría el spool apenas vuelve el broker sin esperar actividad, pero agrega un proceso residente que hoy no existe; descartado por simplicidad.
- **Tope de spool con descarte de lo más viejo:** acota disco pero rompe la garantía de no-pérdida; descartado.

## Consecuencias

- La sesión del coder tiene un techo de latencia conocido e imperceptible (500 ms solo en falla); el costo de un broker caído se paga en reconciliación diferida, nunca en fricción del trabajo.
- Un spool sin límites puede crecer durante una caída larga con mucha actividad; es texto local y se vacía entero al volver el contacto — se acepta el riesgo por la garantía de no-pérdida.
- Sin flusher residente, un spool con pendientes puede esperar hasta el próximo hook para vaciarse; sin actividad no hay nada urgente que mostrar, así que la espera es inocua.
- Los valores concretos (500 ms, 1 s/30 s) son puntos de partida razonados; si la implementación demuestra otra cosa, se revisan editando este EDR (cambio menor, git lleva el historial).

## Relacionados

- `depende-de` → [../delivery/deployment-topology.md](../delivery/deployment-topology.md) — el daemon local único es la premisa de todos los presupuestos (latencia localhost, caídas = reinicios).
- `relacionado-con` → [../contracts/api-contract.md](../contracts/api-contract.md) — el WebSocket de lectura cuya reconexión gobierna este EDR.
- `relacionado-con` → [concurrency-async.md](concurrency-async.md) — el hub de suscripciones donde vive la reanudación por posición.
