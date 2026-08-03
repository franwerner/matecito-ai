package sync

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// syncStateFile es el contrato JSON del archivo de estado del sincronizador.
// Se persiste en SyncStatePath() con escritura atómica (temp + rename).
type syncStateFile struct {
	LastCheck string `json:"lastCheck"` // RFC3339

	// PendingSync marca que una corrida anterior de la interfaz interactiva
	// terminó reemplazando el propio ejecutable y difirió los componentes
	// dependientes del payload embebido (deploy, config ecosistema) al
	// arranque siguiente. omitempty mantiene el archivo idéntico a hoy en
	// el caso común (sin marca).
	PendingSync bool `json:"pendingSync,omitempty"`
}

// SyncStatePath devuelve la ruta canónica de sync-state.json: ~/.matecito-ai/sync-state.json.
func SyncStatePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".matecito-ai", "sync-state.json"), nil
}

// ShouldCheck es una función pura que determina si debe realizarse una nueva
// comprobación de versiones.
//
// Retorna true cuando:
//   - lastCheck es el zero time (nunca se chequeó).
//   - El tiempo transcurrido desde lastCheck es mayor o igual a interval.
func ShouldCheck(now, lastCheck time.Time, interval time.Duration) bool {
	if lastCheck.IsZero() {
		return true
	}
	return now.Sub(lastCheck) >= interval
}

// loadStateFile lee y parsea el archivo de estado completo. Un archivo
// ausente (ENOENT) se lee como un valor fresco (err == nil): el estado
// natural antes de que se haya persistido nada. Cualquier otro problema —
// permiso denegado, JSON truncado o inválido — vuelve como un err no-nil con
// valor cero; los llamadores de read-modify-write deben tratar eso como
// "el estado existe pero no se puede confiar en él" y negarse a pisarlo
// (no-clobber), en vez de partir de un estado fresco.
func loadStateFile(path string) (syncStateFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return syncStateFile{}, nil
		}
		return syncStateFile{}, err
	}
	var f syncStateFile
	if err := json.Unmarshal(data, &f); err != nil {
		return syncStateFile{}, err
	}
	return f, nil
}

// saveStateFile persiste f en path con escritura atómica: escribe a un
// archivo temporal en el mismo directorio y luego lo renombra. Crea el
// directorio padre con permisos 0o755 si no existe. El archivo resultante
// tiene permisos 0o644.
func saveStateFile(path string, f syncStateFile) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.Marshal(f)
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// mutateStateFile implementa read-modify-write sobre el archivo de estado
// compartido: carga (o parte de un valor fresco si está ausente), aplica
// mutate, y guarda. Honra el no-clobber: cuando el archivo existe pero no se
// puede leer ni interpretar, devuelve ese error sin escribir nada — perder
// esta escritura es preferible a pisar el resto del estado persistido con
// una hoja en blanco.
func mutateStateFile(path string, mutate func(*syncStateFile)) error {
	f, err := loadStateFile(path)
	if err != nil {
		return err
	}
	mutate(&f)
	return saveStateFile(path, f)
}

// LoadSyncState lee el timestamp de la última comprobación desde path.
//
// Ausente, ilegible o corrupto → (zero time, nil): en los tres casos se
// trata como si no hubiera estado previo; ShouldCheck lo interpretará como
// "hay que chequear".
func LoadSyncState(path string) (time.Time, error) {
	f, err := loadStateFile(path)
	if err != nil {
		return time.Time{}, nil
	}
	if f.LastCheck == "" {
		return time.Time{}, nil
	}

	t, err := time.Parse(time.RFC3339, f.LastCheck)
	if err != nil {
		// timestamp malformado → tratar como ausente
		return time.Time{}, nil
	}
	return t, nil
}

// SaveSyncState persiste t como lastCheck vía read-modify-write, preservando
// el resto del estado (p.ej. la marca de sincronización pendiente). Antes
// pisaba el archivo entero con un syncStateFile nuevo.
func SaveSyncState(path string, t time.Time) error {
	return mutateStateFile(path, func(f *syncStateFile) {
		f.LastCheck = t.UTC().Format(time.RFC3339)
	})
}

// LoadPendingSync reporta si la marca de sincronización diferida está
// puesta en el archivo de estado en path. Ausente, ilegible o corrupto leen
// todos como false — sin marca — igual que LoadSyncState tolera las mismas
// condiciones para lastCheck.
func LoadPendingSync(path string) bool {
	f, err := loadStateFile(path)
	if err != nil {
		return false
	}
	return f.PendingSync
}

// SetPendingSync actualiza la marca de sincronización diferida vía
// read-modify-write, preservando el resto del estado persistido (p.ej.
// lastCheck). Honra el no-clobber: si el archivo existe pero no se puede
// leer ni interpretar, devuelve ese error sin escribir nada.
func SetPendingSync(path string, pending bool) error {
	return mutateStateFile(path, func(f *syncStateFile) {
		f.PendingSync = pending
	})
}
