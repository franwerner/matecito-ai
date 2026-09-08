package platform

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFlushPendingInput_NonFileReaderIsNoop(t *testing.T) {
	r := strings.NewReader("y\n")
	FlushPendingInput(r)
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != "y\n" {
		t.Fatalf("data consumed: got %q, want %q", got, "y\n")
	}
}

func TestFlushPendingInput_RegularFileIsNoop(t *testing.T) {
	path := filepath.Join(t.TempDir(), "in")
	if err := os.WriteFile(path, []byte("y\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	FlushPendingInput(f)
	got, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != "y\n" {
		t.Fatalf("data consumed: got %q, want %q", got, "y\n")
	}
}

func TestFlushPendingInput_PipeIsNoop(t *testing.T) {
	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer pr.Close()
	if _, err := pw.WriteString("y\n"); err != nil {
		t.Fatal(err)
	}
	pw.Close()
	FlushPendingInput(pr)
	got, err := io.ReadAll(pr)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != "y\n" {
		t.Fatalf("data consumed: got %q, want %q", got, "y\n")
	}
}
