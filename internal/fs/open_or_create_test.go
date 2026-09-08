//go:build unix

package fs

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

// Concurrent first-time creators of one lock file must all succeed. macOS answers a plain
// O_CREAT open with ENOENT for most racers, which OpenLockFile must hide.
func TestOpenLockFileConcurrentCreation(t *testing.T) {
	root := t.TempDir()
	for trial := range 100 {
		path := filepath.Join(root, fmt.Sprintf("race-%d.lock", trial))
		var wg sync.WaitGroup
		errs := make(chan error, 4)
		for range 4 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				f, err := OpenLockFile(root, path)
				if err != nil {
					errs <- err
					return
				}
				errs <- f.Close()
			}()
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			if err != nil {
				t.Fatalf("trial %d: %v", trial, err)
			}
		}
	}
}

func TestOpenAppendFileConcurrentCreation(t *testing.T) {
	root := t.TempDir()
	for trial := range 100 {
		path := filepath.Join(root, fmt.Sprintf("race-%d.log", trial))
		var wg sync.WaitGroup
		errs := make(chan error, 4)
		for range 4 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				f, err := OpenAppendFile(root, path)
				if err != nil {
					errs <- err
					return
				}
				errs <- f.Close()
			}()
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			if err != nil {
				t.Fatalf("trial %d: %v", trial, err)
			}
		}
	}
}
