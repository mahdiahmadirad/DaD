package dad

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RootError struct {
	Message string
}

func (e RootError) Error() string {
	return e.Message
}

func ResolveRoot(explicit string, environment map[string]string, cwd string) (string, error) {
	candidate := explicit
	if candidate == "" {
		candidate = environment["DAD_ROOT"]
	}
	if candidate != "" {
		root, err := absoluteExistingDirectory(candidate, cwd)
		if err != nil {
			return "", RootError{Message: err.Error()}
		}
		if !hasRootMarkers(root) {
			return "", RootError{
				Message: fmt.Sprintf(
					"%s is not a DaD root: PROJECT-VISION.md and AGENTS.md are required",
					root,
				),
			}
		}
		return root, nil
	}

	current, err := absoluteExistingDirectory(cwd, cwd)
	if err != nil {
		return "", RootError{Message: err.Error()}
	}
	for {
		if hasRootMarkers(current) {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", RootError{
		Message: "DaD root not found: use --root or set DAD_ROOT",
	}
}

func ResolveInitTarget(value, cwd string) (string, error) {
	if value == "" {
		value = cwd
	}
	root, err := absoluteExistingDirectory(value, cwd)
	if err != nil {
		return "", RootError{Message: err.Error()}
	}
	return root, nil
}

func absoluteExistingDirectory(value, cwd string) (string, error) {
	if !filepath.IsAbs(value) {
		value = filepath.Join(cwd, value)
	}
	absolute, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("resolve path %q: %w", value, err)
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return "", fmt.Errorf("access path %q: %w", absolute, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path %q is not a directory", absolute)
	}
	evaluated, err := filepath.EvalSymlinks(absolute)
	if err == nil {
		absolute = evaluated
	}
	return filepath.Clean(absolute), nil
}

func hasRootMarkers(root string) bool {
	return nonEmptyFile(filepath.Join(root, "PROJECT-VISION.md")) &&
		nonEmptyFile(filepath.Join(root, "AGENTS.md"))
}

func nonEmptyFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular() && info.Size() > 0
}

func WithinRoot(root, path string) bool {
	root = canonicalForComparison(root)
	path = canonicalForComparison(path)
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return relative != ".." &&
		!strings.HasPrefix(relative, ".."+string(filepath.Separator)) &&
		!filepath.IsAbs(relative)
}

func canonicalForComparison(path string) string {
	absolute, err := filepath.Abs(path)
	if err == nil {
		path = absolute
	}
	if evaluated, err := filepath.EvalSymlinks(path); err == nil {
		return filepath.Clean(evaluated)
	}
	var missing []string
	current := filepath.Clean(path)
	for {
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		missing = append(missing, filepath.Base(current))
		current = parent
		if evaluated, err := filepath.EvalSymlinks(current); err == nil {
			current = evaluated
			for index := len(missing) - 1; index >= 0; index-- {
				current = filepath.Join(current, missing[index])
			}
			return filepath.Clean(current)
		}
	}
	return filepath.Clean(path)
}

func ResolveReference(root, source, target string) (string, error) {
	if target == "" {
		return "", errors.New("empty reference")
	}
	if strings.HasPrefix(target, "#") {
		return source, nil
	}
	target = strings.SplitN(target, "#", 2)[0]
	target = filepath.FromSlash(target)
	resolved := filepath.Clean(filepath.Join(filepath.Dir(source), target))
	if !WithinRoot(root, resolved) {
		return "", fmt.Errorf("reference resolves outside repository root")
	}
	if evaluated, err := filepath.EvalSymlinks(resolved); err == nil {
		if !WithinRoot(root, evaluated) {
			return "", fmt.Errorf("reference resolves outside repository root")
		}
		resolved = evaluated
	}
	return resolved, nil
}
