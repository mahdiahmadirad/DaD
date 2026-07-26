package dad

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

const (
	ExitOK        = 0
	ExitCheck     = 1
	ExitUsage     = 2
	ExitRoot      = 3
	ExitConflict  = 4
	ExitIO        = 5
	ExitInternal  = 6
	ExitInterrupt = 130
)

type Diagnostic struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Path     string `json:"path,omitempty"`
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
}

func (d Diagnostic) String() string {
	location := d.Path
	if d.Line > 0 {
		location += fmt.Sprintf(":%d", d.Line)
		if d.Column > 0 {
			location += fmt.Sprintf(":%d", d.Column)
		}
	}
	if location != "" {
		return fmt.Sprintf("%s: %s: %s", location, d.Severity, d.Message)
	}
	return fmt.Sprintf("%s: %s", d.Severity, d.Message)
}

func SortDiagnostics(diagnostics []Diagnostic) {
	rank := map[string]int{"error": 0, "warning": 1, "info": 2}
	sort.SliceStable(diagnostics, func(i, j int) bool {
		left, right := diagnostics[i], diagnostics[j]
		if left.Path != right.Path {
			return left.Path < right.Path
		}
		if left.Line != right.Line {
			return left.Line < right.Line
		}
		if left.Column != right.Column {
			return left.Column < right.Column
		}
		if rank[left.Severity] != rank[right.Severity] {
			return rank[left.Severity] < rank[right.Severity]
		}
		return left.Code < right.Code
	})
}

type DocumentType string

const (
	ADR  DocumentType = "ADR"
	SPEC DocumentType = "SPEC"
	TASK DocumentType = "TASK"
)

var typeOrder = map[DocumentType]int{ADR: 0, SPEC: 1, TASK: 2}

var allowedStatuses = map[DocumentType]map[string]bool{
	ADR: {
		"Proposed":   true,
		"Accepted":   true,
		"Rejected":   true,
		"Superseded": true,
	},
	SPEC: {
		"Draft":      true,
		"Approved":   true,
		"Withdrawn":  true,
		"Superseded": true,
		"Retired":    true,
	},
	TASK: {
		"Proposed":    true,
		"Ready":       true,
		"In Progress": true,
		"Blocked":     true,
		"Complete":    true,
		"Cancelled":   true,
	},
}

func ParseType(value string) (DocumentType, bool) {
	switch strings.ToUpper(value) {
	case "ADR":
		return ADR, true
	case "SPEC":
		return SPEC, true
	case "TASK":
		return TASK, true
	default:
		return "", false
	}
}

func StatusAllowed(documentType DocumentType, status string) bool {
	return allowedStatuses[documentType][status]
}

func AnyStatusAllowed(status string) bool {
	for _, statuses := range allowedStatuses {
		if statuses[status] {
			return true
		}
	}
	return false
}

func AuthoritativeStatus(documentType DocumentType) string {
	switch documentType {
	case ADR:
		return "Accepted"
	case SPEC:
		return "Approved"
	default:
		return ""
	}
}

type MarkdownLink struct {
	Label  string
	Target string
	Line   int
	Column int
}

type Document struct {
	ID       string            `json:"id"`
	Type     DocumentType      `json:"type"`
	Number   int               `json:"-"`
	Title    string            `json:"title"`
	Status   string            `json:"status"`
	Path     string            `json:"path"`
	AbsPath  string            `json:"-"`
	Content  string            `json:"-"`
	Sections map[string]string `json:"-"`
	Links    []MarkdownLink    `json:"-"`
}

func (d Document) CanonicalFilename() string {
	return d.ID + ".md"
}

func (d Document) ExpectedPath() string {
	switch d.Type {
	case ADR:
		return filepath.ToSlash(filepath.Join("docs", "adr", d.CanonicalFilename()))
	case SPEC:
		return filepath.ToSlash(filepath.Join("docs", "specs", d.CanonicalFilename()))
	case TASK:
		return d.CanonicalFilename()
	default:
		return ""
	}
}

func SortDocuments(documents []Document) {
	sort.SliceStable(documents, func(i, j int) bool {
		left, right := documents[i], documents[j]
		if typeOrder[left.Type] != typeOrder[right.Type] {
			return typeOrder[left.Type] < typeOrder[right.Type]
		}
		if left.Number != right.Number {
			return left.Number < right.Number
		}
		return left.Path < right.Path
	})
}

type Repository struct {
	Root        string
	Documents   []Document
	Diagnostics []Diagnostic
	ByID        map[string][]int
}

func (r Repository) Document(id string) (Document, bool) {
	indexes := r.ByID[strings.ToUpper(id)]
	if len(indexes) != 1 {
		return Document{}, false
	}
	return r.Documents[indexes[0]], true
}

func RelativePath(root, path string) string {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relative)
}
