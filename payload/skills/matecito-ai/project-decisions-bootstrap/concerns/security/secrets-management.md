---
name: secrets-management
depth: light
domain: security
tipo: política
adr-output: secrets-management
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

ADR `secrets-management` con: backend elegido, política de rotación (manual / automática / no definida), lista explícita de qué NUNCA se commitea (ej: `.env`, archivos de certificados, API keys), y qué NUNCA se loggea (passwords, tokens, PII). Incluir referencia a `.gitignore` o escaneo de secretos en CI si aplica.
