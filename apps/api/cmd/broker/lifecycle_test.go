package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func noEnv(string) string { return "" }

func TestRun_InvalidConfigAbortsBeforeServing(t *testing.T) {
	persistDir := t.TempDir()
	var out strings.Builder

	code := run(context.Background(), []string{"-port=999999", "-persist-dir=" + persistDir}, noEnv, &out)

	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if !strings.Contains(out.String(), "invalid configuration") {
		t.Errorf("expected a message naming the offending value, got: %s", out.String())
	}
}

func TestRun_ServesHealthAndShutsDownCleanly(t *testing.T) {
	persistDir := t.TempDir()
	port := freePort(t)

	ctx, cancel := context.WithCancel(context.Background())
	args := []string{fmt.Sprintf("-port=%d", port), "-persist-dir=" + persistDir}

	done := make(chan int, 1)
	go func() {
		done <- run(ctx, args, noEnv, io.Discard)
	}()

	healthURL := fmt.Sprintf("http://127.0.0.1:%d/health", port)
	if err := waitForHealth(healthURL, 2*time.Second); err != nil {
		cancel()
		t.Fatalf("broker never became healthy: %v", err)
	}

	resp, err := http.Get(healthURL)
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	var body map[string]string
	decodeErr := json.NewDecoder(resp.Body).Decode(&body)
	func() { _ = resp.Body.Close() }()
	if decodeErr != nil {
		t.Fatalf("invalid JSON body: %v", decodeErr)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if body["process"] != "ok" || body["store"] != "ok" {
		t.Errorf("unexpected health body: %+v", body)
	}

	cancel()

	select {
	case code := <-done:
		if code != 0 {
			t.Errorf("exit code = %d, want 0 after clean shutdown", code)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("run() did not return after ctx cancellation (shutdown hung)")
	}

	if resp2, err := http.Get(healthURL); err == nil {
		func() { _ = resp2.Body.Close() }()
		t.Error("expected /health to be unreachable after shutdown")
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve free port: %v", err)
	}
	defer func() { _ = ln.Close() }()
	return ln.Addr().(*net.TCPAddr).Port
}

func waitForHealth(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			func() { _ = resp.Body.Close() }()
			return nil
		}
		lastErr = err
		time.Sleep(20 * time.Millisecond)
	}
	return lastErr
}
