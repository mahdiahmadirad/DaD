package dad

import (
	"reflect"
	"testing"
)

func TestCheckValidRepository(t *testing.T) {
	root := fixtureRepository(t)
	writeFixture(t, root, "docs/adr/ADR-0001.md", validADR("ADR-0001", "Accepted", ""))
	writeFixture(
		t, root, "docs/specs/SPEC-0001.md",
		validSPEC("SPEC-0001", "Approved", "[ADR-0001](../adr/ADR-0001.md) governs."),
	)
	writeFixture(
		t, root, "docs/tasks/TASK-0001.md",
		validTASK(
			"TASK-0001", "Ready",
			"[SPEC-0001](../specs/SPEC-0001.md) governs.",
			"- All fixture checks pass with the approved behavior.", "Not yet recorded.",
		),
	)
	result, diagnostics, code := Check(root, nil, false)
	if code != ExitOK || result.Summary.Errors != 0 || len(diagnostics) != 0 {
		t.Fatalf("result=%#v diagnostics=%#v code=%d", result, diagnostics, code)
	}
}

func TestCheckFindings(t *testing.T) {
	longDuplicate := "This exact normative fixture paragraph is intentionally long enough " +
		"to trigger deterministic duplicate detection across two governed documents " +
		"without relying on semantic comparison or probabilistic analysis."
	tests := []struct {
		name string
		code string
		set  func(*testing.T, string)
	}{
		{
			name: "duplicate identifier", code: "DAD-DOC-002",
			set: func(t *testing.T, root string) {
				writeFixture(t, root, "docs/adr/ADR-0001.md", validADR("ADR-0001", "Accepted", ""))
				writeFixture(t, root, "docs/adr/ADR-0002.md", validADR("ADR-0001", "Accepted", ""))
			},
		},
		{
			name: "noncanonical filename", code: "DAD-CHECK-006",
			set: func(t *testing.T, root string) {
				writeFixture(t, root, "docs/adr/ADR-0001-title.md", validADR("ADR-0001", "Accepted", ""))
			},
		},
		{
			name: "invalid status", code: "DAD-DOC-003",
			set: func(t *testing.T, root string) {
				writeFixture(t, root, "docs/specs/SPEC-0001.md", validSPEC("SPEC-0001", "Ready", ""))
			},
		},
		{
			name: "broken reference", code: "DAD-CHECK-016",
			set: func(t *testing.T, root string) {
				writeFixture(
					t, root, "docs/tasks/TASK-0001.md",
					validTASK(
						"TASK-0001", "Proposed",
						"[SPEC-0001](../specs/SPEC-0001.md) governs.",
						"- Acceptance remains proposed.", "Not yet recorded.",
					),
				)
			},
		},
		{
			name: "non-authoritative governing reference", code: "DAD-CHECK-013",
			set: func(t *testing.T, root string) {
				writeFixture(t, root, "docs/specs/SPEC-0001.md", validSPEC("SPEC-0001", "Draft", ""))
				writeFixture(
					t, root, "docs/tasks/TASK-0001.md",
					validTASK(
						"TASK-0001", "Ready",
						"[SPEC-0001](../specs/SPEC-0001.md) governs.",
						"- Acceptance is substantive and locally verifiable.", "Not yet recorded.",
					),
				)
			},
		},
		{
			name: "superseded without replacement", code: "DAD-CHECK-014",
			set: func(t *testing.T, root string) {
				writeFixture(t, root, "docs/adr/ADR-0001.md", validADR("ADR-0001", "Superseded", ""))
			},
		},
		{
			name: "complete without evidence", code: "DAD-CHECK-011",
			set: func(t *testing.T, root string) {
				writeFixture(
					t, root, "docs/tasks/TASK-0001.md",
					validTASK(
						"TASK-0001", "Complete", "",
						"- Acceptance is satisfied.", "To be recorded after implementation.",
					),
				)
			},
		},
		{
			name: "reference escapes root", code: "DAD-CHECK-015",
			set: func(t *testing.T, root string) {
				writeFixture(
					t, root, "docs/tasks/TASK-0001.md",
					validTASK(
						"TASK-0001", "Proposed", "[outside](../../../outside.md)",
						"- Acceptance remains proposed.", "Not yet recorded.",
					),
				)
			},
		},
		{
			name: "ready without acceptance", code: "DAD-CHECK-010",
			set: func(t *testing.T, root string) {
				writeFixture(t, root, "docs/tasks/TASK-0001.md", validTASK("TASK-0001", "Ready", "", "", ""))
			},
		},
		{
			name: "blocked without note", code: "DAD-CHECK-012",
			set: func(t *testing.T, root string) {
				writeFixture(
					t, root, "docs/tasks/TASK-0001.md",
					validTASK(
						"TASK-0001", "Blocked", "",
						"- Acceptance remains defined while blocked.", "Not yet recorded.",
					),
				)
			},
		},
		{
			name: "malformed template", code: "DAD-CHECK-017",
			set: func(t *testing.T, root string) {
				writeFixture(t, root, "docs/templates/SPEC-TEMPLATE.md", "# SPEC-NNNN: Title\n")
			},
		},
		{
			name: "unresolved template identifier", code: "DAD-DOC-001",
			set: func(t *testing.T, root string) {
				writeFixture(t, root, "docs/specs/SPEC-0001.md", specTemplate)
			},
		},
		{
			name: "exact duplicate", code: "DAD-CHECK-018",
			set: func(t *testing.T, root string) {
				first := validADR("ADR-0001", "Accepted", "") + "\n" + longDuplicate + "\n"
				second := validSPEC("SPEC-0001", "Approved", "") + "\n" + longDuplicate + "\n"
				writeFixture(t, root, "docs/adr/ADR-0001.md", first)
				writeFixture(t, root, "docs/specs/SPEC-0001.md", second)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := fixtureRepository(t)
			test.set(t, root)
			_, diagnostics, _ := Check(root, nil, false)
			if !hasDiagnostic(diagnostics, test.code) {
				t.Fatalf("missing %s in %#v", test.code, diagnostics)
			}
		})
	}
}

