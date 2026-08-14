# EDR — La discusión lateral es bloqueante o consultiva, declarada al abrir; el pickup es siempre por consulta

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Mientras la sesión lateral discute, el principal puede seguir trabajando o no, según si lo que se está discutiendo es una premisa de la que depende el trabajo siguiente. No existe en este entorno ningún mecanismo de push, webhook o señal entre sesiones — sólo consulta.

## Decisión
Cada discusión es bloqueante o consultiva, y el usuario lo dice al abrirla; si no lo dice, el orquestador pregunta, nunca lo elige por su cuenta. Bloqueante: el principal se detiene y no avanza sobre nada que dependa de la discusión — el trabajo no relacionado puede seguir. Consultiva: el principal sigue trabajando y recoge la conclusión cuando el usuario avisa o cuando llega a un punto donde la necesita. En ambos casos la conclusión se lee sólo por consulta, nunca por notificación; si la clave está ausente, el orquestador lo dice y ofrece esperar, reabrir la discusión con el mismo traspaso, o traer la pregunta al chat principal — sin adivinar. No hay timeout tras el cual el principal siga con su propia lectura.

## Reglas verificables
- **[manual]** El usuario declara el tipo (blocking | consultive) al abrir la discusión; si no lo dice, el orquestador pregunta — nunca lo infiere ni lo elige por su cuenta.
- **[manual]** Blocking: el principal no avanza sobre nada que dependa de la discusión hasta tener la conclusión; el trabajo no relacionado puede seguir.
- **[manual]** Consultive: el principal sigue trabajando y recoge la conclusión cuando el usuario avisa que terminó o cuando el principal llega a un punto que la necesita.
- **[manual]** La conclusión se lee sólo por consulta, nunca por notificación; si la clave todavía no existe, el orquestador lo dice y ofrece esperar, reabrir la discusión con el mismo traspaso, o traer la pregunta al chat principal — sin elegir un default.
- **[manual]** No hay timeout tras el cual el principal siga adelante con su propia lectura de la pregunta abierta.

## Alternativas consideradas
Un solo modo, sin espera, para toda discusión — descartado: obliga al principal a seguir construyendo sobre una premisa que la discusión puede tumbar, que es justo la razón de ser de una discusión bloqueante. Inferir el tipo de la forma de la pregunta — descartado por el mismo motivo que el orquestador no detecta discusiones: es al usuario a quien le toca decidirlo, y una inferencia equivocada cuesta o trabajo parado sin necesidad o trabajo construido sobre una premisa muerta.

## Consecuencias
Nunca se pierde trabajo por avanzar sobre una premisa bloqueante todavía sin resolver. El costo: una discusión bloqueante sin cerrar puede frenar trabajo dependiente indefinidamente si el usuario no vuelve — no hay timeout que lo destrabe solo.

## Relacionados
- `relacionado-con` → [side-discussion-launch-boundary.md](side-discussion-launch-boundary.md) — la otra decisión de este cambio ajustada en el mismo gate.
