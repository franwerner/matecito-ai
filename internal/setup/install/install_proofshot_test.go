package install_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/franwerner/matecito-ai/internal/setup/install"
)

// TestInstallProofshot_BinaryPresentAfterInstall verifica el happy path: npm
// sale 0 y proofshot queda en PATH → InstallProofshot retorna nil.
// El post-step `proofshot install` también sale 0 (mismo stub).
func TestInstallProofshot_BinaryPresentAfterInstall(t *testing.T) {
	d := tempDir(t)
	// npm config get prefix retorna ruta non-system → ensureUserNpmPrefix
	// reconfigura bin dir y llama EnsurePathInShell. npm install -g sale 0.
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	writeBin(t, d, "proofshot", 0, "")
	isolatedPATH(t, d)

	var out bytes.Buffer
	err := install.InstallProofshot(install.Options{Stdout: &out, Stderr: &out})
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

// TestInstallProofshot_BinaryAbsentAfterNpmSuccess verifica que cuando npm sale 0
// pero el binario proofshot no está en PATH, InstallProofshot retorna un error
// que menciona "proofshot" y "PATH".
func TestInstallProofshot_BinaryAbsentAfterNpmSuccess(t *testing.T) {
	d := tempDir(t)
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	// proofshot NO escrito en d → exec.LookPath("proofshot") debe fallar.
	isolatedPATH(t, d)

	var out bytes.Buffer
	err := install.InstallProofshot(install.Options{Stdout: &out, Stderr: &out})
	if err == nil {
		t.Fatal("expected error when proofshot binary is absent after npm install, got nil")
	}
	if !strings.Contains(err.Error(), "proofshot") {
		t.Errorf("error should mention 'proofshot'; got: %v", err)
	}
	if !strings.Contains(err.Error(), "PATH") {
		t.Errorf("error should mention 'PATH'; got: %v", err)
	}
}

// TestInstallProofshot_NpmMissing verifica que cuando npm no está en PATH,
// InstallProofshot retorna un error que menciona "npm".
func TestInstallProofshot_NpmMissing(t *testing.T) {
	d := tempDir(t)
	// npm NO escrito → ausente del PATH.
	writeBin(t, d, "proofshot", 0, "")
	isolatedPATH(t, d)

	var out bytes.Buffer
	err := install.InstallProofshot(install.Options{Stdout: &out, Stderr: &out})
	if err == nil {
		t.Fatal("expected error when npm is missing, got nil")
	}
	if !strings.Contains(err.Error(), "npm") {
		t.Errorf("error should mention 'npm'; got: %v", err)
	}
}

// TestInstallProofshot_NpmFails verifica que cuando npm install sale non-zero,
// InstallProofshot retorna error antes de llegar al LookPath de proofshot.
func TestInstallProofshot_NpmFails(t *testing.T) {
	d := tempDir(t)
	writeBin(t, d, "npm", 1, "")
	writeBin(t, d, "proofshot", 0, "")
	isolatedPATH(t, d)

	var out bytes.Buffer
	err := install.InstallProofshot(install.Options{Stdout: &out, Stderr: &out})
	if err == nil {
		t.Fatal("expected error when npm fails, got nil")
	}
	// No debe ser el error de LookPath guard.
	if strings.Contains(err.Error(), "PATH") {
		t.Errorf("error should be from npm, not LookPath guard; got: %v", err)
	}
}

// TestInstallProofshot_PostStepFails_ContinuesWithManualInstruction verifica que
// si el post-step `proofshot install` falla, InstallProofshot continúa (retorna nil)
// y emite una instrucción manual.
func TestInstallProofshot_PostStepFails_ContinuesWithManualInstruction(t *testing.T) {
	d := tempDir(t)
	writeBin(t, d, "npm", 0, "/tmp/npm-global")
	// proofshot stub: el primer llamado (LookPath tras install) lo resuelve porque
	// el binario existe; el segundo (proofshot install) sale 1 → post-step falla.
	writeBin(t, d, "proofshot", 1, "")
	isolatedPATH(t, d)

	var out bytes.Buffer
	err := install.InstallProofshot(install.Options{Stdout: &out, Stderr: &out})
	// InstallProofshot debe retornar nil (política best-effort para el post-step).
	if err != nil {
		t.Fatalf("expected nil (post-step failure is best-effort), got: %v", err)
	}
	// Debe emitir la instrucción manual.
	outStr := out.String()
	if !strings.Contains(outStr, "proofshot install") {
		t.Errorf("expected manual-instruction message mentioning 'proofshot install'; got: %q", outStr)
	}
}

// TestEnsureUserNpmPrefix_AlreadyUserLocal_EnsuresPathInShell verifica el fix del
// bug de PATH: cuando el prefix ya es user-local (no system-owned), ensureUserNpmPrefix
// igual debe asegurar el bin dir en PATH (sin retornar early).
//
// La forma de observar el efecto desde fuera es verificar que la variable de entorno
// PATH del proceso fue modificada para incluir el bin dir del prefix user-local,
// lo que es el efecto colateral visible de la llamada a os.Setenv dentro de ensureUserNpmPrefix.
func TestEnsureUserNpmPrefix_AlreadyUserLocal_EnsuresPathInShell(t *testing.T) {
	d := tempDir(t)
	// Prefix user-local: la ruta de home del usuario. Usamos un tempDir como prefix
	// simulado que no es un path del sistema.
	userPrefix := filepath.Join(d, "npm-prefix")
	if err := os.MkdirAll(filepath.Join(userPrefix, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}

	// npm config get prefix retorna userPrefix (non-system).
	// npm install -g sale 0.
	// proofshot está presente en d (para que LookPath no falle).
	writeBin(t, d, "proofshot", 0, "")

	// Escribimos el stub de npm que imprime el userPrefix en stdout.
	var sb strings.Builder
	sb.WriteString("#!/bin/sh\n")
	sb.WriteString("echo ")
	sb.WriteString(userPrefix)
	sb.WriteString("\n")
	npmPath := filepath.Join(d, "npm")
	if err := os.WriteFile(npmPath, []byte(sb.String()), 0o755); err != nil {
		t.Fatal(err)
	}

	// Solo d en PATH; userPrefix/bin está fuera del PATH inicial.
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", d)
	t.Cleanup(func() { os.Setenv("PATH", origPath) })

	expectedBinDir := filepath.Join(userPrefix, "bin")

	// Confirmar que expectedBinDir no está en PATH antes de la llamada.
	if strings.Contains(os.Getenv("PATH"), expectedBinDir) {
		t.Skip("expectedBinDir ya está en PATH antes del test; condición inicial no cumplida")
	}

	var out bytes.Buffer
	err := install.InstallProofshot(install.Options{Stdout: &out, Stderr: &out})
	if err != nil {
		t.Fatalf("InstallProofshot failed: %v", err)
	}

	// El PATH del proceso debe incluir ahora el bin dir del prefix user-local.
	if !strings.Contains(os.Getenv("PATH"), expectedBinDir) {
		t.Errorf("ensureUserNpmPrefix did not add user-local bin dir %q to PATH; got PATH=%q",
			expectedBinDir, os.Getenv("PATH"))
	}
}
