package dad

import "path/filepath"

type StatusResult struct {
	Root            string                    `json:"root"`
	Counts          map[string]map[string]int `json:"counts"`
	ActionableTasks []Document                `json:"actionable_tasks"`
	BlockedTasks    []Document                `json:"blocked_tasks"`
	CheckSummary    CheckSummary              `json:"check_summary"`
	CoreFiles       map[string]bool           `json:"core_files"`
}

func Status(root string) (StatusResult, []Diagnostic, int) {
	repository := LoadRepository(root)
	checkResult, diagnostics, exit := Check(root, nil, false)
	result := StatusResult{
		Root:         root,
		Counts:       make(map[string]map[string]int),
		CheckSummary: checkResult.Summary,
		CoreFiles: map[string]bool{
			"PROJECT-VISION.md": nonEmptyFile(filepath.Join(root, "PROJECT-VISION.md")),
			"AGENTS.md":         nonEmptyFile(filepath.Join(root, "AGENTS.md")),
			"docs/DOCUMENTATION.md": nonEmptyFile(
				filepath.Join(root, "docs", "DOCUMENTATION.md"),
			),
			"docs/templates/ADR-TEMPLATE.md": nonEmptyFile(
				filepath.Join(root, "docs", "templates", "ADR-TEMPLATE.md"),
			),
			"docs/templates/SPEC-TEMPLATE.md": nonEmptyFile(
				filepath.Join(root, "docs", "templates", "SPEC-TEMPLATE.md"),
			),
			"docs/templates/TASK-TEMPLATE.md": nonEmptyFile(
				filepath.Join(root, "docs", "templates", "TASK-TEMPLATE.md"),
			),
		},
	}
	for _, document := range repository.Documents {
		if !StatusAllowed(document.Type, document.Status) ||
			len(repository.ByID[document.ID]) != 1 {
			continue
		}
		typeName := string(document.Type)
		if result.Counts[typeName] == nil {
			result.Counts[typeName] = make(map[string]int)
		}
		result.Counts[typeName][document.Status]++
		if document.Type != TASK {
			continue
		}
		switch document.Status {
		case "Ready", "In Progress":
			result.ActionableTasks = append(result.ActionableTasks, document)
		case "Blocked":
			result.BlockedTasks = append(result.BlockedTasks, document)
		}
	}
	SortDocuments(result.ActionableTasks)
	SortDocuments(result.BlockedTasks)
	return result, diagnostics, exit
}
