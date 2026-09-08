//go:build linux

package platform

import "golang.org/x/sys/unix"

func flushTTYInput(fd int) {
	if _, err := unix.IoctlGetTermios(fd, unix.TCGETS); err != nil {
		return
	}
	_ = unix.IoctlSetInt(fd, unix.TCFLSH, unix.TCIFLUSH)
}
