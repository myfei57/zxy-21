package store

import (
	"encoding/json"

	"lims/internal/report"
)

// SaveReport writes one report entity.
func (fs *FileStore) SaveReport(item report.Report) error {
	return fs.writeJSON("reports", item.ID, item)
}

// LoadReport reads one report entity.
func (fs *FileStore) LoadReport(id string) (report.Report, error) {
	var item report.Report
	if err := fs.readJSON("reports", id, &item); err != nil {
		return item, err
	}
	return item, nil
}

// ListReports reads every report entity.
func (fs *FileStore) ListReports() ([]report.Report, error) {
	rows, err := fs.listAll("reports")
	if err != nil {
		return nil, err
	}
	out := make([]report.Report, 0, len(rows))
	for _, row := range rows {
		var item report.Report
		if err := json.Unmarshal(row, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}
