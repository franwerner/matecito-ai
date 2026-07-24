# golang.org/x/sync

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** runtime
- **Date:** 2026-07-23

## Por qué la elegimos

Su grupo de goroutines con cancelación propagada (errgroup) orquesta el lifecycle del broker —arranque, espera y cancelación coordinada del escritor, el hub y los tailers— sin escribir a mano el boilerplate de esperas y cancelación.

## Alternativas descartadas

- **Coordinación manual de esperas y cancelación (stdlib sync):** más verbosa y propensa a fugas de goroutines para el mismo resultado.

## Notas

Usada en: runtime/concurrency-async.
