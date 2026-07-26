package dad

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type CheckSummary struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Info     int `json:"info"`
}

type CheckResult struct {
	Summary          CheckSummary `json:"summary"`
	CheckedDocuments []string     `json:"checked_documents"`
}

func Check(root string, identifiers []string, strict bool) (
	CheckResult, []Diagnostic, int,
) {
	repository := LoadRepository(root)
	diagnostics := append([]Diagnostic(nil), repository.Diagnostics...)
	selected := selectedDocuments(repository, identifiers, &diagnostics)
	full := len(identifiers) == 0

	if full {
		diagnostics = append(diagnostics, checkCoreFiles(root)...)
		diagnostics = append(diagnostics, checkTemplates(root)...)
	}
	for _, document := range selected {
		diagnostics = append(diagnostics, checkDocument(repository, document)...)
	}
	diagnostics = append(diagnostics, checkExactDuplicates(root, selected)...)
	diagnostics = append(diagnostics, checkLinks(root, selected, full)...)
	SortDiagnostics(diagnostics)

	result := CheckResult{}
	for _, document := range selected {
		result.CheckedDocuments = append(result.CheckedDocuments, document.ID)
	}
	sort.Strings(result.CheckedDocuments)
	for _, diagnostic := range diagnostics {
		switch diagnostic.Severity {
		case "error":
			result.Summary.Errors++
		case "warning":
			result.Summary.Warnings++
		default:
			result.Summary.Info++
		}
	}
	exit := ExitOK
	if result.Summary.Errors > 0 || strict && result.Summary.Warnings > 0 {
		exit = ExitCheck
	}
	return result, diagnostics, exit
}

func selectedDocuments(
	repository Repository,
	identifiers []string,
	diagnostics *[]Diagnostic,
) []Document {
	if len(identifiers) == 0 {
		return append([]Document(nil), repository.Documents...)
	}
	var selected []Document
	seen := make(map[string]bool)
	for _, value := range identifiers {
		id := strings.ToUpper(value)
		if seen[id] {
			continue
		}
		seen[id] = true
		document, ok := repository.Document(id)
		if !ok {
			*diagnostics = append(*diagnostics, Diagnostic{
				Code: "DAD-CHECK-002", Severity: "error",
				Message: fmt.Sprintf("Document %s does not exist or is duplicated.", id),
			})
			continue
		}
		selected = append(selected, document)
	}
	for index := 0; index < len(selected); index++ {
		source := selected[index]
		for _, link := range sectionLinks(source, "References") {
			if externalLink(link.Target) || strings.HasPrefix(link.Target, "#") {
				continue
			}
			targetPath, err := ResolveReference(repository.Root, source.AbsPath, link.Target)
			if err != nil {
				continue
			}
			target, ok := documentByPath(repository, targetPath)
			if !ok || seen[target.ID] {
				continue
			}
			seen[target.ID] = true
			selected = append(selected, target)
		}
	}
	SortDocuments(selected)
	return selected
}

func checkCoreFiles(root string) []Diagnostic {
	var diagnostics []Diagnostic
	for _, path := range []string{"PROJECT-VISION.md", "AGENTS.md"} {
		if !nonEmptyFile(filepath.Join(root, path)) {
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-CHECK-003", Severity: "error",
				Message: "Required core document is missing or empty.",
				Path:    path,
			})
		}
	}
	return diagnostics
}

var requiredTemplateSections = map[DocumentType][]string{
	ADR:  {"Status", "Context", "Decision", "Consequences", "Alternatives considered"},
	SPEC: {"Status", "Purpose", "Scope", "Out of scope", "Requirements", "Acceptance criteria"},
	TASK: {"Status", "Outcome", "In scope", "Out of scope", "Constraints",
		"Acceptance conditions", "Completion evidence"},
}

