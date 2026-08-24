package report

import "fmt"

// Revise writes the next report file durably before advancing the version.
func Revise(store Store, files FileStore, current Report, content []byte) (Report, error) {
	next := NextVersion(current)
	current.Version = next
	if err := store.SaveReport(current); err != nil {
		return current, fmt.Errorf("persist revised report: %w", err)
	}
	path, err := WriteContent(files, current.ID, next, content)
	if err != nil {
		return current, err
	}
	current.ContentPath = path
	return current, nil
}
