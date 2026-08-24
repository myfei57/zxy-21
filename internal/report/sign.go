package report

import (
	"fmt"
	"time"
)

// Sign writes the report file durably before marking the report signed, so a
// report whose status is "signed" always has downloadable content.
func Sign(store Store, files FileStore, current Report, content []byte) (Report, error) {
	if !IsDraft(current) {
		return current, fmt.Errorf("only draft reports can be signed, current status %s", current.Status)
	}
	// Persist the file first. If this fails the report stays a draft and no
	// "signed" state is ever committed without a backing file.
	path, err := WriteContent(files, current.ID, current.Version, content)
	if err != nil {
		return current, err
	}
	signedAt := time.Now().UTC()
	current.Status = StatusSigned
	current.SignedAt = &signedAt
	current.ContentPath = path
	if err := store.SaveReport(current); err != nil {
		return current, fmt.Errorf("persist signed report: %w", err)
	}
	return current, nil
}
