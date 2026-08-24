package store

import (
	"encoding/json"

	"lims/internal/sample"
)

// SaveSample writes one sample entity.
func (fs *FileStore) SaveSample(item sample.Sample) error {
	return fs.writeJSON("samples", item.ID, item)
}

// LoadSample reads one sample entity.
func (fs *FileStore) LoadSample(id string) (sample.Sample, error) {
	var item sample.Sample
	if err := fs.readJSON("samples", id, &item); err != nil {
		return item, err
	}
	return item, nil
}

// ListSamples reads every sample entity.
func (fs *FileStore) ListSamples() ([]sample.Sample, error) {
	rows, err := fs.listAll("samples")
	if err != nil {
		return nil, err
	}
	out := make([]sample.Sample, 0, len(rows))
	for _, row := range rows {
		var item sample.Sample
		if err := json.Unmarshal(row, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

// SaveRegistration writes one registration record.
func (fs *FileStore) SaveRegistration(record sample.RegistrationRecord) error {
	return fs.writeJSON("registrations", record.SampleID, record)
}

// LoadRegistration reads the registration record of one sample.
func (fs *FileStore) LoadRegistration(sampleID string) (sample.RegistrationRecord, error) {
	var record sample.RegistrationRecord
	if err := fs.readJSON("registrations", sampleID, &record); err != nil {
		return record, err
	}
	return record, nil
}

// ListRegistrations reads every registration record.
func (fs *FileStore) ListRegistrations() ([]sample.RegistrationRecord, error) {
	rows, err := fs.listAll("registrations")
	if err != nil {
		return nil, err
	}
	out := make([]sample.RegistrationRecord, 0, len(rows))
	for _, row := range rows {
		var record sample.RegistrationRecord
		if err := json.Unmarshal(row, &record); err != nil {
			return nil, err
		}
		out = append(out, record)
	}
	return out, nil
}
