package store

import (
	"fmt"
	"os"
	"path/filepath"

	"lims/internal/errs"
)

// WriteFile persists report content under report-files and returns the
// relative path used inside the report entity.
func (fs *FileStore) WriteFile(kind, id string, version int, content []byte) (string, error) {
	name := fmt.Sprintf("%s-v%d.txt", id, version)
	path := filepath.Join(fs.root, "report-files", name)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return "", err
	}
	return "report-files/" + name, nil
}

// ReadFile loads report content from a stored relative path.
func (fs *FileStore) ReadFile(path string) ([]byte, error) {
	full := filepath.Join(fs.root, path)
	data, err := os.ReadFile(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", errs.ErrNotFound, path)
		}
		return nil, err
	}
	return data, nil
}

// FileExists reports whether a stored relative path is present.
func (fs *FileStore) FileExists(path string) (bool, error) {
	full := filepath.Join(fs.root, path)
	_, err := os.Stat(full)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
