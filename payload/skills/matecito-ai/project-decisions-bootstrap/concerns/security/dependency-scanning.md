---
name: dependency-scanning
depth: light
domain: security
tipo: política
adr-output: dependency-scanning
source: OWASP Top 10 A06:2021 (Vulnerable and Outdated Components)
---

# Fase: Escaneo de dependencias

## Qué decide

Qué herramienta detecta vulnerabilidades conocidas en dependencias y si ese escaneo corre de forma automática en el pipeline de CI.

## Preguntas

### 1. Herramienta y punto de ejecución

> Las dependencias con CVEs conocidos son el vector A06 del OWASP Top 10. Sin escaneo automático, las vulnerabilidades se acumulan silenciosamente.

- **`npm audit` / `pip audit` / `cargo audit` (nativo del ecosistema) en CI** — *default para proyectos que no usan GitHub.*
- Dependabot (GitHub) — PRs automáticos con actualizaciones de seguridad; cero configuración en repos de GitHub.
- Snyk — escaneo con más contexto de explotabilidad; requiere cuenta.
- No aplica — proyecto sin dependencias externas o sin entorno productivo.
- No sé, recomendame.

## Notas de lógica (para el motor)

- Si el repositorio ya está en GitHub, sugerir Dependabot como opción de menor fricción.
- Si el usuario elige "No aplica" para un proyecto productivo, registrar la razón explícita en el ADR.

## Tech a registrar

Si se elige Snyk u otra herramienta externa, registrarla en el catálogo `tech/`.

## Qué materializar

ADR `dependency-scanning` con: herramienta elegida, si corre en CI (y en qué etapa: PR check, merge a main, schedule), política ante vulnerabilidades críticas (bloquear build vs notificar), y si existe un proceso de revisión periódica de dependencias desactualizadas.
