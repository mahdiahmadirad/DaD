package dad

import (
	"bytes"
	"os"
	"testing"
)

func TestEmbeddedSupportFilesMatchAuthoritativeSources(t *testing.T) {
	files, err := SupportFiles()
	if err != nil {
		t.Fatal(err)
	}
	for path, embedded := range files {
		source, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if !bytes.Equal(embedded, source) {
			t.Errorf("embedded %s differs from source", path)
		}
	}
}

func TestOfficialPromptsAreOrderedAndReadable(t *testing.T) {
	prompts := Prompts()
	if len(prompts) != 5 {
		t.Fatalf("got %d prompts, want 5", len(prompts))
	}
	wantFirst := "project-bootstrap"
	wantLast := "documentation-reconciliation"
	if prompts[0].Name != wantFirst || prompts[len(prompts)-1].Name != wantLast {
		t.Fatalf("unexpected prompt order: %#v", prompts)
	}
	for _, prompt := range prompts {
		info, content, found, err := Prompt(prompt.Name)
		if err != nil {
			t.Fatal(err)
		}
		if !found || info.Name != prompt.Name || len(content) == 0 {
			t.Errorf("prompt %q is not readable", prompt.Name)
		}
	}
	if _, _, found, err := Prompt("unknown"); err != nil || found {
		t.Fatalf("unknown prompt: found=%v err=%v", found, err)
	}
}
