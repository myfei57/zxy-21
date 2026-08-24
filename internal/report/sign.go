package report

import (
	"fmt"
	"time"
)

// Sign writes the report file durably before marking the report signed.
func Sign(store Store, files FileStore, current Report, content []byte) (Report, error) {
	if !IsDraft(current) {
		return current, fmt.Errorf("only draft reports can be signed, current status %s", current.Status)
	}
	path, err := WriteContent(files, current.ID, current.Version, content)
	if err != nil {
		return current, err
	}
	current.Status = StatusSigned
	current.ContentPath = path
	signedAt := time.Now().UTC()
	current.SignedAt = &signedAt
	if err := store.SaveReport(current); err != nil {
		return current, fmt.Errorf("persist signed report: %w", err)
	}
	return current, nil
}
