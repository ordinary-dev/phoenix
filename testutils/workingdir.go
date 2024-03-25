package testutils

import (
	"os"
	"path"
	"runtime"
)

// Change the current directory to the project directory.
// Useful for tests that work with files.
func ResetWorkingDir() error {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	return os.Chdir(dir)
}
