# EDR — CI y quality gates

- **Status:** Pending
- **Type:** decision
- **Date:** 2026-07-23

## Razón de omisión / aplazamiento

**Status:** Pending

El gate de tests en CI se define cuando se arme el CI; hoy no existe pipeline de CI. Trigger esperado: cuando se arme el CI.

Acá también aterriza el hook de build de la UI antes del release (el bundle de la UI debe construirse antes del build de Go para embeberse en el binario), tal como lo anticipa la topología de despliegue.
