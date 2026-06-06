---
name: secrets-management
depth: light
domain: security
type: policy
source: OWASP ASVS v4 §2.10 (Service Authentication) · 12-factor app §III (Config)
---

# Fase: Manejo de secretos

## Qué decide

Dónde se almacenan los secretos del sistema (credenciales, tokens, claves), cómo se rotan, y qué está explícitamente prohibido commitear o loggear.

## Preguntas

### 1. Backend de almacenamiento

> Un secreto en el repositorio es una brecha permanente — el historial de git persiste aunque se borre. OWASP ASVS 2.10.4 prohíbe credenciales hardcodeadas. 12-factor exige separación estricta de config y código.

- **Variables de entorno inyectadas en el proceso** — *default para proyectos simples; suficiente si el entorno es confiable.*
- Secret manager gestionado (AWS Secrets Manager, GCP Secret Manager, HashiCorp Vault) — para producción con múltiples secretos o rotación automática.
- Archivo `.env` local + variable de entorno en CI/CD — *solo para desarrollo; nunca se commitea.*
- No sé, recomendame.

## Notas de lógica (para el motor)

- Si el proyecto es `script` o `libreria` sin runtime propio, adaptar la pregunta al contexto del consumidor.
- Si el usuario elige secret manager, preguntar si ya tiene uno provisionado o si necesita elegir.

## Qué materializar

ADR `secrets-management` materializado según `~/.claude/references/adr/templates/adr.md`. Esta es una decisión de tipo `policy`; sus reglas deben quedar especialmente accionables. Debe contener:

- **Contexto**: por qué un secreto en el repositorio es una brecha permanente (el historial de git persiste aunque se borre), citando OWASP ASVS 2.10.4 (prohíbe credenciales hardcodeadas) y 12-factor §III (separación estricta de config y código).
- **Decisión**: backend de almacenamiento elegido (variables de entorno, secret manager gestionado, o `.env` local + var en CI/CD), y política de rotación (manual, automática, o no definida).
- **Reglas verificables** (cada una con su mecanismo):
  - `[manual]` lista explícita de qué NUNCA se commitea (ej. `.env`, archivos de certificados, API keys); todos esos patrones figuran en `.gitignore`.
  - `[manual]` qué NUNCA se loggea: passwords, tokens, PII.
  - `[tool: <escáner de secretos en CI>]` ningún secreto hardcodeado entra al repositorio (si hay escaneo en CI; omitir el mecanismo `[manual]` correspondiente si no lo hay).
  - `[manual]` los secretos se rotan según la política decidida (cadencia o trigger documentado).
- **Alternativas consideradas**: los otros backends evaluados y por qué no se eligieron para este nivel de sensibilidad.
- **Consecuencias**: confianza requerida en el entorno (para env vars) o dependencia del secret manager provisionado.
