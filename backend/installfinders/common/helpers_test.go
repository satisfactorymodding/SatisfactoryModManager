package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOsPathEqualSameFile(t *testing.T) {
	dir := t.TempDir()

	// Same path is trivially equal.
	if !OsPathEqual(dir, dir) {
		t.Error("identical paths should be equal")
	}

	// Same directory reached through a symlink must be equal.
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(dir, link); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}
	if !OsPathEqual(dir, link) {
		t.Error("symlinked paths to the same directory should be equal")
	}

	// Different directories must not be equal.
	other := t.TempDir()
	if OsPathEqual(dir, other) {
		t.Error("different directories should not be equal")
	}
}
