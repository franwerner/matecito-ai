# Capability — Configurar la búsqueda de registros

- **Status:** Accepted
- **Date:** 2026-09-02
- **Components:** cli

## Propósito

Configuración única por proyecto del camino de búsqueda semántica: una colección por store de registros, sus exclusiones, sus credenciales de embedding y el primer índice. Se ejecuta una vez por proyecto y de nuevo cuando los stores del proyecto cambian — nunca como parte de una búsqueda.

## Actores

- **La persona que mantiene el proyecto** — este proceso es declarativo, no un motor usado por otros flujos.

## Requisitos

### Requisito: Una colección por store, nombrada exactamente como mapea

El sistema DEBE crear exactamente una colección por store, y DEBE usar el nombre de colección que el mapeo de stores devolvió, al pie de la letra. NO DEBE agrupar stores por tipo ni por componente: el alcance por-store es lo que impide que una búsqueda de comportamiento devuelva decisiones. El nombre mapeado está prefijado por proyecto porque un índice de búsqueda único sirve cada proyecto desde una configuración — un nombre sin prefijo declarado por un segundo proyecto reemplaza el del primer proyecto, silenciosamente, dejando una colección apuntando al repositorio equivocado.

#### Scenario: un repositorio con cuatro stores

- **GIVEN** un proyecto con un store de decisiones raíz, uno de comportamiento raíz, y dos stores de decisiones de componentes
- **WHEN** la configuración se ejecuta
- **THEN** existen cuatro colecciones, cada una nombrada como el mapeo la devolvió

#### Scenario: el nombre viene del mapeo, no de la carpeta

- **GIVEN** dos proyectos cuyos stores de decisiones usan nombres de carpeta distintos para la misma clase de registro
- **WHEN** la configuración de cada proyecto se ejecuta
- **THEN** ambas colecciones son nombradas por clase y componente, no por nombre de carpeta, y ninguna reemplaza la otra

### Requisito: Los índices de navegación están excluidos de cada colección

Cada colección DEBE excluir el índice de navegación de su store. Esos archivos hacen match con todo y desplazan registros reales de la parte alta de un conjunto de resultados; el índice se lee entero, no se busca.

#### Scenario: el índice del store no aparece en resultados

- **GIVEN** una colección configurada sobre un store que tiene un índice de navegación
- **WHEN** una búsqueda semántica se ejecuta contra ella
- **THEN** el índice de navegación no está entre los resultados

### Requisito: Las credenciales vienen del ambiente, nunca de la configuración versionada

Las credenciales de embedding DEBEN ser suministradas por el ambiente o por la configuración global no-versionada de la persona. NO DEBEN ser escritas en el archivo de configuración versionado del proyecto. Sin credenciales, la configuración DEBE parar antes de indexar y decirlo, en lugar de comenzar y dejar un índice a mitad de construir.

#### Scenario: credenciales presentes

- **GIVEN** el endpoint de embedding, la clave y el modelo suministrados por el ambiente
- **WHEN** la configuración se ejecuta
- **THEN** procede a indexar, y ninguna credencial se escribe en archivo versionado alguno

#### Scenario: credenciales ausentes

- **GIVEN** sin credenciales de embedding disponibles
- **WHEN** la configuración se ejecuta
- **THEN** para antes de indexar, dice que faltan credenciales, y no deja índice parcial alguno

### Requisito: La configuración completa solo después de indexar e incrustar

Después de las colecciones y las credenciales, el sistema DEBE recorrer los archivos y luego generar los vectores. La configuración está completa solo una vez que ambas han terminado.

#### Scenario: configuración reportada como completa

- **GIVEN** colecciones creadas y credenciales disponibles
- **WHEN** la configuración termina
- **THEN** tanto el recorrido de archivos como la generación de vectores han ocurrido, y el camino semántico responde consultas

### Requisito: Los stores que cambian hacen obsoletas las colecciones, y re-ejecutar lo detecta

Cuando el proyecto gana o pierde un componente con su propio store, las colecciones se desalinean y los registros de ese componente dejan de aparecer — silenciosamente. Re-ejecutar el mapeo de stores DEBE reportar un store sin colección coincidente, que es el disparador para ejecutar esta configuración de nuevo.

#### Scenario: un store nuevo de componente aparece

- **GIVEN** un proyecto que agrega un componente con su propio store de registros después de que la configuración se ejecutó
- **WHEN** el mapeo de stores se ejecuta de nuevo
- **THEN** muestra un store sin colección correspondiente

## Referencias

- **Flow** → [`../flow/find-durable-records.md`](../flow/find-durable-records.md) — la capacidad que consume el camino que esta configuración activa, y que degrada cuando no está disponible.
