# gofumpt

- **Category:** Other
- **Version:** v0.11.0
- **Status:** Accepted
- **Decided in phase:** ci-quality-gates
- **Date:** 2026-07-28

## Por qué la elegimos

Superset estricto del formateador estándar: aplica todo lo que este aplica y además normaliza formas que el estándar deja libradas al autor. No tiene opciones de configuración, así que el formato deja de ser discutible y de generar diffs de estilo entre ediciones.

## Alternativas descartadas

- **Quedarse con el formateador estándar:** es el piso, no discrimina entre dos formas igual de válidas de escribir lo mismo; sobrevivían diferencias de estilo entre archivos.

## Notas

Al adoptarlo, el código existente ya cumplía: la conversión no produjo un solo cambio de formato.

Se invoca con la versión pineada desde el módulo, sin instalación previa, de modo que el pre-commit funcione en cualquier máquina que tenga el toolchain del lenguaje.
