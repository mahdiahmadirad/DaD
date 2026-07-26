//go:build windows

package dad

import (
	"path/filepath"
	"testing"
)

func TestWindowsRootSupportsNativeUnicodePath(t *testing.T) {
	root := fixtureRepository(t)
	nested := filepath.Join(root, "space and ünicode")
	writeFixture(t, nested, "PROJECT-VISION.md", "# Vision\n")
	writeFixture(t, nested, "AGENTS.md", "# Agreement\n")
	resolved, err := ResolveRoot(nested, nil, root)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(resolved) != filepath.Clean(nested) {
		t.Fatalf("resolved %q, want %q", resolved, nested)
	}
}
