# EDR — El lanzamiento de la sesión lateral se prueba por cinco propiedades, sin nombrar herramienta

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
La sesión lateral necesita una terminal nueva e interactiva, abierta automáticamente por el orquestador. Nombrar una herramienta concreta —aunque sea "la que anda hoy"— acopla el mecanismo a un launcher que puede no estar instalado o puede cambiar, exactamente el problema que el índice de exploración de este ecosistema ya resuelve por capacidad en vez de por nombre fijo.

## Decisión
El requisito del lanzamiento se declara como cinco propiedades verificables, y ningún archivo del payload nombra una herramienta concreta —ni como ejemplo, ni como default, ni como "la que funciona hoy": (1) el orquestador la invoca sin que el usuario haga nada; (2) lo que abre es una sesión interactiva que el usuario puede leer y escribir; (3) se sabe si arrancó o no; (4) abre parada en el árbol donde ya trabaja la sesión que la lanza —el change workspace si el aislamiento a nivel de cambio está activo para el cambio en curso, el árbol de trabajo común si no; (5) no crea worktree ni checkout propio. Cómo se cumplen esas cinco propiedades se resuelve al momento de uso, contra lo que el entorno ofrezca, y no queda fijado en ningún archivo. Si nada resuelve en un lanzamiento que cumpla las cinco, el mecanismo no está disponible: se dice en una línea, se ofrece traer la pregunta al hilo principal nombrando el costo (gasta el contexto que la discusión existía para ahorrar), no se entrega ningún comando, y se menciona una vez que un lanzador restauraría la función.

## Reglas verificables
- **[manual]** Ningún archivo del payload nombra una herramienta concreta de lanzamiento — ni como ejemplo, ni como default, ni como la opción vigente.
- **[manual]** El lanzamiento cumple las cinco propiedades: automático, interactivo, arranque verificable, hereda el directorio de la sesión que lo lanza, y no crea worktree ni checkout propio.
- **[manual]** Cuando nada resuelve un lanzamiento que cumpla las cinco, el orquestador lo dice en una línea, ofrece traer la pregunta al hilo principal nombrando el costo, y no entrega ningún comando.

## Alternativas consideradas
Nombrar el launcher que funciona hoy, como default o como instancia ilustrativa del test — descartado: el ejemplo se convierte en la implementación y se desactualiza igual que un nombre hardcodeado. Sondear un conjunto de nombres candidatos — descartado por ser el mismo acoplamiento, escondido. Caer de vuelta a entregarle al usuario el comando para que lo corra — descartado: es el defecto que esta corrección elimina. Abrir la sesión lateral en un checkout propio — descartado: rompe la propiedad (5) y le quita a la discusión la visibilidad del trabajo sin commitear del hilo principal.

## Consecuencias
El mecanismo puede reportarse no disponible en un entorno que en verdad podría haberlo servido, porque la resolución no enumera nada — ese es el costo aceptado de no nombrar herramientas. A cambio, el mecanismo no se rompe ni queda desactualizado cuando cambia qué CLI está instalado o disponible.

## Relacionados
- `relacionado-con` → [side-discussion-launch-boundary.md](side-discussion-launch-boundary.md) — la otra decisión de este cambio sobre el mismo lanzamiento, ajustada en el mismo gate con el costo mayor.
