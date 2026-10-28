package files_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jonatassales/tao/exphere/cli/src/infrastructure/files"
)

func TestLock(t *testing.T) {
	t.Parallel()

	path := filepath.Join(tempDir(t), "state", "cluster.lock")

	unlock, err := files.Lock(path)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := files.Lock(path); err == nil || !strings.Contains(err.Error(), "another operation holds cluster.lock") {
		t.Fatalf("second Lock() error = %v, want the held lock reported", err)
	}

	unlock()
	unlock()

	again, err := files.Lock(path)
	if err != nil {
		t.Fatalf("Lock() after release error = %v", err)
	}
	defer again()

	assertMode(t, path, 0o600)
	assertMode(t, filepath.Dir(path), 0o700)
}

func TestLock_RefusesLinks(t *testing.T) {
	t.Parallel()

	root := tempDir(t)
	target := writeFile(t, filepath.Join(root, "target"), "")

	if _, err := files.Lock(symlink(t, target, filepath.Join(root, "cluster.lock"))); err == nil {
		t.Error("Lock followed a symlinked lock file")
	}

	linked := symlink(t, mkdir(t, filepath.Join(root, "real")), filepath.Join(root, "linked"))
	if _, err := files.Lock(filepath.Join(linked, "cluster.lock")); err == nil {
		t.Error("Lock used a symlinked directory")
	}
}

// BenchmarkLock measures taking and releasing the project lock, and refusing
// it while another holder has it. The lock never waits, so contention costs
// one refused attempt.
func BenchmarkLock(b *testing.B) {
	path := filepath.Join(tempDir(b), "state", "cluster.lock")

	b.Run("uncontended", func(b *testing.B) {
		for range b.N {
			unlock, err := files.Lock(path)
			if err != nil {
				b.Fatal(err)
			}

			unlock()
		}
	})

	b.Run("contended", func(b *testing.B) {
		unlock, err := files.Lock(path)
		if err != nil {
			b.Fatal(err)
		}
		defer unlock()

		for range b.N {
			if _, err := files.Lock(path); err == nil {
				b.Fatal("a held lock was granted twice")
			}
		}
	})
}
