package platform

import (
	"io"
	"os"
)

// FlushPendingInput discards any unread input queued on r's terminal. It is a
// silent no-op when r is not an *os.File, the file is not a terminal, or the
// flush fails: it never reports an error, so a prompt can call it unconditionally.
func FlushPendingInput(r io.Reader) {
	f, ok := r.(*os.File)
	if !ok {
		return
	}
	// Not f.Fd(): that would switch the file to blocking mode as a side effect.
	sc, err := f.SyscallConn()
	if err != nil {
		return
	}
	_ = sc.Control(func(fd uintptr) { flushTTYInput(int(fd)) })
}
