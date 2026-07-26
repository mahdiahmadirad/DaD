package release

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type target struct {
	os         string
	arch       string
	extension  string
	executable string
}

var releaseTargets = []target{
	{os: "linux", arch: "amd64", extension: "tar.gz", executable: "dad"},
	{os: "linux", arch: "arm64", extension: "tar.gz", executable: "dad"},
	{os: "darwin", arch: "amd64", extension: "tar.gz", executable: "dad"},
	{os: "darwin", arch: "arm64", extension: "tar.gz", executable: "dad"},
	{os: "windows", arch: "amd64", extension: "zip", executable: "dad.exe"},
}

func VerifyArtifacts(directory, version string) error {
	version = strings.TrimPrefix(version, "v")
	expectedArchives := make(map[string]target, len(releaseTargets))
	for _, item := range releaseTargets {
		name := fmt.Sprintf(
			"dad_%s_%s_%s.%s",
			version, item.os, item.arch, item.extension,
		)
		expectedArchives[name] = item
	}
	checksumName := fmt.Sprintf("dad_%s_checksums.txt", version)

	actualArchives, err := artifactFiles(directory)
	if err != nil {
		return err
	}
	expectedNames := make([]string, 0, len(expectedArchives)+1)
	for name := range expectedArchives {
		expectedNames = append(expectedNames, name)
	}
	expectedNames = append(expectedNames, checksumName)
	sort.Strings(expectedNames)
	if strings.Join(actualArchives, "\n") != strings.Join(expectedNames, "\n") {
		return fmt.Errorf(
			"release artifacts differ:\ngot  %v\nwant %v",
			actualArchives, expectedNames,
		)
	}

	checksums, err := readChecksums(filepath.Join(directory, checksumName))
	if err != nil {
		return err
	}
	if len(checksums) != len(expectedArchives) {
		return fmt.Errorf(
			"checksum file covers %d artifacts, want %d",
			len(checksums), len(expectedArchives),
		)
	}

	for _, name := range expectedNames {
		item, archive := expectedArchives[name]
		if !archive {
			continue
		}
		archivePath := filepath.Join(directory, name)
		if err := verifyChecksum(archivePath, checksums[name]); err != nil {
			return err
		}
		execute := item.os == runtime.GOOS && item.arch == runtime.GOARCH
		files, executable, err := inspectArchive(archivePath, item, execute)
		if err != nil {
			return err
		}
		if executable != "" {
			defer os.RemoveAll(filepath.Dir(executable))
		}
		expectedFiles := []string{
			"CHANGELOG.md", "LICENSE", "README.md", item.executable,
		}
		sort.Strings(files)
		if strings.Join(files, "\n") != strings.Join(expectedFiles, "\n") {
			return fmt.Errorf(
				"%s contents differ: got %v, want %v",
				name, files, expectedFiles,
			)
		}
		if execute {
			if err := verifyExecutable(executable, version); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}
		}
	}
	return nil
}

func artifactFiles(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".tar.gz") ||
			strings.HasSuffix(name, ".zip") ||
			strings.HasSuffix(name, "_checksums.txt") {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	if len(names) == 0 {
		return nil, fmt.Errorf("no release artifacts found in %s", directory)
	}
	return names, nil
}

