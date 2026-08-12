# EDR — La instrucción de registro de prosa vive en un solo lugar

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
La instrucción de cómo redactar un ítem de buzón gate-facing o un resumen ejecutivo podía terminar copiada en cada uno de los cinco contratos de retorno que declaran el split summary/rationale, con el riesgo de desincronizarse en la próxima edición — el mismo patrón de duplicación que ya forzó a centralizar el contrato de retorno completo en una única sección.

## Decisión
La instrucción de registro vive en un único lugar — el protocolo compartido de fases (la sección que ya centraliza el contrato de retorno, para el ítem de buzón; y por referencia, para el resumen ejecutivo) — y cada contrato de retorno que declara el split summary/rationale apunta ahí con un puntero, sin restablecer el criterio.

## Reglas verificables
- **[manual]** Los cinco contratos de retorno que declaran el split summary/rationale contienen un puntero a la sección compartida, no una copia del criterio de registro.
- **[manual]** El protocolo compartido de fases es el único lugar donde el criterio de registro está escrito en extenso.

## Alternativas consideradas
Copiar la instrucción en cada uno de los cinco contratos. Descartado: reproduce el mismo patrón de duplicación que la centralización del contrato de retorno ya existe para evitar — la próxima edición del criterio desincroniza cinco copias en vez de una.

## Consecuencias
Editar el criterio de registro es un cambio en un solo archivo; los cinco contratos quedan como lectores, no como copias.

## Relacionados
- `relacionado-con` → [contract-pair-in-templates.md](contract-pair-in-templates.md) — el mismo razonamiento — un solo lugar de verdad, todo lo demás apunta — aplicado antes a la forma del contrato en vez de al registro de su prosa.
