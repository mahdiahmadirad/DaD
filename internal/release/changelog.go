package release

import (
	"fmt"
	"regexp"
	"strings"
)

var releaseHeading = regexp.MustCompile(
	`(?m)^## \[([0-9A-Za-z.+-]+)\] - ([0-9]{4}-[0-9]{2}-[0-9]{2})[ \t]*$`,
)

func ChangelogNotes(content, version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	matches := releaseHeading.FindAllStringSubmatchIndex(content, -1)
	for index, match := range matches {
		if content[match[2]:match[3]] != version {
			continue
		}
		end := len(content)
		if index+1 < len(matches) {
			end = matches[index+1][0]
		} else if next := strings.Index(content[match[1]:], "\n## ["); next >= 0 {
			end = match[1] + next
		}
		notes := strings.TrimSpace(content[match[1]:end])
		if !substantiveNotes(notes) {
			return "", fmt.Errorf("changelog entry for %s is empty", version)
		}
		return notes + "\n", nil
	}
	return "", fmt.Errorf("dated changelog entry for %s does not exist", version)
}

func substantiveNotes(notes string) bool {
	for _, line := range strings.Split(notes, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") && len(line) > 2 {
			return true
		}
	}
	return false
}
