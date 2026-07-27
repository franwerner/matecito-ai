package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/franwerner/matecito-ai/apps/api/internal/config"
)

func TestNewLogger_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := newLogger(config.Config{LogFormat: "text", LogLevel: "info"}, &buf)
	logger.Info("hello")

	if json.Valid(buf.Bytes()) {
		t.Fatalf("expected non-JSON text output, got valid JSON: %s", buf.String())
	}
}

func TestNewLogger_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := newLogger(config.Config{LogFormat: "json", LogLevel: "info"}, &buf)
	logger.Info("hello")

	if !json.Valid(buf.Bytes()) {
		t.Fatalf("expected valid JSON output, got: %s", buf.String())
	}
}

func TestNewLogger_LevelFilter(t *testing.T) {
	var buf bytes.Buffer
	logger := newLogger(config.Config{LogFormat: "text", LogLevel: "info"}, &buf)

	logger.Debug("hidden")
	if buf.Len() != 0 {
		t.Fatalf("expected debug log to be filtered out, got: %s", buf.String())
	}

	logger.Warn("visible")
	if buf.Len() == 0 {
		t.Fatal("expected warn log to be emitted")
	}
}
