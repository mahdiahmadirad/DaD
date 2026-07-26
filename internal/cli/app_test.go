package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testEnvelope struct {
	SchemaVersion string          `json:"schema_version"`
	Command       string          `json:"command"`
	Success       bool            `json:"success"`
	Data          json.RawMessage `json:"data"`
	Diagnostics   []struct {
		Code string `json:"code"`
	} `json:"diagnostics"`
}

func runCLI(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := Run(args, &stdout, &stderr, []string{})
	return code, stdout.String(), stderr.String()
}

func writeCLIFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestJSONUsageErrorUsesEnvelopeAndNoStderr(t *testing.T) {
	code, stdout, stderr := runCLI(t, "--format", "json", "unknown")
	if code != 2 || stderr != "" {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
	var envelope testEnvelope
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Success || envelope.SchemaVersion != "1" ||
		len(envelope.Diagnostics) != 1 {
		t.Fatalf("envelope=%#v", envelope)
	}
}

func TestBooleanGlobalOptionRejectsValue(t *testing.T) {
	code, _, stderr := runCLI(t, "--quiet=false", "prompt", "list")
	if code != 2 || !strings.Contains(stderr, "does not accept a value") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestPromptCommandsNeedNoRepository(t *testing.T) {
	code, stdout, stderr := runCLI(t, "prompt", "list")
	if code != 0 || stderr != "" {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "project-bootstrap") ||
		!strings.Contains(stdout, "documentation-reconciliation") {
		t.Fatalf("unexpected prompt list: %q", stdout)
	}
	code, stdout, stderr = runCLI(t, "prompt", "show", "task-implementation")
	if code != 0 || stderr != "" || !strings.HasPrefix(stdout, "# Task Implementation Prompt") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestVersionOutputsAgree(t *testing.T) {
	previous := releaseVersion
	releaseVersion = "1.2.3-rc.1"
	t.Cleanup(func() {
		releaseVersion = previous
	})

	code, stdout, stderr := runCLI(t, "--version")
	if code != 0 || stdout != "dad 1.2.3-rc.1\n" || stderr != "" {
		t.Fatalf("text code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	code, stdout, stderr = runCLI(t, "--format", "json", "--version")
	if code != 0 || stderr != "" {
		t.Fatalf("JSON code=%d stderr=%q", code, stderr)
	}
	var versionEnvelope struct {
		Data struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &versionEnvelope); err != nil {
		t.Fatal(err)
	}
	if versionEnvelope.Data.Version != "1.2.3-rc.1" {
		t.Fatalf("JSON version=%q", versionEnvelope.Data.Version)
	}

	code, stdout, stderr = runCLI(
		t, "--format", "json", "prompt", "show", "task-implementation",
	)
	if code != 0 || stderr != "" {
		t.Fatalf("prompt code=%d stderr=%q", code, stderr)
	}
	var promptEnvelope struct {
		Data struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &promptEnvelope); err != nil {
		t.Fatal(err)
	}
	if promptEnvelope.Data.Version != "1.2.3-rc.1" {
		t.Fatalf("prompt version=%q", promptEnvelope.Data.Version)
	}
}

func TestEndToEndRepositoryCommands(t *testing.T) {
	root := t.TempDir()
	writeCLIFile(t, root, "PROJECT-VISION.md", "# Vision\n")
	writeCLIFile(t, root, "AGENTS.md", "# Agreement\n")

	code, stdout, stderr := runCLI(t, "--root", root, "--dry-run", "init")
	if code != 0 || stderr != "" || !strings.Contains(stdout, "would create") {
		t.Fatalf("dry init code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if _, err := os.Stat(filepath.Join(root, "docs")); !os.IsNotExist(err) {
		t.Fatalf("dry init wrote docs: %v", err)
	}

	code, _, stderr = runCLI(t, "init", root)
	if code != 0 || stderr != "" {
		t.Fatalf("init code=%d stderr=%q", code, stderr)
	}
	code, stdout, stderr = runCLI(
		t, "--root", root, "new", "task", "--title", "Implement fixture",
	)
	if code != 0 || stderr != "" || !strings.Contains(stdout, "TASK-0001") {
		t.Fatalf("new code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}

	code, stdout, stderr = runCLI(t, "--root", root, "--format", "json", "list", "task")
	if code != 0 || stderr != "" {
		t.Fatalf("list code=%d stderr=%q", code, stderr)
	}
	var listEnvelope testEnvelope
	if err := json.Unmarshal([]byte(stdout), &listEnvelope); err != nil {
		t.Fatal(err)
	}
	if !listEnvelope.Success || listEnvelope.Command != "list" ||
		!bytes.Contains(listEnvelope.Data, []byte("TASK-0001")) {
		t.Fatalf("list envelope=%s", stdout)
	}

	code, stdout, stderr = runCLI(t, "--root", root, "check")
	if code != 0 || stderr != "" || !strings.Contains(stdout, "0 errors") {
		t.Fatalf("check code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	code, stdout, stderr = runCLI(t, "--root", root, "status")
	if code != 0 || stderr != "" || !strings.Contains(stdout, "TASK Proposed 1") {
		t.Fatalf("status code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	code, stdout, stderr = runCLI(t, "--root", root, "context", "task-0001")
	if code != 0 || stderr != "" || !strings.Contains(stdout, "project-vision") {
		t.Fatalf("context code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func TestNewDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	writeCLIFile(t, root, "PROJECT-VISION.md", "# Vision\n")
	writeCLIFile(t, root, "AGENTS.md", "# Agreement\n")
	code, _, stderr := runCLI(t, "init", root)
	if code != 0 {
		t.Fatalf("init code=%d stderr=%q", code, stderr)
	}
	code, stdout, stderr := runCLI(
		t, "--root", root, "--dry-run", "new", "adr", "--title", "Preview",
	)
	if code != 0 || stderr != "" || !strings.Contains(stdout, "would create ADR-0001") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "adr")); !os.IsNotExist(err) {
		t.Fatalf("dry new created directory: %v", err)
	}
}

func TestInitRejectsPathAndRootTogether(t *testing.T) {
	root := t.TempDir()
	code, _, stderr := runCLI(t, "--root", root, "init", root)
	if code != 2 || !strings.Contains(stderr, "either") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestDADRootEnvironmentIsUsed(t *testing.T) {
	root := t.TempDir()
	writeCLIFile(t, root, "PROJECT-VISION.md", "# Vision\n")
	writeCLIFile(t, root, "AGENTS.md", "# Agreement\n")
	code, _, stderr := runCLI(t, "init", root)
	if code != 0 {
		t.Fatalf("init code=%d stderr=%q", code, stderr)
	}
	var stdout bytes.Buffer
	var errorOutput bytes.Buffer
	code = Run(
		[]string{"list"}, &stdout, &errorOutput,
		[]string{"DAD_ROOT=" + root},
	)
	if code != 0 || errorOutput.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, errorOutput.String())
	}
}
