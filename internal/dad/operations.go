package dad

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

type InitResult struct {
	Root      string   `json:"root"`
	Created   []string `json:"created"`
	Unchanged []string `json:"unchanged"`
	Conflicts []string `json:"conflicts"`
}

func Init(root string, supportFiles map[string][]byte, dryRun bool) (
	InitResult, []Diagnostic, int,
) {
	result := InitResult{Root: root}
	if !nonEmptyFile(filepath.Join(root, "PROJECT-VISION.md")) ||
		!nonEmptyFile(filepath.Join(root, "AGENTS.md")) {
		return result, []Diagnostic{{
			Code: "DAD-ROOT-002", Severity: "error",
			Message: "PROJECT-VISION.md and AGENTS.md must exist and be non-empty; " +
				"start with `dad prompt show project-bootstrap`.",
		}}, ExitRoot
	}
	paths := make([]string, 0, len(supportFiles))
	for path := range supportFiles {
		paths = append(paths, path)
	}
	sortPaths(paths)
	for _, relative := range paths {
		target := filepath.Join(root, filepath.FromSlash(relative))
		if !WithinRoot(root, target) {
			result.Conflicts = append(result.Conflicts, relative)
			continue
		}
		existing, err := os.ReadFile(target)
		switch {
		case err == nil && bytes.Equal(existing, supportFiles[relative]):
			result.Unchanged = append(result.Unchanged, relative)
		case err == nil:
			result.Conflicts = append(result.Conflicts, relative)
		case errors.Is(err, os.ErrNotExist):
			result.Created = append(result.Created, relative)
		default:
			return result, []Diagnostic{{
				Code: "DAD-IO-002", Severity: "error",
				Message: fmt.Sprintf("Cannot inspect initialization target: %v", err),
				Path:    relative,
			}}, ExitIO
		}
	}
	if len(result.Conflicts) > 0 {
		diagnostics := make([]Diagnostic, 0, len(result.Conflicts))
		for _, path := range result.Conflicts {
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-INIT-001", Severity: "error",
				Message: "Existing support file differs from the canonical resource.",
				Path:    path,
			})
		}
		result.Created = nil
		return result, diagnostics, ExitConflict
	}
	if dryRun {
		return result, nil, ExitOK
	}
	planned := append([]string(nil), result.Created...)
	result.Created = nil
	for _, relative := range planned {
		target := filepath.Join(root, filepath.FromSlash(relative))
		if err := AtomicWriteNew(target, supportFiles[relative]); err != nil {
			return result, []Diagnostic{{
				Code: "DAD-IO-003", Severity: "error",
				Message: fmt.Sprintf("Cannot create support file: %v", err),
				Path:    relative,
			}}, ExitIO
		}
		result.Created = append(result.Created, relative)
	}
	return result, nil, ExitOK
}

type NewOptions struct {
	Type     DocumentType
	Title    string
	Number   int
	DryRun   bool
	Template []byte
}

type NewResult struct {
	ID      string       `json:"id"`
	Type    DocumentType `json:"type"`
	Title   string       `json:"title"`
	Path    string       `json:"path"`
	Status  string       `json:"status"`
	Created bool         `json:"created"`
}

