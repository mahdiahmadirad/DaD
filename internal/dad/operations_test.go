package dad

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"
)

func TestInitDryRunApplyAndIdempotence(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "PROJECT-VISION.md", "# Vision\n")
	writeFixture(t, root, "AGENTS.md", "# Agreement\n")
	files := map[string][]byte{
		"docs/DOCUMENTATION.md":           []byte("# Documentation\n"),
		"docs/templates/ADR-TEMPLATE.md":  []byte(adrTemplate),
		"docs/templates/SPEC-TEMPLATE.md": []byte(specTemplate),
		"docs/templates/TASK-TEMPLATE.md": []byte(taskTemplate),
	}
	preview, diagnostics, code := Init(root, files, true)
	if code != ExitOK || len(diagnostics) != 0 || len(preview.Created) != 4 {
		t.Fatalf("preview=%#v diagnostics=%#v code=%d", preview, diagnostics, code)
	}
	if _, err := os.Stat(filepath.Join(root, "docs")); !os.IsNotExist(err) {
		t.Fatalf("dry run created docs: %v", err)
	}

	applied, diagnostics, code := Init(root, files, false)
	if code != ExitOK || len(diagnostics) != 0 || len(applied.Created) != 4 {
		t.Fatalf("applied=%#v diagnostics=%#v code=%d", applied, diagnostics, code)
	}
	again, diagnostics, code := Init(root, files, false)
	if code != ExitOK || len(diagnostics) != 0 ||
		len(again.Created) != 0 || len(again.Unchanged) != 4 {
		t.Fatalf("again=%#v diagnostics=%#v code=%d", again, diagnostics, code)
	}
}

func TestInitRefusesConflictsBeforeWriting(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "PROJECT-VISION.md", "# Vision\n")
	writeFixture(t, root, "AGENTS.md", "# Agreement\n")
	writeFixture(t, root, "docs/DOCUMENTATION.md", "different\n")
	files := map[string][]byte{
		"docs/DOCUMENTATION.md":          []byte("# Documentation\n"),
		"docs/templates/ADR-TEMPLATE.md": []byte(adrTemplate),
	}
	result, diagnostics, code := Init(root, files, false)
	if code != ExitConflict || len(result.Conflicts) != 1 ||
		!hasDiagnostic(diagnostics, "DAD-INIT-001") {
		t.Fatalf("result=%#v diagnostics=%#v code=%d", result, diagnostics, code)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "templates")); !os.IsNotExist(err) {
		t.Fatalf("conflicting init wrote files: %v", err)
	}
}

func TestAtomicWriteNewNeverOverwrites(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "document.md")
	if err := AtomicWriteNew(path, []byte("first\n")); err != nil {
		t.Fatal(err)
	}
	if err := AtomicWriteNew(path, []byte("second\n")); !errors.Is(err, os.ErrExist) {
		t.Fatalf("got %v, want os.ErrExist", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "first\n" {
		t.Fatalf("content overwritten: %q", content)
	}
}

func TestNewDocumentDryRunAndExplicitNumber(t *testing.T) {
	root := fixtureRepository(t)
	preview, diagnostics, code := NewDocument(root, NewOptions{
		Type: TASK, Title: "Preview", DryRun: true,
	})
	if code != ExitOK || len(diagnostics) != 0 ||
		preview.ID != "TASK-0001" || preview.Created {
		t.Fatalf("preview=%#v diagnostics=%#v code=%d", preview, diagnostics, code)
	}
	if _, err := os.Stat(filepath.Join(root, preview.Path)); !os.IsNotExist(err) {
		t.Fatalf("dry run wrote document: %v", err)
	}

	result, diagnostics, code := NewDocument(root, NewOptions{
		Type: SPEC, Title: "Contract", Number: 7,
	})
	if code != ExitOK || len(diagnostics) != 0 ||
		result.ID != "SPEC-0007" || result.Status != "Draft" || !result.Created {
		t.Fatalf("result=%#v diagnostics=%#v code=%d", result, diagnostics, code)
	}
	content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(result.Path)))
	if err != nil {
		t.Fatal(err)
	}
	if string(content[:len("# SPEC-0007: Contract")]) != "# SPEC-0007: Contract" {
		t.Fatalf("unexpected content: %q", content)
	}
}

func TestNewDocumentRefusesMalformedRepository(t *testing.T) {
	root := fixtureRepository(t)
	writeFixture(
		t, root, "docs/specs/SPEC-0001.md",
		validSPEC("SPEC-0001", "Ready", ""),
	)
	_, diagnostics, code := NewDocument(root, NewOptions{
		Type: SPEC, Title: "Must not be created",
	})
	if code != ExitConflict || !hasDiagnostic(diagnostics, "DAD-DOC-003") {
		t.Fatalf("diagnostics=%#v code=%d", diagnostics, code)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "specs", "SPEC-0002.md")); !os.IsNotExist(err) {
		t.Fatalf("malformed repository was mutated: %v", err)
	}
}

func TestNewDocumentRefusesDirectoryOutsideRoot(t *testing.T) {
	root := fixtureRepository(t)
	outside := t.TempDir()
	link := filepath.Join(root, "docs", "adr")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	_, diagnostics, code := NewDocument(root, NewOptions{
		Type: ADR, Title: "Must remain inside",
	})
	if code != ExitConflict || !hasDiagnostic(diagnostics, "DAD-PATH-002") {
		t.Fatalf("diagnostics=%#v code=%d", diagnostics, code)
	}
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("outside directory was modified: %#v", entries)
	}
}

func TestConcurrentNewDocumentAllocatesUniqueNumbers(t *testing.T) {
	root := fixtureRepository(t)
	const count = 8
	var wait sync.WaitGroup
	results := make(chan NewResult, count)
	errorsChannel := make(chan []Diagnostic, count)
	for index := 0; index < count; index++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			result, diagnostics, code := NewDocument(root, NewOptions{
				Type: TASK, Title: "Concurrent",
			})
			if code != ExitOK {
				errorsChannel <- diagnostics
				return
			}
			results <- result
		}()
	}
	wait.Wait()
	close(results)
	close(errorsChannel)
	for diagnostics := range errorsChannel {
		t.Errorf("creation failed: %#v", diagnostics)
	}
	var identifiers []string
	for result := range results {
		identifiers = append(identifiers, result.ID)
	}
	sort.Strings(identifiers)
	if len(identifiers) != count {
		t.Fatalf("created %d documents, want %d", len(identifiers), count)
	}
	for index, identifier := range identifiers {
		want := "TASK-000" + string(rune('1'+index))
		if identifier != want {
			t.Fatalf("identifiers=%v; at %d got %s want %s", identifiers, index, identifier, want)
		}
	}
}
