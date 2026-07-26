package dad

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadRepository(root string) Repository {
	paths, pathDiagnostics := candidateDocumentPaths(root)
	repository := Repository{
		Root:        root,
		ByID:        make(map[string][]int),
		Diagnostics: pathDiagnostics,
	}
	for _, path := range paths {
		relative := RelativePath(root, path)
		document, diagnostics := ParseDocument(path, relative)
		repository.Diagnostics = append(repository.Diagnostics, diagnostics...)
		if document.ID == "" {
			continue
		}
		index := len(repository.Documents)
		repository.Documents = append(repository.Documents, document)
		repository.ByID[strings.ToUpper(document.ID)] = append(
			repository.ByID[strings.ToUpper(document.ID)], index,
		)
	}
	SortDocuments(repository.Documents)
	repository.reindex()
	for id, indexes := range repository.ByID {
		if len(indexes) < 2 {
			continue
		}
		for _, index := range indexes {
			repository.Diagnostics = append(repository.Diagnostics, Diagnostic{
				Code: "DAD-DOC-002", Severity: "error",
				Message: fmt.Sprintf("Duplicate document identifier %s.", id),
				Path:    repository.Documents[index].Path,
			})
		}
	}
	SortDiagnostics(repository.Diagnostics)
	return repository
}

func (r *Repository) reindex() {
	r.ByID = make(map[string][]int, len(r.Documents))
	for index := range r.Documents {
		id := strings.ToUpper(r.Documents[index].ID)
		r.ByID[id] = append(r.ByID[id], index)
	}
}

func candidateDocumentPaths(root string) ([]string, []Diagnostic) {
	var paths []string
	var diagnostics []Diagnostic
	addPath := func(path string) {
		if !WithinRoot(root, path) {
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-PATH-003", Severity: "error",
				Message: "Governed document path resolves outside the repository root.",
				Path:    RelativePath(root, path),
			})
			return
		}
		paths = append(paths, path)
	}
	addDirectoryFiles := func(directory string, expectedType DocumentType) {
		if info, err := os.Lstat(directory); err == nil && info.Mode()&os.ModeSymlink != 0 &&
			!WithinRoot(root, directory) {
			diagnostics = append(diagnostics, Diagnostic{
				Code: "DAD-PATH-003", Severity: "error",
				Message: "Governed document directory resolves outside the repository root.",
				Path:    RelativePath(root, directory),
			})
			return
		}
		entries, err := os.ReadDir(directory)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
				continue
			}
			if strings.HasPrefix(
				strings.ToUpper(entry.Name()), string(expectedType)+"-",
			) {
				addPath(filepath.Join(directory, entry.Name()))
			}
		}
	}
	addDirectoryFiles(filepath.Join(root, "docs", "adr"), ADR)
	addDirectoryFiles(filepath.Join(root, "docs", "specs"), SPEC)
	entries, err := os.ReadDir(root)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
				continue
			}
			if strings.HasPrefix(strings.ToUpper(entry.Name()), "TASK-") {
				addPath(filepath.Join(root, entry.Name()))
			}
		}
	}
	sortPaths(paths)
	SortDiagnostics(diagnostics)
	return paths, diagnostics
}

func sortPaths(paths []string) {
	for i := 1; i < len(paths); i++ {
		for j := i; j > 0 && paths[j] < paths[j-1]; j-- {
			paths[j], paths[j-1] = paths[j-1], paths[j]
		}
	}
}

func (r Repository) MaxNumber(documentType DocumentType) int {
	maximum := 0
	for _, document := range r.Documents {
		if document.Type == documentType && document.Number > maximum {
			maximum = document.Number
		}
	}
	return maximum
}