func NewDocument(root string, options NewOptions) (NewResult, []Diagnostic, int) {
	title := strings.TrimSpace(options.Title)
	if title == "" || strings.ContainsAny(title, "\r\n") {
		return NewResult{}, []Diagnostic{{
			Code: "DAD-NEW-001", Severity: "error",
			Message: "Title must be non-empty, trimmed, and contain one line.",
		}}, ExitUsage
	}
	localTemplate := filepath.Join(
		root, "docs", "templates", string(options.Type)+"-TEMPLATE.md",
	)
	template := options.Template
	if content, err := os.ReadFile(localTemplate); err == nil {
		template = content
	} else if !errors.Is(err, os.ErrNotExist) {
		return NewResult{}, []Diagnostic{{
			Code: "DAD-IO-004", Severity: "error",
			Message: fmt.Sprintf("Cannot read local template: %v", err),
			Path:    RelativePath(root, localTemplate),
		}}, ExitIO
	}
	if !utf8.Valid(template) || !validTemplate(options.Type, string(template)) {
		return NewResult{}, []Diagnostic{{
			Code: "DAD-NEW-002", Severity: "error",
			Message: "Selected template is malformed.",
			Path:    RelativePath(root, localTemplate),
		}}, ExitConflict
	}

	directory := root
	switch options.Type {
	case ADR:
		directory = filepath.Join(root, "docs", "adr")
	case SPEC:
		directory = filepath.Join(root, "docs", "specs")
	}
	if !WithinRoot(root, directory) {
		return NewResult{}, []Diagnostic{{
			Code: "DAD-PATH-002", Severity: "error",
			Message: "Document directory resolves outside the repository root.",
			Path:    RelativePath(root, directory),
		}}, ExitConflict
	}
	if options.DryRun {
		repository := LoadRepository(root)
		if hasErrorDiagnostic(repository.Diagnostics) {
			return NewResult{}, repository.Diagnostics, ExitConflict
		}
		number, diagnostic := requestedNumber(repository, options.Type, options.Number)
		if diagnostic != nil {
			return NewResult{}, []Diagnostic{*diagnostic}, ExitConflict
		}
		return buildNewResult(options.Type, title, number, false), nil, ExitOK
	}

	release, err := acquireCreationLock(directory, options.Type)
	if err != nil {
		return NewResult{}, []Diagnostic{{
			Code: "DAD-NEW-003", Severity: "error",
			Message: fmt.Sprintf("Cannot acquire document creation lock: %v", err),
		}}, ExitConflict
	}
	defer release()

	repository := LoadRepository(root)
	if hasErrorDiagnostic(repository.Diagnostics) {
		return NewResult{}, repository.Diagnostics, ExitConflict
	}
	number, diagnostic := requestedNumber(repository, options.Type, options.Number)
	if diagnostic != nil {
		return NewResult{}, []Diagnostic{*diagnostic}, ExitConflict
	}
	result := buildNewResult(options.Type, title, number, true)
	content := renderTemplate(options.Type, result.ID, title, string(template))
	target := filepath.Join(root, filepath.FromSlash(result.Path))
	if !WithinRoot(root, target) {
		return NewResult{}, []Diagnostic{{
			Code: "DAD-PATH-001", Severity: "error",
			Message: "Document target resolves outside the repository root.",
			Path:    result.Path,
		}}, ExitConflict
	}
	if err := AtomicWriteNew(target, []byte(content)); err != nil {
		exit := ExitIO
		code := "DAD-IO-005"
		if errors.Is(err, os.ErrExist) {
			exit, code = ExitConflict, "DAD-NEW-004"
		}
		return NewResult{}, []Diagnostic{{
			Code: code, Severity: "error",
			Message: fmt.Sprintf("Cannot create document: %v", err),
			Path:    result.Path,
		}}, exit
	}
	return result, nil, ExitOK
}

func hasErrorDiagnostic(diagnostics []Diagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == "error" {
			return true
		}
	}
	return false
}

func requestedNumber(repository Repository, documentType DocumentType, requested int) (
	int, *Diagnostic,
) {
	maximum := repository.MaxNumber(documentType)
	if requested == 0 {
		requested = maximum + 1
		if documentType == TASK && requested == 0 {
			requested = 1
		}
	}
	if requested < 1 || requested > 9999 ||
		(documentType == TASK && requested == 0) {
		return 0, &Diagnostic{
			Code: "DAD-NEW-005", Severity: "error",
			Message: "Requested number must be between 0001 and 9999.",
		}
	}
	if requested <= maximum {
		return 0, &Diagnostic{
			Code: "DAD-NEW-006", Severity: "error",
			Message: fmt.Sprintf(
				"Requested number %04d must be greater than existing maximum %04d.",
				requested, maximum,
			),
		}
	}
	return requested, nil
}

func buildNewResult(documentType DocumentType, title string, number int, created bool) NewResult {
	id := fmt.Sprintf("%s-%04d", documentType, number)
	path := id + ".md"
	status := "Proposed"
	switch documentType {
	case ADR:
		path = filepath.ToSlash(filepath.Join("docs", "adr", path))
	case SPEC:
		path = filepath.ToSlash(filepath.Join("docs", "specs", path))
		status = "Draft"
	}
	return NewResult{
		ID: id, Type: documentType, Title: title, Path: path,
		Status: status, Created: created,
	}
}

func validTemplate(documentType DocumentType, template string) bool {
	title := fmt.Sprintf("# %s-NNNN: Title", documentType)
	if !strings.Contains(template, title) || !strings.Contains(template, "## Status") {
		return false
	}
	requiredStatus := "Proposed"
	if documentType == SPEC {
		requiredStatus = "Draft"
	}
	if !strings.Contains(template, "\n"+requiredStatus+"\n") {
		return false
	}
	for _, section := range requiredTemplateSections[documentType] {
		if !strings.Contains(template, "## "+section) {
			return false
		}
	}
	return true
}

func renderTemplate(documentType DocumentType, id, title, template string) string {
	placeholder := fmt.Sprintf("%s-NNNN: Title", documentType)
	rendered := strings.Replace(template, placeholder, id+": "+title, 1)
	rendered = strings.ReplaceAll(rendered, "\r\n", "\n")
	rendered = strings.ReplaceAll(rendered, "\r", "\n")
	if !strings.HasSuffix(rendered, "\n") {
		rendered += "\n"
	}
	return rendered
}

func ParseExplicitNumber(value string) (int, error) {
	if len(value) != 4 {
		return 0, fmt.Errorf("number must contain exactly four digits")
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return 0, fmt.Errorf("number must contain exactly four digits")
		}
	}
	number, _ := strconv.Atoi(value)
	if number == 0 {
		return 0, fmt.Errorf("number 0000 is reserved")
	}
	return number, nil
}
