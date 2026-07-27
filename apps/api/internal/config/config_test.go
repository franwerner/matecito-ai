package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func noEnv(string) string { return "" }

func TestLoad_Defaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home) // os.UserHomeDir reads $HOME on linux

	cfg, err := Load(nil, noEnv)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Port != 4300 {
		t.Errorf("Port = %d, want 4300", cfg.Port)
	}
	want := filepath.Join(home, ".matecito-ai")
	if cfg.PersistDir != want {
		t.Errorf("PersistDir = %q, want %q", cfg.PersistDir, want)
	}
	if cfg.LogFormat != "text" {
		t.Errorf("LogFormat = %q, want text", cfg.LogFormat)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want info", cfg.LogLevel)
	}
	if _, err := os.Stat(want); err != nil {
		t.Errorf("expected default persist dir to be created: %v", err)
	}
}

func TestLoad_FlagOverridesDefault(t *testing.T) {
	persistDir := t.TempDir()
	cfg, err := Load([]string{"-port=5555", "-persist-dir=" + persistDir}, noEnv)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Port != 5555 {
		t.Errorf("Port = %d, want 5555", cfg.Port)
	}
	if cfg.PersistDir != persistDir {
		t.Errorf("PersistDir = %q, want %q", cfg.PersistDir, persistDir)
	}
}

func TestLoad_EnvOverridesDefault(t *testing.T) {
	persistDir := t.TempDir()
	env := map[string]string{
		"BROKER_PORT":        "6000",
		"BROKER_PERSIST_DIR": persistDir,
	}
	cfg, err := Load(nil, func(k string) string { return env[k] })
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Port != 6000 {
		t.Errorf("Port = %d, want 6000", cfg.Port)
	}
}

func TestLoad_FlagWinsOverEnv(t *testing.T) {
	persistDir := t.TempDir()
	env := map[string]string{"BROKER_PORT": "6000"}
	cfg, err := Load([]string{"-port=7000", "-persist-dir=" + persistDir}, func(k string) string { return env[k] })
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Port != 7000 {
		t.Errorf("Port = %d, want 7000 (flag must win over env)", cfg.Port)
	}
}

func TestLoad_InvalidPortAborts(t *testing.T) {
	persistDir := t.TempDir()
	_, err := Load([]string{"-port=70000", "-persist-dir=" + persistDir}, noEnv)
	if !errors.Is(err, ErrInvalidPort) {
		t.Fatalf("err = %v, want wrapping ErrInvalidPort", err)
	}
}

func TestLoad_UnwritablePersistDirAborts(t *testing.T) {
	// A regular file where a directory is expected makes MkdirAll fail
	// reliably regardless of the running user (unlike a permission-bit
	// check, which a root process would bypass).
	tmp := t.TempDir()
	blocked := filepath.Join(tmp, "not-a-dir")
	if err := os.WriteFile(blocked, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := Load([]string{"-persist-dir=" + blocked}, noEnv)
	if !errors.Is(err, ErrPersistDirNotWritable) {
		t.Fatalf("err = %v, want wrapping ErrPersistDirNotWritable", err)
	}
}

func TestLoad_InvalidLogFormatAborts(t *testing.T) {
	persistDir := t.TempDir()
	_, err := Load([]string{"-persist-dir=" + persistDir, "-log-format=xml"}, noEnv)
	if !errors.Is(err, ErrInvalidLogFormat) {
		t.Fatalf("err = %v, want wrapping ErrInvalidLogFormat", err)
	}
}

func TestLoad_InvalidLogLevelAborts(t *testing.T) {
	persistDir := t.TempDir()
	_, err := Load([]string{"-persist-dir=" + persistDir, "-log-level=verbose"}, noEnv)
	if !errors.Is(err, ErrInvalidLogLevel) {
		t.Fatalf("err = %v, want wrapping ErrInvalidLogLevel", err)
	}
}

func TestLoad_SentinelSurvivesWrapping(t *testing.T) {
	persistDir := t.TempDir()
	_, err := Load([]string{"-persist-dir=" + persistDir, "-port=abc"}, noEnv)
	if err == nil {
		t.Fatal("expected error")
	}
	doubleWrapped := fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", err))
	if !errors.Is(doubleWrapped, ErrInvalidPort) {
		t.Fatal("errors.Is failed after extra wrapping layers")
	}
}
