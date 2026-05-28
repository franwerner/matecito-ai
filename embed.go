package matecitoai

import "embed"

// PayloadFS contiene los archivos de payload/ (skills, agents, CLAUDE.md)
// embebidos en el binario. El prefijo all: incluye archivos ocultos y
// directorios que empiezan con "." o "_".
//
//go:embed all:payload
var PayloadFS embed.FS
