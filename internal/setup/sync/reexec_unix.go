//go:build !windows

package sync

import (
	"fmt"
	"os"
	"syscall"
)

// ReExec reemplaza el proceso actual por una nueva invocación del mismo binario
// usando syscall.Exec (exec(2)). Solo se llama después de reemplazar el binario
// en disco y cuando el proceso corre en una TTY Unix.
//
// Verifica que el ejecutable exista y sea ejecutable antes de intentar el exec.
// Si la verificación falla retorna el error; el caller decide cómo manejarlo.
func ReExec() error {
	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolviendo ejecutable: %w", err)
	}

	fi, err := os.Stat(self)
	if err != nil {
		return fmt.Errorf("verificando ejecutable %s: %w", self, err)
	}
	// Verifica que el archivo sea regular y tenga bit de ejecución para el dueño.
	if fi.Mode()&0o100 == 0 {
		return fmt.Errorf("ejecutable %s no tiene bit de ejecución", self)
	}

	return syscall.Exec(self, os.Args, os.Environ())
}
