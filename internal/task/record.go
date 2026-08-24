package task

import (
	"fmt"
	"time"
)

// AssignmentRecord is the durable proof that a task was dispatched to an instrument.
type AssignmentRecord struct {
	TaskID       string    `json:"task_id"`
	SampleID     string    `json:"sample_id"`
	InstrumentID string    `json:"instrument_id"`
	AssignedAt   time.Time `json:"assigned_at"`
	Operator     string    `json:"operator"`
}

// NewAssignmentRecord builds the assignment evidence for one task.
func NewAssignmentRecord(taskID, sampleID, instrumentID, operator string, at time.Time) AssignmentRecord {
	return AssignmentRecord{
		TaskID:       taskID,
		SampleID:     sampleID,
		InstrumentID: instrumentID,
		AssignedAt:   at,
		Operator:     operator,
	}
}

// SaveAssignmentRecord persists the assignment record and surfaces write errors.
func SaveAssignmentRecord(store Store, record AssignmentRecord) error {
	if err := store.SaveAssignment(record); err != nil {
		return fmt.Errorf("persist assignment record: %w", err)
	}
	return nil
}
