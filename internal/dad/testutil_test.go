package dad

import (
	"os"
	"path/filepath"
	"testing"
)

const adrTemplate = `# ADR-NNNN: Title

## Status

Proposed

## Context

Context placeholder.

## Decision

Decision placeholder.

## Consequences

Consequences placeholder.

## Alternatives considered

Alternatives placeholder.

## References

References placeholder.
`

const specTemplate = `# SPEC-NNNN: Title

## Status

Draft

## Purpose

Purpose placeholder.

## Scope

Scope placeholder.

## Out of scope

Out of scope placeholder.

## Requirements

Requirements placeholder.

## Interfaces and behavior

Interfaces placeholder.

## Acceptance criteria

Acceptance placeholder.

## References

References placeholder.
`

const taskTemplate = `# TASK-NNNN: Title

## Status

Proposed

## Outcome

Outcome placeholder.

## In scope

In scope placeholder.

## Out of scope

Out of scope placeholder.

## Constraints

Constraints placeholder.

## References

References placeholder.

## Acceptance conditions

Acceptance placeholder.

## Completion evidence

Completion placeholder.
`

func fixtureRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeFixture(t, root, "PROJECT-VISION.md", "# Project Vision\n\nFixture.\n")
	writeFixture(t, root, "AGENTS.md", "# Working Agreement\n\nFixture.\n")
	writeFixture(t, root, "docs/DOCUMENTATION.md", "# Documentation\n")
	writeFixture(t, root, "docs/templates/ADR-TEMPLATE.md", adrTemplate)
	writeFixture(t, root, "docs/templates/SPEC-TEMPLATE.md", specTemplate)
	writeFixture(t, root, "docs/templates/TASK-TEMPLATE.md", taskTemplate)
	return root
}

func writeFixture(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func validADR(id, status, references string) string {
	return "# " + id + `: Decision

## Status

` + status + `

## Context

The fixture needs one consequential decision with enough context for review.

## Decision

Use the fixture decision in every applicable implementation.

## Consequences

The fixture remains deterministic and the decision remains visible.

## Alternatives considered

The alternative was rejected because it did not preserve the fixture.

## References

` + references + "\n"
}

func validSPEC(id, status, references string) string {
	return "# " + id + `: Behavior

## Status

` + status + `

## Purpose

Define stable fixture behavior for command and repository tests.

## Scope

The fixture document and its observable result are in scope.

## Out of scope

Unrelated production behavior and external services are out of scope.

## Requirements

The fixture must be deterministic and locally verifiable.

## Interfaces and behavior

The fixture accepts one input and returns one stable result.

## Acceptance criteria

Focused tests demonstrate the stable result and failure behavior.

## References

` + references + "\n"
}

func validTASK(id, status, references, acceptance, completion string) string {
	return "# " + id + `: Work

## Status

` + status + `

## Outcome

Deliver one bounded and verifiable fixture outcome.

## In scope

- Implement only the fixture outcome.

## Out of scope

- Do not change unrelated fixture behavior.

## Constraints

- Preserve deterministic behavior and existing fixture files.

## References

` + references + `

## Acceptance conditions

` + acceptance + `

## Completion evidence

` + completion + "\n"
}

func hasDiagnostic(diagnostics []Diagnostic, code string) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			return true
		}
	}
	return false
}
