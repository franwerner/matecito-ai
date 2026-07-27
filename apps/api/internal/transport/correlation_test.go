package transport

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestWithCorrelation_NoValuesYet(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	correlated := WithCorrelation(logger, "", "", "")
	correlated.Info("boot")

	out := buf.String()
	if strings.Contains(out, "project=") || strings.Contains(out, "change=") || strings.Contains(out, "agent=") {
		t.Fatalf("expected no correlation fields, got: %s", out)
	}
}

func TestWithCorrelation_WithValues(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	correlated := WithCorrelation(logger, "proj-1", "change-1", "agent-1")
	correlated.Info("event")

	out := buf.String()
	for _, want := range []string{"project=proj-1", "change=change-1", "agent=agent-1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected log to contain %q, got: %s", want, out)
		}
	}
}
