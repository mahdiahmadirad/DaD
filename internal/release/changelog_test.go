package release

import "testing"

func TestChangelogNotes(t *testing.T) {
	t.Parallel()
	content := `# Changelog

## [Unreleased]

### Added

- Pending work.

## [1.2.0] - 2026-07-27

### Added

- A useful feature.

## [1.1.0] - 2026-07-20

### Fixed

- An older fix.
`
	got, err := ChangelogNotes(content, "v1.2.0")
	if err != nil {
		t.Fatal(err)
	}
	want := "### Added\n\n- A useful feature.\n"
	if got != want {
		t.Fatalf("notes = %q, want %q", got, want)
	}
}

func TestChangelogNotesRejectsMissingOrEmptyEntry(t *testing.T) {
	t.Parallel()
	content := "# Changelog\n\n## [1.2.0] - 2026-07-27\n\n### Added\n"
	if _, err := ChangelogNotes(content, "1.2.0"); err == nil {
		t.Fatal("expected empty entry error")
	}
	if _, err := ChangelogNotes(content, "1.3.0"); err == nil {
		t.Fatal("expected missing entry error")
	}
}
