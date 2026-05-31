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

// LoadSyncState lee el timestamp de la última comprobación desde path.
//
// ENOENT → (zero time, nil): primer arranque, ShouldCheck lo tratará como check.
// JSON inválido → (zero time, nil): tratamos corrupción como si no hubiera estado.
func LoadSyncState(path string) (time.Time, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	var f syncStateFile
	if err := json.Unmarshal(data, &f); err != nil {
		// archivo corrupto → tratar como ausente
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

// SaveSyncState persiste t en path como JSON con escritura atómica:
// escribe a un archivo temporal en el mismo directorio y luego lo renombra.
// Crea el directorio padre con permisos 0o755 si no existe.
// El archivo resultante tiene permisos 0o644.
func SaveSyncState(path string, t time.Time) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	f := syncStateFile{LastCheck: t.UTC().Format(time.RFC3339)}
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
