package dad

import "testing"

func TestContextResolvesExplicitAuthorityChain(t *testing.T) {
	root := fixtureRepository(t)
	writeFixture(t, root, "docs/adr/ADR-0001.md", validADR("ADR-0001", "Accepted", ""))
	writeFixture(
		t, root, "docs/specs/SPEC-0001.md",
		validSPEC(
			"SPEC-0001", "Approved",
			"[ADR-0001](../adr/ADR-0001.md) constrains the behavior.",
		),
	)
	writeFixture(
		t, root, "docs/tasks/TASK-0000.md",
		validTASK(
			"TASK-0000", "Complete", "",
			"- Bootstrap acceptance condition is satisfied.",
			"Bootstrap verification passed and the fixture result was recorded.",
		),
	)
	writeFixture(
		t, root, "docs/tasks/TASK-0001.md",
		validTASK(
			"TASK-0001", "Ready",
			"- [TASK-0000](TASK-0000.md) is the completed prerequisite.\n"+
				"- [SPEC-0001](../specs/SPEC-0001.md) governs behavior.",
			"- The approved behavior is implemented and verified.",
			"Not yet recorded.",
		),
	)
	result, diagnostics, code := Context(root, "task-0001")
	if code != ExitOK || len(diagnostics) != 0 {
		t.Fatalf("result=%#v diagnostics=%#v code=%d", result, diagnostics, code)
	}
	var roles []string
	for _, document := range result.Documents {
		roles = append(roles, document.Role)
	}
	want := []string{
		"project-vision", "working-agreement", "task", "referenced-task",
		"architecture-decision", "specification",
	}
	if len(roles) != len(want) {
		t.Fatalf("roles=%v, want %v", roles, want)
	}
	for index := range want {
		if roles[index] != want[index] {
			t.Fatalf("roles=%v, want %v", roles, want)
		}
	}
}

func TestContextRejectsDraftSpecification(t *testing.T) {
	root := fixtureRepository(t)
	writeFixture(
		t, root, "docs/specs/SPEC-0001.md",
		validSPEC("SPEC-0001", "Draft", ""),
	)
	writeFixture(
		t, root, "docs/tasks/TASK-0001.md",
		validTASK(
			"TASK-0001", "Ready",
			"[SPEC-0001](../specs/SPEC-0001.md) governs behavior.",
			"- The approved behavior is implemented.", "Not yet recorded.",
		),
	)
	_, diagnostics, code := Context(root, "TASK-0001")
	if code != ExitCheck || !hasDiagnostic(diagnostics, "DAD-CONTEXT-004") {
		t.Fatalf("diagnostics=%#v code=%d", diagnostics, code)
	}
}
