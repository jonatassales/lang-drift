package files

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// Lock takes an advisory, nonblocking, cross-process lock on path. The lock
// file stays after release; the kernel releases the lock if the process dies.
// Every writer of the protected state must lock the same path.
func Lock(path string) (unlock func(), err error) {
	directory := filepath.Dir(path)

	if err := RequireNoSymlinks(directory); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(directory, privateDirectoryMode); err != nil {
		return nil, err
	}

	descriptor, err := syscall.Open(path, syscall.O_CREAT|syscall.O_RDWR|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, privateFileMode)
	if err != nil {
		return nil, fmt.Errorf("open lock: %w", err)
	}

	file := os.NewFile(uintptr(descriptor), path)

	if err := syscall.Flock(descriptor, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = file.Close()

		return nil, fmt.Errorf("another operation holds %s: %w", filepath.Base(path), err)
	}

	released := false

	return func() {
		if released {
			return
		}

		released = true
		_ = syscall.Flock(descriptor, syscall.LOCK_UN)
		_ = file.Close()
	}, nil
}
