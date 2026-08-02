package deploy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// registryFileName is the deploy registry's name under StateDir().
const registryFileName = "deployed.json"

// Origin describes which of the three real payload sources produced a FileOp:
// a domain's own components, the shared/ tree (no active-domain gate), or the
// composed matecito-ai.md (core + domain index, generated in memory). It is
// descriptive only — deletion is decided by the op's FUENTE (registry vs the
// embedded legacy catalog, see FileOp and applyRemoval in deploy.go), never by
// Origin.
type Origin string

const (
	OriginShared   Origin = "shared"
	OriginComposed Origin = "composed"
)

// domainOrigin builds the "domain:<id>" Origin for a domain-sourced FileOp.
func domainOrigin(id string) Origin {
	return Origin("domain:" + id)
}

// Entry is one tracked destination in the registry: the path deploy actually
// wrote (relative to claudeHome, "/"-separated for a stable, OS-independent
// format), the hash of the bytes written, and its origin.
type Entry struct {
	Path   string `json:"path"`
	Hash   string `json:"hash"`
	Origin Origin `json:"origin"`
}

// Registry is the persisted record of everything the last successful Apply
// wrote. It holds no file content — only enough to plan a future deletion
// safely (path + hash) and to describe where each entry came from.
type Registry struct {
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updatedAt"`
	Entries   []Entry   `json:"entries"`
}

// StateDir devuelve ~/.matecito-ai — la raíz de estado del ecosistema donde
// vive el registro de despliegue (deployed.json), separada de claudeHome
// (~/.claude, el destino real de los archivos).
func StateDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".matecito-ai"), nil
}

// LoadRegistry lee el registro de stateDir. Ausente o corrupto (JSON inválido)
// se tratan por igual: found=false, sin error — esa corrida es una migración.
// Cualquier otro error de lectura (permisos, etc.) sí se propaga.
func LoadRegistry(stateDir string) (Registry, bool, error) {
	data, err := os.ReadFile(filepath.Join(stateDir, registryFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return Registry{}, false, nil
		}
		return Registry{}, false, err
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		// Corrupt registry == absent (New Decision: "Registro corrupto se
		// trata como ausente"). Never block a run over a broken file we own.
		return Registry{}, false, nil
	}
	return reg, true, nil
}

// SaveRegistry reescribe stateDir/deployed.json atómicamente (tmp + rename,
// 0o644), la misma convención que agentmodel.Save.
func SaveRegistry(stateDir string, reg Registry) error {
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return fmt.Errorf("deploy: mkdir %s: %w", stateDir, err)
	}
	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return fmt.Errorf("deploy: marshal registry: %w", err)
	}
	data = append(data, '\n')

	path := filepath.Join(stateDir, registryFileName)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("deploy: write tmp %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("deploy: rename %s → %s: %w", tmp, path, err)
	}
	return nil
}

// hashBytes hashea con SHA-256 (misma convención que releasedl.verifyChecksum).
func hashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// hashDiskFile hashea el contenido actual de path. El error de os.ReadFile se
// propaga tal cual (incluido "no existe"): el caller decide qué significa un
// archivo ausente en su contexto (Apply lo trata como "ya no está, sale del
// registro sin error").
func hashDiskFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return hashBytes(data), nil
}

// relToClaudeHome devuelve target relativo a claudeHome con separadores "/",
// el formato de Entry.Path.
func relToClaudeHome(target, claudeHome string) (string, error) {
	rel, err := filepath.Rel(claudeHome, target)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}
