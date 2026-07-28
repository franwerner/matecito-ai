# Capability — Navegación del cockpit

- **Status:** Accepted
- **Date:** 2026-07-27

## Propósito

Permitir al usuario recorrer las secciones del cockpit y conocer de un vistazo la identidad del producto y el estado de la corrida en curso.

## Actores

- **Usuario** — persona que opera el cockpit en el navegador; dispara la navegación.

## Flujo principal

1. El usuario abre el cockpit y ve tres zonas: encabezado con la identidad del producto y el estado de la corrida, rail de navegación a la izquierda, y área de contenido con franja de timeline.
2. El rail ofrece cinco destinos en orden fijo: Canvas, Diffs, Decisiones, Cola y Ajustes (Ajustes anclado al pie); exactamente uno está marcado como activo.
3. El usuario elige un destino (con mouse o teclado) y el cockpit lo registra como la nueva sección activa.

## Ramas / flujos alternativos

- **Corrida activa transmitiendo en vivo** → el encabezado muestra su identificador y el indicador "en vivo".
- **Transmisión detenida** → el indicador pasa a estado neutro/detenido; no se oculta.
- **Sin corrida activa** → el lugar del identificador muestra un placeholder; el resto del encabezado no se altera.

## Casos borde

- **Contenido más alto que el viewport** → scrollea solo el área de contenido; encabezado, rail y timeline permanecen visibles.
- **Sección activa inválida (fuera de los cinco destinos)** → ningún destino aparece activo y el rail sigue operable.
- **Accesibilidad** → cada destino tiene nombre accesible en español, es operable por teclado con foco visible, y el activo se distingue por un medio no cromático además del color.

## Errores de cara al actor

- **Sin datos de corrida disponibles** → el cockpit se muestra completo con los estados neutro/placeholder del encabezado; la navegación no se bloquea ni se muestran errores.

## Escenarios

### Scenario: Selección de destino

- **GIVEN** el cockpit abierto con "Canvas" como sección activa
- **WHEN** el usuario elige "Diffs" en el rail
- **THEN** el cockpit registra "Diffs" como la nueva sección activa y el rail la marca como tal

### Scenario: Operación por teclado

- **GIVEN** el foco del teclado dentro del cockpit
- **WHEN** el usuario tabula hasta un destino del rail y presiona Enter o Space
- **THEN** el destino muestra foco visible y queda registrado como sección activa

### Scenario: Corrida en vivo

- **GIVEN** una corrida activa transmitiendo en vivo
- **WHEN** el usuario mira el encabezado
- **THEN** ve el identificador de la corrida y el indicador "en vivo"

### Scenario: Sin corrida activa

- **GIVEN** ninguna corrida activa
- **WHEN** el usuario mira el encabezado
- **THEN** ve un placeholder en el lugar del identificador y el encabezado conserva su disposición

### Scenario: Scroll del contenido

- **GIVEN** un contenido más alto que el viewport
- **WHEN** el usuario scrollea
- **THEN** solo se desplaza el área de contenido, con encabezado, rail y timeline siempre visibles

### Scenario: Sección activa inválida

- **GIVEN** una sección activa que no pertenece a los cinco destinos
- **WHEN** se muestra el rail
- **THEN** ningún destino aparece activo y el usuario puede seguir eligiendo destinos

## Referencias

- **EDR** → [`../../../apps/ui/.matecito-ai/edr/frontend/accessibility.md`](../../../apps/ui/.matecito-ai/edr/frontend/accessibility.md) — criterios WCAG 2.2 AA que gobiernan el comportamiento accesible del rail y el foco.
- **EDR** → [`../../../apps/ui/.matecito-ai/edr/frontend/routing.md`](../../../apps/ui/.matecito-ai/edr/frontend/routing.md) — el cableado de la navegación a rutas reales es una decisión posterior; este flujo no lo presupone.
- **Contexto de negocio** → cambio SDD `ui-cockpit-shell` (artefactos en Engram, topic `sdd/ui-cockpit-shell/*`).