func TestCheckStrictMakesWarningFail(t *testing.T) {
	root := fixtureRepository(t)
	writeFixture(
		t, root, "docs/tasks/TASK-0001.md",
		validTASK(
			"TASK-0001", "Blocked", "",
			"- Acceptance remains defined while blocked.", "Not yet recorded.",
		),
	)
	_, _, code := Check(root, nil, false)
	if code != ExitOK {
		t.Fatalf("warning-only check exited %d", code)
	}
	_, _, code = Check(root, nil, true)
	if code != ExitCheck {
		t.Fatalf("strict warning check exited %d", code)
	}
}

func TestCheckSelectedDocumentIncludesReferencedGovernedDocuments(t *testing.T) {
	root := fixtureRepository(t)
	writeFixture(t, root, "docs/adr/ADR-0001.md", validADR("ADR-0001", "Accepted", ""))
	writeFixture(
		t, root, "docs/specs/SPEC-0001.md",
		validSPEC("SPEC-0001", "Approved", "[ADR-0001](../adr/ADR-0001.md) governs."),
	)
	writeFixture(
		t, root, "docs/tasks/TASK-0001.md",
		validTASK(
			"TASK-0001", "In Progress",
			"[SPEC-0001](../specs/SPEC-0001.md) governs.",
			"- The selected document scope is checked.", "Not yet recorded.",
		),
	)

	result, diagnostics, code := Check(root, []string{"TASK-0001"}, false)
	if code != ExitOK || len(diagnostics) != 0 {
		t.Fatalf("result=%#v diagnostics=%#v code=%d", result, diagnostics, code)
	}
	want := []string{"ADR-0001", "SPEC-0001", "TASK-0001"}
	if !reflect.DeepEqual(result.CheckedDocuments, want) {
		t.Fatalf("CheckedDocuments = %v, want %v", result.CheckedDocuments, want)
	}
}
