package report

import "fmt"

// Archive moves a signed report into the archived state after confirming
// that its current file is still present.
func Archive(store Store, files FileStore, current Report) (Report, error) {
	if current.Status != StatusSigned {
		return current, fmt.Errorf("only signed reports can be archived, current status %s", current.Status)
	}
	exists, err := FileExists(files, current.ContentPath)
	if err != nil {
		return current, fmt.Errorf("check report file: %w", err)
	}
	if !exists {
		return current, fmt.Errorf("report file missing: %s", current.ContentPath)
	}
	current.Status = StatusArchived
	if err := store.SaveReport(current); err != nil {
		return current, fmt.Errorf("persist archived report: %w", err)
	}
	return current, nil
}