func readChecksums(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid checksum line %q", scanner.Text())
		}
		name := strings.TrimPrefix(fields[1], "*")
		if filepath.Base(name) != name {
			return nil, fmt.Errorf("checksum contains non-asset path %q", name)
		}
		if _, exists := result[name]; exists {
			return nil, fmt.Errorf("duplicate checksum for %s", name)
		}
		result[name] = fields[0]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func verifyChecksum(path, expected string) error {
	if expected == "" {
		return fmt.Errorf("checksum missing for %s", filepath.Base(path))
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	actual := hex.EncodeToString(hash.Sum(nil))
	if actual != expected {
		return fmt.Errorf(
			"checksum mismatch for %s: got %s, want %s",
			filepath.Base(path), actual, expected,
		)
	}
	return nil
}

func inspectArchive(
	path string,
	item target,
	extractExecutable bool,
) ([]string, string, error) {
	switch item.extension {
	case "zip":
		return inspectZip(path, item, extractExecutable)
	default:
		return inspectTarGzip(path, item, extractExecutable)
	}
}

func inspectZip(
	path string,
	item target,
	extractExecutable bool,
) ([]string, string, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, "", err
	}
	defer archive.Close()

	var names []string
	var executable string
	for _, file := range archive.File {
		if file.FileInfo().IsDir() {
			return nil, "", fmt.Errorf("%s contains directory %q", path, file.Name)
		}
		if filepath.Base(file.Name) != file.Name || strings.Contains(file.Name, `\`) {
			return nil, "", fmt.Errorf("%s contains unsafe path %q", path, file.Name)
		}
		names = append(names, file.Name)
		if extractExecutable && file.Name == item.executable {
			executable, err = extractFile(file.Open, file.Mode(), item.executable)
			if err != nil {
				return nil, "", err
			}
		}
	}
	return names, executable, nil
}

func inspectTarGzip(
	path string,
	item target,
	extractExecutable bool,
) ([]string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	compressed, err := gzip.NewReader(file)
	if err != nil {
		return nil, "", err
	}
	defer compressed.Close()

	var names []string
	var executable string
	reader := tar.NewReader(compressed)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, "", err
		}
		if header.Typeflag != tar.TypeReg {
			return nil, "", fmt.Errorf("%s contains non-file %q", path, header.Name)
		}
		if filepath.Base(header.Name) != header.Name ||
			strings.Contains(header.Name, `\`) {
			return nil, "", fmt.Errorf("%s contains unsafe path %q", path, header.Name)
		}
		names = append(names, header.Name)
		if header.Name == item.executable {
			if header.FileInfo().Mode()&0o111 == 0 {
				return nil, "", fmt.Errorf("%s is not executable", header.Name)
			}
			if extractExecutable {
				executable, err = extractFile(
					func() (io.ReadCloser, error) {
						return io.NopCloser(reader), nil
					},
					header.FileInfo().Mode(),
					item.executable,
				)
				if err != nil {
					return nil, "", err
				}
			}
		}
	}
	return names, executable, nil
}

func extractFile(
	open func() (io.ReadCloser, error),
	mode os.FileMode,
	name string,
) (string, error) {
	source, err := open()
	if err != nil {
		return "", err
	}
	defer source.Close()
	directory, err := os.MkdirTemp("", "dad-release-verify-*")
	if err != nil {
		return "", err
	}
	path := filepath.Join(directory, name)
	target, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode|0o700)
	if err != nil {
		os.RemoveAll(directory)
		return "", err
	}
	if _, err := io.Copy(target, source); err != nil {
		target.Close()
		os.RemoveAll(directory)
		return "", err
	}
	if err := target.Close(); err != nil {
		os.RemoveAll(directory)
		return "", err
	}
	return path, nil
}

func verifyExecutable(path, version string) error {
	if path == "" {
		return fmt.Errorf("archive has no executable")
	}
	defer os.RemoveAll(filepath.Dir(path))

	output, err := exec.Command(path, "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("text version failed: %w: %s", err, output)
	}
	if got, want := string(output), "dad "+version+"\n"; got != want {
		return fmt.Errorf("text version = %q, want %q", got, want)
	}

	output, err = exec.Command(path, "--format", "json", "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("JSON version failed: %w: %s", err, output)
	}
	var envelope struct {
		Success bool `json:"success"`
		Data    struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal(output, &envelope); err != nil {
		return fmt.Errorf("invalid JSON version output: %w", err)
	}
	if !envelope.Success || envelope.Data.Version != version {
		return fmt.Errorf(
			"JSON version success=%t version=%q, want true and %q",
			envelope.Success, envelope.Data.Version, version,
		)
	}
	return nil
}
