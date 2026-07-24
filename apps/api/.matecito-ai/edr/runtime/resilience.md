# EDR — Resiliencia

- **Status:** Pending
- **Type:** decision
- **Date:** 2026-07-23

## Razón de omisión / aplazamiento

**Status:** Pending

La resiliencia del broker (reconexión del WebSocket de la UI y reanudación del tail del transcript por offset) se define cuando aterrice la lógica del broker, porque depende de cómo queden implementados el hub y los tailers. Trigger esperado: cuando aterrice la lógica del broker.
