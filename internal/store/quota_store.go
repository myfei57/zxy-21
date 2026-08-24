package store

import (
	"encoding/json"

	"lims/internal/quota"
)

// SaveLimit writes one quota row.
func (fs *FileStore) SaveLimit(limit quota.Limit) error {
	return fs.writeJSON("quota", limit.InstrumentID, limit)
}

// LoadLimit reads one quota row.
func (fs *FileStore) LoadLimit(instrumentID string) (quota.Limit, error) {
	var limit quota.Limit
	if err := fs.readJSON("quota", instrumentID, &limit); err != nil {
		return limit, err
	}
	return limit, nil
}

// ListLimits reads every quota row.
func (fs *FileStore) ListLimits() ([]quota.Limit, error) {
	rows, err := fs.listAll("quota")
	if err != nil {
		return nil, err
	}
	out := make([]quota.Limit, 0, len(rows))
	for _, row := range rows {
		var limit quota.Limit
		if err := json.Unmarshal(row, &limit); err != nil {
			return nil, err
		}
		out = append(out, limit)
	}
	return out, nil
}
