# EDR — Los atajos que rodean 'pregunta abierta = bloqueado' se nombran por su forma, no se restablece el principio

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
El principio de que una pregunta abierta bloquea el avance ya estaba declarado dos veces, y un agente que lo tenía en contexto decidió por su cuenta igual — dos veces: inventando la forma de un contrato JSON público, y haciendo un refactor no solicitado. Su propio relato de la razón: "como sub-agente no puedo preguntar; hago lo mínimo correcto y lo reporto", y "asumí que el lane ya había sido elegido y era `direct`". Ninguno de los dos es un desacuerdo con la regla; los dos son atajos que la rodean sin nunca contradecirla.

## Decisión
En vez de restablecer el principio una tercera vez, se nombran los dos atajos por su forma exacta, como reglas propias: "I cannot ask is never a licence to decide" (no tener canal para preguntar hace la pregunta MÁS bloqueante, no menos) y "A gate you did not watch resolve was not resolved" (una tarea en el prompt no es evidencia de que una bifurcación fue ofrecida o una pregunta contestada).

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### Open question = blocked, not permission`, nombra explícitamente los dos atajos por su forma ("no puedo preguntar" y "asumir que un gate ya se resolvió") en vez de repetir el principio general una tercera vez.

## Alternativas consideradas
Repetir el principio una tercera vez con otras palabras — descartada: el agente ya lo tenía en contexto dos veces y decidió por su cuenta de todos modos; una regla que no atrapa el razonamiento que un agente realmente usa no atrapa nada.

## Consecuencias
La sección crece con los dos atajos nombrados, pero la regla queda más difícil de sortear sin contradecirla explícitamente en vez de rodearla con una justificación que suena razonable.
