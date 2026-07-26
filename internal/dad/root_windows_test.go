//go:build windows

package dad

import (
	"os"
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
	resolvedInfo, err := os.Stat(resolved)
	if err != nil {
		t.Fatal(err)
	}
	nestedInfo, err := os.Stat(nested)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(resolvedInfo, nestedInfo) {
		t.Fatalf("resolved %q does not identify %q", resolved, nested)
	}
	if filepath.Base(resolved) != "space and ünicode" {
		t.Fatalf("resolved path lost native Unicode name: %q", resolved)
	}
}
