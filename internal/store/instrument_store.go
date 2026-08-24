package store

import (
	"encoding/json"
	"fmt"

	"lims/internal/instrument"
)

// SaveInstrument writes one instrument entity.
func (fs *FileStore) SaveInstrument(item instrument.Instrument) error {
	return fs.writeJSON("instruments", item.ID, item)
}

// LoadInstrument reads one instrument entity.
func (fs *FileStore) LoadInstrument(id string) (instrument.Instrument, error) {
	var item instrument.Instrument
	if err := fs.readJSON("instruments", id, &item); err != nil {
		return item, err
	}
	return item, nil
}

// ListInstruments reads every instrument entity.
func (fs *FileStore) ListInstruments() ([]instrument.Instrument, error) {
	rows, err := fs.listAll("instruments")
	if err != nil {
		return nil, err
	}
	out := make([]instrument.Instrument, 0, len(rows))
	for _, row := range rows {
		var item instrument.Instrument
		if err := json.Unmarshal(row, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

// SaveDispatch appends one instrument dispatch record.
func (fs *FileStore) SaveDispatch(assignment instrument.InstrumentAssignment) error {
	return fs.writeJSON("instrument-dispatches", assignment.ID, assignment)
}

// HasAssignment reports whether any dispatch record references the task.
func (fs *FileStore) HasAssignment(taskID string) (bool, error) {
	dispatches, err := fs.ListDispatches()
	if err != nil {
		return false, err
	}
	for _, dispatch := range dispatches {
		if dispatch.TaskID == taskID {
			return true, nil
		}
	}
	return false, nil
}

// ListDispatches reads every instrument dispatch record.
func (fs *FileStore) ListDispatches() ([]instrument.InstrumentAssignment, error) {
	rows, err := fs.listAll("instrument-dispatches")
	if err != nil {
		return nil, err
	}
	out := make([]instrument.InstrumentAssignment, 0, len(rows))
	for _, row := range rows {
		var assignment instrument.InstrumentAssignment
		if err := json.Unmarshal(row, &assignment); err != nil {
			return nil, err
		}
		out = append(out, assignment)
	}
	return out, nil
}

// DispatchFileCount returns how many dispatch files exist for one task,
// used by the console to surface duplicate dispatches.
func (fs *FileStore) DispatchFileCount(taskID string) (int, error) {
	entries, err := fs.listAll("instrument-dispatches")
	if err != nil {
		return 0, err
	}
	count := 0
	for _, row := range entries {
		var assignment instrument.InstrumentAssignment
		if err := json.Unmarshal(row, &assignment); err != nil {
			return 0, fmt.Errorf("decode dispatch: %w", err)
		}
		if assignment.TaskID == taskID {
			count++
		}
	}
	return count, nil
}
