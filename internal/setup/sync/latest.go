package sync

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/franwerner/matecito-ai/internal/setup/releasedl"
)

// npmPackageCodeGraph es el paquete npm que instala el CLI de CodeGraph.
// Debe coincidir con el argumento que usa InstallCodegraph en install.go.
const npmPackageCodeGraph = "@colbymchenry/codegraph"

// fetchLatestMatecito obtiene la versión más reciente de matecito-ai desde GitHub.
// Devuelve la versión sin prefijo "v" (ej: "1.2.3").
// En caso de error retorna ("", err); el caller marca Unknown=true.
func fetchLatestMatecito(timeout time.Duration) (string, error) {
	plat, err := releasedl.Detect()
	if err != nil {
		return "", fmt.Errorf("detectando plataforma: %w", err)
	}
	rel, err := releasedl.LatestReleaseWithTimeout(releasedl.MatecitoRepo, plat, timeout)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(rel.Tag, "v"), nil
}

// fetchLatestEngram obtiene la versión más reciente de Engram desde GitHub.
// Devuelve la versión sin prefijo "v" (ej: "1.16.1").
// En caso de error retorna ("", err); el caller marca Unknown=true.
func fetchLatestEngram(timeout time.Duration) (string, error) {
	plat, err := releasedl.Detect()
	if err != nil {
		return "", fmt.Errorf("detectando plataforma: %w", err)
	}
	rel, err := releasedl.LatestReleaseWithTimeout(releasedl.EngramRepo, plat, timeout)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(rel.Tag, "v"), nil
}

// fetchLatestCodeGraph obtiene la versión más reciente del paquete npm de CodeGraph
// consultando el registry de npm directamente (HTTP GET, sin invocar el CLI de npm).
// Devuelve la versión semver (ej: "0.9.7").
// En caso de error retorna ("", err); el caller marca Unknown=true.
func fetchLatestCodeGraph(timeout time.Duration) (string, error) {
	encoded := strings.ReplaceAll(npmPackageCodeGraph, "/", "%2F")
	url := fmt.Sprintf("https://registry.npmjs.org/%s/latest", encoded)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("consultando npm registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("npm registry devolvió status %d para %s", resp.StatusCode, npmPackageCodeGraph)
	}

	var payload struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("parseando respuesta de npm registry: %w", err)
	}
	if payload.Version == "" {
		return "", fmt.Errorf("npm registry no devolvió versión para %s", npmPackageCodeGraph)
	}
	return payload.Version, nil
}
