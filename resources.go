package dad

import (
	"embed"
	"fmt"
)

// resources embeds authoritative repository files directly. The embedded
// bytes are build artifacts; the Markdown files remain their only maintained
// source.
//
//go:embed docs/DOCUMENTATION.md docs/templates/*.md prompts/*.md
var resources embed.FS

var supportPaths = []string{
	"docs/DOCUMENTATION.md",
	"docs/templates/ADR-TEMPLATE.md",
	"docs/templates/SPEC-TEMPLATE.md",
	"docs/templates/TASK-TEMPLATE.md",
}

type PromptInfo struct {
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
	path    string
}

var promptInfos = []PromptInfo{
	{
		Name:    "project-bootstrap",
		Purpose: "Establish the minimum durable DaD context in a repository.",
		path:    "prompts/PROJECT-BOOTSTRAP.md",
	},
	{
		Name:    "architecture-discovery",
		Purpose: "Build an evidence-based understanding of an existing architecture.",
		path:    "prompts/ARCHITECTURE-DISCOVERY.md",
	},
	{
		Name:    "documentation-audit",
		Purpose: "Evaluate documentation authority, connections, currency, and usefulness.",
		path:    "prompts/DOCUMENTATION-AUDIT.md",
	},
	{
		Name:    "task-implementation",
		Purpose: "Implement one approved DaD TASK within its boundaries.",
		path:    "prompts/TASK-IMPLEMENTATION.md",
	},
	{
		Name:    "documentation-reconciliation",
		Purpose: "Restore agreement between documentation, implementation, and evidence.",
		path:    "prompts/DOCUMENTATION-RECONCILIATION.md",
	},
}

func SupportFiles() (map[string][]byte, error) {
	files := make(map[string][]byte, len(supportPaths))
	for _, path := range supportPaths {
		content, err := resources.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read embedded support file %s: %w", path, err)
		}
		files[path] = append([]byte(nil), content...)
	}
	return files, nil
}

func Prompts() []PromptInfo {
	result := make([]PromptInfo, len(promptInfos))
	copy(result, promptInfos)
	return result
}

func Prompt(name string) (PromptInfo, []byte, bool, error) {
	for _, info := range promptInfos {
		if info.Name != name {
			continue
		}
		content, err := resources.ReadFile(info.path)
		if err != nil {
			return PromptInfo{}, nil, false, fmt.Errorf(
				"read embedded prompt %s: %w", name, err,
			)
		}
		return info, append([]byte(nil), content...), true, nil
	}
	return PromptInfo{}, nil, false, nil
}