func checkTemplates(root string) []Diagnostic {
	var diagnostics []Diagnostic
	for _, documentType := range []DocumentType{ADR, SPEC, TASK} {
		relative := filepath.ToSlash(filepath.Join(
			"docs", "templates", string(documentType)+"-TEMPLATE.md",
		))
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-CHECK-004", Severity: "error",
				Message: "Canonical template is missing or unreadable.",
				Path:    relative,
			})
			continue
		}
		text := string(content)
		if !validTemplate(documentType, text) {
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-CHECK-017", Severity: "warning",
				Message: "Template omits a canonical placeholder, status, or section.",
				Path:    relative,
			})
		}
		for _, section := range requiredTemplateSections[documentType] {
			if !strings.Contains(text, "## "+section) {
				diagnostics = append(diagnostics, Diagnostic{
					Code: "DAD-CHECK-005", Severity: "warning",
					Message: fmt.Sprintf("Template omits required section %q.", section),
					Path:    relative,
				})
			}
		}
	}
	return diagnostics
}

func checkDocument(repository Repository, document Document) []Diagnostic {
	var diagnostics []Diagnostic
	if document.Path != document.ExpectedPath() {
		diagnostics = append(diagnostics, Diagnostic{
			Code: "DAD-CHECK-006", Severity: "error",
			Message: fmt.Sprintf(
				"Non-canonical document path; expected %s.", document.ExpectedPath(),
			),
			Path: document.Path, Line: 1, Column: 1,
		})
	}
	for _, section := range requiredTemplateSections[document.Type] {
		if _, exists := document.Sections[section]; !exists {
			severity := "error"
			if document.Type == SPEC && section == "Interfaces and behavior" {
				severity = "warning"
			}
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-CHECK-009", Severity: severity,
				Message: fmt.Sprintf("Document omits required section %q.", section),
				Path:    document.Path,
			})
		}
	}
	if document.Type == TASK {
		if (document.Status == "Ready" || document.Status == "In Progress") &&
			!substantive(document.Sections["Acceptance conditions"]) {
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-CHECK-010", Severity: "warning",
				Message: "Actionable TASK has no substantive acceptance conditions.",
				Path:    document.Path,
			})
		}
		if document.Status == "Complete" &&
			!substantiveCompletion(document.Sections["Completion evidence"]) {
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-CHECK-011", Severity: "error",
				Message: "Complete TASK has no substantive completion evidence.",
				Path:    document.Path,
			})
		}
	}
	if statusNeedsNote(document.Status) &&
		!substantive(document.Sections["Status note"]) {
		diagnostics = append(diagnostics, Diagnostic{
			Code: "DAD-CHECK-012", Severity: "warning",
			Message: fmt.Sprintf("Status %q has no explanatory Status note.", document.Status),
			Path:    document.Path,
		})
	}
	diagnostics = append(diagnostics, checkGoverningReferences(repository, document)...)
	if document.Status == "Superseded" {
		diagnostics = append(diagnostics, checkReplacement(repository, document)...)
	}
	return diagnostics
}

func substantive(value string) bool {
	value = strings.TrimSpace(value)
	return len(value) >= 20
}

func substantiveCompletion(value string) bool {
	if !substantive(value) {
		return false
	}
	lower := strings.ToLower(value)
	return !strings.Contains(lower, "not yet recorded") &&
		!strings.Contains(lower, "to be recorded") &&
		!strings.Contains(lower, "do not claim")
}

func statusNeedsNote(status string) bool {
	switch status {
	case "Blocked", "Cancelled", "Rejected", "Withdrawn", "Retired":
		return true
	default:
		return false
	}
}

func checkGoverningReferences(repository Repository, document Document) []Diagnostic {
	var diagnostics []Diagnostic
	for _, link := range sectionLinks(document, "References") {
		if externalLink(link.Target) || strings.HasPrefix(link.Target, "#") {
			continue
		}
		targetPath, err := ResolveReference(
			repository.Root, document.AbsPath, link.Target,
		)
		if err != nil {
			continue
		}
		target, ok := documentByPath(repository, targetPath)
		if !ok || target.Type != ADR && target.Type != SPEC {
			continue
		}
		required := AuthoritativeStatus(target.Type)
		if target.Status != required {
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-CHECK-013", Severity: "error",
				Message: fmt.Sprintf(
					"Governing reference %s has status %q; %q is required.",
					target.ID, target.Status, required,
				),
				Path: document.Path, Line: link.Line, Column: link.Column,
			})
		}
	}
	return diagnostics
}

