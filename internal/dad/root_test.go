package dad

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveRootPrecedenceAndDiscovery(t *testing.T) {
	parent := fixtureRepository(t)
	nestedRoot := filepath.Join(parent, "nested repo ü")
	if err := os.MkdirAll(nestedRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, nestedRoot, "PROJECT-VISION.md", "# Nested\n")
	writeFixture(t, nestedRoot, "AGENTS.md", "# Nested\n")
	deep := filepath.Join(nestedRoot, "a", "b")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}

	root, err := ResolveRoot("", nil, deep)
	if err != nil {
		t.Fatal(err)
	}
	expectedNested, _ := filepath.EvalSymlinks(nestedRoot)
	if root != expectedNested {
		t.Fatalf("discovered %q, want %q", root, expectedNested)
	}

	root, err = ResolveRoot(parent, map[string]string{"DAD_ROOT": nestedRoot}, deep)
	if err != nil {
		t.Fatal(err)
	}
	expectedParent, _ := filepath.EvalSymlinks(parent)
	if root != expectedParent {
		t.Fatalf("explicit root %q, want %q", root, expectedParent)
	}

	root, err = ResolveRoot("", map[string]string{"DAD_ROOT": parent}, deep)
	if err != nil {
		t.Fatal(err)
	}
	if root != expectedParent {
		t.Fatalf("environment root %q, want %q", root, expectedParent)
	}
}

func TestResolveRootDoesNotRequireGit(t *testing.T) {
	root := fixtureRepository(t)
	if _, err := os.Stat(filepath.Join(root, ".git")); !os.IsNotExist(err) {
		t.Fatalf("fixture unexpectedly has .git: %v", err)
	}
	resolved, err := ResolveRoot(root, nil, root)
	if err != nil {
		t.Fatal(err)
	}
	expected, _ := filepath.EvalSymlinks(root)
	if resolved != expected {
		t.Fatalf("resolved %q, want %q", resolved, expected)
	}
}

func TestResolveRootRejectsMissingMarkers(t *testing.T) {
	root := t.TempDir()
	if _, err := ResolveRoot(root, nil, root); err == nil {
		t.Fatal("expected root error")
	}
}

func TestResolveReferenceRejectsEscape(t *testing.T) {
	root := fixtureRepository(t)
	source := filepath.Join(root, "docs", "tasks", "TASK-0001.md")
	if _, err := ResolveReference(root, source, "../../../outside.md"); err == nil {
		t.Fatal("expected escape error")
	}
}
