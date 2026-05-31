//go:build windows

package sync

// ReExec en Windows no puede reemplazar el proceso en ejecución con exec(2)
// porque la plataforma no lo soporta. Retorna siempre nil; el caller en Sync
// imprimirá el aviso de reinicio manual al usuario.
func ReExec() error {
	return nil
}