func checkReplacement(repository Repository, document Document) []Diagnostic {
	for _, link := range sectionLinks(document, "References") {
		targetPath, err := ResolveReference(repository.Root, document.AbsPath, link.Target)
		if err != nil {
			continue
		}
		target, ok := documentByPath(repository, targetPath)
		if ok && target.Type == document.Type &&
			target.Status == AuthoritativeStatus(document.Type) {
			return nil
		}
	}
	return []Diagnostic{{
		Code: "DAD-CHECK-014", Severity: "error",
		Message: "Superseded document has no authoritative same-type replacement reference.",
		Path:    document.Path,
	}}
}

func documentByPath(repository Repository, path string) (Document, bool) {
	clean := canonicalForComparison(path)
	for _, document := range repository.Documents {
		if canonicalForComparison(document.AbsPath) == clean {
			return document, true
		}
	}
	return Document{}, false
}

func checkLinks(root string, selected []Document, full bool) []Diagnostic {
	var paths []string
	if full {
		paths = markdownPaths(root)
	} else {
		for _, document := range selected {
			paths = append(paths, document.AbsPath)
		}
	}
	var diagnostics []Diagnostic
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		relative := RelativePath(root, path)
		for _, link := range parseLinks(string(content)) {
			if externalLink(link.Target) || strings.HasPrefix(link.Target, "#") {
				continue
			}
			target, err := ResolveReference(root, path, link.Target)
			if err != nil {
				diagnostics = append(diagnostics, Diagnostic{
					Code: "DAD-CHECK-015", Severity: "error",
					Message: err.Error(), Path: relative,
					Line: link.Line, Column: link.Column,
				})
				continue
			}
			if _, err := os.Stat(target); err != nil {
				diagnostics = append(diagnostics, Diagnostic{
					Code: "DAD-CHECK-016", Severity: "error",
					Message: "Referenced local path does not exist.",
					Path:    relative, Line: link.Line, Column: link.Column,
				})
			}
		}
	}
	return diagnostics
}

func checkExactDuplicates(root string, documents []Document) []Diagnostic {
	type occurrence struct {
		path string
	}
	ignored := make(map[string]bool)
	for _, documentType := range []DocumentType{ADR, SPEC, TASK} {
		path := filepath.Join(
			root, "docs", "templates", string(documentType)+"-TEMPLATE.md",
		)
		if content, err := os.ReadFile(path); err == nil {
			for _, paragraph := range normalizedLongParagraphs(string(content)) {
				ignored[paragraph] = true
			}
		}
	}
	seen := make(map[string]occurrence)
	var diagnostics []Diagnostic
	for _, document := range documents {
		for _, normalized := range normalizedLongParagraphs(document.Content) {
			if ignored[normalized] {
				continue
			}
			if first, exists := seen[normalized]; exists && first.path != document.Path {
				diagnostics = append(diagnostics, Diagnostic{
					Code: "DAD-CHECK-018", Severity: "warning",
					Message: fmt.Sprintf(
						"Exact long-form wording is duplicated from %s.", first.path,
					),
					Path: document.Path,
				})
				continue
			}
			seen[normalized] = occurrence{path: document.Path}
		}
	}
	return diagnostics
}

func normalizedLongParagraphs(content string) []string {
	var result []string
	for _, paragraph := range strings.Split(content, "\n\n") {
		normalized := strings.Join(strings.Fields(paragraph), " ")
		if len(normalized) < 120 || strings.HasPrefix(normalized, "#") ||
			strings.HasPrefix(normalized, "```") {
			continue
		}
		result = append(result, normalized)
	}
	return result
}

func externalLink(target string) bool {
	lower := strings.ToLower(target)
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "mailto:")
}

func markdownPaths(root string) []string {
	var paths []string
	entries, _ := os.ReadDir(root)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			paths = append(paths, filepath.Join(root, entry.Name()))
		}
	}
	for _, directory := range []string{
		filepath.Join(root, "docs"),
		filepath.Join(root, "prompts"),
	} {
		_ = filepath.WalkDir(directory, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(path), ".md") {
				paths = append(paths, path)
			}
			return nil
		})
	}
	sortPaths(paths)
	return paths
}
