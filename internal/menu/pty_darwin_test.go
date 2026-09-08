//go:build darwin

package menu

import (
	"os"
	"testing"
	"unsafe"

	"golang.org/x/sys/unix"
)

// openPTY opens a pseudo terminal pair through the macOS ptmx ioctls, which differ from the
// Linux TIOCSPTLCK/TIOCGPTN sequence in pty_linux_test.go.
func openPTY(t *testing.T) (master, slave *os.File) {
	t.Helper()
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	fd := int(master.Fd())
	for _, req := range []uint{unix.TIOCPTYGRANT, unix.TIOCPTYUNLK} {
		if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(req), 0); errno != 0 {
			master.Close()
			t.Fatal(errno)
		}
	}
	name, err := ptsName(master)
	if err != nil {
		master.Close()
		t.Fatal(err)
	}
	slave, err = os.OpenFile(name, os.O_RDWR, 0)
	if err != nil {
		master.Close()
		t.Fatal(err)
	}
	ws := unix.Winsize{Row: 40, Col: 120}
	_ = unix.IoctlSetWinsize(int(slave.Fd()), unix.TIOCSWINSZ, &ws)
	return master, slave
}

func ptsName(master *os.File) (string, error) {
	var buf [128]byte
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, master.Fd(), uintptr(unix.TIOCPTYGNAME), uintptr(unsafe.Pointer(&buf[0]))); errno != 0 {
		return "", errno
	}
	n := 0
	for n < len(buf) && buf[n] != 0 {
		n++
	}
	return string(buf[:n]), nil
}
