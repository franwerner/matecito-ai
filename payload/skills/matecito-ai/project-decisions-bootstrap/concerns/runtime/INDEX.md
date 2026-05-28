# Dominio: `runtime`

Cómo se comporta el sistema mientras corre: errores, concurrencia, trabajo en background, resiliencia ante fallos y caching. Define la *dinámica* del sistema.

## Criterio de pertenencia

Un concern nuevo va en `runtime` si trata sobre lo que pasa durante la ejecución (timing, fallos, paralelismo, estado efímero). Si trata sobre cómo se persisten los datos, va en `data`.

## Concerns en este dominio

| Concern | Prof. | Tipo | Qué decide |
|---|---|---|---|
| [background-jobs](background-jobs.md) | light | decisión | Si el proyecto necesita procesar trabajo fuera del ciclo request/response, y con qué mecanismo: cola, scheduler, o ninguno. |
| [caching](caching.md) | light | decisión | Qué se cachea, dónde, y cómo se invalida. Mal hecho sirve datos viejos; bien hecho define latencia y costo. |
| [concurrency-async](concurrency-async.md) | light | decisión | Cómo el proyecto maneja operaciones que no son estrictamente secuenciales: async nativo del lenguaje, threads, workers de proceso, o simplemente síncrono dir... |
| [error-handling](error-handling.md) | deep | decisión | Cómo se representan, propagan y responden los errores en todo el sistema. Es de las decisiones más transversales: toca dominio, infraestructura y bordes. |
| [resilience](resilience.md) | light | decisión | Cómo el sistema se comporta cuando una dependencia externa (DB, API tercera, cola) es lenta o falla. Sin política explícita, el default es "esperar indefinid... |
