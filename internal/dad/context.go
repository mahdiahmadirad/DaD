package dad

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ContextDocument struct {
	ID         string `json:"id,omitempty"`
	Path       string `json:"path"`
	Role       string `json:"role"`
	Status     string `json:"status,omitempty"`
	RequiredBy string `json:"required_by,omitempty"`
}

type ContextResult struct {
	TaskStatus string            `json:"task_status"`
	Documents  []ContextDocument `json:"documents"`
}

func Context(root, taskID string) (ContextResult, []Diagnostic, int) {
	repository := LoadRepository(root)
	if len(repository.Diagnostics) > 0 {
		return ContextResult{}, repository.Diagnostics, ExitCheck
	}
	task, ok := repository.Document(strings.ToUpper(taskID))
	if !ok || task.Type != TASK {
		return ContextResult{}, []Diagnostic{{
			Code: "DAD-CONTEXT-001", Severity: "error",
			Message: fmt.Sprintf("TASK %s does not exist or is duplicated.", taskID),
		}}, ExitCheck
	}
	result := ContextResult{TaskStatus: task.Status}
	result.Documents = append(result.Documents, ContextDocument{
		Path: "PROJECT-VISION.md", Role: "project-vision",
		RequiredBy: task.ID,
	})
	for _, path := range applicableAgreements(root, filepath.Dir(task.AbsPath)) {
		result.Documents = append(result.Documents, ContextDocument{
			Path: RelativePath(root, path), Role: "working-agreement",
			RequiredBy: task.ID,
		})
	}
	result.Documents = append(result.Documents, ContextDocument{
		ID: task.ID, Path: task.Path, Role: "task", Status: task.Status,
	})

	var relatedTasks, adrs, specs []ContextDocument
	seen := map[string]bool{task.ID: true}
	var diagnostics []Diagnostic
	for _, link := range sectionLinks(task, "References") {
		target, found, diagnostic := referencedDocument(repository, task, link)
		if diagnostic != nil {
			diagnostics = append(diagnostics, *diagnostic)
			continue
		}
		if !found || seen[target.ID] {
			continue
		}
		seen[target.ID] = true
		item := ContextDocument{
			ID: target.ID, Path: target.Path, Status: target.Status,
			RequiredBy: task.ID,
		}
		switch target.Type {
		case TASK:
			item.Role = "referenced-task"
			relatedTasks = append(relatedTasks, item)
		case ADR:
			item.Role = "architecture-decision"
			if target.Status != "Accepted" {
				diagnostics = append(diagnostics, nonAuthoritativeDiagnostic(
					task, link, target, "Accepted",
				))
			}
			adrs = append(adrs, item)
		case SPEC:
			item.Role = "specification"
			if target.Status != "Approved" {
				diagnostics = append(diagnostics, nonAuthoritativeDiagnostic(
					task, link, target, "Approved",
				))
			}
			specs = append(specs, item)
		}
	}

	for _, specItem := range specs {
		spec, _ := repository.Document(specItem.ID)
		for _, link := range sectionLinks(spec, "References") {
			target, found, diagnostic := referencedDocument(repository, spec, link)
			if diagnostic != nil {
				diagnostics = append(diagnostics, *diagnostic)
				continue
			}
			if !found || target.Type != ADR || seen[target.ID] {
				continue
			}
			seen[target.ID] = true
			if target.Status != "Accepted" {
				diagnostics = append(diagnostics, nonAuthoritativeDiagnostic(
					spec, link, target, "Accepted",
				))
			}
			adrs = append(adrs, ContextDocument{
				ID: target.ID, Path: target.Path, Role: "architecture-decision",
				Status: target.Status, RequiredBy: spec.ID,
			})
		}
	}

	sortContextDocuments(relatedTasks)
	sortContextDocuments(adrs)
	sortContextDocuments(specs)
	result.Documents = append(result.Documents, relatedTasks...)
	result.Documents = append(result.Documents, adrs...)
	result.Documents = append(result.Documents, specs...)
	SortDiagnostics(diagnostics)
	if len(diagnostics) > 0 {
		return result, diagnostics, ExitCheck
	}
	return result, nil, ExitOK
}

func applicableAgreements(root, taskDirectory string) []string {
	var directories []string
	current := filepath.Clean(taskDirectory)
	for {
		if !WithinRoot(root, current) && current != filepath.Clean(root) {
			break
		}
		directories = append(directories, current)
		if current == filepath.Clean(root) {
			break
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	for left, right := 0, len(directories)-1; left < right; left, right = left+1, right-1 {
		directories[left], directories[right] = directories[right], directories[left]
	}
	var paths []string
	for _, directory := range directories {
		path := filepath.Join(directory, "AGENTS.md")
		if info, err := os.Stat(path); err == nil && info.Mode().IsRegular() {
			paths = append(paths, path)
		}
	}
	return paths
}

func referencedDocument(
	repository Repository,
	source Document,
	link MarkdownLink,
) (Document, bool, *Diagnostic) {
	if externalLink(link.Target) || strings.HasPrefix(link.Target, "#") {
		return Document{}, false, nil
	}
	targetPath, err := ResolveReference(repository.Root, source.AbsPath, link.Target)
	if err != nil {
		diagnostic := Diagnostic{
			Code: "DAD-CONTEXT-002", Severity: "error",
			Message: err.Error(), Path: source.Path,
			Line: link.Line, Column: link.Column,
		}
		return Document{}, false, &diagnostic
	}
	target, ok := documentByPath(repository, targetPath)
	if !ok {
		if _, err := os.Stat(targetPath); err != nil {
			diagnostic := Diagnostic{
				Code: "DAD-CONTEXT-003", Severity: "error",
				Message: "Referenced document does not exist.",
				Path:    source.Path, Line: link.Line, Column: link.Column,
			}
			return Document{}, false, &diagnostic
		}
		return Document{}, false, nil
	}
	return target, true, nil
}

func nonAuthoritativeDiagnostic(
	source Document,
	link MarkdownLink,
	target Document,
	required string,
) Diagnostic {
	return Diagnostic{
		Code: "DAD-CONTEXT-004", Severity: "error",
		Message: fmt.Sprintf(
			"%s has status %q; %q is required for governing context.",
			target.ID, target.Status, required,
		),
		Path: source.Path, Line: link.Line, Column: link.Column,
	}
}

func sortContextDocuments(documents []ContextDocument) {
	sort.SliceStable(documents, func(i, j int) bool {
		return documents[i].ID < documents[j].ID
	})
}
