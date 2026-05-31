package releasedl

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const httpTimeout = 60 * time.Second

// Repo describe un repositorio de GitHub que publica releases con assets
// goreleaser-style: <binary>_<version>_<os>_<arch>.(tar.gz|zip) + checksums.txt.
type Repo struct {
	Owner  string
	Name   string
	Binary string // nombre del ejecutable dentro del archive (sin extensión)
}

// Repos predefinidos.
var (
	EngramRepo = Repo{
		Owner:  "Gentleman-Programming",
		Name:   "engram",
		Binary: "engram",
	}
	MatecitoRepo = Repo{
		Owner:  "franwerner",
		Name:   "matecito-ai",
		Binary: "matecito-ai",
	}
)

type Release struct {
	Tag       string
	AssetURL  string
	AssetName string
	SumsURL   string
}

type Platform struct {
	OS   string
	Arch string
	Ext  string
}

// Detect resuelve el OS/arch actual al formato usado por los assets goreleaser
// (ej: linux_amd64, darwin_arm64, windows_amd64).
func Detect() (Platform, error) {
	p := Platform{OS: runtime.GOOS, Arch: runtime.GOARCH}
	switch p.OS {
	case "linux", "darwin":
		p.Ext = "tar.gz"
	case "windows":
		p.Ext = "zip"
	default:
		return p, fmt.Errorf("OS no soportado: %s", p.OS)
	}
	switch p.Arch {
	case "amd64", "arm64":
	default:
		return p, fmt.Errorf("arquitectura no soportada: %s", p.Arch)
	}
	return p, nil
}

// LatestRelease consulta la GitHub API de repo y devuelve la URL del asset
// que corresponde a la plataforma indicada y la URL del checksums.txt.
// Usa el timeout predeterminado de 60s.
func LatestRelease(repo Repo, p Platform) (Release, error) {
	return LatestReleaseWithTimeout(repo, p, httpTimeout)
}

// LatestReleaseWithTimeout es idéntica a LatestRelease pero acepta un timeout
// configurable; permite a la TUI usar un timeout corto (~5s) sin afectar
// las rutas de install/update que necesitan más tiempo.
// El parámetro variadic apiBaseURL existe solo para tests (httptest.Server);
// en producción se omite y se usa la API real de GitHub.
func LatestReleaseWithTimeout(repo Repo, p Platform, timeout time.Duration, apiBaseURL ...string) (Release, error) {
	var rel Release

	base := "https://api.github.com"
	if len(apiBaseURL) > 0 && apiBaseURL[0] != "" {
		base = apiBaseURL[0]
	}
	apiURL := fmt.Sprintf("%s/repos/%s/%s/releases/latest", base, repo.Owner, repo.Name)

	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return rel, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	token := githubToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return rel, fmt.Errorf("consultando GitHub API: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			hint := "esperá unos minutos o seteá GITHUB_TOKEN/GH_TOKEN para subir el límite (60 → 5000 req/h)"
			if token != "" {
				hint = "ya hay un token seteado; puede ser un secondary rate limit, esperá unos minutos"
			}
			return rel, fmt.Errorf("rate limit de GitHub API (status %d) consultando releases de %s/%s — %s", resp.StatusCode, repo.Owner, repo.Name, hint)
		}
		return rel, fmt.Errorf("GitHub API devolvió status %d consultando releases de %s/%s", resp.StatusCode, repo.Owner, repo.Name)
	}

	var payload struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return rel, fmt.Errorf("parseando respuesta de GitHub: %w", err)
	}

	rel.Tag = payload.TagName
	version := strings.TrimPrefix(payload.TagName, "v")
	assetName := fmt.Sprintf("%s_%s_%s_%s.%s", repo.Binary, version, p.OS, p.Arch, p.Ext)

	for _, a := range payload.Assets {
		if a.Name == assetName {
			rel.AssetURL = a.URL
			rel.AssetName = a.Name
		}
		if a.Name == "checksums.txt" {
			rel.SumsURL = a.URL
		}
	}
	releasesURL := fmt.Sprintf("https://github.com/%s/%s/releases", repo.Owner, repo.Name)
	if rel.AssetURL == "" {
		available := make([]string, 0, len(payload.Assets))
		for _, a := range payload.Assets {
			available = append(available, a.Name)
		}
		return rel, fmt.Errorf(
			"no se encontró el asset %q en la última release (%s) de %s/%s — el esquema de nombres del release pudo haber cambiado.\nAssets disponibles: %s\nDescargá manualmente desde %s",
			assetName, payload.TagName, repo.Owner, repo.Name,
			strings.Join(available, ", "), releasesURL)
	}
	if rel.SumsURL == "" {
		return rel, fmt.Errorf(
			"no se encontró checksums.txt en la última release (%s) de %s/%s.\nDescargá manualmente desde %s",
			payload.TagName, repo.Owner, repo.Name, releasesURL)
	}
	return rel, nil
}

