//go:build darwin

package platform

import "golang.org/x/sys/unix"

func flushTTYInput(fd int) {
	if _, err := unix.IoctlGetTermios(fd, unix.TIOCGETA); err != nil {
		return
	}
	// TIOCFLUSH takes a pointer to FREAD/FWRITE flags; darwin's TCIFLUSH equals FREAD.
	_ = unix.IoctlSetPointerInt(fd, unix.TIOCFLUSH, unix.TCIFLUSH)
}
