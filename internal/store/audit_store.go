package store

import (
	"encoding/json"

	"lims/internal/audit"
)

// SaveEvent writes one audit event.
func (fs *FileStore) SaveEvent(event audit.Event) error {
	return fs.writeJSON("audit", event.ID, event)
}

// LoadEvent reads one audit event.
func (fs *FileStore) LoadEvent(id string) (audit.Event, error) {
	var event audit.Event
	if err := fs.readJSON("audit", id, &event); err != nil {
		return event, err
	}
	return event, nil
}

// ListEvents reads every audit event.
func (fs *FileStore) ListEvents() ([]audit.Event, error) {
	rows, err := fs.listAll("audit")
	if err != nil {
		return nil, err
	}
	out := make([]audit.Event, 0, len(rows))
	for _, row := range rows {
		var event audit.Event
		if err := json.Unmarshal(row, &event); err != nil {
			return nil, err
		}
		out = append(out, event)
	}
	return out, nil
}
