package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"lims/internal/audit"
	"lims/internal/errs"
	"lims/internal/instrument"
	"lims/internal/ns"
	"lims/internal/quota"
	"lims/internal/report"
	"lims/internal/result"
	"lims/internal/sample"
	"lims/internal/task"
)

var dataKinds = []string{
	"ns",
	"samples",
	"registrations",
	"tasks",
	"assignments",
	"instruments",
	"instrument-dispatches",
	"results",
	"qc",
	"acks",
	"reports",
	"report-files",
	"audit",
	"quota",
}

// FileStore persists every LIMS entity as JSON files under one data root.
type FileStore struct {
	root string
}

// New creates a file store and ensures all data directories exist.
func New(root string) (*FileStore, error) {
	fs := &FileStore{root: root}
	for _, kind := range dataKinds {
		dir := filepath.Join(root, kind)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create data dir %s: %w", kind, err)
		}
	}
	return fs, nil
}

// Root returns the base directory of the file store.
func (fs *FileStore) Root() string {
	return fs.root
}

func (fs *FileStore) writeJSON(kind, id string, value any) error {
	path := filepath.Join(fs.root, kind, id+".json")
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (fs *FileStore) readJSON(kind, id string, value any) error {
	path := filepath.Join(fs.root, kind, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s/%s", errs.ErrNotFound, kind, id)
		}
		return err
	}
	return json.Unmarshal(data, value)
}

func (fs *FileStore) listAll(kind string) ([][]byte, error) {
	dir := filepath.Join(fs.root, kind)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make([][]byte, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		out = append(out, data)
	}
	return out, nil
}

// Compile-time checks that the file store satisfies every domain interface.
var (
	_ sample.Store       = (*FileStore)(nil)
	_ task.Store         = (*FileStore)(nil)
	_ instrument.Store   = (*FileStore)(nil)
	_ result.Store       = (*FileStore)(nil)
	_ report.Store       = (*FileStore)(nil)
	_ report.FileStore   = (*FileStore)(nil)
	_ audit.Store        = (*FileStore)(nil)
	_ quota.LimitStore   = (*FileStore)(nil)
	_ ns.Store           = (*FileStore)(nil)
)