func githubToken() string {
	if t := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); t != "" {
		return t
	}
	return strings.TrimSpace(os.Getenv("GH_TOKEN"))
}

// Download baja el asset, verifica SHA256 contra checksums.txt, extrae el
// binario (repo.Binary o repo.Binary.exe en Windows) y lo deja en destBinary.
func Download(repo Repo, rel Release, destBinary string, out io.Writer) error {
	tmpDir, err := os.MkdirTemp("", "releasedl-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, rel.AssetName)
	fmt.Fprintf(out, "  Descargando %s\n", rel.AssetName)
	if err := downloadFile(rel.AssetURL, archivePath); err != nil {
		return fmt.Errorf("descargando asset: %w", err)
	}

	fmt.Fprintln(out, "  Verificando checksum SHA256")
	sumsPath := filepath.Join(tmpDir, "checksums.txt")
	if err := downloadFile(rel.SumsURL, sumsPath); err != nil {
		return fmt.Errorf("descargando checksums.txt: %w", err)
	}
	if err := verifyChecksum(archivePath, sumsPath, rel.AssetName); err != nil {
		return err
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return err
	}
	binaryName := repo.Binary
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	extracted, err := extractBinary(archivePath, extractDir, binaryName)
	if err != nil {
		return fmt.Errorf("extrayendo binario: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destBinary), 0o755); err != nil {
		return err
	}
	if err := moveFile(extracted, destBinary); err != nil {
		return fmt.Errorf("instalando binario en %s: %w", destBinary, err)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(destBinary, 0o755); err != nil {
			return err
		}
	}
	fmt.Fprintf(out, "  Instalado: %s (release %s)\n", destBinary, rel.Tag)
	return nil
}

// DefaultBinaryPath devuelve la ruta canónica donde instalar el binario:
// ~/.local/bin/<name> en POSIX, %LOCALAPPDATA%\matecito-ai\bin\<name>.exe en Windows.
func DefaultBinaryPath(repo Repo) (string, error) {
	name := repo.Binary
	if runtime.GOOS == "windows" {
		name += ".exe"
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", errors.New("LOCALAPPDATA no está definido")
		}
		return filepath.Join(localAppData, "matecito-ai", "bin", name), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "bin", name), nil
}

func downloadFile(url, dest string) error {
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d bajando %s", resp.StatusCode, url)
	}
	tmp := dest + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dest)
}

func verifyChecksum(archivePath, sumsPath, assetName string) error {
	data, err := os.ReadFile(sumsPath)
	if err != nil {
		return err
	}
	var expected string
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == assetName {
			expected = fields[0]
			break
		}
	}
	if expected == "" {
		return fmt.Errorf("checksum para %s no encontrado en checksums.txt", assetName)
	}
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(got, expected) {
		return fmt.Errorf("checksum no coincide: esperado %s, obtenido %s", expected, got)
	}
	return nil
}

func extractBinary(archivePath, destDir, binaryName string) (string, error) {
	if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir, binaryName)
	}
	return extractTarGz(archivePath, destDir, binaryName)
}

func extractTarGz(archivePath, destDir, binaryName string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if filepath.Base(hdr.Name) != binaryName {
			continue
		}
		dest := filepath.Join(destDir, binaryName)
		out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return "", err
		}
		if err := out.Close(); err != nil {
			return "", err
		}
		return dest, nil
	}
	return "", fmt.Errorf("binario %q no encontrado dentro del tar.gz", binaryName)
}

func extractZip(archivePath, destDir, binaryName string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()
	for _, zf := range r.File {
		if filepath.Base(zf.Name) != binaryName {
			continue
		}
		rc, err := zf.Open()
		if err != nil {
			return "", err
		}
		dest := filepath.Join(destDir, binaryName)
		out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
		if err != nil {
			rc.Close()
			return "", err
		}
		if _, err := io.Copy(out, rc); err != nil {
			rc.Close()
			out.Close()
			return "", err
		}
		rc.Close()
		if err := out.Close(); err != nil {
			return "", err
		}
		return dest, nil
	}
	return "", fmt.Errorf("binario %q no encontrado dentro del zip", binaryName)
}

func moveFile(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Remove(src)
}
