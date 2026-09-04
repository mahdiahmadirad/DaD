package dad

import "testing"

func TestRepositoryDiscoversTasksOnlyFromDocsTasks(t *testing.T) {
	root := fixtureRepository(t)
	writeFixture(
		t, root, "TASK-0001.md",
		validTASK(
			"TASK-0001", "Ready", "",
			"- Legacy root task should not be discovered.", "Not yet recorded.",
		),
	)
	writeFixture(
		t, root, "docs/tasks/TASK-0002.md",
		validTASK(
			"TASK-0002", "Ready", "",
			"- Canonical task should be discovered.", "Not yet recorded.",
		),
	)

	repository := LoadRepository(root)
	if len(repository.Diagnostics) != 0 {
		t.Fatalf("diagnostics=%#v", repository.Diagnostics)
	}
	if len(repository.Documents) != 1 {
		t.Fatalf("documents=%#v", repository.Documents)
	}
	if repository.Documents[0].ID != "TASK-0002" ||
		repository.Documents[0].Path != "docs/tasks/TASK-0002.md" {
		t.Fatalf("document=%#v", repository.Documents[0])
	}
	if _, found := repository.ByID["TASK-0001"]; found {
		t.Fatal("root-level TASK was discovered")
	}
}
