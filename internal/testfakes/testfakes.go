// Package testfakes writes the fake command-line programs that tests place on PATH.
//
// It exists for one platform quirk: macOS charges a few hundred milliseconds of policy checks on
// the first execve of a script whose inode changed, and tests that write a fake and probe it under
// a sub-second budget then time out before the fake even starts. The check is keyed on the file
// as written, so the warm-up must execute the final content: WriteExecutable inserts a guard
// after the shebang that exits immediately for one reserved argument, writes the script, and
// runs it once with that argument. The package is only imported from _test files.
package testfakes

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

// WarmupArg is the reserved first argument that makes a script written by WriteExecutable exit
// before doing anything. Fakes must not use it themselves.
const WarmupArg = "--testfakes-warmup"

// WriteExecutable writes a shell script to path and runs it once so that its first real
// invocation carries no first-execution latency. On Windows it only writes the file.
func WriteExecutable(t testing.TB, path string, content []byte) {
	t.Helper()
	if runtime.GOOS == "windows" {
		if err := os.WriteFile(path, content, 0o755); err != nil {
			t.Fatal(err)
		}
		return
	}
	guard := []byte("if [ \"${1:-}\" = \"" + WarmupArg + "\" ]; then exit 0; fi\n")
	var script []byte
	if bytes.HasPrefix(content, []byte("#!")) {
		newline := bytes.IndexByte(content, '\n')
		if newline < 0 {
			newline = len(content) - 1
		}
		script = append(append(append([]byte{}, content[:newline+1]...), guard...), content[newline+1:]...)
	} else {
		script = append(append([]byte("#!/bin/sh\n"), guard...), content...)
	}
	if err := os.WriteFile(path, script, 0o755); err != nil {
		t.Fatal(err)
	}
	// Best effort: a fake written to fail at launch (for example with a missing interpreter)
	// fails here too, and that is the behavior the test wants to observe.
	_ = exec.Command(path, WarmupArg).Run()
}
