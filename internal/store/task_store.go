package store

import (
	"encoding/json"

	"lims/internal/task"
)

// SaveTask writes one task entity.
func (fs *FileStore) SaveTask(item task.Task) error {
	return fs.writeJSON("tasks", item.ID, item)
}

// LoadTask reads one task entity.
func (fs *FileStore) LoadTask(id string) (task.Task, error) {
	var item task.Task
	if err := fs.readJSON("tasks", id, &item); err != nil {
		return item, err
	}
	return item, nil
}

// ListTasks reads every task entity.
func (fs *FileStore) ListTasks() ([]task.Task, error) {
	rows, err := fs.listAll("tasks")
	if err != nil {
		return nil, err
	}
	out := make([]task.Task, 0, len(rows))
	for _, row := range rows {
		var item task.Task
		if err := json.Unmarshal(row, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

// SaveAssignment writes the assignment record of one task.
func (fs *FileStore) SaveAssignment(record task.AssignmentRecord) error {
	return fs.writeJSON("assignments", record.TaskID, record)
}

// LoadAssignment reads the assignment record of one task.
func (fs *FileStore) LoadAssignment(taskID string) (task.AssignmentRecord, error) {
	var record task.AssignmentRecord
	if err := fs.readJSON("assignments", taskID, &record); err != nil {
		return record, err
	}
	return record, nil
}

// ListAssignments reads every task assignment record.
func (fs *FileStore) ListAssignments() ([]task.AssignmentRecord, error) {
	rows, err := fs.listAll("assignments")
	if err != nil {
		return nil, err
	}
	out := make([]task.AssignmentRecord, 0, len(rows))
	for _, row := range rows {
		var record task.AssignmentRecord
		if err := json.Unmarshal(row, &record); err != nil {
			return nil, err
		}
		out = append(out, record)
	}
	return out, nil
}
