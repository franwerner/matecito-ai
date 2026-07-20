//go:build windows

package sync

import (
	"os"
	"os/exec"
)

// ReExec spawns a new invocation of the same binary and waits for it, since
// Windows lacks exec(2)-style process replacement. A *exec.ExitError means
// the child ran and exited non-zero — not a spawn failure — so its exit
// code is propagated via os.Exit instead of being returned as an error.
func ReExec() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(self, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = resumeEnv(os.Environ())

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return err
	}
	os.Exit(cmd.ProcessState.ExitCode())
	return nil
}
