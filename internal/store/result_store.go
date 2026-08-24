package store

import (
	"encoding/json"
	"fmt"

	"lims/internal/errs"
	"lims/internal/result"
)

// SaveResult writes one result entity.
func (fs *FileStore) SaveResult(item result.Result) error {
	return fs.writeJSON("results", item.ID, item)
}

// LoadResult reads one result entity.
func (fs *FileStore) LoadResult(id string) (result.Result, error) {
	var item result.Result
	if err := fs.readJSON("results", id, &item); err != nil {
		return item, err
	}
	return item, nil
}

// LoadResultStatus reads only the status field of one result.
func (fs *FileStore) LoadResultStatus(resultID string) (string, error) {
	item, err := fs.LoadResult(resultID)
	if err != nil {
		return "", err
	}
	return string(item.Status), nil
}

// ListResults reads every result entity.
func (fs *FileStore) ListResults() ([]result.Result, error) {
	rows, err := fs.listAll("results")
	if err != nil {
		return nil, err
	}
	out := make([]result.Result, 0, len(rows))
	for _, row := range rows {
		var item result.Result
		if err := json.Unmarshal(row, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

// SaveQCRecord writes the QC record of one result.
func (fs *FileStore) SaveQCRecord(record result.QCRecord) error {
	return fs.writeJSON("qc", record.ResultID, record)
}

// LoadQCRecord reads the QC record of one result.
func (fs *FileStore) LoadQCRecord(resultID string) (result.QCRecord, error) {
	var record result.QCRecord
	if err := fs.readJSON("qc", resultID, &record); err != nil {
		return record, err
	}
	return record, nil
}

// ListQCRecords reads every QC record.
func (fs *FileStore) ListQCRecords() ([]result.QCRecord, error) {
	rows, err := fs.listAll("qc")
	if err != nil {
		return nil, err
	}
	out := make([]result.QCRecord, 0, len(rows))
	for _, row := range rows {
		var record result.QCRecord
		if err := json.Unmarshal(row, &record); err != nil {
			return nil, err
		}
		out = append(out, record)
	}
	return out, nil
}

// SaveAck writes one station acknowledgement.
func (fs *FileStore) SaveAck(ack result.AckRecord) error {
	name := ack.SampleID + "-" + ack.Station
	return fs.writeJSON("acks", name, ack)
}

// LoadAck reads the latest acknowledgement of one sample.
func (fs *FileStore) LoadAck(sampleID string) (result.AckRecord, error) {
	rows, err := fs.listAll("acks")
	if err != nil {
		return result.AckRecord{}, err
	}
	var latest result.AckRecord
	found := false
	for _, row := range rows {
		var ack result.AckRecord
		if err := json.Unmarshal(row, &ack); err != nil {
			return result.AckRecord{}, err
		}
		if ack.SampleID == sampleID && (!found || ack.ReceivedAt.After(latest.ReceivedAt)) {
			latest = ack
			found = true
		}
	}
	if !found {
		return result.AckRecord{}, fmt.Errorf("%w: ack for sample %s", errs.ErrNotFound, sampleID)
	}
	return latest, nil
}

// ListAcks reads every station acknowledgement.
func (fs *FileStore) ListAcks() ([]result.AckRecord, error) {
	rows, err := fs.listAll("acks")
	if err != nil {
		return nil, err
	}
	out := make([]result.AckRecord, 0, len(rows))
	for _, row := range rows {
		var ack result.AckRecord
		if err := json.Unmarshal(row, &ack); err != nil {
			return nil, err
		}
		out = append(out, ack)
	}
	return out, nil
}
