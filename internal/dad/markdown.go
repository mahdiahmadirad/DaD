package dad

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	titlePattern      = regexp.MustCompile(`(?m)^# (ADR|SPEC|TASK)-([0-9]{4}): (.+?)\s*$`)
	headingPattern    = regexp.MustCompile(`(?m)^## (.+?)\s*$`)
	linkPattern       = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	inlineCodePattern = regexp.MustCompile("`[^`]*`")
)

func ParseDocument(path, relative string) (Document, []Diagnostic) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return Document{}, []Diagnostic{{
			Code: "DAD-IO-001", Severity: "error",
			Message: fmt.Sprintf("Cannot read document: %v", err), Path: relative,
		}}
	}
	if !utf8.Valid(contentBytes) {
		return Document{}, []Diagnostic{{
			Code: "DAD-TEXT-001", Severity: "error",
			Message: "Document is not valid UTF-8.", Path: relative,
		}}
	}
	content := string(contentBytes)
	match := titlePattern.FindStringSubmatchIndex(content)
	if match == nil || match[0] != 0 {
		return Document{Path: relative, AbsPath: path, Content: content},
			[]Diagnostic{{
				Code: "DAD-DOC-001", Severity: "error",
				Message: "Canonical document heading is missing or malformed.",
				Path:    relative, Line: 1, Column: 1,
			}}
	}
	documentType, _ := ParseType(content[match[2]:match[3]])
	numberText := content[match[4]:match[5]]
	number, _ := strconv.Atoi(numberText)
	title := strings.TrimSpace(content[match[6]:match[7]])
	sections := parseSections(content)
	status := firstContentLine(sections["Status"])
	document := Document{
		ID:       fmt.Sprintf("%s-%s", documentType, numberText),
		Type:     documentType,
		Number:   number,
		Title:    title,
		Status:   status,
		Path:     relative,
		AbsPath:  path,
		Content:  content,
		Sections: sections,
		Links:    parseLinks(content),
	}
	var diagnostics []Diagnostic
	if !StatusAllowed(document.Type, document.Status) {
		diagnostics = append(diagnostics, Diagnostic{
			Code: "DAD-DOC-003", Severity: "error",
			Message: fmt.Sprintf(
				"Status %q is not valid for %s.", document.Status, document.Type,
			),
			Path: relative,
		})
	}
	return document, diagnostics
}

func parseSections(content string) map[string]string {
	matches := headingPattern.FindAllStringSubmatchIndex(content, -1)
	sections := make(map[string]string, len(matches))
	for index, match := range matches {
		name := strings.TrimSpace(content[match[2]:match[3]])
		start := match[1]
		end := len(content)
		if index+1 < len(matches) {
			end = matches[index+1][0]
		}
		sections[name] = strings.TrimSpace(content[start:end])
	}
	return sections
}

func firstContentLine(section string) string {
	for _, line := range strings.Split(section, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
	return ""
}

func parseLinks(content string) []MarkdownLink {
	var links []MarkdownLink
	lines := strings.Split(content, "\n")
	inFence := false
	var fence string
	for index, original := range lines {
		trimmed := strings.TrimSpace(original)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			marker := trimmed[:3]
			if !inFence {
				inFence, fence = true, marker
			} else if marker == fence {
				inFence, fence = false, ""
			}
			continue
		}
		if inFence {
			continue
		}
		line := inlineCodePattern.ReplaceAllString(original, "")
		for _, match := range linkPattern.FindAllStringSubmatchIndex(line, -1) {
			target := line[match[4]:match[5]]
			if decoded, err := url.PathUnescape(target); err == nil {
				target = decoded
			}
			links = append(links, MarkdownLink{
				Label:  line[match[2]:match[3]],
				Target: target,
				Line:   index + 1,
				Column: match[0] + 1,
			})
		}
	}
	return links
}

func sectionLinks(document Document, section string) []MarkdownLink {
	matches := headingPattern.FindAllStringSubmatchIndex(document.Content, -1)
	for index, match := range matches {
		name := strings.TrimSpace(document.Content[match[2]:match[3]])
		if name != section {
			continue
		}
		start := match[1]
		end := len(document.Content)
		if index+1 < len(matches) {
			end = matches[index+1][0]
		}
		lineOffset := strings.Count(document.Content[:start], "\n")
		links := parseLinks(document.Content[start:end])
		for linkIndex := range links {
			links[linkIndex].Line += lineOffset
		}
		return links
	}
	return nil
}
