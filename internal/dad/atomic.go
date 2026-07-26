package dad

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func AtomicWriteNew(path string, content []byte) error {
	if _, err := os.Lstat(path); err == nil {
		return os.ErrExist
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(directory, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	cleanup := func() {
		_ = temporary.Close()
		_ = os.Remove(temporaryPath)
	}
	if err := temporary.Chmod(0o644); err != nil {
		cleanup()
		return err
	}
	if _, err := temporary.Write(content); err != nil {
		cleanup()
		return err
	}
	if err := temporary.Sync(); err != nil {
		cleanup()
		return err
	}
	if err := temporary.Close(); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}
	if _, err := os.Lstat(path); err == nil {
		_ = os.Remove(temporaryPath)
		return os.ErrExist
	} else if !errors.Is(err, os.ErrNotExist) {
		_ = os.Remove(temporaryPath)
		return err
	}
	// Linking a fully written sibling temp file creates the target atomically
	// and, unlike rename on POSIX systems, never replaces an existing target.
	if err := os.Link(temporaryPath, path); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}
	_ = os.Remove(temporaryPath)
	return nil
}

func acquireCreationLock(directory string, documentType DocumentType) (func(), error) {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(directory, fmt.Sprintf(".dad-new-%s.lock", documentType))
	for attempt := 0; attempt < 100; attempt++ {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			_, _ = fmt.Fprintf(file, "%d\n", os.Getpid())
			_ = file.Close()
			return func() { _ = os.Remove(path) }, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil, fmt.Errorf("document creation is locked by another process")
}
