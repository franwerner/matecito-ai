package sync

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// resumeEnvVar signals a resumed invocation that the on-disk binary was
// already swapped, so it must skip self-update and auto-confirm.
const resumeEnvVar = "MATECITO_RESUME"

// reExecFn is a substitution seam: FinishSelfReplace calls it instead of
// ReExec directly so tests can assert invocation/fallback without replacing
// the test process. The real handoff (syscall.Exec on unix, spawn+wait on
// windows) is not unit-testable — it would replace/spawn the test process
// itself — so it is exercised only via this seam plus manual verification.
var reExecFn = ReExec

// resumeEnv returns env with resumeEnvVar set, without mutating env or
// duplicating the flag if it is already present.
func resumeEnv(env []string) []string {
	flag := resumeEnvVar + "=1"
	for _, e := range env {
		if e == flag {
			return append([]string{}, env...)
		}
	}
	return append(append([]string{}, env...), flag)
}

// ResumeRequested reports whether this process was launched by a
// self-replace re-exec/spawn.
func ResumeRequested() bool {
	return os.Getenv(resumeEnvVar) == "1"
}

// FinishSelfReplace triggers the re-exec/spawn after a self-replace. It is a
// no-op when replaced is false. The engine never calls this — only CLI
// callers do, so the TUI can never self-exec.
//
// A reExecFn failure degrades to the manual-rerun message instead of a hard
// error, because the on-disk binary was already swapped.
func FinishSelfReplace(out, errOut io.Writer, replaced bool) error {
	if !replaced {
		return nil
	}
	if err := reExecFn(); err != nil {
		self, _ := exec.LookPath(os.Args[0])
		if self == "" {
			self = "matecito-ai"
		}
		fmt.Fprintf(out, "matecito-ai actualizado — re-ejecutá %s para usar la nueva versión.\n", self)
	}
	return nil
}
