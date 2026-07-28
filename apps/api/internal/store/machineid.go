package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// machineIDFileName is the machine identity file inside the daemon's
// persist-dir. Never the hostname (renames it, collides across machines) —
// a UUID v7 generated once on first startup and reused after (ROADMAP-2
// anexo, "project_paths").
const machineIDFileName = "machine-id"

// EnsureMachineID returns the daemon's machine_id, generating and
// persisting a new UUID v7 into persistDir on first run.
func EnsureMachineID(persistDir string) (string, error) {
	path := filepath.Join(persistDir, machineIDFileName)

	existing, err := os.ReadFile(path)
	switch {
	case err == nil:
		id := strings.TrimSpace(string(existing))
		if id == "" {
			return "", fmt.Errorf("store: %s is empty", path)
		}
		return id, nil
	case !os.IsNotExist(err):
		return "", fmt.Errorf("store: read %s: %w", path, err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("store: generate machine id: %w", err)
	}
	if err := os.WriteFile(path, []byte(id.String()+"\n"), 0o600); err != nil {
		return "", fmt.Errorf("store: write %s: %w", path, err)
	}
	return id.String(), nil
}
