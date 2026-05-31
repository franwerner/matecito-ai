package releasedl_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/franwerner/matecito-ai/internal/setup/releasedl"
)

// apiPayload construye una respuesta mínima válida de la GitHub API de releases.
func apiPayload(tag, assetName, assetURL, sumsURL string) []byte {
	payload := map[string]any{
		"tag_name": tag,
		"assets": []map[string]any{
			{"name": assetName, "browser_download_url": assetURL},
			{"name": "checksums.txt", "browser_download_url": sumsURL},
		},
	}
	b, _ := json.Marshal(payload)
	return b
}

func TestLatestReleaseWithTimeout_NormalResponse(t *testing.T) {
	repo := releasedl.Repo{Owner: "owner", Name: "myrepo", Binary: "mybinary"}
	plat := releasedl.Platform{OS: "linux", Arch: "amd64", Ext: "tar.gz"}

	tag := "v1.2.3"
	version := strings.TrimPrefix(tag, "v")
	assetName := "mybinary_" + version + "_linux_amd64.tar.gz"
	assetURL := "https://example.com/download/" + assetName
	sumsURL := "https://example.com/download/checksums.txt"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(apiPayload(tag, assetName, assetURL, sumsURL))
	}))
	defer srv.Close()

	rel, err := releasedl.LatestReleaseWithTimeout(repo, plat, 5*time.Second, srv.URL)
	if err != nil {
		t.Fatalf("respuesta normal: error inesperado: %v", err)
	}
	if rel.Tag != tag {
		t.Errorf("tag: got %q, want %q", rel.Tag, tag)
	}
	if rel.AssetURL != assetURL {
		t.Errorf("assetURL: got %q, want %q", rel.AssetURL, assetURL)
	}
	if rel.SumsURL != sumsURL {
		t.Errorf("sumsURL: got %q, want %q", rel.SumsURL, sumsURL)
	}
}

func TestLatestReleaseWithTimeout_TimeoutError(t *testing.T) {
	repo := releasedl.Repo{Owner: "owner", Name: "myrepo", Binary: "mybinary"}
	plat := releasedl.Platform{OS: "linux", Arch: "amd64", Ext: "tar.gz"}

	// El servidor duerme más que el timeout antes de responder.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := releasedl.LatestReleaseWithTimeout(repo, plat, 50*time.Millisecond, srv.URL)
	if err == nil {
		t.Fatal("se esperaba error de timeout pero no hubo error")
	}
}
