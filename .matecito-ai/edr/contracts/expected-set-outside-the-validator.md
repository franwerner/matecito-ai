# EDR — El conjunto esperado de propuestas rechazadas vive del lado del orquestador, no del validador

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
sdd-apply declara un veredicto por cada propuesta rechazada que efectivamente chequeó — pero eso sólo prueba que lo declarado es honesto, nunca que la corrida chequeó TODO lo que el orquestador forwardeó. La completitud (¿faltó algún `record:`?) necesita el conjunto de lo forwardeado, un dato que sólo el orquestador tiene: es quien construyó el prompt de despacho. El spec permite tocar `validate-return.js`/`render-return.js` sólo si el mecanismo no es expresable hoy con el motor existente, y acá la mitad de FORMA (¿el ítem que está tiene los dos tokens legales?) sí es expresable hoy — el motor de `render-return.js` ya exige ambos tokens por ítem y `validate-return.js` ya recorre ítems multi-token — pero la mitad de COMPLETITUD necesita un dato (el conjunto forwardeado) que el retorno nunca carga y que ningún script recibe.

## Decisión
El conjunto esperado queda enteramente del lado del orquestador — comparación literal de cada `· record: <domain>/<slug>` devuelto contra los que el propio orquestador forwardeó como rechazados en el prompt de despacho de esa tarea. Ni `render-return.js` ni `validate-return.js` se tocan: ninguno de los dos gana un flag `--expect-records` ni ningún concepto nuevo. La forma de cada ítem sigue validándose donde siempre se validó (el motor); la completitud del conjunto se valida donde vive el único dato que la hace posible.

## Reglas verificables
- **[manual]** `domains/development.md` declara explícitamente: "The expected set is yours, not the return's" — el chequeo de completitud es responsabilidad exclusiva del orquestador.
- **[auto]** `validate-return-tokens.test.js`, `validate-return.js` y `render-return.js` no cambian de comportamiento por este mecanismo — mismos tests, mismo exit code, sin flag nuevo (verificable corriendo `node --test payload/domains/development/dev-tests/` antes y después de este cambio).

## Alternativas consideradas
Un flag `--expect-records` en `validate-return.js`: descartado — el spec permite tocar el script sólo si el mecanismo no es expresable hoy; la mitad de forma ya lo es, y la mitad de completitud necesita un dato (el prompt de despacho) que el script nunca recibe ni debería recibir — mezclaría una responsabilidad de auditoría de datos con un validador de forma de retorno.

## Consecuencias
Cero superficie nueva en los dos scripts: ambos siguen siendo exactamente lo que `structure/two-scripts-render-and-validate.md` fija (uno renderiza, el otro valida forma). La responsabilidad de completitud queda donde tiene que estar — junto al único participante que conoce el conjunto forwardeado —, sin depender de que el motor la adivine.

## Relacionados
- `relacionado-con` → [rejected-proposal-verdict-token.md](rejected-proposal-verdict-token.md) — el token cuya completitud este chequeo audita.
- `relacionado-con` → [../structure/verdict-classified-by-the-orchestrator.md](../structure/verdict-classified-by-the-orchestrator.md) — la misma lectura del orquestador que también resuelve esta mitad.
